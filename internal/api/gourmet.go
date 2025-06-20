package api

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

// GourmetResponse represents the HotPepper API response
type GourmetResponse struct {
	Results struct {
		Shop []struct {
			Name         string `json:"name"`
			MobileAccess string `json:"mobile_access"`
			URLs         struct {
				PC string `json:"pc"`
			} `json:"urls"`
		} `json:"shop"`
	} `json:"results"`
}

// GetGourmet searches for restaurants
func (c *Client) GetGourmet(address, keyword string) (string, error) {
	if address == "" {
		address = "渋谷駅"
	}

	// Replace commas with spaces in keyword
	if keyword != "" {
		keyword = strings.ReplaceAll(keyword, ",", " ")
		keyword = strings.ReplaceAll(keyword, "、", " ")
	}

	requestURL := c.buildURL(c.config.HotPepperAPIHost, map[string]string{
		"key":     c.config.RecruitAPIKey,
		"address": address,
		"keyword": keyword,
		"count":   strconv.Itoa(100),
		"format":  "json",
	})

	var response GourmetResponse
	if err := c.makeGetRequest(requestURL, &response); err != nil {
		return "", fmt.Errorf("failed to get gourmet info: %w", err)
	}

	if len(response.Results.Shop) == 0 {
		return "ごめんね、お店見つけられなかったよ……", nil
	}

	// Select random shop
	shop := response.Results.Shop[rand.Intn(len(response.Results.Shop))]

	message := fmt.Sprintf("%sで探してみたよ！ こことかどうかなー！\n", address)
	message += fmt.Sprintf("%s 『%s』\n", shop.MobileAccess, shop.Name)
	message += shop.URLs.PC

	return message, nil
}
