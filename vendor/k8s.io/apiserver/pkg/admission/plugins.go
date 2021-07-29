/*
Copyright 2014 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package admission

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"reflect"
	"sort"
	"strings"
	"sync"

	"k8s.io/klog/v2"
)

// Factory is a function that returns an Interface for admission decisions.
// The config parameter provides an io.Reader handler to the factory in
// order to load specific configurations. If no configuration is provided
// the parameter is nil.
// 其中Factory为准入控制器实现的接口定义，它接收准入控制器的config配置信息，
// 通过--admission-control-config-file参数指定准入控制器的配置文件，返回准入控制器的插件实现。
// Plugins 数据结构提供了Register方法，为外部提供了准入控制器的注册方法
type Factory func(config io.Reader) (Interface, error)

// kube-apiserver在启动时注册所有准入控制器，准入控制器通过Plugins数据结构统一注册、存放、管理所有的准入控制器。
type Plugins struct {
	// 用于保护registry字段的并发一致性。
	lock sync.Mutex
	// 以键值对形式存放插件，
	// key为准入控制器的名称，例如 AlwaysPullImages、LimitRanger等;
	// value 为对应的准入控制器名称的代码实现。
	registry map[string]Factory
}

func NewPlugins() *Plugins {
	return &Plugins{}
}

// All registered admission options.
var (
	// PluginEnabledFn checks whether a plugin is enabled.  By default, if you ask about it, it's enabled.
	PluginEnabledFn = func(name string, config io.Reader) bool {
		return true
	}
)

// PluginEnabledFunc is a function type that can provide an external check on whether an admission plugin may be enabled
type PluginEnabledFunc func(name string, config io.Reader) bool

// Registered enumerates the names of all registered plugins.
func (ps *Plugins) Registered() []string {
	ps.lock.Lock()
	defer ps.lock.Unlock()
	keys := []string{}
	for k := range ps.registry {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// Register registers a plugin Factory by name. This
// is expected to happen during app startup.
func (ps *Plugins) Register(name string, plugin Factory) {
	ps.lock.Lock()
	defer ps.lock.Unlock()
	if ps.registry != nil {
		_, found := ps.registry[name]
		if found {
			klog.Fatalf("Admission plugin %q was registered twice", name)
		}
	} else {
		ps.registry = map[string]Factory{}
	}

	klog.V(1).InfoS("Registered admission plugin", "plugin", name)
	ps.registry[name] = plugin
}

// getPlugin creates an instance of the named plugin.  It returns `false` if the
// the name is not known. The error is returned only when the named provider was
// known but failed to initialize.  The config parameter specifies the io.Reader
// handler of the configuration file for the cloud provider, or nil for no configuration.
func (ps *Plugins) getPlugin(name string, config io.Reader) (Interface, bool, error) {
	ps.lock.Lock()
	defer ps.lock.Unlock()
	f, found := ps.registry[name]
	if !found {
		return nil, false, nil
	}

	config1, config2, err := splitStream(config)
	if err != nil {
		return nil, true, err
	}
	if !PluginEnabledFn(name, config1) {
		return nil, true, nil
	}

	ret, err := f(config2)
	return ret, true, err
}

// splitStream reads the stream bytes and constructs two copies of it.
func splitStream(config io.Reader) (io.Reader, io.Reader, error) {
	if config == nil || reflect.ValueOf(config).IsNil() {
		return nil, nil, nil
	}

	configBytes, err := ioutil.ReadAll(config)
	if err != nil {
		return nil, nil, err
	}

	return bytes.NewBuffer(configBytes), bytes.NewBuffer(configBytes), nil
}

// NewFromPlugins returns an admission.Interface that will enforce admission control decisions of all
// the given plugins.
func (ps *Plugins) NewFromPlugins(pluginNames []string, configProvider ConfigProvider, pluginInitializer PluginInitializer, decorator Decorator) (Interface, error) {
	handlers := []Interface{}
	mutationPlugins := []string{}
	validationPlugins := []string{}
	for _, pluginName := range pluginNames {
		pluginConfig, err := configProvider.ConfigFor(pluginName)
		if err != nil {
			return nil, err
		}

		plugin, err := ps.InitPlugin(pluginName, pluginConfig, pluginInitializer)
		if err != nil {
			return nil, err
		}
		if plugin != nil {
			if decorator != nil {
				handlers = append(handlers, decorator.Decorate(plugin, pluginName))
			} else {
				handlers = append(handlers, plugin)
			}

			if _, ok := plugin.(MutationInterface); ok {
				mutationPlugins = append(mutationPlugins, pluginName)
			}
			if _, ok := plugin.(ValidationInterface); ok {
				validationPlugins = append(validationPlugins, pluginName)
			}
		}
	}
	if len(mutationPlugins) != 0 {
		klog.Infof("Loaded %d mutating admission controller(s) successfully in the following order: %s.", len(mutationPlugins), strings.Join(mutationPlugins, ","))
	}
	if len(validationPlugins) != 0 {
		klog.Infof("Loaded %d validating admission controller(s) successfully in the following order: %s.", len(validationPlugins), strings.Join(validationPlugins, ","))
	}
	return newReinvocationHandler(chainAdmissionHandler(handlers)), nil
}

// InitPlugin creates an instance of the named interface.
func (ps *Plugins) InitPlugin(name string, config io.Reader, pluginInitializer PluginInitializer) (Interface, error) {
	if name == "" {
		klog.Info("No admission plugin specified.")
		return nil, nil
	}

	plugin, found, err := ps.getPlugin(name, config)
	if err != nil {
		return nil, fmt.Errorf("couldn't init admission plugin %q: %v", name, err)
	}
	if !found {
		return nil, fmt.Errorf("unknown admission plugin: %s", name)
	}

	pluginInitializer.Initialize(plugin)
	// ensure that plugins have been properly initialized
	if err := ValidateInitialization(plugin); err != nil {
		return nil, fmt.Errorf("failed to initialize admission plugin %q: %v", name, err)
	}

	return plugin, nil
}

// ValidateInitialization will call the InitializationValidate function in each plugin if they implement
// the InitializationValidator interface.
func ValidateInitialization(plugin Interface) error {
	if validater, ok := plugin.(InitializationValidator); ok {
		err := validater.ValidateInitialization()
		if err != nil {
			return err
		}
	}
	return nil
}

type PluginInitializers []PluginInitializer

func (pp PluginInitializers) Initialize(plugin Interface) {
	for _, p := range pp {
		p.Initialize(plugin)
	}
}
