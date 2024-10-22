package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/mike-jacks/neo/model"
	"github.com/mike-jacks/neo/utils"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
)

// SetupNeo4jDriver creates a new Neo4j driver with the given context
func SetupNeo4jDriver() (neo4j.DriverWithContext, error) {
	dbUri := os.Getenv("NEO4J_URI")
	dbUser := os.Getenv("NEO4J_USERNAME")
	dbPassword := os.Getenv("NEO4J_PASSWORD")

	driver, err := neo4j.NewDriverWithContext(dbUri, neo4j.BasicAuth(dbUser, dbPassword, ""))
	if err != nil {
		return nil, err
	}

	if err := driver.VerifyConnectivity(context.Background()); err != nil {
		return nil, err
	}

	log.Println("Neo4j connection established")
	return driver, nil
}

type Neo4jDatabase struct {
	Driver neo4j.DriverWithContext
}

// Database interface implementation
func (db *Neo4jDatabase) GetDriver() neo4j.DriverWithContext {
	return db.Driver
}

func (db *Neo4jDatabase) CreateObjectNode(ctx context.Context, domain string, name string, typeArg string, labels []string, properties []*model.PropertyInput) (*model.Response, error) {
	session := db.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	domain = strings.Trim(strings.ToUpper(domain), " ")
	name = strings.Trim(strings.ToUpper(name), " ")
	typeArg = strings.Trim(strings.ToUpper(typeArg), " ")
	labelFromTypeArg := strings.ReplaceAll(typeArg, " ", "_")
	for i, label := range labels {
		labels[i] = strings.ReplaceAll(strings.Trim(strings.ToUpper(label), " "), " ", "_")
	}
	for i := range properties {
		properties[i].Key = strings.ReplaceAll(strings.Trim(strings.ToLower(properties[i].Key), " "), " ", "_")
	}

	query := fmt.Sprintf(`
		CREATE CONSTRAINT IF NOT EXISTS
		FOR (n:%v)
		REQUIRE (n.name, n.type, n.domain) IS NODE KEY
		`, labelFromTypeArg)

	_, err := session.Run(ctx, query, nil)
	if err != nil {
		return nil, err
	}

	for _, label := range labels {
		query = fmt.Sprintf(`
		CREATE CONSTRAINT IF NOT EXISTS
		FOR (n:%v)
		REQUIRE (n.name, n.type, n.domain) IS NODE KEY
		`, label)
		_, err = session.Run(ctx, query, nil)
		if err != nil {
			return nil, err
		}
	}

	query = fmt.Sprintf("CREATE (o:%v", labelFromTypeArg)
	for _, label := range labels {
		query += fmt.Sprintf(":%v", label)
	}
	query += " {name: $name, type: $typeArg, domain: $domain, "
	query = utils.CreatePropertiesQuery(query, properties)
	query = strings.TrimSuffix(query, ", ")
	query += "}) RETURN o"

	parameters := map[string]any{
		"name":    name,
		"typeArg": typeArg,
		"domain":  domain,
	}

	result, err := session.Run(ctx, query, parameters)
	if err != nil {
		return nil, err
	}

	if result.Next(ctx) {
		record := result.Record()
		node, ok := record.Get("o")
		if !ok {
			return nil, fmt.Errorf("failed to retrieve the created node")
		}

		neo4jNode, ok := node.(dbtype.Node)
		if !ok {
			return nil, fmt.Errorf("unexpected type for node: %T", node)
		}

		nodeProperties := make(map[string]interface{})
		for key, value := range neo4jNode.Props {
			nodeProperties[key] = value
		}

		data := []map[string]interface{}{}
		data = append(data, map[string]interface{}{
			"name":       utils.PopString(nodeProperties, "name"),
			"type":       utils.PopString(nodeProperties, "type"),
			"domain":     utils.PopString(nodeProperties, "domain"),
			"labels":     neo4jNode.Labels,
			"properties": nodeProperties,
		})
		message := "Object node created successfully"
		return &model.Response{Success: true, Message: &message, Data: data}, nil
	}

	return nil, err
}

