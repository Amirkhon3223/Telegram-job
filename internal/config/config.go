package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	BotToken         string
	ChannelID        int64
	DatabaseURL      string
	APIPort          string
	AdminTelegramIDs map[int64]bool
	JobMaxDays       int
}

func Load() (*Config, error) {
	channelID, err := strconv.ParseInt(os.Getenv("CHANNEL_ID"), 10, 64)
	if err != nil {
		return nil, err
	}

	adminIDs := parseAdminIDs(os.Getenv("ADMIN_TELEGRAM_IDS"))

	port := os.Getenv("API_PORT")
	if port == "" {
		port = "8080"
	}

	maxDays := 40 // default
	if days := os.Getenv("JOB_MAX_DAYS"); days != "" {
		if d, err := strconv.Atoi(days); err == nil {
			maxDays = d
		}
	}

	return &Config{
		BotToken:         os.Getenv("BOT_TOKEN"),
		ChannelID:        channelID,
		DatabaseURL:      os.Getenv("DATABASE_URL"),
		APIPort:          port,
		AdminTelegramIDs: adminIDs,
		JobMaxDays:       maxDays,
	}, nil
}

func parseAdminIDs(s string) map[int64]bool {
	result := make(map[int64]bool)
	if s == "" {
		return result
	}
	for _, idStr := range strings.Split(s, ",") {
		idStr = strings.TrimSpace(idStr)
		if id, err := strconv.ParseInt(idStr, 10, 64); err == nil {
			result[id] = true
		}
	}
	return result
}

func (c *Config) IsAdmin(telegramID int64) bool {
	return c.AdminTelegramIDs[telegramID]
}
