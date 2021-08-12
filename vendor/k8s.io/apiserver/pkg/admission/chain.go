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

import "context"

// chainAdmissionHandler is an instance of admission.NamedHandler that performs admission control using
// a chain of admission handlers
// kube-apisever中的所有已启用的准入控制器(Admit 方法及Validate 方法)由chainAdmissionHandler []Interface数据结构管理
// 假设kube-apisever开启了AlwaysPullImages和PodNodeSelector准入控制器，
// 当客户端发送请求给kube-apisever时，该请求会进入Admission Controller Handler函数(处理准入控制器相关的Handler 函数)。
// 在Admission Controller Handler中，会遍历已启用的准入控制器列表，按顺序尝试执行每个准入控制器，执行所有的变更操作。
type chainAdmissionHandler []Interface

// NewChainHandler creates a new chain handler from an array of handlers. Used for testing.
func NewChainHandler(handlers ...Interface) chainAdmissionHandler {
	return chainAdmissionHandler(handlers)
}

// Admit performs an admission control check using a chain of handlers, and returns immediately on first error
// Admit函数会遍历已启用的准入控制器列表，并执行变更操作的准入控制器(即拥有Admit方法的准入控制器)。
func (admissionHandler chainAdmissionHandler) Admit(ctx context.Context, a Attributes, o ObjectInterfaces) error {
	for _, handler := range admissionHandler {
		if !handler.Handles(a.GetOperation()) {
			continue
		}
		if mutator, ok := handler.(MutationInterface); ok {
			err := mutator.Admit(ctx, a, o)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Validate performs an admission control check using a chain of handlers, and returns immediately on first error
// 以同样的方式执行Validate函数，Validate函数会遍历已启用的准入控制器列表，并执行验证操作的准入控制器(即拥有Validate方法的准入控制器)
func (admissionHandler chainAdmissionHandler) Validate(ctx context.Context, a Attributes, o ObjectInterfaces) error {
	for _, handler := range admissionHandler {
		if !handler.Handles(a.GetOperation()) {
			continue
		}
		if validator, ok := handler.(ValidationInterface); ok {
			err := validator.Validate(ctx, a, o)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Handles will return true if any of the handlers handles the given operation
func (admissionHandler chainAdmissionHandler) Handles(operation Operation) bool {
	for _, handler := range admissionHandler {
		if handler.Handles(operation) {
			return true
		}
	}
	return false
}
