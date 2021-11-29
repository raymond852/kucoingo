package kucoingo

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	publicTokenEndPoint  = "https://api.kucoin.com/api/v1/bullet-public"
	privateTokenEndPoint = "https://api.kucoin.com/api/v1/bullet-private"
)

type SubscribeTopicRequest struct {
	Id             int    `json:"id"`
	Type           string `json:"type"`
	Topic          string `json:"topic"`
	Response       bool   `json:"response"`
	PrivateChannel bool   `json:"privateChannel"`
}

type SymbolTickerStreamResponse struct {
	Type    string                         `json:"type"`
	Topic   string                         `json:"topic"`
	Subject string                         `json:"subject"`
	Data    SymbolTickerStreamResponseData `json:"data"`
}

type SymbolTickerStreamResponseData struct {
	Sequence    string `json:"sequence"`
	Price       string `json:"price"`
	Size        string `json:"size"`
	BestAsk     string `json:"bestAsk"`
	BestAskSize string `json:"bestAskSize"`
	BestBid     string `json:"bestBid"`
	BestBidSize string `json:"bestBidSize"`
}

type webSocketConnection struct {
	rawConn     *websocket.Conn
	readHandler WebSocketReadHandler
}

type WebSocketReadHandler func(msgBytes []byte, err error) bool

type wsServerTokenResponse struct {
	Code string                         `json:"code"`
	Data wsServerTokenResponseFieldData `json:"data"`
}
type wsServerTokenResponseFieldData struct {
	InstanceServers []wsServerTokenResponseFieldDataFieldInstanceServers `json:"instanceServers"`
	Token           string                                               `json:"token"`
}

type wsServerTokenResponseFieldDataFieldInstanceServers struct {
	Endpoint     string `json:"endpoint"`
	Protocol     string `json:"protocol"`
	Encrypt      bool   `json:"encrypt"`
	PingInterval int    `json:"pingInterval"`
	PingTimeout  int    `json:"pingTimeout"`
}

func getPublicTokenWSEndpoint() (wsServerTokenResponse, error) {
	ret := wsServerTokenResponse{}
	if req, err := http.NewRequest(http.MethodPost, publicTokenEndPoint, nil); err != nil {
		return ret, err
	} else if resp, err := http.DefaultClient.Do(req); err != nil {
		return ret, err
	} else if bodyBytes, err := ioutil.ReadAll(resp.Body); err != nil {
		return ret, err
	} else if err := json.Unmarshal(bodyBytes, &ret); err != nil {
		return ret, err
	} else {
		return ret, err
	}
}

func getPrivateTokenWSEndpoint(kucoin Kucoin) (wsServerTokenResponse, error) {
	ret := wsServerTokenResponse{}
	if req, err := http.NewRequest(http.MethodPost, privateTokenEndPoint, nil); err != nil {
		return ret, err
	} else {
		kucoin.signRequest(req, time.Now(), req.URL.Path, nil)
		if resp, err := http.DefaultClient.Do(req); err != nil {
			return ret, err
		} else if bodyBytes, err := ioutil.ReadAll(resp.Body); err != nil {
			return ret, err
		} else if err := json.Unmarshal(bodyBytes, &ret); err != nil {
			return ret, err
		} else {
			return ret, err
		}
	}
}

func createConnection(wsTokenResp wsServerTokenResponse, readHandler WebSocketReadHandler) (*webSocketConnection, error) {
	dialer := &websocket.Dialer{}
	conn, _, err := dialer.Dial(wsTokenResp.Data.InstanceServers[0].Endpoint+"?token="+wsTokenResp.Data.Token, nil)
	if err != nil {
		return nil, err
	} else {
		ret := &webSocketConnection{
			rawConn:     conn,
			readHandler: readHandler,
		}
		ret.startHeartBeat(time.Duration(wsTokenResp.Data.InstanceServers[0].PingInterval) * time.Millisecond)
		if err := ret.handleWelcomeMsg(); err != nil {
			return nil, err
		}
		ret.startRead()
		return ret, nil
	}
}

func (ws *webSocketConnection) SubscribeTickerSymbol(symbol ...string) error {
	if len(symbol) == 0 {
		return BadRequestError{
			IncorrectFields: []string{"symbol"},
		}
	}
	topic := "/market/ticker:" + strings.Join(symbol, ",")
	return ws.subscribeTopic(topic, false)
}

func (ws *webSocketConnection) SubscribeAllTickers() error {
	return ws.subscribeTopic("/market/ticker:all", false)
}

func (ws *webSocketConnection) SubscribeOrderChangeEvent() error {
	return ws.subscribeTopic("/spotMarket/tradeOrders", true)
}

func (ws *webSocketConnection) subscribeTopic(topic string, isPrivate bool) error {
	req := SubscribeTopicRequest{
		Id:             int(time.Now().UnixMilli()),
		Type:           "subscribe",
		Topic:          topic,
		PrivateChannel: isPrivate,
	}
	return ws.rawConn.WriteJSON(req)
}

func (ws *webSocketConnection) handleHeartbeatMsg(bytes []byte) bool {
	payload := string(bytes)
	if strings.Contains(payload, "\"pong\"") {
		return true
	}
	return false
}

func (ws *webSocketConnection) startHeartBeat(heartBeatInterval time.Duration) {
	go func() {
		for {
			time.Sleep(heartBeatInterval)
			if err := ws.rawConn.WriteJSON(struct {
				Id   string `json:"id"`
				Type string `json:"type"`
			}{
				Id:   strconv.Itoa(int(time.Now().UnixMilli())),
				Type: "ping",
			}); err != nil {
				break
			}
		}
	}()
}

func (ws *webSocketConnection) handleWelcomeMsg() error {
	if _, bytes, err := ws.rawConn.ReadMessage(); err != nil {
		return err
	} else if !strings.Contains(string(bytes), "welcome") {
		return errors.New(fmt.Sprintf("Unregonized messge:%s, expecting welcome message", string(bytes)))
	} else {
		return nil
	}
}

func (ws *webSocketConnection) startRead() {
	go func() {
		for {
			_, bytes, err := ws.rawConn.ReadMessage()
			if !ws.handleHeartbeatMsg(bytes) {
				if continueRead := ws.readHandler(bytes, err); !continueRead {
					break
				}
			}
		}
	}()
}
