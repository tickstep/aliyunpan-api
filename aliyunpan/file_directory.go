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

package aliyunpan

import (
	"encoding/json"
	"fmt"
	"github.com/tickstep/aliyunpan-api/aliyunpan/apierror"
	"github.com/tickstep/aliyunpan-api/aliyunpan/apiutil"
	"github.com/tickstep/library-go/logger"
	"path"
	"strings"
)

type (
	// HandleFileDirectoryFunc 处理文件或目录的元信息, 返回值控制是否退出递归
	HandleFileDirectoryFunc func(depth int, fdPath string, fd *FileEntity, apierr *apierror.ApiError) bool

	// FileListParam 文件列表参数
	FileListParam struct {
		DriveId               string `json:"drive_id"`
		ParentFileId          string `json:"parent_file_id"`
		Limit                 int    `json:"limit"`
	}

	FileList []*FileEntity

	// FileEntity 文件/文件夹信息
	FileEntity struct {
		// 网盘ID
		DriveId         string    `json:"driveId"`
		// 域ID
		DomainId        string    `json:"domainId"`
		// FileId 文件ID
		FileId string `json:"fileId"`
		// FileName 文件名
		FileName string `json:"fileName"`
		// FileSize 文件大小
		FileSize uint64 `json:"fileSize"`
		// 文件类别 folder / file
		FileType            string    `json:"fileType"`
		// 创建时间
		CreatedAt       string `json:"createdAt"`
		// 最后修改时间
		UpdatedAt       string `json:"updatedAt"`
		// 后缀名，例如：dmg
		FileExtension   string    `json:"fileExtension"`
		// 文件上传ID
		UploadId        string    `json:"uploadId"`
		// 父文件夹ID
		ParentFileId    string    `json:"parentFileId"`
		// 内容CRC64校验值，只有文件才会有
		Crc64Hash       string    `json:"crc64Hash"`
		// 内容Hash值，只有文件才会有
		ContentHash     string    `json:"contentHash"`
		// 内容Hash计算方法，只有文件才会有，默认为：sha1
		ContentHashName string    `json:"contentHashName"`
		// FilePath 文件的完整路径
		Path string `json:"path"`
	}

	fileEntityResult struct {
		DriveId         string    `json:"drive_id"`
		DomainId        string    `json:"domain_id"`
		FileId          string    `json:"file_id"`
		Name            string    `json:"name"`
		Type            string    `json:"type"`
		ContentType     string    `json:"content_type"`
		CreatedAt       string `json:"created_at"`
		UpdatedAt       string `json:"updated_at"`
		FileExtension   string    `json:"file_extension"`
		MimeType        string    `json:"mime_type"`
		MimeExtension   string    `json:"mime_extension"`
		Hidden          bool      `json:"hidden"`
		Size            uint64       `json:"size"`
		Starred         bool      `json:"starred"`
		Status          string    `json:"status"`
		UploadId        string    `json:"upload_id"`
		ParentFileId    string    `json:"parent_file_id"`
		Crc64Hash       string    `json:"crc64_hash"`
		ContentHash     string    `json:"content_hash"`
		ContentHashName string    `json:"content_hash_name"`
		DownloadUrl     string    `json:"download_Url"`
		Url             string    `json:"Url"`
		Category        string    `json:"category"`
		EncryptMode     string    `json:"encrypt_mode"`
		PunishFlag      int       `json:"punish_flag"`
	}

	fileListResult struct {
		Items []*fileEntityResult `json:"items"`
		NextMarker string `json:"next_marker"`
	}
)

const (
	DefaultRootParentFileId string = "root"
)

// NewFileEntityForRootDir 创建根目录"/"的默认文件信息
func NewFileEntityForRootDir() *FileEntity {
	return &FileEntity {
		FileId: DefaultRootParentFileId,
		FileType: "folder",
		FileName: "/",
		ParentFileId: "",
		Path: "/",
	}
}

