package main

import (
	"flag"
	"net/http"
	"time"

	"github.com/lonewolf101101/Architect-betting/backend/cmd/web/app"
	"github.com/lonewolf101101/Architect-betting/backend/pkg/actionlogman"
	"github.com/lonewolf101101/Architect-betting/backend/pkg/customerman"
)

func main() {

	mode := flag.String("mode", "debug", "Choose mode. debug, test or production")
	configPath := flag.String("conf", "../confs/web.yaml", "Configuration file path")
	flag.Parse()

	app.Init(*configPath, *mode)
	defer app.Close()

	if err := app.DB.AutoMigrate(
		&customerman.Customer{},
		&actionlogman.ActionLog{},
	); err != nil {
		app.ErrorLog.Panic(err)
	}

	// app.FrontendWS.OnConnect = socket.OnFrontendWSConnect

	defineProducer()

	srv := &http.Server{
		Addr:         app.Config.Port,
		ErrorLog:     app.ErrorLog,
		Handler:      routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Minute,
		WriteTimeout: 5 * time.Minute,
		// MaxHeaderBytes: 1 << 20,
		MaxHeaderBytes: 300 << 20,
	}

	if err := app.KafkaServer.StartNonBlocking(); err != nil {
		app.ErrorLog.Fatal(err)
	}

	app.InfoLog.Printf("Starting server on %s", app.Config.Port)
	app.ErrorLog.Fatal(srv.ListenAndServe())
}

// func main() {
// 	// The JSON string (use the provided JSON string here)
// 	str := "Artificial intelligence (AI) refers to the capability of computational systems to perform tasks typically associated with human intelligence, such as learning, reasoning, problem-solving, perception, and decision-making. It is a field of research in computer science that develops and studies methods and software that enable machines to perceive their environment and use learning and intelligence to take actions that maximize their chances of achieving defined goals.[1] Such machines may be called AIs.High-profile applications of AI include advanced web search engines (e.g., Google Search); recommendation systems (used by YouTube, Amazon, and Netflix); virtual assistants (e.g., Google Assistant, Siri, and Alexa); autonomous vehicles (e.g., Waymo); generative and creative tools (e.g., ChatGPT and AI art); and superhuman play and analysis in strategy games (e.g., chess and Go). However, many AI applications are not perceived as AI: A lot of cutting edge AI has filtered into general applications, often without being called AI because once something becomes useful enough and common enough it's not labeled AI anymore."
// 	mode := flag.String("mode", "debug", "Choose mode. debug, test or production")
// 	configPath := flag.String("conf", "../confs/web.yaml", "Configuration file path")
// 	flag.Parse()

// 	app.Init(*configPath, *mode)
// 	defer app.Close()
// 	generateQuiz(str)
// }
