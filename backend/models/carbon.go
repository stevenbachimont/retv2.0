package models

type TransportMode struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	Factor float64 `json:"factor"`
}

type CarbonCalculation struct {
	Mode     string  `json:"mode"`
	Distance float64 `json:"distance"`
	Result   float64 `json:"result"`
} 