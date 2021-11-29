package kucoingo

import (
	"context"
	"fmt"
	"os"
	"testing"
)



func TestEndpoint_GetTradeHistories(t *testing.T) {
	var APIKey = os.Getenv("KUCOIN_API_KEY")
	var APISecret = os.Getenv("KUCOIN_API_SECRET")
	var APIPassPhrase = os.Getenv("KUCOIN_API_PASS_PHRASE")

	kubcoin := NewKucoinClient(APIKey, APISecret, APIPassPhrase)
	resp, err := kubcoin.MarketDataEndpoint().GetTradeHistories(context.Background(), GetTradeHistoriesRequest{
		Symbol: "NUM-USDT",
	})
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(resp)
	}
}

func TestMarketEndpoint_GetKLine(t *testing.T) {
	var APIKey = os.Getenv("KUCOIN_API_KEY")
	var APISecret = os.Getenv("KUCOIN_API_SECRET")
	var APIPassPhrase = os.Getenv("KUCOIN_API_PASS_PHRASE")

	kubcoin := NewKucoinClient(APIKey, APISecret, APIPassPhrase)
	resp, err := kubcoin.MarketDataEndpoint().GetKLine(context.Background(), GetKLineRequest{
		Symbol:  "BASIC-USDT",
		StartAt: StringPointer("1636822800"),
		EndAt:   StringPointer("1636823400"),
		Type:    EnumKlineType1min,
	})
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(resp)
	}
}