func (db *Neo4jDatabase) UpdateObjectNode(ctx context.Context, domain string, name string, typeArg string, updateObjectNodeInput model.UpdateObjectNodeInput) (*model.Response, error) {
	session := db.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	domain = strings.Trim(strings.ToUpper(domain), " ")
	name = strings.Trim(strings.ToUpper(name), " ")
	typeArg = strings.Trim(strings.ToUpper(typeArg), " ")
	labelFromTypeArg := strings.ReplaceAll(typeArg, " ", "_")

	if updateObjectNodeInput.Properties != nil {
		for i, property := range updateObjectNodeInput.Properties {
			updateObjectNodeInput.Properties[i].Key = strings.ReplaceAll(strings.Trim(strings.ToLower(property.Key), " "), " ", "_")
		}
	}

	var newName string
	if updateObjectNodeInput.Name != nil {
		newName = strings.Trim(strings.ToUpper(*updateObjectNodeInput.Name), " ")
	}
	var newTypeArg string
	if updateObjectNodeInput.Type != nil {
		newTypeArg = strings.Trim(strings.ToUpper(*updateObjectNodeInput.Type), " ")
	}
	var newDomain string
	if updateObjectNodeInput.Domain != nil {
		newDomain = strings.Trim(strings.ToUpper(*updateObjectNodeInput.Domain), " ")
	}
	newLabelFromTypeArg := strings.ReplaceAll(newTypeArg, " ", "_")

	replaceLabelQuery := ""
	query := "MATCH (o{name: $name, type: $typeArg, domain: $domain})\n"
	if newName != "" || newTypeArg != "" || newDomain != "" || len(updateObjectNodeInput.Labels) > 0 || len(updateObjectNodeInput.Properties) > 0 {
		query += "SET "
	}
	if newName != "" {
		query += fmt.Sprintf("o.name = \"%v\", ", newName)
	}
	if newTypeArg != "" {
		query += fmt.Sprintf("o.type = \"%v\", ", newTypeArg)
		replaceLabelQuery += fmt.Sprintf("MATCH (o{name: $name, type: \"%v\", domain: $domain}) REMOVE o:%v SET o:%v RETURN o", newTypeArg, labelFromTypeArg, newLabelFromTypeArg)
	}
	if newDomain != "" {
		query += fmt.Sprintf("o.domain = \"%v\", ", newDomain)
	}
	if len(updateObjectNodeInput.Labels) > 0 {
		query += "o"
		for _, label := range updateObjectNodeInput.Labels {
			query += fmt.Sprintf(":%v", strings.Trim(strings.ToUpper(label), " "))
		}
		query += ", "
	}
	if len(updateObjectNodeInput.Properties) > 0 {
		query = utils.CreatePropertiesQuery(query, updateObjectNodeInput.Properties, "o")
	}
	query = strings.TrimSuffix(query, ", ")
	query += " RETURN o;"
	fmt.Println(query)
	fmt.Println(replaceLabelQuery)

	parameters := map[string]any{
		"name":    name,
		"typeArg": typeArg,
		"domain":  domain,
	}

	_, err := session.Run(ctx, query, parameters)
	if err != nil {
		message := "Failed to update object node"
		return &model.Response{Success: false, Message: &message, Data: nil}, err
	}

	result, err := session.Run(ctx, replaceLabelQuery, parameters)
	if err != nil {
		return nil, err
	}

	if result.Next(ctx) {
		record := result.Record()
		node, ok := record.Get("o")
		if !ok {
			message := "Failed to retrieve the updated node"
			return &model.Response{Success: false, Message: &message, Data: nil}, nil
		}

		neo4jNode, ok := node.(dbtype.Node)
		if !ok {
			return nil, fmt.Errorf("unexpected type for node: %T", node)
		}

		nodeProperties := make(map[string]interface{})
		for key, value := range neo4jNode.Props {
			nodeProperties[key] = value
		}

		data := []map[string]interface{}{}
		data = append(data, map[string]interface{}{
			"name":       utils.PopString(nodeProperties, "name"),
			"type":       utils.PopString(nodeProperties, "type"),
			"domain":     utils.PopString(nodeProperties, "domain"),
			"labels":     neo4jNode.Labels,
			"properties": nodeProperties,
		})
		message := "Object node updated successfully"
		return &model.Response{Success: true, Message: &message, Data: data}, nil
	}
	message := "Failed to update object node"

	return &model.Response{Success: false, Message: &message, Data: nil}, fmt.Errorf("failed to update object node")
}

func (db *Neo4jDatabase) DeleteObjectNode(ctx context.Context, domain string, name string, typeArg string) (*model.Response, error) {
	session := db.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	domain = strings.Trim(strings.ToUpper(domain), " ")
	name = strings.Trim(strings.ToUpper(name), " ")
	typeArg = strings.Trim(strings.ToUpper(typeArg), " ")

	query := fmt.Sprintf("MATCH (o:%v {name: $name, type: $typeArg, domain: $domain}) DETACH DELETE o RETURN count(o) as deletedCount", typeArg)
	parameters := map[string]any{
		"name":    name,
		"typeArg": typeArg,
		"domain":  domain,
	}

	fmt.Println(query)

	result, err := session.Run(ctx, query, parameters)
	if err != nil {
		message := "Failed to delete object node"
		return &model.Response{Success: false, Message: &message, Data: nil}, err
	}

	if result.Next(ctx) {
		record := result.Record()
		deletedCount, ok := record.Get("deletedCount")
		if !ok {
			message := "Failed to retrieve the deleted count"
			return &model.Response{Success: false, Message: &message, Data: nil}, nil
		}
		deletedCountInt := deletedCount.(int64)
		if deletedCountInt > 0 {
			message := fmt.Sprintf("Successfully deleted object node: %v", name)
			return &model.Response{Success: true, Message: &message, Data: nil}, nil
		} else {
			message := fmt.Sprintf("Object node: %v does not exist", name)
			return &model.Response{Success: false, Message: &message, Data: nil}, nil
		}
	}
	return nil, fmt.Errorf("failed to delete object node")
}

