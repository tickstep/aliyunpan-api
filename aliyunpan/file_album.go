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
	// AlbumListParam 相册列表参数
	AlbumListParam struct {
		OrderBy        AlbumOrderBy        `json:"order_by"`
		OrderDirection AlbumOrderDirection `json:"order_direction"`
		Limit          int                 `json:"limit"`
		// Marker 下一页参数
		Marker string `json:"marker"`
	}

	// AlbumEntity 相薄实体
	AlbumEntity struct {
		Owner       string `json:"owner"`
		Name        string `json:"name"`
		Description string `json:"description"`
		AlbumId     string `json:"albumId"`
		FileCount   int    `json:"fileCount"`
		ImageCount  int    `json:"imageCount"`
		VideoCount  int    `json:"videoCount"`
		CreatedAt   string `json:"createdAt"`
		UpdatedAt   string `json:"updatedAt"`
		IsSharing   bool   `json:"isSharing"`
	}
	AlbumList []*AlbumEntity

	AlbumListResult struct {
		Items AlbumList `json:"items"`
		// NextMarker 不为空，说明还有下一页
		NextMarker string `json:"nextMarker"`
	}

	albumEntityResult struct {
		Owner       string `json:"owner"`
		Name        string `json:"name"`
		Description string `json:"description"`
		AlbumId     string `json:"album_id"`
		FileCount   int    `json:"file_count"`
		ImageCount  int    `json:"image_count"`
		VideoCount  int    `json:"video_count"`
		CreatedAt   int64  `json:"created_at"`
		UpdatedAt   int64  `json:"updated_at"`
		IsSharing   bool   `json:"is_sharing"`
	}
	albumListResult struct {
		Items []*albumEntityResult `json:"items"`
		// NextMarker 不为空，说明还有下一页
		NextMarker string `json:"next_marker"`
	}

	AlbumOrderBy        string
	AlbumOrderDirection string

	// AlbumEditParam 相簿编辑参数
	AlbumEditParam struct {
		AlbumId     string `json:"albumId"`
		Description string `json:"description"`
		Name        string `json:"name"`
	}

	// AlbumDeleteParam 相簿删除参数
	AlbumDeleteParam struct {
		AlbumId string `json:"albumId"`
	}

	// AlbumGetParam 相簿查询参数
	AlbumGetParam struct {
		AlbumId string `json:"albumId"`
	}
)

const (
	AlbumOrderByCreatedAt AlbumOrderBy = "created_at"
	AlbumOrderByUpdatedAt AlbumOrderBy = "updated_at"
	AlbumOrderByFileCount AlbumOrderBy = "file_count"

	// AlbumOrderDirectionDesc 降序
	AlbumOrderDirectionDesc AlbumOrderDirection = "DESC"
	// AlbumOrderDirectionAsc 升序
	AlbumOrderDirectionAsc AlbumOrderDirection = "ASC"
)

func createAlbumEntity(f *albumEntityResult) *AlbumEntity {
	if f == nil {
		return nil
	}
	return &AlbumEntity{
		Owner:       f.Owner,
		Name:        f.Name,
		Description: f.Description,
		AlbumId:     f.AlbumId,
		FileCount:   f.FileCount,
		ImageCount:  f.ImageCount,
		VideoCount:  f.VideoCount,
		CreatedAt:   apiutil.UnixTime2LocalFormat(f.CreatedAt),
		UpdatedAt:   apiutil.UnixTime2LocalFormat(f.UpdatedAt),
		IsSharing:   f.IsSharing,
	}
}

// AlbumListGetAll 获取所有相册列表
func (p *PanClient) AlbumListGetAll(param *AlbumListParam) (AlbumList, *apierror.ApiError) {
	internalParam := &AlbumListParam{
		OrderBy:        param.OrderBy,
		OrderDirection: param.OrderDirection,
		Limit:          param.Limit,
		Marker:         param.Marker,
	}
	if internalParam.Limit <= 0 {
		internalParam.Limit = 100
	}

	fileList := AlbumList{}
	result, err := p.AlbumList(internalParam)
	if err != nil || result == nil {
		return nil, err
	}
	fileList = append(fileList, result.Items...)

	// more page?
	for len(result.NextMarker) > 0 {
		internalParam.Marker = result.NextMarker
		result, err = p.AlbumList(internalParam)
		if err == nil && result != nil {
			fileList = append(fileList, result.Items...)
		} else {
			break
		}
	}
	return fileList, nil
}

// AlbumList 获取相册列表
func (p *PanClient) AlbumList(param *AlbumListParam) (*AlbumListResult, *apierror.ApiError) {
	result := &AlbumListResult{
		Items:      AlbumList{},
		NextMarker: "",
	}
	if flr, err := p.albumListReq(param); err == nil {
		for k := range flr.Items {
			if flr.Items[k] == nil {
				continue
			}
			result.Items = append(result.Items, createAlbumEntity(flr.Items[k]))
		}
		result.NextMarker = flr.NextMarker
	}
	return result, nil
}