func createFileEntity(f *fileEntityResult) *FileEntity {
	return &FileEntity{
		DriveId: f.DriveId,
		DomainId: f.DomainId,
		FileId: f.FileId,
		FileName: f.Name,
		FileSize: f.Size,
		FileType: f.Type,
		CreatedAt: apiutil.UTCTimeFormat(f.CreatedAt),
		UpdatedAt: apiutil.UTCTimeFormat(f.UpdatedAt),
		FileExtension: f.FileExtension,
		UploadId: f.UploadId,
		ParentFileId: f.ParentFileId,
		Crc64Hash: f.Crc64Hash,
		ContentHash: f.ContentHash,
		ContentHashName: f.ContentHashName,
		Path: f.Name,
	}
}

// IsFolder 是否是文件夹
func (f *FileEntity) IsFolder() bool {
	return f.FileType == "folder"
}

// 是否是文件
func (f *FileEntity) IsFile() bool {
	return f.FileType == "file"
}

// 是否是网盘根目录
func (f *FileEntity) IsDriveRootFolder() bool {
	return f.FileId == DefaultRootParentFileId
}

// 文件展示信息
func (f *FileEntity) String() string {
	builder := &strings.Builder{}
	builder.WriteString("文件ID: " + f.FileId + "\n")
	builder.WriteString("文件名: " + f.FileName + "\n")
	if f.IsFolder() {
		builder.WriteString("文件类型: 目录\n")
	} else {
		builder.WriteString("文件类型: 文件\n")
	}
	builder.WriteString("文件路径: " + f.Path + "\n")
	return builder.String()
}

// TotalSize 获取目录下文件的总大小
func (fl FileList) TotalSize() uint64 {
	var size uint64
	for k := range fl {
		if fl[k] == nil {
			continue
		}

		size += fl[k].FileSize
	}
	return size
}

// Count 获取文件总数和目录总数
func (fl FileList) Count() (fileN, directoryN int64) {
	for k := range fl {
		if fl[k] == nil {
			continue
		}

		if fl[k].IsFolder() {
			directoryN++
		} else {
			fileN++
		}
	}
	return
}

// FileList 获取文件列表
func (p *PanClient) FileList(param *FileListParam) (FileList, *apierror.ApiError) {
	result := FileList{}
	if flr,err := p.fileListReq(param); err == nil {
		for k := range flr.Items {
			if flr.Items[k] == nil {
				continue
			}

			result = append(result, createFileEntity(flr.Items[k]))
		}
	}
	return result, nil
}

