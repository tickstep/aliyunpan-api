package aliyunpan_open

import (
	"github.com/tickstep/aliyunpan-api/aliyunpan"
	"github.com/tickstep/aliyunpan-api/aliyunpan/apierror"
	"github.com/tickstep/aliyunpan-api/aliyunpan_open/openapi"
)

// FileMove 移动文件
func (p *OpenPanClient) FileMove(param []*aliyunpan.FileMoveParam) ([]*aliyunpan.FileMoveResult, *apierror.ApiError) {
	retryTime := 0
	returnResult := []*aliyunpan.FileMoveResult{}

	for _, v := range param {
	RetryBegin:
		opParam := &openapi.FileMoveParam{
			DriveId:        v.DriveId,
			FileId:         v.FileId,
			ToParentFileId: v.ToParentFileId,
		}
		if result, err := p.apiClient.FileMove(opParam); err == nil {
			returnResult = append(returnResult, &aliyunpan.FileMoveResult{
				FileId:  result.FileId,
				Success: true,
			})
		} else {
			// handle common error
			if apiErrorHandleResp := p.HandleAliApiError(err, &retryTime); apiErrorHandleResp.NeedRetry {
				goto RetryBegin
			} else {
				returnResult = append(returnResult, &aliyunpan.FileMoveResult{
					FileId:  result.FileId,
					Success: false,
				})
			}
		}
	}
	return returnResult, nil
}