func (db *Neo4jDatabase) AddLabelsToObjectNode(ctx context.Context, domain string, name string, typeArg string, labels []string) (*model.Response, error) {
	session := db.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	domain = strings.Trim(strings.ToUpper(domain), " ")
	name = strings.Trim(strings.ToUpper(name), " ")
	typeArg = strings.Trim(strings.ToUpper(typeArg), " ")
	for i, label := range labels {
		labels[i] = strings.ReplaceAll(strings.Trim(strings.ToUpper(label), " "), " ", "_")
	}

	if len(labels) == 0 {
		message := "No labels provided"
		return &model.Response{Success: false, Message: &message, Data: nil}, nil
	}

	query := fmt.Sprintf("MATCH (o:%v {name: $name, type: $typeArg, domain: $domain}) SET o", typeArg)
	for _, label := range labels {
		query += fmt.Sprintf(":%v", label)
	}
	query += " RETURN o"

	parameters := map[string]any{
		"name":    name,
		"typeArg": typeArg,
		"domain":  domain,
	}

	fmt.Println(query)
	result, err := session.Run(ctx, query, parameters)
	if err != nil {
		message := "Failed to add labels to object node"
		return &model.Response{Success: false, Message: &message, Data: nil}, err
	}
	if result.Next(ctx) {
		record := result.Record()
		node, ok := record.Get("o")
		if !ok {
			message := "Failed to retrieve the updated node"
			return &model.Response{Success: false, Message: &message, Data: nil}, nil
		}
		neo4jNode, ok := node.(dbtype.Node)
		if !ok {
			return nil, fmt.Errorf("unexpected type for node: %T", node)
		}

		nodeProperties := make(map[string]interface{})
		for key, value := range neo4jNode.Props {
			nodeProperties[key] = value
		}
		data := []map[string]interface{}{}
		data = append(data, map[string]interface{}{
			"name":       utils.PopString(nodeProperties, "name"),
			"type":       utils.PopString(nodeProperties, "type"),
			"domain":     utils.PopString(nodeProperties, "domain"),
			"labels":     neo4jNode.Labels,
			"properties": nodeProperties,
		})
		message := "Labels added to object node successfully"
		return &model.Response{Success: true, Message: &message, Data: data}, nil
	}
	return nil, fmt.Errorf("failed to add labels to object node")
}

func (db *Neo4jDatabase) RemoveLabelsFromObjectNode(ctx context.Context, domain string, name string, typeArg string, labels []string) (*model.Response, error) {
	session := db.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	domain = strings.Trim(strings.ToUpper(domain), " ")
	name = strings.Trim(strings.ToUpper(name), " ")
	typeArg = strings.Trim(strings.ToUpper(typeArg), " ")
	typeArgAsLabel := strings.ReplaceAll(strings.Trim(strings.ToUpper(typeArg), " "), " ", "_")
	for i := 0; i < len(labels); i++ {
		label := strings.ReplaceAll(strings.Trim(strings.ToUpper(labels[i]), " "), " ", "_")
		if label == typeArgAsLabel {
			labels = append(labels[:i], labels[i+1:]...)
			i--
			continue
		}
		labels[i] = label
	}

	if len(labels) == 0 {
		message := "No labels provided"
		return &model.Response{Success: false, Message: &message, Data: nil}, nil
	}

	query := fmt.Sprintf("MATCH (o:%v {name: $name, type: $typeArg, domain: $domain}) REMOVE o", typeArgAsLabel)
	for _, label := range labels {
		query += fmt.Sprintf(":%v", label)
	}
	query += " RETURN o"

	parameters := map[string]any{
		"name":    name,
		"typeArg": typeArg,
		"domain":  domain,
	}

	fmt.Println(query)
	result, err := session.Run(ctx, query, parameters)
	if err != nil {
		message := "Failed to add labels to object node"
		return &model.Response{Success: false, Message: &message, Data: nil}, err
	}
	if result.Next(ctx) {
		record := result.Record()
		node, ok := record.Get("o")
		if !ok {
			message := "Failed to retrieve the updated node"
			return &model.Response{Success: false, Message: &message, Data: nil}, nil
		}
		neo4jNode, ok := node.(dbtype.Node)
		if !ok {
			return nil, fmt.Errorf("unexpected type for node: %T", node)
		}

		nodeProperties := make(map[string]interface{})
		for key, value := range neo4jNode.Props {
			nodeProperties[key] = value
		}
		data := []map[string]interface{}{}
		data = append(data, map[string]interface{}{
			"name":       utils.PopString(nodeProperties, "name"),
			"type":       utils.PopString(nodeProperties, "type"),
			"domain":     utils.PopString(nodeProperties, "domain"),
			"labels":     neo4jNode.Labels,
			"properties": nodeProperties,
		})
		message := "Labels removed from object node successfully"
		return &model.Response{Success: true, Message: &message, Data: data}, nil
	}
	return nil, fmt.Errorf("failed to remove labels from object node")
}

