package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
)

type flags struct {
	ConfigFile string
}

type app struct {
	api *Api
}

func main() {
	app := app{}
	flags := parseFlags()

	cfg, err := NewConfig(flags.ConfigFile)
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
		ConfigFile: *configFile,
	}
}

func panicIfError(err error) {
	if err != nil {
		panic(err)
	}
}

func (a *app) Setup(config Config) {
	a.api = NewApi(config)
}

func (a *app) Start() {
	a.api.Start()
}

func (a *app) Stop() {
	a.api.Stop()
}