package aliyunpan

import (
	"encoding/json"
	"fmt"
	"github.com/tickstep/aliyunpan-api/aliyunpan/apierror"
	"github.com/tickstep/aliyunpan-api/aliyunpan/apiutil"
	"github.com/tickstep/library-go/logger"
	"strings"
)

type (
	// 请求参数
	BatchRequest struct {
		Id     string `json:"id"`
		Method string `json:"method"`
		Url    string `json:"url"`

		Headers map[string]string      `json:"headers"`
		Body    map[string]interface{} `json:"body"`
	}
	BatchRequestList  []*BatchRequest
	BatchRequestParam struct {
		Requests BatchRequestList `json:"requests"`
		Resource string           `json:"resource"`
	}

	// 响应结果
	BatchResponse struct {
		Id     string                 `json:"id"`
		Status int                    `json:"status"`
		Body   map[string]interface{} `json:"body"`
	}
	BatchResponseList   []*BatchResponse
	BatchResponseResult struct {
		Responses BatchResponseList `json:"responses"`
	}
)

// BatchTask 批量请求任务。多选操作基本都是批量任务
func (p *PanClient) BatchTask(url string, param *BatchRequestParam) (*BatchResponseResult, *apierror.ApiError) {
	if param == nil {
		return nil, apierror.NewFailedApiError("参数不能为空")
	}

	// header
	header := map[string]string{
		"authorization": p.webToken.GetAuthorizationStr(),
	}

	// url
	fullUrl := &strings.Builder{}
	fmt.Fprintf(fullUrl, "%s", url)
	logger.Verboseln("do request url: " + fullUrl.String())

	// data
	postData := param

	// request
	body, err := p.client.Fetch("POST", fullUrl.String(), postData, p.AddSignatureHeader(apiutil.AddCommonHeader(header)))
	if err != nil {
		logger.Verboseln("batch request error ", err)
		return nil, apierror.NewFailedApiError(err.Error())
	}

	// handler common error
	if err1 := apierror.ParseCommonApiError(body); err1 != nil {
		return nil, err1
	}

	// parse result
	r := &BatchResponseResult{}
	if err2 := json.Unmarshal(body, r); err2 != nil {
		logger.Verboseln("batch result json error ", err2)
		return nil, apierror.NewFailedApiError(err2.Error())
	}
	return r, nil
}
