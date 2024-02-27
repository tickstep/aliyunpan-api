package aliyunpan_open

import (
	"github.com/tickstep/aliyunpan-api/aliyunpan"
	"github.com/tickstep/aliyunpan-api/aliyunpan/apierror"
	"github.com/tickstep/aliyunpan-api/aliyunpan/apiutil"
	"github.com/tickstep/aliyunpan-api/aliyunpan_open/openapi"
	"strings"
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
	retryTime := 0

RetryBegin:
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
		// handle common error
		if apiErrorHandleResp := p.HandleAliApiError(err, &retryTime); apiErrorHandleResp.NeedRetry {
			goto RetryBegin
		} else {
			return nil, apiErrorHandleResp.ApiErr
		}
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

// FileInfoById 通过FileId获取文件信息
func (p *OpenPanClient) FileInfoById(driveId, fileId string) (*aliyunpan.FileEntity, *apierror.ApiError) {
	retryTime := 0

RetryBegin:
	opParam := &openapi.FileIdentityPair{
		DriveId: driveId,
		FileId:  fileId,
	}
	if result, err := p.apiClient.FileGetDetailInfo(opParam); err == nil {
		return createFileEntity(result), nil
	} else {
		// handle common error
		if apiErrorHandleResp := p.HandleAliApiError(err, &retryTime); apiErrorHandleResp.NeedRetry {
			goto RetryBegin
		} else {
			return nil, apiErrorHandleResp.ApiErr
		}
	}
}

// FileInfoByPath 通过路径获取文件详情，pathStr是绝对路径
func (p *OpenPanClient) FileInfoByPath(driveId string, pathStr string) (fileInfo *aliyunpan.FileEntity, error *apierror.ApiError) {
	retryTime := 0

RetryBegin:
	opParam := &openapi.FilePathPair{
		DriveId:  driveId,
		FilePath: pathStr,
	}
	if result, err := p.apiClient.FileGetDetailInfoByPath(opParam); err == nil {
		return createFileEntity(result), nil
	} else {
		// handle common error
		if apiErrorHandleResp := p.HandleAliApiError(err, &retryTime); apiErrorHandleResp.NeedRetry {
			goto RetryBegin
		} else {
			return nil, apiErrorHandleResp.ApiErr
		}
	}
}

// FilesDirectoriesRecurseList 递归获取目录下的文件和目录列表
func (p *OpenPanClient) FilesDirectoriesRecurseList(driveId string, path string, handleFileDirectoryFunc aliyunpan.HandleFileDirectoryFunc) aliyunpan.FileList {
	targetFileInfo, er := p.FileInfoByPath(driveId, path)
	if er != nil {
		if handleFileDirectoryFunc != nil {
			handleFileDirectoryFunc(0, path, nil, er)
		}
		return nil
	}
	if targetFileInfo.IsFolder() {
		// folder
		if handleFileDirectoryFunc != nil {
			handleFileDirectoryFunc(0, path, targetFileInfo, nil)
		}
	} else {
		// file
		if handleFileDirectoryFunc != nil {
			handleFileDirectoryFunc(0, path, targetFileInfo, nil)
		}
		return aliyunpan.FileList{targetFileInfo}
	}

	fld := &aliyunpan.FileList{}
	ok := p.recurseList(driveId, targetFileInfo, 1, handleFileDirectoryFunc, fld)
	if !ok {
		return nil
	}
	return *fld
}

func (p *OpenPanClient) recurseList(driveId string, folderInfo *aliyunpan.FileEntity, depth int, handleFileDirectoryFunc aliyunpan.HandleFileDirectoryFunc, fld *aliyunpan.FileList) bool {
	flp := &aliyunpan.FileListParam{
		DriveId:      driveId,
		ParentFileId: folderInfo.FileId,
	}
	r, apiError := p.FileListGetAll(flp, 0)
	if apiError != nil {
		if handleFileDirectoryFunc != nil {
			handleFileDirectoryFunc(depth, folderInfo.Path, nil, apiError)
		}
		return false
	}
	ok := true
	for _, fi := range r {
		fi.Path = strings.ReplaceAll(folderInfo.Path+aliyunpan.PathSeparator+fi.FileName, "//", "/")
		*fld = append(*fld, fi)
		if fi.IsFolder() {
			if handleFileDirectoryFunc != nil {
				ok = handleFileDirectoryFunc(depth, fi.Path, fi, nil)
			}
			ok = p.recurseList(driveId, fi, depth+1, handleFileDirectoryFunc, fld)
		} else {
			if handleFileDirectoryFunc != nil {
				ok = handleFileDirectoryFunc(depth, fi.Path, fi, nil)
			}
		}
		if !ok {
			return false
		}
	}
	return true
}
