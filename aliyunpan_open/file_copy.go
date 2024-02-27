package aliyunpan_open

import (
	"github.com/tickstep/aliyunpan-api/aliyunpan"
	"github.com/tickstep/aliyunpan-api/aliyunpan/apierror"
	"github.com/tickstep/aliyunpan-api/aliyunpan_open/openapi"
)

// FileCopy 同网盘内复制文件或文件夹
func (p *OpenPanClient) FileCopy(param *aliyunpan.FileCopyParam) (*aliyunpan.FileAsyncTaskResult, *apierror.ApiError) {
	retryTime := 0

RetryBegin:
	opParam := &openapi.FileCopyParam{
		DriveId:        param.DriveId,
		FileId:         param.FileId,
		ToDriveId:      param.DriveId,
		ToParentFileId: param.ToParentFileId,
		AutoRename:     true,
	}
	if result, err := p.apiClient.FileCopy(opParam); err == nil {
		return &aliyunpan.FileAsyncTaskResult{
			DriveId:     result.DriveId,
			FileId:      result.FileId,
			AsyncTaskId: result.AsyncTaskId,
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
