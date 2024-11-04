package db

import (
	"context"

	"github.com/mike-jacks/neo/model"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Database interface {
	GetDriver() neo4j.DriverWithContext

	CreateObjectNode(ctx context.Context, domain string, name string, typeArg string, labels []string, properties []*model.PropertyInput) (*model.ObjectNodeResponse, error)
	RenameObjectNode(ctx context.Context, id string, newName string) (*model.ObjectNodeResponse, error)
	DeleteObjectNode(ctx context.Context, id string) (*model.ObjectNodeResponse, error)

	AddLabelsOnObjectNode(ctx context.Context, id string, labels []string) (*model.ObjectNodeResponse, error)
	RemoveLabelsFromObjectNode(ctx context.Context, id string, labels []string) (*model.ObjectNodeResponse, error)

	UpdatePropertiesOnObjectNode(ctx context.Context, id string, properties []*model.PropertyInput) (*model.ObjectNodeResponse, error)
	RemovePropertiesFromObjectNode(ctx context.Context, id string, properties []string) (*model.ObjectNodeResponse, error)

	GetObjectNode(ctx context.Context, id string) (*model.ObjectNodeResponse, error)
	GetObjectNodes(ctx context.Context, domain *string, typeArg *string) (*model.ObjectNodesResponse, error)

	CreateObjectRelationship(ctx context.Context, relationshipName string, properties []*model.PropertyInput, fromObjectNodeId string, toObjectNodeId string) (*model.ObjectRelationshipResponse, error)
	DeleteObjectRelationship(ctx context.Context, id string) (*model.ObjectRelationshipResponse, error)

	UpdatePropertiesOnObjectRelationship(ctx context.Context, id string, properties []*model.PropertyInput) (*model.ObjectRelationshipResponse, error)
	RemovePropertiesFromObjectRelationship(ctx context.Context, id string, properties []string) (*model.ObjectRelationshipResponse, error)

	GetObjectNodeRelationship(ctx context.Context, id string) (*model.ObjectRelationshipResponse, error)
	GetObjectNodeOutgoingRelationships(ctx context.Context, fromObjectNodeId string) (*model.ObjectRelationshipsResponse, error)
	GetObjectNodeIncomingRelationships(ctx context.Context, toObjectNodeId string) (*model.ObjectRelationshipsResponse, error)

	CreateDomainSchemaNode(ctx context.Context, domain string) (*model.DomainSchemaNodeResponse, error)
	RenameDomainSchemaNode(ctx context.Context, id string, newName string) (*model.DomainSchemaNodeResponse, error)
	DeleteDomainSchemaNode(ctx context.Context, id string) (*model.DomainSchemaNodeResponse, error)
	
	GetDomainSchemaNodes(ctx context.Context) (*model.DomainSchemaNodesResponse, error)
	GetDomainSchemaNode(ctx context.Context, id string) (*model.DomainSchemaNodeResponse, error)

	CreateTypeSchemaNode(ctx context.Context, domain string, name string) (*model.TypeSchemaNodeResponse, error)
	RenameTypeSchemaNode(ctx context.Context, id string, newName string) (*model.TypeSchemaNodeResponse, error)
	UpdatePropertiesOnTypeSchemaNode(ctx context.Context, id string, properties []*model.PropertyInput) (*model.TypeSchemaNodeResponse, error)
	RenamePropertyOnTypeSchemaNode(ctx context.Context, id string, oldPropertyName string, newPropertyName string) (*model.TypeSchemaNodeResponse, error)
	RemovePropertiesFromTypeSchemaNode(ctx context.Context, id string, properties []string) (*model.TypeSchemaNodeResponse, error)
	DeleteTypeSchemaNode(ctx context.Context, id string) (*model.TypeSchemaNodeResponse, error)

	GetTypeSchemaNodes(ctx context.Context, domain *string) (*model.TypeSchemaNodesResponse, error)
	GetTypeSchemaNode(ctx context.Context, id string) (*model.TypeSchemaNodeResponse, error)

	CreateRelationshipSchemaNode(ctx context.Context, relationshipName string, domain string, fromTypeSchemaNodeName string, toTypeSchemaNodeName string) (*model.Response, error)
	UpdatePropertiesOnRelationshipSchemaNode(ctx context.Context, relationshipName string, domain string, fromTypeSchemaNodeName string, toTypeSchemaNodeName string, properties []*model.PropertyInput) (*model.Response, error)
	RenamePropertyOnRelationshipSchemaNode(ctx context.Context, relationshipName string, domain string, fromTypeSchemaNodeName string, toTypeSchemaNodeName string, oldPropertyName string, newPropertyName string) (*model.Response, error)
	RemovePropertiesFromRelationshipSchemaNode(ctx context.Context, relationshipName string, domain string, fromTypeSchemaNodeName string, toTypeSchemaNodeName string, properties []string) (*model.Response, error)

	CypherQuery(ctx context.Context, cypherStatement string) (*model.ObjectNodesOrRelationshipNodesResponse, error)
	CypherMutation(ctx context.Context, cypherStatement string) (*model.ObjectNodesOrRelationshipNodesResponse, error)
}
