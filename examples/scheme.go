package main

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func main() {
	// KnownType external
	coreGV := schema.GroupVersion{"", "v1"}
	extensionsGV := schema.GroupVersion{"extensions", "v1beta1"}

	// KnownType internal
	coreInternalGV := schema.GroupVersion{"", runtime.APIVersionInternal}

	// UnversionedType
	Unversioned := schema.GroupVersion{"", "v1"}

	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(coreGV, &corev1.Pod{})
	scheme.AddKnownTypes(extensionsGV, &appsv1.DaemonSet{})
	scheme.AddKnownTypes(coreInternalGV, &corev1.Pod{})
	scheme.AddKnownTypes(Unversioned, &metav1.Status{})
}
