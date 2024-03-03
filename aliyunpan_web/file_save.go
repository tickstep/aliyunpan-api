package aliyunpan_web

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/tickstep/aliyunpan-api/aliyunpan/apierror"
	"github.com/tickstep/aliyunpan-api/aliyunpan/apiutil"
	"github.com/tickstep/library-go/logger"
)

type GetShareByAnonymous struct {
	CreatorID           string    `json:"creator_id"`
	CreatorName         string    `json:"creator_name"`
	CreatorPhone        string    `json:"creator_phone"`
	Expiration          string    `json:"expiration"`
	UpdatedAt           time.Time `json:"updated_at"`
	VIP                 string    `json:"vip"`
	Avatar              string    `json:"avatar"`
	ShareName           string    `json:"share_name"`
	FileCount           int       `json:"file_count"`
	IsCreatorFollowable bool      `json:"is_creator_followable"`
	IsFollowingCreator  bool      `json:"is_following_creator"`
	DisplayName         string    `json:"display_name"`
	ShareTitle          string    `json:"share_title"`
	HasPassword         bool      `json:"has_pwd"`
	SaveButton          struct {
		Text          string `json:"text"`
		SelectAllText string `json:"select_all_text"`
	} `json:"save_button"`
}

func (p *PanClient) GetShareInfo(shareID string) (*GetShareByAnonymous, *apierror.ApiError) {
	// header
	header := map[string]string{
		"authorization": p.webToken.GetAuthorizationStr(),
	}

	// url
	fullUrl := &strings.Builder{}
	fmt.Fprintf(fullUrl, "%s/adrive/v3/share_link/get_share_by_anonymous?share_id=%s", API_URL, shareID)
	logger.Verboseln("do request url: " + fullUrl.String())

	// data
	postData := map[string]interface{}{
		"share_id": shareID,
	}

	// request
	body, err := p.client.Fetch("POST", fullUrl.String(), postData, p.AddSignatureHeader(apiutil.AddCommonHeader(header)))
	logger.Verboseln(string(body))
	if err != nil {
		logger.Verboseln("get share by anonymous error ", err)
		return nil, apierror.NewFailedApiError(err.Error())
	}

	// handler common error
	if err1 := apierror.ParseCommonApiError(body); err1 != nil {
		return nil, err1
	}

	// parse result
	r := &GetShareByAnonymous{}
	if err2 := json.Unmarshal(body, r); err2 != nil {
		logger.Verboseln("parse share by anonymous json error ", err2)
		return nil, apierror.NewFailedApiError(err2.Error())
	}
	return r, nil
}

type GetShareTokenResult struct {
	ExpireTime time.Time `json:"expire_time"`
	ExpiresIn  int       `json:"expires_in"`
	ShareToken string    `json:"share_token"`
}

func (p *PanClient) GetShareToken(shareID, sharePwd string) (*GetShareTokenResult, *apierror.ApiError) {
	// header
	header := map[string]string{
		"authorization": p.webToken.GetAuthorizationStr(),
	}

	// url
	fullUrl := &strings.Builder{}
	fmt.Fprintf(fullUrl, "%s/v2/share_link/get_share_token", API_URL)
	logger.Verboseln("do request url: " + fullUrl.String())

	// data
	postData := map[string]interface{}{
		"share_id":  shareID,
		"share_pwd": sharePwd,
	}

	// request
	body, err := p.client.Fetch("POST", fullUrl.String(), postData, p.AddSignatureHeader(apiutil.AddCommonHeader(header)))
	logger.Verboseln(string(body))
	if err != nil {
		logger.Verboseln("get share token error ", err)
		return nil, apierror.NewFailedApiError(err.Error())
	}

	// handler common error
	if err1 := apierror.ParseCommonApiError(body); err1 != nil {
		return nil, err1
	}

	// parse result
	r := &GetShareTokenResult{}
	if err2 := json.Unmarshal(body, r); err2 != nil {
		logger.Verboseln("parse share token json error ", err2)
		return nil, apierror.NewFailedApiError(err2.Error())
	}
	return r, nil
}

type ListByShareResult struct {
	Items      []*ListByShareItem `json:"items"`
	NextMarker string             `json:"next_marker"`
}

