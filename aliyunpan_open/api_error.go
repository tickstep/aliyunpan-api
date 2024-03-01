package aliyunpan_open

import (
	"github.com/tickstep/aliyunpan-api/aliyunpan/apierror"
	"github.com/tickstep/aliyunpan-api/aliyunpan_open/openapi"
	"github.com/tickstep/library-go/logger"
	"time"
)

type (
	ApiErrorHandleResp struct {
		// NeedRetry 是否需要重试
		NeedRetry bool
		// ApiErr 错误
		ApiErr *apierror.ApiError
	}
)

func NewApiErrorHandleResp(needRetry bool, apiErr *apierror.ApiError) *ApiErrorHandleResp {
	return &ApiErrorHandleResp{
		NeedRetry: needRetry,
		ApiErr:    apiErr,
	}
}

// ParseAliApiError 解析阿里接口返回错误，封装成本地的统一实体
func (p *OpenPanClient) ParseAliApiError(respErr *openapi.AliApiErrResult) *apierror.ApiError {
	if respErr == nil {
		return nil
	}

	switch respErr.HttpStatusCode {
	case 200:
		return apierror.NewFailedApiError(respErr.Message)
	case 400:
		if respErr.Code == "NotFound.File" {
			return apierror.NewApiError(apierror.ApiCodeFileNotFoundCode, respErr.Message)
		}
	case 401:
		if respErr.Code == "AccessTokenExpired" {
			return apierror.NewApiError(apierror.ApiCodeTokenExpiredCode, respErr.Message)
		} else if respErr.Code == "RefreshTokenExpired" {
			return apierror.NewApiError(apierror.ApiCodeRefreshTokenExpiredCode, respErr.Message)
		}
	case 403:
		if respErr.Code == "PermissionDenied" {
			return apierror.NewApiError(apierror.ApiCodePermissionDenied, respErr.Message)
		} else if respErr.Code == "UserNotAllowedAccessDrive" {
			return apierror.NewApiError(apierror.ApiCodeUserNotAllowedAccessDrive, respErr.Message)
		}
	case 404:
		if respErr.Code == "NotFound.FileId" {
			return apierror.NewApiError(apierror.ApiCodeFileNotFoundCode, respErr.Message)
		}
	case 409:
		if respErr.Code == "TooManyRequests" {
			return apierror.NewApiError(apierror.ApiCodeTooManyRequests, respErr.Message)
		}
	}
	return apierror.NewFailedApiError(respErr.Message)
}

// HandleAliApiError 处理公共错误
func (p *OpenPanClient) HandleAliApiError(respErr *openapi.AliApiErrResult, retryTime *int) *ApiErrorHandleResp {
	// handle error, retry, token refresh
	myApiErr := p.ParseAliApiError(respErr)
	if myApiErr.Code == apierror.ApiCodeAccessTokenInvalid {
		// get new access token
		time.Sleep(time.Duration(1) * time.Second)
		if tokenErr := p.RefreshNewAccessToken(); tokenErr != nil {
			logger.Verboseln("get new access token from server error: ", tokenErr)
			time.Sleep(time.Duration(2) * time.Second)
		}
		// retry check
		if *retryTime < ApiRetryMaxTimes {
			*retryTime++
			return NewApiErrorHandleResp(true, myApiErr)
		}
	} else if myApiErr.Code == apierror.ApiCodeTooManyRequests {
		// sleep 3s
		time.Sleep(time.Duration(int64(*retryTime+1)*2) * time.Second)
		// retry check
		if *retryTime < ApiRetryMaxTimes {
			*retryTime++
			return NewApiErrorHandleResp(true, myApiErr)
		}
	}
	return NewApiErrorHandleResp(false, myApiErr)
}
