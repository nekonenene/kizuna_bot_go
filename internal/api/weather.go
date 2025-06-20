package api

import (
	"fmt"
	"strconv"
)

// WeatherResponse represents the weather API response
type WeatherResponse struct {
	Location struct {
		City string `json:"city"`
	} `json:"location"`
	Forecasts []struct {
		DateLabel string `json:"dateLabel"`
		Date      string `json:"date"`
		Telop     string `json:"telop"`
		Temperature struct {
			Max *struct {
				Celsius string `json:"celsius"`
			} `json:"max"`
		} `json:"temperature"`
	} `json:"forecasts"`
	PinpointLocations []struct {
		Name string `json:"name"`
		Link string `json:"link"`
	} `json:"pinpointLocations"`
}

// GetWeather fetches weather information for Tokyo
func (c *Client) GetWeather() (string, error) {
	url := c.buildURL(c.config.LivedoorWeatherAPIHost, map[string]string{
		"city": strconv.Itoa(c.config.TokyoCityID),
	})

	var response WeatherResponse
	if err := c.makeGetRequest(url, &response); err != nil {
		return "", fmt.Errorf("failed to get weather: %w", err)
	}

	// Find Shibuya link
	shibuyaLink := ""
	for _, location := range response.PinpointLocations {
		if location.Name == "渋谷区" {
			shibuyaLink = location.Link
			break
		}
	}

	// Build message
	message := ""
	for _, forecast := range response.Forecasts {
		maxTemp := ""
		if forecast.Temperature.Max != nil {
			maxTemp = forecast.Temperature.Max.Celsius
		}

		message += fmt.Sprintf("%s（%s）の%sの天気は「%s」",
			forecast.DateLabel, forecast.Date, response.Location.City, forecast.Telop)

		if maxTemp != "" {
			message += fmt.Sprintf("、最高気温は%s℃", maxTemp)
		}
		message += "\n"
	}

	if shibuyaLink != "" {
		message += fmt.Sprintf("詳しくはこちら → %s", shibuyaLink)
	}

	return message, nil
}
