package config

import (
	"github.com/pelletier/go-toml"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Routes       map[string]ServiceInfo `toml:"routes"`
	AccessKey    map[string]ServiceInfo `toml:"accesskey"`
	Natsurl      string                 `toml:"natsUrl"`
	Port         string                 `toml:"port"`
	GateWayTopic string                 `toml:"gatewayTopic"`
}

type ServiceInfo struct {
	Address   string `toml:"address"`
	AccessKey string `toml:"accesskey"`
}

var appConf *Config = &Config{}

func InitAppConfig() {
	logrus.Info("AppConfig Loading...")
	err := loadConfig("config.toml")

	if err != nil {
		logrus.Error("Config Load Error : ", err)
	}

}

func loadConfig(filename string) error {
	config, err := toml.LoadFile(filename)

	if err != nil {
		logrus.Error("Config Load Error : ", err)
		return err
	}

	if err := config.Unmarshal(appConf); err != nil {
		logrus.Error("Config unmarshal Error : ", err)
		return err
	}

	return nil
}

func GetAppConfig() *Config {
	if appConf != nil {
		return appConf
	}

	err := loadConfig("config.toml")

	if err != nil {
		logrus.Error("Config Load Error : ", err)
	}

	return appConf

}
