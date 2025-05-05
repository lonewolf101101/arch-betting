package main

import (
	"flag"

	"github.com/lonewolf101101/Architect-betting/backend/cmd/superset-saver/app"
	"github.com/lonewolf101101/Architect-betting/backend/common"
	"github.com/lonewolf101101/Architect-betting/backend/pkg/actionlogman"
	"github.com/lonewolf101101/Architect-betting/backend/pkg/customerman"
)

func main() {
	mode := flag.String("mode", "debug", "Choose mode. debug, test or production")
	configPath := flag.String("conf", "../confs/superset-saver.yaml", "Configuration file path")
	flag.Parse()

	app.Init(*configPath, *mode)
	defer app.Close()
	common.CloseOnSignalInterrupt(app.Close)

	app.DB.AutoMigrate(
		new(customerman.Customer),
		new(actionlogman.ActionLog),
	)

	// add handlers to event relay
	defineProducersAndConsumers()

	app.InfoLog.Println("superset-saver starting...")
	if err := app.KafkaServer.Start(); err != nil {
		app.ErrorLog.Fatal(err)
	}
}
