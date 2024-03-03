// Copyright (c) 2020 tickstep.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package aliyunpan_web

import (
	"encoding/json"
	"fmt"
	"github.com/tickstep/aliyunpan-api/aliyunpan"
	"github.com/tickstep/aliyunpan-api/aliyunpan/apierror"
	"github.com/tickstep/aliyunpan-api/aliyunpan/apiutil"
	"github.com/tickstep/library-go/logger"
	"path"
	"strings"
	"time"
)

type (
	fileEntityResult struct {
		DriveId         string `json:"drive_id"`
		DomainId        string `json:"domain_id"`
		FileId          string `json:"file_id"`
		Name            string `json:"name"`
		Type            string `json:"type"`
		ContentType     string `json:"content_type"`
		CreatedAt       string `json:"created_at"`
		UpdatedAt       string `json:"updated_at"`
		FileExtension   string `json:"file_extension"`
		MimeType        string `json:"mime_type"`
		MimeExtension   string `json:"mime_extension"`
		Hidden          bool   `json:"hidden"`
		Size            int64  `json:"size"`
		Starred         bool   `json:"starred"`
		Status          string `json:"status"`
		UploadId        string `json:"upload_id"`
		ParentFileId    string `json:"parent_file_id"`
		Crc64Hash       string `json:"crc64_hash"`
		ContentHash     string `json:"content_hash"`
		ContentHashName string `json:"content_hash_name"`
		DownloadUrl     string `json:"download_Url"`
		Url             string `json:"url"`
		Category        string `json:"category"`
		EncryptMode     string `json:"encrypt_mode"`
		PunishFlag      int    `json:"punish_flag"`
		SyncFlag        bool   `json:"sync_flag"`
		SyncMeta        string `json:"sync_meta"`
	}

	fileListResult struct {
		Items []*fileEntityResult `json:"items"`
		// NextMarker 不为空，说明还有下一页
		NextMarker string `json:"next_marker"`
	}
)

func createFileEntity(f *fileEntityResult) *aliyunpan.FileEntity {
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
		UploadId:        f.UploadId,
		ParentFileId:    f.ParentFileId,
		Crc64Hash:       f.Crc64Hash,
		ContentHash:     f.ContentHash,
		ContentHashName: f.ContentHashName,
		Path:            f.Name,
		Category:        f.Category,
		SyncFlag:        f.SyncFlag,
		SyncMeta:        f.SyncMeta,
	}
}

// FileList 获取文件列表
func (p *WebPanClient) FileList(param *aliyunpan.FileListParam) (*aliyunpan.FileListResult, *apierror.ApiError) {
	result := &aliyunpan.FileListResult{
		FileList:   aliyunpan.FileList{},
		NextMarker: "",
	}
	retryCount := int64(1)
retry:
	if flr, err := p.fileListReq(param); err == nil {
		for k := range flr.Items {
			if flr.Items[k] == nil {
				continue
			}

			result.FileList = append(result.FileList, createFileEntity(flr.Items[k]))
		}
		result.NextMarker = flr.NextMarker
	} else {
		if err.Code == apierror.ApiCodeTooManyRequests {
			if retryCount <= aliyunpan.MaxRequestRetryCount {
				logger.Verboseln("too many request error, sleep and retry later")
				time.Sleep(time.Duration(retryCount*2) * time.Second)
				retryCount++
				goto retry
			}
		} else if err.Code == apierror.ApiCodeDeviceSessionSignatureInvalid {
			logger.Verboseln("device session signature invalid, updating new session signature")
			time.Sleep(time.Duration(2 * time.Second))
			if r, e := p.CreateSession(nil); e != nil {
				logger.Verboseln("update session signature error")
				logger.Verboseln(r)
			} else {
				logger.Verboseln("update session signature success")
				goto retry
			}
		}
		return nil, err
	}
	return result, nil
}

