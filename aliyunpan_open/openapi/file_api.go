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

	// FileIdentityPair 文件唯一标识对
	FileIdentityPair struct {
		// DriveId 网盘id
		DriveId string `json:"drive_id"`
		// FileId 文件ID
		FileId string `json:"file_id"`
	}

	// FilePathPair 文件路径唯一标识对
	FilePathPair struct {
		// DriveId 网盘id
		DriveId string `json:"drive_id"`
		// FilePath 文件绝对路径，文件或者文件夹。 例子： 1. 根目录下有 a.jepg 文件，传 /a.jpeg 获取文件属性 2. 根目录下有 bb 文件夹，传 /bb 获取文件夹属性
		FilePath string `json:"file_path"`
	}

	// FileDownloadUrlParam 获取文件下载链接参数
	FileDownloadUrlParam struct {
		// DriveId 网盘id
		DriveId string `json:"drive_id"`
		// FileId 文件ID
		FileId string `json:"file_id"`
		// ExpireSec 下载地址过期时间，单位为秒，默认为 900 秒,最长4h（14400秒，需要申请）
		ExpireSec int64 `json:"expire_sec"`
	}
	// FileDownloadUrlResult 文件下载链接返回值
	FileDownloadUrlResult struct {
		Method          string `json:"method"`
		Url             string `json:"url"`
		Expiration      string `json:"expiration"`
		Size            int64  `json:"size"`
		StreamsUrl      string `json:"streamsUrl"`
		ContentHash     string `json:"content_hash"`
		ContentHashName string `json:"content_hash_name"`
		FileId          string `json:"file_id"`
	}

	// FileUpdateParam 文件更新参数
	FileUpdateParam struct {
		// DriveId 网盘id
		DriveId string `json:"drive_id"`
		// FileId 文件ID
		FileId string `json:"file_id"`
		// Name 新的文件名，支持后缀名更改
		Name string `json:"name"`
		// CheckNameMode auto_rename-自动重命名 refuse-同名不创建 ignore-同名文件可创建。 默认ignore
		CheckNameMode string `json:"check_name_mode"`
		// Starred 收藏 true，移除收藏 false
		Starred bool `json:"starred"`
	}

	// FileMoveParam 文件移动参数
	FileMoveParam struct {
		// DriveId 网盘id
		DriveId string `json:"drive_id"`
		// FileId 文件ID
		FileId string `json:"file_id"`
		// ToParentFileId 目标目录ID、根目录为 root
		ToParentFileId string `json:"to_parent_file_id"`
		// CheckNameMode auto_rename-自动重命名 refuse-同名不创建 ignore-同名文件可创建。 默认ignore
		CheckNameMode string `json:"check_name_mode"`
		// NewName 当云端存在同名文件时，使用的新名字
		NewName string `json:"new_name"`
	}
	// FileMoveResult 文件移动返回值
	FileMoveResult struct {
		// Exist 文件是否已存在
		Exist bool `json:"exist"`
		// DriveId 网盘id
		DriveId string `json:"drive_id"`
		// FileId 文件ID
		FileId string `json:"file_id"`
		// AsyncTaskId 异步任务id。 如果返回为空字符串，表示直接移动成功。 如果返回非空字符串，表示需要经过异步处理。
		AsyncTaskId string `json:"async_task_id"`
	}

	// FileCopyParam 文件复制参数
	FileCopyParam struct {
		// DriveId 网盘id
		DriveId string `json:"drive_id"`
		// FileId 文件ID
		FileId string `json:"file_id"`
		// ToDriveId 目标网盘id
		ToDriveId string `json:"to_drive_id"`
		// ToParentFileId 目标目录ID、根目录为 root
		ToParentFileId string `json:"to_parent_file_id"`
		// AutoRename 当目标文件夹下存在同名文件时，是否自动重命名，默认为 false，默认允许同名文件
		AutoRename bool `json:"auto_rename"`
	}

	// FileAsyncTaskResult 文件异步操作返回值
	FileAsyncTaskResult struct {
		// DriveId 网盘id
		DriveId string `json:"drive_id"`
		// FileId 文件ID
		FileId string `json:"file_id"`
		// AsyncTaskId 异步任务id。 如果返回为空字符串，表示直接移动成功。 如果返回非空字符串，表示需要经过异步处理。
		AsyncTaskId string `json:"async_task_id"`
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

// FileGetDetailInfo 获取文件详情
func (a *AliPanClient) FileGetDetailInfo(param *FileIdentityPair) (*FileItem, *AliApiErrResult) {
	fullUrl := &strings.Builder{}
	fmt.Fprintf(fullUrl, "%s/adrive/v1.0/openFile/get", OPENAPI_URL)
	logger.Verboseln("do request url: " + fullUrl.String())

	// parameters
	postData := param

	// request
	resp, err := a.httpclient.Req("POST", fullUrl.String(), postData, a.Headers())
	if err != nil {
		logger.Verboseln("get file detail info error ", err)
		return nil, NewAliApiHttpError(err.Error())
	}

	// handler common error
	var body []byte
	var apiErrResult *AliApiErrResult
	if body, apiErrResult = ParseCommonOpenApiError(resp); apiErrResult != nil {
		return nil, apiErrResult
	}

	// parse result
	r := &FileItem{}
	if err2 := json.Unmarshal(body, r); err2 != nil {
		logger.Verboseln("parse file detail info result json error ", err2)
		return nil, NewAliApiAppError(err2.Error())
	}
	return r, nil
}

// FileGetDetailInfoByPath 文件路径查找文件
func (a *AliPanClient) FileGetDetailInfoByPath(param *FilePathPair) (*FileItem, *AliApiErrResult) {
	fullUrl := &strings.Builder{}
	fmt.Fprintf(fullUrl, "%s/adrive/v1.0/openFile/get_by_path", OPENAPI_URL)
	logger.Verboseln("do request url: " + fullUrl.String())

	// parameters
	postData := param

	// request
	resp, err := a.httpclient.Req("POST", fullUrl.String(), postData, a.Headers())
	if err != nil {
		logger.Verboseln("get file detail by path error ", err)
		return nil, NewAliApiHttpError(err.Error())
	}

	// handler common error
	var body []byte
	var apiErrResult *AliApiErrResult
	if body, apiErrResult = ParseCommonOpenApiError(resp); apiErrResult != nil {
		return nil, apiErrResult
	}

	// parse result
	r := &FileItem{}
	if err2 := json.Unmarshal(body, r); err2 != nil {
		logger.Verboseln("parse file detail info by path result json error ", err2)
		return nil, NewAliApiAppError(err2.Error())
	}
	return r, nil
}

// FileGetDetailInfoBatch 批量获取文件详情
func (a *AliPanClient) FileGetDetailInfoBatch(param []*FileIdentityPair) (*FileListResult, *AliApiErrResult) {
	fullUrl := &strings.Builder{}
	fmt.Fprintf(fullUrl, "%s/adrive/v1.0/openFile/batch/get", OPENAPI_URL)
	logger.Verboseln("do request url: " + fullUrl.String())

	// parameters
	postData := map[string]interface{}{
		"file_list": param,
	}

	// request
	resp, err := a.httpclient.Req("POST", fullUrl.String(), postData, a.Headers())
	if err != nil {
		logger.Verboseln("batch get file detail info error ", err)
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
		logger.Verboseln("parse batch file detail info result json error ", err2)
		return nil, NewAliApiAppError(err2.Error())
	}
	return r, nil
}

// FileGetDownloadUrl 获取文件下载链接
func (a *AliPanClient) FileGetDownloadUrl(param *FileDownloadUrlParam) (*FileDownloadUrlResult, *AliApiErrResult) {
	fullUrl := &strings.Builder{}
	fmt.Fprintf(fullUrl, "%s/adrive/v1.0/openFile/getDownloadUrl", OPENAPI_URL)
	logger.Verboseln("do request url: " + fullUrl.String())

	// parameters
	postData := param
	if postData.ExpireSec <= 0 || postData.ExpireSec >= 14400 {
		postData.ExpireSec = 14400
	}

	// request
	resp, err := a.httpclient.Req("POST", fullUrl.String(), postData, a.Headers())
	if err != nil {
		logger.Verboseln("get file download url error ", err)
		return nil, NewAliApiHttpError(err.Error())
	}

	// handler common error
	var body []byte
	var apiErrResult *AliApiErrResult
	if body, apiErrResult = ParseCommonOpenApiError(resp); apiErrResult != nil {
		return nil, apiErrResult
	}

	// parse result
	r := &FileDownloadUrlResult{}
	if err2 := json.Unmarshal(body, r); err2 != nil {
		logger.Verboseln("parse get file download url result json error ", err2)
		return nil, NewAliApiAppError(err2.Error())
	}
	return r, nil
}

// FileUpdate 文件更新
func (a *AliPanClient) FileUpdate(param *FileUpdateParam) (*FileItem, *AliApiErrResult) {
	fullUrl := &strings.Builder{}
	fmt.Fprintf(fullUrl, "%s/adrive/v1.0/openFile/update", OPENAPI_URL)
	logger.Verboseln("do request url: " + fullUrl.String())

	// parameters
	postData := map[string]interface{}{
		"drive_id": param.DriveId,
		"file_id":  param.FileId,
	}
	if len(param.Name) > 0 {
		postData["name"] = param.Name
		postData["check_name_mode"] = param.CheckNameMode
		if len(param.CheckNameMode) == 0 {
			postData["check_name_mode"] = "refuse"
		}
	} else {
		postData["starred"] = param.Starred
	}

	// request
	resp, err := a.httpclient.Req("POST", fullUrl.String(), postData, a.Headers())
	if err != nil {
		logger.Verboseln("file update error ", err)
		return nil, NewAliApiHttpError(err.Error())
	}

	// handler common error
	var body []byte
	var apiErrResult *AliApiErrResult
	if body, apiErrResult = ParseCommonOpenApiError(resp); apiErrResult != nil {
		return nil, apiErrResult
	}

	// parse result
	r := &FileItem{}
	if err2 := json.Unmarshal(body, r); err2 != nil {
		logger.Verboseln("parse file update result json error ", err2)
		return nil, NewAliApiAppError(err2.Error())
	}
	return r, nil
}

// FileMove 移动文件或文件夹
func (a *AliPanClient) FileMove(param *FileMoveParam) (*FileMoveResult, *AliApiErrResult) {
	fullUrl := &strings.Builder{}
	fmt.Fprintf(fullUrl, "%s/adrive/v1.0/openFile/move", OPENAPI_URL)
	logger.Verboseln("do request url: " + fullUrl.String())

	// parameters
	postData := param
	if len(postData.CheckNameMode) == 0 {
		postData.CheckNameMode = "auto_rename"
	}

	// request
	resp, err := a.httpclient.Req("POST", fullUrl.String(), postData, a.Headers())
	if err != nil {
		logger.Verboseln("file move error ", err)
		return nil, NewAliApiHttpError(err.Error())
	}

	// handler common error
	var body []byte
	var apiErrResult *AliApiErrResult
	if body, apiErrResult = ParseCommonOpenApiError(resp); apiErrResult != nil {
		return nil, apiErrResult
	}

	// parse result
	r := &FileMoveResult{}
	if err2 := json.Unmarshal(body, r); err2 != nil {
		logger.Verboseln("parse file move result json error ", err2)
		return nil, NewAliApiAppError(err2.Error())
	}
	return r, nil
}

// FileCopy 复制文件或文件夹
func (a *AliPanClient) FileCopy(param *FileCopyParam) (*FileAsyncTaskResult, *AliApiErrResult) {
	fullUrl := &strings.Builder{}
	fmt.Fprintf(fullUrl, "%s/adrive/v1.0/openFile/copy", OPENAPI_URL)
	logger.Verboseln("do request url: " + fullUrl.String())

	// parameters
	postData := param

	// request
	resp, err := a.httpclient.Req("POST", fullUrl.String(), postData, a.Headers())
	if err != nil {
		logger.Verboseln("file copy error ", err)
		return nil, NewAliApiHttpError(err.Error())
	}

	// handler common error
	var body []byte
	var apiErrResult *AliApiErrResult
	if body, apiErrResult = ParseCommonOpenApiError(resp); apiErrResult != nil {
		return nil, apiErrResult
	}

	// parse result
	r := &FileAsyncTaskResult{}
	if err2 := json.Unmarshal(body, r); err2 != nil {
		logger.Verboseln("parse file copy result json error ", err2)
		return nil, NewAliApiAppError(err2.Error())
	}
	return r, nil
}

// FileTrash 把文件或文件夹放入回收站
func (a *AliPanClient) FileTrash(param *FileIdentityPair) (*FileAsyncTaskResult, *AliApiErrResult) {
	fullUrl := &strings.Builder{}
	fmt.Fprintf(fullUrl, "%s/adrive/v1.0/openFile/recyclebin/trash", OPENAPI_URL)
	logger.Verboseln("do request url: " + fullUrl.String())

	// parameters
	postData := param

	// request
	resp, err := a.httpclient.Req("POST", fullUrl.String(), postData, a.Headers())
	if err != nil {
		logger.Verboseln("file trash error ", err)
		return nil, NewAliApiHttpError(err.Error())
	}

	// handler common error
	var body []byte
	var apiErrResult *AliApiErrResult
	if body, apiErrResult = ParseCommonOpenApiError(resp); apiErrResult != nil {
		return nil, apiErrResult
	}

	// parse result
	r := &FileAsyncTaskResult{}
	if err2 := json.Unmarshal(body, r); err2 != nil {
		logger.Verboseln("parse file trash result json error ", err2)
		return nil, NewAliApiAppError(err2.Error())
	}
	return r, nil
}

// FileDelete 文件直接删除，不放到回收站直接删除
func (a *AliPanClient) FileDelete(param *FileIdentityPair) (*FileAsyncTaskResult, *AliApiErrResult) {
	fullUrl := &strings.Builder{}
	fmt.Fprintf(fullUrl, "%s/adrive/v1.0/openFile/delete", OPENAPI_URL)
	logger.Verboseln("do request url: " + fullUrl.String())

	// parameters
	postData := param

	// request
	resp, err := a.httpclient.Req("POST", fullUrl.String(), postData, a.Headers())
	if err != nil {
		logger.Verboseln("file delete error ", err)
		return nil, NewAliApiHttpError(err.Error())
	}

	// handler common error
	var body []byte
	var apiErrResult *AliApiErrResult
	if body, apiErrResult = ParseCommonOpenApiError(resp); apiErrResult != nil {
		return nil, apiErrResult
	}

	// parse result
	r := &FileAsyncTaskResult{}
	if err2 := json.Unmarshal(body, r); err2 != nil {
		logger.Verboseln("parse file delete result json error ", err2)
		return nil, NewAliApiAppError(err2.Error())
	}
	return r, nil
}