func (db *Neo4jDatabase) AddPropertiesToObjectNode(ctx context.Context, domain string, name string, typeArg string, properties []*model.PropertyInput) (*model.Response, error) {
	for _, property := range properties {
		property.Key = strings.ReplaceAll(strings.Trim(strings.ToLower(property.Key), " "), " ", "_")
	}
	session := db.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	domain = strings.Trim(strings.ToUpper(domain), " ")
	name = strings.Trim(strings.ToUpper(name), " ")
	typeArg = strings.Trim(strings.ToUpper(typeArg), " ")

	query := "MATCH (o{name: $name, type: $typeArg, domain: $domain}) SET o"

	for _, property := range properties {
		if property.Type == model.PropertyTypeString {
			query += fmt.Sprintf(".%v = \"%v\", ", property.Key, property.Value)
		} else {
			query += fmt.Sprintf(".%v = %v, ", property.Key, property.Value)
		}
	}
	query = strings.TrimSuffix(query, ", ")
	query += " RETURN o"

	parameters := map[string]any{
		"name":    name,
		"typeArg": typeArg,
		"domain":  domain,
	}

	fmt.Println(query)
	result, err := session.Run(ctx, query, parameters)
	if err != nil {
		return nil, err
	}

	if result.Next(ctx) {
		record := result.Record()
		node, ok := record.Get("o")
		if !ok {
			return nil, fmt.Errorf("failed to retrieve the updated node")
		}
		neo4jNode, ok := node.(dbtype.Node)
		if !ok {
			return nil, fmt.Errorf("unexpected type for node: %T", node)
		}

		nodeProperties := make(map[string]interface{})
		for key, value := range neo4jNode.Props {
			nodeProperties[key] = value
		}

		data := []map[string]interface{}{}
		data = append(data, map[string]interface{}{
			"name":       utils.PopString(nodeProperties, "name"),
			"type":       utils.PopString(nodeProperties, "type"),
			"domain":     utils.PopString(nodeProperties, "domain"),
			"labels":     neo4jNode.Labels,
			"properties": nodeProperties,
		})
		message := "Properties added to object node successfully"
		return &model.Response{Success: true, Message: &message, Data: data}, nil
	}
	return nil, fmt.Errorf("failed to add properties to object node")
}

func (db *Neo4jDatabase) RemovePropertiesFromObjectNode(ctx context.Context, domain string, name string, typeArg string, properties []string) (*model.Response, error) {
	for i, property := range properties {
		properties[i] = strings.ReplaceAll(strings.Trim(strings.ToLower(property), " "), " ", "_")
	}
	session := db.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	domain = strings.Trim(strings.ToUpper(domain), " ")
	name = strings.Trim(strings.ToUpper(name), " ")
	typeArg = strings.Trim(strings.ToUpper(typeArg), " ")

	query := "MATCH (o{name: $name, type: $typeArg, domain: $domain}) REMOVE "

	for _, property := range properties {
		query += fmt.Sprintf("o.%v, ", property)
	}
	query = strings.TrimSuffix(query, ", ")
	query += " RETURN o"

	parameters := map[string]any{
		"name":    name,
		"typeArg": typeArg,
		"domain":  domain,
	}

	fmt.Println(query)
	result, err := session.Run(ctx, query, parameters)
	if err != nil {
		return nil, err
	}

	if result.Next(ctx) {
		record := result.Record()
		node, ok := record.Get("o")
		if !ok {
			return nil, fmt.Errorf("failed to retrieve the updated node")
		}
		neo4jNode, ok := node.(dbtype.Node)
		if !ok {
			return nil, fmt.Errorf("unexpected type for node: %T", node)
		}

		nodeProperties := make(map[string]interface{})
		for key, value := range neo4jNode.Props {
			nodeProperties[key] = value
		}

		data := []map[string]interface{}{}
		data = append(data, map[string]interface{}{
			"name":       utils.PopString(nodeProperties, "name"),
			"type":       utils.PopString(nodeProperties, "type"),
			"domain":     utils.PopString(nodeProperties, "domain"),
			"labels":     neo4jNode.Labels,
			"properties": nodeProperties,
		})
		message := "Properties removed from object node successfully"
		return &model.Response{Success: true, Message: &message, Data: data}, nil
	}
	return nil, fmt.Errorf("failed to remove properties from object node")
}

func (db *Neo4jDatabase) GetObjectNode(ctx context.Context, domain string, name string, typeArg string) (*model.Response, error) {
	session := db.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	domain = strings.Trim(strings.ToUpper(domain), " ")
	name = strings.Trim(strings.ToUpper(name), " ")
	typeArg = strings.Trim(strings.ToUpper(typeArg), " ")

	query := "MATCH (o{name: $name, type: $typeArg, domain: $domain}) RETURN o"

	parameters := map[string]any{
		"name":    name,
		"typeArg": typeArg,
		"domain":  domain,
	}

	result, err := session.Run(ctx, query, parameters)
	if err != nil {
		return nil, err
	}

	if result.Next(ctx) {
		record := result.Record()
		node, ok := record.Get("o")
		if !ok {
			return nil, fmt.Errorf("failed to retrieve the node")
		}
		neo4jNode, ok := node.(dbtype.Node)
		if !ok {
			return nil, fmt.Errorf("unexpected type for node: %T", node)
		}

		nodeProperties := make(map[string]interface{})
		for key, value := range neo4jNode.Props {
			nodeProperties[key] = value
		}
		data := []map[string]interface{}{}
		data = append(data, map[string]interface{}{
			"name":       utils.PopString(nodeProperties, "name"),
			"type":       utils.PopString(nodeProperties, "type"),
			"domain":     utils.PopString(nodeProperties, "domain"),
			"labels":     neo4jNode.Labels,
			"properties": nodeProperties,
		})
		message := "Object node retrieved successfully"
		return &model.Response{Success: true, Message: &message, Data: data}, nil
	}
	return nil, fmt.Errorf("failed to get object node")
}

