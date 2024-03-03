package aliyunpan_web

import (
	"encoding/json"
	"fmt"
	"github.com/tickstep/aliyunpan-api/aliyunpan/apierror"
	"github.com/tickstep/aliyunpan-api/aliyunpan/apiutil"
	"github.com/tickstep/library-go/logger"
	"strings"
)

type (
	AsyncTaskQueryStatusParam struct {
		AsyncTaskId string `json:"async_task_id"`
	}

	AsyncTaskQueryStatusResult struct {
		AsyncTaskId       string `json:"async_task_id"`
		State             string `json:"state"`
		Status            string `json:"status"`
		TotalProcess      int    `json:"total_process"`
		ConsumedProcess   int    `json:"consumed_process"`
		SkippedProcess    int    `json:"skipped_process"`
		FailedProcess     int    `json:"failed_process"`
		PunishedFileCount int    `json:"punished_file_count"`
	}
)

// AsyncTaskQueryStatus 查询异步任务进度和状态
func (p *WebPanClient) AsyncTaskQueryStatus(param *AsyncTaskQueryStatusParam) (*AsyncTaskQueryStatusResult, *apierror.ApiError) {
	header := map[string]string{
		"authorization": p.webToken.GetAuthorizationStr(),
		"referer":       "https://www.aliyundrive.com/",
		"origin":        "https://www.aliyundrive.com",
	}

	fullUrl := &strings.Builder{}
	fmt.Fprintf(fullUrl, "%s/v2/async_task/get", API_URL)
	logger.Verboseln("do request url: " + fullUrl.String())

	postData := map[string]interface{}{
		"async_task_id": param.AsyncTaskId,
	}
	// request
	body, err := p.client.Fetch("POST", fullUrl.String(), postData, p.AddSignatureHeader(apiutil.AddCommonHeader(header)))
	if err != nil {
		logger.Verboseln("async task query status error ", err)
		return nil, apierror.NewFailedApiError(err.Error())
	}

	// handler common error
	if err1 := apierror.ParseCommonApiError(body); err1 != nil {
		return nil, err1
	}

	// parse result
	r := &AsyncTaskQueryStatusResult{}
	if err2 := json.Unmarshal(body, r); err2 != nil {
		logger.Verboseln("parse async task query status result json error ", err2)
		return nil, apierror.NewFailedApiError(err2.Error())
	}
	return r, nil
}
