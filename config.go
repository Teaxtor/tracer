package main

import (
	"github.com/spf13/viper"
	"strings"
)

type Config struct {
	Port string
}

func NewConfig(configFile string) (Config, error) {
	vpr := viper.GetViper()
	vpr.SetEnvPrefix("TRACER")
	vpr.AutomaticEnv()

	if err := readConfig(vpr, configFile); err != nil {
		return Config{}, err
	}

	return Config{
		Port: vpr.GetString("api_port"),
	}, nil
}

func setDefaults(vpr *viper.Viper) {
	vpr.SetDefault("api_port", 80)
}

func readConfig(vpr *viper.Viper, file string) error {
	if len(file) == 0 {
		return nil
	}

	index := strings.LastIndexAny(file, ".")

	if index != -1 {
		file = file[0:index]
	}

	vpr.SetConfigName(file) // name of config without file extension
	vpr.AddConfigPath(".")

	err := vpr.ReadInConfig()

	return err
}
