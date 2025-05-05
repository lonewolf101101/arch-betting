package supersetman

import (
	"time"

	"github.com/lonewolf101101/Architect-betting/backend/pkg/customerman"
)

type CustomerCreated struct {
	Action   string               `json:"action"`
	Customer customerman.Customer `json:"customer"`
}

type CustomerVisit struct {
	Action    string    `json:"action"`
	TimeStamp time.Time `json:"timestamp"`
}

type Error struct {
	Action  string `json:"action"`
	Message string `json:"message"`
}

type CustomerLogin struct {
	Action   string               `json:"action"`
	Customer customerman.Customer `json:"customer"`
}
