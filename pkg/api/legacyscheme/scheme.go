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

package legacyscheme

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

// 在legacyscheme包中，定义了Scheme资源注册表、Codec编解码器及ParameterCodec参数编解码器。
// 它们被定义为全局变量，这些全局变量在kube-apiserver的任何地方都可以被调用，服务于KubeAPIServer。
var (
	// Scheme is the default instance of runtime.Scheme to which types in the Kubernetes API are already registered.
	// NOTE: If you are copying this file to start a new api group, STOP! Copy the
	// extensions group instead. This Scheme is special and should appear ONLY in
	// the api group, unless you really know what you're doing.
	// TODO(lavalamp): make the above error impossible.
	Scheme = runtime.NewScheme()

	// Codecs provides access to encoding and decoding for the scheme
	Codecs = serializer.NewCodecFactory(Scheme)

	// ParameterCodec handles versioning of objects that are converted to query parameters.
	ParameterCodec = runtime.NewParameterCodec(Scheme)
)

// 提示：除将KubeAPIServer ( API核心服务)注册至legacyscheme.Scheme 资源注册表以外，还需要了解APIExtensionsServer和AggregatorServer资源注册过程。
// 将APIExtensionsServer ( API扩展服务)注册至extensionsapiserver.Scheme资源注册表，注册过程定义在vendor/k8s.io/apiextensions-apiserver/pkg/apiserver/apiserver.go中。
// 将AggregatorServer ( API聚合服务)注册至aggregatorscheme Scheme资源注册表，注册过程定义在vendor/k8s.io/kube-aggregator/pkg/apiserver/scheme/scheme.go中。
