package aliyunpan_web

import (
	"encoding/json"
	"fmt"
	"github.com/tickstep/aliyunpan-api/aliyunpan"
	"github.com/tickstep/aliyunpan-api/aliyunpan/apierror"
	"github.com/tickstep/aliyunpan-api/aliyunpan/apiutil"
	"github.com/tickstep/library-go/logger"
	"github.com/tickstep/library-go/requester"
	"strings"
)

type (
	completeUploadFileReqResult struct {
		DriveId         string `json:"drive_id"`
		DomainId        string `json:"domain_id"`
		FileId          string `json:"file_id"`
		Name            string `json:"name"`
		Type            string `json:"type"`
		ContentType     string `json:"content_type"`
		CreatedAt       string `json:"created_at"`
		UpdatedAt       string `json:"updated_at"`
		FileExtension   string `json:"file_extension"`
		Hidden          bool   `json:"hidden"`
		Size            int64  `json:"size"`
		Starred         bool   `json:"starred"`
		Status          string `json:"status"`
		UploadId        string `json:"upload_id"`
		ParentFileId    string `json:"parent_file_id"`
		Crc64Hash       string `json:"crc64_hash"`
		ContentHash     string `json:"content_hash"`
		ContentHashName string `json:"content_hash_name"`
		Category        string `json:"category"`
		EncryptMode     string `json:"encrypt_mode"`
		Location        string `json:"location"`
	}
)

const ()

// CreateUploadFile 创建上传文件，如果文件已经上传过则会直接秒传
func (p *WebPanClient) CreateUploadFile(param *aliyunpan.CreateFileUploadParam) (*aliyunpan.CreateFileUploadResult, *apierror.ApiError) {
	// header
	header := map[string]string{
		"authorization": p.webToken.GetAuthorizationStr(),
	}

	// url
	fullUrl := &strings.Builder{}
	fmt.Fprintf(fullUrl, "%s/adrive/v2/file/createWithFolders", API_URL)
	logger.Verboseln("do request url: " + fullUrl.String())

	// data
	postData := param

	if len(postData.PartInfoList) == 0 {
		blockSize := aliyunpan.DefaultChunkSize
		if param.BlockSize > 0 {
			blockSize = param.BlockSize
		}
		postData.PartInfoList = aliyunpan.GenerateFileUploadPartInfoListWithChunkSize(param.Size, blockSize)
	}
	if postData.ContentHashName == "" {
		postData.ContentHashName = "sha1"
	}
	if postData.ParentFileId == "" {
		postData.ParentFileId = aliyunpan.DefaultRootParentFileId
	}
	if postData.ProofVersion == "" {
		postData.ProofVersion = "v1"
	}
	if postData.CheckNameMode == "" {
		postData.CheckNameMode = "auto_rename"
	}
	postData.Type = "file"

	// request
	resp, err := p.client.Req("POST", fullUrl.String(), postData, p.AddSignatureHeader(apiutil.AddCommonHeader(header)))
	if err != nil {
		logger.Verboseln("create upload file error ", err)
		return nil, apierror.NewFailedApiError(err.Error())
	}

	// handler common error
	body, err1 := apierror.ParseCommonResponseApiError(resp)
	if err1 != nil {
		return nil, err1
	}

	// parse result
	r := &aliyunpan.CreateFileUploadResult{}
	if err2 := json.Unmarshal(body, r); err2 != nil {
		logger.Verboseln("parse create upload file result json error ", err2)
		return nil, apierror.NewFailedApiError(err2.Error())
	}
	return r, nil
}

// GetUploadUrl 获取上传数据链接参数
// 因为有些文件过大，或者暂定上传后，然后过段时间再继续上传，这时候之前的上传链接可能已经失效了，所以需要重新获取上传数据的链接
// 如果该文件已经上传完毕，则该接口返回错误
func (p *WebPanClient) GetUploadUrl(param *aliyunpan.GetUploadUrlParam) (*aliyunpan.GetUploadUrlResult, *apierror.ApiError) {
	// header
	header := map[string]string{
		"authorization": p.webToken.GetAuthorizationStr(),
	}

	// url
	fullUrl := &strings.Builder{}
	fmt.Fprintf(fullUrl, "%s/v2/file/get_upload_url", API_URL)
	logger.Verboseln("do request url: " + fullUrl.String())

	// data
	postData := param

	// request
	body, err := p.client.Fetch("POST", fullUrl.String(), postData, p.AddSignatureHeader(apiutil.AddCommonHeader(header)))
	if err != nil {
		logger.Verboseln("get upload url error ", err)
		return nil, apierror.NewFailedApiError(err.Error())
	}

	// handler common error
	if err1 := apierror.ParseCommonApiError(body); err1 != nil {
		return nil, err1
	}

	// parse result
	r := &aliyunpan.GetUploadUrlResult{}
	if err2 := json.Unmarshal(body, r); err2 != nil {
		logger.Verboseln("parse get upload url result json error ", err2)
		return nil, apierror.NewFailedApiError(err2.Error())
	}
	r.CreateAt = apiutil.UtcTime2LocalFormat(r.CreateAt)
	return r, nil
}

