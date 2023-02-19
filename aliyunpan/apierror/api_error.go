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
	"io/ioutil"
	"net/http"
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
	// ApiCodeFeatureTemporaryDisabled 功能维护中
	ApiCodeFeatureTemporaryDisabled ApiCode = 27
	// ApiCodeForbiddenFileInTheRecycleBin 文件已经被删除
	ApiCodeForbiddenFileInTheRecycleBin ApiCode = 28
	// ApiCodeBadGateway 502网关错误，一般代表请求被限流了
	ApiCodeBadGateway ApiCode = 29
	// ApiCodeTooManyRequests 429 Too Many Requests错误，一般代表请求被限流了
	ApiCodeTooManyRequests ApiCode = 30
	// ApiCodeUserDeviceOffline 客户端离线，阿里云盘单账户最多只允许同时登录 10 台设备
	ApiCodeUserDeviceOffline ApiCode = 31
	// ApiCodeDeviceSessionSignatureInvalid 签名过期，需要更新签名密钥
	ApiCodeDeviceSessionSignatureInvalid ApiCode = 32
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
	if string(data) == "Bad Gateway" {
		// 	HTTP/1.1 502 Bad Gateway
		return NewApiError(ApiCodeBadGateway, "网关错误，你的请求可能被临时限流了")
	}
	errResp := &ErrorResp{}
	if err := json.Unmarshal(data, errResp); err == nil {
		if errResp.ErrorCode != "" {
			if "AccessTokenInvalid" == errResp.ErrorCode {
				return NewApiError(ApiCodeAccessTokenInvalid, errResp.GetErrorMsg())
			} else if "NotFound.File" == errResp.ErrorCode || "NotFound.FileId" == errResp.ErrorCode {
				return NewApiError(ApiCodeFileNotFoundCode, errResp.GetErrorMsg())
			} else if "AlreadyExist.File" == errResp.ErrorCode {
				return NewApiError(ApiCodeFileAlreadyExisted, errResp.GetErrorMsg())
			} else if "BadRequest" == errResp.ErrorCode {
				return NewApiError(ApiCodeFailed, errResp.GetErrorMsg())
			} else if "InvalidParameter.RefreshToken" == errResp.ErrorCode {
				return NewApiError(ApiCodeRefreshTokenExpiredCode, errResp.GetErrorMsg())
			} else if "FileShareNotAllowed" == errResp.ErrorCode {
				return NewApiError(ApiCodeFileShareNotAllowed, errResp.GetErrorMsg())
			} else if "InvalidRapidProof" == errResp.ErrorCode {
				return NewApiError(ApiCodeInvalidRapidProof, errResp.GetErrorMsg())
			} else if "NotFound.View" == errResp.ErrorCode {
				return NewApiError(ApiCodeNotFoundView, errResp.GetErrorMsg())
			} else if "BadRequest" == errResp.ErrorCode {
				return NewApiError(ApiCodeBadRequest, errResp.GetErrorMsg())
			} else if "InvalidResource.FileTypeFolder" == errResp.ErrorCode {
				return NewApiError(ApiCodeInvalidResource, errResp.GetErrorMsg())
			} else if "NotFound.VideoPreviewInfo" == errResp.ErrorCode {
				return NewApiError(ApiCodeVideoPreviewInfoNotFound, errResp.GetErrorMsg())
			} else if "FeatureTemporaryDisabled" == errResp.ErrorCode {
				return NewApiError(ApiCodeFeatureTemporaryDisabled, errResp.GetErrorMsg())
			} else if "ForbiddenFileInTheRecycleBin" == errResp.ErrorCode {
				return NewApiError(ApiCodeForbiddenFileInTheRecycleBin, errResp.GetErrorMsg())
			} else if "UserDeviceOffline" == errResp.ErrorCode {
				return NewApiError(ApiCodeUserDeviceOffline, "你账号已超出最大登录设备数量，请先下线一台设备，然后重启本应用，才可以继续使用")
			} else if "DeviceSessionSignatureInvalid" == errResp.ErrorCode {
				return NewApiError(ApiCodeDeviceSessionSignatureInvalid, "签名过期，需要更新签名密钥")
			}
			return NewFailedApiError(errResp.GetErrorMsg())
		}
	}
	return nil
}

// ParseCommonResponseApiError 解析阿里云盘错误，如果没有错误则返回nil
func ParseCommonResponseApiError(resp *http.Response) ([]byte, *ApiError) {
	if resp == nil {
		return nil, nil
	}

	switch resp.StatusCode {
	case 502:
		return nil, NewApiError(ApiCodeBadGateway, "网关错误，你的请求可能被临时限流了")
	case 429:
		return nil, NewApiError(ApiCodeTooManyRequests, "太频繁请求错误，你的请求可能被临时限流了")
	}
	data, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		return nil, NewFailedApiError(e.Error())
	}
	return data, ParseCommonApiError(data)
}
