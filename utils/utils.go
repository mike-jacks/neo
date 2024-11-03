package utils

import (
	"context"
	"fmt"
	"math/rand/v2"
	"regexp"
	"strings"
	"sync/atomic"
	"time"

	"github.com/mike-jacks/neo/model"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/nrednav/cuid2"
)

var specialProps = map[string]bool{
	"_originalrelationshipname": true,
	"_relationshipname":         true,
	"_domain":                   true,
	"_name":                     true,
	"_type":                     true,
	"_originalname":             true,
	"_fromtypeschemanodename":   true,
	"_totypeschemanodename":     true,
	"_id":                       true,
	"_fromobjectnodename":       true,
	"_toobjectnodename":         true,
	"_fromObjectNodeId":         true,
	"_toobjectnodeid":           true,
	"_properties":               true,
}

func DereferenceOrNilString(s *string) interface{} {
	if s == nil {
		return nil
	}
	return *s
}

func CreateConstraint(ctx context.Context, driver neo4j.DriverWithContext, queries []*string) error {
	session := driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	for _, query := range queries {
		_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
			_, err := tx.Run(ctx, *query, nil)
			return nil, err
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func PopString(m map[string]interface{}, key string) string {
	value, ok := m[key]
	if !ok {
		return ""
	}
	delete(m, key)
	return value.(string)
}

func CreatePropertiesQuery(query string, properties []*model.PropertyInput, prefix ...string) string {
	for _, property := range properties {
		if specialProps[property.Key] {
			continue
		}
		if len(prefix) > 0 {
			query += fmt.Sprintf("%v.", prefix[0])
			if property.Type.String() == "STRING" {
				query += fmt.Sprintf("%v = \"%v\", ", property.Key, property.Value)
			} else if property.Type.String() == "BOOLEAN" {
				query += fmt.Sprintf("%v = %v, ", property.Key, property.Value)
			} else if property.Type.String() == "NUMBER" {
				query += fmt.Sprintf("%v = %v, ", property.Key, property.Value)
			} else if property.Type.String() == "ARRAY_STRING" {
				query += fmt.Sprintf("%v = [", property.Key)
				for _, value := range property.Value.([]interface{}) {
					query += fmt.Sprintf("\"%v\", ", value)
				}
				query = strings.TrimSuffix(query, ", ")
				query += "], "
			} else if property.Type.String() == "ARRAY_NUMBER" {
				query += fmt.Sprintf("%v = [", property.Key)
				for _, value := range property.Value.([]interface{}) {
					query += fmt.Sprintf("%v, ", value)
				}
				query = strings.TrimSuffix(query, ", ")
				query += "], "
			} else if property.Type.String() == "ARRAY_BOOLEAN" {
				query += fmt.Sprintf("%v = [", property.Key)
				for _, value := range property.Value.([]interface{}) {
					query += fmt.Sprintf("%v, ", value)
				}
				query = strings.TrimSuffix(query, ", ")
				query += "], "
			}
		} else {
			if property.Type.String() == "STRING" {
				query += fmt.Sprintf("%v: \"%v\", ", property.Key, property.Value)
			} else if property.Type.String() == "BOOLEAN" {
				query += fmt.Sprintf("%v: %v, ", property.Key, property.Value)
			} else if property.Type.String() == "NUMBER" {
				query += fmt.Sprintf("%v: %v, ", property.Key, property.Value)
			} else if property.Type.String() == "ARRAY_STRING" {
				query += fmt.Sprintf("%v: [", property.Key)
				for _, value := range property.Value.([]interface{}) {
					query += fmt.Sprintf("\"%v\", ", value)
				}
				query = strings.TrimSuffix(query, ", ")
				query += "], "
			} else if property.Type.String() == "ARRAY_NUMBER" {
				query += fmt.Sprintf("%v: [", property.Key)
				for _, value := range property.Value.([]interface{}) {
					query += fmt.Sprintf("%v, ", value)
				}
				query = strings.TrimSuffix(query, ", ")
				query += "], "
			} else if property.Type.String() == "ARRAY_BOOLEAN" {
				query += fmt.Sprintf("%v: [", property.Key)
				for _, value := range property.Value.([]interface{}) {
					query += fmt.Sprintf("%v, ", value)
				}
				query = strings.TrimSuffix(query, ", ")
				query += "], "
			}
		}

	}
	return query
}

func RemovePropertiesQuery(query string, properties []string, prefix ...string) string {
	if len(prefix) > 0 {
		for _, property := range properties {
			if specialProps[property] {
				continue
			}
			query += fmt.Sprintf("%v.%v = null, ", prefix[0], property)
		}
	} else {
		for _, property := range properties {
			if specialProps[property] {
				continue
			}
			query += fmt.Sprintf("%v: null, ", property)
		}
	}
	return query
}

func CleanUpRelationshipName(relationshipName string) string {
	return strings.ReplaceAll(strings.TrimSpace(strings.ToUpper(relationshipName)), " ", "_")
}

func CleanUpObjectNode(objectNode *model.ObjectNodeInput) error {
	if objectNode == nil {
		return fmt.Errorf("objectNode is required")
	}
	objectNode.Domain = strings.TrimSpace(objectNode.Domain)
	objectNode.Name = strings.TrimSpace(strings.ToUpper(objectNode.Name))
	objectNode.Type = strings.TrimSpace(strings.ToUpper(objectNode.Type))
	return nil
}

func CleanUpPropertyObjects(properties *[]*model.PropertyInput) error {
	result := []*model.PropertyInput{}
	for _, property := range *properties {
		cleanPropKey := RemoveSpacesAndLowerCase(property.Key)
		if !specialProps[cleanPropKey] {
			property.Key = cleanPropKey
			result = append(result, property)
		}
	}

	*properties = result

	if len(*properties) == 0 {
		return fmt.Errorf("properties are required")
	}

	return nil
}

func CleanUpPropertyKeys(properties *[]string) error {
	result := []string{}
	for _, property := range *properties {
		cleanProp := RemoveSpacesAndLowerCase(property)
		if !specialProps[cleanProp] {
			result = append(result, cleanProp)
		}
	}

	*properties = result

	if len(*properties) == 0 {
		return fmt.Errorf("properties are required")
	}

	return nil
}

func RenamePropertyQuery(query string, oldPropertyName string, newPropertyName string, prefix ...string) string {
	if len(prefix) > 0 {
		query += fmt.Sprintf("%v.%v = %v.%v, ", prefix[0], newPropertyName, prefix[0], oldPropertyName)
		query += fmt.Sprintf("%v.%v = null, ", prefix[0], oldPropertyName)
	} else {
		query += fmt.Sprintf("%v: $newPropertyName, ", newPropertyName)
		query += fmt.Sprintf("%v: null, ", oldPropertyName)
	}
	return query
}

func SanitizeStringToLower(s string) string {
	reg := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	return reg.ReplaceAllString(strings.ToLower(strings.TrimSpace(s)), "_")
}

func SanitizeStringToUpper(s string) string {
	reg := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	return reg.ReplaceAllString(strings.ToUpper(strings.TrimSpace(s)), "_")
}

// Generate unique cuid2 ids

var generator func() string

func init() {
	var err error
	startValue := time.Now().UnixNano()
	counter := NewCounter(startValue)

	generator, err = cuid2.Init(
		cuid2.WithRandomFunc(rand.Float64),
		cuid2.WithLength(32),
		cuid2.WithFingerprint("All Your Base Are Belong To Us"),
		cuid2.WithSessionCounter(counter),
	)
	if err != nil {
		panic(err)
	}
}

type Counter struct {
	value int64
}

func NewCounter(initialCount int64) *Counter {
	return &Counter{value: initialCount}
}

func (c *Counter) Increment() int64 {
	return atomic.AddInt64(&c.value, 1)
}

func GenerateId() string {
	return generator()
}

func ExtractPropertiesFromNeo4jNode(properties map[string]interface{}) []*model.Property {
	extractedProperties := []*model.Property{}
	for key, value := range properties {
		switch value.(type) {
		case string:
			extractedProperties = append(extractedProperties, &model.Property{Key: key, Value: value, Type: model.PropertyTypeString})
		case int64:
			extractedProperties = append(extractedProperties, &model.Property{Key: key, Value: value, Type: model.PropertyTypeNumber})
		case bool:
			extractedProperties = append(extractedProperties, &model.Property{Key: key, Value: value, Type: model.PropertyTypeBoolean})
		case float64:
			extractedProperties = append(extractedProperties, &model.Property{Key: key, Value: value, Type: model.PropertyTypeNumber})
		case []interface{}:
			switch value.([]interface{})[0].(type) {
			case string:
				extractedProperties = append(extractedProperties, &model.Property{Key: key, Value: value, Type: model.PropertyTypeArrayString})
			case int64:
				extractedProperties = append(extractedProperties, &model.Property{Key: key, Value: value, Type: model.PropertyTypeArrayNumber})
			case bool:
				extractedProperties = append(extractedProperties, &model.Property{Key: key, Value: value, Type: model.PropertyTypeArrayBoolean})
			case float64:
				extractedProperties = append(extractedProperties, &model.Property{Key: key, Value: value, Type: model.PropertyTypeArrayNumber})
			}
		}
	}
	return extractedProperties
}

func RemoveSpacesAndHyphens(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(strings.TrimSpace(s), " ", "_"), "-", "_")
}

func RemoveSpacesAndLowerCase(s string) string {
	return strings.ReplaceAll(strings.TrimSpace(strings.ToLower(s)), " ", "_")
}

func RemoveSpacesAndUpperCase(s string) string {
	return strings.ReplaceAll(strings.TrimSpace(strings.ToUpper(s)), " ", "_")
}
