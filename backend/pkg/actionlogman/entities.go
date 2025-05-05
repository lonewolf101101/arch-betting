package actionlogman

import (
	"github.com/lonewolf101101/Architect-betting/backend/pkg/customerman"
	"github.com/lonewolf101101/Architect-betting/backend/pkg/entities"
)

const (
	ACTION_LOGGED_IN  = "logged_in"
	ACTION_LOGGED_OUT = "logged_out"

	ACTION_CUSTOMER_ADDED       = "customer_added"
	ACTION_CUSTOMER_EDITED      = "customer_edited"
	ACTION_CUSTOMER_DELETED     = "customer_deleted"
	ACTION_CUSTOMER_ACTIVATED   = "customer_activated"
	ACTION_CUSTOMER_DEACTIVATED = "customer_deactivated"

	ACTION_FEEDBACK_ADDED = "feedback_added"
)

type ActionLog struct {
	entities.Model
	Action     string                `json:"action"`
	CustomerID int                   `json:"customer_id"`
	Customer   *customerman.Customer `json:"customer"`
	RefID      int                   `json:"ref_id"`
	Data       string                `json:"data"`
}
