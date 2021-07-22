package aliyunpan

import (
	"encoding/json"
	"fmt"
	"github.com/tickstep/aliyunpan-api/aliyunpan/apierror"
	"github.com/tickstep/aliyunpan-api/aliyunpan/apiutil"
	"github.com/tickstep/library-go/logger"
	"strings"
)

type(
	ShareEntity struct {
		Creator       string    `json:"creator"`
		DriveId       string    `json:"drive_id"`
		ShareId       string    `json:"share_id"`
		ShareName     string    `json:"share_name"`
		// SharePwd 密码，为空代表没有密码
		SharePwd      string    `json:"share_pwd"`
		ShareUrl      string    `json:"share_url"`
		FileIdList    []string  `json:"file_id_list"`
		SaveCount     int       `json:"save_count"`
		// Expiration 过期时间，为空代表永不过期
		Expiration    string `json:"expiration"`
		UpdatedAt     string `json:"updated_at"`
		CreatedAt     string `json:"created_at"`
	}

	shareEntityResult struct {
		CreatedAt     string `json:"created_at"`
		Creator       string    `json:"creator"`
		Description   string    `json:"description"`
		DownloadCount int       `json:"download_count"`
		DriveId       string    `json:"drive_id"`
		Expiration    string `json:"expiration"`
		Expired       bool      `json:"expired"`
		FileId        string    `json:"file_id"`
		FileIdList    []string  `json:"file_id_list"`
		PreviewCount  int       `json:"preview_count"`
		SaveCount     int       `json:"save_count"`
		ShareId       string    `json:"share_id"`
		ShareMsg      string    `json:"share_msg"`
		ShareName     string    `json:"share_name"`
		SharePolicy   string    `json:"share_policy"`
		SharePwd      string    `json:"share_pwd"`
		ShareUrl      string    `json:"share_url"`
		Status        string    `json:"status"`
		UpdatedAt     string `json:"updated_at"`
	}

	shareListResult struct {
		Items []*shareEntityResult `json:"items"`
		NextMarker string `json:"next_marker"`
	}

	ShareCancelResult struct {
		// 分享ID
		Id string
		// 是否成功
		Success bool
	}
)

// ShareList 获取分享链接列表
func (p *PanClient) ShareLinkList(userId string) ([]*ShareEntity, *apierror.ApiError) {
	resultList := []*ShareEntity{}
	if r,e := p.getShareLinkListReq(userId); e == nil {
		for _,item := range r.Items {
			resultList = append(resultList, &ShareEntity{
				Creator: item.Creator,
				DriveId: item.DriveId,
				ShareId: item.ShareId,
				ShareName: item.ShareName,
				SharePwd: item.SharePwd,
				ShareUrl: item.ShareUrl,
				FileIdList: item.FileIdList,
				SaveCount: item.SaveCount,
				Expiration: apiutil.UTCTimeFormat(item.Expiration),
				UpdatedAt: apiutil.UTCTimeFormat(item.UpdatedAt),
				CreatedAt: apiutil.UTCTimeFormat(item.CreatedAt),
			})
		}
	} else {
		return nil, e
	}
	return resultList, nil
}

// ShareLinkCancel 取消分享链接
func (p *PanClient) ShareLinkCancel(shareIdList []string) ([]*ShareCancelResult, *apierror.ApiError) {
	// url
	fullUrl := &strings.Builder{}
	fmt.Fprintf(fullUrl, "%s/adrive/v2/batch", API_URL)
	logger.Verboseln("do request url: " + fullUrl.String())

	// param
	pr := BatchRequestList{}
	for _,shareId := range shareIdList {
		pr = append(pr, &BatchRequest{
			Id:      shareId,
			Method:  "POST",
			Url:     "/share_link/cancel",
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body:    map[string]interface{}{
				"share_id": shareId,
			},
		})
	}

	batchParam := BatchRequestParam{
		Requests: pr,
		Resource: "file",
	}

	// request
	result,err := p.BatchTask(fullUrl.String(), &batchParam)
	if err != nil {
		logger.Verboseln("share cancel error ", err)
		return nil, apierror.NewFailedApiError(err.Error())
	}

	// parse result
	r := []*ShareCancelResult{}
	for _,item := range result.Responses{
		r = append(r, &ShareCancelResult{
			Id: item.Id,
			Success: item.Status == 204,
		})
	}
	return r, nil
}

func (p *PanClient) getShareLinkListReq(userId string) (*shareListResult, *apierror.ApiError) {
	// header
	header := map[string]string {
		"authorization": p.webToken.GetAuthorizationStr(),
	}

	// url
	fullUrl := &strings.Builder{}
	fmt.Fprintf(fullUrl, "%s/adrive/v2/share_link/list", API_URL)
	logger.Verboseln("do request url: " + fullUrl.String())

	// data
	postData := map[string]interface{} {
		"creator": userId,
		"include_canceled": false,
		"order_by": "created_at",
		"order_direction": "DESC",
	}

	// request
	body, err := client.Fetch("POST", fullUrl.String(), postData, apiutil.AddCommonHeader(header))
	if err != nil {
		logger.Verboseln("get share list error ", err)
		return nil, apierror.NewFailedApiError(err.Error())
	}

	// handler common error
	if err1 := apierror.ParseCommonApiError(body); err1 != nil {
		return nil, err1
	}

	// parse result
	r := &shareListResult{}
	if err2 := json.Unmarshal(body, r); err2 != nil {
		logger.Verboseln("parse share list result json error ", err2)
		return nil, apierror.NewFailedApiError(err2.Error())
	}
	return r, nil
}