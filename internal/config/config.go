package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/pkg/errors"
)

type CliConfig struct {
	TG                 `yaml:"tg"`
	Comagic            `yaml:"comagic"`
	CallReport         `yaml:"call_report"`
	CampaignReport     `yaml:"campaign_report"`
	CampaignConditions `yaml:"campaign_conditions"`
	BQ                 `yaml:"bq"`
	CS                 `yaml:"cloud_storage"`
}

func (cfg CliConfig) Check() error {
	if cfg.Comagic.Version == "" {
		return errors.New("can't find Comagic.Version")
	}
	if cfg.Comagic.Token == "" {
		return errors.New("can't find Comagic.Token")
	}
	if cfg.CallReport.DatasetID == "" {
		return errors.New("can't find CallReport.DatasetID")
	}
	if cfg.CallReport.TableID == "" {
		return errors.New("can't find CallReport.TableID")
	}
	if cfg.CampaignReport.DatasetID == "" {
		return errors.New("can't find CampaignReport.TableID")
	}
	if cfg.CampaignReport.TableID == "" {
		return errors.New("can't find CampaignReport.TableID")
	}
	if cfg.BQ.ServiceKeyPath == "" {
		return errors.New("can't find BQ.ServiceKeyPath")
	}
	if cfg.BQ.ProjectID == "" {
		return errors.New("can't find BQ.ProjectID")
	}
	if cfg.BQ.DatasetID == "" {
		return errors.New("can't find BQ.DatasetID")
	}
	if cfg.CS.ServiceKeyPath == "" {
		return errors.New("can't find CS.ServiceKeyPath")
	}
	if cfg.CS.BucketName == "" {
		return errors.New("can't find CS.BucketName")
	}
	return nil
}

type Comagic struct {
	Version string `yaml:"version" env:"COMAGIC_VERSION"`
	Token   string `yaml:"token"`
}

type BQ struct {
	ServiceKeyPath string `yaml:"service_key_path"`
	ProjectID      string `yaml:"project_id"`
	DatasetID      string `yaml:"dataset_id"`
}

type CS struct {
	ServiceKeyPath string `yaml:"service_key_path"`
	BucketName     string `yaml:"bucket_name"`
}

type CallReport struct {
	DatasetID string `yaml:"dataset_id"`
	TableID   string `yaml:"table_id"`
}

type CampaignReport struct {
	DatasetID string `yaml:"dataset_id"`
	TableID   string `yaml:"table_id"`
}

type CampaignConditions struct {
	DatasetID string `yaml:"dataset_id"`
	TableID   string `yaml:"table_id"`
}

type TG struct {
	IsEnabled bool   `yaml:"is_enabled" env:"TG_ENABLED"`
	Token     string `yaml:"token" env:"TG_TOKEN"`
	Chat      int64  `yaml:"chat" env:"TG_CHAT"`
}

func NewCliConfig(filePath string) (*CliConfig, error) {
	cfg := &CliConfig{}
	err := cleanenv.ReadConfig(filePath, cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

type ServerConfig struct {
	TG      `yaml:"tg"`
	GRPC    `yaml:"grpc"`
	Comagic `yaml:"comagic"`
}

type GRPC struct {
	IP   string `yaml:"ip" env:"GRPC_IP"`
	Port int    `yaml:"port" env:"GRPC_PORT"`
}

func NewServerConfig(filePath string) (*ServerConfig, error) {
	cfg := &ServerConfig{}
	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

type ScheduleConfig struct {
	Time               `yaml:"time"`
	TG                 `yaml:"tg"`
	Comagic            `yaml:"comagic"`
	CallReport         `yaml:"call_report"`
	CampaignReport     `yaml:"campaign_report"`
	CampaignConditions `yaml:"campaign_conditions"`
	BQ                 `yaml:"bq"`
	CS                 `yaml:"cloud_storage"`
}

type Time struct {
	Calls string `yaml:"calls" `
}

func NewScheduleConfig(filePath string) (*ScheduleConfig, error) {
	cfg := &ScheduleConfig{}
	err := cleanenv.ReadConfig(filePath, cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	return cfg, nil
}

func (cfg ScheduleConfig) Check() error {
	if cfg.Time.Calls == "" {
		return errors.New("can't find Time.Calls")
	}
	if cfg.Comagic.Version == "" {
		return errors.New("can't find Comagic.Version")
	}
	if cfg.Comagic.Token == "" {
		return errors.New("can't find Comagic.Token")
	}
	if cfg.CallReport.DatasetID == "" {
		return errors.New("can't find CallReport.DatasetID")
	}
	if cfg.CallReport.TableID == "" {
		return errors.New("can't find CallReport.TableID")
	}
	if cfg.CampaignReport.DatasetID == "" {
		return errors.New("can't find CampaignReport.TableID")
	}
	if cfg.CampaignReport.TableID == "" {
		return errors.New("can't find CampaignReport.TableID")
	}
	if cfg.BQ.ServiceKeyPath == "" {
		return errors.New("can't find BQ.ServiceKeyPath")
	}
	if cfg.BQ.ProjectID == "" {
		return errors.New("can't find BQ.ProjectID")
	}
	if cfg.BQ.DatasetID == "" {
		return errors.New("can't find BQ.DatasetID")
	}
	if cfg.CS.ServiceKeyPath == "" {
		return errors.New("can't find CS.ServiceKeyPath")
	}
	if cfg.CS.BucketName == "" {
		return errors.New("can't find CS.BucketName")
	}
	return nil
}