func (p *PanClient) fileListReq(param *FileListParam) (*fileListResult, *apierror.ApiError) {
	header := map[string]string {
		"accept": "application/json, text/plain, */*",
		"referer": "https://www.aliyundrive.com/",
		"origin": "https://www.aliyundrive.com",
		"content-type": "application/json;charset=UTF-8",
		"user-agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
		"authorization": p.webToken.GetAuthorizationStr(),
	}

	fullUrl := &strings.Builder{}
	fmt.Fprintf(fullUrl, "%s/v2/file/list", API_URL)
	logger.Verboseln("do request url: " + fullUrl.String())

	pFileId := param.ParentFileId
	if pFileId == "" {
		pFileId = DefaultRootParentFileId
	}
	limit := param.Limit
	if limit <= 0 {
		limit = 100
	}
	postData := map[string]interface{} {
		"drive_id": param.DriveId,
		"parent_file_id": pFileId,
		"limit": limit,
		"all": false,
		"url_expire_sec": 1600,
		"image_thumbnail_process": "image/resize,w_400/format,jpeg",
		"image_url_process": "image/resize,w_1920/format,jpeg",
		"video_thumbnail_process": "video/snapshot,t_0,f_jpg,ar_auto,w_800",
		"fields": "*",
		"order_by": "updated_at",
		"order_direction": "DESC",
	}

	// request
	body, err := client.Fetch("POST", fullUrl.String(), postData, header)
	if err != nil {
		logger.Verboseln("get file list error ", err)
		return nil, apierror.NewFailedApiError(err.Error())
	}

	// handler common error
	if err1 := apierror.ParseCommonApiError(body); err1 != nil {
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
func (p *PanClient) FileInfoById(driveId, fileId string) (*FileEntity, *apierror.ApiError) {
	header := map[string]string {
		"accept": "application/json, text/plain, */*",
		"referer": "https://www.aliyundrive.com/",
		"origin": "https://www.aliyundrive.com",
		"content-type": "application/json;charset=UTF-8",
		"user-agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
		"authorization": p.webToken.GetAuthorizationStr(),
	}

	fullUrl := &strings.Builder{}
	fmt.Fprintf(fullUrl, "%s/v2/file/get", API_URL)
	logger.Verboseln("do request url: " + fullUrl.String())

	pFileId := fileId
	if pFileId == "" {
		pFileId = DefaultRootParentFileId
	}
	postData := map[string]interface{} {
		"drive_id": driveId,
		"file_id": pFileId,
	}

	// request
	body, err := client.Fetch("POST", fullUrl.String(), postData, header)
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
func (p *PanClient) FileInfoByPath(driveId string, pathStr string) (fileInfo *FileEntity, error *apierror.ApiError) {
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

	var pathSlice []string
	if pathStr == "/" {
		pathSlice = []string{""}
	} else {
		pathSlice = strings.Split(pathStr, PathSeparator)
		if pathSlice[0] != "" {
			return nil, apierror.NewFailedApiError("pathStr必须是绝对路径")
		}
	}
	fileInfo,error = p.getFileInfoByPath(driveId, 0, &pathSlice, nil)
	if fileInfo != nil {
		fileInfo.Path = pathStr
	}
	return fileInfo, error
}

func (p *PanClient) getFileInfoByPath(driveId string, index int, pathSlice *[]string, parentFileInfo *FileEntity) (*FileEntity, *apierror.ApiError)  {
	if parentFileInfo == nil {
		// default root "/" entity
		parentFileInfo = NewFileEntityForRootDir()
		if index == 0 && len(*pathSlice) == 1 {
			// root path "/"
			return parentFileInfo, nil
		}
		return p.getFileInfoByPath(driveId, index + 1, pathSlice, parentFileInfo)
	}

	if index >= len(*pathSlice) {
		return parentFileInfo, nil
	}

	fileListParam := &FileListParam{
		DriveId: driveId,
		ParentFileId: parentFileInfo.FileId,
	}
	fileResult, err := p.FileList(fileListParam)
	if err != nil {
		return nil, err
	}

	if fileResult == nil || len(fileResult) == 0 {
		return nil, apierror.NewApiError(apierror.ApiCodeFileNotFoundCode, "文件不存在")
	}
	for _, fileEntity := range fileResult {
		if fileEntity.FileName == (*pathSlice)[index] {
			return p.getFileInfoByPath(driveId, index + 1, pathSlice, fileEntity)
		}
	}
	return nil, apierror.NewApiError(apierror.ApiCodeFileNotFoundCode, "文件不存在")
}

// FilesDirectoriesRecurseList 递归获取目录下的文件和目录列表
func (p *PanClient) FilesDirectoriesRecurseList(driveId string, path string, handleFileDirectoryFunc HandleFileDirectoryFunc) FileList {
	targetFileInfo, er := p.FileInfoByPath(driveId, path)
	if er != nil {
		if handleFileDirectoryFunc != nil {
			handleFileDirectoryFunc(0, path, nil, er)
		}
		return nil
	}
	if !targetFileInfo.IsFolder() {
		if handleFileDirectoryFunc != nil {
			handleFileDirectoryFunc(0, path, targetFileInfo, nil)
		}
		return FileList{targetFileInfo}
	}

	fld := &FileList{}
	ok := p.recurseList(driveId, targetFileInfo, 1, handleFileDirectoryFunc, fld)
	if !ok {
		return nil
	}
	return *fld
}

func (p *PanClient) recurseList(driveId string, folderInfo *FileEntity, depth int, handleFileDirectoryFunc HandleFileDirectoryFunc, fld *FileList) bool {
	flp := &FileListParam{
		DriveId: driveId,
		ParentFileId: folderInfo.FileId,
	}
	r, apiError := p.FileList(flp)
	if apiError != nil {
		if handleFileDirectoryFunc != nil {
			handleFileDirectoryFunc(depth, folderInfo.Path, nil, apiError)
		}
		return false
	}
	ok := true
	for _, fi := range r {
		fi.Path = strings.ReplaceAll(folderInfo.Path + PathSeparator + fi.FileName, "//", "/")
		*fld = append(*fld, fi)
		if fi.IsFolder() {
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