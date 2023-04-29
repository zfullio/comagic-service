package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
)

type Comagic struct {
	Version string `yaml:"version" env:"COMAGIC_VERSION"`
	Token   string `yaml:"token"`
}

type TG struct {
	IsEnabled bool   `yaml:"is_enabled" env:"TG_ENABLED"`
	Token     string `yaml:"token" env:"TG_TOKEN"`
	Chat      int64  `yaml:"chat" env:"TG_CHAT"`
}

type ServerConfig struct {
	KeysDir string `yaml:"keys_dir" env:"KEYS_DIR"`
	GRPC    `yaml:"grpc"`
	Comagic `yaml:"comagic"`
	TG      `yaml:"tg"`
}

type GRPC struct {
	IP   string `yaml:"ip" env:"GRPC_IP"`
	Port int    `yaml:"port" env:"GRPC_PORT"`
}

func NewServerConfig(filePath string, useEnv bool) (*ServerConfig, error) {
	cfg := &ServerConfig{}

	if useEnv {
		err := cleanenv.ReadEnv(cfg)
		if err != nil {
			return nil, err
		}
	} else {
		err := cleanenv.ReadConfig(filePath, cfg)
		if err != nil {
			return nil, fmt.Errorf("config error: %w", err)
		}
	}

	return cfg, nil
}

type Report struct {
	ObjectName          string `yaml:"object"`
	ComagicToken        string `yaml:"comagic_token"`
	GoogleServiceKey    string `yaml:"google_service_key"`
	ProjectID           string `yaml:"project_id"`
	DatasetID           string `yaml:"dataset_id"`
	BucketName          string `yaml:"bucket_name"`
	OfflineMessageTable string `yaml:"offline_message_table"`
	CallsTable          string `yaml:"calls_table"`
}

type ScheduleConfig struct {
	Time    string `yaml:"time"`
	GRPC    `yaml:"grpc"`
	Reports []Report `yaml:"reports"`
}

func NewScheduleConfig(filePath string) (*ScheduleConfig, error) {
	cfg := &ScheduleConfig{}

	err := cleanenv.ReadConfig(filePath, cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	return cfg, nil
}
