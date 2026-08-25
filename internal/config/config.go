package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Address   string
	Database  string
	DevMode   bool
	MaxBody   int64
	SiteTitle string
}

func Load() Config {
	return LoadFrom(func(key string) string { return os.Getenv(key) })
}

func LoadFrom(get func(string) string) Config {
	address := strings.TrimSpace(get("QIN_ADDRESS"))
	if address == "" {
		address = ":8080"
	}
	database := strings.TrimSpace(get("QIN_DATABASE"))
	if database == "" {
		database = "qin-culture.db"
	}
	maxBody := int64(1 << 20)
	if raw := strings.TrimSpace(get("QIN_MAX_BODY")); raw != "" {
		if parsed, err := strconv.ParseInt(raw, 10, 64); err == nil && parsed > 0 {
			maxBody = parsed
		}
	}
	dev := strings.EqualFold(strings.TrimSpace(get("QIN_DEV")), "true")
	return Config{Address: address, Database: database, DevMode: dev, MaxBody: maxBody, SiteTitle: "听见山水：古琴艺术专题"}
}

func (c Config) Validate() error {
	if strings.TrimSpace(c.Address) == "" {
		return fmt.Errorf("address is required")
	}
	if strings.TrimSpace(c.Database) == "" {
		return fmt.Errorf("database path is required")
	}
	if c.MaxBody <= 0 {
		return fmt.Errorf("max body must be positive")
	}
	return nil
}

func (c Config) Summary() string {
	mode := "production"
	if c.DevMode {
		mode = "development"
	}
	return fmt.Sprintf("%s on %s (%s)", c.SiteTitle, c.Address, mode)
}
