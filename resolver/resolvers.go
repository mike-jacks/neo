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
func (r *mutationResolver) CreateObjectNode(ctx context.Context, domain string, name string, typeArg string, labels []string, properties []*model.PropertyInput) (*model.ObjectNodeResponse, error) {
	result, err := r.Database.CreateObjectNode(ctx, domain, name, typeArg, labels, properties)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// RenameObjectNode is the resolver for the renameObjectNode field.
func (r *mutationResolver) RenameObjectNode(ctx context.Context, id string, newName string) (*model.ObjectNodeResponse, error) {
	result, err := r.Database.RenameObjectNode(ctx, id, newName)
	if err != nil {
		return nil, err

	}
	return result, nil
}

// DeleteObjectNode is the resolver for the deleteObjectNode field.
func (r *mutationResolver) DeleteObjectNode(ctx context.Context, id string) (*model.ObjectNodeResponse, error) {
	result, err := r.Database.DeleteObjectNode(ctx, id)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// AddLabelsToObjectNode is the resolver for the addLabelsToObjectNode field.
func (r *mutationResolver) AddLabelsOnObjectNode(ctx context.Context, id string, labels []string) (*model.ObjectNodeResponse, error) {
	result, err := r.Database.AddLabelsOnObjectNode(ctx, id, labels)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// RemoveLabelsFromObjectNode is the resolver for the removeLabelsFromObjectNode field.
func (r *mutationResolver) RemoveLabelsFromObjectNode(ctx context.Context, id string, labels []string) (*model.ObjectNodeResponse, error) {
	result, err := r.Database.RemoveLabelsFromObjectNode(ctx, id, labels)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// AddPropertiesToObjectNode is the resolver for the addPropertiesToObjectNode field.
func (r *mutationResolver) UpdatePropertiesOnObjectNode(ctx context.Context, id string, properties []*model.PropertyInput) (*model.ObjectNodeResponse, error) {
	result, err := r.Database.UpdatePropertiesOnObjectNode(ctx, id, properties)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// RemovePropertiesFromObjectNode is the resolver for the removePropertiesFromObjectNode field.
func (r *mutationResolver) RemovePropertiesFromObjectNode(ctx context.Context, id string, properties []string) (*model.ObjectNodeResponse, error) {
	result, err := r.Database.RemovePropertiesFromObjectNode(ctx, id, properties)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// CreateObjectRelationship is the resolver for the createObjectRelationship field.
func (r *mutationResolver) CreateObjectRelationship(ctx context.Context, name string, properties []*model.PropertyInput, fromObjectNodeID string, toObjectNodeID string) (*model.ObjectRelationshipResponse, error) {
	result, err := r.Database.CreateObjectRelationship(ctx, name, properties, fromObjectNodeID, toObjectNodeID)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// UpdatePropertiesOnObjectRelationship is the resolver for the updatePropertiesOnObjectRelationship field.
func (r *mutationResolver) UpdatePropertiesOnObjectRelationship(ctx context.Context, id string, properties []*model.PropertyInput) (*model.ObjectRelationshipResponse, error) {
	result, err := r.Database.UpdatePropertiesOnObjectRelationship(ctx, id, properties)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// RemovePropertiesFromObjectRelationship is the resolver for the removePropertiesFromObjectRelationship field.
func (r *mutationResolver) RemovePropertiesFromObjectRelationship(ctx context.Context, id string, properties []string) (*model.ObjectRelationshipResponse, error) {
	result, err := r.Database.RemovePropertiesFromObjectRelationship(ctx, id, properties)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// DeleteObjectRelationship is the resolver for the deleteObjectRelationship field.
func (r *mutationResolver) DeleteObjectRelationship(ctx context.Context, id string) (*model.ObjectRelationshipResponse, error) {
	result, err := r.Database.DeleteObjectRelationship(ctx, id)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// CreateDomainSchemaNode is the resolver for the createDomainSchemaNode field.
func (r *mutationResolver) CreateDomainSchemaNode(ctx context.Context, domain string) (*model.DomainSchemaNodeResponse, error) {
	result, err := r.Database.CreateDomainSchemaNode(ctx, domain)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// RenameDomainSchemaNode is the resolver for the renameDomainSchemaNode field.
func (r *mutationResolver) RenameDomainSchemaNode(ctx context.Context, id string, newName string) (*model.DomainSchemaNodeResponse, error) {
	result, err := r.Database.RenameDomainSchemaNode(ctx, id, newName)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// DeleteDomainSchemaNode is the resolver for the deleteDomainSchemaNode field.
func (r *mutationResolver) DeleteDomainSchemaNode(ctx context.Context, id string) (*model.DomainSchemaNodeResponse, error) {
	result, err := r.Database.DeleteDomainSchemaNode(ctx, id)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// CreateTypeSchemaNode is the resolver for the createTypeSchemaNode field.
func (r *mutationResolver) CreateTypeSchemaNode(ctx context.Context, domain string, name string) (*model.TypeSchemaNodeResponse, error) {
	result, err := r.Database.CreateTypeSchemaNode(ctx, domain, name)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// RenameTypeSchemaNode is the resolver for the renameTypeSchemaNode field.
func (r *mutationResolver) RenameTypeSchemaNode(ctx context.Context, id string, newName string) (*model.TypeSchemaNodeResponse, error) {
	result, err := r.Database.RenameTypeSchemaNode(ctx, id, newName)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// UpdatePropertiesOnTypeSchemaNode is the resolver for the updatePropertiesOnTypeSchemaNode field.
func (r *mutationResolver) UpdatePropertiesOnTypeSchemaNode(ctx context.Context, id string, properties []*model.PropertyInput) (*model.TypeSchemaNodeResponse, error) {
	result, err := r.Database.UpdatePropertiesOnTypeSchemaNode(ctx, id, properties)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// RenamePropertyOnTypeSchemaNode is the resolver for the renamePropertyOnTypeSchemaNode field.
func (r *mutationResolver) RenamePropertyOnTypeSchemaNode(ctx context.Context, id string, oldPropertyName string, newPropertyName string) (*model.TypeSchemaNodeResponse, error) {
	result, err := r.Database.RenamePropertyOnTypeSchemaNode(ctx, id, oldPropertyName, newPropertyName)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// RemovePropertiesFromTypeSchemaNode is the resolver for the removePropertiesFromTypeSchemaNode field.
func (r *mutationResolver) RemovePropertiesFromTypeSchemaNode(ctx context.Context, id string, properties []string) (*model.TypeSchemaNodeResponse, error) {
	result, err := r.Database.RemovePropertiesFromTypeSchemaNode(ctx, id, properties)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// DeleteTypeSchemaNode is the resolver for the deleteTypeSchemaNode field.
func (r *mutationResolver) DeleteTypeSchemaNode(ctx context.Context, id string) (*model.TypeSchemaNodeResponse, error) {
	result, err := r.Database.DeleteTypeSchemaNode(ctx, id)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// CreateRelationshipSchemaNode is the resolver for the createRelationshipSchemaNode field.
func (r *mutationResolver) CreateRelationshipSchemaNode(ctx context.Context, name string, domain string, fromTypeSchemaNodeID string, toTypeSchemaNodeID string) (*model.RelationshipSchemaNodeResponse, error) {
	result, err := r.Database.CreateRelationshipSchemaNode(ctx, name, domain, fromTypeSchemaNodeID, toTypeSchemaNodeID)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// RenameRelationshipSchemaNode is the resolver for the renameRelationshipSchemaNode field.
func (r *mutationResolver) RenameRelationshipSchemaNode(ctx context.Context, id string, newName string) (*model.RelationshipSchemaNodeResponse, error) {
	result, err := r.Database.RenameRelationshipSchemaNode(ctx, id, newName)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// UpdatePropertiesOnRelationshipSchemaNode is the resolver for the updatePropertiesOnRelationshipSchemaNode field.
func (r *mutationResolver) UpdatePropertiesOnRelationshipSchemaNode(ctx context.Context, id string, properties []*model.PropertyInput) (*model.RelationshipSchemaNodeResponse, error) {
	result, err := r.Database.UpdatePropertiesOnRelationshipSchemaNode(ctx, id, properties)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// RenamePropertyOnRelationshipSchemaNode is the resolver for the renamePropertyOnRelationshipSchemaNode field.
func (r *mutationResolver) RenamePropertyOnRelationshipSchemaNode(ctx context.Context, id string, oldPropertyName string, newPropertyName string) (*model.RelationshipSchemaNodeResponse, error) {
	result, err := r.Database.RenamePropertyOnRelationshipSchemaNode(ctx, id, oldPropertyName, newPropertyName)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// RemovePropertiesFromRelationshipSchemaNode is the resolver for the removePropertiesFromRelationshipSchemaNode field.
func (r *mutationResolver) RemovePropertiesFromRelationshipSchemaNode(ctx context.Context, id string, properties []string) (*model.RelationshipSchemaNodeResponse, error) {
	result, err := r.Database.RemovePropertiesFromRelationshipSchemaNode(ctx, id, properties)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// DeleteRelationshipSchemaNode is the resolver for the deleteRelationshipSchemaNode field.
func (r *mutationResolver) DeleteRelationshipSchemaNode(ctx context.Context, id string) (*model.RelationshipSchemaNodeResponse, error) {
	result, err := r.Database.DeleteRelationshipSchemaNode(ctx, id)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// CypherMutation is the resolver for the cypherMutation field.
func (r *mutationResolver) CypherMutation(ctx context.Context, cypherStatement string) (*model.ObjectNodesOrRelationshipNodesResponse, error) {
	result, err := r.Database.CypherMutation(ctx, cypherStatement)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetObjectNode is the resolver for the getObjectNode field.
func (r *queryResolver) GetObjectNode(ctx context.Context, id string) (*model.ObjectNodeResponse, error) {
	result, err := r.Database.GetObjectNode(ctx, id)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetObjectNodes is the resolver for the getObjectNodes field.
func (r *queryResolver) GetObjectNodes(ctx context.Context, domain *string, typeArg *string) (*model.ObjectNodesResponse, error) {
	result, err := r.Database.GetObjectNodes(ctx, domain, typeArg)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetObjectNodeRelationship is the resolver for the getObjectNodeRelationship field.
func (r *queryResolver) GetObjectNodeRelationship(ctx context.Context, id string) (*model.ObjectRelationshipResponse, error) {
	result, err := r.Database.GetObjectNodeRelationship(ctx, id)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetObjectNodeOutgoingRelationships is the resolver for the getObjectNodeOutgoingRelationships field.
func (r *queryResolver) GetObjectNodeOutgoingRelationships(ctx context.Context, fromObjectNodeID string) (*model.ObjectRelationshipsResponse, error) {
	result, err := r.Database.GetObjectNodeOutgoingRelationships(ctx, fromObjectNodeID)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetObjectNodeIncomingRelationships is the resolver for the getObjectNodeIncomingRelationships field.
func (r *queryResolver) GetObjectNodeIncomingRelationships(ctx context.Context, toObjectNodeID string) (*model.ObjectRelationshipsResponse, error) {
	result, err := r.Database.GetObjectNodeIncomingRelationships(ctx, toObjectNodeID)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetDomainSchemaNode is the resolver for the getDomainSchemaNode field.
func (r *queryResolver) GetDomainSchemaNode(ctx context.Context, id string) (*model.DomainSchemaNodeResponse, error) {
	result, err := r.Database.GetDomainSchemaNode(ctx, id)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetDomainSchemaNodes is the resolver for the getDomainSchemaNodes field.
func (r *queryResolver) GetDomainSchemaNodes(ctx context.Context) (*model.DomainSchemaNodesResponse, error) {
	result, err := r.Database.GetDomainSchemaNodes(ctx)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetTypeSchemaNode is the resolver for the getTypeSchemaNode field.
func (r *queryResolver) GetTypeSchemaNode(ctx context.Context, id string) (*model.TypeSchemaNodeResponse, error) {
	result, err := r.Database.GetTypeSchemaNode(ctx, id)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetTypeSchemaNodes is the resolver for the getTypeSchemaNodes field.
func (r *queryResolver) GetTypeSchemaNodes(ctx context.Context, domain *string) (*model.TypeSchemaNodesResponse, error) {
	result, err := r.Database.GetTypeSchemaNodes(ctx, domain)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetTypeSchemaNodeOutgoingRelationships is the resolver for the getTypeSchemaNodeOutgoingRelationships field.
func (r *queryResolver) GetTypeSchemaNodeOutgoingRelationships(ctx context.Context, id string) (*model.RelationshipSchemaNodesResponse, error) {
	panic(fmt.Errorf("not implemented: GetTypeSchemaNodeOutgoingRelationships - getTypeSchemaNodeOutgoingRelationships"))
}

// GetTypeSchemaNodeIncomingRelationships is the resolver for the getTypeSchemaNodeIncomingRelationships field.
func (r *queryResolver) GetTypeSchemaNodeIncomingRelationships(ctx context.Context, id string) (*model.RelationshipSchemaNodesResponse, error) {
	panic(fmt.Errorf("not implemented: GetTypeSchemaNodeIncomingRelationships - getTypeSchemaNodeIncomingRelationships"))
}

// GetRelationshipSchemaNode is the resolver for the getRelationshipSchemaNode field.
func (r *queryResolver) GetRelationshipSchemaNode(ctx context.Context, id string) (*model.RelationshipSchemaNodeResponse, error) {
	panic(fmt.Errorf("not implemented: GetRelationshipSchemaNode - getRelationshipSchemaNode"))
}

// GetRelationshipSchemaNodes is the resolver for the getRelationshipSchemaNodes field.
func (r *queryResolver) GetRelationshipSchemaNodes(ctx context.Context, domain *string) (*model.RelationshipSchemaNodesResponse, error) {
	panic(fmt.Errorf("not implemented: GetRelationshipSchemaNodes - getRelationshipSchemaNodes"))
}

// CypherQuery is the resolver for the cypherQuery field.
func (r *queryResolver) CypherQuery(ctx context.Context, cypherStatement string) (*model.ObjectNodesOrRelationshipNodesResponse, error) {
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

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
/*
	func (r *queryResolver) GetTypeSchemaNodeRelationships(ctx context.Context, id string) (*model.RelationshipSchemaNodesResponse, error) {
	result, err := r.Database.GetTypeSchemaNodeRelationships(ctx, id)
	if err != nil {
		return nil, err
	}
	return result, nil
}
*/
