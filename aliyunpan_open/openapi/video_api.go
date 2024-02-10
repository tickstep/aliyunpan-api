package openapi

import (
	"encoding/json"
	"fmt"
	"github.com/tickstep/library-go/logger"
	"strings"
)

type (
	// VideoGetPreviewPlayInfoParam 获取文件播放详情参数
	VideoGetPreviewPlayInfoParam struct {
		// DriveId 网盘id
		DriveId string `json:"drive_id"`
		// ParentFileId 根目录为root
		FileId string `json:"file_id"`
		// Category live_transcoding 边转边播
		Category string `json:"category"`
	}

	LiveTranscodingTask struct {
		TemplateId     string `json:"template_id"`
		TemplateName   string `json:"template_name"`
		TemplateWidth  int    `json:"template_width"`
		TemplateHeight int    `json:"template_height"`
		// Status 状态。 枚举值如下： finished, 索引完成，可以获取到url, running, 正在索引，请稍等片刻重试, failed, 转码失败，请检查是否媒体文件，如果有疑问请联系客服
		Status string `json:"status"`
		Stage  string `json:"stage"`
		Url    string `json:"url"`
	}
	// VideoGetPreviewPlayInfoResult 获取文件播放详情返回值
	VideoGetPreviewPlayInfoResult struct {
		DriveId              string `json:"drive_id"`
		FileId               string `json:"file_id"`
		VideoPreviewPlayInfo struct {
			Category string `json:"category"`
			Meta     struct {
				Duration float64 `json:"duration"`
				Width    int     `json:"width"`
				Height   int     `json:"height"`
			} `json:"meta"`
			LiveTranscodingTaskList []*LiveTranscodingTask `json:"live_transcoding_task_list"`
		} `json:"video_preview_play_info"`
	}
)

// VideoGetPreviewPlayInfo 获取文件播放详情
func (a *AliPanClient) VideoGetPreviewPlayInfo(param *VideoGetPreviewPlayInfoParam) (*VideoGetPreviewPlayInfoResult, *AliApiErrResult) {
	fullUrl := &strings.Builder{}
	fmt.Fprintf(fullUrl, "%s/adrive/v1.0/openFile/getVideoPreviewPlayInfo", OPENAPI_URL)
	logger.Verboseln("do request url: " + fullUrl.String())

	// parameters
	postData := param

	// request
	resp, err := a.httpclient.Req("POST", fullUrl.String(), postData, a.Headers())
	if err != nil {
		logger.Verboseln("video get preview play info error ", err)
		return nil, NewAliApiHttpError(err.Error())
	}

	// handler common error
	var body []byte
	var apiErrResult *AliApiErrResult
	if body, apiErrResult = ParseCommonOpenApiError(resp); apiErrResult != nil {
		return nil, apiErrResult
	}

	// parse result
	r := &VideoGetPreviewPlayInfoResult{}
	if err2 := json.Unmarshal(body, r); err2 != nil {
		logger.Verboseln("parse video get preview play info result json error ", err2)
		return nil, NewAliApiAppError(err2.Error())
	}
	return r, nil
}
