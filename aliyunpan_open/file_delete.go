package aliyunpan_open

import (
	"github.com/tickstep/aliyunpan-api/aliyunpan"
	"github.com/tickstep/aliyunpan-api/aliyunpan/apierror"
	"github.com/tickstep/aliyunpan-api/aliyunpan_open/openapi"
)

// FileDelete 删除文件到回收站
func (p *OpenPanClient) FileDelete(param *aliyunpan.FileBatchActionParam) (*aliyunpan.FileBatchActionResult, *apierror.ApiError) {
	retryTime := 0

RetryBegin:
	opParam := &openapi.FileIdentityPair{
		DriveId: param.DriveId,
		FileId:  param.FileId,
	}
	if result, err := p.apiClient.FileTrash(opParam); err == nil {
		return &aliyunpan.FileBatchActionResult{
			FileId:  result.FileId,
			Success: true,
		}, nil
	} else {
		// handle common error
		if apiErrorHandleResp := p.HandleAliApiError(err, &retryTime); apiErrorHandleResp.NeedRetry {
			goto RetryBegin
		} else {
			return &aliyunpan.FileBatchActionResult{
				FileId:  result.FileId,
				Success: false,
			}, apiErrorHandleResp.ApiErr
		}
	}
}
