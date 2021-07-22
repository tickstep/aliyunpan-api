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
	"fmt"
	"github.com/tickstep/aliyunpan-api/aliyunpan/apierror"
	"github.com/tickstep/aliyunpan-api/aliyunpan/apiutil"
	"github.com/tickstep/library-go/logger"
	"strings"
)

type(
	FileDeleteParam struct {
		// 网盘ID
		DriveId        string `json:"drive_id"`
		// 文件ID
		FileId        string `json:"file_id"`
	}

	FileDeleteResult struct {
		// 文件ID
		FileId string
		// 是否成功
		Success bool
	}
)

// FileDelete 删除文件
func (p *PanClient) FileDelete(param []*FileDeleteParam) ([]*FileDeleteResult, *apierror.ApiError) {
	// url
	fullUrl := &strings.Builder{}
	fmt.Fprintf(fullUrl, "%s/v2/batch", API_URL)
	logger.Verboseln("do request url: " + fullUrl.String())

	// data
	requests,e := p.getFileDeleteBatchRequestList(param)
	if e != nil {
		return nil, e
	}
	batchParam := BatchRequestParam{
		Requests: requests,
		Resource: "file",
	}

	// request
	result,err := p.BatchTask(fullUrl.String(), &batchParam)
	if err != nil {
		logger.Verboseln("file delete error ", err)
		return nil, apierror.NewFailedApiError(err.Error())
	}

	// parse result
	r := []*FileDeleteResult{}
	for _,item := range result.Responses{
		r = append(r, &FileDeleteResult{
			FileId: item.Id,
			Success: item.Status == 204,
		})
	}
	return r, nil
}

func (p *PanClient) getFileDeleteBatchRequestList(param []*FileDeleteParam) (BatchRequestList, *apierror.ApiError) {
	if param == nil {
		return nil, apierror.NewFailedApiError("参数不能为空")
	}

	r := BatchRequestList{}
	for _,item := range param {
		r = append(r, &BatchRequest{
			Id:      item.FileId,
			Method:  "POST",
			Url:     "/recyclebin/trash",
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body:    apiutil.GetMapSet(item),
		})
	}
	return r, nil
}