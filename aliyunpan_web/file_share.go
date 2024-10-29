package aliyunpan_web

import (
	"encoding/json"
	"fmt"
	"github.com/tickstep/aliyunpan-api/aliyunpan"
	"github.com/tickstep/aliyunpan-api/aliyunpan/apierror"
	"github.com/tickstep/aliyunpan-api/aliyunpan/apiutil"
	"github.com/tickstep/library-go/logger"
	"strings"
	"time"
)

type (
	shareEntityResult struct {
		CreatedAt   string `json:"created_at"`
		Creator     string `json:"creator"`
		Description string `json:"description"`
		// 下载次数
		DownloadCount int      `json:"download_count"`
		DriveId       string   `json:"drive_id"`
		Expiration    string   `json:"expiration"`
		Expired       bool     `json:"expired"`
		FileId        string   `json:"file_id"`
		FileIdList    []string `json:"file_id_list"`
		// 浏览次数
		PreviewCount int `json:"preview_count"`
		// 转存次数
		SaveCount   int    `json:"save_count"`
		ShareId     string `json:"share_id"`
		ShareMsg    string `json:"share_msg"`
		ShareName   string `json:"share_name"`
		SharePolicy string `json:"share_policy"`
		SharePwd    string `json:"share_pwd"`
		ShareUrl    string `json:"share_url"`
		Status      string `json:"status"`
		UpdatedAt   string `json:"updated_at"`

		FirstFile *fileEntityResult `json:"first_file"`
	}

	ShareListParam struct {
		Creator string `json:"creator"`
		Limit   int64  `json:"limit"`
		Marker  string `json:"marker"`
	}

	ShareListResult struct {
		Items      []*shareEntityResult `json:"items"`
		NextMarker string               `json:"next_marker"`
	}

	ShareCancelResult struct {
		// 分享ID
		Id string
		// 是否成功
		Success bool
	}

	// FastShareCreateParam 创建快传分享
	FastShareCreateParam struct {
		DriveId    string   `json:"drive_id"`
		FileIdList []string `json:"file_id_list"`
	}

	FastShareFileItem struct {
		DriveId string `json:"drive_id"`
		FileId  string `json:"file_id"`
	}

	FastShareCreateResult struct {
		Expiration    string              `json:"expiration"`
		Thumbnail     string              `json:"thumbnail"`
		ShareName     string              `json:"share_name"`
		ShareId       string              `json:"share_id"`
		ShareUrl      string              `json:"share_url"`
		DriveFileList []FastShareFileItem `json:"drive_file_list"`
		FullShareMsg  string              `json:"full_share_msg"`
		ShareTitle    string              `json:"share_title"`
		ShareSubtitle string              `json:"share_subtitle"`
		Expired       bool                `json:"expired"`
	}
)

func createShareEntity(item *shareEntityResult) *aliyunpan.ShareEntity {
	if item == nil {
		return nil
	}
	return &aliyunpan.ShareEntity{
		Creator:    item.Creator,
		DriveId:    item.DriveId,
		ShareId:    item.ShareId,
		ShareName:  item.ShareName,
		SharePwd:   item.SharePwd,
		ShareUrl:   item.ShareUrl,
		FileIdList: item.FileIdList,
		SaveCount:  item.SaveCount,
		Status:     item.Status,
		Expiration: apiutil.UtcTime2LocalFormat(item.Expiration),
		UpdatedAt:  apiutil.UtcTime2LocalFormat(item.UpdatedAt),
		CreatedAt:  apiutil.UtcTime2LocalFormat(item.CreatedAt),
		FirstFile:  createFileEntity(item.FirstFile),
	}
}

// ShareLinkList 获取所有分享链接列表
func (p *WebPanClient) ShareLinkList(userId string) ([]*aliyunpan.ShareEntity, *apierror.ApiError) {
	resultList := []*aliyunpan.ShareEntity{}
	param := ShareListParam{
		Creator: userId,
		Limit:   100,
		Marker:  "",
	}
	for {
		if r, e := p.GetShareLinkListReq(param); e == nil {
			for _, item := range r.Items {
				resultList = append(resultList, createShareEntity(item))
			}

			// next page?
			if r.NextMarker != "" {
				param.Marker = r.NextMarker
				time.Sleep(500 * time.Millisecond)
			} else {
				break
			}
		} else {
			return nil, e
		}
	}
	return resultList, nil
}

// ShareLinkCancel 取消分享链接
func (p *WebPanClient) ShareLinkCancel(shareIdList []string) ([]*ShareCancelResult, *apierror.ApiError) {
	// url
	fullUrl := &strings.Builder{}
	fmt.Fprintf(fullUrl, "%s/adrive/v4/batch", API_URL)
	logger.Verboseln("do request url: " + fullUrl.String())

	// param
	pr := BatchRequestList{}
	for _, shareId := range shareIdList {
		pr = append(pr, &BatchRequest{
			Id:     shareId,
			Method: "POST",
			Url:    "/share_link/cancel",
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body: map[string]interface{}{
				"share_id": shareId,
			},
		})
	}

	batchParam := BatchRequestParam{
		Requests: pr,
		Resource: "file",
	}

	// request
	result, err := p.BatchTask(fullUrl.String(), &batchParam)
	if err != nil {
		logger.Verboseln("share cancel error ", err)
		return nil, apierror.NewFailedApiError(err.Error())
	}

	// parse result
	r := []*ShareCancelResult{}
	for _, item := range result.Responses {
		r = append(r, &ShareCancelResult{
			Id:      item.Id,
			Success: item.Status == 204,
		})
	}
	return r, nil
}

