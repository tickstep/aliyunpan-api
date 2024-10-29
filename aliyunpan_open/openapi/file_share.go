package openapi

import (
	"encoding/json"
	"fmt"
	"github.com/tickstep/library-go/logger"
	"strings"
)

type (
	// FileShareCreateParam 创建文件分享参数
	FileShareCreateParam struct {
		// DriveId 网盘id
		DriveId string `json:"driveId"`
		// 文件ID数据，元数数量 [1, 100]
		FileIdList []string `json:"fileIdList"`
		// 分享过期时间，格式：2024-09-19T09:32:50.000Z
		Expiration string `json:"expiration"`
		// 分享提取码
		SharePwd string `json:"sharePwd"`
	}
	// FileShareCreateResult 创建文件分享返回值
	FileShareCreateResult struct {
		// 分享ID
		ShareId string `json:"share_id"`
		// 分享过期时间
		Expiration string `json:"expiration"`
		// 分享是否已过期
		Expired bool `json:"expired"`
		// 分享提取码
		SharePwd string `json:"share_pwd"`
		// 分享链接地址
		ShareUrl string `json:"share_url"`
		// 分享创建者ID
		Creator string `json:"creator"`
		// 分享当前状态
		Status string `json:"status"`
		// 分享创建时间，格式：2024-09-14T02:11:34.264Z
		CreatedAt string `json:"created_at"`
		// 分享更新时间，格式：2024-09-14T02:11:34.264Z
		UpdatedAt string `json:"update_at"`
	}
)

// FileShareCreate 创建文件分享
func (a *AliPanClient) FileShareCreate(param *FileShareCreateParam) (*FileShareCreateResult, *AliApiErrResult) {
	fullUrl := &strings.Builder{}
	fmt.Fprintf(fullUrl, "%s/adrive/v1.0/openFile/createShare", OPENAPI_URL)
	logger.Verboseln("do request url: " + fullUrl.String())

	// parameters
	postData := param

	// request
	resp, err := a.httpclient.Req("POST", fullUrl.String(), postData, a.Headers())
	if err != nil {
		logger.Verboseln("create file share error ", err)
		return nil, NewAliApiHttpError(err.Error())
	}

	// handler common error
	var body []byte
	var apiErrResult *AliApiErrResult
	if body, apiErrResult = ParseCommonOpenApiError(resp); apiErrResult != nil {
		return nil, apiErrResult
	}

	// parse result
	r := &FileShareCreateResult{}
	if err2 := json.Unmarshal(body, r); err2 != nil {
		logger.Verboseln("parse file share result json error ", err2)
		return nil, NewAliApiAppError(err2.Error())
	}
	return r, nil
}
