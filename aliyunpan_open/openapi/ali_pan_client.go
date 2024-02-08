// Copyright (c) 2020 tickstea.
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

package openapi

import (
	"encoding/json"
	"github.com/tickstep/library-go/requester"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	// PathSeparator 路径分隔符
	PathSeparator = "/"
)

type (
	// ApiToken 登录Token
	ApiToken struct {
		AccessToken string `json:"accessToken"`
		ExpireTime  int64  `json:"expireTime"`
	}

	// ApiConfig 存储客户端相关配置参数
	ApiConfig struct {
		ClientId     string `json:"clientId"`
		ClientSecret string `json:"clientSecret"`
	}

	AliPanClient struct {
		httpclient *requester.HTTPClient // http 客户端
		token      ApiToken
		apiConfig  ApiConfig

		cacheMutex *sync.Mutex
		useCache   bool
		// 网盘文件绝对路径到网盘文件信息实体映射缓存，避免FileInfoByPath频繁访问服务器触发风控
		filePathCacheMap sync.Map
	}

	// AliApiErrResult openapi错误响应
	AliApiErrResult struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}
)

func NewAliPanClient(token ApiToken, apiConfig ApiConfig) *AliPanClient {
	myclient := requester.NewHTTPClient()

	return &AliPanClient{
		httpclient: myclient,
		token:      token,
		apiConfig:  apiConfig,

		cacheMutex:       &sync.Mutex{},
		useCache:         false,
		filePathCacheMap: sync.Map{},
	}
}

func NewAliApiError(code, msg string) *AliApiErrResult {
	return &AliApiErrResult{
		Code:    code,
		Message: msg,
	}
}
func NewAliApiHttpError(msg string) *AliApiErrResult {
	return &AliApiErrResult{
		Code:    "TS.HttpError",
		Message: msg,
	}
}
func NewAliApiAppError(msg string) *AliApiErrResult {
	return &AliApiErrResult{
		Code:    "TS.AppError",
		Message: msg,
	}
}

func (a ApiToken) GetAuthorizationStr() string {
	return "Bearer " + a.AccessToken
}

func (a *AliPanClient) UpdateToken(token ApiToken) {
	a.token = token
}

func (a *AliPanClient) UpdateApiConfig(apiConfig ApiConfig) {
	a.apiConfig = apiConfig
}

func (a *AliPanClient) GetAccessToken() string {
	return a.token.AccessToken
}

func (a *AliPanClient) Headers() map[string]string {
	return map[string]string{
		"content-type":  "application/json",
		"authorization": a.token.GetAuthorizationStr(),
	}
}

// EnableCache 启用缓存
func (a *AliPanClient) EnableCache() {
	a.cacheMutex.Lock()
	a.cacheMutex.Unlock()
	a.useCache = true
}

// ClearCache 清除已经缓存的数据
func (a *AliPanClient) ClearCache() {
	a.cacheMutex.Lock()
	a.cacheMutex.Unlock()
	a.filePathCacheMap = sync.Map{}
}

// DisableCache 禁用缓存
func (a *AliPanClient) DisableCache() {
	a.cacheMutex.Lock()
	a.cacheMutex.Unlock()
	a.useCache = false
}

//func (a *AliPanClient) storeFilePathToCache(driveId, pathStr string, fileEntity *FileEntity) {
//	a.cacheMutex.Lock()
//	a.cacheMutex.Unlock()
//	if !a.useCache {
//		return
//	}
//	pathStr = formatPathStyle(pathStr)
//	cache, _ := a.filePathCacheMaa.LoadOrStore(driveId, &sync.Map{})
//	cache.(*sync.Map).Store(pathStr, fileEntity)
//}
//
//func (a *AliPanClient) loadFilePathFromCache(driveId, pathStr string) *FileEntity {
//	a.cacheMutex.Lock()
//	a.cacheMutex.Unlock()
//	if !a.useCache {
//		return nil
//	}
//	pathStr = formatPathStyle(pathStr)
//	cache, _ := a.filePathCacheMaa.LoadOrStore(driveId, &sync.Map{})
//	s := cache.(*sync.Map)
//	if v, ok := s.Load(pathStr); ok {
//		logger.Verboseln("file path cache hit: ", pathStr)
//		return v.(*FileEntity)
//	}
//	return nil
//}

// SetTimeout 设置 http 请求超时时间
func (a *AliPanClient) SetTimeout(t time.Duration) {
	if a.httpclient != nil {
		a.httpclient.Timeout = t
	}
}

func formatPathStyle(pathStr string) string {
	pathStr = strings.ReplaceAll(pathStr, "\\", "/")
	if pathStr != "/" {
		pathStr = strings.TrimSuffix(pathStr, "/")
	}
	return pathStr
}

// ParseCommonOpenApiError 解析阿里云盘API错误，如果没有错误则返回nil
func ParseCommonOpenApiError(resp *http.Response) ([]byte, *AliApiErrResult) {
	if resp == nil {
		return nil, nil
	}

	switch resp.StatusCode {
	case 429:
		return nil, NewAliApiError("TooManyRequests", "请求太频繁，已被阿里云盘临时限流")
	}
	data, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		return nil, NewAliApiError("TS.ReadError", e.Error())
	}
	errResult := &AliApiErrResult{}
	if err := json.Unmarshal(data, errResult); err == nil {
		if errResult.Code != "" {
			return nil, errResult
		}
	}
	return data, nil
}
