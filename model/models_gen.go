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

type Mutation struct {
}

type ObjectNode struct {
	Domain     string      `json:"domain"`
	Name       string      `json:"name"`
	Type       string      `json:"type"`
	Labels     []string    `json:"labels,omitempty"`
	Properties []*Property `json:"properties,omitempty"`
}

type ObjectNodeInput struct {
	Domain     string           `json:"domain"`
	Name       string           `json:"name"`
	Type       string           `json:"type"`
	Labels     []string         `json:"labels,omitempty"`
	Properties []*PropertyInput `json:"properties,omitempty"`
}

type ObjectRelationship struct {
	Type           string      `json:"type"`
	FromObjectNode *ObjectNode `json:"fromObjectNode"`
	ToObjectNode   *ObjectNode `json:"toObjectNode"`
	Properties     []*Property `json:"properties,omitempty"`
}

type Property struct {
	Key   string       `json:"key"`
	Value string       `json:"value"`
	Type  PropertyType `json:"type"`
}

type PropertyInput struct {
	Key   string       `json:"key"`
	Value string       `json:"value"`
	Type  PropertyType `json:"type"`
}

type Query struct {
}

type Response struct {
	Success bool                     `json:"success"`
	Message *string                  `json:"message,omitempty"`
	Data    []map[string]interface{} `json:"data,omitempty"`
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
	PropertyTypeInteger      PropertyType = "INTEGER"
	PropertyTypeFloat        PropertyType = "FLOAT"
	PropertyTypeBoolean      PropertyType = "BOOLEAN"
	PropertyTypeArrayString  PropertyType = "ARRAY_STRING"
	PropertyTypeArrayInteger PropertyType = "ARRAY_INTEGER"
	PropertyTypeArrayFloat   PropertyType = "ARRAY_FLOAT"
	PropertyTypeArrayBoolean PropertyType = "ARRAY_BOOLEAN"
)

var AllPropertyType = []PropertyType{
	PropertyTypeString,
	PropertyTypeInteger,
	PropertyTypeFloat,
	PropertyTypeBoolean,
	PropertyTypeArrayString,
	PropertyTypeArrayInteger,
	PropertyTypeArrayFloat,
	PropertyTypeArrayBoolean,
}

func (e PropertyType) IsValid() bool {
	switch e {
	case PropertyTypeString, PropertyTypeInteger, PropertyTypeFloat, PropertyTypeBoolean, PropertyTypeArrayString, PropertyTypeArrayInteger, PropertyTypeArrayFloat, PropertyTypeArrayBoolean:
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
