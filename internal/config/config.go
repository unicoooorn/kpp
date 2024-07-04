package config

import (
	"errors"
	"os"
	"strings"
	"time"

	validatorPkg "github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type FileMonitoringType string

const WhitelistMode FileMonitoringType = "whitelist"
const BlacklistMode FileMonitoringType = "blacklist"

type Config struct {
	MonitoringPeriod time.Duration        `mapstructure:"monitoring_period" validate:"required"`
	DiskUsage        DiskUsageConfig      `mapstructure:"disk_usage"`
	FileMonitoring   FileMonitoringConfig `mapstructure:"file_monitoring"`
  Strat            ActionStrat          `mapstructure:"strat"`
}

type ActionStrat string

const StratKill ActionStrat = "kill"
const StratPause ActionStrat = "pause"
const StratStop ActionStrat = "stop"
const StratRestart ActionStrat = "restart"

type DiskUsageConfig struct {
	Max int64 `mapstructure:"max"`
}

type FileMonitoringConfig struct {
	Type  FileMonitoringType `mapstructure:"type"`
	Files []string           `mapstructure:"files"`
}

func LoadApp(configPath string) (*Config, error) {
	if _, err := os.Stat(configPath); errors.Is(err, os.ErrNotExist) {
		return nil, err
	}

	viper.SetConfigFile(configPath)

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.SetDefault("logger.instance", os.Getenv("HOSTNAME"))
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, validatorPkg.New().Struct(&config)
}
