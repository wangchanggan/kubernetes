/*
Copyright 2018 The Kubernetes Authors.

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

package rest

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"
	genericapiserver "k8s.io/apiserver/pkg/server"
	serverstorage "k8s.io/apiserver/pkg/server/storage"

	"k8s.io/kube-aggregator/pkg/apis/apiregistration"
	v1 "k8s.io/kube-aggregator/pkg/apis/apiregistration/v1"
	"k8s.io/kube-aggregator/pkg/apis/apiregistration/v1beta1"
	aggregatorscheme "k8s.io/kube-aggregator/pkg/apiserver/scheme"
	apiservicestorage "k8s.io/kube-aggregator/pkg/registry/apiservice/etcd"
)

// NewRESTStorage returns an APIGroupInfo object that will work against apiservice.
func NewRESTStorage(apiResourceConfigSource serverstorage.APIResourceConfigSource, restOptionsGetter generic.RESTOptionsGetter, shouldServeBeta bool) genericapiserver.APIGroupInfo {
	apiGroupInfo := genericapiserver.NewDefaultAPIGroupInfo(apiregistration.GroupName, aggregatorscheme.Scheme, metav1.ParameterCodec, aggregatorscheme.Codecs)

	// AggregatorServer会先判断 apiregistration.k8s.io/v1 资源组/资源版本是否已启用，
	// 如果其已启用，则将该资源组/资源版本下的资源与资源存储对象进行映射，并将其存储至APIGroupInfo对象的VersionedResourcesStorageMap 字段中。
	if shouldServeBeta && apiResourceConfigSource.VersionEnabled(v1beta1.SchemeGroupVersion) {
		storage := map[string]rest.Storage{}
		// 每个资源(包括子资源)都通过类似于NewREST的函数创建资源存储对象(即RESTStorage )。
		// kube-apiserver将RESTStorage封装成HTTP Handler方法，资源存储对象以RESTful的方式运行，一个RESTStorage对象负责一个资源的增、删、改、查操作。
		// 当操作apiservices 资源数据时，通过对应的RESTStorage资源存储对象与genericregistry.Store进行交互。
		apiServiceREST := apiservicestorage.NewREST(aggregatorscheme.Scheme, restOptionsGetter)
		storage["apiservices"] = apiServiceREST
		storage["apiservices/status"] = apiservicestorage.NewStatusREST(aggregatorscheme.Scheme, apiServiceREST)
		apiGroupInfo.VersionedResourcesStorageMap["v1beta1"] = storage
	}

	if apiResourceConfigSource.VersionEnabled(v1.SchemeGroupVersion) {
		storage := map[string]rest.Storage{}
		apiServiceREST := apiservicestorage.NewREST(aggregatorscheme.Scheme, restOptionsGetter)
		storage["apiservices"] = apiServiceREST
		storage["apiservices/status"] = apiservicestorage.NewStatusREST(aggregatorscheme.Scheme, apiServiceREST)
		apiGroupInfo.VersionedResourcesStorageMap["v1"] = storage
	}

	return apiGroupInfo
}
