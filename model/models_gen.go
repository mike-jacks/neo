// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"fmt"
	"io"
	"strconv"
)

type CreateObjectRelationshipInput struct {
	FromObjectNode *ObjectNodeInput `json:"fromObjectNode"`
	ToObjectNode   *ObjectNodeInput `json:"toObjectNode"`
}

type DeleteObjectNodeInput struct {
	Domain string `json:"domain"`
	Name   string `json:"name"`
	Type   string `json:"type"`
}

type DomainSchemaNode struct {
	ID         string                 `json:"id"`
	Domain     string                 `json:"domain"`
	Name       string                 `json:"name"`
	Type       string                 `json:"type"`
	Labels     []string               `json:"labels,omitempty"`
	Properties map[string]interface{} `json:"properties,omitempty"`
}

type DomainSchemaNodeResponse struct {
	Success          bool              `json:"success"`
	Message          *string           `json:"message,omitempty"`
	DomainSchemaNode *DomainSchemaNode `json:"domainSchemaNode"`
}

type DomainSchemaNodesResponse struct {
	Success           bool                `json:"success"`
	Message           *string             `json:"message,omitempty"`
	DomainSchemaNodes []*DomainSchemaNode `json:"domainSchemaNodes,omitempty"`
}

type Mutation struct {
}

type ObjectNode struct {
	ID           string      `json:"id"`
	Domain       string      `json:"domain"`
	Name         string      `json:"name"`
	Type         string      `json:"type"`
	OriginalName string      `json:"originalName"`
	Labels       []string    `json:"labels,omitempty"`
	Properties   []*Property `json:"properties,omitempty"`
}

type ObjectNodeInput struct {
	Domain     string           `json:"domain"`
	Name       string           `json:"name"`
	Type       string           `json:"type"`
	Labels     []string         `json:"labels,omitempty"`
	Properties []*PropertyInput `json:"properties,omitempty"`
}

type ObjectNodeResponse struct {
	Success    bool        `json:"success"`
	Message    *string     `json:"message,omitempty"`
	ObjectNode *ObjectNode `json:"objectNode"`
}

type ObjectNodesResponse struct {
	Success     bool          `json:"success"`
	Message     *string       `json:"message,omitempty"`
	ObjectNodes []*ObjectNode `json:"objectNodes,omitempty"`
}

type ObjectRelationship struct {
	ID               string      `json:"id"`
	RelationshipName string      `json:"relationshipName"`
	FromObjectNode   *ObjectNode `json:"fromObjectNode"`
	ToObjectNode     *ObjectNode `json:"toObjectNode"`
	Properties       []*Property `json:"properties,omitempty"`
}

type ObjectRelationshipObjectNode struct {
	ID               string                  `json:"id"`
	FromObjectNode   *ObjectNode             `json:"fromObjectNode"`
	RelationshipNode *RelationshipSchemaNode `json:"relationshipNode"`
	ToObjectNode     *ObjectNode             `json:"toObjectNode"`
}

type ObjectRelationshipObjectNodeResponse struct {
	Success                       bool                          `json:"success"`
	Message                       *string                       `json:"message,omitempty"`
	ObjectRelationshipObjectNodes *ObjectRelationshipObjectNode `json:"objectRelationshipObjectNodes"`
}

type ObjectRelationshipObjectNodesResponse struct {
	Success                       bool                            `json:"success"`
	Message                       *string                         `json:"message,omitempty"`
	ObjectRelationshipObjectNodes []*ObjectRelationshipObjectNode `json:"objectRelationshipObjectNodes,omitempty"`
}

type Property struct {
	Key   string       `json:"key"`
	Value any          `json:"value"`
	Type  PropertyType `json:"type"`
}

type PropertyInput struct {
	Key   string       `json:"key"`
	Value any          `json:"value"`
	Type  PropertyType `json:"type"`
}

type Query struct {
}

type RelationshipSchemaNode struct {
	ID                     string                 `json:"id"`
	RelationshipName       string                 `json:"relationshipName"`
	Domain                 string                 `json:"domain"`
	FromTypeSchemaNodeName string                 `json:"fromTypeSchemaNodeName"`
	ToTypeSchemaNodeName   string                 `json:"toTypeSchemaNodeName"`
	Properties             map[string]interface{} `json:"properties,omitempty"`
}

type RelationshipSchemaNodeResponse struct {
	Success                bool                    `json:"success"`
	Message                *string                 `json:"message,omitempty"`
	RelationshipSchemaNode *RelationshipSchemaNode `json:"relationshipSchemaNode"`
}

type RelationshipSchemaNodesResponse struct {
	Success                 bool                      `json:"success"`
	Message                 *string                   `json:"message,omitempty"`
	RelationshipSchemaNodes []*RelationshipSchemaNode `json:"relationshipSchemaNodes,omitempty"`
}

type Response struct {
	Success bool                     `json:"success"`
	Message *string                  `json:"message,omitempty"`
	Data    []map[string]interface{} `json:"data,omitempty"`
}

type TypeSchemaNode struct {
	ID           string                 `json:"id"`
	Domain       string                 `json:"domain"`
	Name         string                 `json:"name"`
	Type         string                 `json:"type"`
	OriginalName string                 `json:"originalName"`
	Labels       []string               `json:"labels,omitempty"`
	Properties   map[string]interface{} `json:"properties,omitempty"`
}

type TypeSchemaNodeResponse struct {
	Success        bool            `json:"success"`
	Message        *string         `json:"message,omitempty"`
	TypeSchemaNode *TypeSchemaNode `json:"typeSchemaNode"`
}

type TypeSchemaNodesResponse struct {
	Success         bool              `json:"success"`
	Message         *string           `json:"message,omitempty"`
	TypeSchemaNodes []*TypeSchemaNode `json:"typeSchemaNodes,omitempty"`
}

type UpdateObjectNodeInput struct {
	Name       *string          `json:"name,omitempty"`
	Labels     []string         `json:"labels,omitempty"`
	Properties []*PropertyInput `json:"properties,omitempty"`
}

type PropertyType string

const (
	PropertyTypeString       PropertyType = "STRING"
	PropertyTypeNumber       PropertyType = "NUMBER"
	PropertyTypeBoolean      PropertyType = "BOOLEAN"
	PropertyTypeArrayString  PropertyType = "ARRAY_STRING"
	PropertyTypeArrayNumber  PropertyType = "ARRAY_NUMBER"
	PropertyTypeArrayBoolean PropertyType = "ARRAY_BOOLEAN"
	PropertyTypeRelationship PropertyType = "RELATIONSHIP"
)

var AllPropertyType = []PropertyType{
	PropertyTypeString,
	PropertyTypeNumber,
	PropertyTypeBoolean,
	PropertyTypeArrayString,
	PropertyTypeArrayNumber,
	PropertyTypeArrayBoolean,
	PropertyTypeRelationship,
}

func (e PropertyType) IsValid() bool {
	switch e {
	case PropertyTypeString, PropertyTypeNumber, PropertyTypeBoolean, PropertyTypeArrayString, PropertyTypeArrayNumber, PropertyTypeArrayBoolean, PropertyTypeRelationship:
		return true
	}
	return false
}

func (e PropertyType) String() string {
	return string(e)
}

func (e *PropertyType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = PropertyType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid PropertyType", str)
	}
	return nil
}

func (e PropertyType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
