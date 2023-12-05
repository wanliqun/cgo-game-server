package config

import (
	"strings"

	"github.com/mcuadros/go-defaults"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Shared instance.
var sharedInstance *Config

// Read system enviroment variables prefixed with "CGS".
// eg., `CGS_LOG_LEVEL` will override "log.level" config item from the config file.
const viperEnvPrefix = "cgs"

func Shared() *Config {
	return sharedInstance
}

type logConfig struct {
	Level      string `default:"info"`
	ForceColor bool   `default:"false"`
}

type serverConfig struct {
	ServerName            string `default:"cgo_game_server"`
	ServerPassword        string `default:"helloworld"`
	TCPEndpoint           string `default:":8765"`
	UDPEndpoint           string `default:":8765"`
	MaxPlayerCapacity     int    `default:"10000"`
	MaxConnectionCapacity int    `default:"15000"`
}

type Config struct {
	Log    logConfig
	Server serverConfig
}

func init() {
	cfg, err := initConfig()
	if err != nil {
		panic(err)
	}

	if err := initLogger(cfg); err != nil {
		panic(err)
	}

	sharedInstance = cfg
}

func initViper() error {
	// Init configuration path
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	// Find and read the config file
	err := viper.ReadInConfig()
	if err != nil {
		return errors.WithMessage(err, "failed to read config file")
	}

	// Bind environment
	viper.SetEnvPrefix(viperEnvPrefix)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	return nil
}

func initConfig() (*Config, error) {
	if err := initViper(); err != nil {
		return nil, errors.WithMessage(err, "failed to init viper")
	}

	cfg := new(Config)
	defaults.SetDefaults(cfg)

	// Unmarshal configurations
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, errors.WithMessage(err, "failed to unmarshal config")
	}

	return cfg, nil
}

func initLogger(cfg *Config) error {
	// Set log level
	level, err := logrus.ParseLevel(cfg.Log.Level)
	if err != nil {
		return errors.WithMessagef(err, "invalid log level configured: %v", cfg.Log.Level)
	}
	logrus.SetLevel(level)

	// Set force color
	if cfg.Log.ForceColor {
		logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true, FullTimestamp: true})
	}

	return nil
}
