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
	"github.com/tickstep/library-go/crypto"
	"github.com/tickstep/library-go/crypto/secp256k1"
	"github.com/tickstep/library-go/logger"
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
	// AppConfig 存储客户端相关配置参数，目前主要是签名需要用的参数
	AppConfig struct {
		AppId string `json:"appId"`
		// DeviceId标识登录客户端，阿里限制：为了保障你的数据隐私安全，阿里云盘最多只允许你同时登录 10 台设备。你已超出最大设备数量，请先选择一台设备下线，才可以继续使用
		DeviceId      string `json:"deviceId"`
		UserId        string `json:"userId"`
		Nonce         int32  `json:"nonce"`
		PublicKey     string `json:"publicKey"`
		SignatureData string `json:"signatureData"`

		PrivKey *secp256k1.PrivKey `json:"-"`
		PubKey  *crypto.PubKey     `json:"-"`
	}

	PanClient struct {
		client    *requester.HTTPClient // http 客户端
		webToken  WebLoginToken
		appToken  AppLoginToken
		appConfig AppConfig

		cacheMutex *sync.Mutex
		useCache   bool
		// 网盘文件绝对路径到网盘文件信息实体映射缓存，避免FileInfoByPath频繁访问服务器触发风控
		filePathCacheMap sync.Map
	}
)

func NewPanClient(webToken WebLoginToken, appToken AppLoginToken, appConfig AppConfig) *PanClient {
	myclient := requester.NewHTTPClient()

	return &PanClient{
		client:           myclient,
		webToken:         webToken,
		appToken:         appToken,
		appConfig:        appConfig,
		cacheMutex:       &sync.Mutex{},
		useCache:         false,
		filePathCacheMap: sync.Map{},
	}
}

func (p *PanClient) UpdateToken(webToken WebLoginToken) {
	p.webToken = webToken
}

func (p *PanClient) UpdateAppConfig(appConfig AppConfig) {
	p.appConfig = appConfig
}

func (p *PanClient) GetAccessToken() string {
	return p.webToken.AccessToken
}

// EnableCache 启用缓存
func (p *PanClient) EnableCache() {
	p.cacheMutex.Lock()
	p.cacheMutex.Unlock()
	p.useCache = true
}

// ClearCache 清除已经缓存的数据
func (p *PanClient) ClearCache() {
	p.cacheMutex.Lock()
	p.cacheMutex.Unlock()
	p.filePathCacheMap = sync.Map{}
}

// DisableCache 禁用缓存
func (p *PanClient) DisableCache() {
	p.cacheMutex.Lock()
	p.cacheMutex.Unlock()
	p.useCache = false
}

func (p *PanClient) storeFilePathToCache(driveId, pathStr string, fileEntity *FileEntity) {
	p.cacheMutex.Lock()
	p.cacheMutex.Unlock()
	if !p.useCache {
		return
	}
	pathStr = formatPathStyle(pathStr)
	cache, _ := p.filePathCacheMap.LoadOrStore(driveId, &sync.Map{})
	cache.(*sync.Map).Store(pathStr, fileEntity)
}

func (p *PanClient) loadFilePathFromCache(driveId, pathStr string) *FileEntity {
	p.cacheMutex.Lock()
	p.cacheMutex.Unlock()
	if !p.useCache {
		return nil
	}
	pathStr = formatPathStyle(pathStr)
	cache, _ := p.filePathCacheMap.LoadOrStore(driveId, &sync.Map{})
	s := cache.(*sync.Map)
	if v, ok := s.Load(pathStr); ok {
		logger.Verboseln("file path cache hit: ", pathStr)
		return v.(*FileEntity)
	}
	return nil
}

// SetTimeout 设置 http 请求超时时间
func (p *PanClient) SetTimeout(t time.Duration) {
	if p.client != nil {
		p.client.Timeout = t
	}
}

func formatPathStyle(pathStr string) string {
	pathStr = strings.ReplaceAll(pathStr, "\\", "/")
	if pathStr != "/" {
		pathStr = strings.TrimSuffix(pathStr, "/")
	}
	return pathStr
}
