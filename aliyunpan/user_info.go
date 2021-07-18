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
	"github.com/tickstep/library-go/logger"
	"strings"
	"time"
)

type (
	UserRole string
	UserStatus string

	// UserInfo 用户信息
	UserInfo struct {
		// DomainId 域ID
		DomainId string `json:"domainId"`
		// DefaultDriveId 网盘ID
		DefaultDriveId string `json:"defaultDriveId"`
		// 用户UID
		UserId string `json:"userId"`
		// UserName 用户名
		UserName       string `json:"userName"`
		// CreatedAt 创建时间
		CreatedAt      string  `json:"createdAt"`
		// Email 邮箱
		Email          string `json:"email"`
		// Phone 手机
		Phone          string `json:"phone"`
		// Role 角色，默认是user
		Role           UserRole `json:"role"`
		// Status 是否被禁用，enable / disable
		Status         UserStatus `json:"status"`
		// Nickname 昵称，如果没有设置则为空
		Nickname string `json:"nickname"`
		// TotalSize 网盘空间总大小
		TotalSize uint64 `json:"totalSize"`
		// UsedSize 网盘已使用空间大小
		UsedSize uint64 `json:"usedSize"`
	}

	// userInfoResult 用户信息返回实体
	userInfoResult struct {
		DomainID       string `json:"domain_id"`
		UserID         string `json:"user_id"`
		Avatar         string `json:"avatar"`
		CreatedAt      int64  `json:"created_at"`
		UpdatedAt      int64  `json:"updated_at"`
		Email          string `json:"email"`
		NickName       string `json:"nick_name"`
		Phone          string `json:"phone"`
		Role           string `json:"role"`
		Status         string `json:"status"`
		UserName       string `json:"user_name"`
		Description    string `json:"description"`
		DefaultDriveID string `json:"default_drive_id"`
		DenyChangePasswordBySelf    bool `json:"deny_change_password_by_self"`
		NeedChangePasswordNextLogin bool `json:"need_change_password_next_login"`
	}

	personalInfoResult struct {
		// 权限
		PersonalRightsInfo struct {
			SpuID      string `json:"spu_id"`
			Name       string `json:"name"`
			IsExpires  bool   `json:"is_expires"`
			Privileges []struct {
				FeatureID     string `json:"feature_id"`
				FeatureAttrID string `json:"feature_attr_id"`
				Quota         int    `json:"quota"`
			} `json:"privileges"`
		} `json:"personal_rights_info"`

		// quota配额
		PersonalSpaceInfo struct {
			UsedSize  uint64   `json:"used_size"`
			TotalSize uint64 `json:"total_size"`
		} `json:"personal_space_info"`
	}
)

const (
	User UserRole = "user"
	UnknownRole UserRole = "unknown"

	Enabled UserStatus = "enable"
	UnknownStatus UserStatus = "unknown"
)

func parseUserRole(role string) UserRole {
	switch role {
	case "user":
		return User
	}
	return UnknownRole
}

func parseUserStatus(status string) UserStatus {
	switch status {
	case "enabled":
		return Enabled
	}
	return UnknownStatus
}

// GetUserInfo 获取用户信息
func (p *PanClient) GetUserInfo() (*UserInfo, *apierror.ApiError) {
	userInfo := &UserInfo{}

	if r,err := p.getUserInfoReq(); err == nil {
		userInfo.DomainId = r.DomainID
		userInfo.DefaultDriveId = r.DefaultDriveID
		userInfo.UserId = r.UserID
		userInfo.UserName = r.UserName
		userInfo.CreatedAt = time.Unix(r.CreatedAt / 1000, 0).Format("2006-01-02 15:04:05")
		userInfo.Email = r.Email
		userInfo.Phone = r.Email
		userInfo.Role = parseUserRole(r.Role)
		userInfo.Status = parseUserStatus(r.Status)
		userInfo.Nickname = r.NickName
	} else {
		return nil, err
	}

	if r,err := p.getPersonalInfoReq(); err == nil {
		userInfo.TotalSize = r.PersonalSpaceInfo.TotalSize
		userInfo.UsedSize = r.PersonalSpaceInfo.UsedSize
	} else {
		return nil, err
	}
	return userInfo, nil
}

// getUserInfoReq 获取用户基本信息
func (p *PanClient) getUserInfoReq() (*userInfoResult, *apierror.ApiError) {
	header := map[string]string {
		"accept": "application/json, text/plain, */*",
		"referer": "https://www.aliyundrive.com/",
		"origin": "https://www.aliyundrive.com",
		"content-type": "application/json;charset=UTF-8",
		"user-agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
		"authorization": p.webToken.GetAuthorizationStr(),
	}

	fullUrl := &strings.Builder{}
	fmt.Fprintf(fullUrl, "%s/v2/user/get", API_URL)
	logger.Verboseln("do request url: " + fullUrl.String())
	postData := map[string]string {}

	// request
	body, err := client.Fetch("POST", fullUrl.String(), postData, header)
	if err != nil {
		logger.Verboseln("get user info error ", err)
		return nil, apierror.NewFailedApiError(err.Error())
	}

	// handler common error
	if err1 := apierror.ParseCommonApiError(body); err1 != nil {
		return nil, err1
	}

	// parse result
	r := &userInfoResult{}
	if err2 := json.Unmarshal(body, r); err2 != nil {
		logger.Verboseln("parse user info result json error ", err2)
		return nil, apierror.NewFailedApiError(err2.Error())
	}
	return r, nil
}

// getPersonalInfoReq 获取用户网盘基本信息，包括配额，上传下载等权限限制
func (p *PanClient) getPersonalInfoReq() (*personalInfoResult, *apierror.ApiError) {
	header := map[string]string {
		"accept": "application/json, text/plain, */*",
		"referer": "https://www.aliyundrive.com/",
		"origin": "https://www.aliyundrive.com",
		"content-type": "application/json;charset=UTF-8",
		"user-agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
		"authorization": p.webToken.GetAuthorizationStr(),
	}

	fullUrl := &strings.Builder{}
	fmt.Fprintf(fullUrl, "%s/v2/databox/get_personal_info", API_URL)
	logger.Verboseln("do request url: " + fullUrl.String())
	postData := map[string]string {}

	// request
	body, err := client.Fetch("POST", fullUrl.String(), postData, header)
	if err != nil {
		logger.Verboseln("get person info error ", err)
		return nil, apierror.NewFailedApiError(err.Error())
	}

	// handler common error
	if err1 := apierror.ParseCommonApiError(body); err1 != nil {
		return nil, err1
	}

	// parse result
	r := &personalInfoResult{}
	if err2 := json.Unmarshal(body, r); err2 != nil {
		logger.Verboseln("parse person info result json error ", err2)
		return nil, apierror.NewFailedApiError(err2.Error())
	}
	return r, nil
}
