package main

import (
	"fmt"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/kubernetes/pkg/apis/core"
	"reflect"
)

func main() {
	// 实例化Pod资源，得到Pod资源对象。
	pod := &core.Pod{
		TypeMeta: v1.TypeMeta{
			Kind: "Pod",
		},
		ObjectMeta: v1.ObjectMeta{
			Labels: map[string]string{"name": "foo"},
		},
	}

	// 通过runtime.Object将Pod资源对象转换成通用资源对象，得到obj
	obj := runtime.Object(pod)

	// 通过断言的方式，将obj通用资源对象转换成Pod资源对象，得到pdd2
	pod2, ok := obj.(*core.Pod)
	if !ok {
		panic("unexpected")
	}

	// 通过reflect反射来验证转换之前和转换之后的资源对象是否相等
	if reflect.DeepEqual(pod, pod2) {
		fmt.Println("expected")
	} else {
		panic("unexpected")
	}
}
