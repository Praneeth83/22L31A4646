package controllers

import (
	"Backend-Test-Submission/config"
	"Backend-Test-Submission/dto"
	"Backend-Test-Submission/models"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

func GetShortUrl(c echo.Context) error {
	var req dto.ShortenRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request"})
	}
	if req.ShortCode == "" || req.OriginalLongURL == "" || req.Validity <= 0 {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Missing fields"})
	}
	var exists models.ShortURL
	if err := config.DB.Where("short_code = ?", req.ShortCode).First(&exists).Error; err == nil {
		return c.JSON(http.StatusConflict, echo.Map{"error": "Shortcode already exists"})
	}
	now := time.Now()
	expiry := now.Add(time.Duration(req.Validity) * 24 * time.Hour)
	short := models.ShortURL{
		OriginalURL: req.OriginalLongURL,
		ShortCode:   req.ShortCode,
		CreatedAt:   now,
		ExpiryDate:  expiry,
	}
	if err := config.DB.Create(&short).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "DB error"})
	}
	resp := dto.ShortenResponse{
		ShortLink:  fmt.Sprintf("http://%s/%s", c.Request().Host, req.ShortCode),
		ExpiryDate: expiry,
	}
	return c.JSON(http.StatusOK, resp)

}

func GetStats(c echo.Context) error {
	code := c.Param("shortcode")
	var short models.ShortURL
	if err := config.DB.Where("short_code = ?", code).Preload("ClickDetails").First(&short).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Shortcode not found"})
	}
	resp := dto.StatsResponse{
		OriginalURL: short.OriginalURL,
		CreatedAt:   short.CreatedAt,
		ExpiryDate:  short.ExpiryDate,
		TotalClicks: short.Clicks,
		ClickData:   short.ClickDetails,
	}
	return c.JSON(http.StatusOK, resp)
}

func UrlRedirecter(c echo.Context) error {
	code := c.Param("code")
	var short models.ShortURL
	if err := config.DB.Where("short_code = ?", code).First(&short).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Shortcode not found"})
	}
	if time.Now().After(short.ExpiryDate) {
		return c.JSON(http.StatusGone, echo.Map{"error": "Link expired"})
	}
	sourceIP := getIP(c.Request())
	location := getCountry(sourceIP)
	click := models.ClickInfo{
		ShortURLID: short.ID,
		Timestamp:  time.Now(),
		SourceIP:   sourceIP,
		Location:   location,
	}
	config.DB.Create(&click)
	config.DB.Model(&short).Update("clicks", short.Clicks+1)
	return c.Redirect(http.StatusFound, short.OriginalURL)
}

func getCountry(ip string) string {
	if ip == "127.0.0.1" || ip == "::1" {
		return "Localhost"
	}
	url := fmt.Sprintf("http://ip-api.com/json/%s?fields=country", ip)
	resp, err := http.Get(url)
	if err != nil {
		return "Unknown"
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Country string `json:"country"`
	}
	if err := json.Unmarshal(body, &data); err != nil {
		return "Unknown"
	}
	if data.Country == "" {
		return "Unknown"
	}
	return data.Country
}

func getIP(r *http.Request) string {
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		ips := strings.Split(xff, ",")
		return strings.TrimSpace(ips[0])
	}
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}
