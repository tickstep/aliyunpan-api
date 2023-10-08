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
	FileCrossCopyParam struct {
		// FromDriveId 源网盘ID
		FromDriveId string `json:"from_drive_id"`
		// FromFileIds 源网盘文件列表ID
		FromFileIds []string `json:"from_file_ids"`
		// ToDriveId 目标网盘ID。必须和源网盘ID不一样，否则会报错
		ToDriveId string `json:"to_drive_id"`
		// ToParentFileId 目标网盘目录ID
		ToParentFileId string `json:"to_parent_fileId"`
	}

	FileCrossCopyResult struct {
		DriveId       string `json:"drive_id"`
		FileId        string `json:"file_id"`
		SourceDriveId string `json:"source_drive_id"`
		SourceFileId  string `json:"source_file_id"`
		// Status 结果状态，201代表成功
		Status int `json:"status"`
	}
)

// FileCrossDriveCopy 跨网盘复制文件，支持资源库和备份盘之间复制文件
func (p *PanClient) FileCrossDriveCopy(param *FileCrossCopyParam) ([]*FileCrossCopyResult, *apierror.ApiError) {
	// header
	header := map[string]string{
		"authorization": p.webToken.GetAuthorizationStr(),
	}

	// url
	fullUrl := &strings.Builder{}
	fmt.Fprintf(fullUrl, "%s/adrive/v2/file/crossDriveCopy", API_URL)
	logger.Verboseln("do request url: " + fullUrl.String())

	// data
	postData := param
	if param.FromDriveId == param.ToDriveId {
		return nil, apierror.NewFailedApiError("目标网盘ID和源网盘ID必须不一样")
	}

	// request
	body, err := p.client.Fetch("POST", fullUrl.String(), postData, p.AddSignatureHeader(apiutil.AddCommonHeader(header)))
	if err != nil {
		logger.Verboseln("do cross drive copy error ", err)
		return nil, apierror.NewFailedApiError(err.Error())
	}
	logger.Verboseln("response: ", string(body))

	// handler common error
	if err1 := apierror.ParseCommonApiError(body); err1 != nil {
		return nil, err1
	}

	// parse result
	result := struct {
		Items []*FileCrossCopyResult `json:"items"`
	}{}
	if err2 := json.Unmarshal(body, &result); err2 != nil {
		logger.Verboseln("parse cross drive copy result json error ", err2)
		return nil, apierror.NewFailedApiError(err2.Error())
	}

	// parse result
	r := []*FileCrossCopyResult{}
	for _, item := range result.Items {
		r = append(r, item)
	}
	return r, nil
}

// FileCrossDriveMove 跨网盘移动文件，只支持从资源库移动到备份盘
func (p *PanClient) FileCrossDriveMove(param *FileCrossCopyParam) ([]*FileCrossCopyResult, *apierror.ApiError) {
	// header
	header := map[string]string{
		"authorization": p.webToken.GetAuthorizationStr(),
	}

	// url
	fullUrl := &strings.Builder{}
	fmt.Fprintf(fullUrl, "%s/adrive/v2/file/crossDriveMove", API_URL)
	logger.Verboseln("do request url: " + fullUrl.String())

	// data
	postData := param
	if param.FromDriveId == param.ToDriveId {
		return nil, apierror.NewFailedApiError("目标网盘ID和源网盘ID必须不一样")
	}

	// request
	body, err := p.client.Fetch("POST", fullUrl.String(), postData, p.AddSignatureHeader(apiutil.AddCommonHeader(header)))
	if err != nil {
		logger.Verboseln("do cross drive copy error ", err)
		return nil, apierror.NewFailedApiError(err.Error())
	}
	logger.Verboseln("response: ", string(body))

	// handler common error
	if err1 := apierror.ParseCommonApiError(body); err1 != nil {
		return nil, err1
	}

	// parse result
	result := struct {
		Items []*FileCrossCopyResult `json:"items"`
	}{}
	if err2 := json.Unmarshal(body, &result); err2 != nil {
		logger.Verboseln("parse cross drive copy result json error ", err2)
		return nil, apierror.NewFailedApiError(err2.Error())
	}

	// parse result
	r := []*FileCrossCopyResult{}
	for _, item := range result.Items {
		r = append(r, item)
	}
	return r, nil
}
