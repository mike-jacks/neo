// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type CreateSchemaLabelInput struct {
	Name                 string `json:"name"`
	Domain               string `json:"domain"`
	ParentSchemaNodeName string `json:"parentSchemaNodeName"`
}

type CreateSchemaNodeInput struct {
	Name   string `json:"name"`
	Domain string `json:"domain"`
}

type CreateSchemaPropertyInput struct {
	Name                 string `json:"name"`
	Type                 string `json:"type"`
	Domain               string `json:"domain"`
	ParentSchemaNodeName string `json:"parentSchemaNodeName"`
}

type CreateSchemaRelationshipInput struct {
	Name                 string `json:"name"`
	Domain               string `json:"domain"`
	TargetSchemaNodeName string `json:"targetSchemaNodeName"`
	ParentSchemaNodeName string `json:"parentSchemaNodeName"`
}

type Mutation struct {
}

type Query struct {
}

type SchemaLabel struct {
	Name                 string `json:"name"`
	Domain               string `json:"domain"`
	ParentSchemaNodeName string `json:"parentSchemaNodeName"`
}

type SchemaNode struct {
	Name          string                `json:"name"`
	Domain        string                `json:"domain"`
	ParentName    *string               `json:"parentName,omitempty"`
	Properties    []*SchemaProperty     `json:"properties,omitempty"`
	Relationships []*SchemaRelationship `json:"relationships,omitempty"`
	Labels        []*SchemaLabel        `json:"labels,omitempty"`
}

type SchemaProperty struct {
	Name                 string `json:"name"`
	Type                 string `json:"type"`
	Domain               string `json:"domain"`
	ParentSchemaNodeName string `json:"parentSchemaNodeName"`
}

type SchemaRelationship struct {
	Name                 string `json:"name"`
	Domain               string `json:"domain"`
	TargetSchemaNodeName string `json:"targetSchemaNodeName"`
	ParentSchemaNodeName string `json:"parentSchemaNodeName"`
}

type UpdateSchemaLabelInput struct {
	Name                 *string `json:"name,omitempty"`
	Domain               *string `json:"domain,omitempty"`
	ParentSchemaNodeName *string `json:"parentSchemaNodeName,omitempty"`
}

type UpdateSchemaNodeInput struct {
	Name   *string `json:"name,omitempty"`
	Domain *string `json:"domain,omitempty"`
}

type UpdateSchemaPropertyInput struct {
	Name                 *string `json:"name,omitempty"`
	Type                 *string `json:"type,omitempty"`
	Domain               *string `json:"domain,omitempty"`
	ParentSchemaNodeName *string `json:"parentSchemaNodeName,omitempty"`
}

type UpdateSchemaRelationshipInput struct {
	Name                 *string `json:"name,omitempty"`
	Domain               *string `json:"domain,omitempty"`
	TargetSchemaNodeName *string `json:"targetSchemaNodeName,omitempty"`
	ParentSchemaNodeName *string `json:"parentSchemaNodeName,omitempty"`
}
