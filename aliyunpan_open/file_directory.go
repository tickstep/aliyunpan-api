package aliyunpan_open

import (
	"github.com/tickstep/aliyunpan-api/aliyunpan"
	"github.com/tickstep/aliyunpan-api/aliyunpan/apierror"
	"github.com/tickstep/aliyunpan-api/aliyunpan/apiutil"
	"github.com/tickstep/aliyunpan-api/aliyunpan_open/openapi"
	"time"
)

func createFileEntity(f *openapi.FileItem) *aliyunpan.FileEntity {
	if f == nil {
		return nil
	}
	return &aliyunpan.FileEntity{
		DriveId:         f.DriveId,
		DomainId:        f.DomainId,
		FileId:          f.FileId,
		FileName:        f.Name,
		FileSize:        f.Size,
		FileType:        f.Type,
		CreatedAt:       apiutil.UtcTime2LocalFormat(f.CreatedAt),
		UpdatedAt:       apiutil.UtcTime2LocalFormat(f.UpdatedAt),
		FileExtension:   f.FileExtension,
		ParentFileId:    f.ParentFileId,
		ContentHash:     f.ContentHash,
		ContentHashName: f.ContentHashName,
		Path:            f.Name,
		Category:        f.Category,
	}
}

// FileList 获取文件列表
func (p *OpenPanClient) FileList(param *aliyunpan.FileListParam) (*aliyunpan.FileListResult, *apierror.ApiError) {
	result := &aliyunpan.FileListResult{
		FileList:   aliyunpan.FileList{},
		NextMarker: "",
	}

	opParam := &openapi.FileListParam{
		DriveId:        param.DriveId,
		ParentFileId:   param.ParentFileId,
		Limit:          param.Limit,
		Marker:         param.Marker,
		OrderBy:        string(param.OrderBy),
		OrderDirection: string(param.OrderDirection),
		Type:           "all",
		Fields:         "*",
	}
	if flr, err := p.apiClient.FileList(opParam); err == nil {
		for k := range flr.Items {
			if flr.Items[k] == nil {
				continue
			}
			result.FileList = append(result.FileList, createFileEntity(flr.Items[k]))
		}
		result.NextMarker = flr.NextMarker
	} else {
		// handle error, retry, token refresh
	}

	return result, nil
}

// FileListGetAll 获取指定目录下的所有文件列表
func (p *OpenPanClient) FileListGetAll(param *aliyunpan.FileListParam, delayMilliseconds int) (aliyunpan.FileList, *apierror.ApiError) {
	internalParam := &aliyunpan.FileListParam{
		OrderBy:        param.OrderBy,
		OrderDirection: param.OrderDirection,
		DriveId:        param.DriveId,
		ParentFileId:   param.ParentFileId,
		Limit:          param.Limit,
		Marker:         param.Marker,
	}
	if internalParam.Limit <= 0 {
		internalParam.Limit = 100
	}

	fileList := aliyunpan.FileList{}
	result, err := p.FileList(internalParam)
	if err != nil || result == nil {
		return nil, err
	}
	fileList = append(fileList, result.FileList...)

	// more page?
	for len(result.NextMarker) > 0 {
		if delayMilliseconds > 0 {
			time.Sleep(time.Duration(delayMilliseconds) * time.Millisecond)
		}
		internalParam.Marker = result.NextMarker
		result, err = p.FileList(internalParam)
		if err == nil && result != nil {
			fileList = append(fileList, result.FileList...)
		} else {
			return nil, err
		}
	}
	return fileList, nil
}
