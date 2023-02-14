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
	"github.com/tickstep/library-go/cachepool"
	"github.com/tickstep/library-go/logger"
	"github.com/tickstep/library-go/requester"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type (
	DownloadFuncCallback func(httpMethod, fullUrl string, headers map[string]string) (resp *http.Response, err error)

	// FileDownloadRange 分片。0-100,101-200,201-300...
	FileDownloadRange struct {
		// 起始值，包含
		Offset int64
		// 结束值，包含
		End int64
	}

	GetFileDownloadUrlParam struct {
		DriveId   string `json:"drive_id"`
		FileId    string `json:"file_id"`
		ExpireSec int    `json:"expire_sec"`
	}

	GetFileDownloadUrlResult struct {
		Method      string `json:"method"`
		Url         string `json:"url"`
		InternalUrl string `json:"internal_url"`
		CdnUrl      string `json:"cdn_url"`
		Expiration  string `json:"expiration"`
		Size        int64  `json:"size"`
		Ratelimit   struct {
			PartSpeed int64 `json:"part_speed"`
			PartSize  int64 `json:"part_size"`
		} `json:"ratelimit"`
	}
)

const (
	// 资源被屏蔽，提示资源非法链接
	IllegalDownloadUrlPrefix = "https://pds-system-file.oss-cn-beijing.aliyuncs.com/illegal"
)

// GetFileDownloadUrl 获取文件下载URL路径
func (p *PanClient) GetFileDownloadUrl(param *GetFileDownloadUrlParam) (*GetFileDownloadUrlResult, *apierror.ApiError) {
	// header
	header := map[string]string{
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
	postData := map[string]interface{}{
		"drive_id":   param.DriveId,
		"file_id":    param.FileId,
		"expire_sec": sec,
	}

	// request
	body, err := p.client.Fetch("POST", fullUrl.String(), postData, p.AddSignatureHeader(apiutil.AddCommonHeader(header)))
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
	// time format
	r.Expiration = apiutil.UtcTime2LocalFormat(r.Expiration)
	return r, nil
}

// DownloadFileData 下载文件内容
func (p *PanClient) DownloadFileData(downloadFileUrl string, fileRange FileDownloadRange, downloadFunc DownloadFuncCallback) *apierror.ApiError {
	// url
	fullUrl := &strings.Builder{}
	fmt.Fprintf(fullUrl, "%s", downloadFileUrl)
	logger.Verboseln("do request url: " + fullUrl.String())

	// header
	headers := map[string]string{
		"user-agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
		"referer":    "https://www.aliyundrive.com/",
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

	// request callback
	_, err := downloadFunc("GET", fullUrl.String(), headers)
	//resp, err := p.client.Req("GET", fullUrl.String(), nil, headers)

	if err != nil {
		logger.Verboseln("download file data response failed")
		return apierror.NewApiErrorWithError(err)
	}
	return nil
}

// DownloadFileDataAndSave 下载文件并存储到指定IO设备里面。该方法是同步阻塞的
func (p *PanClient) DownloadFileDataAndSave(downloadFileUrl string, fileRange FileDownloadRange, writerAt io.WriterAt) *apierror.ApiError {
	var resp *http.Response
	var err error
	var client = requester.NewHTTPClient()

	apierr := p.DownloadFileData(
		downloadFileUrl,
		fileRange,
		func(httpMethod, fullUrl string, headers map[string]string) (*http.Response, error) {
			resp, err = client.Req(httpMethod, fullUrl, nil, headers)
			if err != nil {
				return nil, err
			}
			return resp, err
		})

	if apierr != nil {
		return apierr
	}

	// close socket defer
	if resp != nil {
		defer func() {
			resp.Body.Close()
		}()
	}

	switch resp.StatusCode {
	case 200, 206:
		// do nothing, continue
		break
	case 416: //Requested Range Not Satisfiable
		fallthrough
	case 403: // Forbidden
		fallthrough
	case 406: // Not Acceptable
		return apierror.NewFailedApiError("")
	case 404:
		return apierror.NewFailedApiError("")
	case 429, 509: // Too Many Requests
		return apierror.NewFailedApiError("")
	default:
		return apierror.NewApiErrorWithError(fmt.Errorf("unexpected http status code, %d, %s", resp.StatusCode, resp.Status))
	}

	// save data
	var (
		buf                       = make([]byte, 4096)
		totalCount, readByteCount int
	)
	defer cachepool.SyncPool.Put(buf)

	var readErr error
	totalCount = 0

	for true {
		readByteCount, readErr = resp.Body.Read(buf)
		logger.Verboseln("get byte piece:", readByteCount)
		if readErr == io.EOF && readByteCount > 0 {
			// the last piece
			writerAt.WriteAt(buf[:readByteCount], fileRange.Offset+int64(totalCount))
			totalCount += readByteCount
			break
		}
		if readErr != nil {
			return apierror.NewApiErrorWithError(readErr)
		}

		// write
		writerAt.WriteAt(buf[:readByteCount], fileRange.Offset+int64(totalCount))
		totalCount += readByteCount
	}
	return nil
}
