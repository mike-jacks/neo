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
func (r *mutationResolver) AddLabelsToObjectNode(ctx context.Context, domain string, name string, typeArg string, labels []string) (*model.Response, error) {
	result, err := r.Database.AddLabelsToObjectNode(ctx, domain, name, typeArg, labels)
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
func (r *mutationResolver) AddPropertiesToObjectNode(ctx context.Context, domain string, name string, typeArg string, properties []*model.PropertyInput) (*model.Response, error) {
	result, err := r.Database.AddPropertiesToObjectNode(ctx, domain, name, typeArg, properties)
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
func (r *mutationResolver) CreateObjectRelationship(ctx context.Context, typeArg string, properties []*model.PropertyInput, fromObjectNode model.ObjectNodeInput, toObjectNode model.ObjectNodeInput) (*model.Response, error) {
	result, err := r.Database.CreateObjectRelationship(ctx, typeArg, properties, fromObjectNode, toObjectNode)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// UpdatePropertiesOnObjectRelationship is the resolver for the updatePropertiesOnObjectRelationship field.
func (r *mutationResolver) UpdatePropertiesOnObjectRelationship(ctx context.Context, typeArg string, properties []*model.PropertyInput, fromObjectNode model.ObjectNodeInput, toObjectNode model.ObjectNodeInput) (*model.Response, error) {
	result, err := r.Database.UpdatePropertiesOnObjectRelationship(ctx, typeArg, properties, fromObjectNode, toObjectNode)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// RemovePropertiesFromObjectRelationship is the resolver for the removePropertiesFromObjectRelationship field.
func (r *mutationResolver) RemovePropertiesFromObjectRelationship(ctx context.Context, typeArg string, properties []string, fromObjectNode model.ObjectNodeInput, toObjectNode model.ObjectNodeInput) (*model.Response, error) {
	panic(fmt.Errorf("not implemented: RemovePropertiesFromObjectRelationship - removePropertiesFromObjectRelationship"))
}

// DeleteObjectRelationship is the resolver for the deleteObjectRelationship field.
func (r *mutationResolver) DeleteObjectRelationship(ctx context.Context, typeArg string, fromObjectNode model.ObjectNodeInput, toObjectNode model.ObjectNodeInput) (*model.Response, error) {
	panic(fmt.Errorf("not implemented: DeleteObjectRelationship - deleteObjectRelationship"))
}

// CypherMutation is the resolver for the cypherMutation field.
func (r *mutationResolver) CypherMutation(ctx context.Context, cypherStatement string) ([]*model.Response, error) {
	panic(fmt.Errorf("not implemented: CypherMutation - cypherMutation"))
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
func (r *queryResolver) GetObjectNodeRelationship(ctx context.Context, name string, fromObjectNode model.ObjectNodeInput, toObjectNode model.ObjectNodeInput) (*model.Response, error) {
	panic(fmt.Errorf("not implemented: GetObjectNodeRelationship - getObjectNodeRelationship"))
}

// GetObjectNodeRelationships is the resolver for the getObjectNodeRelationships field.
func (r *queryResolver) GetObjectNodeRelationships(ctx context.Context, fromObjectNode model.ObjectNodeInput) (*model.Response, error) {
	panic(fmt.Errorf("not implemented: GetObjectNodeRelationships - getObjectNodeRelationships"))
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
