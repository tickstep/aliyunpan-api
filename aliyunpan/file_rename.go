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

// FileRename 重命名文件
func (p *PanClient) FileRename(driveId, renameFileId, newName string) (bool, *apierror.ApiError) {
	if renameFileId == "" {
		return false, apierror.NewFailedApiError("请指定命名的文件")
	}
	// header
	header := map[string]string{
		"authorization": p.webToken.GetAuthorizationStr(),
	}

	// url
	fullUrl := &strings.Builder{}
	fmt.Fprintf(fullUrl, "%s/adrive/v3/file/update", API_URL)
	logger.Verboseln("do request url: " + fullUrl.String())

	// data
	postData := map[string]interface{}{
		"drive_id":        driveId,
		"file_id":         renameFileId,
		"name":            newName,
		"check_name_mode": "refuse",
	}

	// request
	body, err := p.client.Fetch("POST", fullUrl.String(), postData, p.AddSignatureHeader(apiutil.AddCommonHeader(header)))
	if err != nil {
		logger.Verboseln("get rename error ", err)
		return false, apierror.NewFailedApiError(err.Error())
	}

	// handler common error
	if err1 := apierror.ParseCommonApiError(body); err1 != nil {
		return false, err1
	}

	// parse result
	r := &FileEntity{}
	if err2 := json.Unmarshal(body, r); err2 != nil {
		logger.Verboseln("parse rename result json error ", err2)
		return false, apierror.NewFailedApiError(err2.Error())
	}
	return true, nil
}
