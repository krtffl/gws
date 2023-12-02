package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/viper"

	"github.com/krtffl/gws/internal/api"
	"github.com/krtffl/gws/internal/config"
)

var configPath = flag.String(
	"config",
	"config/config.yaml",
	"path from where the config file will be loaded",
)

func main() {
	flag.Parse()

	cfg := config.LoadConfig(viper.New(), *configPath)

	api := api.New(cfg)

	go api.Run()
	handleSignals(api)
}

func handleSignals(api *api.GWS) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	<-signalChan
	api.Shutdown()
}
