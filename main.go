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

func main() {
	flags := parseFlags()

	cfg, err := NewConfig(flags.ConfigFile)
	panicIfError(err)

	api := NewApi(cfg)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM)
	signal.Notify(sig, syscall.SIGINT)

	api.Start()

	<-sig

	api.Stop()
}

func parseFlags() flags {
	configFile := flag.String("config", "", "path to a config file")

	return flags{
		ConfigFile: *configFile,
	}
}

func panicIfError(err error) {
	panic(err)
}