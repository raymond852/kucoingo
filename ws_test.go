package kucoingo

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestWebSocketConnection_SubscribeTickerSymbol(t *testing.T) {
	ws, err := Kucoin{}.NewPublicWebSocketConnection(func(msgBytes []byte, err error) bool {
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(string(msgBytes))
		}
		return true
	})
	if err != nil {
		fmt.Println(err)
	} else {
		ws.SubscribeTickerSymbol("BTC-USDT")
		time.Sleep(10 * time.Minute)
	}
}

func TestWebSocketConnection_SubscribeAllTickers(t *testing.T) {
	ws, err := Kucoin{}.NewPublicWebSocketConnection(func(msgBytes []byte, err error) bool {
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(string(msgBytes))
		}
		return true
	})
	if err != nil {
		fmt.Println(err)
	} else {
		ws.SubscribeAllTickers()
		time.Sleep(10 * time.Minute)
	}
}

func TestWebSocketConnection_SubscribeOrderChangeEvent(t *testing.T) {
	var APIKey = os.Getenv("KUCOIN_API_KEY")
	var APISecret = os.Getenv("KUCOIN_API_SECRET")
	var APIPassPhrase = os.Getenv("KUCOIN_API_PASS_PHRASE")

	ws, err := NewKucoinClient(APIKey, APISecret, APIPassPhrase).NewPrivateWebSocketConnection(func(msgBytes []byte, err error) bool {
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(string(msgBytes))
		}
		return true
	})
	if err != nil {
		fmt.Println(err)
	} else if err := ws.SubscribeOrderChangeEvent(); err != nil {
		fmt.Println(err)
	} else {
		time.Sleep(10 * time.Minute)
	}
}