func (p *WebPanClient) fileListReq(param *aliyunpan.FileListParam) (*fileListResult, *apierror.ApiError) {
	header := map[string]string{
		"authorization": p.webToken.GetAuthorizationStr(),
	}

	fullUrl := &strings.Builder{}
	fmt.Fprintf(fullUrl, "%s/adrive/v3/file/list", API_URL)
	logger.Verboseln("do request url: " + fullUrl.String())

	pFileId := param.ParentFileId
	if pFileId == "" {
		pFileId = aliyunpan.DefaultRootParentFileId
	}
	limit := param.Limit
	if limit <= 0 {
		limit = 100
	}
	if param.OrderBy == "" {
		param.OrderBy = aliyunpan.FileOrderByUpdatedAt
	}
	if param.OrderDirection == "" {
		param.OrderDirection = aliyunpan.FileOrderDirectionDesc
	}
	postData := map[string]interface{}{
		"drive_id":                param.DriveId,
		"parent_file_id":          pFileId,
		"limit":                   limit,
		"all":                     false,
		"url_expire_sec":          1600,
		"image_thumbnail_process": "image/resize,w_400/format,jpeg",
		"image_url_process":       "image/resize,w_1920/format,jpeg",
		"video_thumbnail_process": "video/snapshot,t_0,f_jpg,ar_auto,w_800",
		"fields":                  "*",
		"order_by":                param.OrderBy,
		"order_direction":         param.OrderDirection,
	}
	if len(param.Marker) > 0 {
		postData["marker"] = param.Marker
	}

	// request
	resp, err := p.client.Req("POST", fullUrl.String(), postData, p.AddSignatureHeader(apiutil.AddCommonHeader(header)))
	//logger.Verboseln("get file list response: ", string(body))
	if err != nil {
		logger.Verboseln("get file list error ", err)
		return nil, apierror.NewApiErrorWithError(err)
	}

	// handler common error
	body, err1 := apierror.ParseCommonResponseApiError(resp)
	if err1 != nil {
		return nil, err1
	}

	// parse result
	r := &fileListResult{}
	if err2 := json.Unmarshal(body, r); err2 != nil {
		logger.Verboseln("parse file list result json error ", err2)
		return nil, apierror.NewFailedApiError(err2.Error())
	}
	return r, nil
}

// FileInfoById 通过FileId获取文件信息
func (p *WebPanClient) FileInfoById(driveId, fileId string) (*aliyunpan.FileEntity, *apierror.ApiError) {
	header := map[string]string{
		"authorization": p.webToken.GetAuthorizationStr(),
	}

	fullUrl := &strings.Builder{}
	fmt.Fprintf(fullUrl, "%s/v2/file/get", API_URL)
	logger.Verboseln("do request url: " + fullUrl.String())

	pFileId := fileId
	if pFileId == "" {
		pFileId = aliyunpan.DefaultRootParentFileId
	}
	postData := map[string]interface{}{
		"drive_id": driveId,
		"file_id":  pFileId,
	}

	// request
	body, err := p.client.Fetch("POST", fullUrl.String(), postData, p.AddSignatureHeader(apiutil.AddCommonHeader(header)))
	if err != nil {
		logger.Verboseln("get file info error ", err)
		return nil, apierror.NewFailedApiError(err.Error())
	}

	// handler common error
	if err1 := apierror.ParseCommonApiError(body); err1 != nil {
		return nil, err1
	}

	// parse result
	r := &fileEntityResult{}
	if err2 := json.Unmarshal(body, r); err2 != nil {
		logger.Verboseln("parse file info result json error ", err2)
		return nil, apierror.NewFailedApiError(err2.Error())
	}
	return createFileEntity(r), nil
}

// FileInfoByPath 通过路径获取文件详情，pathStr是绝对路径
func (p *WebPanClient) FileInfoByPath(driveId string, pathStr string) (fileInfo *aliyunpan.FileEntity, error *apierror.ApiError) {
	if pathStr == "" {
		pathStr = "/"
	}
	//pathStr = path.Clean(pathStr)
	if !path.IsAbs(pathStr) {
		return nil, apierror.NewFailedApiError("pathStr必须是绝对路径")
	}
	if len(pathStr) > 1 {
		pathStr = path.Clean(pathStr)
	}

	// try cache
	if v := p.loadFilePathFromCache(driveId, pathStr); v != nil {
		return v, nil
	}

	var pathSlice []string
	if pathStr == "/" {
		pathSlice = []string{""}
	} else {
		pathSlice = strings.Split(pathStr, aliyunpan.PathSeparator)
		if pathSlice[0] != "" {
			return nil, apierror.NewFailedApiError("pathStr必须是绝对路径")
		}
	}
	fileInfo, error = p.getFileInfoByPath(driveId, 0, &pathSlice, nil)
	if fileInfo != nil {
		fileInfo.Path = pathStr
	}
	p.storeFilePathToCache(driveId, pathStr, fileInfo)
	return fileInfo, error
}

