package aliyunpan

import (
	"encoding/json"
	"fmt"
	"github.com/tickstep/aliyunpan-api/aliyunpan/apierror"
	"github.com/tickstep/aliyunpan-api/aliyunpan/apiutil"
	"github.com/tickstep/library-go/logger"
	"math"
	"strings"
)

type(

	// 上传文件分片参数。从1开始，最大为 10000
	FileUploadPartInfoParam struct {
		PartNumber        int    `json:"part_number"`
	}

	// 创建上传文件参数
	CreateFileUploadParam struct {
		Name            string `json:"name"`
		DriveId         string `json:"drive_id"`
		ParentFileId    string `json:"parent_file_id"`
		Size            int64    `json:"size"`
		// 上传文件分片参数，最大为 10000
		PartInfoList    []FileUploadPartInfoParam `json:"part_info_list"`
		ContentHash     string `json:"content_hash"`
		// 默认为 sha1
		ContentHashName string `json:"content_hash_name"`
		// 默认为 file
		Type            string `json:"type"`
		// 默认为 auto_rename
		CheckNameMode   string `json:"check_name_mode"`
	}

	FileUploadPartInfoResult struct {
		PartNumber        int    `json:"part_number"`
		UploadURL         string `json:"upload_url"`
		InternalUploadURL string `json:"internal_upload_url"`
		ContentType       string `json:"content_type"`
	}

	// 创建上传文件返回值
	CreateFileUploadResult struct {
		ParentFileId string `json:"parent_file_id"`
		PartInfoList []FileUploadPartInfoResult `json:"part_info_list"`
		UploadId    string `json:"upload_id"`
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
)

const(
	// 默认分片大小，512KB
	DefaultChunkSize = int64(524288)

	// 最大分片数量大小
	MaxPartNum = 10000
)

// GenerateFileUploadPartInfoList 根据文件大小自动生成分片
func GenerateFileUploadPartInfoList(size int64) []FileUploadPartInfoParam {
	r := []FileUploadPartInfoParam{}
	if size <= DefaultChunkSize {
		r = append(r, FileUploadPartInfoParam{
			PartNumber: 1,
		})
	} else {
		pageSize := int(math.Ceil(float64(size) / float64(DefaultChunkSize)))
		for i := 1; i <= pageSize; i++ {
			r = append(r, FileUploadPartInfoParam{
				PartNumber: i,
			})
		}
	}
	return r
}

// CreateUploadFile 创建上传文件，如果文件已经上传过则会直接秒传
func (p *PanClient) CreateUploadFile(param *CreateFileUploadParam) (*CreateFileUploadResult, *apierror.ApiError) {
	// header
	header := map[string]string {
		"authorization": p.webToken.GetAuthorizationStr(),
	}

	// url
	fullUrl := &strings.Builder{}
	fmt.Fprintf(fullUrl, "%s/v2/file/create", API_URL)
	logger.Verboseln("do request url: " + fullUrl.String())

	// data
	postData := param
	if postData.ContentHashName == "" {
		postData.ContentHashName = "sha1"
	}
	if postData.ParentFileId == "" {
		postData.ParentFileId = DefaultRootParentFileId
	}
	postData.Type = "file"
	postData.CheckNameMode = "auto_rename"

	// request
	body, err := client.Fetch("POST", fullUrl.String(), postData, apiutil.AddCommonHeader(header))
	if err != nil {
		logger.Verboseln("create upload file error ", err)
		return nil, apierror.NewFailedApiError(err.Error())
	}

	// handler common error
	if err1 := apierror.ParseCommonApiError(body); err1 != nil {
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