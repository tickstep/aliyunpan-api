package aliyunpan_open

import (
	"github.com/tickstep/aliyunpan-api/aliyunpan"
	"github.com/tickstep/aliyunpan-api/aliyunpan_open/openapi"
)

// VideoGetPreviewPlayInfo 获取视频预览信息，调用该接口会触发视频云端转码
func (p *OpenPanClient) VideoGetPreviewPlayInfo(param *aliyunpan.VideoGetPreviewPlayInfoParam) (*aliyunpan.VideoGetPreviewPlayInfoResult, error) {
	retryTime := 0

RetryBegin:
	opParam := &openapi.VideoGetPreviewPlayInfoParam{
		DriveId:  param.DriveId,
		FileId:   param.FileId,
		Category: "live_transcoding",
	}
	if result, err := p.apiClient.VideoGetPreviewPlayInfo(opParam); err == nil {
		return &aliyunpan.VideoGetPreviewPlayInfoResult{
			DriveId: result.DriveId,
			FileId:  result.FileId,
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
