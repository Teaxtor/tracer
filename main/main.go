package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"tracer"
)

type flags struct {
	configFile string
}

type app struct {
	api *Api
	tracer *tracer.Tracer
}

func main() {
	app := app{}
	flags := parseFlags()

	cfg, err := NewConfig(flags.configFile)
	panicIfError(err)

	app.Setup(cfg)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM)
	signal.Notify(sig, syscall.SIGINT)

	app.Start()

	<-sig

	app.Stop()
}

func parseFlags() flags {
	configFile := flag.String("config", "", "path to a config file")

	return flags{
		configFile: *configFile,
	}
}

func panicIfError(err error) {
	if err != nil {
		panic(err)
	}
}

func (a *app) Setup(config Config) {
	a.api = NewApi(config.Port)
	a.tracer = tracer.New(config.Browser, config.ProxyInfo, config.RemotePort)
}

func (a *app) Start() {
	a.api.Start()
}

func (a *app) Stop() {
	a.api.Stop()
	a.tracer.Stop()
}