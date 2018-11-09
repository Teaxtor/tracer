package app

import (
	"github.com/spf13/viper"
	"log"
	"strings"
	"time"
	"tracer/pkg"
)

type Config struct {
	Browser    pkg.BrowserConfig
	Port       int
	ProxyInfo  pkg.ProxyInfo
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

	proxy := pkg.ProxyInfo{
		DefaultKey: vpr.GetString("proxy_default_endpoint"),
		User:       vpr.GetString("proxy_user"),
		Password:   vpr.GetString("proxy_password"),
		Endpoints:  vpr.GetStringMapString("proxy_endpoints"),
	}

	browser := pkg.BrowserConfig{
		ScreenWidth:         vpr.GetInt("browser_screen_width"),
		ScreenHeight:        vpr.GetInt("browser_screen_height"),
		UseMobile:           vpr.GetBool("browser_use_mobile"),
		UserAgent:           vpr.GetString("browser_user_agent"),
		Timeout:             vpr.GetInt("browser_timeout"),
		RemoteConnects:      vpr.GetInt("browser_remote_connects"),
		WaitBetweenConnects: vpr.GetDuration("browser_wait_between_connects") * time.Millisecond,
		Headless: vpr.GetBool("browser_headless"),
	}

	return Config{
		Port: vpr.GetInt("api_port"),
		ProxyInfo: proxy,
		RemotePort: vpr.GetInt("remote_port"),
		Browser: browser,
	}, nil
}

func setDefaults(vpr *viper.Viper) {
	vpr.SetDefault("api_port", 80)
	vpr.SetDefault("proxy_default_endpoint", "")
	vpr.SetDefault("proxy_user", "")
	vpr.SetDefault("proxy_password", "")
	vpr.SetDefault("proxy_endpoints", map[string]string{})
	vpr.SetDefault("remote_port", 9000)
	vpr.SetDefault("browser_screen_width", 1080)
	vpr.SetDefault("browser_screen_height", 1920)
	vpr.SetDefault("browser_use_mobile", false)
	vpr.SetDefault("browser_user_agent", "")
	vpr.SetDefault("browser_timeout", 15)
	vpr.SetDefault("browser_remote_connects", 10)
	vpr.SetDefault("browser_wait_between_connects", 500)
}

func readConfig(vpr *viper.Viper, file string) error {
	if len(file) == 0 {
		return nil
	}

	log.Println("Using config file", file)

	index := strings.LastIndexAny(file, ".")

	if index != -1 {
		file = file[0:index]
	}

	vpr.SetConfigName(file) // name of config without file extension
	vpr.AddConfigPath(".")

	err := vpr.ReadInConfig()

	return err
}
