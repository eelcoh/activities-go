package config

import (
	"errors"
	"fmt"
	"strconv"
)

type Config struct {
	Port       string
	ProjectID  string
	Passphrase string
}

func Load(args []string, getenv func(string) string) (Config, error) {
	projectID := getenv("PROJECTID")
	if projectID == "" {
		return Config{}, errors.New("PROJECTID must be set")
	}

	return Config{
		Port:       resolvePort(args, getenv),
		ProjectID:  projectID,
		Passphrase: getenv("PASSPHRASE"),
	}, nil
}

func (c Config) Addr() string {
	return fmt.Sprintf(":%s", c.Port)
}

func resolvePort(args []string, getenv func(string) string) string {
	if port := getenv("PORT"); port != "" {
		return port
	}

	if len(args) > 0 {
		if _, err := strconv.Atoi(args[0]); err == nil {
			return args[0]
		}
	}

	return "8080"
}
