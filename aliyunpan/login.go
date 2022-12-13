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

// WEB端API
package aliyunpan

import (
	"encoding/json"
	"fmt"
	"github.com/tickstep/aliyunpan-api/aliyunpan/apierror"
	"github.com/tickstep/aliyunpan-api/aliyunpan/apiutil"
	"github.com/tickstep/library-go/logger"
	"github.com/tickstep/library-go/requester"
	"strings"
	"time"
)

const (
)

type (
	refreshTokenResult struct {
		AccessToken string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn int `json:"expires_in"`
		TokenType string `json:"token_type"`
		UserId string `json:"user_id"`
		UserName string `json:"user_name"`
		NickName string `json:"nick_name"`
		DefaultDriveId string `json:"default_drive_id"`
		DefaultSboxDriveId string `json:"default_sbox_drive_id"`
		Role string `json:"role"`
		Status string `json:"status"`
		ExpireTime string `json:"expire_time"`
		DeviceId string `json:"device_id"`
	}

	WebLoginToken struct {
		AccessTokenType string `json:"accessTokenType"`
		AccessToken string `json:"accessToken"`
		RefreshToken string `json:"refreshToken"`
		ExpiresIn int `json:"expiresIn"`
		ExpireTime string `json:"expireTime"`
	}
)

func (w *WebLoginToken) GetAuthorizationStr() string {
	return w.AccessTokenType + " " + w.AccessToken
}

func (w *WebLoginToken) IsAccessTokenExpired() bool {
	local, _ := time.LoadLocation("Local")
	expireTime, _ := time.ParseInLocation("2006-01-02 15:04:05", w.ExpireTime, local)
	now := time.Now()

	return (expireTime.Unix() - now.Unix()) < 60
}

func GetAccessTokenFromRefreshToken(refreshToken string) (*WebLoginToken, *apierror.ApiError) {
	myclient := requester.NewHTTPClient()

	header := map[string]string {}

	fullUrl := &strings.Builder{}
	fmt.Fprintf(fullUrl, "%s/v2/account/token", AUTH_URL)
	logger.Verboseln("do request url: " + fullUrl.String())
	postData := map[string]string {
		"refresh_token": refreshToken,
		"api_id": "pJZInNHN2dZWk8qg",
		"grant_type": "refresh_token",
	}

	body, err := myclient.Fetch("POST", fullUrl.String(), postData, apiutil.AddCommonHeader(header))
	if err != nil {
		logger.Verboseln("get access token error ", err)
		return nil, apierror.NewFailedApiError(err.Error())
	}

	// handler common error
	if err1 := apierror.ParseCommonApiError(body); err1 != nil {
		return nil, err1
	}

	r := &refreshTokenResult{}
	if err1 := json.Unmarshal(body, r); err1 != nil {
		logger.Verboseln("parse refresh token result json error ", err1)
		return nil, apierror.NewFailedApiError(err1.Error())
	}

	result := &WebLoginToken{
		r.TokenType,
		r.AccessToken,
		r.RefreshToken,
		r.ExpiresIn,
		apiutil.UtcTime2LocalFormat(r.ExpireTime),
	}
	return result, nil
}