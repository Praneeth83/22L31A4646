package dto

import (
	"Backend-Test-Submission/models"
	"time"
)

type ShortenRequest struct {
	OriginalLongURL string `json:"url"`
	Validity        int    `json:"validity"`  // in days
	ShortCode       string `json:"shortcode"` // custom code
}

type ShortenResponse struct {
	ShortLink  string    `json:"shortLink"`
	ExpiryDate time.Time `json:"expiry"`
}

type StatsResponse struct {
	OriginalURL string             `json:"original_url"`
	CreatedAt   time.Time          `json:"creation_date"`
	ExpiryDate  time.Time          `json:"expiry_date"`
	TotalClicks int                `json:"total_clicks"`
	ClickData   []models.ClickInfo `json:"detailed_click_data"`
}
