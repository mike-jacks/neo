package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.54

import (
	"context"
	"fmt"

	"github.com/mike-jacks/neo/generated"
	"github.com/mike-jacks/neo/model"
)

// CreateObjectNode is the resolver for the createObjectNode field.
func (r *mutationResolver) CreateObjectNode(ctx context.Context, domain string, name string, typeArg string, labels []string, properties []*model.PropertyInput) (*model.Response, error) {
	result, err := r.Database.CreateObjectNode(ctx, domain, name, typeArg, labels, properties)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// UpdateObjectNode is the resolver for the updateObjectNode field.
func (r *mutationResolver) UpdateObjectNode(ctx context.Context, domain string, name string, typeArg string, updateObjectNodeInput model.UpdateObjectNodeInput) (*model.Response, error) {
	result, err := r.Database.UpdateObjectNode(ctx, domain, name, typeArg, updateObjectNodeInput)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// DeleteObjectNode is the resolver for the deleteObjectNode field.
func (r *mutationResolver) DeleteObjectNode(ctx context.Context, domain string, name string, typeArg string) (*model.Response, error) {
	result, err := r.Database.DeleteObjectNode(ctx, domain, name, typeArg)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// AddLabelsToObjectNode is the resolver for the addLabelsToObjectNode field.
func (r *mutationResolver) UpdateLabelsOnObjectNode(ctx context.Context, domain string, name string, typeArg string, labels []string) (*model.Response, error) {
	result, err := r.Database.UpdateLabelsOnObjectNode(ctx, domain, name, typeArg, labels)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// RemoveLabelsFromObjectNode is the resolver for the removeLabelsFromObjectNode field.
func (r *mutationResolver) RemoveLabelsFromObjectNode(ctx context.Context, domain string, name string, typeArg string, labels []string) (*model.Response, error) {
	result, err := r.Database.RemoveLabelsFromObjectNode(ctx, domain, name, typeArg, labels)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// AddPropertiesToObjectNode is the resolver for the addPropertiesToObjectNode field.
func (r *mutationResolver) UpdatePropertiesOnObjectNode(ctx context.Context, domain string, name string, typeArg string, properties []*model.PropertyInput) (*model.Response, error) {
	result, err := r.Database.UpdatePropertiesOnObjectNode(ctx, domain, name, typeArg, properties)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// RemovePropertiesFromObjectNode is the resolver for the removePropertiesFromObjectNode field.
func (r *mutationResolver) RemovePropertiesFromObjectNode(ctx context.Context, domain string, name string, typeArg string, properties []string) (*model.Response, error) {
	result, err := r.Database.RemovePropertiesFromObjectNode(ctx, domain, name, typeArg, properties)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// CreateObjectRelationship is the resolver for the createObjectRelationship field.
func (r *mutationResolver) CreateObjectRelationship(ctx context.Context, relationshipName string, properties []*model.PropertyInput, fromObjectNode model.ObjectNodeInput, toObjectNode model.ObjectNodeInput) (*model.Response, error) {
	result, err := r.Database.CreateObjectRelationship(ctx, relationshipName, properties, fromObjectNode, toObjectNode)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// UpdatePropertiesOnObjectRelationship is the resolver for the updatePropertiesOnObjectRelationship field.
func (r *mutationResolver) UpdatePropertiesOnObjectRelationship(ctx context.Context, relationshipName string, properties []*model.PropertyInput, fromObjectNode model.ObjectNodeInput, toObjectNode model.ObjectNodeInput) (*model.Response, error) {
	result, err := r.Database.UpdatePropertiesOnObjectRelationship(ctx, relationshipName, properties, fromObjectNode, toObjectNode)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// RemovePropertiesFromObjectRelationship is the resolver for the removePropertiesFromObjectRelationship field.
func (r *mutationResolver) RemovePropertiesFromObjectRelationship(ctx context.Context, relationshipName string, properties []string, fromObjectNode model.ObjectNodeInput, toObjectNode model.ObjectNodeInput) (*model.Response, error) {
	result, err := r.Database.RemovePropertiesFromObjectRelationship(ctx, relationshipName, properties, fromObjectNode, toObjectNode)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// DeleteObjectRelationship is the resolver for the deleteObjectRelationship field.
func (r *mutationResolver) DeleteObjectRelationship(ctx context.Context, relationshipName string, fromObjectNode model.ObjectNodeInput, toObjectNode model.ObjectNodeInput) (*model.Response, error) {
	result, err := r.Database.DeleteObjectRelationship(ctx, relationshipName, fromObjectNode, toObjectNode)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// CreateDomainSchemaNode is the resolver for the createDomainSchemaNode field.
func (r *mutationResolver) CreateDomainSchemaNode(ctx context.Context, domain string) (*model.Response, error) {
	result, err := r.Database.CreateDomainSchemaNode(ctx, domain)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// RenameDomainSchemaNode is the resolver for the renameDomainSchemaNode field.
func (r *mutationResolver) RenameDomainSchemaNode(ctx context.Context, domain string, newName string) (*model.Response, error) {
	result, err := r.Database.RenameDomainSchemaNode(ctx, domain, newName)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// DeleteDomainSchemaNode is the resolver for the deleteDomainSchemaNode field.
func (r *mutationResolver) DeleteDomainSchemaNode(ctx context.Context, domain string) (*model.Response, error) {
	result, err := r.Database.DeleteDomainSchemaNode(ctx, domain)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// NewDeleteDomainSchemaNode is the resolver for the newDeleteDomainSchemaNode field.
func (r *mutationResolver) NewDeleteDomainSchemaNode(ctx context.Context, domain string) (*model.DomainSchemaNodeResponse, error) {
	panic(fmt.Errorf("not implemented: NewDeleteDomainSchemaNode - newDeleteDomainSchemaNode"))
}

// CreateTypeSchemaNode is the resolver for the createTypeSchemaNode field.
func (r *mutationResolver) CreateTypeSchemaNode(ctx context.Context, domain string, name string) (*model.Response, error) {
	result, err := r.Database.CreateTypeSchemaNode(ctx, domain, name)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// RenameTypeSchemaNode is the resolver for the renameTypeSchemaNode field.
func (r *mutationResolver) RenameTypeSchemaNode(ctx context.Context, domain string, existingName string, newName string) (*model.Response, error) {
	result, err := r.Database.RenameTypeSchemaNode(ctx, domain, existingName, newName)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// UpdatePropertiesOnTypeSchemaNode is the resolver for the updatePropertiesOnTypeSchemaNode field.
func (r *mutationResolver) UpdatePropertiesOnTypeSchemaNode(ctx context.Context, domain string, name string, properties []*model.PropertyInput) (*model.Response, error) {
	result, err := r.Database.UpdatePropertiesOnTypeSchemaNode(ctx, domain, name, properties)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// RenamePropertyOnTypeSchemaNode is the resolver for the renamePropertyOnTypeSchemaNode field.
func (r *mutationResolver) RenamePropertyOnTypeSchemaNode(ctx context.Context, domain string, name string, oldPropertyName string, newPropertyName string) (*model.Response, error) {
	result, err := r.Database.RenamePropertyOnTypeSchemaNode(ctx, domain, name, oldPropertyName, newPropertyName)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// RemovePropertiesFromTypeSchemaNode is the resolver for the removePropertiesFromTypeSchemaNode field.
func (r *mutationResolver) RemovePropertiesFromTypeSchemaNode(ctx context.Context, domain string, name string, properties []string) (*model.Response, error) {
	result, err := r.Database.RemovePropertiesFromTypeSchemaNode(ctx, domain, name, properties)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// DeleteTypeSchemaNode is the resolver for the deleteTypeSchemaNode field.
func (r *mutationResolver) DeleteTypeSchemaNode(ctx context.Context, domain string, name string) (*model.Response, error) {
	result, err := r.Database.DeleteTypeSchemaNode(ctx, domain, name)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// CreateRelationshipSchemaNode is the resolver for the createRelationshipSchemaNode field.
func (r *mutationResolver) CreateRelationshipSchemaNode(ctx context.Context, relationshipName string, domain string, fromTypeSchemaNodeName string, toTypeSchemaNodeName string) (*model.Response, error) {
	result, err := r.Database.CreateRelationshipSchemaNode(ctx, relationshipName, domain, fromTypeSchemaNodeName, toTypeSchemaNodeName)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// UpdatePropertiesOnRelationshipSchemaNode is the resolver for the updatePropertiesOnRelationshipSchemaNode field.
func (r *mutationResolver) UpdatePropertiesOnRelationshipSchemaNode(ctx context.Context, relationshipName string, domain string, fromTypeSchemaNodeName string, toTypeSchemaNodeName string, properties []*model.PropertyInput) (*model.Response, error) {
	result, err := r.Database.UpdatePropertiesOnRelationshipSchemaNode(ctx, relationshipName, domain, fromTypeSchemaNodeName, toTypeSchemaNodeName, properties)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// RenamePropertyOnRelationshipSchemaNode is the resolver for the renamePropertyOnRelationshipSchemaNode field.
func (r *mutationResolver) RenamePropertyOnRelationshipSchemaNode(ctx context.Context, relationshipName string, domain string, fromTypeSchemaNodeName string, toTypeSchemaNodeName string, oldPropertyName string, newPropertyName string) (*model.Response, error) {
	result, err := r.Database.RenamePropertyOnRelationshipSchemaNode(ctx, relationshipName, domain, fromTypeSchemaNodeName, toTypeSchemaNodeName, oldPropertyName, newPropertyName)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// RemovePropertiesFromRelationshipSchemaNode is the resolver for the removePropertiesFromRelationshipSchemaNode field.
func (r *mutationResolver) RemovePropertiesFromRelationshipSchemaNode(ctx context.Context, relationshipName string, domain string, fromTypeSchemaNodeName string, toTypeSchemaNodeName string, properties []string) (*model.Response, error) {
	result, err := r.Database.RemovePropertiesFromRelationshipSchemaNode(ctx, relationshipName, domain, fromTypeSchemaNodeName, toTypeSchemaNodeName, properties)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// DeleteRelationshipSchemaNode is the resolver for the deleteRelationshipSchemaNode field.
func (r *mutationResolver) DeleteRelationshipSchemaNode(ctx context.Context, relationshipName string, domain string, fromTypeSchemaNodeName string, toTypeSchemaNodeName string) (*model.Response, error) {
	panic(fmt.Errorf("not implemented: DeleteRelationshipSchemaNode - deleteRelationshipSchemaNode"))
}

// CypherMutation is the resolver for the cypherMutation field.
func (r *mutationResolver) CypherMutation(ctx context.Context, cypherStatement string) ([]*model.Response, error) {
	result, err := r.Database.CypherMutation(ctx, cypherStatement)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetObjectNode is the resolver for the getObjectNode field.
func (r *queryResolver) GetObjectNode(ctx context.Context, domain string, name string, typeArg string) (*model.Response, error) {
	result, err := r.Database.GetObjectNode(ctx, domain, name, typeArg)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetObjectNodes is the resolver for the getObjectNodes field.
func (r *queryResolver) GetObjectNodes(ctx context.Context, domain *string, name *string, typeArg *string, labels []string) (*model.Response, error) {
	result, err := r.Database.GetObjectNodes(ctx, domain, name, typeArg, labels)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetObjectNodeRelationship is the resolver for the getObjectNodeRelationship field.
func (r *queryResolver) GetObjectNodeRelationship(ctx context.Context, relationshipName string, fromObjectNode model.ObjectNodeInput, toObjectNode model.ObjectNodeInput) (*model.Response, error) {
	result, err := r.Database.GetObjectNodeRelationship(ctx, relationshipName, fromObjectNode, toObjectNode)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetObjectNodeOutgoingRelationships is the resolver for the getObjectNodeOutgoingRelationships field.
func (r *queryResolver) GetObjectNodeOutgoingRelationships(ctx context.Context, fromObjectNode model.ObjectNodeInput) (*model.Response, error) {
	result, err := r.Database.GetObjectNodeOutgoingRelationships(ctx, fromObjectNode)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetObjectNodeIncomingRelationships is the resolver for the getObjectNodeIncomingRelationships field.
func (r *queryResolver) GetObjectNodeIncomingRelationships(ctx context.Context, toObjectNode model.ObjectNodeInput) (*model.Response, error) {
	result, err := r.Database.GetObjectNodeIncomingRelationships(ctx, toObjectNode)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetAllDomainSchemaNodes is the resolver for the getAllDomainSchemaNodes field.
func (r *queryResolver) GetAllDomainSchemaNodes(ctx context.Context) (*model.Response, error) {
	result, err := r.Database.GetAllDomainSchemaNodes(ctx)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetAllTypeSchemaNodes is the resolver for the getAllTypeSchemaNodes field.
func (r *queryResolver) GetAllTypeSchemaNodes(ctx context.Context, domain string) (*model.Response, error) {
	result, err := r.Database.GetAllTypeSchemaNodes(ctx, domain)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetAllRelationshipsFromTypeSchemaNode is the resolver for the getAllRelationshipsFromTypeSchemaNode field.
func (r *queryResolver) GetAllRelationshipsFromTypeSchemaNode(ctx context.Context, domain string, typeSchemaNodeName string) (*model.Response, error) {
	panic(fmt.Errorf("not implemented: GetAllRelationshipsFromTypeSchemaNode - getAllRelationshipsFromTypeSchemaNode"))
}

// NewGetAllDomainSchemaNodes is the resolver for the newGetAllDomainSchemaNodes field.
func (r *queryResolver) NewGetAllDomainSchemaNodes(ctx context.Context) (*model.DomainSchemaNodeResponse, error) {
	result, err := r.Database.NewGetAllDomainSchemaNodes(ctx)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// CypherQuery is the resolver for the cypherQuery field.
func (r *queryResolver) CypherQuery(ctx context.Context, cypherStatement string) ([]*model.Response, error) {
	result, err := r.Database.CypherQuery(ctx, cypherStatement)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
