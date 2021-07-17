// Copyright (c) 2020 tickstep.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package aliyunpan

import (
	"github.com/tickstep/library-go/requester"
)

const (
	// PathSeparator 路径分隔符
	PathSeparator = "/"
)

var (
)

type (
	PanClient struct {
		client     *requester.HTTPClient // http 客户端
		webToken WebLoginToken
		appToken AppLoginToken
	}
)


func NewPanClient(webToken WebLoginToken, appToken AppLoginToken) *PanClient {
	client := requester.NewHTTPClient()

	return &PanClient{
		client: client,
		webToken: webToken,
		appToken: appToken,
	}
}