type ListByShareItem struct {
	DriveID            string    `json:"drive_id"`
	DomainID           string    `json:"domain_id"`
	FileID             string    `json:"file_id"`
	ShareID            string    `json:"share_id"`
	Name               string    `json:"name"`
	Type               string    `json:"type"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	FileExtension      string    `json:"file_extension,omitempty"`
	MimeType           string    `json:"mime_type,omitempty"`
	MimeExtension      string    `json:"mime_extension,omitempty"`
	Size               int       `json:"size,omitempty"`
	ParentFileID       string    `json:"parent_file_id,omitempty"`
	Thumbnail          string    `json:"thumbnail,omitempty"`
	Category           string    `json:"category,omitempty"`
	ImageMediaMetadata struct {
		Width  int    `json:"width"`
		Height int    `json:"height"`
		Exif   string `json:"exif"`
	} `json:"image_media_metadata,omitempty"`
	PunishFlag int    `json:"punish_flag,omitempty"`
	RevisionID string `json:"revision_id,omitempty"`
}

func (p *PanClient) GetListByShare(shareToken, shareID, marker string) (*ListByShareResult, *apierror.ApiError) {
	// header
	header := map[string]string{
		"authorization": p.webToken.GetAuthorizationStr(),
		"x-share-token": shareToken,
	}

	// url
	fullUrl := &strings.Builder{}
	fmt.Fprintf(fullUrl, "%s/adrive/v2/file/list_by_share", API_URL)
	logger.Verboseln("do request url: " + fullUrl.String())

	// data
	postData := map[string]interface{}{
		"share_id":                shareID,
		"parent_file_id":          "root",
		"limit":                   20,
		"image_thumbnail_process": "image/resize,w_256/format,jpeg",
		"image_url_process":       "image/resize,w_1920/format,jpeg/interlace,1",
		"video_thumbnail_process": "video/snapshot,t_1000,f_jpg,ar_auto,w_256",
		"order_by":                "name",
		"order_direction":         "DESC",
	}
	if marker != "" {
		postData["marker"] = marker
	}

	// request
	body, err := p.client.Fetch("POST", fullUrl.String(), postData, p.AddSignatureHeader(apiutil.AddCommonHeader(header)))
	logger.Verboseln(string(body))
	if err != nil {
		logger.Verboseln("get list by share error ", err)
		return nil, apierror.NewFailedApiError(err.Error())
	}

	// handler common error
	if err1 := apierror.ParseCommonApiError(body); err1 != nil {
		return nil, err1
	}

	// parse result
	r := &ListByShareResult{}
	if err2 := json.Unmarshal(body, r); err2 != nil {
		logger.Verboseln("parse list by share json error ", err2)
		return nil, apierror.NewFailedApiError(err2.Error())
	}
	return r, nil
}

type (
	FileSaveParam struct {
		ShareID        string `json:"share_id"`
		FileId         string `json:"file_id"`
		AutoRename     bool   `json:"auto_rename"` // default: true
		ToDriveId      string `json:"to_drive_id"`
		ToParentFileId string `json:"to_parent_file_id"`
	}

	FileSaveResult struct {
		DomainID    string `json:"domain_id"`
		FileId      string `json:"file_id"`
		DriveId     string `json:"tdrive_id"`
		AsyncTaskId string `json:"async_task_id,omitempty"`
		Status      int
	}
)

func (p *PanClient) FileCopy(shareToken string, param []*FileSaveParam) ([]*FileSaveResult, *apierror.ApiError) {
	// url
	fullUrl := &strings.Builder{}
	fmt.Fprintf(fullUrl, "%s/adrive/v2/batch", API_URL)
	logger.Verboseln("do request url: " + fullUrl.String())

	// data
	requests, e := p.getFileCopyBatchRequestList(param)
	if e != nil {
		return nil, e
	}
	batchParam := BatchRequestParam{
		Requests: requests,
		Resource: "file",
	}

	// request
	result, err := p.BatchTask(fullUrl.String(), &batchParam, [2]string{"x-share-token", shareToken})
	if err != nil {
		logger.Verboseln("file copy error ", err)
		return nil, apierror.NewFailedApiError(err.Error())
	}

	// parse result
	r := []*FileSaveResult{}
	for _, item := range result.Responses {
		var result FileSaveResult
		if domainID, ok := item.Body["domain_id"]; ok {
			result.DomainID, _ = domainID.(string)
		}
		if DriveId, ok := item.Body["drive_id"]; ok {
			result.DriveId, _ = DriveId.(string)
		}
		if FileId, ok := item.Body["file_id"]; ok {
			result.FileId, _ = FileId.(string)
		}
		if AsyncTaskId, ok := item.Body["async_task_id"]; ok {
			result.AsyncTaskId, _ = AsyncTaskId.(string)
		}
		result.Status = item.Status
		r = append(r, &result)
	}

	return r, nil
}

func (p *PanClient) getFileCopyBatchRequestList(param []*FileSaveParam) (BatchRequestList, *apierror.ApiError) {
	if param == nil {
		return nil, apierror.NewFailedApiError("参数不能为空")
	}

	r := BatchRequestList{}
	for i, item := range param {
		r = append(r, &BatchRequest{
			Id:     strconv.Itoa(i),
			Method: "POST",
			Url:    "/file/copy",
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body: apiutil.GetMapSet(item),
		})
	}
	return r, nil
}

type AsyncTaskGetResult struct {
	AsyncTaskId string
	Success     bool
}

func (p *PanClient) AsyncTaskGet(shareToken string, asyncTaskIds []string) ([]*AsyncTaskGetResult, *apierror.ApiError) {
	// url
	fullUrl := &strings.Builder{}
	fmt.Fprintf(fullUrl, "%s/adrive/v2/batch", API_URL)
	logger.Verboseln("do request url: " + fullUrl.String())

	// data
	requests, e := p.getAsyncTaskGetBatchRequestList(asyncTaskIds)
	if e != nil {
		return nil, e
	}
	batchParam := BatchRequestParam{
		Requests: requests,
		Resource: "file",
	}

	// request
	result, err := p.BatchTask(fullUrl.String(), &batchParam, [2]string{"x-share-token", shareToken})
	if err != nil {
		logger.Verboseln("async task get error ", err)
		return nil, apierror.NewFailedApiError(err.Error())
	}

	// parse result
	r := []*AsyncTaskGetResult{}
	for _, item := range result.Responses {
		r = append(r, &AsyncTaskGetResult{
			AsyncTaskId: item.Id,
			Success:     item.Status == 200,
		})
	}

	return r, nil
}

func (p *PanClient) getAsyncTaskGetBatchRequestList(param []string) (BatchRequestList, *apierror.ApiError) {
	if param == nil {
		return nil, apierror.NewFailedApiError("参数不能为空")
	}

	r := BatchRequestList{}
	for _, item := range param {
		r = append(r, &BatchRequest{
			Id:     item,
			Method: "POST",
			Url:    "/async_task/get",
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body: map[string]interface{}{"async_task_id": item},
		})
	}
	return r, nil
}
