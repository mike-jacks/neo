package resolver

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

// Resolver struct holds dependencies like the database driver
type Resolver struct {
	Driver neo4j.DriverWithContext
}

// NewResolver creates a new resolver with the given Neo4j driver
func NewResolver(driver neo4j.DriverWithContext) *Resolver {
	return &Resolver{
		Driver: driver,
	}
}
