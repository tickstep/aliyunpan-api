package openapi

import (
	"encoding/json"
	"fmt"
	"github.com/tickstep/aliyunpan-api/aliyunpan"
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

	// ShareAlbumListFileParam 获取共享相册文件列表参数
	ShareAlbumListFileParam struct {
		// AlbumId 共享相册唯一ID
		AlbumId string `json:"sharedAlbumId"`
		// OrderBy 排序字段，当前仅支持joined_at
		OrderBy string `json:"order_by"`
		// OrderDirection 排序方向，默认 DESC。ASC 升序，DESC 降序。
		OrderDirection string `json:"order_direction"`
		// Marker 分页标记
		Marker string `json:"marker"`
		// Limit 返回文件数量，默认50
		Limit int `json:"limit"`
		// ImageThumbnailWidth 生成的图片缩略图宽度，默认480px
		ImageThumbnailWidth int `json:"image_thumbnail_width"`
	}
	// ShareAlbumListFileResult 获取共享相册文件列表返回值
	ShareAlbumListFileResult struct {
		// Items 文件列表
		Items []*FileItem `json:"items"`
		// NextMarker 不为空代表还有下一页
		NextMarker string `json:"nextMarker"`
	}

	// ShareAlbumGetFileUrlParam 获取共享相册下文件下载地址参数
	ShareAlbumGetFileUrlParam struct {
		// AlbumId 共享相册唯一ID
		AlbumId string `json:"sharedAlbumId"`
		// DriveId 文件所属drive
		DriveId string `json:"drive_id"`
		// FileId 文件id
		FileId string `json:"file_id"`
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

// ShareAlbumListFile 获取共享相册包含图片视频文件列表
func (a *AliPanClient) ShareAlbumListFile(param *ShareAlbumListFileParam) (*ShareAlbumListFileResult, *AliApiErrResult) {
	fullUrl := &strings.Builder{}
	fmt.Fprintf(fullUrl, "%s/adrive/v1.0/sharedAlbum/listFile", OPENAPI_URL)
	logger.Verboseln("do request url: " + fullUrl.String())

	// parameters
	postData := param

	// request
	resp, err := a.httpclient.Req("POST", fullUrl.String(), postData, a.Headers())
	if err != nil {
		logger.Verboseln("list file of share album error ", err)
		return nil, NewAliApiHttpError(err.Error())
	}

	// handler common error
	var body []byte
	var apiErrResult *AliApiErrResult
	if body, apiErrResult = ParseCommonOpenApiError(resp); apiErrResult != nil {
		return nil, apiErrResult
	}

	// parse result
	r := &ShareAlbumListFileResult{}
	if err2 := json.Unmarshal(body, r); err2 != nil {
		logger.Verboseln("parse list file of share album result json error ", err2)
		return nil, NewAliApiAppError(err2.Error())
	}
	// 补全 album ID
	if r.Items != nil {
		for _, item := range r.Items {
			item.AlbumId = param.AlbumId
		}
	}
	return r, nil
}

// ShareAlbumGetFileDownloadUrl 获取共享相册下文件下载地址
func (a *AliPanClient) ShareAlbumGetFileDownloadUrl(param *ShareAlbumGetFileUrlParam) (*aliyunpan.ShareAlbumGetFileUrlResult, *AliApiErrResult) {
	fullUrl := &strings.Builder{}
	fmt.Fprintf(fullUrl, "%s/adrive/v1.0/sharedAlbum/getDownloadUrl", OPENAPI_URL)
	logger.Verboseln("do request url: " + fullUrl.String())

	// parameters
	postData := param

	// request
	resp, err := a.httpclient.Req("POST", fullUrl.String(), postData, a.Headers())
	if err != nil {
		logger.Verboseln("get share album file download url error ", err)
		return nil, NewAliApiHttpError(err.Error())
	}

	// handler common error
	var body []byte
	var apiErrResult *AliApiErrResult
	if body, apiErrResult = ParseCommonOpenApiError(resp); apiErrResult != nil {
		return nil, apiErrResult
	}

	// parse result
	r := &aliyunpan.ShareAlbumGetFileUrlResult{}
	if err2 := json.Unmarshal(body, r); err2 != nil {
		logger.Verboseln("parse download url of share album file result json error ", err2)
		return nil, NewAliApiAppError(err2.Error())
	}
	return r, nil
}
