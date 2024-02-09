package openapi

import (
	"encoding/json"
	"fmt"
	"github.com/tickstep/library-go/logger"
	"strings"
)

type (
	// FileListParam 获取文件列表参数
	FileListParam struct {
		// DriveId 网盘id
		DriveId string `json:"drive_id"`
		// ParentFileId 根目录为root
		ParentFileId string `json:"parent_file_id"`
		// Limit 返回文件数量，默认 50，最大 100
		Limit int `json:"limit"`
		// Marker 分页标记
		Marker string `json:"marker"`
		// OrderBy created_at,updated_at,name,size,name_enhanced（对数字编号的文件友好，排序结果为 1、2、3...99 而不是 1、10、11...2、21...9、91...99）
		OrderBy string `json:"order_by"`
		// OrderDirection DESC ASC
		OrderDirection string `json:"order_direction"`
		// Type all | file | folder，默认所有类型
		Type string `json:"type"`
		// Fields 当填 * 时，返回文件所有字段。或某些字段，逗号分隔： id_path,name_path
		Fields string `json:"fields"`
	}

	// FileItem 文件信息
	FileItem struct {
		// DriveId 网盘id
		DriveId string `json:"drive_id"`
		// ParentFileId 根目录为root
		ParentFileId string `json:"parent_file_id"`
		// FileId 文件ID
		FileId string `json:"file_id"`
		// Name 文件名
		Name string `json:"name"`
		// FileExtension 文件扩展名
		FileExtension string `json:"file_extension"`
		// Type 文件类型 file | folder
		Type string `json:"type"`
		// Size 文件大小
		Size int64 `json:"size"`
		// Category 文件类别，例如：doc
		Category string `json:"category"`
		Hidden   bool   `json:"hidden"`
		// Status 状态，available-可用
		Status string `json:"status"`
		// Url 图片预览图地址、小于 5MB 文件的下载地址。超过5MB 请使用 /getDownloadUrl
		Url string `json:"url"`
		// Starred 收藏标志，true-被收藏，false-未收藏
		Starred bool `json:"starred"`
		// MimeType 多媒体类型，例如：application/pdf
		MimeType string `json:"mime_type"`
		// ContentHashName 内容hash计算算法 sha1
		ContentHashName string `json:"content_hash_name"`
		// ContentHash 内容hash
		ContentHash string `json:"content_hash"`
		DomainId    string `json:"domain_id"`
		// ContentType 内容类型，例如：application/oct-stream
		ContentType string `json:"content_type"`
		// CreatedAt 创建时间
		CreatedAt string `json:"created_at"`
		// UpdatedAt 修改时间
		UpdatedAt string `json:"updated_at"`
		// PunishFlag 文件禁止标志，0-正常，103-禁止下载
		PunishFlag int `json:"punish_flag"`
	}

	// FileListResult 获取文件列表返回值
	FileListResult struct {
		Items      []*FileItem `json:"items"`
		NextMarker string      `json:"next_marker"`
	}

	// FileSearchParam 文件搜索参数
	FileSearchParam struct {
		// DriveId 网盘id
		DriveId string `json:"drive_id"`
		// Query 查询语句，query拼接的条件 <= 5个 例如： parent_file_id = 'root' and name = '123' and category = 'video'
		Query string `json:"query"`
		// Limit 返回文件数量，默认 50，最大 100
		Limit int `json:"limit"`
		// Marker 分页标记
		Marker string `json:"marker"`
		// OrderBy 排序，例如：created_at ASC | DESC
		OrderBy string `json:"order_by"`
		// ReturnTotalCount 是否返回总数
		ReturnTotalCount bool `json:"return_total_count"`
	}
	// FileSearchResult 文件搜索参数返回值
	FileSearchResult struct {
		Items      []*FileItem `json:"items"`
		NextMarker string      `json:"next_marker"`
		TotalCount int64       `json:"total_count"`
	}

	// FileStarredListParam 获取收藏文件列表参数
	FileStarredListParam struct {
		// DriveId 网盘id
		DriveId string `json:"drive_id"`
		// Limit 返回文件数量，默认 50，最大 100
		Limit int `json:"limit"`
		// Marker 分页标记
		Marker string `json:"marker"`
		// OrderBy created_at,updated_at,name,size,name_enhanced（对数字编号的文件友好，排序结果为 1、2、3...99 而不是 1、10、11...2、21...9、91...99）
		OrderBy string `json:"order_by"`
		// OrderDirection DESC ASC
		OrderDirection string `json:"order_direction"`
		// Type all | file | folder，默认所有类型
		Type string `json:"type"`
	}
)

