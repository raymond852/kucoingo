package kucoingo

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type MarketEndpoint struct {
	conf Kucoin
}

type GetTradeHistoriesRequest struct {
	Symbol string
}

type GetTradeHistoriesResponse struct {
	Code string                          `json:"code"`
	Data []GetTradeHistoriesResponseItem `json:"data"`
}

type GetTradeHistoriesResponseItem struct {
	Sequence string `json:"sequence"`
	Price    string `json:"price"`
	Size     string `json:"size"`
	Side     string `json:"side"`
	Time     int    `json:"time"`
}

func (e MarketEndpoint) GetTradeHistories(ctx context.Context, req GetTradeHistoriesRequest) (*GetTradeHistoriesResponse, error) {
	path := "/api/v1/market/histories"
	httpReq, _ := http.NewRequestWithContext(ctx, http.MethodGet, e.conf.GetKucoinUrl()+path, nil)
	q := httpReq.URL.Query()
	q.Set("symbol", req.Symbol)
	httpReq.URL.RawQuery = q.Encode()

	e.conf.signRequest(httpReq, time.Now(), path, nil)

	var ret GetTradeHistoriesResponse
	if resp, err := http.DefaultClient.Do(httpReq); err != nil {
		return nil, err
	} else if bodyByte, err := ioutil.ReadAll(resp.Body); err != nil {
		return nil, err
	} else if resp.StatusCode != http.StatusOK {
		return nil, ResponseError{
			Body:       bodyByte,
			StatusCode: resp.StatusCode,
		}
	} else if err := json.Unmarshal(bodyByte, &ret); err != nil {
		fmt.Println(string(bodyByte))
		return nil, err
	} else {
		return &ret, nil
	}
}

type GetKLineRequest struct {
	Symbol  string
	StartAt *string
	EndAt   *string
	Type    EnumKLineType
}

type EnumKLineType string

const (
	EnumKlineType1min    = "1min"
	EnumKlineType3min     = "3min"
	EnumKlineType5min     = "5min"
	EnumKlineType15min    = "15min"
	EnumKlineType30min    = "30min"
	EnumKlineType1hour    = "1hour"
	EnumKlineType2hour    = "2hour"
	EnumKlineType4hour    = "4hour"
	EnumKlineType6hour     = "6hour"
	EnumKlineType8hour   = "8hour"
	EnumKlineType12hour   = "12hour"
	EnumKlineType1day  = "1day"
	EnumKlineType1week = "1week"
)

type GetKLineResponse struct {
	StartTime         int
	OpenPrice         float64
	ClosePrice        float64
	HighestPrice      float64
	LowestPrice       float64
	TransactionVolume float64
	TransactionAmount float64
}

func (e MarketEndpoint) GetKLine(ctx context.Context, req GetKLineRequest) ([]GetKLineResponse, error) {
	path := "/api/v1/market/candles"
	httpReq, _ := http.NewRequestWithContext(ctx, http.MethodGet, e.conf.GetKucoinUrl()+path, nil)
	q := httpReq.URL.Query()
	q.Set("symbol", req.Symbol)
	if req.StartAt != nil {
		q.Set("startAt", *req.StartAt)
	}
	if req.EndAt != nil {
		q.Set("endAt", *req.EndAt)
	}
	q.Set("type", string(req.Type))
	httpReq.URL.RawQuery = q.Encode()

	e.conf.signRequest(httpReq, time.Now(), path, nil)

	var respBody struct {
		Code int        `json:"string"`
		Data [][]string `json:"data"`
	}

	if resp, err := http.DefaultClient.Do(httpReq); err != nil {
		return nil, err
	} else if bodyByte, err := ioutil.ReadAll(resp.Body); err != nil {
		return nil, err
	} else if resp.StatusCode != http.StatusOK {
		return nil, ResponseError{
			Body:       bodyByte,
			StatusCode: resp.StatusCode,
		}
	} else if err := json.Unmarshal(bodyByte, &respBody); err != nil {
		fmt.Println(string(bodyByte))
		return nil, err
	} else {
		ret := make([]GetKLineResponse, len(respBody.Data))
		for i, item := range respBody.Data {
			retItem := GetKLineResponse{}
			if startTime, err := strconv.Atoi(item[0]); err != nil {
				return nil, err
			} else {
				retItem.StartTime = startTime
			}
			if openPrice, err := strconv.ParseFloat(item[1], 64); err != nil {
				return nil, err
			} else {
				retItem.OpenPrice = openPrice
			}
			if closePrice, err := strconv.ParseFloat(item[2], 64); err != nil {
				return nil, err
			} else {
				retItem.ClosePrice = closePrice
			}
			if highestPrice, err := strconv.ParseFloat(item[3], 64); err != nil {
				return nil, err
			} else {
				retItem.HighestPrice = highestPrice
			}
			if lowestPrice, err := strconv.ParseFloat(item[4], 64); err != nil {
				return nil, err
			} else {
				retItem.LowestPrice = lowestPrice
			}
			if transactionVolume, err := strconv.ParseFloat(item[5], 64); err != nil {
				return nil, err
			} else {
				retItem.TransactionVolume = transactionVolume
			}
			if transactionAmount, err := strconv.ParseFloat(item[6], 64); err != nil {
				return nil, err
			} else {
				retItem.TransactionAmount = transactionAmount
			}
			ret[i] = retItem
		}
		return ret, nil
	}
}
