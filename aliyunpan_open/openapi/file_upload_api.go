package openapi

import (
	"encoding/json"
	"fmt"
	"github.com/tickstep/library-go/logger"
	"strings"
)

type (
	PartInfoItem struct {
		// PartNumber 分片序列号，从 1 开始。单个文件分片最大限制5GB，最小限制100KB
		PartNumber int `json:"part_number"`
		// UploadUrl 分片上传URL地址
		UploadUrl string `json:"upload_url"`
		// PartSize 分片大小
		PartSize int64 `json:"part_size"`
	}

	UploadedPartItem struct {
		// Etag 在上传分片结束后，服务端会返回这个分片的Etag，在complete的时候可以在uploadInfo指定分片的Etag，服务端会在合并时对每个分片Etag做校验
		Etag string `json:"etag"`
		// PartNumber 分片序列号，从 1 开始。单个文件分片最大限制5GB，最小限制100KB
		PartNumber int `json:"part_number"`
		// PartSize 分片大小
		PartSize int64 `json:"part_size"`
	}

	// FileUploadCheckPreHashParam 文件PreHash检测参数
	FileUploadCheckPreHashParam struct {
		// DriveId 网盘ID
		DriveId string `json:"drive_id"`
		// ParentFileId 父目录id，上传到根目录时填写 root
		ParentFileId string `json:"parent_file_id"`
		// Name 文件名称，按照 utf8 编码最长 1024 字节，不能以 / 结尾
		Name string `json:"name"`
		// Type file | folder
		Type string `json:"type"`
		// CheckNameMode auto_rename 自动重命名，存在并发问题 ,refuse 同名不创建 ,ignore 同名文件可创建
		CheckNameMode string `json:"check_name_mode"`
		// Size 文件大小，单位为 byte。秒传必须
		Size int64 `json:"size"`
		// PreHash 针对大文件sha1计算非常耗时的情况， 可以先在读取文件的前1k的sha1， 如果前1k的sha1没有匹配的， 那么说明文件无法做秒传， 如果1ksha1有匹配再计算文件sha1进行秒传，这样有效边避免无效的sha1计算。
		PreHash string `json:"pre_hash"`
	}

	// FileUploadCreateParam 文件创建参数
	FileUploadCreateParam struct {
		// DriveId 网盘ID
		DriveId string `json:"drive_id"`
		// ParentFileId 父目录id，上传到根目录时填写 root
		ParentFileId string `json:"parent_file_id"`
		// Name 文件名称，按照 utf8 编码最长 1024 字节，不能以 / 结尾
		Name string `json:"name"`
		// Type file | folder
		Type string `json:"type"`
		// CheckNameMode auto_rename 自动重命名，存在并发问题 ,refuse 同名不创建 ,ignore 同名文件可创建
		CheckNameMode string `json:"check_name_mode"`
		// Size 文件大小，单位为 byte。秒传必须
		Size int64 `json:"size"`
		// 最大分片数量 10000
		PartInfoList []*PartInfoItem `json:"part_info_list"`
		// ContentHash 文件内容 hash 值，需要根据 content_hash_name 指定的算法计算，当前都是sha1算法
		ContentHash string `json:"content_hash"`
		// ContentHashName 秒传必须 ,默认都是 sha1
		ContentHashName string `json:"content_hash_name"`
		// ProofCode 防伪码，秒传必须
		ProofCode string `json:"proof_code"`
		// ProofVersion 固定 v1
		ProofVersion string `json:"proof_version"`
		// LocalCreatedAt 本地创建时间，只对文件有效，格式yyyy-MM-dd'T'HH:mm:ss.SSS'Z'
		LocalCreatedAt string `json:"local_created_at"`
		// LocalModifiedAt 本地修改时间，只对文件有效，格式yyyy-MM-dd'T'HH:mm:ss.SSS'Z'
		LocalModifiedAt string `json:"local_modified_at"`
	}

	FileUploadCreateResult struct {
		// DriveId 网盘ID
		DriveId string `json:"drive_id"`
		// ParentFileId 父目录id，上传到根目录时填写 root
		ParentFileId string `json:"parent_file_id"`
		// FileId 文件ID
		FileId string `json:"file_id"`
		// FileName 文件名称
		FileName string `json:"file_name"`
		// Status
		Status string `json:"status"`
		// UploadId 创建文件夹返回空
		UploadId string `json:"upload_id"`
		// Available
		Available bool `json:"available"`
		// Exist 是否存在同名文件
		Exist bool `json:"exist"`
		// RapidUpload 是否能秒传
		RapidUpload bool `json:"rapid_upload"`
		// 最大分片数量 10000
		PartInfoList []*PartInfoItem `json:"part_info_list"`
	}

	// FileUploadGetUploadUrlParam 刷新获取上传地址参数
	FileUploadGetUploadUrlParam struct {
		// DriveId 网盘ID
		DriveId string `json:"drive_id"`
		// FileId
		FileId string `json:"file_id"`
		// UploadId 文件创建获取的upload_id
		UploadId string `json:"upload_id"`
		// 最大分片数量 10000
		PartInfoList []*PartInfoItem `json:"part_info_list"`
	}
	// FileUploadGetUploadUrlResult 刷新获取上传地址返回值
	FileUploadGetUploadUrlResult struct {
		// DriveId 网盘ID
		DriveId string `json:"drive_id"`
		// FileId
		FileId string `json:"file_id"`
		// UploadId 文件创建获取的upload_id
		UploadId string `json:"upload_id"`
		// CreatedAt 创建时间，格式yyyy-MM-dd'T'HH:mm:ss.SSS'Z'
		CreatedAt string `json:"created_at"`
		// 最大分片数量 10000
		PartInfoList []*PartInfoItem `json:"part_info_list"`
	}

	// FileUploadListUploadedPartsParam 列举已上传分片参数
	FileUploadListUploadedPartsParam struct {
		// DriveId 网盘ID
		DriveId string `json:"drive_id"`
		// FileId
		FileId string `json:"file_id"`
		// UploadId 文件创建获取的upload_id
		UploadId string `json:"upload_id"`
		// PartNumberMarker 分页标记
		PartNumberMarker string `json:"part_number_marker"`
	}
	// FileUploadListUploadedPartsResult 列举已上传分片返回值
	FileUploadListUploadedPartsResult struct {
		// DriveId 网盘ID
		DriveId string `json:"drive_id"`
		// UploadId 文件创建获取的upload_id
		UploadId string `json:"upload_id"`
		// ParallelUpload 是否并行上传
		ParallelUpload bool `json:"parallelUpload"`
		// UploadedParts 已经上传分片列表
		UploadedParts []*UploadedPartItem `json:"uploaded_parts"`
		// NextPartNumberMarker	下一页起始资源标识符, 最后一页该值为空。
		NextPartNumberMarker string `json:"next_part_number_marker"`
	}

	// FileUploadCompleteParam 上传完毕参数
	FileUploadCompleteParam struct {
		// DriveId 网盘ID
		DriveId string `json:"drive_id"`
		// FileId
		FileId string `json:"file_id"`
		// UploadId 文件创建获取的upload_id
		UploadId string `json:"upload_id"`
	}
	// FileUploadCompleteResult 上传完毕返回值
	FileUploadCompleteResult struct {
		// DriveId 网盘ID
		DriveId string `json:"drive_id"`
		// ParentFileId 父目录id
		ParentFileId string `json:"parent_file_id"`
		// FileId 文件id
		FileId string `json:"file_id"`
		// Name 文件名称，按照 utf8 编码最长 1024 字节，不能以 / 结尾
		Name string `json:"name"`
		// Type file | folder
		Type string `json:"type"`
		// Size 文件大小，单位为 byte。秒传必须
		Size int64 `json:"size"`
		// Category 文件类别
		Category string `json:"category"`
		// FileExtension 文件扩展名
		FileExtension string `json:"file_extension"`
		// ContentHash 文件内容 hash 值，需要根据 content_hash_name 指定的算法计算，当前都是sha1算法
		ContentHash string `json:"content_hash"`
		// ContentHashName 秒传必须 ,默认都是 sha1
		ContentHashName string `json:"content_hash_name"`
		// LocalCreatedAt 本地创建时间，只对文件有效，格式yyyy-MM-dd'T'HH:mm:ss.SSS'Z'
		LocalCreatedAt string `json:"local_created_at"`
		// LocalModifiedAt 本地修改时间，只对文件有效，格式yyyy-MM-dd'T'HH:mm:ss.SSS'Z'
		LocalModifiedAt string `json:"local_modified_at"`
	}
)

