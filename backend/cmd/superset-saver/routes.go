package main

import (
	"git.bolor.net/bolorsoft/micro"
	"github.com/lonewolf101101/Architect-betting/backend/cmd/superset-saver/app"
	"github.com/lonewolf101101/Architect-betting/backend/pkg/supersetman"
)

func defineProducersAndConsumers() {
	app.KafkaServer.DefineProducers(app.Config.Topic)

	app.KafkaServer.DefineConsumer(app.Config.Topic, app.Config.GroupID, 1, func(cinfo micro.ConsumerInfo) {
		cinfo.Handle(supersetman.CustomerCreated{}, onCustomerCreated)
		cinfo.Handle(supersetman.CustomerLogin{}, onCustomerLogin)
		cinfo.Handle(supersetman.CustomerVisit{}, onCustomerVisit)
		cinfo.Handle(supersetman.Error{}, onError)
	})
}
