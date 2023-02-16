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
	RecycleBinFileListParam struct {
		DriveId string `json:"drive_id"`
		Limit   int    `json:"limit"`
		Marker  string `json:"marker"`
	}
)

// RecycleBinFileList 获取回收站文件列表
func (p *PanClient) RecycleBinFileList(param *RecycleBinFileListParam) (*FileListResult, *apierror.ApiError) {
	result := &FileListResult{
		FileList:   FileList{},
		NextMarker: "",
	}
	if flr, err := p.recycleBinFileListReq(param); err == nil {
		for k := range flr.Items {
			if flr.Items[k] == nil {
				continue
			}

			result.FileList = append(result.FileList, createFileEntity(flr.Items[k]))
		}
		result.NextMarker = flr.NextMarker
	}
	return result, nil
}

// RecycleBinFileListGetAll 获取所有列表文件
func (p *PanClient) RecycleBinFileListGetAll(param *RecycleBinFileListParam) (FileList, *apierror.ApiError) {
	internalParam := &RecycleBinFileListParam{
		DriveId: param.DriveId,
		Limit:   param.Limit,
		Marker:  param.Marker,
	}
	if internalParam.Limit <= 0 {
		internalParam.Limit = 100
	}

	fileList := FileList{}
	result, err := p.RecycleBinFileList(internalParam)
	if err != nil || result == nil {
		return nil, err
	}
	fileList = append(fileList, result.FileList...)

	// more page?
	for len(result.NextMarker) > 0 {
		internalParam.Marker = result.NextMarker
		result, err = p.RecycleBinFileList(internalParam)
		if err == nil && result != nil {
			fileList = append(fileList, result.FileList...)
		} else {
			break
		}
	}
	return fileList, nil
}

func (p *PanClient) recycleBinFileListReq(param *RecycleBinFileListParam) (*fileListResult, *apierror.ApiError) {
	header := map[string]string{
		"authorization": p.webToken.GetAuthorizationStr(),
		"referer":       "https://www.aliyundrive.com/",
		"origin":        "https://www.aliyundrive.com",
	}

	fullUrl := &strings.Builder{}
	fmt.Fprintf(fullUrl, "%s/v2/recyclebin/list", API_URL)
	logger.Verboseln("do request url: " + fullUrl.String())

	limit := param.Limit
	if limit <= 0 {
		limit = 100
	}
	postData := map[string]interface{}{
		"drive_id":                param.DriveId,
		"limit":                   limit,
		"image_thumbnail_process": "image/resize,w_400/format,jpeg",
		"video_thumbnail_process": "video/snapshot,t_0,f_jpg,ar_auto,w_800",
		"order_by":                "name",
		"order_direction":         "DESC",
	}
	if len(param.Marker) > 0 {
		postData["marker"] = param.Marker
	}

	// request
	body, err := p.client.Fetch("POST", fullUrl.String(), postData, p.AddSignatureHeader(apiutil.AddCommonHeader(header)))
	if err != nil {
		logger.Verboseln("get recycle bin file list error ", err)
		return nil, apierror.NewFailedApiError(err.Error())
	}

	// handler common error
	if err1 := apierror.ParseCommonApiError(body); err1 != nil {
		return nil, err1
	}

	// parse result
	r := &fileListResult{}
	if err2 := json.Unmarshal(body, r); err2 != nil {
		logger.Verboseln("parse recycle bin file list result json error ", err2)
		return nil, apierror.NewFailedApiError(err2.Error())
	}
	return r, nil
}
