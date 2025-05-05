package main

import (
	"net/http"
	"time"

	"github.com/lonewolf101101/Architect-betting/backend/cmd/web/app"
	"github.com/lonewolf101101/Architect-betting/backend/common/oapi"
	"github.com/lonewolf101101/Architect-betting/backend/pkg/customerman"
	"github.com/lonewolf101101/Architect-betting/backend/pkg/supersetman"
)

// #region Authorization

func Me(w http.ResponseWriter, r *http.Request) {
	customer := r.Context().Value(app.ContextKeyAuthCustomer).(*customerman.Customer)
	oapi.SendResp(w, customer)
}

func logout(w http.ResponseWriter, r *http.Request) {
	app.Session.Remove(r, "email")
	app.Session.Remove(r, "auth_user_id")
	app.Session.Remove(r, "oauth2_provider_name")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func customerVisit(w http.ResponseWriter, r *http.Request) {
	go func() {
		if err := app.KafkaServer.Push(app.Config.Topic, supersetman.CustomerVisit{Action: "customer_visit", TimeStamp: time.Now()}, app.Config.Topic, nil); err != nil {
			app.ErrorLog.Println("Failed to produce customer event:", err)
		}
	}()

	oapi.SendResp(w, nil)
}
