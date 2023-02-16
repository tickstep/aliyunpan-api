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
	VideoGetPreviewPlayInfoParam struct {
		DriveId string `json:"drive_id"`
		// FileId 视频文件ID
		FileId string `json:"file_id"`
	}

	VideoGetPreviewPlayInfoResult struct {
		DomainId             string `json:"domain_id"`
		DriveId              string `json:"drive_id"`
		FileId               string `json:"file_id"`
		VideoPreviewPlayInfo struct {
			Category string `json:"category"`
			Meta     struct {
				Duration            float64 `json:"duration"`
				Width               int     `json:"width"`
				Height              int     `json:"height"`
				LiveTranscodingMeta struct {
					TsSegment    int `json:"ts_segment"`
					TsTotalCount int `json:"ts_total_count"`
					TsPreCount   int `json:"ts_pre_count"`
				} `json:"live_transcoding_meta"`
			} `json:"meta"`
			LiveTranscodingTaskList []struct {
				TemplateId     string `json:"template_id"`
				TemplateName   string `json:"template_name"`
				TemplateWidth  int    `json:"template_width"`
				TemplateHeight int    `json:"template_height"`
				Status         string `json:"status"`
				Stage          string `json:"stage"`
				URL            string `json:"url"`
			} `json:"live_transcoding_task_list"`
		} `json:"video_preview_play_info"`
	}
)

// VideoGetPreviewPlayInfo 获取视频预览信息，调用该接口会触发视频云端转码
func (p *PanClient) VideoGetPreviewPlayInfo(param *VideoGetPreviewPlayInfoParam) (*VideoGetPreviewPlayInfoResult, error) {
	header := map[string]string{
		"authorization": p.webToken.GetAuthorizationStr(),
	}

	fullUrl := &strings.Builder{}
	fmt.Fprintf(fullUrl, "%s/v2/file/get_video_preview_play_info", API_URL)
	logger.Verboseln("do request url: " + fullUrl.String())

	postData := map[string]interface{}{
		"category":    "live_transcoding",
		"drive_id":    param.DriveId,
		"file_id":     param.FileId,
		"template_id": "",
	}

	// request
	body, err := p.client.Fetch("POST", fullUrl.String(), postData, p.AddSignatureHeader(apiutil.AddCommonHeader(header)))
	logger.Verboseln("response: " + string(body))
	if err != nil {
		logger.Verboseln("get video preview play info error ", err)
		return nil, apierror.NewFailedApiError(err.Error())
	}

	// handler common error
	if err1 := apierror.ParseCommonApiError(body); err1 != nil {
		return nil, err1
	}

	// parse result
	r := &VideoGetPreviewPlayInfoResult{}
	if err2 := json.Unmarshal(body, r); err2 != nil {
		logger.Verboseln("parse video preview play info json error ", err2)
		return nil, apierror.NewFailedApiError(err2.Error())
	}
	return r, nil
}
