package main

import (
	"encoding/json"

	"git.bolor.net/bolorsoft/micro"
	"github.com/lonewolf101101/Architect-betting/backend/cmd/superset-saver/app"
	"github.com/lonewolf101101/Architect-betting/backend/pkg/actionlogman"
	"github.com/lonewolf101101/Architect-betting/backend/pkg/supersetman"
)

func onCustomerCreated(key micro.ConsumerInfo, event *micro.Event) error {
	var data *supersetman.CustomerCreated
	if err := json.Unmarshal(event.Data, &data); err != nil {
		app.ErrorLog.Println(err)
		return err
	}

	if data.Action == "" {
		app.InfoLog.Println("ignoring empty item")
		return nil
	}

	_, err := app.Customers.Save(&data.Customer)
	if err != nil {
		app.ErrorLog.Println("error on customer save:", err)
		return err
	}

	go AddActionLog("customer_created", data.Customer.ID, data)
	return nil
}

func onError(key micro.ConsumerInfo, event *micro.Event) error {
	var data *supersetman.Error
	if err := json.Unmarshal(event.Data, &data); err != nil {
		app.ErrorLog.Println(err)
		return err
	}

	if data.Action == "" {
		app.InfoLog.Println("ignoring empty item")
		return nil
	}

	go AddActionLog("error", 0, data)
	return nil
}

func onCustomerLogin(key micro.ConsumerInfo, event *micro.Event) error {
	var data *supersetman.CustomerLogin
	if err := json.Unmarshal(event.Data, &data); err != nil {
		app.ErrorLog.Println(err)
		return err
	}

	if data.Action == "" {
		app.InfoLog.Println("ignoring empty item")
		return nil
	}

	go AddActionLog("customer_login", data.Customer.ID, data)
	return nil
}

func onCustomerVisit(key micro.ConsumerInfo, event *micro.Event) error {
	var data *supersetman.CustomerVisit
	if err := json.Unmarshal(event.Data, &data); err != nil {
		app.ErrorLog.Println(err)
		return err
	}

	if data.Action == "" {
		app.InfoLog.Println("ignoring empty item")
		return nil
	}

	go AddActionLog("customer_visit", 0, data)
	return nil
}

//#region Customer Actionlogger

func AddActionLog(action string, refID int, data interface{}) {
	app.InfoLog.Println("Writing log:", action, refID)

	jsonData, err := json.Marshal(data)
	if err != nil {
		app.ErrorLog.Println("addActionLog:", err)
		return
	}

	actionLog := &actionlogman.ActionLog{
		Action: action,
		RefID:  refID,
		Data:   string(jsonData),
	}

	if _, err := app.Actions.Save(actionLog); err != nil {
		app.ErrorLog.Println("addActionLog:", err)
		return
	}
}
