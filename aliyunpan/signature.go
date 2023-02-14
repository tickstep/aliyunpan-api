package aliyunpan

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tickstep/aliyunpan-api/aliyunpan/apierror"
	"github.com/tickstep/aliyunpan-api/aliyunpan/apiutil"
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

// CalcSignature 计算接口签名
func (p *PanClient) CalcSignature() error {
	max := 32
	key := randomString(max)
	data := fmt.Sprintf("%s:%s:%s:%d", p.appConfig.AppId, p.appConfig.DeviceId, p.appConfig.UserId, p.appConfig.Nonce)
	var privKey = secp256k1.PrivKey(key)
	pubKey := privKey.PubKey()
	p.appConfig.PublicKey = "04" + hex.EncodeToString(pubKey.Bytes())
	signature, err := privKey.Sign([]byte(data))
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
		logger.Verboseln("get file download url error ", err)
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