func (db *Neo4jDatabase) GetObjectNodes(ctx context.Context, domain *string, name *string, typeArg *string, labels []string) (*model.Response, error) {
	session := db.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	if domain != nil {
		*domain = strings.Trim(strings.ToUpper(*domain), " ")
	}
	if name != nil {
		*name = strings.Trim(strings.ToUpper(*name), " ")
	}
	if typeArg != nil {
		*typeArg = strings.Trim(strings.ToUpper(*typeArg), " ")
	}
	for i := 0; i < len(labels); i++ {
		labels[i] = strings.ReplaceAll(strings.Trim(strings.ToUpper(labels[i]), " "), " ", "_")
	}

	query := "MATCH (o"

	if len(labels) > 0 {
		for _, label := range labels {
			query += fmt.Sprintf(":%v", label)
		}
	}
	query += "{"
	if domain != nil {
		query += "domain: $domain, "
	}
	if name != nil {
		query += "name: $name, "
	}
	if typeArg != nil {
		query += "type: $typeArg, "
	}
	query = strings.TrimSuffix(query, ", ")
	query += "}) RETURN o"

	fmt.Println(query)

	parameters := map[string]any{}
	if domain != nil {
		parameters["domain"] = *domain
	}
	if name != nil {
		parameters["name"] = *name
	}
	if typeArg != nil {
		parameters["typeArg"] = *typeArg
	}

	result, err := session.Run(ctx, query, parameters)
	if err != nil {
		return nil, err
	}

	data := []map[string]interface{}{}
	for result.Next(ctx) {
		record := result.Record()
		node, ok := record.Get("o")
		if !ok {
			return nil, fmt.Errorf("failed to retrieve the node")
		}
		neo4jNode, ok := node.(dbtype.Node)
		if !ok {
			return nil, fmt.Errorf("unexpected type for node: %T", node)
		}

		nodeProperties := make(map[string]interface{})
		for key, value := range neo4jNode.Props {
			nodeProperties[key] = value
		}

		data = append(data, map[string]interface{}{
			"name":       utils.PopString(nodeProperties, "name"),
			"type":       utils.PopString(nodeProperties, "type"),
			"domain":     utils.PopString(nodeProperties, "domain"),
			"labels":     neo4jNode.Labels,
			"properties": nodeProperties,
		})
	}
	message := "Object nodes retrieved successfully"
	return &model.Response{Success: true, Message: &message, Data: data}, nil
}

func (db *Neo4jDatabase) CypherQuery(ctx context.Context, cypherStatement string) ([]*model.Response, error) {
	session := db.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	result, err := session.Run(ctx, cypherStatement, nil)
	if err != nil {
		return nil, err
	}

	data := []map[string]interface{}{}
	for result.Next(ctx) {
		record := result.Record()
		keys := record.Keys
		for _, key := range keys {
			node, ok := record.Get(key)
			if !ok {
				return nil, fmt.Errorf("failed to retrieve the node")
			}

			var nodeData map[string]interface{}

			switch v := node.(type) {
			case dbtype.Node:
				nodeProperties := make(map[string]interface{})
				for key, value := range v.Props {
					nodeProperties[key] = value
				}
				nodeData = map[string]interface{}{
					"name":       utils.PopString(nodeProperties, "name"),
					"type":       utils.PopString(nodeProperties, "type"),
					"domain":     utils.PopString(nodeProperties, "domain"),
					"labels":     v.Labels,
					"properties": nodeProperties,
				}
			case dbtype.Relationship:
				nodeData = map[string]interface{}{
					"type":       v.Type,
					"properties": v.Props,
				}
			default:
				return nil, fmt.Errorf("unexpected type for node: %T", node)
			}

			data = append(data, nodeData)
		}
	}
	message := "Cypher query executed successfully"
	return []*model.Response{{Success: true, Message: &message, Data: data}}, nil
}

