package kucoingo

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"strconv"
	"time"
)

const (
	baseAPIMainURL = "https://api.kucoin.com"

	HeaderAPIKey        = "KC-API-KEY"
	HeaderAPISign       = "KC-API-SIGN"
	HeaderAPITimeStamp  = "KC-API-TIMESTAMP"
	HeaderAPIPassPhase  = "KC-API-PASSPHRASE"
	HeaderAPIKeyVersion = "KC-API-KEY-VERSION"
)

type Kucoin struct {
	apiKey           string
	apiSecretKey     string
	apiKeyPassPhrase string
	kucoinUrl        string
}

func NewKucoinClient(apiKey string, apiSecretKey string, apiKeyPassPhrase string) Kucoin {
	return Kucoin{
		apiKey:           apiKey,
		apiSecretKey:     apiSecretKey,
		apiKeyPassPhrase: apiKeyPassPhrase,
		kucoinUrl:        baseAPIMainURL,
	}
}

func (k Kucoin) NewPublicWebSocketConnection(readHandler WebSocketReadHandler) (*webSocketConnection, error) {
	wsTokenResp, err := getPublicTokenWSEndpoint()
	if err != nil {
		return nil, err
	}
	return createConnection(wsTokenResp, readHandler)
}

func (k Kucoin) NewPrivateWebSocketConnection(readHandler WebSocketReadHandler) (*webSocketConnection, error) {
	wsTokenResp, err := getPrivateTokenWSEndpoint(k)
	if err != nil {
		return nil, err
	}
	return createConnection(wsTokenResp, readHandler)
}

func (k Kucoin) MarketDataEndpoint() MarketEndpoint {
	return MarketEndpoint{
		conf: k,
	}
}

func (k Kucoin) GetAPIKey() string {
	return k.apiKey
}

func (k Kucoin) GetAPIPassPhrase() string {
	return k.apiKeyPassPhrase
}

func (k Kucoin) GetAPISecret() string {
	return k.apiSecretKey
}

func (k Kucoin) GetKucoinUrl() string {
	return k.kucoinUrl
}

func (k Kucoin) signRequest(httpReq *http.Request, t time.Time, path string, body *string) {
	timestamp := strconv.Itoa(int(t.UnixMilli()))
	h := hmac.New(sha256.New, []byte(k.apiSecretKey))
	b := ""
	if body != nil {
		b = *body
	}
	h.Write([]byte(k.apiKeyPassPhrase))
	httpReq.Header.Set(HeaderAPIPassPhase, base64.StdEncoding.EncodeToString(h.Sum(nil)))
	httpReq.Header.Set(HeaderAPIKeyVersion, "2")

	h.Reset()
	h.Write([]byte(timestamp + httpReq.Method + path + b))
	httpReq.Header.Set(HeaderAPIKey, k.apiKey)
	httpReq.Header.Set(HeaderAPISign, base64.StdEncoding.EncodeToString(h.Sum(nil)))
	httpReq.Header.Set(HeaderAPITimeStamp, timestamp)

}
