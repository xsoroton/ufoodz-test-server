package models

// Quotes ...
type Quotes struct {
	Quotes []Quote `json:"quotes"`
}

// example https://api.myjson.com/bins/7c0qw
// {
//  "quote": {
//    "name": "Service A Rental",
//    "amount": "55.00",
//    "currency": "AUD"
//  }
//}

// Quote ...
type Quote struct {
	Quote QuoteDetails `json:"quote"`
}

// QuoteDetails ...
type QuoteDetails struct {
	Name     string `json:"name"`
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
}
