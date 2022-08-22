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

import (
	"encoding/json"
	"fmt"
	"strings"
)

const (
	/* ------------------------------- 默认错误码 -------------------------------*/

	// ApiCodeOk 成功
	ApiCodeOk ApiCode = 0
	// ApiCodeFailed 失败
	ApiCodeFailed ApiCode = 999

	/* ------------------------------- 系统错误码(800-899) -------------------------------*/

	// ApiCodeNetError 网络错误
	ApiCodeNetError ApiCode = 800

	/* ------------------------------- 阿里云盘错误码(10-799) -------------------------------*/

	// ApiCodeNeedCaptchaCode 验证码
	ApiCodeNeedCaptchaCode ApiCode = 10
	// ApiCodeTokenExpiredCode 会话/Token已过期
	ApiCodeTokenExpiredCode ApiCode = 11
	// ApiCodeFileNotFoundCode 文件不存在 NotFound.File / NotFound.FileId
	ApiCodeFileNotFoundCode ApiCode = 12
	// ApiCodeUploadFileStatusVerifyFailed 上传文件失败
	ApiCodeUploadFileStatusVerifyFailed = 13
	// ApiCodeUploadOffsetVerifyFailed 上传文件数据偏移值校验失败
	ApiCodeUploadOffsetVerifyFailed = 14
	// ApiCodeUploadFileNotFound 服务器上传文件不存在
	ApiCodeUploadFileNotFound = 15
	// ApiCodeFileAlreadyExisted 文件已存在 AlreadyExist.File
	ApiCodeFileAlreadyExisted = 16
	// ApiCodeUserDayFlowOverLimited 上传达到日数量上限
	ApiCodeUserDayFlowOverLimited = 17
	// ApiCodeAccessTokenInvalid Token无效或者已过期 AccessTokenInvalid
	ApiCodeAccessTokenInvalid = 18
	// ApiCodeForbidden 被禁止 Forbidden
	ApiCodeForbidden = 19
	// ApiCodeRefreshTokenExpiredCode RefreshToken已过期
	ApiCodeRefreshTokenExpiredCode ApiCode = 20
	// ApiCodeFileShareNotAllowed 文件不允许分享
	ApiCodeFileShareNotAllowed ApiCode = 21
	// ApiCodeInvalidRapidProof 文件上传水印码错误
	ApiCodeInvalidRapidProof ApiCode = 22
	// ApiCodeNotFoundView 资源不存在
	ApiCodeNotFoundView ApiCode = 23
	// ApiCodeBadRequest 请求非法
	ApiCodeBadRequest ApiCode = 24
	// ApiCodeInvalidResource 请求无效资源
	ApiCodeInvalidResource ApiCode = 25
	// ApiCodeVideoPreviewInfoNotFound 视频预览信息不存在
	ApiCodeVideoPreviewInfoNotFound ApiCode = 26
)

type ApiCode int

type ApiError struct {
	Code ApiCode
	Err  string
}

func NewApiError(code ApiCode, err string) *ApiError {
	return &ApiError{
		code,
		err,
	}
}

func NewApiErrorWithError(err error) *ApiError {
	if err == nil {
		return NewApiError(ApiCodeOk, "")
	} else {
		if IsNetErr(err) {
			return NewApiError(ApiCodeNetError, err.Error())
		}
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

func (a *ApiError) String() string {
	sb := &strings.Builder{}
	fmt.Fprintf(sb, "Code=%d, Err=%s", a.Code, a.Err)
	return sb.String()
}

// ParseCommonApiError 解析阿里云盘错误，如果没有错误则返回nil
func ParseCommonApiError(data []byte) *ApiError {
	errResp := &ErrorResp{}
	if err := json.Unmarshal(data, errResp); err == nil {
		if errResp.ErrorCode != "" {
			if "AccessTokenInvalid" == errResp.ErrorCode {
				return NewApiError(ApiCodeAccessTokenInvalid, errResp.ErrorMsg)
			} else if "NotFound.File" == errResp.ErrorCode || "NotFound.FileId" == errResp.ErrorCode {
				return NewApiError(ApiCodeFileNotFoundCode, errResp.ErrorMsg)
			} else if "AlreadyExist.File" == errResp.ErrorCode {
				return NewApiError(ApiCodeFileAlreadyExisted, errResp.ErrorMsg)
			} else if "BadRequest" == errResp.ErrorCode {
				return NewApiError(ApiCodeFailed, errResp.ErrorMsg)
			} else if "InvalidParameter.RefreshToken" == errResp.ErrorCode {
				return NewApiError(ApiCodeRefreshTokenExpiredCode, errResp.ErrorMsg)
			} else if "FileShareNotAllowed" == errResp.ErrorCode {
				return NewApiError(ApiCodeFileShareNotAllowed, errResp.ErrorMsg)
			} else if "InvalidRapidProof" == errResp.ErrorCode {
				return NewApiError(ApiCodeInvalidRapidProof, errResp.ErrorMsg)
			} else if "NotFound.View" == errResp.ErrorCode {
				return NewApiError(ApiCodeNotFoundView, errResp.ErrorMsg)
			} else if "BadRequest" == errResp.ErrorCode {
				return NewApiError(ApiCodeBadRequest, errResp.ErrorMsg)
			} else if "InvalidResource.FileTypeFolder" == errResp.ErrorCode {
				return NewApiError(ApiCodeInvalidResource, errResp.ErrorMsg)
			} else if "NotFound.VideoPreviewInfo" == errResp.ErrorCode {
				return NewApiError(ApiCodeVideoPreviewInfoNotFound, errResp.ErrorMsg)
			}
			return NewFailedApiError(errResp.ErrorMsg)
		}
	}
	return nil
}