// FileList 获取文件列表
func (a *AliPanClient) FileList(param *FileListParam) (*FileListResult, *AliApiErrResult) {
	fullUrl := &strings.Builder{}
	fmt.Fprintf(fullUrl, "%s/adrive/v1.0/openFile/list", OPENAPI_URL)
	logger.Verboseln("do request url: " + fullUrl.String())

	// parameters
	postData := param
	if postData.Limit > 100 {
		postData.Limit = 100
	}

	// request
	resp, err := a.httpclient.Req("POST", fullUrl.String(), postData, a.Headers())
	if err != nil {
		logger.Verboseln("get file list info error ", err)
		return nil, NewAliApiHttpError(err.Error())
	}

	// handler common error
	var body []byte
	var apiErrResult *AliApiErrResult
	if body, apiErrResult = ParseCommonOpenApiError(resp); apiErrResult != nil {
		return nil, apiErrResult
	}

	// parse result
	r := &FileListResult{}
	if err2 := json.Unmarshal(body, r); err2 != nil {
		logger.Verboseln("parse file list info result json error ", err2)
		return nil, NewAliApiAppError(err2.Error())
	}
	return r, nil
}

// FileSearch 搜索文件
func (a *AliPanClient) FileSearch(param *FileSearchParam) (*FileSearchResult, *AliApiErrResult) {
	fullUrl := &strings.Builder{}
	fmt.Fprintf(fullUrl, "%s/adrive/v1.0/openFile/search", OPENAPI_URL)
	logger.Verboseln("do request url: " + fullUrl.String())

	// parameters
	postData := param
	if postData.Limit > 100 {
		postData.Limit = 100
	}

	// request
	resp, err := a.httpclient.Req("POST", fullUrl.String(), postData, a.Headers())
	if err != nil {
		logger.Verboseln("get file search error ", err)
		return nil, NewAliApiHttpError(err.Error())
	}

	// handler common error
	var body []byte
	var apiErrResult *AliApiErrResult
	if body, apiErrResult = ParseCommonOpenApiError(resp); apiErrResult != nil {
		return nil, apiErrResult
	}

	// parse result
	r := &FileSearchResult{}
	if err2 := json.Unmarshal(body, r); err2 != nil {
		logger.Verboseln("parse file search result json error ", err2)
		return nil, NewAliApiAppError(err2.Error())
	}
	return r, nil
}

// FileStarredList 获取收藏文件列表
func (a *AliPanClient) FileStarredList(param *FileStarredListParam) (*FileListResult, *AliApiErrResult) {
	fullUrl := &strings.Builder{}
	fmt.Fprintf(fullUrl, "%s/adrive/v1.0/openFile/starredList", OPENAPI_URL)
	logger.Verboseln("do request url: " + fullUrl.String())

	// parameters
	postData := param
	if postData.Limit > 100 {
		postData.Limit = 100
	}

	// request
	resp, err := a.httpclient.Req("POST", fullUrl.String(), postData, a.Headers())
	if err != nil {
		logger.Verboseln("get file starred list error ", err)
		return nil, NewAliApiHttpError(err.Error())
	}

	// handler common error
	var body []byte
	var apiErrResult *AliApiErrResult
	if body, apiErrResult = ParseCommonOpenApiError(resp); apiErrResult != nil {
		return nil, apiErrResult
	}

	// parse result
	r := &FileListResult{}
	if err2 := json.Unmarshal(body, r); err2 != nil {
		logger.Verboseln("parse file starred list result json error ", err2)
		return nil, NewAliApiAppError(err2.Error())
	}
	return r, nil
}
