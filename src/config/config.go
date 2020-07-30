package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

var Configs *Config

type MetaConfig struct {
	RunMode string `yaml:"run-mode"`
	LogPath string `yaml:"log-path"`
	LogName string `yaml:"log-name"`
}

type BotConfig struct {
	Token         string `yaml:"token"`
	PollerTimeout uint64 `yaml:"poller-timeout"`
}

type MysqlConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Name     string `json:"name"`
	User     string `json:"user"`
	Password string `json:"password"`
	LogMode  bool   `json:"log-mode"`
}

type RedisConfig struct {
	Host           string `yaml:"host"`
	Port           int32  `yaml:"port"`
	Db             int32  `yaml:"db"`
	Password       string `yaml:"password"`
	ConnectTimeout int32  `yaml:"connect-timeout"`
	ReadTimeout    int32  `yaml:"read-timeout"`
	WriteTimeout   int32  `yaml:"write-timeout"`
}

type TaskConfig struct {
	ActivityDuration uint64 `yaml:"activity-duration"`
	IssueDuration    uint64 `yaml:"issue-duration"`
}

type Config struct {
	Meta  *MetaConfig  `yaml:"meta"`
	Bot   *BotConfig   `yaml:"bot"`
	Mysql *MysqlConfig `yaml:"mysql"`
	Redis *RedisConfig `yaml:"redis"`
	Task  *TaskConfig  `yaml:"task"`
}

func Load(configPath string) error {
	f, err := ioutil.ReadFile(configPath)
	if err != nil {
		return err
	}

	Configs = &Config{}
	err = yaml.Unmarshal(f, Configs)
	if err != nil {
		return err
	}
	return nil
}
