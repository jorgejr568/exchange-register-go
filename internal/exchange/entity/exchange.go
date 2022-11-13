package entity

import "time"

type Exchange struct {
	ID             uint64  `ksql:"id"`
	BaseCurrency   string  `ksql:"base_currency"`
	TargetCurrency string  `ksql:"target_currency"`
	Rate           float64 `ksql:"rate"`

	CreatedAt time.Time  `ksql:"created_at"`
	UpdatedAt *time.Time `ksql:"updated_at"`
}

type ExchangeResponse struct {
	ID             uint64  `json:"id"`
	SourceCurrency string  `json:"source_currency"`
	TargetCurrency string  `json:"target_currency"`
	Rate           float64 `json:"rate"`
}
