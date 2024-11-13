package resolver

import (
	"github.com/mike-jacks/neo/db"
	"github.com/mike-jacks/neo/subscriptions"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Database      db.Database
	Subscriptions *subscriptions.SubscriptionManager
}

func NewResolver(Database db.Database) *Resolver {
	return &Resolver{
		Database:      Database,
		Subscriptions: subscriptions.NewSubscriptionManager(),
	}
}
