package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"tracer/app"
)

type flags struct {
	configFile string
}

func main() {
	flags := parseFlags()

	cfg, err := app.NewConfig(flags.configFile)
	panicIfError(err)

	a := app.New(cfg)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM)
	signal.Notify(sig, syscall.SIGINT)

	err = a.Start()
	panicIfError(err)

	<-sig

	a.Stop()
}

func parseFlags() flags {
	configFile := flag.String("config", "", "path to a config file")
	flag.Parse()

	return flags{
		configFile: *configFile,
	}
}

func panicIfError(err error) {
	if err != nil {
		panic(err)
	}
}