func (db *Neo4jDatabase) CypherMutation(ctx context.Context, cypherStatement string) ([]*model.Response, error) {
	session := db.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	result, err := session.Run(ctx, cypherStatement, nil)
	if err != nil {
		return nil, err
	}

	data := []map[string]interface{}{}
	for result.Next(ctx) {
		record := result.Record()
		keys := record.Keys
		for _, key := range keys {
			node, ok := record.Get(key)
			if !ok {
				return nil, fmt.Errorf("failed to retrieve the node")
			}
			neo4jNode, ok := node.(dbtype.Node)
			if !ok {
				return nil, fmt.Errorf("unexpected type for node: %T", node)
			}
			nodeProperties := make(map[string]interface{})
			for key, value := range neo4jNode.Props {
				nodeProperties[key] = value
			}
			data = append(data, map[string]interface{}{
				"name":       utils.PopString(nodeProperties, "name"),
				"type":       utils.PopString(nodeProperties, "type"),
				"domain":     utils.PopString(nodeProperties, "domain"),
				"labels":     neo4jNode.Labels,
				"properties": nodeProperties,
			})
		}
	}
	message := "Cypher query executed successfully"
	return []*model.Response{{Success: true, Message: &message, Data: data}}, nil
}

func (db *Neo4jDatabase) CreateObjectRelationship(ctx context.Context, relationshipName string, properties []*model.PropertyInput, fromObjectNode model.ObjectNodeInput, toObjectNode model.ObjectNodeInput) (*model.Response, error) {
	session := db.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	fromObjectNode.Domain = strings.Trim(strings.ToUpper(fromObjectNode.Domain), " ")
	fromObjectNode.Name = strings.Trim(strings.ToUpper(fromObjectNode.Name), " ")
	fromObjectNode.Type = strings.Trim(strings.ToUpper(fromObjectNode.Type), " ")

	toObjectNode.Domain = strings.Trim(strings.ToUpper(toObjectNode.Domain), " ")
	toObjectNode.Name = strings.Trim(strings.ToUpper(toObjectNode.Name), " ")
	toObjectNode.Type = strings.Trim(strings.ToUpper(toObjectNode.Type), " ")
	relationshipName = strings.ReplaceAll(strings.Trim(strings.ToUpper(relationshipName), " "), " ", "_")

	for i, property := range properties {
		properties[i].Key = strings.ReplaceAll(strings.Trim(strings.ToLower(property.Key), " "), " ", "_")
	}

	query := fmt.Sprintf("MATCH (fromObjectNode{name: $fromName, type: $fromType, domain: $fromDomain}), (toObjectNode{name: $toName, type: $toType, domain: $toDomain}) MERGE (fromObjectNode)-[relationship:%v]->(toObjectNode)", relationshipName)
	if len(properties) > 0 {
		query += " SET "
		query = utils.CreatePropertiesQuery(query, properties, "relationship")
		query = strings.TrimSuffix(query, ", ")
	}
	query += " WITH toObjectNode, relationship, fromObjectNode RETURN toObjectNode, relationship, fromObjectNode"

	parameters := map[string]any{
		"fromName":   fromObjectNode.Name,
		"fromType":   fromObjectNode.Type,
		"fromDomain": fromObjectNode.Domain,
		"toName":     toObjectNode.Name,
		"toType":     toObjectNode.Type,
		"toDomain":   toObjectNode.Domain,
	}

	fmt.Println(query)

	result, err := session.Run(ctx, query, parameters)
	if err != nil {
		return nil, err
	}

	if result.Next(ctx) {
		record := result.Record()
		toObjectNode, ok := record.Get("toObjectNode")
		if !ok {
			return nil, fmt.Errorf("failed to retrieve the toObjectNode")
		}
		relationship, ok := record.Get("relationship")
		if !ok {
			return nil, fmt.Errorf("failed to retrieve the relationship")
		}
		fromObjectNode, ok := record.Get("fromObjectNode")
		if !ok {
			return nil, fmt.Errorf("failed to retrieve the fromObjectNode")
		}
		neo4jFromObjectNode, ok := fromObjectNode.(dbtype.Node)
		if !ok {
			return nil, fmt.Errorf("unexpected type for fromObjectNode: %T", fromObjectNode)
		}
		neo4jRelationship, ok := relationship.(dbtype.Relationship)
		if !ok {
			return nil, fmt.Errorf("unexpected type for relationship: %T", relationship)
		}
		neo4jToObjectNode, ok := toObjectNode.(dbtype.Node)
		if !ok {
			return nil, fmt.Errorf("unexpected type for toObjectNode: %T", toObjectNode)
		}
		data := []map[string]interface{}{}
		data = append(data, map[string]interface{}{
			"fromObjectNode": map[string]interface{}{
				"name":       utils.PopString(neo4jFromObjectNode.GetProperties(), "name"),
				"type":       utils.PopString(neo4jFromObjectNode.GetProperties(), "type"),
				"domain":     utils.PopString(neo4jFromObjectNode.GetProperties(), "domain"),
				"properties": neo4jFromObjectNode.GetProperties(),
				"labels":     neo4jFromObjectNode.Labels,
			},
			"relationship": map[string]interface{}{
				"relationshipName": neo4jRelationship.Type,
				"properties":       neo4jRelationship.GetProperties(),
			},
			"toObjectNode": map[string]interface{}{
				"name":       utils.PopString(neo4jToObjectNode.GetProperties(), "name"),
				"type":       utils.PopString(neo4jToObjectNode.GetProperties(), "type"),
				"domain":     utils.PopString(neo4jToObjectNode.GetProperties(), "domain"),
				"properties": neo4jToObjectNode.GetProperties(),
				"labels":     neo4jToObjectNode.Labels,
			},
		})
		message := "Object relationship created successfully"
		return &model.Response{Success: true, Message: &message, Data: data}, nil
	} else {
		message := "Object relationship creation failed"
		return &model.Response{Success: false, Message: &message, Data: nil}, nil
	}
}

