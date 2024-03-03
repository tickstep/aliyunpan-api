package aliyunpan_web

import (
	"encoding/json"
	"fmt"
	"github.com/tickstep/aliyunpan-api/aliyunpan"
	"github.com/tickstep/aliyunpan-api/aliyunpan/apierror"
	"github.com/tickstep/aliyunpan-api/aliyunpan/apiutil"
	"github.com/tickstep/library-go/logger"
	"strings"
)

type ()

// VideoGetPreviewPlayInfo 获取视频预览信息，调用该接口会触发视频云端转码
func (p *PanClient) VideoGetPreviewPlayInfo(param *aliyunpan.VideoGetPreviewPlayInfoParam) (*aliyunpan.VideoGetPreviewPlayInfoResult, error) {
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
	r := &aliyunpan.VideoGetPreviewPlayInfoResult{}
	if err2 := json.Unmarshal(body, r); err2 != nil {
		logger.Verboseln("parse video preview play info json error ", err2)
		return nil, apierror.NewFailedApiError(err2.Error())
	}
	return r, nil
}
