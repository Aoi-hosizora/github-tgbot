package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

// _configs represents the global config.Config.
var _configs *Config

func Configs() *Config {
	return _configs
}

type Config struct {
	Meta   *MetaConfig   `yaml:"meta"`
	Bot    *BotConfig    `yaml:"bot"`
	SQLite *SQLiteConfig `yaml:"sqlite"`
	Redis  *RedisConfig  `yaml:"redis"`
	Task   *TaskConfig   `yaml:"task"`
}

type MetaConfig struct {
	RunMode string `yaml:"run-mode"`
	LogName string `yaml:"log-name"`
}

type BotConfig struct {
	Token         string `yaml:"token"`
	PollerTimeout uint64 `yaml:"poller-timeout"`
}

type SQLiteConfig struct {
	Database string `yaml:"database"`
	LogMode  bool   `yaml:"log-mode"`
}

type RedisConfig struct {
	Host         string `yaml:"host"`
	Port         int32  `yaml:"port"`
	DB           int32  `yaml:"db"`
	Password     string `yaml:"password"`
	LogMode      bool   `yaml:"log-mode"`
	DialTimeout  int32  `yaml:"dial-timeout"`
	ReadTimeout  int32  `yaml:"read-timeout"`
	WriteTimeout int32  `yaml:"write-timeout"`

	MaxOpen     int32 `yaml:"max-open"`
	MaxLifetime int32 `yaml:"max-lifetime"`
	MaxIdletime int32 `yaml:"max-idletime"`
}

type TaskConfig struct {
	Activity string `yaml:"activity"`
	Issue    string `yaml:"issue"`
}

func Load(configPath string) error {
	f, err := ioutil.ReadFile(configPath)
	if err != nil {
		return err
	}

	cfg := &Config{}
	err = yaml.Unmarshal(f, &cfg)
	if err != nil {
		return err
	}

	_configs = cfg
	return nil
}
