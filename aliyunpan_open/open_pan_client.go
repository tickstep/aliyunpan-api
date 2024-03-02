package aliyunpan_open

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tickstep/aliyunpan-api/aliyunpan"
	"github.com/tickstep/aliyunpan-api/aliyunpan_open/openapi"
	"github.com/tickstep/library-go/logger"
	"github.com/tickstep/library-go/requester"
	"strings"
	"sync"
	"time"
)

const (
	// ApiRetryMaxTimes API失败重试次数
	ApiRetryMaxTimes int = 3
)

type (
	// AccessTokenRefreshCallback Token刷新回调
	AccessTokenRefreshCallback func(newToken openapi.ApiToken) error

	// OpenPanClient 开放接口客户端
	OpenPanClient struct {
		httpClient *requester.HTTPClient // http 客户端
		apiClient  *openapi.AliPanClient

		accessTokenRefreshCallback AccessTokenRefreshCallback

		// 缓存
		cacheMutex *sync.Mutex
		useCache   bool
		// 网盘文件绝对路径到网盘文件信息实体映射缓存，避免FileInfoByPath频繁访问服务器触发风控
		filePathCacheMap sync.Map
	}
)

// NewOpenPanClient 创建开放接口客户端
func NewOpenPanClient(apiConfig openapi.ApiConfig, apiToken openapi.ApiToken, tokenCallback AccessTokenRefreshCallback) *OpenPanClient {
	myclient := requester.NewHTTPClient()

	return &OpenPanClient{
		httpClient:                 myclient,
		apiClient:                  openapi.NewAliPanClient(apiToken, apiConfig),
		accessTokenRefreshCallback: tokenCallback,
		cacheMutex:                 &sync.Mutex{},
		useCache:                   false,
		filePathCacheMap:           sync.Map{},
	}
}

// SetTimeout 设置 http 请求超时时间
func (p *OpenPanClient) SetTimeout(t time.Duration) {
	if p.apiClient != nil {
		p.apiClient.SetTimeout(t)
	}

	if p.httpClient != nil {
		p.httpClient.Timeout = t
	}
}

// GetAccessToken 获取AccessToken鉴权字符串
func (p *OpenPanClient) GetAccessToken() string {
	return p.apiClient.GetAccessToken()
}

// RefreshNewAccessToken 获取新的AccessToken
func (p *OpenPanClient) RefreshNewAccessToken() error {
	fullUrl := &strings.Builder{}
	fmt.Fprintf(fullUrl, "https://api.tickstep.com/auth/tickstep/aliyunpan/token/openapi/%s/refresh?userId=%s",
		p.apiClient.GetApiConfig().TicketId, p.apiClient.GetApiConfig().UserId)
	logger.Verboseln("do request url: " + fullUrl.String())

	// request
	data, err := p.httpClient.Fetch("GET", fullUrl.String(), nil, p.apiClient.Headers())
	if err != nil {
		logger.Verboseln("get new access token error ", err)
		return err
	}

	// parse result
	type respEntity struct {
		Code int               `json:"code"`
		Data *openapi.ApiToken `json:"data"`
		Msg  string            `json:"msg"`
	}
	r := &respEntity{}
	if err2 := json.Unmarshal(data, r); err2 != nil {
		logger.Verboseln("parse access token result json error ", err2)
		return err2
	}
	if r.Code != 0 {
		return errors.New(r.Msg)
	}
	token := *r.Data
	p.apiClient.UpdateToken(token)
	if p.accessTokenRefreshCallback != nil {
		p.accessTokenRefreshCallback(token)
	}
	return nil
}

// UpdateUserId 更新用户ID
func (p *OpenPanClient) UpdateUserId(userId string) {
	c := p.apiClient.GetApiConfig()
	c.UserId = userId
	p.apiClient.UpdateApiConfig(c)
}

// EnableCache 启用缓存
func (p *OpenPanClient) EnableCache() {
	p.cacheMutex.Lock()
	p.cacheMutex.Unlock()
	p.useCache = true
}

// ClearCache 清除已经缓存的数据
func (p *OpenPanClient) ClearCache() {
	p.cacheMutex.Lock()
	p.cacheMutex.Unlock()
	p.filePathCacheMap = sync.Map{}
}

// DisableCache 禁用缓存
func (p *OpenPanClient) DisableCache() {
	p.cacheMutex.Lock()
	p.cacheMutex.Unlock()
	p.useCache = false
}

// storeFilePathToCache 存储文件信息到缓存
func (p *OpenPanClient) storeFilePathToCache(driveId, pathStr string, fileEntity *aliyunpan.FileEntity) {
	p.cacheMutex.Lock()
	p.cacheMutex.Unlock()
	if !p.useCache {
		return
	}
	pathStr = formatPathStyle(pathStr)
	cache, _ := p.filePathCacheMap.LoadOrStore(driveId, &sync.Map{})
	cache.(*sync.Map).Store(pathStr, fileEntity)
}

// loadFilePathFromCache 从缓存获取文件信息
func (p *OpenPanClient) loadFilePathFromCache(driveId, pathStr string) *aliyunpan.FileEntity {
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
		return v.(*aliyunpan.FileEntity)
	}
	return nil
}

func formatPathStyle(pathStr string) string {
	pathStr = strings.ReplaceAll(pathStr, "\\", "/")
	if pathStr != "/" {
		pathStr = strings.TrimSuffix(pathStr, "/")
	}
	return pathStr
}
