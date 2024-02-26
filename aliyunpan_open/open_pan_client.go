package aliyunpan_open

import (
	"github.com/tickstep/aliyunpan-api/aliyunpan_open/openapi"
	"github.com/tickstep/library-go/requester"
	"time"
)

type (
	// OpenPanClient 开放接口客户端
	OpenPanClient struct {
		httpClient *requester.HTTPClient // http 客户端
		apiClient  *openapi.AliPanClient
	}
)

// NewOpenPanClient 创建开放接口客户端
func NewOpenPanClient(apiConfig openapi.ApiConfig, apiToken openapi.ApiToken) *OpenPanClient {
	myclient := requester.NewHTTPClient()

	return &OpenPanClient{
		httpClient: myclient,
		apiClient:  openapi.NewAliPanClient(apiToken, apiConfig),
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
