package main

import (
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main() {
	// 在RBAC Example代码示例中，通过rbacv1 .Role创建了一个名为PodReader的角色，该角色对资源v1/pods 拥有get、 list、 watch 操作权限。
	roles := &rbacv1.Role{
		ObjectMeta: metav1.ObjectMeta{Namespace: "default", Name: "PodReader"},
		Rules: []rbacv1.PolicyRule{
			{
				Verbs:     []string{"get", "list", "watch"},
				APIGroups: []string{"v1"},
				Resources: []string{"pods"},
			},
		},
	}

	// 通过rbacv1.RoleBinding将角色与用户绑定，绑定的用户为Derek并只被授予default命名空间的权限。
	roleBindings := &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{Namespace: "default"},
		Subjects: []rbacv1.Subject{
			{APIGroup: rbacv1.GroupName, Kind: rbacv1.UserKind, Name: "Derek"},
		},
		RoleRef: rbacv1.RoleRef{APIGroup: rbacv1.GroupName, Kind: "Role", Name: "PodReader"},
	}

	// 完成上述操作后，Derek 用户对default命名空间下的v1/pods资源拥有了get、list、watch操作权限，但Derek用户并没有其他命名空间下任何资源的操作权限。
}
