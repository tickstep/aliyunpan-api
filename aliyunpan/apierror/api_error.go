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

const (
	// 成功
	ApiCodeOk ApiCode = 0
	// 成功
	ApiCodeNeedCaptchaCode ApiCode = 10
	// 会话/Token已过期
	ApiCodeTokenExpiredCode ApiCode = 11
	// 文件不存在
	ApiCodeFileNotFoundCode ApiCode = 12
	// 上传文件失败
	ApiCodeUploadFileStatusVerifyFailed = 13
	// 上传文件数据偏移值校验失败
	ApiCodeUploadOffsetVerifyFailed = 14
	// 服务器上传文件不存在
	ApiCodeUploadFileNotFound = 15
	// 文件已存在
	ApiCodeFileAlreadyExisted = 16
	// 上传达到日数量上限
	UserDayFlowOverLimited = 17
	// 失败
	ApiCodeFailed ApiCode = 999
)

type ApiCode int

type ApiError struct {
	Code ApiCode
	Err string
}

func NewApiError(code ApiCode, err string) *ApiError {
	return &ApiError {
		code,
		err,
	}
}

func NewApiErrorWithError(err error) *ApiError {
	if err == nil {
		return NewApiError(ApiCodeOk, "")
	} else {
		return NewApiError(ApiCodeFailed, err.Error())
	}
}

func NewOkApiError() *ApiError {
	return NewApiError(ApiCodeOk, "")
}

func NewFailedApiError(err string) *ApiError {
	return NewApiError(ApiCodeFailed, err)
}

func (a *ApiError) SetErr(code ApiCode, err string) {
	a.Code = code
	a.Err = err
}

func (a *ApiError) Error() string {
	return a.Err
}

func (a *ApiError) ErrCode() ApiCode {
	return a.Code
}
