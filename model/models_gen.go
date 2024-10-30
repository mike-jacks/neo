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
	Domain     string                 `json:"domain"`
	Name       string                 `json:"name"`
	Type       string                 `json:"type"`
	Labels     []string               `json:"labels,omitempty"`
	Properties map[string]interface{} `json:"properties,omitempty"`
}

type DomainSchemaNodeResponse struct {
	Success bool                `json:"success"`
	Message *string             `json:"message,omitempty"`
	Data    []*DomainSchemaNode `json:"data,omitempty"`
}

type Mutation struct {
}

type ObjectNode struct {
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
	Success bool          `json:"success"`
	Message *string       `json:"message,omitempty"`
	Data    []*ObjectNode `json:"data,omitempty"`
}

type ObjectRelationship struct {
	RelationshipName string      `json:"relationshipName"`
	FromObjectNode   *ObjectNode `json:"fromObjectNode"`
	ToObjectNode     *ObjectNode `json:"toObjectNode"`
	Properties       []*Property `json:"properties,omitempty"`
}

type ObjectRelationshipObjectNode struct {
	FromObjectNode   *ObjectNode             `json:"fromObjectNode"`
	RelationshipNode *RelationshipSchemaNode `json:"relationshipNode"`
	ToObjectNode     *ObjectNode             `json:"toObjectNode"`
}

type ObjectRelationshipObjectNodeResponse struct {
	Success bool                            `json:"success"`
	Message *string                         `json:"message,omitempty"`
	Data    []*ObjectRelationshipObjectNode `json:"data,omitempty"`
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
	RelationshipName       string                 `json:"relationshipName"`
	Domain                 string                 `json:"domain"`
	FromTypeSchemaNodeName string                 `json:"fromTypeSchemaNodeName"`
	ToTypeSchemaNodeName   string                 `json:"toTypeSchemaNodeName"`
	Properties             map[string]interface{} `json:"properties,omitempty"`
}

type RelationshipSchemaNodeResponse struct {
	Success bool                      `json:"success"`
	Message *string                   `json:"message,omitempty"`
	Data    []*RelationshipSchemaNode `json:"data,omitempty"`
}

type Response struct {
	Success bool                     `json:"success"`
	Message *string                  `json:"message,omitempty"`
	Data    []map[string]interface{} `json:"data,omitempty"`
}

type TypeSchemaNode struct {
	Domain       string                 `json:"domain"`
	Name         string                 `json:"name"`
	Type         string                 `json:"type"`
	OriginalName string                 `json:"originalName"`
	Labels       []string               `json:"labels,omitempty"`
	Properties   map[string]interface{} `json:"properties,omitempty"`
}

type TypeSchemaNodeResponse struct {
	Success bool              `json:"success"`
	Message *string           `json:"message,omitempty"`
	Data    []*TypeSchemaNode `json:"data,omitempty"`
}

type UpdateObjectNodeInput struct {
	Domain     *string          `json:"domain,omitempty"`
	Name       *string          `json:"name,omitempty"`
	Type       *string          `json:"type,omitempty"`
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
