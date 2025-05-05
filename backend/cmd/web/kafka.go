package main

import "github.com/lonewolf101101/Architect-betting/backend/cmd/web/app"

func defineProducer() {
	app.KafkaServer.DefineProducers(app.Config.Topic)
}

// type PagesCreated struct {
// 	Action string `json:"action"`
// 	Pages  int    `json:"pages"`
// }
