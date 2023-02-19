package aliyunpan

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/tickstep/aliyunpan-api/aliyunpan/apierror"
	"github.com/tickstep/aliyunpan-api/aliyunpan/apiutil"
	"github.com/tickstep/library-go/crypto/secp256k1"
	"github.com/tickstep/library-go/logger"
	"math/rand"
	"strings"
	"time"
)

type (
	CreateSessionParam struct {
		DeviceName string `json:"deviceName"`
		ModelName  string `json:"modelName"`
		PubKey     string `json:"pubKey"`
	}

	CreateSessionResult struct {
		Result  bool   `json:"result"`
		Success bool   `json:"success"`
		Code    string `json:"code"`
		Message string `json:"message"`
	}
)

const (
	NONCE_MIN = int32(0)
	NONCE_MAX = int32(2147483647)
)

func randomString(l int) []byte {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		rand.NewSource(time.Now().UnixNano())
		bytes[i] = byte(randInt(1, 2^256-1))
	}
	return bytes
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func getNextNonce(nonce int32) int32 {
	next := nonce + 1
	if next > NONCE_MAX {
		return NONCE_MIN
	} else {
		return next
	}
}

// CalcSignature 生成新的密钥并计算接口签名
func (p *PanClient) CalcSignature() error {
	max := 32
	key := randomString(max)
	p.appConfig.Nonce = 0
	data := fmt.Sprintf("%s:%s:%s:%d", p.appConfig.AppId, p.appConfig.DeviceId, p.appConfig.UserId, p.appConfig.Nonce)
	var privKey = secp256k1.PrivKey(key)
	p.appConfig.PrivKey = &privKey
	pubKey := privKey.PubKey()
	p.appConfig.PubKey = &pubKey
	p.appConfig.PublicKey = "04" + hex.EncodeToString(pubKey.Bytes())
	signature, err := privKey.Sign([]byte(data))
	if err != nil {
		return err
	}
	p.appConfig.SignatureData = hex.EncodeToString(signature) + "01"
	return nil
}

// CalcNextSignature 使用已有的密钥并生成新的签名
func (p *PanClient) CalcNextSignature() error {
	p.appConfig.Nonce = getNextNonce(p.appConfig.Nonce)
	data := fmt.Sprintf("%s:%s:%s:%d", p.appConfig.AppId, p.appConfig.DeviceId, p.appConfig.UserId, p.appConfig.Nonce)
	signature, err := p.appConfig.PrivKey.Sign([]byte(data))
	if err != nil {
		return err
	}
	p.appConfig.SignatureData = hex.EncodeToString(signature) + "01"
	return nil
}

// AddSignatureHeader 增加接口签名header
func (p *PanClient) AddSignatureHeader(headers map[string]string) map[string]string {
	if p.appConfig.PublicKey == "" {
		p.CalcSignature()
	}

	if headers == nil {
		return headers
	}

	// add signature
	//headers["x-canary"] = "client=web,app=adrive,version=v3.17.0"
	headers["x-device-id"] = p.appConfig.DeviceId
	headers["x-signature"] = p.appConfig.SignatureData
	return headers
}

// CreateSession 上传会话签名秘钥给服务器
func (p *PanClient) CreateSession(param *CreateSessionParam) (*CreateSessionResult, *apierror.ApiError) {
	// header
	header := map[string]string{
		"authorization": p.webToken.GetAuthorizationStr(),
	}

	// url
	fullUrl := &strings.Builder{}
	fmt.Fprintf(fullUrl, "%s/users/v1/users/device/create_session", API_URL)
	logger.Verboseln("do request url: " + fullUrl.String())

	// data
	if p.appConfig.PublicKey == "" {
		p.CalcSignature()
	}
	postData := map[string]interface{}{
		"deviceName": param.DeviceName,
		"modelName":  param.ModelName,
		"pubKey":     p.appConfig.PublicKey,
	}

	// request
	body, err := p.client.Fetch("POST", fullUrl.String(), postData, p.AddSignatureHeader(apiutil.AddCommonHeader(header)))
	if err != nil {
		logger.Verboseln("do create session error ", err)
		return nil, apierror.NewFailedApiError(err.Error())
	}

	// handler common error
	if err1 := apierror.ParseCommonApiError(body); err1 != nil {
		return nil, err1
	}

	// parse result
	r := &CreateSessionResult{}
	if err2 := json.Unmarshal(body, r); err2 != nil {
		logger.Verboseln("parse create session result json error ", err2)
		return nil, apierror.NewFailedApiError(err2.Error())
	}
	return r, nil
}

// RenewSession 刷新签名秘钥，如果刷新失败则需要调用CreateSession重新上传新秘钥
func (p *PanClient) RenewSession() (*CreateSessionResult, *apierror.ApiError) {
	// header
	header := map[string]string{
		"authorization": p.webToken.GetAuthorizationStr(),
	}

	// url
	fullUrl := &strings.Builder{}
	fmt.Fprintf(fullUrl, "%s/users/v1/users/device/renew_session", API_URL)
	logger.Verboseln("do request url: " + fullUrl.String())

	// request
	data := map[string]string{}
	body, err := p.client.Fetch("POST", fullUrl.String(), data, p.AddSignatureHeader(apiutil.AddCommonHeader(header)))
	if err != nil {
		logger.Verboseln("do renew session error ", err)
		return nil, apierror.NewFailedApiError(err.Error())
	}

	// handler common error
	if err1 := apierror.ParseCommonApiError(body); err1 != nil {
		return nil, err1
	}

	// parse result
	r := &CreateSessionResult{}
	if err2 := json.Unmarshal(body, r); err2 != nil {
		logger.Verboseln("parse renew session result json error ", err2)
		return nil, apierror.NewFailedApiError(err2.Error())
	}
	return r, nil
}