func (db *Neo4jDatabase) UpdatePropertiesOnObjectRelationship(ctx context.Context, relationshipName string, properties []*model.PropertyInput, fromObjectNode model.ObjectNodeInput, toObjectNode model.ObjectNodeInput) (*model.Response, error) {
	session := db.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	relationshipName = strings.ReplaceAll(strings.Trim(strings.ToUpper(relationshipName), " "), " ", "_")
	if len(properties) > 0 {
		for i, property := range properties {
			properties[i].Key = strings.ReplaceAll(strings.Trim(strings.ToLower(property.Key), " "), " ", "_")
		}
	} else {
		return nil, fmt.Errorf("properties are required")
	}

	fromObjectNode.Domain = strings.Trim(strings.ToUpper(fromObjectNode.Domain), " ")
	fromObjectNode.Name = strings.Trim(strings.ToUpper(fromObjectNode.Name), " ")
	fromObjectNode.Type = strings.Trim(strings.ToUpper(fromObjectNode.Type), " ")

	toObjectNode.Domain = strings.Trim(strings.ToUpper(toObjectNode.Domain), " ")
	toObjectNode.Name = strings.Trim(strings.ToUpper(toObjectNode.Name), " ")
	toObjectNode.Type = strings.Trim(strings.ToUpper(toObjectNode.Type), " ")

	query := fmt.Sprintf("MATCH (fromObjectNode{name: $fromName, type: $fromType, domain: $fromDomain}), (toObjectNode{name: $toName, type: $toType, domain: $toDomain}) MATCH (fromObjectNode)-[relationship:%v]->(toObjectNode) SET ", relationshipName)
	query = utils.CreatePropertiesQuery(query, properties, "relationship")
	query = strings.TrimSuffix(query, ", ")
	query += " WITH toObjectNode, relationship, fromObjectNode RETURN toObjectNode, relationship, fromObjectNode"

	fmt.Println(query)

	parameters := map[string]any{
		"fromName":   fromObjectNode.Name,
		"fromType":   fromObjectNode.Type,
		"fromDomain": fromObjectNode.Domain,
		"toName":     toObjectNode.Name,
		"toType":     toObjectNode.Type,
		"toDomain":   toObjectNode.Domain,
	}

	result, err := session.Run(ctx, query, parameters)
	if err != nil {
		return nil, err
	}

	if result.Next(ctx) {
		record := result.Record()
		toObjectNode, ok := record.Get("toObjectNode")
		if !ok {
			return nil, fmt.Errorf("failed to retrieve the toObjectNode")
		}
		relationship, ok := record.Get("relationship")
		if !ok {
			return nil, fmt.Errorf("failed to retrieve the relationship")
		}
		fromObjectNode, ok := record.Get("fromObjectNode")
		if !ok {
			return nil, fmt.Errorf("failed to retrieve the fromObjectNode")
		}
		neo4jToObjectNode, ok := toObjectNode.(dbtype.Node)
		if !ok {
			return nil, fmt.Errorf("unexpected type for toObjectNode: %T", toObjectNode)
		}
		neo4jRelationship, ok := relationship.(dbtype.Relationship)
		if !ok {
			return nil, fmt.Errorf("unexpected type for relationship: %T", relationship)
		}
		neo4jFromObjectNode, ok := fromObjectNode.(dbtype.Node)
		if !ok {
			return nil, fmt.Errorf("unexpected type for fromObjectNode: %T", fromObjectNode)
		}
		data := []map[string]interface{}{}
		data = append(data, map[string]interface{}{
			"fromObjectNode": map[string]interface{}{
				"name":       utils.PopString(neo4jFromObjectNode.GetProperties(), "name"),
				"type":       utils.PopString(neo4jFromObjectNode.GetProperties(), "type"),
				"domain":     utils.PopString(neo4jFromObjectNode.GetProperties(), "domain"),
				"properties": neo4jFromObjectNode.GetProperties(),
				"labels":     neo4jFromObjectNode.Labels,
			},
			"relationship": map[string]interface{}{
				"relationshipName": neo4jRelationship.Type,
				"properties":       neo4jRelationship.GetProperties(),
			},
			"toObjectNode": map[string]interface{}{
				"name":       utils.PopString(neo4jToObjectNode.GetProperties(), "name"),
				"type":       utils.PopString(neo4jToObjectNode.GetProperties(), "type"),
				"domain":     utils.PopString(neo4jToObjectNode.GetProperties(), "domain"),
				"properties": neo4jToObjectNode.GetProperties(),
				"labels":     neo4jToObjectNode.Labels,
			},
		})
		message := "Object relationship created successfully"
		return &model.Response{Success: true, Message: &message, Data: data}, nil
	} else {
		message := "Object relationship update failed"
		return &model.Response{Success: false, Message: &message, Data: nil}, nil
	}
}

