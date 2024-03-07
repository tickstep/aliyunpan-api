package aliyunpan_open

import (
	"github.com/tickstep/aliyunpan-api/aliyunpan"
	"github.com/tickstep/aliyunpan-api/aliyunpan/apierror"
	"github.com/tickstep/aliyunpan-api/aliyunpan_open/openapi"
)

// FileMove 移动文件
func (p *OpenPanClient) FileMove(param *aliyunpan.FileMoveParam) (*aliyunpan.FileMoveResult, *apierror.ApiError) {
	retryTime := 0

RetryBegin:
	opParam := &openapi.FileMoveParam{
		DriveId:        param.DriveId,
		FileId:         param.FileId,
		ToParentFileId: param.ToParentFileId,
	}
	if result, err := p.apiClient.FileMove(opParam); err == nil {
		return &aliyunpan.FileMoveResult{
			FileId:  result.FileId,
			Success: true,
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
