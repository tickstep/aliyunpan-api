package aliyunpan_open

import (
	"github.com/tickstep/aliyunpan-api/aliyunpan"
	"github.com/tickstep/aliyunpan-api/aliyunpan/apierror"
	"github.com/tickstep/aliyunpan-api/aliyunpan/apiutil"
	"github.com/tickstep/aliyunpan-api/aliyunpan_open/openapi"
)

// ShareLinkCreate 创建文件分享
func (p *OpenPanClient) ShareLinkCreate(param aliyunpan.ShareCreateParam) (*aliyunpan.ShareEntity, *apierror.ApiError) {
	retryTime := 0

RetryBegin:
	opParam := &openapi.FileShareCreateParam{
		DriveId:    param.DriveId,
		FileIdList: param.FileIdList,
		Expiration: param.Expiration,
		SharePwd:   param.SharePwd,
	}
	// format time
	if opParam.Expiration != "" {
		opParam.Expiration = apiutil.LocalTime2UtcFormat(param.Expiration)
	}
	if result, err := p.apiClient.FileShareCreate(opParam); err == nil {
		return &aliyunpan.ShareEntity{
			Creator:    result.Creator,
			DriveId:    param.DriveId,
			ShareId:    result.ShareId,
			ShareName:  "",
			SharePwd:   result.SharePwd,
			ShareUrl:   result.ShareUrl,
			FileIdList: nil,
			SaveCount:  0,
			Expiration: apiutil.UtcTime2LocalFormat(result.Expiration),
			UpdatedAt:  apiutil.UtcTime2LocalFormat(result.UpdatedAt),
			CreatedAt:  apiutil.UtcTime2LocalFormat(result.CreatedAt),
			Status:     result.Status,
			FirstFile:  nil,
		}, nil
	} else {
		// handle common error
		if apiErrorHandleResp := p.HandleAliApiError(err, &retryTime); apiErrorHandleResp.NeedRetry {
			goto RetryBegin
		} else {
			return nil, apiErrorHandleResp.ApiErr
		}
	}
}
