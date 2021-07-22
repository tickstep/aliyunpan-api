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

)

// FileStarred 收藏文件
func (p *PanClient) FileStarred(param []*FileBatchActionParam) ([]*FileBatchActionResult, *apierror.ApiError) {
	return p.doFileStarredBatchRequestList(true, param)
}

// FileUnstarred 取消收藏文件
func (p *PanClient) FileUnstarred(param []*FileBatchActionParam) ([]*FileBatchActionResult, *apierror.ApiError) {
	return p.doFileStarredBatchRequestList(false, param)
}

func (p *PanClient) doFileStarredBatchRequestList(starred bool, param []*FileBatchActionParam) ([]*FileBatchActionResult, *apierror.ApiError) {
	if param == nil {
		return nil, apierror.NewFailedApiError("参数不能为空")
	}

	// url
	fullUrl := &strings.Builder{}
	fmt.Fprintf(fullUrl, "%s/v2/batch", API_URL)
	logger.Verboseln("do request url: " + fullUrl.String())

	// param
	pr := BatchRequestList{}
	for _,item := range param {
		body := apiutil.GetMapSet(item)
		if starred {
			body["starred"] = true
			body["custom_index_key"] = "starred_yes"
		} else {
			body["starred"] = false
			body["custom_index_key"] = ""
		}

		pr = append(pr, &BatchRequest{
			Id:      item.FileId,
			Method:  "PUT",
			Url:     "/file/update",
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body:    body,
		})
	}

	batchParam := BatchRequestParam{
		Requests: pr,
		Resource: "file",
	}

	// request
	result,err := p.BatchTask(fullUrl.String(), &batchParam)
	if err != nil {
		logger.Verboseln("file starred error ", err)
		return nil, apierror.NewFailedApiError(err.Error())
	}

	// parse result
	r := []*FileBatchActionResult{}
	for _,item := range result.Responses{
		r = append(r, &FileBatchActionResult{
			FileId: item.Id,
			Success: item.Status == 200,
		})
	}
	return r, nil
}