// FileUploadCheckPreHash 文件PreHash检测
func (a *AliPanClient) FileUploadCheckPreHash(param *FileUploadCheckPreHashParam) (bool, *AliApiErrResult) {
	fullUrl := &strings.Builder{}
	fmt.Fprintf(fullUrl, "%s/adrive/v1.0/openFile/create", OPENAPI_URL)
	logger.Verboseln("do request url: " + fullUrl.String())

	// parameters
	postData := param

	// request
	resp, err := a.httpclient.Req("POST", fullUrl.String(), postData, a.Headers())
	if err != nil {
		logger.Verboseln("file check pre hash error ", err)
		return false, NewAliApiHttpError(err.Error())
	}

	// handler common error
	var apiErrResult *AliApiErrResult
	if _, apiErrResult = ParseCommonOpenApiError(resp); apiErrResult != nil {
		if apiErrResult.Code == "PreHashMatched" {
			return true, nil
		}
		return false, apiErrResult
	}
	return false, nil
}

// FileUploadCreate 文件（文件夹）创建
func (a *AliPanClient) FileUploadCreate(param *FileUploadCreateParam) (*FileUploadCreateResult, *AliApiErrResult) {
	fullUrl := &strings.Builder{}
	fmt.Fprintf(fullUrl, "%s/adrive/v1.0/openFile/create", OPENAPI_URL)
	logger.Verboseln("do request url: " + fullUrl.String())

	// parameters
	postData := map[string]interface{}{
		"drive_id":        param.DriveId,
		"parent_file_id":  param.ParentFileId,
		"name":            param.Name,
		"type":            param.Type,
		"check_name_mode": param.CheckNameMode,
	}
	if len(param.CheckNameMode) == 0 {
		postData["check_name_mode"] = "auto_rename"
	}

	if strings.ToLower(param.Type) == "folder" {
		// 创建文件夹
	} else if strings.ToLower(param.Type) == "file" {
		// 创建文件上传任务
		postData["size"] = param.Size

		// 分片
		if param.PartInfoList != nil {
			parts := []map[string]int{}
			for _, v := range param.PartInfoList {
				parts = append(parts, map[string]int{"part_number": v.PartNumber})
			}
			postData["part_info_list"] = parts
		} else {
			// 默认只有一个分片
			postData["part_info_list"] = []map[string]int{
				{"part_number": 1},
			}
		}

		// 文件hash
		if len(param.ContentHashName) > 0 {
			postData["content_hash_name"] = param.ContentHashName
		}
		if len(param.ContentHash) > 0 {
			postData["content_hash"] = param.ContentHash
		}

		if len(param.ProofVersion) > 0 {
			postData["proof_version"] = param.ProofVersion
		}
		if len(param.ProofCode) > 0 {
			postData["proof_code"] = param.ProofCode
		}

		// 时间
		if len(param.LocalCreatedAt) > 0 {
			postData["local_created_at"] = param.LocalCreatedAt
		}
		if len(param.LocalModifiedAt) > 0 {
			postData["local_modified_at"] = param.LocalModifiedAt
		}
	}

	// request
	resp, err := a.httpclient.Req("POST", fullUrl.String(), postData, a.Headers())
	if err != nil {
		logger.Verboseln("file create error ", err)
		return nil, NewAliApiHttpError(err.Error())
	}

	// handler common error
	var body []byte
	var apiErrResult *AliApiErrResult
	if body, apiErrResult = ParseCommonOpenApiError(resp); apiErrResult != nil {
		return nil, apiErrResult
	}

	// parse result
	r := &FileUploadCreateResult{}
	if err2 := json.Unmarshal(body, r); err2 != nil {
		logger.Verboseln("parse file create result json error ", err2)
		return nil, NewAliApiAppError(err2.Error())
	}
	return r, nil
}

