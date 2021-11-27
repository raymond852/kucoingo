package kucoingo

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
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

func NewMarketEndPoint(conf Kucoin) MarketEndpoint {
	return MarketEndpoint{
		conf: conf,
	}
}

type GetTradeHistoriesRequest struct {
	Symbol string
}

type GetTradeHistoriesResponse struct {
	Code string                           `json:"code"`
	Data []GetTradeHistoriesResponse_Item `json:"data"`
}

type GetTradeHistoriesResponse_Item struct {
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
	timestamp := strconv.Itoa(int(time.Now().Unix()))
	h := hmac.New(sha256.New, []byte(e.conf.GetAPISecret()))
	h.Write([]byte(timestamp + http.MethodGet + path))
	httpReq.Header.Set(HeaderAPIKey, e.conf.GetAPIKey())
	httpReq.Header.Set(HeaderAPISign, base64.StdEncoding.EncodeToString(h.Sum(nil)))
	httpReq.Header.Set(HeaderAPITimeStamp, timestamp)

	h1 := hmac.New(sha256.New, []byte(e.conf.GetAPISecret()))
	h1.Write([]byte(e.conf.GetAPIPassPhrase()))
	httpReq.Header.Set(HeaderAPIPassPhase, base64.StdEncoding.EncodeToString(h.Sum(nil)))
	httpReq.Header.Set(HeaderAPIKeyVersion, "V2")
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
	ENUM_KLINE_TYPE_1MIN   = "1min"
	ENUM_KLINE_TYPE_3MIN   = "3min"
	ENUM_KLINE_TYPE_5MIN   = "5min"
	ENUM_KLINE_TYPE_15MIN  = "15min"
	ENUM_KLINE_TYPE_30MIN  = "30min"
	ENUM_KLINE_TYPE_1HOUR  = "1hour"
	ENUM_KLINE_TYPE_2HOUR  = "2hour"
	ENUM_KLINE_TYPE_4HOUR  = "4hour"
	ENUM_KLINE_TYPE_6HOUR  = "6hour"
	ENUM_KLINE_TYPE_8HOUR  = "8hour"
	ENUM_KLINE_TYPE_12HOUR = "12hour"
	ENUM_KLINE_TYPE_1DAY   = "1day"
	ENUM_KLINE_TYPE_1WEEK  = "1week"
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
	timestamp := strconv.Itoa(int(time.Now().Unix()))
	h := hmac.New(sha256.New, []byte(e.conf.GetAPISecret()))
	h.Write([]byte(timestamp + http.MethodGet + path))
	httpReq.Header.Set(HeaderAPIKey, e.conf.GetAPIKey())
	httpReq.Header.Set(HeaderAPISign, base64.StdEncoding.EncodeToString(h.Sum(nil)))
	httpReq.Header.Set(HeaderAPITimeStamp, timestamp)

	h1 := hmac.New(sha256.New, []byte(e.conf.GetAPISecret()))
	h1.Write([]byte(e.conf.GetAPIPassPhrase()))
	httpReq.Header.Set(HeaderAPIPassPhase, base64.StdEncoding.EncodeToString(h.Sum(nil)))
	httpReq.Header.Set(HeaderAPIKeyVersion, "V2")
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
			if closePrice, err := strconv.ParseFloat(item[1], 64); err != nil {
				return nil, err
			} else {
				retItem.ClosePrice = closePrice
			}
			if highestPrice, err := strconv.ParseFloat(item[1], 64); err != nil {
				return nil, err
			} else {
				retItem.HighestPrice = highestPrice
			}
			if lowestPrice, err := strconv.ParseFloat(item[1], 64); err != nil {
				return nil, err
			} else {
				retItem.LowestPrice = lowestPrice
			}
			if transactionVolume, err := strconv.ParseFloat(item[1], 64); err != nil {
				return nil, err
			} else {
				retItem.TransactionVolume = transactionVolume
			}
			if transactionAmount, err := strconv.ParseFloat(item[1], 64); err != nil {
				return nil, err
			} else {
				retItem.TransactionAmount = transactionAmount
			}
			ret[i] = retItem
		}
		return ret, nil
	}
}
