package openapi

import (
	"encoding/json"
	"fmt"
	"github.com/tickstep/library-go/logger"
	"strings"
)

type (
	// AsyncTaskQueryStatusParam 查询异步任务状态参数
	AsyncTaskQueryStatusParam struct {
		// AsyncTaskId 异步任务ID
		AsyncTaskId string `json:"async_task_id"`
	}
	// AsyncTaskQueryStatusResult 查询异步任务状态返回值
	AsyncTaskQueryStatusResult struct {
		// State Succeed 成功，Running 处理中，Failed 已失败
		State string `json:"state"`
		// AsyncTaskId 异步任务ID
		AsyncTaskId string `json:"async_task_id"`
	}
)

// AsyncTaskQueryStatus 获取异步任务状态
func (a *AliPanClient) AsyncTaskQueryStatus(param *AsyncTaskQueryStatusParam) (*AsyncTaskQueryStatusResult, *AliApiErrResult) {
	fullUrl := &strings.Builder{}
	fmt.Fprintf(fullUrl, "%s/adrive/v1.0/openFile/async_task/get", OPENAPI_URL)
	logger.Verboseln("do request url: " + fullUrl.String())

	// parameters
	postData := param

	// request
	resp, err := a.httpclient.Req("POST", fullUrl.String(), postData, a.Headers())
	if err != nil {
		logger.Verboseln("async task status error ", err)
		return nil, NewAliApiHttpError(err.Error())
	}

	// handler common error
	var body []byte
	var apiErrResult *AliApiErrResult
	if body, apiErrResult = ParseCommonOpenApiError(resp); apiErrResult != nil {
		return nil, apiErrResult
	}

	// parse result
	r := &AsyncTaskQueryStatusResult{}
	if err2 := json.Unmarshal(body, r); err2 != nil {
		logger.Verboseln("parse async task status result json error ", err2)
		return nil, NewAliApiAppError(err2.Error())
	}
	return r, nil
}
