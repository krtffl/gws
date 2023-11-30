package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"

	getwellsoon "github.com/krtffl/get-well-soon"
	"github.com/krtffl/get-well-soon/internal/logger"
)

const (
	CommonFormat string = "common"
	JSONFormat   string = "json"
)

type Config struct {
	// Global configuration
	Port uint `mapstructure:"port" yaml:"port"`

	// Database configuration
	Database Database `mapstructure:"database" yaml:"database"`

	// Logger configuration
	Logger Logger `mapstructure:"logger" yaml:"logger"`
}

type Logger struct {
	Format   string `mapstructure:"format" yaml:"format"`
	Level    string `mapstructure:"level"  yaml:"level"`
	FilePath string `mapstructure:"path"   yaml:"path"`
}

type Database struct {
	Host     string `mapstructure:"host"     yaml:"host"`
	Port     uint   `mapstructure:"port"     yaml:"port"`
	User     string `mapstructure:"user"     yaml:"user"`
	Password string `mapstructure:"password" yaml:"password"`
	Name     string `mapstructure:"name"     yaml:"name"`
	SSLMode  string `mapstructure:"ssl"      yaml:"ssl"`
}

// LoadConfig loads a custom configuration
// or creates a new one from the default configuration
// if a custom one is not found
func LoadConfig(v *viper.Viper, path string) *Config {
	v.SetConfigFile(path)

	if err := v.ReadInConfig(); err != nil {
		log.Printf("could not load configuration file: %v. creating default one", err)
		if err := os.MkdirAll(filepath.Dir(path), 0770); err != nil {
			logger.Fatal("could not create default configuration path: %v", err)
		}

		f, err := os.Create(path)
		if err != nil {
			logger.Fatal("could not create default configuration file: %v", err)
		}
		defer f.Close()

		if _, err = f.Write(getwellsoon.DefaultConfig); err != nil {
			logger.Fatal("could not write default configuration: %v", err)
		}

		// Try to reaload the created configuration
		if err := v.ReadInConfig(); err != nil {
			logger.Fatal("could not load default configuration file: %v", err)
		}
	}

	config := &Config{}
	if err := v.Unmarshal(&config); err != nil {
		logger.Fatal("could not unmarshal configuration file: %v", err)
	}
	if config.Logger.Format != CommonFormat && config.Logger.Format != JSONFormat {
		config.Logger.Format = CommonFormat
	}

	initializeLogger(config.Logger)
	return config
}

func initializeLogger(config Logger) {
	logConfg := logger.Configuration{
		EnableConsole: true,
		ConsoleLevel:  logger.GetLevel(config.Level),
		FileLevel:     logger.GetLevel(config.Level),
		EnableFile:    len(config.FilePath) > 0,
		FileLocation:  config.FilePath,
	}

	if config.Format == JSONFormat {
		logConfg.ConsoleJSONFormat = true
		logConfg.FileJSONFormat = true
	}

	logger.NewLogger(logConfg)
}
