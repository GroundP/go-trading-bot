package model

type PriceTicker struct {
    Symbol string `json:"symbol"`
    Price  string `json:"price"`
}

type Price struct {
    Asset string
    Price  string
}