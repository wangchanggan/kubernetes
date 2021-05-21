package main

import (
	"fmt"
	"k8s.io/api/apps/v1"
	"k8s.io/api/apps/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/kubernetes/pkg/apis/apps"
)

func main() {
	// 实例化一个空的Scheme资源注册表，将v1beta1资源版本、v1资源版本及内部版本（__internal）的Deployment资源注册到Scheme资源注册表中。
	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(v1beta1.SchemeGroupVersion, &v1beta1.Deployment{})
	scheme.AddKnownTypes(v1.SchemeGroupVersion, &v1.Deployment{})
	scheme.AddKnownTypes(apps.SchemeGroupVersion, &v1.Deployment{})
	metav1.AddToGroupVersion(scheme, v1beta1.SchemeGroupVersion)
	metav1.AddToGroupVersion(scheme, v1.SchemeGroupVersion)

	// 实例化v1beta1Deployment资源对象，
	v1beta1Deployment := &v1beta1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1betal",
		},
	}

	// 通过scheme.ConvertToVersion将其转换为目标资源版本(即__internal版本)，得到objInternal资源对象
	// vlbetal -> __internal
	objInternal, err := scheme.ConvertToVersion(v1beta1Deployment, apps.SchemeGroupVersion)
	if err != nil {
		panic(err)
	}

	// objInternal资源对象的GVK输出为“/, Kind=”，当资源对象的GVK输出为“/, Kind=”时，同样认为它是内部版本的资源对象
	fmt.Println("GVK: ", objInternal.GetObjectKind().GroupVersionKind().String())
	// output:
	// GVK: /, Kind=

	// 将objInternal资源对象通过scheme.ConvertToVersion转换为目标资源版本（即v1资源服本），得到objV1资源对象
	// __internal -> v1
	objV1, err := scheme.ConvertToVersion(objInternal, v1.SchemeGroupVersion)
	if err != nil {
		panic(err)
	}

	//通过断言的方式来验证是否转换成功
	v1Deployment, ok := objV1.(*v1.Deployment)
	if !ok {
		panic("Got wrong type")
	}

	// objV1资源对象的GVK输出为“apps/vl, Kind=Deployment”
	fmt.Println("GVK: ", v1Deployment.GetObjectKind().GroupVersionKind().String())
	// output:
	// GVK: apps/vl, Kind=Deployment
}
