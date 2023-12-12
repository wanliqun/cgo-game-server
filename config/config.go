package config

import (
	"strings"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/mcuadros/go-defaults"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type LogConfig struct {
	Level      string `default:"info"`
	ForceColor bool   `default:"true"`
}

type ServerConfig struct {
	Name                  string `default:"cgo_game_server"`
	Password              string `default:"helloworld"`
	TCPEndpoint           string `default:":8765"`
	UDPEndpoint           string `default:":8765"`
	MaxPlayerCapacity     int    `default:"10000"`
	MaxConnectionCapacity int    `default:"15000"`
}

type Config struct {
	Log    LogConfig
	Server ServerConfig
}

func NewConfig(configYaml string) (*Config, error) {
	k := koanf.New(".")

	// Load YAML config.
	if err := k.Load(file.Provider(configYaml), yaml.Parser()); err != nil {
		return nil, err
	}

	// Read system enviroment variables prefixed with "CGS".
	// eg., `CGS_LOG_LEVEL` will override "log.level" config item from the config file.
	k.Load(env.Provider("CGS_", ".", func(s string) string {
		return strings.Replace(strings.ToLower(
			strings.TrimPrefix(s, "CGS_")), "_", ".", -1)
	}), nil)

	cfg := new(Config)
	defaults.SetDefaults(cfg)

	if err := k.Unmarshal("", &cfg); err != nil {
		return nil, errors.WithMessage(err, "failed to unmarshal config")
	}

	return cfg, nil
}

func InitLogger(cfg *LogConfig) error {
	// Set log level
	level, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		return errors.WithMessagef(err, "invalid log level: %v", cfg.Level)
	}
	logrus.SetLevel(level)

	// Set force color
	if cfg.ForceColor {
		logrus.SetFormatter(&logrus.TextFormatter{
			ForceColors:   true,
			FullTimestamp: true,
		})
	}

	return nil
}
