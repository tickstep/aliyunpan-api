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

package apierror

import "encoding/xml"

// ErrorResp 默认的错误信息
type ErrorResp struct {
	ErrorCode       string `json:"code"`
	ErrorMsg        string `json:"message"`
	ErrorDisplayMsg string `json:"display_message"`
}

type ErrorXmlResp struct {
	XMLName xml.Name `xml:"Error"`
	Code    string   `xml:"Code"`
	Message string   `xml:"Message"`
}

func (e *ErrorResp) GetErrorMsg() string {
	if e.ErrorDisplayMsg != "" {
		return e.ErrorDisplayMsg
	}
	return e.ErrorMsg
}
