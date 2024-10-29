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
	"github.com/tickstep/library-go/requester"
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
		ExpiredAt   int64  `json:"expired"`
	}

	// ApiConfig 存储客户端相关配置参数
	ApiConfig struct {
		TicketId     string `json:"ticket_id"`
		UserId       string `json:"user_id"`
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
		//"X-Canary":      "label=gray", // 标记灰度测试header
	}
}

func (a *AliPanClient) GetApiConfig() ApiConfig {
	return a.apiConfig
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
