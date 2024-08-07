package aliyunpan_open

import (
	"github.com/tickstep/aliyunpan-api/aliyunpan"
	"github.com/tickstep/aliyunpan-api/aliyunpan/apierror"
	"github.com/tickstep/aliyunpan-api/aliyunpan/apiutil"
)

// GetUserInfo 获取用户信息
func (p *OpenPanClient) GetUserInfo() (*aliyunpan.UserInfo, *apierror.ApiError) {
	retryTime := 0
	returnResult := &aliyunpan.UserInfo{
		DomainId:        "",
		FileDriveId:     "",
		SafeBoxDriveId:  "",
		AlbumDriveId:    "",
		ResourceDriveId: "",
		UserId:          "",
		UserName:        "",
		CreatedAt:       "",
		Email:           "",
		Phone:           "",
		Role:            "",
		Status:          "",
		Nickname:        "",
		TotalSize:       0,
		UsedSize:        0,
	}

RetryBegin:
	// user basic info
	if result, err := p.apiClient.UserGetDriveInfo(); err == nil {
		returnResult = &aliyunpan.UserInfo{
			DomainId:            "",
			FileDriveId:         result.BackupDriveId,
			SafeBoxDriveId:      "",
			AlbumDriveId:        "",
			ResourceDriveId:     result.ResourceDriveId,
			UserId:              result.UserId,
			UserName:            "",
			CreatedAt:           "",
			Email:               "",
			Phone:               "",
			Role:                "",
			Status:              "",
			Nickname:            result.Name,
			TotalSize:           0,
			UsedSize:            0,
			ThirdPartyVip:       false,
			ThirdPartyVipExpire: "",
		}
	} else {
		// handle common error
		if apiErrorHandleResp := p.HandleAliApiError(err, &retryTime); apiErrorHandleResp.NeedRetry {
			goto RetryBegin
		} else {
			return nil, apiErrorHandleResp.ApiErr
		}
	}

	// user vip info
	if result, err := p.apiClient.UserGetVipInfo(); err == nil {
		returnResult.ThirdPartyVip = result.ThirdPartyVip
		if result.ThirdPartyVipExpire > 0 {
			returnResult.ThirdPartyVipExpire = apiutil.UnixTime2LocalFormat(result.ThirdPartyVipExpire)
		}
	} else {
		// handle common error
		if apiErrorHandleResp := p.HandleAliApiError(err, &retryTime); apiErrorHandleResp.NeedRetry {
			goto RetryBegin
		} else {
			return nil, apiErrorHandleResp.ApiErr
		}
	}

	// drive spaces
	if result, err := p.apiClient.UserGetSpaceInfo(); err == nil {
		returnResult.TotalSize = uint64(result.TotalSize)
		returnResult.UsedSize = uint64(result.UsedSize)
	} else {
		// handle common error
		if apiErrorHandleResp := p.HandleAliApiError(err, &retryTime); apiErrorHandleResp.NeedRetry {
			goto RetryBegin
		} else {
			return nil, apiErrorHandleResp.ApiErr
		}
	}

	return returnResult, nil
}