// ShareLinkCreate 创建分享
func (p *WebPanClient) ShareLinkCreate(param aliyunpan.ShareCreateParam) (*aliyunpan.ShareEntity, *apierror.ApiError) {
	// header
	header := map[string]string{
		"authorization": p.webToken.GetAuthorizationStr(),
	}

	// url
	fullUrl := &strings.Builder{}
	fmt.Fprintf(fullUrl, "%s/adrive/v2/share_link/create", API_URL)
	logger.Verboseln("do request url: " + fullUrl.String())

	// data
	postData := param

	// check pwd
	if postData.SharePwd != "" && len(postData.SharePwd) != 4 {
		return nil, apierror.NewFailedApiError("密码必须是4个字符")
	}

	// format time
	if postData.Expiration != "" {
		postData.Expiration = apiutil.LocalTime2UtcFormat(param.Expiration)
	}

	// request
	body, err := p.client.Fetch("POST", fullUrl.String(), postData, p.AddSignatureHeader(apiutil.AddCommonHeader(header)))
	if err != nil {
		logger.Verboseln("create share list error ", err)
		return nil, apierror.NewFailedApiError(err.Error())
	}
	logger.Verboseln("response: ", string(body))

	// handler common error
	if err1 := apierror.ParseCommonApiError(body); err1 != nil {
		return nil, err1
	}

	// parse result
	r := &shareEntityResult{}
	if err2 := json.Unmarshal(body, r); err2 != nil {
		logger.Verboseln("parse share create result json error ", err2)
		return nil, apierror.NewFailedApiError(err2.Error())
	}
	return createShareEntity(r), nil
}

func (p *WebPanClient) GetShareLinkListReq(param ShareListParam) (*ShareListResult, *apierror.ApiError) {
	// header
	header := map[string]string{
		"authorization": p.webToken.GetAuthorizationStr(),
	}

	// url
	fullUrl := &strings.Builder{}
	fmt.Fprintf(fullUrl, "%s/adrive/v3/share_link/list", API_URL)
	logger.Verboseln("do request url: " + fullUrl.String())

	if param.Limit <= 0 {
		param.Limit = 100
	}
	// data
	postData := map[string]interface{}{
		"category":         "file,album",
		"creator":          param.Creator,
		"include_canceled": false,
		"order_by":         "created_at",
		"order_direction":  "DESC",
		"limit":            param.Limit,
	}
	if param.Marker != "" {
		postData["marker"] = param.Marker
	}

	// request
	body, err := p.client.Fetch("POST", fullUrl.String(), postData, p.AddSignatureHeader(apiutil.AddCommonHeader(header)))
	logger.Verboseln(string(body))
	if err != nil {
		logger.Verboseln("get share list error ", err)
		return nil, apierror.NewFailedApiError(err.Error())
	}

	// handler common error
	if err1 := apierror.ParseCommonApiError(body); err1 != nil {
		return nil, err1
	}

	// parse result
	r := &ShareListResult{}
	if err2 := json.Unmarshal(body, r); err2 != nil {
		logger.Verboseln("parse share list result json error ", err2)
		return nil, apierror.NewFailedApiError(err2.Error())
	}
	return r, nil
}

// FastShareLinkCreate 创建快传分享
func (p *WebPanClient) FastShareLinkCreate(param FastShareCreateParam) (*FastShareCreateResult, *apierror.ApiError) {
	// header
	header := map[string]string{
		"authorization": p.webToken.GetAuthorizationStr(),
	}

	// url
	fullUrl := &strings.Builder{}
	fmt.Fprintf(fullUrl, "%s/adrive/v1/share/create", API_URL)
	logger.Verboseln("do request url: " + fullUrl.String())

	// data
	var fileList []FastShareFileItem
	for _, fileId := range param.FileIdList {
		fileList = append(fileList, FastShareFileItem{DriveId: param.DriveId, FileId: fileId})
	}
	postData := struct {
		DriveFileList []FastShareFileItem `json:"drive_file_list"`
	}{DriveFileList: fileList}

	// request
	body, err := p.client.Fetch("POST", fullUrl.String(), postData, p.AddSignatureHeader(apiutil.AddCommonHeader(header)))
	if err != nil {
		logger.Verboseln("create fast share list error ", err)
		return nil, apierror.NewFailedApiError(err.Error())
	}
	logger.Verboseln("response: ", string(body))

	// handler common error
	if err1 := apierror.ParseCommonApiError(body); err1 != nil {
		return nil, err1
	}

	// parse result
	r := &FastShareCreateResult{}
	if err2 := json.Unmarshal(body, r); err2 != nil {
		logger.Verboseln("parse fast share create result json error ", err2)
		return nil, apierror.NewFailedApiError(err2.Error())
	}
	r.Expiration = apiutil.UtcTime2LocalFormat(r.Expiration)
	return r, nil
}
