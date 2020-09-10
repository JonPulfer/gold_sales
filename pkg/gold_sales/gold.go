package gold_sales

import "time"

// GoldPayment details for a Gold spend.
type GoldPayment struct {
	Spender      Spender   `json:"spender"`
	Description  string    `json:"description"`
	Amount       float64   `json:"amount"`
	Rate         float64   `json:"rate"`
	ToCurrency   string    `json:"toCurrency"`
	FromCurrency string    `json:"fromCurrency"`
	Date         time.Time `json:"date"`
	GramWeight   float64   `json:"gramWeight"`
}

const GoldSpend = "CARD SPEND"

const GoldCurrencyCode = "GGM"
