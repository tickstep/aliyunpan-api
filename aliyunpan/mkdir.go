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
	"strings"
)

type (
	MkdirResult struct {
		ParentFileId string `json:"parent_file_id"`
		Type         string `json:"type"`
		FileId       string `json:"file_id"`
		DomainId     string `json:"domain_id"`
		DriveId      string `json:"drive_id"`
		FileName     string `json:"file_name"`
		EncryptMode  string `json:"encrypt_mode"`
	}
)

// Mkdir 创建文件夹
func (p *PanClient) Mkdir(driveId, parentFileId, dirName string) (*MkdirResult, *apierror.ApiError) {
	if parentFileId == "" {
		// 默认根目录
		parentFileId = DefaultRootParentFileId
	}
	header := map[string]string{
		"authorization": p.webToken.GetAuthorizationStr(),
	}

	fullUrl := &strings.Builder{}
	fmt.Fprintf(fullUrl, "%s/adrive/v2/file/createWithFolders", API_URL)
	logger.Verboseln("do request url: " + fullUrl.String())

	postData := map[string]interface{}{
		"drive_id":        driveId,
		"parent_file_id":  parentFileId,
		"name":            dirName,
		"check_name_mode": "refuse",
		"type":            "folder",
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
	r := &MkdirResult{}
	if err2 := json.Unmarshal(body, r); err2 != nil {
		logger.Verboseln("parse file info result json error ", err2)
		return nil, apierror.NewFailedApiError(err2.Error())
	}
	return r, nil
}

func (p *PanClient) MkdirByFullPath(driveId, fullPath string) (*MkdirResult, *apierror.ApiError) {
	fullPath = strings.ReplaceAll(fullPath, "//", "/")
	pathSlice := strings.Split(fullPath, "/")
	return p.MkdirRecursive(driveId, "", "", 0, pathSlice)
}

func (p *PanClient) MkdirRecursive(driveId, parentFileId string, fullPath string, index int, pathSlice []string) (*MkdirResult, *apierror.ApiError) {
	r := &MkdirResult{}
	if parentFileId == "" {
		// default root "/" entity
		parentFileId = NewFileEntityForRootDir().FileId
		if index == 0 && len(pathSlice) == 1 {
			// root path "/"
			r.FileId = parentFileId
			return r, nil
		}

		fullPath = ""
		return p.MkdirRecursive(driveId, parentFileId, fullPath, index+1, pathSlice)
	}

	if index >= len(pathSlice) {
		r.FileId = parentFileId
		return r, nil
	}

	listFilePath := &FileListParam{}
	listFilePath.DriveId = driveId
	listFilePath.ParentFileId = parentFileId
	fileResult, err := p.FileListGetAll(listFilePath, 0)
	if err != nil {
		r.FileId = ""
		return r, err
	}

	// existed?
	for _, fileEntity := range fileResult {
		if fileEntity.FileName == pathSlice[index] {
			return p.MkdirRecursive(driveId, fileEntity.FileId, fullPath+"/"+pathSlice[index], index+1, pathSlice)
		}
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
