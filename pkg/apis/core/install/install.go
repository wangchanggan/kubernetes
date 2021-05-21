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

// Package install installs the v1 monolithic api, making it available as an
// option to all of the API encoding/decoding machinery.
package install

import (
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/kubernetes/pkg/api/legacyscheme"
	"k8s.io/kubernetes/pkg/apis/core"
	"k8s.io/kubernetes/pkg/apis/core/v1"
)

func init() {
	// legacyscheme.Scheme是kube-apiserver组件的全局资源注册表，Kubernetes的所有资源信息都交给资源注册表统一管理。
	Install(legacyscheme.Scheme)
}

// Install registers the API group and adds types to a scheme
func Install(scheme *runtime.Scheme) {
	// 注册core资源组内部版本的资源
	utilruntime.Must(core.AddToScheme(scheme))
	// 注册 core 资源组外部版本的资源
	utilruntime.Must(v1.AddToScheme(scheme))
	// 注册资源组的版本顺序，如有多个资源版本，排在最前面的为资源首选版本。
	utilruntime.Must(scheme.SetVersionPriority(v1.SchemeGroupVersion))
}
