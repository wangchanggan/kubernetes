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

package tokentest

import (
	"context"

	"k8s.io/apiserver/pkg/authentication/authenticator"
	"k8s.io/apiserver/pkg/authentication/user"
)

type TokenAuthenticator struct {
	Tokens map[string]*user.DefaultInfo
}

func New() *TokenAuthenticator {
	return &TokenAuthenticator{
		Tokens: make(map[string]*user.DefaultInfo),
	}
}

// Token认证接口定义了AuthenticateToken 方法，该方法接收token字符串。
// 若验证失败，bool 值会为false; 若验证成功，bool 值会为true, 并返回*authenticator.Response,
// *authenticator.Response中携带了身份验证用户的信息，例如Name、UID、Groups、Extra 等信息。
func (a *TokenAuthenticator) AuthenticateToken(ctx context.Context, value string) (*authenticator.Response, bool, error) {
	//在进行Token认证时，a.tokens 中存储了服务端的Token列表，通过a.tokens查询客户端提供的Token
	// 如果查询不到，则认证失败返回false, 反之则认证成功返回true。
	user, ok := a.Tokens[value]
	if !ok {
		return nil, false, nil
	}
	return &authenticator.Response{User: user}, true, nil
}