// FileUploadGetUploadUrl 刷新获取上传地址
func (a *AliPanClient) FileUploadGetUploadUrl(param *FileUploadGetUploadUrlParam) (*FileUploadGetUploadUrlResult, *AliApiErrResult) {
	fullUrl := &strings.Builder{}
	fmt.Fprintf(fullUrl, "%s/adrive/v1.0/openFile/getUploadUrl", OPENAPI_URL)
	logger.Verboseln("do request url: " + fullUrl.String())

	// parameters
	postData := param

	// request
	resp, err := a.httpclient.Req("POST", fullUrl.String(), postData, a.Headers())
	if err != nil {
		logger.Verboseln("file create error ", err)
		return nil, NewAliApiHttpError(err.Error())
	}

	// handler common error
	var body []byte
	var apiErrResult *AliApiErrResult
	if body, apiErrResult = ParseCommonOpenApiError(resp); apiErrResult != nil {
		return nil, apiErrResult
	}

	// parse result
	r := &FileUploadGetUploadUrlResult{}
	if err2 := json.Unmarshal(body, r); err2 != nil {
		logger.Verboseln("parse file create result json error ", err2)
		return nil, NewAliApiAppError(err2.Error())
	}
	return r, nil
}

// FileUploadListUploadedParts 列举已上传分片
func (a *AliPanClient) FileUploadListUploadedParts(param *FileUploadListUploadedPartsParam) (*FileUploadListUploadedPartsResult, *AliApiErrResult) {
	fullUrl := &strings.Builder{}
	fmt.Fprintf(fullUrl, "%s/adrive/v1.0/openFile/listUploadedParts", OPENAPI_URL)
	logger.Verboseln("do request url: " + fullUrl.String())

	// parameters
	postData := map[string]interface{}{
		"drive_id":  "19519221",
		"file_id":   "65db23c016b484ac4a0f4d629653442b2e6d9ef9",
		"upload_id": "ADC37FC345574D7BB94E71B1C143F5D3",
	}
	if len(param.PartNumberMarker) > 0 {
		postData["part_number_marker"] = param.PartNumberMarker
	}

	// request
	resp, err := a.httpclient.Req("POST", fullUrl.String(), postData, a.Headers())
	if err != nil {
		logger.Verboseln("file list uploaded parts error ", err)
		return nil, NewAliApiHttpError(err.Error())
	}

	// handler common error
	var body []byte
	var apiErrResult *AliApiErrResult
	if body, apiErrResult = ParseCommonOpenApiError(resp); apiErrResult != nil {
		return nil, apiErrResult
	}

	// parse result
	r := &FileUploadListUploadedPartsResult{}
	if err2 := json.Unmarshal(body, r); err2 != nil {
		logger.Verboseln("parse file create result json error ", err2)
		return nil, NewAliApiAppError(err2.Error())
	}
	return r, nil
}

// FileUploadComplete 上传完毕
func (a *AliPanClient) FileUploadComplete(param *FileUploadCompleteParam) (*FileUploadCompleteResult, *AliApiErrResult) {
	fullUrl := &strings.Builder{}
	fmt.Fprintf(fullUrl, "%s/adrive/v1.0/openFile/complete", OPENAPI_URL)
	logger.Verboseln("do request url: " + fullUrl.String())

	// parameters
	postData := param

	// request
	resp, err := a.httpclient.Req("POST", fullUrl.String(), postData, a.Headers())
	if err != nil {
		logger.Verboseln("file complete error ", err)
		return nil, NewAliApiHttpError(err.Error())
	}

	// handler common error
	var body []byte
	var apiErrResult *AliApiErrResult
	if body, apiErrResult = ParseCommonOpenApiError(resp); apiErrResult != nil {
		return nil, apiErrResult
	}

	// parse result
	r := &FileUploadCompleteResult{}
	if err2 := json.Unmarshal(body, r); err2 != nil {
		logger.Verboseln("parse file complete result json error ", err2)
		return nil, NewAliApiAppError(err2.Error())
	}
	return r, nil
}
