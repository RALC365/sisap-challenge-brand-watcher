package config

import (
        "fmt"
        "os"
        "strconv"
)

type Config struct {
        DatabaseURL string
        Port        int
}

func Load() (*Config, error) {
        cfg := &Config{}

        cfg.DatabaseURL = os.Getenv("DATABASE_URL")
        if cfg.DatabaseURL == "" {
                return nil, fmt.Errorf("DATABASE_URL environment variable is required")
        }

        portStr := os.Getenv("PORT")
        if portStr == "" {
                cfg.Port = 8080
        } else {
                port, err := strconv.Atoi(portStr)
                if err != nil {
                        return nil, fmt.Errorf("PORT must be a valid integer: %w", err)
                }
                cfg.Port = port
        }

        return cfg, nil
}

func (c *Config) GetAddr() string {
        return fmt.Sprintf(":%d", c.Port)
}
