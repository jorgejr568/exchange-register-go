package exchangerate

type GetExchangeRateRequest struct {
	From string
	To   string
}

type GetExchangeRateResponse struct {
	Rate float64 `json:"result"`
}
