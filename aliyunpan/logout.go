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
	DeviceLogoutResult struct {
		Result  bool   `json:"result"`
		Success bool   `json:"success"`
		Code    string `json:"code"`
		Message string `json:"message"`
	}
)

// DeviceLogout 退出登录，登录的设备会同步注销
func (p *PanClient) DeviceLogout() (*DeviceLogoutResult, *apierror.ApiError) {
	// header
	header := map[string]string{
		"authorization": p.webToken.GetAuthorizationStr(),
	}

	// url
	fullUrl := &strings.Builder{}
	fmt.Fprintf(fullUrl, "%s/users/v1/users/device_logout", API_URL)
	logger.Verboseln("do request url: " + fullUrl.String())

	// data
	postData := map[string]interface{}{}

	// request
	body, err := p.client.Fetch("POST", fullUrl.String(), postData, p.AddSignatureHeader(apiutil.AddCommonHeader(header)))
	if err != nil {
		logger.Verboseln("device logout error ", err)
		return nil, apierror.NewFailedApiError(err.Error())
	}
	logger.Verboseln("device logout response: " + string(body))

	// handler common error
	if err1 := apierror.ParseCommonApiError(body); err1 != nil {
		return nil, err1
	}

	// parse result
	r := &DeviceLogoutResult{}
	if err2 := json.Unmarshal(body, r); err2 != nil {
		logger.Verboseln("parse device logout result json error ", err2)
		return nil, apierror.NewFailedApiError(err2.Error())
	}
	return r, nil
}
