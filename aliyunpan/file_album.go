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
		AlbumId     string `json:"album_id"`
		FileCount   int    `json:"file_count"`
		ImageCount  int    `json:"image_count"`
		VideoCount  int    `json:"video_count"`
		CreatedAt   int64  `json:"created_at"`
		UpdatedAt   int64  `json:"updated_at"`
		IsSharing   bool   `json:"is_sharing"`
	}
	AlbumList []*AlbumEntity

	AlbumListResult struct {
		Items AlbumList `json:"items"`
		// NextMarker 不为空，说明还有下一页
		NextMarker string `json:"next_marker"`
	}

	AlbumOrderBy        string
	AlbumOrderDirection string

	// AlbumCreateParam 相簿创建参数
	AlbumCreateParam struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

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

	// AlbumShareCreateParam 创建相簿分享
	AlbumShareCreateParam struct {
		AlbumId string `json:"album_id"`
		// 分享密码，4个字符，为空代码公开分享
		SharePwd string `json:"share_pwd"`
		// 过期时间，为空代表永不过期。时间格式必须是这种：2021-07-23 09:22:19
		Expiration string `json:"expiration"`
	}

	// AlbumShareCreateResult 创建相簿分享返回值
	AlbumShareCreateResult struct {
		Album             *AlbumEntity `json:"album"`
		Popularity        int          `json:"popularity"`
		ShareID           string       `json:"share_id"`
		ShareMsg          string       `json:"share_msg"`
		ShareName         string       `json:"share_name"`
		Description       string       `json:"description"`
		Expiration        string       `json:"expiration"`
		Expired           bool         `json:"expired"`
		SharePwd          string       `json:"share_pwd"`
		ShareURL          string       `json:"share_url"`
		Creator           string       `json:"creator"`
		DriveID           string       `json:"drive_id"`
		FileID            string       `json:"file_id"`
		AlbumID           string       `json:"album_id"`
		PreviewCount      int          `json:"preview_count"`
		SaveCount         int          `json:"save_count"`
		DownloadCount     int          `json:"download_count"`
		Status            string       `json:"status"`
		CreatedAt         string       `json:"created_at"`
		UpdatedAt         string       `json:"updated_at"`
		IsPhotoCollection bool         `json:"is_photo_collection"`
		SyncToHomepage    bool         `json:"sync_to_homepage"`
		PopularityStr     string       `json:"popularity_str"`
		FullShareMsg      string       `json:"full_share_msg"`
		DisplayName       string       `json:"display_name"`
	}

	// AlbumListFileParam 相簿查询包含的文件列表
	AlbumListFileParam struct {
		AlbumId string `json:"albumId"`
		Limit   int    `json:"limit"`
		// Marker 下一页参数
		Marker string `json:"marker"`
	}

	// AlbumDeleteFileParam 相簿删除文件参数
	AlbumDeleteFileParam struct {
		AlbumId       string                 `json:"album_id"`
		DriveFileList []FileBatchActionParam `json:"drive_file_list"`
	}

	// AlbumAddFileParam 相簿增加文件参数
	AlbumAddFileParam struct {
		AlbumId       string                 `json:"album_id"`
		DriveFileList []FileBatchActionParam `json:"drive_file_list"`
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

func (a *AlbumEntity) CreatedAtStr() string {
	return apiutil.UnixTime2LocalFormat(a.CreatedAt)
}
func (a *AlbumEntity) UpdatedAtStr() string {
	return apiutil.UnixTime2LocalFormat(a.UpdatedAt)
}

func (a *AlbumDeleteFileParam) AddFileItem(driveId, fileId string) {
	if a.DriveFileList == nil {
		a.DriveFileList = []FileBatchActionParam{}
	}
	a.DriveFileList = append(a.DriveFileList, FileBatchActionParam{
		DriveId: driveId,
		FileId:  fileId,
	})
}

func (a *AlbumAddFileParam) AddFileItem(driveId, fileId string) {
	if a.DriveFileList == nil {
		a.DriveFileList = []FileBatchActionParam{}
	}
	a.DriveFileList = append(a.DriveFileList, FileBatchActionParam{
		DriveId: driveId,
		FileId:  fileId,
	})
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
			result.Items = append(result.Items, flr.Items[k])
		}
		result.NextMarker = flr.NextMarker
	}
	return result, nil
}

func (p *PanClient) albumListReq(param *AlbumListParam) (*AlbumListResult, *apierror.ApiError) {
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
	body, err := p.client.Fetch("POST", fullUrl.String(), postData, p.AddSignatureHeader(apiutil.AddCommonHeader(header)))
	if err != nil {
		logger.Verboseln("get album list error ", err)
		return nil, apierror.NewFailedApiError(err.Error())
	}

	// handler common error
	if err1 := apierror.ParseCommonApiError(body); err1 != nil {
		return nil, err1
	}

	// parse result
	r := &AlbumListResult{}
	if err2 := json.Unmarshal(body, r); err2 != nil {
		logger.Verboseln("parse album list result json error ", err2)
		return nil, apierror.NewFailedApiError(err2.Error())
	}
	return r, nil
}