func (p *PanClient) albumListReq(param *AlbumListParam) (*albumListResult, *apierror.ApiError) {
	header := map[string]string{
		"authorization": p.webToken.GetAuthorizationStr(),
	}

	fullUrl := &strings.Builder{}
	fmt.Fprintf(fullUrl, "%s/adrive/v1/album/list", API_URL)
	logger.Verboseln("do request url: " + fullUrl.String())

	limit := param.Limit
	if limit <= 0 {
		limit = 100
	}
	if param.OrderBy == "" {
		param.OrderBy = AlbumOrderByCreatedAt
	}
	if param.OrderDirection == "" {
		param.OrderDirection = AlbumOrderDirectionAsc
	}
	postData := map[string]interface{}{
		"limit":           limit,
		"order_by":        param.OrderBy,
		"order_direction": param.OrderDirection,
	}
	if len(param.Marker) > 0 {
		postData["marker"] = param.Marker
	}

	// request
	body, err := client.Fetch("POST", fullUrl.String(), postData, apiutil.AddCommonHeader(header))
	if err != nil {
		logger.Verboseln("get album list error ", err)
		return nil, apierror.NewFailedApiError(err.Error())
	}

	// handler common error
	if err1 := apierror.ParseCommonApiError(body); err1 != nil {
		return nil, err1
	}

	// parse result
	r := &albumListResult{}
	if err2 := json.Unmarshal(body, r); err2 != nil {
		logger.Verboseln("parse album list result json error ", err2)
		return nil, apierror.NewFailedApiError(err2.Error())
	}
	return r, nil
}

// AlbumEdit 相簿编辑
func (p *PanClient) AlbumEdit(param *AlbumEditParam) (*AlbumEntity, *apierror.ApiError) {
	header := map[string]string{
		"authorization": p.webToken.GetAuthorizationStr(),
	}

	fullUrl := &strings.Builder{}
	fmt.Fprintf(fullUrl, "%s/adrive/v1/album/update", API_URL)
	logger.Verboseln("do request url: " + fullUrl.String())

	if param.AlbumId == "" {
		return nil, apierror.NewFailedApiError("album id cannot be empty")
	}
	if param.Name == "" {
		return nil, apierror.NewFailedApiError("album name cannot be empty")
	}

	postData := map[string]interface{}{
		"album_id":    param.AlbumId,
		"name":        param.Name,
		"description": param.Description,
	}

	// request
	body, err := client.Fetch("POST", fullUrl.String(), postData, apiutil.AddCommonHeader(header))
	if err != nil {
		logger.Verboseln("edit album error ", err)
		return nil, apierror.NewFailedApiError(err.Error())
	}

	// handler common error
	if err1 := apierror.ParseCommonApiError(body); err1 != nil {
		return nil, err1
	}

	// parse result
	r := &AlbumEntity{}
	if err2 := json.Unmarshal(body, r); err2 != nil {
		logger.Verboseln("parse album edit result json error ", err2)
		return nil, apierror.NewFailedApiError(err2.Error())
	}
	return r, nil
}

// AlbumDelete 相簿删除
func (p *PanClient) AlbumDelete(param *AlbumDeleteParam) (bool, *apierror.ApiError) {
	header := map[string]string{
		"authorization": p.webToken.GetAuthorizationStr(),
	}

	fullUrl := &strings.Builder{}
	fmt.Fprintf(fullUrl, "%s/adrive/v1/album/delete", API_URL)
	logger.Verboseln("do request url: " + fullUrl.String())

	if param.AlbumId == "" {
		return false, apierror.NewFailedApiError("album id cannot be empty")
	}

	postData := map[string]interface{}{
		"album_id": param.AlbumId,
	}

	// request
	body, err := client.Fetch("POST", fullUrl.String(), postData, apiutil.AddCommonHeader(header))
	if err != nil {
		logger.Verboseln("delete album error ", err)
		return false, apierror.NewFailedApiError(err.Error())
	}

	// handler common error
	if err1 := apierror.ParseCommonApiError(body); err1 != nil {
		return false, err1
	}

	return true, nil
}

// AlbumGet 相簿删除
func (p *PanClient) AlbumGet(param *AlbumGetParam) (*AlbumEntity, *apierror.ApiError) {
	header := map[string]string{
		"authorization": p.webToken.GetAuthorizationStr(),
	}

	fullUrl := &strings.Builder{}
	fmt.Fprintf(fullUrl, "%s/adrive/v1/album/get", API_URL)
	logger.Verboseln("do request url: " + fullUrl.String())

	if param.AlbumId == "" {
		return nil, apierror.NewFailedApiError("album id cannot be empty")
	}

	postData := map[string]interface{}{
		"album_id": param.AlbumId,
	}

	// request
	body, err := client.Fetch("POST", fullUrl.String(), postData, apiutil.AddCommonHeader(header))
	if err != nil {
		logger.Verboseln("get album error ", err)
		return nil, apierror.NewFailedApiError(err.Error())
	}

	// handler common error
	if err1 := apierror.ParseCommonApiError(body); err1 != nil {
		return nil, err1
	}

	// parse result
	r := &AlbumEntity{}
	if err2 := json.Unmarshal(body, r); err2 != nil {
		logger.Verboseln("parse album get result json error ", err2)
		return nil, apierror.NewFailedApiError(err2.Error())
	}
	return r, nil
}
