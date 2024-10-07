package db

import (
	"context"

	"github.com/mike-jacks/neo/model"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Database interface {
	CreateSchemaNode(ctx context.Context, sourceSchemaNodeName *string, createSchemaNodeInput model.CreateSchemaNodeInput) (*model.SchemaNode, error)
	UpdateSchemaNode(ctx context.Context, domain string, name string, updateSchemaNodeInput model.UpdateSchemaNodeInput) ([]*model.SchemaNode, error)
	DeleteSchemaNode(ctx context.Context, domain string, name string) (bool, error)
	GetDriver() neo4j.DriverWithContext
}
