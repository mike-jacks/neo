package resolver

import (
	"github.com/mike-jacks/neo/db"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Database db.Database
}

func NewResolver(Database db.Database) *Resolver {
	return &Resolver{
		Database: Database,
	}
}