// UploadFileData 上传文件数据
func (p *WebPanClient) UploadFileData(uploadUrl string, uploadFunc aliyunpan.UploadFunc) *apierror.ApiError {
	// header
	header := map[string]string{
		"referer": "https://www.aliyundrive.com/",
	}

	// url
	fullUrl := &strings.Builder{}
	fmt.Fprintf(fullUrl, "%s", uploadUrl)
	logger.Verboseln("do request url: " + fullUrl.String())

	// request
	if uploadFunc != nil {
		resp, err := uploadFunc("PUT", fullUrl.String(), header)
		if err != nil || (resp != nil && resp.StatusCode != 200) {
			logger.Verboseln("upload file data chunk error ", err)
			return apierror.NewFailedApiError("update data error")
		}
	}
	return nil
}

// UploadDataChunk 上传数据。该方法是同步阻塞的
func (p *WebPanClient) UploadDataChunk(url string, data *aliyunpan.FileUploadChunkData) *apierror.ApiError {
	var client = requester.NewHTTPClient()

	// header
	header := map[string]string{
		"referer": "https://www.aliyundrive.com/",
	}

	// url
	fullUrl := &strings.Builder{}
	fmt.Fprintf(fullUrl, "%s", url)
	logger.Verboseln("do request url: " + fullUrl.String())

	// data
	if data == nil || data.Reader == nil || data.Len() == 0 {
		return apierror.NewFailedApiError("数据块错误")
	}
	// request
	resp, err := client.Req("PUT", fullUrl.String(), data, header)
	if err != nil || resp.StatusCode != 200 {
		logger.Verboseln("upload file data chunk error ", err)
		return apierror.NewFailedApiError(err.Error())
	}
	return nil
}

// CompleteUploadFile 完成文件上传确认。完成文件数据上传后，需要调用该接口文件才会显示再网盘中
func (p *WebPanClient) CompleteUploadFile(param *aliyunpan.CompleteUploadFileParam) (*aliyunpan.CompleteUploadFileResult, *apierror.ApiError) {
	// header
	header := map[string]string{
		"authorization": p.webToken.GetAuthorizationStr(),
	}

	// url
	fullUrl := &strings.Builder{}
	fmt.Fprintf(fullUrl, "%s/v2/file/complete", API_URL)
	logger.Verboseln("do request url: " + fullUrl.String())

	// data
	postData := map[string]interface{}{
		"ignoreError": true,
		"drive_id":    param.DriveId,
		"file_id":     param.FileId,
		"upload_id":   param.UploadId,
	}

	// request
	body, err := p.client.Fetch("POST", fullUrl.String(), postData, p.AddSignatureHeader(apiutil.AddCommonHeader(header)))
	if err != nil {
		logger.Verboseln("complete upload file error ", err)
		return nil, apierror.NewFailedApiError(err.Error())
	}

	// handler common error
	if err1 := apierror.ParseCommonApiError(body); err1 != nil {
		return nil, err1
	}

	// parse result
	r := &completeUploadFileReqResult{}
	if err2 := json.Unmarshal(body, r); err2 != nil {
		logger.Verboseln("parse complete upload file result json error ", err2)
		return nil, apierror.NewFailedApiError(err2.Error())
	}

	return &aliyunpan.CompleteUploadFileResult{
		DriveId:         r.DriveId,
		DomainId:        r.DomainId,
		FileId:          r.FileId,
		Name:            r.Name,
		Type:            r.Type,
		Size:            r.Size,
		UploadId:        r.UploadId,
		ParentFileId:    r.ParentFileId,
		Crc64Hash:       r.Crc64Hash,
		ContentHash:     r.ContentHash,
		ContentHashName: r.ContentHashName,
		CreatedAt:       apiutil.UtcTime2LocalFormat(r.CreatedAt),
	}, nil
}
