package config

import (
	"fmt"
	"os"
	"time"

	"github.com/hedhyw/jsoncjson"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	Environment string      `json:"environment"`
	Log         LogConfig   `json:"log"`
	API         API         `json:"api"`
	DB          DB          `json:"db"`
	TraceConfig TraceConfig `json:"trace_config"`
}

type LogConfig struct {
	Level  string `json:"level"`
	Format string `json:"format"`
}

type API struct {
	Address         string        `json:"address"`
	ReadTimeout     time.Duration `json:"read_timeout"`
	WriteTimeout    time.Duration `json:"write_timeout"`
	ShutdownTimeout time.Duration `json:"shutdown_timeout"`
	MainPath        string        `json:"main_path"`
	BotToken        string        `json:"bot_token"`
	TokenTTl        int64         `json:"token_ttl"`
	JwtKey          string        `json:"jwt_key"`
}

type DB struct {
	URL          string `json:"url"`
	SchemaName   string `json:"schema_name"`
	MaxOpenConns int    `json:"max_open_conns"`
	MaxIdleConns int    `json:"max_idle_conns"`
}

type TraceConfig struct {
	Disabled      bool    `json:"disabled"`
	ServiceName   string  `json:"service_name"`
	AgentHostPort string  `json:"agent_host_port"`
	SamplerType   string  `json:"sampler_type"`
	SamplerParam  float64 `json:"sampler_param"`
}

func DefaultConfig() (cfg *Config) {
	return &Config{
		Environment: "development",
		Log: LogConfig{
			Level:  "debug",
			Format: "json",
		},
		API: API{
			Address:         ":5001",
			ReadTimeout:     30 * time.Second,
			WriteTimeout:    30 * time.Second,
			ShutdownTimeout: 5 * time.Second,
			MainPath:        "/api/v1",
			BotToken:        "5365934075:AAEr-kvUCy__jdvgcg6HwjDeG0KHut8Mbpk",
		},
		DB: DB{
			URL:          "postgres://rpguser:rpgpass@10.10.15.90:5432/RpgDB?sslmode=disable",
			SchemaName:   "RpgDB",
			MaxOpenConns: 2,
			MaxIdleConns: 2,
		},
		TraceConfig: TraceConfig{
			Disabled:      true,
			ServiceName:   "go-aut-registration-user-telegram",
			AgentHostPort: "",
			SamplerType:   "remote",
			SamplerParam:  1,
		},
	}
}

func LoadConfig(filePath, envPrefix string, logger *logrus.Logger) (cfg *Config, err error) {
	cfg = DefaultConfig()

	viper.SetConfigType("json")
	viper.SetEnvPrefix(envPrefix)
	viper.AutomaticEnv()

	var f *os.File
	f, err = os.Open(filePath)
	if err != nil {
		return cfg, fmt.Errorf("opening: %w", err)
	}

	defer func() {
		if err := f.Close(); err != nil {
			logger.Error("err = %s", err.Error())
		}
	}()

	r := jsoncjson.NewReader(f)

	if err = viper.ReadConfig(r); err != nil {
		return cfg, fmt.Errorf("reading config: %w", err)
	}

	if err = viper.Unmarshal(cfg); err != nil {
		return cfg, fmt.Errorf("unmarshalling: %w", err)
	}

	return cfg, nil
}
