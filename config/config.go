package config

import (
	"os"
	"strings"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/posflag"
	"github.com/knadh/koanf/v2"
	"github.com/mcuadros/go-defaults"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

var (
	sharedK *koanf.Koanf

	defaultConfigSearchPaths = []string{"config.yml", "config/config.yml"}
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
	HTTPEndpoint          string `default:":8787"`
	MaxPlayerCapacity     int    `default:"10000"`
	MaxConnectionCapacity int    `default:"15000"`
}

type CGOConfig struct {
	Enabled     bool   `default:":false"`
	ResourceDir string `default:"./resources"`
}

type Config struct {
	Log    LogConfig
	Server ServerConfig
	CGO    CGOConfig
}

func init() {
	sharedK = koanf.New(".")
}

func NewConfigFromKoanf() (*Config, error) {
	cfg := new(Config)
	defaults.SetDefaults(cfg)

	if err := sharedK.Unmarshal("", &cfg); err != nil {
		return nil, errors.WithMessage(err, "failed to unmarshal config")
	}

	return cfg, nil
}

func InitKoanf(configYaml string, flag *pflag.FlagSet) error {
	if err := InitKoanfFromFile(configYaml); err != nil {
		return err
	}

	if err := InitKoanfFromEnv(); err != nil {
		return err
	}

	return InitKoanfFromPflag(flag)
}

func InitKoanfFromFile(configYaml string) (err error) {
	if len(configYaml) > 0 {
		// Load YAML config.
		return sharedK.Load(file.Provider(configYaml), yaml.Parser())
	}

	// We also define some default search path if user doesn't provide some.
	for _, cpath := range defaultConfigSearchPaths {
		err = sharedK.Load(file.Provider(cpath), yaml.Parser())
		if err == nil {
			logrus.WithField("configYaml", cpath).Info("Used default config file")
			return nil
		}

		if !errors.Is(err, os.ErrNotExist) {
			return err
		}
	}

	return nil
}

func InitKoanfFromEnv() error {
	// Read system enviroment variables prefixed with "CGS".
	// eg., `CGS_LOG_LEVEL` will override "log.level" config item from the config file.
	return sharedK.Load(env.Provider("CGS_", ".", func(s string) string {
		return strings.Replace(strings.ToLower(
			strings.TrimPrefix(s, "CGS_")), "_", ".", -1)
	}), nil)
}

func InitKoanfFromPflag(flag *pflag.FlagSet) error {
	return sharedK.Load(posflag.Provider(flag, ".", sharedK), nil)
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
