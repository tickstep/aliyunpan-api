package aliyunpan_open

import (
	"github.com/tickstep/aliyunpan-api/aliyunpan"
	"github.com/tickstep/aliyunpan-api/aliyunpan/apierror"
	"github.com/tickstep/aliyunpan-api/aliyunpan/apiutil"
	"github.com/tickstep/aliyunpan-api/aliyunpan_open/openapi"
	"strings"
)

// Mkdir 创建文件夹
func (p *OpenPanClient) Mkdir(driveId, parentFileId, dirName string) (*aliyunpan.MkdirResult, *apierror.ApiError) {
	retryTime := 0
	if !apiutil.CheckFileNameValid(dirName) {
		return nil, apierror.NewFailedApiError("文件夹名不能包含特殊字符：" + apiutil.FileNameSpecialChars)
	}

RetryBegin:
	opParam := &openapi.FileUploadCreateParam{
		DriveId:       driveId,
		ParentFileId:  parentFileId,
		Name:          dirName,
		Type:          "folder",
		CheckNameMode: "auto_rename",
	}
	if result, err := p.apiClient.FileUploadCreate(opParam); err == nil {
		return &aliyunpan.MkdirResult{
			ParentFileId: result.ParentFileId,
			Type:         "folder",
			FileId:       result.FileId,
			DriveId:      result.DriveId,
			FileName:     result.FileName,
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

// MkdirByFullPath 通过绝对路径创建文件夹
func (p *OpenPanClient) MkdirByFullPath(driveId, fullPath string) (*aliyunpan.MkdirResult, *apierror.ApiError) {
	fullPath = strings.ReplaceAll(fullPath, "//", "/")
	fullPath = strings.Trim(fullPath, " ")
	if fullPath == "/" {
		return &aliyunpan.MkdirResult{
			ParentFileId: "",
			Type:         "folder",
			FileId:       "root",
			DomainId:     "",
			DriveId:      "",
			FileName:     "",
			EncryptMode:  "",
		}, nil
	}
	pathSlice := strings.Split(fullPath, "/")
	return p.MkdirRecursive(driveId, "", "", 0, pathSlice)
}

func (p *OpenPanClient) MkdirRecursive(driveId, parentFileId string, fullPath string, index int, pathSlice []string) (*aliyunpan.MkdirResult, *apierror.ApiError) {
	r := &aliyunpan.MkdirResult{}
	if parentFileId == "" {
		// default root "/" entity
		parentFileId = "root"
		if index == 0 && len(pathSlice) == 1 {
			// root path "/"
			r.FileId = parentFileId
			return r, nil
		}

		fullPath = ""
		return p.MkdirRecursive(driveId, parentFileId, fullPath, index+1, pathSlice)
	}
	if index >= len(pathSlice) {
		r.ParentFileId = "root"
		r.FileId = parentFileId
		r.Type = "folder"
		r.FileName = pathSlice[index-1]
		return r, nil
	}

	// existed?
	thisDirPath := fullPath + "/" + pathSlice[index]
	fileEntity, e := p.FileInfoByPath(driveId, thisDirPath)
	if e != nil && e.Code != apierror.ApiCodeFileNotFoundCode {
		return nil, e
	}
	if fileEntity != nil {
		// existed
		if fileEntity.IsFile() {
			return nil, apierror.NewFailedApiError("the fileName is a file not a folder")
		}
		return p.MkdirRecursive(driveId, fileEntity.FileId, fullPath+"/"+pathSlice[index], index+1, pathSlice)
	}

	// not existed, mkdir dir
	name := pathSlice[index]
	if !apiutil.CheckFileNameValid(name) {
		r.FileId = ""
		return r, apierror.NewFailedApiError("文件夹名不能包含特殊字符：" + apiutil.FileNameSpecialChars)
	}

	rs, err := p.Mkdir(driveId, parentFileId, name)
	if err != nil {
		r.FileId = ""
		return r, err
	}

	if (index + 1) >= len(pathSlice) {
		return rs, nil
	} else {
		return p.MkdirRecursive(driveId, rs.FileId, fullPath+"/"+pathSlice[index], index+1, pathSlice)
	}
}
