package fxratesapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Rates struct {
	Rub float64 `json:"RUB"`
	Usd float64 `json:"USD"`
	Eur float64 `json:"EUR"`
	Try float64 `json:"TRY"`
	Sar float64 `json:"SAR"`
}

type LatestResponse struct {
	Rates Rates `json:"rates"`
}

func New(token string) *FxRatesApiClient {
	return &FxRatesApiClient{
		token: token,
	}
}

type FxRatesApiClient struct {
	token string
}

func (c *FxRatesApiClient) GetCurrencyRate(ctx context.Context, base string) (map[string]float64, error) {
	resp, err := http.Get(fmt.Sprintf("https://api.fxratesapi.com/latest?api_key=%s&base=%s", c.token, base))
	if err != nil {
		return nil, fmt.Errorf("error get latest: %w", err)
	}
	var ratesResp LatestResponse
	err = json.NewDecoder(resp.Body).Decode(&ratesResp)
	if err != nil {
		return nil, fmt.Errorf("error unmarhall latest: %w", err)
	}
	return map[string]float64{
		"USD": ratesResp.Rates.Usd,
		"RUB": ratesResp.Rates.Rub,
		"EUR": ratesResp.Rates.Eur,
		"TRY": ratesResp.Rates.Try,
		"SAR": ratesResp.Rates.Sar,
	}, nil
}
