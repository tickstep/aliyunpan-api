package aliyunpan_open

import (
	"fmt"
	"github.com/tickstep/aliyunpan-api/aliyunpan"
	"github.com/tickstep/aliyunpan-api/aliyunpan/apierror"
	"github.com/tickstep/aliyunpan-api/aliyunpan_open/openapi"
	"github.com/tickstep/library-go/logger"
	"strings"
)

// CheckUploadFilePreHash 文件PreHash检测，当PreHash检查为false的文件肯定不支持秒传
func (p *OpenPanClient) CheckUploadFilePreHash(param *aliyunpan.FileUploadCheckPreHashParam) (bool, *apierror.ApiError) {
	retryTime := 0

RetryBegin:
	opParam := &openapi.FileUploadCheckPreHashParam{
		DriveId:       param.DriveId,
		ParentFileId:  param.ParentFileId,
		Name:          param.Name,
		Type:          "file",
		CheckNameMode: "ignore",
		Size:          param.Size,
		PreHash:       param.PreHash,
	}
	if result, err := p.apiClient.FileUploadCheckPreHash(opParam); err == nil {
		return result, nil
	} else {
		// handle common error
		if apiErrorHandleResp := p.HandleAliApiError(err, &retryTime); apiErrorHandleResp.NeedRetry {
			goto RetryBegin
		} else {
			return false, apiErrorHandleResp.ApiErr
		}
	}
}

// CreateUploadFile 创建上传文件，如果文件已经上传过可以直接秒传
func (p *OpenPanClient) CreateUploadFile(param *aliyunpan.CreateFileUploadParam) (*aliyunpan.CreateFileUploadResult, *apierror.ApiError) {
	retryTime := 0

	// 计算分片数量
	partInfoListParam := param.PartInfoList
	if len(param.PartInfoList) == 0 {
		blockSize := aliyunpan.DefaultChunkSize
		if param.BlockSize > 0 {
			blockSize = param.BlockSize
		}
		partInfoListParam = aliyunpan.GenerateFileUploadPartInfoListWithChunkSize(param.Size, blockSize)
	}
	realPartInfoList := []*openapi.PartInfoItem{}
	for _, v := range partInfoListParam {
		realPartInfoList = append(realPartInfoList, &openapi.PartInfoItem{
			PartNumber: v.PartNumber,
		})
	}

RetryBegin:
	opParam := &openapi.FileUploadCreateParam{
		DriveId:         param.DriveId,
		ParentFileId:    param.ParentFileId,
		Name:            param.Name,
		Type:            "file",
		CheckNameMode:   param.CheckNameMode,
		Size:            param.Size,
		PartInfoList:    realPartInfoList,
		ContentHash:     param.ContentHash,
		ContentHashName: param.ContentHashName,
		ProofCode:       param.ProofCode,
		ProofVersion:    param.ProofVersion,
		LocalCreatedAt:  param.LocalCreatedAt,
		LocalModifiedAt: param.LocalModifiedAt,
	}
	if opParam.ContentHashName == "" {
		opParam.ContentHashName = "sha1"
	}
	if opParam.ParentFileId == "" {
		opParam.ParentFileId = aliyunpan.DefaultRootParentFileId
	}
	if opParam.ProofVersion == "" {
		opParam.ProofVersion = "v1"
	}
	if opParam.CheckNameMode == "" {
		opParam.CheckNameMode = "auto_rename"
	}

	if result, err := p.apiClient.FileUploadCreate(opParam); err == nil {
		partInfoListResult := []aliyunpan.FileUploadPartInfoResult{}
		for _, v := range result.PartInfoList {
			partInfoListResult = append(partInfoListResult, aliyunpan.FileUploadPartInfoResult{
				PartNumber:        v.PartNumber,
				UploadURL:         v.UploadUrl,
				InternalUploadURL: "",
				ContentType:       "",
			})
		}
		return &aliyunpan.CreateFileUploadResult{
			ParentFileId: result.ParentFileId,
			PartInfoList: partInfoListResult,
			UploadId:     result.UploadId,
			RapidUpload:  result.RapidUpload,
			Type:         "",
			FileId:       result.FileId,
			DomainId:     "",
			DriveId:      result.DriveId,
			FileName:     result.FileName,
			EncryptMode:  "",
			Location:     "",
		}, nil
	} else {
		// handle common error
		if apiErrorHandleResp := p.HandleAliApiError(err, &retryTime); apiErrorHandleResp.NeedRetry {
			goto RetryBegin
		} else {
			return nil, apiErrorHandleResp.ApiErr
		}
	}
}

