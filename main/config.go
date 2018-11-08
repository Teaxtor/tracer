package main

import (
	"github.com/spf13/viper"
	"strings"
	"time"
	"tracer"
)

type Config struct {
	Browser tracer.BrowserConfig
	Port int
	ProxyInfo tracer.ProxyInfo
	RemotePort int
}

func NewConfig(configFile string) (Config, error) {
	vpr := viper.GetViper()
	vpr.SetEnvPrefix("TRACER")
	vpr.AutomaticEnv()

	setDefaults(vpr)

	if err := readConfig(vpr, configFile); err != nil {
		return Config{}, err
	}

	proxy := tracer.ProxyInfo{
		DefaultKey: vpr.GetString("proxy_default_endpoint"),
		User:       vpr.GetString("proxy_user"),
		Password:   vpr.GetString("proxy_password"),
		Endpoints:  vpr.GetStringMapString("proxy_endpoints"),
	}

	return Config{
		Port: vpr.GetInt("api_port"),
		ProxyInfo: proxy,
		RemotePort: vpr.GetInt("remote_port"),
		Browser: tracer.BrowserConfig{
			UserAgent:           vpr.GetString("browser_user_agent"),
			Timeout:             vpr.GetInt("browser_timeout"),
			RemoteConnects:      vpr.GetInt("browser_remote_connects"),
			WaitBetweenConnects: vpr.GetDuration("browser_wait_between_connects"),
		},
	}, nil
}

func setDefaults(vpr *viper.Viper) {
	vpr.SetDefault("api_port", 80)
	vpr.SetDefault("proxy_default_endpoint", "")
	vpr.SetDefault("proxy_user", "")
	vpr.SetDefault("proxy_password", "")
	vpr.SetDefault("proxy_endpoints", map[string]string{})
	vpr.SetDefault("remote_port", 9000)
	vpr.SetDefault("browser_user_agent", "")
	vpr.SetDefault("browser_timeout", 15)
	vpr.SetDefault("browser_remote_connects", 10)
	vpr.SetDefault("browser_wait_between_connects", 500 * time.Millisecond)
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
