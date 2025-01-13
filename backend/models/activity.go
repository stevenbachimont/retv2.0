package models

import "time"

type Activity struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	ActivityID   string    `json:"activityId"`
	Quantity     float64   `json:"quantity"`
	Date         time.Time `json:"date"`
	CarbonAmount float64   `json:"carbonAmount"`
}

type CarbonFactor struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	Category     string  `json:"category"`
	CarbonFactor float64 `json:"carbonFactor"`
	Unit         string  `json:"unit"`
} 