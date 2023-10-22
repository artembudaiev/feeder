package config

import (
	"fmt"
	"github.com/hashicorp/go-multierror"
	"os"
	"strconv"
)

type AppManager interface {
	GetConfig() AppConfig
}

type ClientManager interface {
	GetConfig() ClientConfig
}

type appManager struct {
	config AppConfig
}

func (m *appManager) GetConfig() AppConfig {
	return m.config
}

type clientManager struct {
	config ClientConfig
}

func (m *clientManager) GetConfig() ClientConfig {
	return m.config
}

func NewEnvAppManager() (AppManager, error) {
	cfg := AppConfig{}
	if err := multierror.Append(
		adjustFromEnv(AppHostEnv, &cfg.AppHost),
		adjustFromEnv(AppPortEnv, &cfg.AppPort),
		adjustFromEnv(DbUrlEnv, &cfg.DbUrl),
	).ErrorOrNil(); err != nil {
		return nil, err
	}
	return &appManager{config: cfg}, nil
}

func NewEnvClientManager() (ClientManager, error) {
	cfg := ClientConfig{}
	spamTimeout := ""
	err := multierror.Append(
		adjustFromEnv(AppHostEnv, &cfg.AppHost),
		adjustFromEnv(AppPortEnv, &cfg.AppPort),
		adjustFromEnv(SpamTimeoutMsEnv, &spamTimeout),
	).ErrorOrNil()
	if err != nil {
		return nil, err
	}
	cfg.SpamTimeoutMs, err = strconv.ParseInt(spamTimeout, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to convert %s to int %w", spamTimeout, err)
	}
	return &clientManager{config: cfg}, nil
}

const (
	AppHostEnv       = "APP_HOST"
	AppPortEnv       = "APP_PORT"
	DbUrlEnv         = "DB_URL"
	SpamTimeoutMsEnv = "SPAM_TIMEOUT_MS"
)

type AppConfig struct {
	AppHost string
	AppPort string
	DbUrl   string
}

type ClientConfig struct {
	AppHost       string
	AppPort       string
	SpamTimeoutMs int64
}

func adjustFromEnv(env string, parameter *string) error {
	value := os.Getenv(env)
	if value == "" {
		return fmt.Errorf("%s environmental value not set", env)
	}
	*parameter = value
	return nil
}
