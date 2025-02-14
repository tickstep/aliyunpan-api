package openapi

import (
	"encoding/json"
	"fmt"
	"github.com/tickstep/library-go/logger"
	"strings"
)

type (
	// ShareAlbumListParam 获取共享相册列表参数
	ShareAlbumListParam struct {
	}
	// ShareAlbumListResult 获取共享相册列表返回值
	ShareAlbumListResult struct {
		// 共享相册项列表
		Items []*ShareAlbumItem `json:"items"`
	}

	// ShareAlbumItem 共享相册项
	ShareAlbumItem struct {
		// SharedAlbumId 共享相册唯一ID
		SharedAlbumId string `json:"sharedAlbumId"`
		// Name 共相册名称
		Name string `json:"name"`
		// Description 共享相册简介
		Description string `json:"description"`
		// CoverThumbnail 封面图地址
		CoverThumbnail string `json:"coverThumbnail"`
		// CreatedAt 分享创建时间
		CreatedAt int64 `json:"createdAt"`
		// UpdatedAt 分享更新时间
		UpdatedAt int64 `json:"updatedAt"`
	}
)

// ShareAlbumList 获取共享相册列表
func (a *AliPanClient) ShareAlbumList(param *ShareAlbumListParam) (*ShareAlbumListResult, *AliApiErrResult) {
	fullUrl := &strings.Builder{}
	fmt.Fprintf(fullUrl, "%s/adrive/v1.0/sharedAlbum/list", OPENAPI_URL)
	logger.Verboseln("do request url: " + fullUrl.String())

	// parameters
	postData := param

	// request
	resp, err := a.httpclient.Req("POST", fullUrl.String(), postData, a.Headers())
	if err != nil {
		logger.Verboseln("list share album error ", err)
		return nil, NewAliApiHttpError(err.Error())
	}

	// handler common error
	var body []byte
	var apiErrResult *AliApiErrResult
	if body, apiErrResult = ParseCommonOpenApiError(resp); apiErrResult != nil {
		return nil, apiErrResult
	}

	// parse result
	r := &ShareAlbumListResult{}
	if err2 := json.Unmarshal(body, r); err2 != nil {
		logger.Verboseln("parse list share album result json error ", err2)
		return nil, NewAliApiAppError(err2.Error())
	}
	return r, nil
}
