package openapi

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type (
	// AliApiErrResult openapi错误响应
	AliApiErrResult struct {
		HttpStatusCode int    `json:"http_status_code"`
		Code           string `json:"code"`
		Message        string `json:"message"`
	}
)

const (
	AliErr = ""
)

func NewAliApiError(httpStatusCode int, code, msg string) *AliApiErrResult {
	return &AliApiErrResult{
		HttpStatusCode: httpStatusCode,
		Code:           code,
		Message:        msg,
	}
}
func NewAliApiHttpError(msg string) *AliApiErrResult {
	return &AliApiErrResult{
		HttpStatusCode: 200,
		Code:           "TS.HttpError",
		Message:        msg,
	}
}
func NewAliApiAppError(msg string) *AliApiErrResult {
	return &AliApiErrResult{
		HttpStatusCode: 200,
		Code:           "TS.AppError",
		Message:        msg,
	}
}

// ParseCommonOpenApiError 解析阿里云盘API错误，如果没有错误则返回nil
func ParseCommonOpenApiError(resp *http.Response) ([]byte, *AliApiErrResult) {
	if resp == nil {
		return nil, nil
	}

	// read response text
	data, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		return nil, NewAliApiError(resp.StatusCode, "TS.ReadError", e.Error())
	}
	errResult := &AliApiErrResult{}
	if err := json.Unmarshal(data, errResult); err == nil {
		if errResult.Code != "" {
			errResult.HttpStatusCode = resp.StatusCode
			return nil, errResult
		}
	}
	return data, nil
}
