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
	"net/http"
	"strconv"
	"strings"
	"time"
)

type (
	DownloadFuncCallback func(httpMethod, fullUrl string, headers map[string]string) (resp *http.Response, err error)

	FileDownloadRange struct {
		// 起始值，包含
		Offset int64
		// 结束值，包含
		End int64
	}

	GetFileDownloadUrlParam struct {
		DriveId   string `json:"drive_id"`
		FileId   string `json:"file_id"`
		ExpireSec int    `json:"expire_sec"`
	}

	GetFileDownloadUrlResult struct {
		Method      string    `json:"method"`
		Url         string    `json:"url"`
		InternalUrl string    `json:"internal_url"`
		CdnUrl      string    `json:"cdn_url"`
		Expiration  time.Time `json:"expiration"`
		Size        int       `json:"size"`
		Ratelimit   struct {
			PartSpeed int `json:"part_speed"`
			PartSize  int `json:"part_size"`
		} `json:"ratelimit"`
	}
)

const(
	// 资源被屏蔽，提示资源非法链接
	IllegalDownloadUrl = "https://pds-system-file.oss-cn-beijing.aliyuncs.com/illegal.mp4"
)

// GetFileDownloadUrl 获取文件下载URL路径
func (p *PanClient) GetFileDownloadUrl(param *GetFileDownloadUrlParam) (*GetFileDownloadUrlResult, *apierror.ApiError) {
	// header
	header := map[string]string {
		"authorization": p.webToken.GetAuthorizationStr(),
	}

	// url
	fullUrl := &strings.Builder{}
	fmt.Fprintf(fullUrl, "%s/v2/file/get_download_url", API_URL)
	logger.Verboseln("do request url: " + fullUrl.String())

	// data
	sec := param.ExpireSec
	if sec <= 0 {
		sec = 14400
	}
	postData := map[string]interface{} {
		"drive_id": param.DriveId,
		"file_id": param.FileId,
		"expire_sec": sec,
	}

	// request
	body, err := client.Fetch("POST", fullUrl.String(), postData, apiutil.AddCommonHeader(header))
	if err != nil {
		logger.Verboseln("get file download url error ", err)
		return nil, apierror.NewFailedApiError(err.Error())
	}

	// handler common error
	if err1 := apierror.ParseCommonApiError(body); err1 != nil {
		return nil, err1
	}

	// parse result
	r := &GetFileDownloadUrlResult{}
	if err2 := json.Unmarshal(body, r); err2 != nil {
		logger.Verboseln("parse file download url result json error ", err2)
		return nil, apierror.NewFailedApiError(err2.Error())
	}
	return r, nil
}

// DownloadFileData 下载文件内容
func (p *PanClient) DownloadFileData(downloadFileUrl string, fileRange FileDownloadRange, downloadFunc DownloadFuncCallback) *apierror.ApiError {
	// url
	fullUrl := &strings.Builder{}
	fmt.Fprintf(fullUrl, "%s", downloadFileUrl)
	logger.Verboseln("do request url: " + fullUrl.String())

	// header
	headers := map[string]string {
		"referer": "https://www.aliyundrive.com/",
	}

	// download data resume
	if fileRange.Offset != 0 || fileRange.End != 0 {
		rangeStr := "bytes=" + strconv.FormatInt(fileRange.Offset, 10) + "-"
		if fileRange.End != 0 {
			rangeStr += strconv.FormatInt(fileRange.End, 10)
		}
		headers["range"] = rangeStr
	}
	logger.Verboseln("do request url: " + fullUrl.String())

	// request
	_, err := downloadFunc("GET", fullUrl.String(), apiutil.AddCommonHeader(headers))
	//resp, err := p.client.Req("GET", fullUrl.String(), nil, headers)

	if err != nil {
		logger.Verboseln("download file data response failed")
		return apierror.NewApiErrorWithError(err)
	}
	return nil
}