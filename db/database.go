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

	UpdateLabelsOnObjectNode(ctx context.Context, domain string, name string, typeArg string, labels []string) (*model.Response, error)
	RemoveLabelsFromObjectNode(ctx context.Context, domain string, name string, typeArg string, labels []string) (*model.Response, error)

	UpdatePropertiesOnObjectNode(ctx context.Context, domain string, name string, typeArg string, properties []*model.PropertyInput) (*model.Response, error)
	RemovePropertiesFromObjectNode(ctx context.Context, domain string, name string, typeArg string, properties []string) (*model.Response, error)

	GetObjectNode(ctx context.Context, domain string, name string, typeArg string) (*model.Response, error)
	GetObjectNodes(ctx context.Context, domain *string, name *string, typeArg *string, labels []string) (*model.Response, error)

	CypherQuery(ctx context.Context, cypherStatement string) ([]*model.Response, error)
	CypherMutation(ctx context.Context, cypherStatement string) ([]*model.Response, error)

	CreateObjectRelationship(ctx context.Context, relationshipName string, properties []*model.PropertyInput, fromObjectNode model.ObjectNodeInput, toObjectNode model.ObjectNodeInput) (*model.Response, error)
	UpdatePropertiesOnObjectRelationship(ctx context.Context, relationshipName string, properties []*model.PropertyInput, fromObjectNode model.ObjectNodeInput, toObjectNode model.ObjectNodeInput) (*model.Response, error)
	RemovePropertiesFromObjectRelationship(ctx context.Context, relationshipName string, properties []string, fromObjectNode model.ObjectNodeInput, toObjectNode model.ObjectNodeInput) (*model.Response, error)
	DeleteObjectRelationship(ctx context.Context, relationshipName string, fromObjectNode model.ObjectNodeInput, toObjectNode model.ObjectNodeInput) (*model.Response, error)

	GetObjectNodeRelationship(ctx context.Context, relationshipName string, fromObjectNode model.ObjectNodeInput, toObjectNode model.ObjectNodeInput) (*model.Response, error)
	GetObjectNodeOutgoingRelationships(ctx context.Context, fromObjectNode model.ObjectNodeInput) (*model.Response, error)
	GetObjectNodeIncomingRelationships(ctx context.Context, toObjectNode model.ObjectNodeInput) (*model.Response, error)

	CreateDomainSchemaNode(ctx context.Context, domain string) (*model.Response, error)
	RenameDomainSchemaNode(ctx context.Context, domain string, newName string) (*model.Response, error)
	DeleteDomainSchemaNode(ctx context.Context, domain string) (*model.Response, error)
	GetAllDomainSchemaNodes(ctx context.Context) (*model.Response, error)

	CreateTypeSchemaNode(ctx context.Context, domain string, name string) (*model.Response, error)
	RenameTypeSchemaNode(ctx context.Context, domain string, existingName string, newName string) (*model.Response, error)
	UpdatePropertiesOnTypeSchemaNode(ctx context.Context, domain string, name string, properties []*model.PropertyInput) (*model.Response, error)
	RenamePropertyOnTypeSchemaNode(ctx context.Context, domain string, name string, oldPropertyName string, newPropertyName string) (*model.Response, error)
	RemovePropertiesFromTypeSchemaNode(ctx context.Context, domain string, name string, properties []string) (*model.Response, error)
	DeleteTypeSchemaNode(ctx context.Context, domain string, name string) (*model.Response, error)
	GetAllTypeSchemaNodes(ctx context.Context, domain string) (*model.Response, error)

	CreateRelationshipSchemaNode(ctx context.Context, relationshipName string, domain string, fromTypeSchemaNodeName string, toTypeSchemaNodeName string) (*model.Response, error)
	UpdatePropertiesOnRelationshipSchemaNode(ctx context.Context, relationshipName string, domain string, fromTypeSchemaNodeName string, toTypeSchemaNodeName string, properties []*model.PropertyInput) (*model.Response, error)
	RemovePropertiesFromRelationshipSchemaNode(ctx context.Context, relationshipName string, domain string, fromTypeSchemaNodeName string, toTypeSchemaNodeName string, properties []string) (*model.Response, error)
}
