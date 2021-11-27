package kucoingo

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

//func (b Kucoin) GeneralEndpoint() general.MarketEndpoint {
//	return general.NewGeneralEndPoint(b)
//}

func (b Kucoin) MarketDataEndpoint() MarketEndpoint {
	return NewMarketEndPoint(b)
}

//func (b Kucoin) AccountDataEndpoint() account.MarketEndpoint {
//	return account.NewAccountEndPoint(b)
//}

func (b Kucoin) GetAPIKey() string {
	return b.apiKey
}

func (b Kucoin) GetAPIPassPhrase() string {
	return b.apiKeyPassPhrase
}

func (b Kucoin) GetAPISecret() string {
	return b.apiSecretKey
}

func (b Kucoin) GetKucoinUrl() string {
	return b.kucoinUrl
}
