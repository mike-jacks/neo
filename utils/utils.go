package utils

import (
	"context"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func DereferenceOrNilString(s *string) interface{} {
	if s == nil {
		return nil
	}
	return *s
}

func CreateConstraint(ctx context.Context, driver neo4j.DriverWithContext, queries []*string) error {
	session := driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	for _, query := range queries {
		_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
			_, err := tx.Run(ctx, *query, nil)
			return nil, err
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func StringPtr(s string) *string {
	return &s
}

func PopString(m map[string]interface{}, key string) string {
	value, ok := m[key]
	if !ok {
		return ""
	}
	delete(m, key)
	return value.(string)
}
