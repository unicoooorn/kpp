package config

import (
	"errors"
	"os"
	"strings"
	"time"

	validatorPkg "github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Config struct {
	MonitoringPeriod time.Duration        `mapstructure:"monitoring_period"`
	DiskUsage        DiskUsageConfig      `mapstructure:"disk_usage"`
	FileMonitoring   FileMonitoringConfig `mapstructure:"file_monitoring"`
}

type DiskUsageConfig struct {
	Max int64 `mapstructure:"max"`
}

type FileMonitoringConfig struct {
	Mode  string   `mapstructure:"type"`
	Files []string `mapstructure:"files"`
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