func (db *Neo4jDatabase) RemovePropertiesFromObjectRelationship(ctx context.Context, relationshipName string, properties []string, fromObjectNode model.ObjectNodeInput, toObjectNode model.ObjectNodeInput) (*model.Response, error) {
	session := db.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	relationshipName = strings.ReplaceAll(strings.Trim(strings.ToUpper(relationshipName), " "), " ", "_")
	if len(properties) > 0 {
		for i, property := range properties {
			properties[i] = strings.ReplaceAll(strings.Trim(strings.ToLower(property), " "), " ", "_")
		}
	} else {
		return nil, fmt.Errorf("properties are required")
	}

	fromObjectNode.Domain = strings.Trim(strings.ToUpper(fromObjectNode.Domain), " ")
	fromObjectNode.Name = strings.Trim(strings.ToUpper(fromObjectNode.Name), " ")
	fromObjectNode.Type = strings.Trim(strings.ToUpper(fromObjectNode.Type), " ")

	toObjectNode.Domain = strings.Trim(strings.ToUpper(toObjectNode.Domain), " ")
	toObjectNode.Name = strings.Trim(strings.ToUpper(toObjectNode.Name), " ")
	toObjectNode.Type = strings.Trim(strings.ToUpper(toObjectNode.Type), " ")

	query := fmt.Sprintf("MATCH (fromObjectNode{name: $fromName, type: $fromType, domain: $fromDomain}), (toObjectNode{name: $toName, type: $toType, domain: $toDomain}) MATCH (fromObjectNode)-[relationship:%v]->(toObjectNode) SET ", relationshipName)
	query = utils.RemovePropertiesQuery(query, properties, "relationship")
	query = strings.TrimSuffix(query, ", ")
	query += " WITH toObjectNode, relationship, fromObjectNode RETURN toObjectNode, relationship, fromObjectNode"

	fmt.Println(query)

	parameters := map[string]any{
		"fromName":   fromObjectNode.Name,
		"fromType":   fromObjectNode.Type,
		"fromDomain": fromObjectNode.Domain,
		"toName":     toObjectNode.Name,
		"toType":     toObjectNode.Type,
		"toDomain":   toObjectNode.Domain,
	}

	result, err := session.Run(ctx, query, parameters)
	if err != nil {
		return nil, err
	}

	if result.Next(ctx) {
		record := result.Record()
		toObjectNode, ok := record.Get("toObjectNode")
		if !ok {
			return nil, fmt.Errorf("failed to retrieve the toObjectNode")
		}
		relationship, ok := record.Get("relationship")
		if !ok {
			return nil, fmt.Errorf("failed to retrieve the relationship")
		}
		fromObjectNode, ok := record.Get("fromObjectNode")
		if !ok {
			return nil, fmt.Errorf("failed to retrieve the fromObjectNode")
		}
		neo4jToObjectNode, ok := toObjectNode.(dbtype.Node)
		if !ok {
			return nil, fmt.Errorf("unexpected type for toObjectNode: %T", toObjectNode)
		}
		neo4jRelationship, ok := relationship.(dbtype.Relationship)
		if !ok {
			return nil, fmt.Errorf("unexpected type for relationship: %T", relationship)
		}
		neo4jFromObjectNode, ok := fromObjectNode.(dbtype.Node)
		if !ok {
			return nil, fmt.Errorf("unexpected type for fromObjectNode: %T", fromObjectNode)
		}
		data := []map[string]interface{}{}
		data = append(data, map[string]interface{}{
			"fromObjectNode": map[string]interface{}{
				"name":       utils.PopString(neo4jFromObjectNode.GetProperties(), "name"),
				"type":       utils.PopString(neo4jFromObjectNode.GetProperties(), "type"),
				"domain":     utils.PopString(neo4jFromObjectNode.GetProperties(), "domain"),
				"properties": neo4jFromObjectNode.GetProperties(),
				"labels":     neo4jFromObjectNode.Labels,
			},
			"relationship": map[string]interface{}{
				"relationshipName": neo4jRelationship.Type,
				"properties":       neo4jRelationship.GetProperties(),
			},
			"toObjectNode": map[string]interface{}{
				"name":       utils.PopString(neo4jToObjectNode.GetProperties(), "name"),
				"type":       utils.PopString(neo4jToObjectNode.GetProperties(), "type"),
				"domain":     utils.PopString(neo4jToObjectNode.GetProperties(), "domain"),
				"properties": neo4jToObjectNode.GetProperties(),
				"labels":     neo4jToObjectNode.Labels,
			},
		})
		message := "Object relationship properties removed successfully"
		return &model.Response{Success: true, Message: &message, Data: data}, nil
	} else {
		message := "Object relationship properties removal failed"
		return &model.Response{Success: false, Message: &message, Data: nil}, nil
	}
}
