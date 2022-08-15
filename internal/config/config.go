package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/hedhyw/jsoncjson"
	"github.com/pkg/errors"
)

type Config struct {
	Environment string      `json:"environment"`
	Log         LogConfig   `json:"log"`
	API         API         `json:"api"`
	DB          DB          `json:"db"`
	TraceConfig TraceConfig `json:"trace_config"`
	GrpcCfg     GrpcCfg     `json:"grpc_cfg"`
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

type GrpcCfg struct {
	AddressUserSrv     string `json:"address"`
	AddressTelegramSrv string `json:"address_telegram"`
}

func defaultConfig() (cfg *Config) {
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
		},
		DB: DB{
			URL:          "postgres://rpguser:rpguser@10.10.14.10:5452/rpgDB?sslmode=disable",
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
		GrpcCfg: GrpcCfg{
			AddressTelegramSrv: "rpg-api-telegram:5012",
		},
	}
}

func LoadConfig(path string) (*Config, error) {

	cfg := defaultConfig()

	confPath := os.Getenv("CONFIG_PATH")
	if confPath != "" {
		path = confPath
	}
	f, err := os.Open(path)
	if err != nil {
		return cfg, fmt.Errorf("opening: %w", err)
	}
	defer func() {
		_ = f.Close()
	}()
	jsoncReader := jsoncjson.NewReader(f)
	if err = json.NewDecoder(jsoncReader).Decode(cfg); err != nil {
		return nil, errors.WithStack(err)
	}
	return cfg, nil
}
