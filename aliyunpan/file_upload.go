package aliyunpan

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/tickstep/aliyunpan-api/aliyunpan/apierror"
	"github.com/tickstep/aliyunpan-api/aliyunpan/apiutil"
	"github.com/tickstep/library-go/logger"
	"github.com/tickstep/library-go/requester"
	"github.com/tickstep/library-go/requester/rio"
	"io"
	"math"
	"math/big"
	"net/http"
	"strings"
)

type (
	// UploadFunc 上传文件处理函数
	UploadFunc func(httpMethod, fullUrl string, headers map[string]string) (resp *http.Response, err error)

	// 上传文件分片参数。从1开始，最大为 10000
	FileUploadPartInfoParam struct {
		PartNumber int `json:"part_number"`
	}

	// 创建上传文件参数
	CreateFileUploadParam struct {
		Name         string `json:"name"`
		DriveId      string `json:"drive_id"`
		ParentFileId string `json:"parent_file_id"`
		Size         int64  `json:"size"`
		// 上传文件分片参数，最大为 10000
		PartInfoList []FileUploadPartInfoParam `json:"part_info_list"`
		ContentHash  string                    `json:"content_hash"`
		// 默认为 sha1。可选：sha1，none
		ContentHashName string `json:"content_hash_name"`
		// 默认为 file
		Type string `json:"type"`
		// 默认为 auto_rename。可选：overwrite-覆盖网盘同名文件，auto_rename-自动重命名，refuse-无需检测
		CheckNameMode string `json:"check_name_mode"`

		ProofCode    string `json:"proof_code"`
		ProofVersion string `json:"proof_version"`

		// 分片大小
		// 不进行json序列化
		BlockSize int64 `json:"-"`
	}

	FileUploadPartInfoResult struct {
		PartNumber        int    `json:"part_number"`
		UploadURL         string `json:"upload_url"`
		InternalUploadURL string `json:"internal_upload_url"`
		ContentType       string `json:"content_type"`
	}

	// 创建上传文件返回值
	CreateFileUploadResult struct {
		ParentFileId string                     `json:"parent_file_id"`
		PartInfoList []FileUploadPartInfoResult `json:"part_info_list"`
		UploadId     string                     `json:"upload_id"`
		// RapidUpload 是否秒传。true-已秒传，false-没有秒传，需要手动上传
		RapidUpload bool   `json:"rapid_upload"`
		Type        string `json:"type"`
		FileId      string `json:"file_id"`
		DomainId    string `json:"domain_id"`
		DriveId     string `json:"drive_id"`
		// FileName 保存在网盘的名称，因为网盘会自动重命名同名的文件
		FileName    string `json:"file_name"`
		EncryptMode string `json:"encrypt_mode"`
		Location    string `json:"location"`
	}

	// 获取上传数据链接参数
	GetUploadUrlParam struct {
		DriveId      string                    `json:"drive_id"`
		FileId       string                    `json:"file_id"`
		PartInfoList []FileUploadPartInfoParam `json:"part_info_list"`
		UploadId     string                    `json:"upload_id"`
	}

	// 获取上传数据链接返回值
	GetUploadUrlResult struct {
		DomainId     string                     `json:"domain_id"`
		DriveId      string                     `json:"drive_id"`
		FileId       string                     `json:"file_id"`
		PartInfoList []FileUploadPartInfoResult `json:"part_info_list"`
		UploadId     string                     `json:"upload_id"`
		CreateAt     string                     `json:"create_at"`
	}

	FileUploadRange struct {
		// 起始值，包含
		Offset int64
		// 总上传长度
		Len int64
	}

	// 文件上传数据块
	FileUploadChunkData struct {
		Reader    io.Reader
		ChunkSize int64

		hasReadCount int64
	}

	// 提交上传文件传输完成参数
	CompleteUploadFileParam struct {
		DriveId  string `json:"drive_id"`
		FileId   string `json:"file_id"`
		UploadId string `json:"upload_id"`
	}

	CompleteUploadFileResult struct {
		DriveId         string `json:"drive_id"`
		DomainId        string `json:"domain_id"`
		FileId          string `json:"file_id"`
		Name            string `json:"name"`
		Type            string `json:"type"`
		Size            int64  `json:"size"`
		UploadId        string `json:"upload_id"`
		ParentFileId    string `json:"parent_file_id"`
		Crc64Hash       string `json:"crc64_hash"`
		ContentHash     string `json:"content_hash"`
		ContentHashName string `json:"content_hash_name"`
		CreatedAt       string `json:"created_at"`
	}

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

const (
	// 默认分片大小，512KB
	DefaultChunkSize = int64(524288)

	// 最大分片数量大小
	MaxPartNum = 10000

	// 0KB文件默认的SHA1哈希值
	DefaultZeroSizeFileContentHash = "DA39A3EE5E6B4B0D3255BFEF95601890AFD80709"
)

func (d *FileUploadChunkData) Read(p []byte) (n int, err error) {
	realReadCount := int64(0)
	var buf []byte = p
	needCopy := false
	if (d.hasReadCount + int64(len(p))) > d.ChunkSize {
		realReadCount = d.ChunkSize - d.hasReadCount
		buf = make([]byte, realReadCount)
		needCopy = true
	}

	n, err = d.Reader.Read(buf)
	if needCopy {
		copy(p, buf)
	}
	d.hasReadCount += int64(n)
	return n, err
}

func (d *FileUploadChunkData) Len() int64 {
	return d.ChunkSize
}

// GenerateFileUploadPartInfoList 根据文件大小自动生成分片
func GenerateFileUploadPartInfoList(size int64) []FileUploadPartInfoParam {
	return GenerateFileUploadPartInfoListWithChunkSize(size, DefaultChunkSize)
}

// GenerateFileUploadPartInfoList 根据文件大小和指定的分片大小自动生成分片
func GenerateFileUploadPartInfoListWithChunkSize(size, chunkSize int64) []FileUploadPartInfoParam {
	r := []FileUploadPartInfoParam{}
	if size <= chunkSize {
		r = append(r, FileUploadPartInfoParam{
			PartNumber: 1,
		})
	} else {
		pageSize := int(math.Ceil(float64(size) / float64(chunkSize)))
		for i := 1; i <= pageSize; i++ {
			r = append(r, FileUploadPartInfoParam{
				PartNumber: i,
			})
		}
	}
	return r
}

// CalcProofCode 计算文件上传防伪码
func CalcProofCode(accessToken string, reader rio.ReaderAtLen64, fileSize int64) string {
	if fileSize == 0 { // empty file
		return ""
	}

	md5w := md5.New()
	md5w.Write([]byte(accessToken))
	md5bytes := md5w.Sum(nil)
	hashCode := hex.EncodeToString(md5bytes)[0:16]
	hashInteger, _ := new(big.Int).SetString(hashCode, 16)

	z := big.NewInt(0)
	startPosInteger := big.NewInt(0)
	z.Div(hashInteger, big.NewInt(fileSize))
	startPosInteger.Sub(hashInteger, big.NewInt(z.Int64()*fileSize))
	startPos := startPosInteger.Int64()

	endPos := startPos + 8
	if endPos > fileSize {
		endPos = fileSize
	}

	// read byte from file
	readCount := endPos - startPos
	proofBytes := make([]byte, readCount)
	reader.ReadAt(proofBytes, startPos)

	// calc the base64 string for read bytes
	return base64.StdEncoding.EncodeToString(proofBytes)
}

// CreateUploadFile 创建上传文件，如果文件已经上传过则会直接秒传
func (p *PanClient) CreateUploadFile(param *CreateFileUploadParam) (*CreateFileUploadResult, *apierror.ApiError) {
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
		blockSize := DefaultChunkSize
		if param.BlockSize > 0 {
			blockSize = param.BlockSize
		}
		postData.PartInfoList = GenerateFileUploadPartInfoListWithChunkSize(param.Size, blockSize)
	}
	if postData.ContentHashName == "" {
		postData.ContentHashName = "sha1"
	}
	if postData.ParentFileId == "" {
		postData.ParentFileId = DefaultRootParentFileId
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
	r := &CreateFileUploadResult{}
	if err2 := json.Unmarshal(body, r); err2 != nil {
		logger.Verboseln("parse create upload file result json error ", err2)
		return nil, apierror.NewFailedApiError(err2.Error())
	}
	return r, nil
}

// GetUploadUrl 获取上传数据链接参数
// 因为有些文件过大，或者暂定上传后，然后过段时间再继续上传，这时候之前的上传链接可能已经失效了，所以需要重新获取上传数据的链接
// 如果该文件已经上传完毕，则该接口返回错误
func (p *PanClient) GetUploadUrl(param *GetUploadUrlParam) (*GetUploadUrlResult, *apierror.ApiError) {
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
	r := &GetUploadUrlResult{}
	if err2 := json.Unmarshal(body, r); err2 != nil {
		logger.Verboseln("parse get upload url result json error ", err2)
		return nil, apierror.NewFailedApiError(err2.Error())
	}
	r.CreateAt = apiutil.UtcTime2LocalFormat(r.CreateAt)
	return r, nil
}

// UploadFileData 上传文件数据
func (p *PanClient) UploadFileData(uploadUrl string, uploadFunc UploadFunc) *apierror.ApiError {
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
func (p *PanClient) UploadDataChunk(url string, data *FileUploadChunkData) *apierror.ApiError {
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
func (p *PanClient) CompleteUploadFile(param *CompleteUploadFileParam) (*CompleteUploadFileResult, *apierror.ApiError) {
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

	return &CompleteUploadFileResult{
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