func (p *WebPanClient) getFileInfoByPath(driveId string, index int, pathSlice *[]string, parentFileInfo *aliyunpan.FileEntity) (*aliyunpan.FileEntity, *apierror.ApiError) {
	if parentFileInfo == nil {
		// default root "/" entity
		parentFileInfo = aliyunpan.NewFileEntityForRootDir()
		if index == 0 && len(*pathSlice) == 1 {
			// root path "/"
			return parentFileInfo, nil
		}
		return p.getFileInfoByPath(driveId, index+1, pathSlice, parentFileInfo)
	}

	if index >= len(*pathSlice) {
		return parentFileInfo, nil
	}

	curPathStr := ""
	for idx := 0; idx <= index; idx++ {
		if (*pathSlice)[idx] == "" {
			continue
		}
		curPathStr += "/" + (*pathSlice)[idx]
	}
	// try cache
	if v := p.loadFilePathFromCache(driveId, curPathStr); v != nil {
		return p.getFileInfoByPath(driveId, index+1, pathSlice, v)
	}

	fileListParam := &aliyunpan.FileListParam{
		DriveId:      driveId,
		ParentFileId: parentFileInfo.FileId,
	}
	fileResult, err := p.FileListGetAll(fileListParam, 0)
	if err != nil {
		return nil, err
	}

	if fileResult == nil || len(fileResult) == 0 {
		return nil, apierror.NewApiError(apierror.ApiCodeFileNotFoundCode, "文件不存在")
	}
	var targetFile *aliyunpan.FileEntity = nil
	curParentPathStr := ""
	for idx := 0; idx <= (index - 1); idx++ {
		if (*pathSlice)[idx] == "" {
			continue
		}
		curParentPathStr += "/" + (*pathSlice)[idx]
	}
	for _, fileEntity := range fileResult {
		// cache all
		fileEntity.Path = curParentPathStr + "/" + fileEntity.FileName
		p.storeFilePathToCache(driveId, fileEntity.Path, fileEntity)

		// find target file
		if fileEntity.FileName == (*pathSlice)[index] {
			targetFile = fileEntity
		}
	}
	if targetFile != nil {
		// return
		return p.getFileInfoByPath(driveId, index+1, pathSlice, targetFile)
	}
	return nil, apierror.NewApiError(apierror.ApiCodeFileNotFoundCode, "文件不存在")
}

// FilesDirectoriesRecurseList 递归获取目录下的文件和目录列表
func (p *WebPanClient) FilesDirectoriesRecurseList(driveId string, path string, handleFileDirectoryFunc aliyunpan.HandleFileDirectoryFunc) aliyunpan.FileList {
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

func (p *WebPanClient) recurseList(driveId string, folderInfo *aliyunpan.FileEntity, depth int, handleFileDirectoryFunc aliyunpan.HandleFileDirectoryFunc, fld *aliyunpan.FileList) bool {
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

// FileListGetAll 获取指定目录下的所有文件列表
func (p *WebPanClient) FileListGetAll(param *aliyunpan.FileListParam, delayMilliseconds int) (aliyunpan.FileList, *apierror.ApiError) {
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

// FileGetPath 通过fileId获取对应的目录层级信息
func (p *WebPanClient) FileGetPath(driveId, fileId string) (*aliyunpan.FileGetPathResult, *apierror.ApiError) {
	header := map[string]string{
		"authorization": p.webToken.GetAuthorizationStr(),
	}

	fullUrl := &strings.Builder{}
	fmt.Fprintf(fullUrl, "%s/adrive/v1/file/get_path", API_URL)
	logger.Verboseln("do request url: " + fullUrl.String())

	postData := map[string]interface{}{
		"drive_id": driveId,
		"file_id":  fileId,
	}

	// request
	body, err := p.client.Fetch("POST", fullUrl.String(), postData, p.AddSignatureHeader(apiutil.AddCommonHeader(header)))
	logger.Verboseln("get file path response: ", string(body))
	if err != nil {
		logger.Verboseln("get file path error ", err)
		return nil, apierror.NewApiErrorWithError(err)
	}

	// handler common error
	if err1 := apierror.ParseCommonApiError(body); err1 != nil {
		return nil, err1
	}

	// parse result
	r := &aliyunpan.FileGetPathResult{}
	if err2 := json.Unmarshal(body, r); err2 != nil {
		logger.Verboseln("parse file path result json error ", err2)
		return nil, apierror.NewFailedApiError(err2.Error())
	}
	return r, nil
}
