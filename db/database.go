package db

import (
	"context"

	"github.com/mike-jacks/neo/model"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Database interface {
	GetDriver() neo4j.DriverWithContext
	CreateObjectNode(ctx context.Context, domain string, name string, typeArg string, labels []string, properties []*model.PropertyInput) (*model.Response, error)
	UpdateObjectNode(ctx context.Context, domain string, name string, typeArg string, updateObjectNodeInput model.UpdateObjectNodeInput) (*model.Response, error)
	DeleteObjectNode(ctx context.Context, domain string, name string, typeArg string) (*model.Response, error)

	AddLabelsToObjectNode(ctx context.Context, domain string, name string, typeArg string, labels []string) (*model.Response, error)
	RemoveLabelsFromObjectNode(ctx context.Context, domain string, name string, typeArg string, labels []string) (*model.Response, error)

	AddPropertiesToObjectNode(ctx context.Context, domain string, name string, typeArg string, properties []*model.PropertyInput) (*model.Response, error)
	RemovePropertiesFromObjectNode(ctx context.Context, domain string, name string, typeArg string, properties []string) (*model.Response, error)

	GetObjectNode(ctx context.Context, domain string, name string, typeArg string) (*model.Response, error)
	GetObjectNodes(ctx context.Context, domain *string, name *string, typeArg *string, labels []string) (*model.Response, error)

	CypherQuery(ctx context.Context, cypherStatement string) ([]*model.Response, error)
	CypherMutation(ctx context.Context, cypherStatement string) ([]*model.Response, error)

	CreateObjectRelationship(ctx context.Context, typeArg string, properties []*model.PropertyInput, fromObjectNode model.ObjectNodeInput, toObjectNode model.ObjectNodeInput) (*model.Response, error)
	UpdatePropertiesOnObjectRelationship(ctx context.Context, typeArg string, properties []*model.PropertyInput, fromObjectNode model.ObjectNodeInput, toObjectNode model.ObjectNodeInput) (*model.Response, error)
}
