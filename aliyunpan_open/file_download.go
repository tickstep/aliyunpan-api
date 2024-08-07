package aliyunpan_open

import (
	"fmt"
	"github.com/tickstep/aliyunpan-api/aliyunpan"
	"github.com/tickstep/aliyunpan-api/aliyunpan/apierror"
	"github.com/tickstep/aliyunpan-api/aliyunpan/apiutil"
	"github.com/tickstep/aliyunpan-api/aliyunpan_open/openapi"
	"github.com/tickstep/library-go/logger"
	"strconv"
	"strings"
)

// GetFileDownloadUrl 获取文件下载URL路径
func (p *OpenPanClient) GetFileDownloadUrl(param *aliyunpan.GetFileDownloadUrlParam) (*aliyunpan.GetFileDownloadUrlResult, *apierror.ApiError) {
	retryTime := 0

RetryBegin:
	opParam := &openapi.FileDownloadUrlParam{
		DriveId:   param.DriveId,
		FileId:    param.FileId,
		ExpireSec: param.ExpireSec,
	}
	if result, err := p.apiClient.FileGetDownloadUrl(opParam); err == nil {
		return &aliyunpan.GetFileDownloadUrlResult{
			Method:      result.Method,
			Url:         result.Url,
			InternalUrl: "",
			CdnUrl:      "",
			Expiration:  apiutil.UtcTime2LocalFormat(result.Expiration),
			Size:        result.Size,
			Description: result.Description,
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

// DownloadFileData 下载文件内容
func (p *OpenPanClient) DownloadFileData(downloadFileUrl string, fileRange aliyunpan.FileDownloadRange, downloadFunc aliyunpan.DownloadFuncCallback) *apierror.ApiError {
	// url
	fullUrl := &strings.Builder{}
	fmt.Fprintf(fullUrl, "%s", downloadFileUrl)
	logger.Verboseln("do request url: " + fullUrl.String())

	// header
	headers := map[string]string{
		"user-agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
		"referer":    "https://www.aliyundrive.com/",
	}

	// download data resume
	if fileRange.Offset != 0 || fileRange.End != 0 {
		rangeStr := "bytes=" + strconv.FormatInt(fileRange.Offset, 10) + "-"
		if fileRange.End != 0 {
			rangeStr += strconv.FormatInt(fileRange.End, 10)
		}
		headers["range"] = rangeStr
	}
	logger.Verboseln("do request url: " + fullUrl.String())

	// request callback
	_, err := downloadFunc("GET", fullUrl.String(), headers)
	//resp, err := p.client.Req("GET", fullUrl.String(), nil, headers)

	if err != nil {
		logger.Verboseln("download file data response failed")
		return apierror.NewApiErrorWithError(err)
	}
	return nil
}