// AlbumEdit 相簿编辑
func (p *PanClient) AlbumCreate(param *AlbumCreateParam) (*AlbumEntity, *apierror.ApiError) {
	header := map[string]string{
		"authorization": p.webToken.GetAuthorizationStr(),
	}

	fullUrl := &strings.Builder{}
	fmt.Fprintf(fullUrl, "%s/adrive/v1/album/create", API_URL)
	logger.Verboseln("do request url: " + fullUrl.String())

	if param.Name == "" {
		return nil, apierror.NewFailedApiError("album name cannot be empty")
	}

	postData := map[string]interface{}{
		"name":        param.Name,
		"description": param.Description,
	}

	// request
	body, err := p.client.Fetch("POST", fullUrl.String(), postData, p.AddSignatureHeader(apiutil.AddCommonHeader(header)))
	if err != nil {
		logger.Verboseln("create album error ", err)
		return nil, apierror.NewFailedApiError(err.Error())
	}

	// handler common error
	if err1 := apierror.ParseCommonApiError(body); err1 != nil {
		return nil, err1
	}

	// parse result
	r := &AlbumEntity{}
	if err2 := json.Unmarshal(body, r); err2 != nil {
		logger.Verboseln("parse album create result json error ", err2)
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
	body, err := p.client.Fetch("POST", fullUrl.String(), postData, p.AddSignatureHeader(apiutil.AddCommonHeader(header)))
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
	body, err := p.client.Fetch("POST", fullUrl.String(), postData, p.AddSignatureHeader(apiutil.AddCommonHeader(header)))
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

// AlbumGet 获取相簿信息
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
	body, err := p.client.Fetch("POST", fullUrl.String(), postData, p.AddSignatureHeader(apiutil.AddCommonHeader(header)))
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

// AlbumShareCreate 相簿创建分享链接
func (p *PanClient) AlbumShareCreate(param *AlbumShareCreateParam) (*AlbumShareCreateResult, *apierror.ApiError) {
	// header
	header := map[string]string{
		"authorization": p.webToken.GetAuthorizationStr(),
	}

	// url
	fullUrl := &strings.Builder{}
	fmt.Fprintf(fullUrl, "%s/adrive/v2/share_link/create", API_URL)
	logger.Verboseln("do request url: " + fullUrl.String())

	// data
	postData := param

	// check pwd
	if postData.SharePwd != "" && len(postData.SharePwd) != 4 {
		return nil, apierror.NewFailedApiError("密码必须是4个字符")
	}

	// format time
	if postData.Expiration != "" {
		postData.Expiration = apiutil.LocalTime2UtcFormat(param.Expiration)
	}

	// request
	body, err := p.client.Fetch("POST", fullUrl.String(), postData, p.AddSignatureHeader(apiutil.AddCommonHeader(header)))
	if err != nil {
		logger.Verboseln("create album share error ", err)
		return nil, apierror.NewFailedApiError(err.Error())
	}
	logger.Verboseln("response: ", string(body))

	// handler common error
	if err1 := apierror.ParseCommonApiError(body); err1 != nil {
		return nil, err1
	}

	// parse result
	r := &AlbumShareCreateResult{}
	if err2 := json.Unmarshal(body, r); err2 != nil {
		logger.Verboseln("parse album share result json error ", err2)
		return nil, apierror.NewFailedApiError(err2.Error())
	}
	return r, nil
}

// AlbumListFileGetAll 获取指定相簿下的所有文件列表
func (p *PanClient) AlbumListFileGetAll(param *AlbumListFileParam) (FileList, *apierror.ApiError) {
	internalParam := &AlbumListFileParam{
		AlbumId: param.AlbumId,
		Limit:   param.Limit,
		Marker:  param.Marker,
	}
	if internalParam.Limit <= 0 {
		internalParam.Limit = 100
	}

	fileList := FileList{}
	result, err := p.AlbumListFile(internalParam)
	if err != nil || result == nil {
		return nil, err
	}
	fileList = append(fileList, result.FileList...)

	// more page?
	for len(result.NextMarker) > 0 {
		internalParam.Marker = result.NextMarker
		result, err = p.AlbumListFile(internalParam)
		if err == nil && result != nil {
			fileList = append(fileList, result.FileList...)
		} else {
			break
		}
	}
	return fileList, nil
}

// AlbumListFile 获取相簿下的文件列表
func (p *PanClient) AlbumListFile(param *AlbumListFileParam) (*FileListResult, *apierror.ApiError) {
	result := &FileListResult{
		FileList:   FileList{},
		NextMarker: "",
	}
	if flr, err := p.albumListFileReq(param); err == nil {
		for k := range flr.Items {
			if flr.Items[k] == nil {
				continue
			}
			result.FileList = append(result.FileList, createFileEntity(flr.Items[k]))
		}
		result.NextMarker = flr.NextMarker
	}
	return result, nil
}

func (p *PanClient) albumListFileReq(param *AlbumListFileParam) (*fileListResult, *apierror.ApiError) {
	header := map[string]string{
		"authorization": p.webToken.GetAuthorizationStr(),
	}

	fullUrl := &strings.Builder{}
	fmt.Fprintf(fullUrl, "%s/adrive/v1/album/list_files", API_URL)
	logger.Verboseln("do request url: " + fullUrl.String())

	limit := param.Limit
	if limit <= 0 {
		limit = 100
	}
	postData := map[string]interface{}{
		"album_id":                param.AlbumId,
		"image_thumbnail_process": "image/resize,w_400/format,jpeg",
		"video_thumbnail_process": "video/snapshot,t_0,f_jpg,ar_auto,w_1000",
		"image_url_process":       "image/resize,w_1920/format,jpeg",
		"filter":                  "",
		"fields":                  "*",
		"limit":                   param.Limit,
		"order_by":                "joined_at",
		"order_direction":         "DESC",
	}
	if len(param.Marker) > 0 {
		postData["marker"] = param.Marker
	}

	// request
	body, err := p.client.Fetch("POST", fullUrl.String(), postData, p.AddSignatureHeader(apiutil.AddCommonHeader(header)))
	if err != nil {
		logger.Verboseln("get album file list error ", err)
		return nil, apierror.NewFailedApiError(err.Error())
	}

	// handler common error
	if err1 := apierror.ParseCommonApiError(body); err1 != nil {
		return nil, err1
	}

	// parse result
	r := &fileListResult{}
	if err2 := json.Unmarshal(body, r); err2 != nil {
		logger.Verboseln("parse album file list result json error ", err2)
		return nil, apierror.NewFailedApiError(err2.Error())
	}
	return r, nil
}

// AlbumDeleteFile 相簿删除文件列表
func (p *PanClient) AlbumDeleteFile(param *AlbumDeleteFileParam) (bool, *apierror.ApiError) {
	header := map[string]string{
		"authorization": p.webToken.GetAuthorizationStr(),
	}

	fullUrl := &strings.Builder{}
	fmt.Fprintf(fullUrl, "%s/adrive/v1/album/delete_files", API_URL)
	logger.Verboseln("do request url: " + fullUrl.String())

	if param.AlbumId == "" {
		return false, apierror.NewFailedApiError("album id cannot be empty")
	}
	postData := param

	// request
	body, err := p.client.Fetch("POST", fullUrl.String(), postData, p.AddSignatureHeader(apiutil.AddCommonHeader(header)))
	if err != nil {
		logger.Verboseln("delete album file error ", err)
		return false, apierror.NewFailedApiError(err.Error())
	}

	// handler common error
	if err1 := apierror.ParseCommonApiError(body); err1 != nil {
		return false, err1
	}

	return true, nil
}

// AlbumAddFile 相簿增加文件列表
func (p *PanClient) AlbumAddFile(param *AlbumAddFileParam) (*FileList, *apierror.ApiError) {
	header := map[string]string{
		"authorization": p.webToken.GetAuthorizationStr(),
	}

	fullUrl := &strings.Builder{}
	fmt.Fprintf(fullUrl, "%s/adrive/v1/album/add_files", API_URL)
	logger.Verboseln("do request url: " + fullUrl.String())

	if param.AlbumId == "" {
		return nil, apierror.NewFailedApiError("album id cannot be empty")
	}
	postData := param

	// request
	body, err := p.client.Fetch("POST", fullUrl.String(), postData, p.AddSignatureHeader(apiutil.AddCommonHeader(header)))
	if err != nil {
		logger.Verboseln("add album file error ", err)
		return nil, apierror.NewFailedApiError(err.Error())
	}

	// handler common error
	if err1 := apierror.ParseCommonApiError(body); err1 != nil {
		return nil, err1
	}

	// parse result
	type fileListResult struct {
		Items []*fileEntityResult `json:"file_list"`
	}
	r := &fileListResult{}
	if err2 := json.Unmarshal(body, r); err2 != nil {
		logger.Verboseln("parse add album file result json error ", err2)
		return nil, apierror.NewFailedApiError(err2.Error())
	}
	fileList := FileList{}
	for k := range r.Items {
		if r.Items[k] == nil {
			continue
		}
		fileList = append(fileList, createFileEntity(r.Items[k]))
	}
	return &fileList, nil
}