// GetUploadUrl 获取上传数据链接参数
// 因为有些文件过大，或者暂定上传后，然后过段时间再继续上传，这时候之前的上传链接可能已经失效了，所以需要重新获取上传数据的链接
func (p *OpenPanClient) GetUploadUrl(param *aliyunpan.GetUploadUrlParam) (*aliyunpan.GetUploadUrlResult, *apierror.ApiError) {
	retryTime := 0

	realPartInfoList := []*openapi.PartInfoItem{}
	for _, v := range param.PartInfoList {
		realPartInfoList = append(realPartInfoList, &openapi.PartInfoItem{
			PartNumber: v.PartNumber,
		})
	}
RetryBegin:
	opParam := &openapi.FileUploadGetUploadUrlParam{
		DriveId:      param.DriveId,
		FileId:       param.FileId,
		UploadId:     param.UploadId,
		PartInfoList: realPartInfoList,
	}
	if result, err := p.apiClient.FileUploadGetUploadUrl(opParam); err == nil {
		partInfoListResult := []aliyunpan.FileUploadPartInfoResult{}
		for _, v := range result.PartInfoList {
			partInfoListResult = append(partInfoListResult, aliyunpan.FileUploadPartInfoResult{
				PartNumber:        v.PartNumber,
				UploadURL:         v.UploadUrl,
				InternalUploadURL: "",
				ContentType:       "",
			})
		}
		return &aliyunpan.GetUploadUrlResult{
			DomainId:     "",
			DriveId:      result.DriveId,
			FileId:       result.FileId,
			PartInfoList: partInfoListResult,
			UploadId:     result.UploadId,
			CreateAt:     result.CreatedAt,
		}, nil
	} else {
		// handle common error
		if apiErrorHandleResp := p.HandleAliApiError(err, &retryTime); apiErrorHandleResp.NeedRetry {
			goto RetryBegin
		} else {
			return nil, apiErrorHandleResp.ApiErr
		}
	}
}

// UploadFileData 上传文件数据
func (p *OpenPanClient) UploadFileData(uploadUrl string, uploadFunc aliyunpan.UploadFunc) *apierror.ApiError {
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

// CompleteUploadFile 完成文件上传确认。完成文件数据上传后，需要调用该接口文件才会显示再网盘中
func (p *OpenPanClient) CompleteUploadFile(param *aliyunpan.CompleteUploadFileParam) (*aliyunpan.CompleteUploadFileResult, *apierror.ApiError) {
	retryTime := 0

RetryBegin:
	opParam := &openapi.FileUploadCompleteParam{
		DriveId:  param.DriveId,
		FileId:   param.FileId,
		UploadId: param.UploadId,
	}
	if result, err := p.apiClient.FileUploadComplete(opParam); err == nil {
		return &aliyunpan.CompleteUploadFileResult{
			DriveId:         result.DriveId,
			DomainId:        "",
			FileId:          result.FileId,
			Name:            result.Name,
			Type:            result.Type,
			Size:            result.Size,
			UploadId:        param.UploadId,
			ParentFileId:    result.ParentFileId,
			Crc64Hash:       "",
			ContentHash:     result.ContentHash,
			ContentHashName: result.ContentHashName,
			CreatedAt:       result.LocalCreatedAt,
		}, nil
	} else {
		// handle common error
		if apiErrorHandleResp := p.HandleAliApiError(err, &retryTime); apiErrorHandleResp.NeedRetry {
			goto RetryBegin
		} else {
			return nil, apiErrorHandleResp.ApiErr
		}
	}
}
