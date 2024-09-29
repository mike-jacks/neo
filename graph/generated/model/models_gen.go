// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type CreateSchemaNodeInput struct {
	Name       string                 `json:"name"`
	Domain     string                 `json:"domain"`
	Type       string                 `json:"type"`
	Labels     []*SchemaLabelInput    `json:"labels,omitempty"`
	Properties []*SchemaPropertyInput `json:"properties,omitempty"`
}

type CreateSchemaPropertyInput struct {
	Name       string `json:"name"`
	Type       string `json:"type"`
	Domain     string `json:"domain"`
	ParentName string `json:"parentName"`
}

type CreateSchemaRelationshipInput struct {
	Name       string                 `json:"name"`
	Domain     string                 `json:"domain"`
	TargetName string                 `json:"targetName"`
	Properties []*SchemaPropertyInput `json:"properties"`
	ParentName string                 `json:"parentName"`
}

type Mutation struct {
}

type Query struct {
}

type SchemaLabel struct {
	Name       string `json:"name"`
	Domain     string `json:"domain"`
	ParentName string `json:"parentName"`
}

type SchemaLabelInput struct {
	Name       string `json:"name"`
	Domain     string `json:"domain"`
	ParentName string `json:"parentName"`
}

type SchemaNode struct {
	Name       string            `json:"name"`
	Domain     string            `json:"domain"`
	Type       string            `json:"type"`
	Labels     []*SchemaLabel    `json:"labels,omitempty"`
	Properties []*SchemaProperty `json:"properties,omitempty"`
}

type SchemaProperty struct {
	Name       string `json:"name"`
	Type       string `json:"type"`
	Domain     string `json:"domain"`
	ParentName string `json:"parentName"`
}

type SchemaPropertyInput struct {
	Name       string `json:"name"`
	Type       string `json:"type"`
	Domain     string `json:"domain"`
	ParentName string `json:"parentName"`
}

type SchemaRelationship struct {
	Name       string            `json:"name"`
	Domain     string            `json:"domain"`
	TargetName string            `json:"targetName"`
	Properties []*SchemaProperty `json:"properties"`
	ParentName string            `json:"parentName"`
}

type UpdateSchemaNodeInput struct {
	Name       *string                `json:"name,omitempty"`
	Domain     *string                `json:"domain,omitempty"`
	Type       *string                `json:"type,omitempty"`
	Labels     []*SchemaLabelInput    `json:"labels,omitempty"`
	Properties []*SchemaPropertyInput `json:"properties,omitempty"`
}

type UpdateSchemaPropertyInput struct {
	Name       *string `json:"name,omitempty"`
	Type       *string `json:"type,omitempty"`
	Domain     *string `json:"domain,omitempty"`
	ParentName *string `json:"parentName,omitempty"`
}

type UpdateSchemaRelationshipInput struct {
	Name       *string                `json:"name,omitempty"`
	Domain     *string                `json:"domain,omitempty"`
	TargetName *string                `json:"targetName,omitempty"`
	Properties []*SchemaPropertyInput `json:"properties,omitempty"`
	ParentName *string                `json:"parentName,omitempty"`
}
