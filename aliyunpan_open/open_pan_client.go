package aliyunpan_open

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tickstep/aliyunpan-api/aliyunpan_open/openapi"
	"github.com/tickstep/library-go/logger"
	"github.com/tickstep/library-go/requester"
	"strings"
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
	}
)

// NewOpenPanClient 创建开放接口客户端
func NewOpenPanClient(apiConfig openapi.ApiConfig, apiToken openapi.ApiToken, tokenCallback AccessTokenRefreshCallback) *OpenPanClient {
	myclient := requester.NewHTTPClient()

	return &OpenPanClient{
		httpClient:                 myclient,
		apiClient:                  openapi.NewAliPanClient(apiToken, apiConfig),
		accessTokenRefreshCallback: tokenCallback,
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
