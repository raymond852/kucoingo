package kucoingo

import (
	"context"
	"fmt"
	"os"
	"testing"
)

var APIKey = os.Getenv("KUCOIN_API_KEY")
var APISecret = os.Getenv("KUCOIN_API_SECRET")
var APIPassPhrase = os.Getenv("KUCOIN_API_PASS_PHRASE")

func TestEndpoint_GetTradeHistories(t *testing.T) {
	kubcoin := NewKucoinClient(APIKey, APISecret, APIPassPhrase)
	resp, err := NewMarketEndPoint(kubcoin).GetTradeHistories(context.Background(), GetTradeHistoriesRequest{
		Symbol: "NUM-USDT",
	})
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(resp)
	}
}

func TestMarketEndpoint_GetKLine(t *testing.T) {
	kubcoin := NewKucoinClient(APIKey, APISecret, APIPassPhrase)
	resp, err := NewMarketEndPoint(kubcoin).GetKLine(context.Background(), GetKLineRequest{
		Symbol: "NUM-USDT",
		Type:   ENUM_KLINE_TYPE_1MIN,
	})
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(resp)
	}
}
