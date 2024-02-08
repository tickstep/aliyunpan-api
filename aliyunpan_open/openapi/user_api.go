package openapi

import (
	"encoding/json"
	"fmt"
	"github.com/tickstep/library-go/logger"
	"strings"
)

type (
	DriveInfoResult struct {
		// UserId 用户ID，具有唯一性
		UserId string `json:"user_id"`
		// Name 昵称
		Name string `json:"name"`
		// Avatar 头像地址
		Avatar string `json:"avatar"`
		// DefaultDriveId 默认drive
		DefaultDriveId string `json:"default_drive_id"`
		// ResourceDriveId 资源库。用户选择了授权才会返回
		ResourceDriveId string `json:"resource_drive_id"`
		// BackupDriveId 备份盘。用户选择了授权才会返回
		BackupDriveId string `json:"backup_drive_id"`
	}

	PersonalSpaceInfoResult struct {
		// UsedSize 使用容量，单位bytes
		UsedSize int64 `json:"used_size"`
		// TotalSize 总容量，单位bytes
		TotalSize int64 `json:"total_size"`
	}

	UserVipInfoResult struct {
		// Identity 枚举：member, vip, svip
		Identity   string `json:"identity"`
		PromotedAt string `json:"promotedAt"`
		// Expire 过期时间，时间戳，单位秒
		Expire int64 `json:"expire"`
	}

	UserScopeList []*UserScopeItem
	UserScopeItem struct {
		// Scope 权限标识
		Scope string `json:"scope"`
	}
)

// UserGetDriveInfo 获取用户drive信息
func (a *AliPanClient) UserGetDriveInfo() (*DriveInfoResult, *AliApiErrResult) {
	fullUrl := &strings.Builder{}
	fmt.Fprintf(fullUrl, "%s/adrive/v1.0/user/getDriveInfo", OPENAPI_URL)
	logger.Verboseln("do request url: " + fullUrl.String())

	// request
	resp, err := a.httpclient.Req("POST", fullUrl.String(), nil, a.Headers())
	if err != nil {
		logger.Verboseln("get drive info error ", err)
		return nil, NewAliApiHttpError(err.Error())
	}

	// handler common error
	var body []byte
	var apiErrResult *AliApiErrResult
	if body, apiErrResult = ParseCommonOpenApiError(resp); apiErrResult != nil {
		return nil, apiErrResult
	}

	// parse result
	r := &DriveInfoResult{}
	if err2 := json.Unmarshal(body, r); err2 != nil {
		logger.Verboseln("parse drive info result json error ", err2)
		return nil, NewAliApiAppError(err2.Error())
	}
	return r, nil
}

// UserGetSpaceInfo 获取用户空间信息
func (a *AliPanClient) UserGetSpaceInfo() (*PersonalSpaceInfoResult, *AliApiErrResult) {
	fullUrl := &strings.Builder{}
	fmt.Fprintf(fullUrl, "%s/adrive/v1.0/user/getSpaceInfo", OPENAPI_URL)
	logger.Verboseln("do request url: " + fullUrl.String())

	// request
	resp, err := a.httpclient.Req("POST", fullUrl.String(), nil, a.Headers())
	if err != nil {
		logger.Verboseln("get space info error ", err)
		return nil, NewAliApiHttpError(err.Error())
	}

	// handler common error
	var body []byte
	var apiErrResult *AliApiErrResult
	if body, apiErrResult = ParseCommonOpenApiError(resp); apiErrResult != nil {
		return nil, apiErrResult
	}

	// parse result
	type personalSpaceInfoData struct {
		Info *PersonalSpaceInfoResult `json:"personal_space_info"`
	}
	r := &personalSpaceInfoData{}
	if err2 := json.Unmarshal(body, r); err2 != nil {
		logger.Verboseln("parse space info result json error ", err2)
		return nil, NewAliApiAppError(err2.Error())
	}
	return r.Info, nil
}

// UserGetVipInfo 获取用户vip信息
func (a *AliPanClient) UserGetVipInfo() (*UserVipInfoResult, *AliApiErrResult) {
	fullUrl := &strings.Builder{}
	fmt.Fprintf(fullUrl, "%s/v1.0/user/getVipInfo", OPENAPI_URL)
	logger.Verboseln("do request url: " + fullUrl.String())

	// request
	resp, err := a.httpclient.Req("POST", fullUrl.String(), nil, a.Headers())
	if err != nil {
		logger.Verboseln("get vip info error ", err)
		return nil, NewAliApiHttpError(err.Error())
	}

	// handler common error
	var body []byte
	var apiErrResult *AliApiErrResult
	if body, apiErrResult = ParseCommonOpenApiError(resp); apiErrResult != nil {
		return nil, apiErrResult
	}

	// parse result
	r := &UserVipInfoResult{}
	if err2 := json.Unmarshal(body, r); err2 != nil {
		logger.Verboseln("parse vip info result json error ", err2)
		return nil, NewAliApiAppError(err2.Error())
	}
	return r, nil
}

// UserScopes 获取用户权限
func (a *AliPanClient) UserScopes() (*UserScopeList, *AliApiErrResult) {
	fullUrl := &strings.Builder{}
	fmt.Fprintf(fullUrl, "%s/oauth/users/scopes", OPENAPI_URL)
	logger.Verboseln("do request url: " + fullUrl.String())

	// request
	resp, err := a.httpclient.Req("GET", fullUrl.String(), nil, a.Headers())
	if err != nil {
		logger.Verboseln("get user scope info error ", err)
		return nil, NewAliApiHttpError(err.Error())
	}

	// handler common error
	var body []byte
	var apiErrResult *AliApiErrResult
	if body, apiErrResult = ParseCommonOpenApiError(resp); apiErrResult != nil {
		return nil, apiErrResult
	}

	// parse result
	type userScopeInfoData struct {
		Id         string         `json:"id"`
		UserScopes *UserScopeList `json:"scopes"`
	}
	r := &userScopeInfoData{}
	if err2 := json.Unmarshal(body, r); err2 != nil {
		logger.Verboseln("parse user scope info result json error ", err2)
		return nil, NewAliApiAppError(err2.Error())
	}
	return r.UserScopes, nil
}
