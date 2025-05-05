package common

import (
	"encoding/json"
	"net/http"

	"github.com/lonewolf101101/Architect-betting/backend/cmd/web/app"
	"github.com/lonewolf101101/Architect-betting/backend/pkg/actionlogman"
	"github.com/lonewolf101101/Architect-betting/backend/pkg/customerman"
)

func AddActionLog(r *http.Request, action string, refID int, data interface{}) {
	app.InfoLog.Println("Writing log:", action, refID)

	var customer *customerman.Customer
	if r != nil {
		customer = r.Context().Value(app.ContextKeyAuthCustomer).(*customerman.Customer)
	}

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

	if customer != nil {
		actionLog.CustomerID = customer.ID
	}

	if _, err := app.ActionLogs.Save(actionLog); err != nil {
		app.ErrorLog.Println("addActionLog:", err)
		return
	}
}
