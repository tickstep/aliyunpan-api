package aliyunpan_open

import (
	"github.com/tickstep/aliyunpan-api/aliyunpan/apierror"
	"github.com/tickstep/aliyunpan-api/aliyunpan_open/openapi"
)

// FileRename 重命名文件
func (p *OpenPanClient) FileRename(driveId, renameFileId, newName string) (bool, *apierror.ApiError) {
	retryTime := 0

RetryBegin:
	opParam := &openapi.FileUpdateParam{
		DriveId:       driveId,
		FileId:        renameFileId,
		Name:          newName,
		CheckNameMode: "refuse",
	}
	if result, err := p.apiClient.FileUpdate(opParam); err == nil {
		return result.Name != "", nil
	} else {
		// handle common error
		if apiErrorHandleResp := p.HandleAliApiError(err, &retryTime); apiErrorHandleResp.NeedRetry {
			goto RetryBegin
		} else {
			return false, apiErrorHandleResp.ApiErr
		}
	}
}
