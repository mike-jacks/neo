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

func (db *Neo4jDatabase) CreateObjectNode(ctx context.Context, domain string, name string, typeArg string, labels []string, properties []*model.PropertyInput) (*model.ObjectNodeResponse, error) {
	session := db.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	id := utils.GenerateId()
	domain = strings.TrimSpace(domain)
	originalName := strings.TrimSpace(name)
	name = strings.TrimSpace(strings.ToUpper(name))
	typeArg = strings.TrimSpace(strings.ToUpper(typeArg))
	labelFromTypeArg := utils.RemoveSpacesAndHyphens(typeArg)
	for i, label := range labels {
		labels[i] = utils.RemoveSpacesAndHyphens(label)
	}

	if properties != nil {
		if err := utils.CleanUpPropertyObjects(&properties); err != nil {
			message := err.Error()
			return &model.ObjectNodeResponse{Success: false, Message: &message, ObjectNode: nil}, nil
		}
	}
	query := fmt.Sprintf(`
	CREATE CONSTRAINT object_node_%s_key IF NOT EXISTS
	FOR (n:%v) REQUIRE (n._id) IS NODE KEY
	`, utils.SanitizeStringToLower(labelFromTypeArg), utils.SanitizeStringToUpper(labelFromTypeArg))

	fmt.Println(query)

	_, err := session.Run(ctx, query, nil)
	if err != nil {
		return nil, err
	}

	query = fmt.Sprintf(`
		CREATE CONSTRAINT object_node_%s_unique IF NOT EXISTS
		FOR (n:%v)
		REQUIRE (n._name, n._type, n._domain) IS UNIQUE
		`, utils.SanitizeStringToLower(labelFromTypeArg), utils.SanitizeStringToUpper(labelFromTypeArg))

	fmt.Println(query)

	_, err = session.Run(ctx, query, nil)
	if err != nil {
		return nil, err
	}

	for _, label := range labels {
		query = fmt.Sprintf(`
		CREATE CONSTRAINT object_node_%s_key IF NOT EXISTS
		FOR (n:%v) REQUIRE (n._id) IS NODE KEY
		`, utils.SanitizeStringToLower(label), utils.SanitizeStringToUpper(label))

		fmt.Println(query)

		_, err = session.Run(ctx, query, nil)
		if err != nil {
			return nil, err
		}

		query = fmt.Sprintf(`
		CREATE CONSTRAINT object_node_%s_unique IF NOT EXISTS
		FOR (n:%v)
		REQUIRE (n._name, n._type, n._domain) IS UNIQUE
		`, utils.SanitizeStringToLower(label), utils.SanitizeStringToUpper(label))

		fmt.Println(query)

		_, err = session.Run(ctx, query, nil)
		if err != nil {
			return nil, err
		}
	}

	query = fmt.Sprintf("CREATE (objectNode:%v", utils.SanitizeStringToUpper(labelFromTypeArg))
	for _, label := range labels {
		query += fmt.Sprintf(":%v", utils.SanitizeStringToUpper(label))
	}
	query += " {_id: $id, _name: $name, _type: $typeArg, _domain: $domain, _originalName: $originalName, "
	query = utils.CreatePropertiesQuery(query, properties)
	query = strings.TrimSuffix(query, ", ")
	query += "}) RETURN objectNode"

	fmt.Println(query)

	parameters := map[string]any{
		"id":           id,
		"name":         name,
		"typeArg":      typeArg,
		"domain":       domain,
		"originalName": originalName,
	}

	result, err := session.Run(ctx, query, parameters)
	if err != nil {
		return nil, err
	}

	if result.Next(ctx) {
		record := result.Record()
		objectNode, ok := record.Get("objectNode")
		if !ok {
			return nil, fmt.Errorf("failed to retrieve the created node")
		}

		neo4jObjectNode, ok := objectNode.(dbtype.Node)
		if !ok {
			return nil, fmt.Errorf("unexpected type for node: %T", objectNode)
		}

		nodeProperties := make(map[string]interface{})
		for key, value := range neo4jObjectNode.Props {
			nodeProperties[key] = value
		}

		data := &model.ObjectNode{
			ID:           utils.PopString(nodeProperties, "_id"),
			Name:         utils.PopString(nodeProperties, "_name"),
			Type:         utils.PopString(nodeProperties, "_type"),
			Domain:       utils.PopString(nodeProperties, "_domain"),
			OriginalName: utils.PopString(nodeProperties, "_originalName"),
			Labels:       neo4jObjectNode.Labels,
			Properties:   utils.ExtractPropertiesFromNeo4jNode(nodeProperties),
		}
		message := "Object node created successfully"
		return &model.ObjectNodeResponse{Success: true, Message: &message, ObjectNode: data}, nil
	}

	message := "Failed to create object node"
	return &model.ObjectNodeResponse{Success: false, Message: &message, ObjectNode: nil}, nil
}

func (db *Neo4jDatabase) RenameObjectNode(ctx context.Context, id string, newName string) (*model.ObjectNodeResponse, error) {
	session := db.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	newOriginalName := strings.TrimSpace(newName)
	newName = strings.TrimSpace(strings.ToUpper(newName))

	query := fmt.Sprintf("MATCH (objectNode{_id: $id}) SET objectNode._name = \"%s\", objectNode._originalName = \"%s\" RETURN objectNode;", newName, newOriginalName)
	fmt.Println(query)

	parameters := map[string]any{
		"id": id,
	}

	result, err := session.Run(ctx, query, parameters)
	if err != nil {
		message := "Failed to update object node"
		return &model.ObjectNodeResponse{Success: false, Message: &message, ObjectNode: nil}, err
	}

	if result.Next(ctx) {
		record := result.Record()
		objectNode, ok := record.Get("objectNode")
		if !ok {
			message := "Failed to retrieve the updated node"
			return &model.ObjectNodeResponse{Success: false, Message: &message, ObjectNode: nil}, nil
		}

		neo4jObjectNode, ok := objectNode.(dbtype.Node)
		if !ok {
			message := "Unexpected type for node"
			return &model.ObjectNodeResponse{Success: false, Message: &message, ObjectNode: nil}, nil
		}

		data := &model.ObjectNode{
			ID:           utils.PopString(neo4jObjectNode.Props, "_id"),
			Name:         utils.PopString(neo4jObjectNode.Props, "_name"),
			Type:         utils.PopString(neo4jObjectNode.Props, "_type"),
			Domain:       utils.PopString(neo4jObjectNode.Props, "_domain"),
			OriginalName: utils.PopString(neo4jObjectNode.Props, "_originalName"),
			Labels:       neo4jObjectNode.Labels,
			Properties:   utils.ExtractPropertiesFromNeo4jNode(neo4jObjectNode.Props),
		}
		message := "Object node updated successfully"
		return &model.ObjectNodeResponse{Success: true, Message: &message, ObjectNode: data}, nil
	}
	message := "Failed to update object node"

	return &model.ObjectNodeResponse{Success: false, Message: &message, ObjectNode: nil}, fmt.Errorf("failed to update object node")
}

func (db *Neo4jDatabase) DeleteObjectNode(ctx context.Context, id string) (*model.ObjectNodeResponse, error) {
	session := db.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	query := "MATCH (objectNode{_id: $id}) WITH objectNode, count(objectNode) as deletedCount, objectNode._id as id DETACH DELETE objectNode RETURN id, deletedCount"
	parameters := map[string]any{
		"id": id,
	}

	fmt.Println(query)

	result, err := session.Run(ctx, query, parameters)
	if err != nil {
		message := "Failed to delete object node"
		return &model.ObjectNodeResponse{Success: false, Message: &message, ObjectNode: nil}, err
	}

	if result.Next(ctx) {
		record := result.Record()
		deletedCount, ok := record.Get("deletedCount")
		if !ok {
			message := "Failed to retrieve the deleted count"
			return &model.ObjectNodeResponse{Success: false, Message: &message, ObjectNode: nil}, nil
		}
		id, ok := record.Get("id")
		if !ok {
			message := "Failed to retrieve the deleted id"
			return &model.ObjectNodeResponse{Success: false, Message: &message, ObjectNode: nil}, nil
		}
		idString, ok := id.(string)
		if !ok {
			message := "Failed to retrieve the deleted id"
			return &model.ObjectNodeResponse{Success: false, Message: &message, ObjectNode: nil}, nil
		}
		deletedCountInt := deletedCount.(int64)
		data := &model.ObjectNode{
			ID: idString,
		}
		if deletedCountInt > 0 {
			message := "Successfully deleted object node."
			return &model.ObjectNodeResponse{Success: true, Message: &message, ObjectNode: data}, nil
		} else {
			message := fmt.Sprintf("Object node with id %v does not exist", idString)
			return &model.ObjectNodeResponse{Success: false, Message: &message, ObjectNode: nil}, nil
		}
	}
	message := "Failed to delete object node"
	return &model.ObjectNodeResponse{Success: false, Message: &message, ObjectNode: nil}, fmt.Errorf("failed to delete object node")
}

func (db *Neo4jDatabase) AddLabelsOnObjectNode(ctx context.Context, id string, labels []string) (*model.ObjectNodeResponse, error) {
	session := db.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	for i, label := range labels {
		labels[i] = utils.RemoveSpacesAndHyphens(label)
	}

	if len(labels) == 0 {
		message := "No labels provided"
		return &model.ObjectNodeResponse{Success: false, Message: &message, ObjectNode: nil}, nil
	}

	query := "MATCH (objectNode{_id: $id}) SET "
	for _, label := range labels {
		query += fmt.Sprintf("objectNode:%s, ", label)
	}
	query = strings.TrimSuffix(query, ", ")
	query += " RETURN objectNode"

	parameters := map[string]any{
		"id": id,
	}

	fmt.Println(query)

	result, err := session.Run(ctx, query, parameters)
	if err != nil {
		message := "Failed to add labels to object node"
		return &model.ObjectNodeResponse{Success: false, Message: &message, ObjectNode: nil}, err
	}
	if result.Next(ctx) {
		record := result.Record()
		objectNode, ok := record.Get("objectNode")
		if !ok {
			message := "Failed to retrieve the updated node"
			return &model.ObjectNodeResponse{Success: false, Message: &message, ObjectNode: nil}, nil
		}
		neo4jObjectNode, ok := objectNode.(dbtype.Node)
		if !ok {
			return nil, fmt.Errorf("unexpected type for node: %T", objectNode)
		}

		data := &model.ObjectNode{
			ID:           utils.PopString(neo4jObjectNode.Props, "_id"),
			Name:         utils.PopString(neo4jObjectNode.Props, "_name"),
			Type:         utils.PopString(neo4jObjectNode.Props, "_type"),
			Domain:       utils.PopString(neo4jObjectNode.Props, "_domain"),
			OriginalName: utils.PopString(neo4jObjectNode.Props, "_originalName"),
			Labels:       neo4jObjectNode.Labels,
			Properties:   utils.ExtractPropertiesFromNeo4jNode(neo4jObjectNode.Props),
		}
		message := "Labels added to object node successfully"
		return &model.ObjectNodeResponse{Success: true, Message: &message, ObjectNode: data}, nil
	}
	return nil, fmt.Errorf("failed to add labels to object node")
}

func (db *Neo4jDatabase) RemoveLabelsFromObjectNode(ctx context.Context, id string, labels []string) (*model.ObjectNodeResponse, error) {
	session := db.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	if len(labels) == 0 {
		message := "No labels provided"
		return &model.ObjectNodeResponse{Success: false, Message: &message, ObjectNode: nil}, nil
	}
	query := "MATCH (objectNode{_id: $id}) return objectNode"
	parameters := map[string]any{
		"id": id,
	}

	fmt.Println(query)

	result, err := session.Run(ctx, query, parameters)
	if err != nil {
		return nil, err
	}

	var current_type_label string
	if !result.Next(ctx) {
		return nil, fmt.Errorf("object node with id %v does not exist", id)
	} else {
		record := result.Record()
		objectNode, ok := record.Get("objectNode")
		if !ok {
			return nil, fmt.Errorf("failed to retrieve the updated node")
		}
		neo4jObjectNode, ok := objectNode.(dbtype.Node)
		if !ok {
			return nil, fmt.Errorf("unexpected type for node: %T", objectNode)
		}
		current_type_label = utils.RemoveSpacesAndHyphens(neo4jObjectNode.Props["_type"].(string))
	}

	query = "MATCH (objectNode{_id: $id}) REMOVE "
	for _, label := range labels {
		label = utils.RemoveSpacesAndHyphens(label)
		if label == current_type_label {
			continue
		}
		query += fmt.Sprintf("objectNode:%v, ", label)
	}
	query = strings.TrimSuffix(query, ", ")
	query += " RETURN objectNode"

	parameters = map[string]any{
		"id": id,
	}

	fmt.Println(query)
	result, err = session.Run(ctx, query, parameters)
	if err != nil {
		message := "Failed to add labels to object node"
		return &model.ObjectNodeResponse{Success: false, Message: &message, ObjectNode: nil}, err
	}
	if result.Next(ctx) {
		record := result.Record()
		objectNode, ok := record.Get("objectNode")
		if !ok {
			message := "Failed to retrieve the updated node"
			return &model.ObjectNodeResponse{Success: false, Message: &message, ObjectNode: nil}, nil
		}
		neo4jObjectNode, ok := objectNode.(dbtype.Node)
		if !ok {
			return nil, fmt.Errorf("unexpected type for node: %T", objectNode)
		}

		data := &model.ObjectNode{
			ID:           utils.PopString(neo4jObjectNode.Props, "_id"),
			Name:         utils.PopString(neo4jObjectNode.Props, "_name"),
			Type:         utils.PopString(neo4jObjectNode.Props, "_type"),
			Domain:       utils.PopString(neo4jObjectNode.Props, "_domain"),
			OriginalName: utils.PopString(neo4jObjectNode.Props, "_originalName"),
			Labels:       neo4jObjectNode.Labels,
			Properties:   utils.ExtractPropertiesFromNeo4jNode(neo4jObjectNode.Props),
		}
		message := "Labels removed from object node successfully"
		return &model.ObjectNodeResponse{Success: true, Message: &message, ObjectNode: data}, nil
	}
	return nil, fmt.Errorf("failed to remove labels from object node")
}

func (db *Neo4jDatabase) UpdatePropertiesOnObjectNode(ctx context.Context, id string, properties []*model.PropertyInput) (*model.ObjectNodeResponse, error) {
	session := db.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	if err := utils.CleanUpPropertyObjects(&properties); err != nil {
		message := err.Error()
		return &model.ObjectNodeResponse{Success: false, Message: &message, ObjectNode: nil}, nil
	}

	query := "MATCH (objectNode{_id: $id}) SET "

	for _, property := range properties {
		if property.Type == model.PropertyTypeString {
			query += fmt.Sprintf("objectNode.%v = \"%v\", ", property.Key, property.Value)
		} else {
			query += fmt.Sprintf("objectNode.%v = %v, ", property.Key, property.Value)
		}
	}
	query = strings.TrimSuffix(query, ", ")
	query += " RETURN objectNode"

	parameters := map[string]any{
		"id": id,
	}

	fmt.Println(query)

	result, err := session.Run(ctx, query, parameters)
	if err != nil {
		return nil, err
	}

	if result.Next(ctx) {
		record := result.Record()
		node, ok := record.Get("objectNode")
		if !ok {
			return nil, fmt.Errorf("failed to retrieve the updated node")
		}
		neo4jNode, ok := node.(dbtype.Node)
		if !ok {
			return nil, fmt.Errorf("unexpected type for node: %T", node)
		}

		data := &model.ObjectNode{
			ID:           utils.PopString(neo4jNode.Props, "_id"),
			Name:         utils.PopString(neo4jNode.Props, "_name"),
			Type:         utils.PopString(neo4jNode.Props, "_type"),
			Domain:       utils.PopString(neo4jNode.Props, "_domain"),
			OriginalName: utils.PopString(neo4jNode.Props, "_originalName"),
			Labels:       neo4jNode.Labels,
			Properties:   utils.ExtractPropertiesFromNeo4jNode(neo4jNode.Props),
		}
		message := "Properties added to object node successfully"
		return &model.ObjectNodeResponse{Success: true, Message: &message, ObjectNode: data}, nil
	}
	return nil, fmt.Errorf("failed to add properties to object node")
}

func (db *Neo4jDatabase) RemovePropertiesFromObjectNode(ctx context.Context, id string, properties []string) (*model.ObjectNodeResponse, error) {
	session := db.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	if err := utils.CleanUpPropertyKeys(&properties); err != nil {
		message := err.Error()
		return &model.ObjectNodeResponse{Success: false, Message: &message, ObjectNode: nil}, nil
	}

	query := "MATCH (objectNode{_id: $id}) REMOVE "

	for _, property := range properties {
		query += fmt.Sprintf("objectNode.%v, ", property)
	}
	query = strings.TrimSuffix(query, ", ")
	query += " RETURN objectNode"

	parameters := map[string]any{
		"id": id,
	}

	fmt.Println(query)
	result, err := session.Run(ctx, query, parameters)
	if err != nil {
		return nil, err
	}

	if result.Next(ctx) {
		record := result.Record()
		node, ok := record.Get("objectNode")
		if !ok {
			return nil, fmt.Errorf("failed to retrieve the updated node")
		}
		neo4jNode, ok := node.(dbtype.Node)
		if !ok {
			return nil, fmt.Errorf("unexpected type for node: %T", node)
		}

		data := &model.ObjectNode{
			ID:           utils.PopString(neo4jNode.Props, "_id"),
			Name:         utils.PopString(neo4jNode.Props, "_name"),
			Type:         utils.PopString(neo4jNode.Props, "_type"),
			Domain:       utils.PopString(neo4jNode.Props, "_domain"),
			OriginalName: utils.PopString(neo4jNode.Props, "_originalName"),
			Labels:       neo4jNode.Labels,
			Properties:   utils.ExtractPropertiesFromNeo4jNode(neo4jNode.Props),
		}
		message := "Properties removed from object node successfully"
		return &model.ObjectNodeResponse{Success: true, Message: &message, ObjectNode: data}, nil
	}
	return nil, fmt.Errorf("failed to remove properties from object node")
}

func (db *Neo4jDatabase) GetObjectNode(ctx context.Context, id string) (*model.ObjectNodeResponse, error) {
	session := db.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	query := "MATCH (objectNode{_id: $id}) RETURN objectNode"

	fmt.Println(query)

	parameters := map[string]any{
		"id": id,
	}

	result, err := session.Run(ctx, query, parameters)
	if err != nil {
		return nil, err
	}

	if result.Next(ctx) {
		record := result.Record()
		node, ok := record.Get("objectNode")
		if !ok {
			return nil, fmt.Errorf("failed to retrieve the object node")
		}
		neo4jNode, ok := node.(dbtype.Node)
		if !ok {
			return nil, fmt.Errorf("unexpected type for node: %T", node)
		}

		data := &model.ObjectNode{
			ID:           utils.PopString(neo4jNode.Props, "_id"),
			Name:         utils.PopString(neo4jNode.Props, "_name"),
			Type:         utils.PopString(neo4jNode.Props, "_type"),
			Domain:       utils.PopString(neo4jNode.Props, "_domain"),
			OriginalName: utils.PopString(neo4jNode.Props, "_originalName"),
			Labels:       neo4jNode.Labels,
			Properties:   utils.ExtractPropertiesFromNeo4jNode(neo4jNode.Props),
		}
		message := "Object node retrieved successfully"
		return &model.ObjectNodeResponse{Success: true, Message: &message, ObjectNode: data}, nil
	}
	return nil, fmt.Errorf("failed to get object node")
}

func (db *Neo4jDatabase) GetObjectNodes(ctx context.Context, domain *string, typeArg *string) (*model.ObjectNodesResponse, error) {
	session := db.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	if domain != nil {
		*domain = strings.TrimSpace(*domain)
	}
	if typeArg != nil {
		*typeArg = strings.TrimSpace(strings.ToUpper(*typeArg))
	}

	query := "MATCH (objectNode"

	query += "{"
	if domain != nil {
		query += "_domain: $domain, "
	}

	if typeArg != nil {
		query += "_type: $typeArg, "
	}

	query = strings.TrimSuffix(query, ", ")
	query += "}) WHERE NOT objectNode:RELATIONSHIP_SCHEMA AND NOT objectNode:DOMAIN_SCHEMA AND NOT objectNode:TYPE_SCHEMA RETURN objectNode"

	fmt.Println(query)

	parameters := map[string]any{}
	if domain != nil {
		parameters["domain"] = *domain
	}

	if typeArg != nil {
		parameters["typeArg"] = *typeArg
	}

	result, err := session.Run(ctx, query, parameters)
	if err != nil {
		return nil, err
	}

	data := []*model.ObjectNode{}
	for result.Next(ctx) {
		record := result.Record()
		node, ok := record.Get("objectNode")
		if !ok {
			return nil, fmt.Errorf("failed to retrieve the object node")
		}
		neo4jNode, ok := node.(dbtype.Node)
		if !ok {
			return nil, fmt.Errorf("unexpected type for node: %T", node)
		}

		data = append(data, &model.ObjectNode{
			ID:           utils.PopString(neo4jNode.Props, "_id"),
			Name:         utils.PopString(neo4jNode.Props, "_name"),
			Type:         utils.PopString(neo4jNode.Props, "_type"),
			Domain:       utils.PopString(neo4jNode.Props, "_domain"),
			OriginalName: utils.PopString(neo4jNode.Props, "_originalName"),
			Labels:       neo4jNode.Labels,
			Properties:   utils.ExtractPropertiesFromNeo4jNode(neo4jNode.Props),
		})
	}
	message := "Object nodes retrieved successfully"
	return &model.ObjectNodesResponse{Success: true, Message: &message, ObjectNodes: data}, nil
}

func (db *Neo4jDatabase) CypherQuery(ctx context.Context, cypherStatement string) (*model.ObjectNodesOrRelationshipNodesResponse, error) {
	return nil, nil
	// session := db.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	// defer session.Close(ctx)

	// result, err := session.Run(ctx, cypherStatement, nil)
	// if err != nil {
	// 	return nil, err
	// }

	// data := []map[string]interface{}{}
	// for result.Next(ctx) {
	// 	record := result.Record()
	// 	keys := record.Keys
	// 	for _, key := range keys {
	// 		node, ok := record.Get(key)
	// 		if !ok {
	// 			return nil, fmt.Errorf("failed to retrieve the node")
	// 		}

	// 		var nodeData *model.ObjectNode

	// 		switch v := node.(type) {
	// 		case dbtype.Node:
	// 			nodeProperties := make(map[string]interface{})
	// 			for key, value := range v.Props {
	// 				nodeProperties[key] = value
	// 			}
	// 			nodeData = map[string]interface{}{
	// 				"_name":         utils.PopString(nodeProperties, "_name"),
	// 				"_type":         utils.PopString(nodeProperties, "_type"),
	// 				"_domain":       utils.PopString(nodeProperties, "_domain"),
	// 				"_originalName": utils.PopString(nodeProperties, "_originalName"),
	// 				"_labels":       v.Labels,
	// 				"_properties":   nodeProperties,
	// 			}
	// 		case dbtype.Relationship:
	// 			nodeData = map[string]interface{}{
	// 				"_relationshipName": v.Type,
	// 				"_properties":       v.Props,
	// 			}
	// 		default:
	// 			return nil, fmt.Errorf("unexpected type for node: %T", node)
	// 		}

	// 		data = append(data, nodeData)
	// 	}
	// }
	// message := "Cypher query executed successfully"
	// return []*model.Response{{Success: true, Message: &message, Data: data}}, nil
}

func (db *Neo4jDatabase) CypherMutation(ctx context.Context, cypherStatement string) (*model.ObjectNodesOrRelationshipNodesResponse, error) {
	return nil, nil
	// session := db.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	// defer session.Close(ctx)

	// result, err := session.Run(ctx, cypherStatement, nil)
	// if err != nil {
	// 	return nil, err
	// }

	// data := []map[string]interface{}{}
	// for result.Next(ctx) {
	// 	record := result.Record()
	// 	keys := record.Keys
	// 	for _, key := range keys {
	// 		node, ok := record.Get(key)
	// 		if !ok {
	// 			return nil, fmt.Errorf("failed to retrieve the node")
	// 		}

	// 		var nodeData map[string]interface{}

	// 		switch v := node.(type) {
	// 		case dbtype.Node:
	// 			nodeProperties := make(map[string]interface{})
	// 			for key, value := range v.Props {
	// 				nodeProperties[key] = value
	// 			}
	// 			nodeData = map[string]interface{}{
	// 				"_name":         utils.PopString(nodeProperties, "_name"),
	// 				"_type":         utils.PopString(nodeProperties, "_type"),
	// 				"_domain":       utils.PopString(nodeProperties, "_domain"),
	// 				"_originalName": utils.PopString(nodeProperties, "_originalName"),
	// 				"_labels":       v.Labels,
	// 				"_properties":   nodeProperties,
	// 			}
	// 		case dbtype.Relationship:
	// 			nodeData = map[string]interface{}{
	// 				"_relationshipName":         v.Type,
	// 				"_originalRelationshipName": utils.PopString(v.Props, "_originalRelationshipName"),
	// 				"_properties":               v.Props,
	// 			}
	// 		default:
	// 			return nil, fmt.Errorf("unexpected type for node: %T", node)
	// 		}

	// 		data = append(data, nodeData)
	// 	}
	// }
	// message := "Cypher mutation executed successfully"
	// return []*model.Response{{Success: true, Message: &message, Data: data}}, nil
}

func (db *Neo4jDatabase) CreateObjectRelationship(ctx context.Context, name string, properties []*model.PropertyInput, fromObjectNodeId string, toObjectNodeId string) (*model.ObjectRelationshipResponse, error) {
	session := db.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	id := utils.GenerateId()
	originalName := strings.Trim(name, " ")
	name = utils.CleanUpRelationshipName(name)
	if name == "" {
		message := "Relationship name is required"
		return &model.ObjectRelationshipResponse{Success: false, Message: &message, ObjectRelationship: nil}, nil
	}

	if len(properties) > 0 {
		if err := utils.CleanUpPropertyObjects(&properties); err != nil {
			message := fmt.Sprintf("Unable to clean up properties. Error: %s", err.Error())
			return &model.ObjectRelationshipResponse{Success: false, Message: &message, ObjectRelationship: nil}, nil
		}
	}

	query := fmt.Sprintf("MATCH (fromObjectNode{_id: $fromObjectNodeId}), (toObjectNode{_id: $toObjectNodeId}) MERGE (fromObjectNode)-[relationship:%v {_id: $id, _name: $name, _originalName: $originalName, _fromObjectNodeId: $fromObjectNodeId, _toObjectNodeId: $toObjectNodeId}]->(toObjectNode)", name)
	if len(properties) > 0 {
		query += " SET "
		query = utils.CreatePropertiesQuery(query, properties, "relationship")
		query = strings.TrimSuffix(query, ", ")
	}
	query += " WITH relationship RETURN relationship"

	parameters := map[string]any{
		"id":               id,
		"name":             name,
		"originalName":     originalName,
		"fromObjectNodeId": fromObjectNodeId,
		"toObjectNodeId":   toObjectNodeId,
	}

	fmt.Println(query)

	result, err := session.Run(ctx, query, parameters)
	if err != nil {
		return nil, err
	}

	if result.Next(ctx) {
		record := result.Record()
		relationship, ok := record.Get("relationship")
		if !ok {
			return nil, fmt.Errorf("failed to retrieve the relationship")
		}
		neo4jRelationship, ok := relationship.(dbtype.Relationship)
		if !ok {
			return nil, fmt.Errorf("unexpected type for relationship: %T", relationship)
		}
		data := &model.ObjectRelationship{
			ID:               utils.PopString(neo4jRelationship.Props, "_id"),
			Name:             utils.PopString(neo4jRelationship.Props, "_name"),
			OriginalName:     utils.PopString(neo4jRelationship.Props, "_originalName"),
			FromObjectNodeID: utils.PopString(neo4jRelationship.Props, "_fromObjectNodeId"),
			ToObjectNodeID:   utils.PopString(neo4jRelationship.Props, "_toObjectNodeId"),
			Properties:       utils.ExtractPropertiesFromNeo4jNode(neo4jRelationship.Props),
		}
		message := "Object relationship created successfully"
		return &model.ObjectRelationshipResponse{Success: true, Message: &message, ObjectRelationship: data}, nil
	} else {
		message := "Object relationship creation failed"
		return &model.ObjectRelationshipResponse{Success: false, Message: &message, ObjectRelationship: nil}, nil
	}
}

func (db *Neo4jDatabase) UpdatePropertiesOnObjectRelationship(ctx context.Context, id string, properties []*model.PropertyInput) (*model.ObjectRelationshipResponse, error) {
	session := db.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	if err := utils.CleanUpPropertyObjects(&properties); err != nil {
		message := fmt.Sprintf("Unable to update properties. Error: %s", err.Error())
		return &model.ObjectRelationshipResponse{Success: false, Message: &message, ObjectRelationship: nil}, nil
	}

	query := "MATCH (fromObjectNode)-[relationship]->(toObjectNode) WHERE relationship._id = $id SET "
	query = utils.CreatePropertiesQuery(query, properties, "relationship")
	query = strings.TrimSuffix(query, ", ")
	query += " WITH relationship RETURN relationship"

	fmt.Println(query)

	parameters := map[string]any{
		"id": id,
	}

	result, err := session.Run(ctx, query, parameters)
	if err != nil {
		return nil, err
	}

	if result.Next(ctx) {
		record := result.Record()
		relationship, ok := record.Get("relationship")
		if !ok {
			return nil, fmt.Errorf("failed to retrieve the relationship")
		}
		neo4jRelationship, ok := relationship.(dbtype.Relationship)
		if !ok {
			return nil, fmt.Errorf("unexpected type for relationship: %T", relationship)
		}
		data := &model.ObjectRelationship{
			ID:               utils.PopString(neo4jRelationship.Props, "_id"),
			Name:             utils.PopString(neo4jRelationship.Props, "_name"),
			OriginalName:     utils.PopString(neo4jRelationship.Props, "_originalName"),
			FromObjectNodeID: utils.PopString(neo4jRelationship.Props, "_fromObjectNodeId"),
			ToObjectNodeID:   utils.PopString(neo4jRelationship.Props, "_toObjectNodeId"),
			Properties:       utils.ExtractPropertiesFromNeo4jNode(neo4jRelationship.Props),
		}
		message := "Object relationship properties updated successfully"
		return &model.ObjectRelationshipResponse{Success: true, Message: &message, ObjectRelationship: data}, nil
	} else {
		message := "Object relationship properties update failed"
		return &model.ObjectRelationshipResponse{Success: false, Message: &message, ObjectRelationship: nil}, nil
	}
}

func (db *Neo4jDatabase) RemovePropertiesFromObjectRelationship(ctx context.Context, id string, properties []string) (*model.ObjectRelationshipResponse, error) {
	session := db.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	if err := utils.CleanUpPropertyKeys(&properties); err != nil {
		message := err.Error()
		return &model.ObjectRelationshipResponse{Success: false, Message: &message, ObjectRelationship: nil}, nil
	}

	query := "MATCH (fromObjectNode)-[relationship]->(toObjectNode) WHERE relationship._id = $id SET "
	query = utils.RemovePropertiesQuery(query, properties, "relationship")
	query = strings.TrimSuffix(query, ", ")
	query += " WITH relationship RETURN relationship"

	fmt.Println(query)

	parameters := map[string]any{
		"id": id,
	}

	result, err := session.Run(ctx, query, parameters)
	if err != nil {
		return nil, err
	}

	if result.Next(ctx) {
		record := result.Record()
		relationship, ok := record.Get("relationship")
		if !ok {
			return nil, fmt.Errorf("failed to retrieve the relationship")
		}
		neo4jRelationship, ok := relationship.(dbtype.Relationship)
		if !ok {
			return nil, fmt.Errorf("unexpected type for relationship: %T", relationship)
		}
		data := &model.ObjectRelationship{
			ID:               utils.PopString(neo4jRelationship.Props, "_id"),
			Name:             utils.PopString(neo4jRelationship.Props, "_name"),
			OriginalName:     utils.PopString(neo4jRelationship.Props, "_originalName"),
			FromObjectNodeID: utils.PopString(neo4jRelationship.Props, "_fromObjectNodeId"),
			ToObjectNodeID:   utils.PopString(neo4jRelationship.Props, "_toObjectNodeId"),
			Properties:       utils.ExtractPropertiesFromNeo4jNode(neo4jRelationship.GetProperties()),
		}

		message := "Object relationship properties removed successfully"
		return &model.ObjectRelationshipResponse{Success: true, Message: &message, ObjectRelationship: data}, nil
	} else {
		message := "Object relationship properties removal failed"
		return &model.ObjectRelationshipResponse{Success: false, Message: &message, ObjectRelationship: nil}, nil
	}
}

func (db *Neo4jDatabase) DeleteObjectRelationship(ctx context.Context, id string) (*model.ObjectRelationshipResponse, error) {
	session := db.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	query := `
		MATCH (fromObjectNode)-[relationship]-(toObjectNode)
		WITH relationship, properties(relationship) as properties, fromObjectNode._id as fromObjectNodeId, toObjectNode._id as toObjectNodeId
		DELETE relationship
		RETURN properties, fromObjectNodeId, toObjectNodeId
	`

	fmt.Println(query)

	parameters := map[string]any{
		"id": id,
	}

	result, err := session.Run(ctx, query, parameters)
	if err != nil {
		return nil, err
	}

	if result.Next(ctx) {
		record := result.Record()
		properties, ok := record.Get("properties")
		if !ok {
			return nil, fmt.Errorf("failed to retrieve the relationship_id")
		}
		propertiesMap, ok := properties.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("unexpected type for relationship_id: %T", properties)
		}
		data := &model.ObjectRelationship{
			ID:               utils.PopString(propertiesMap, "_id"),
			Name:             utils.PopString(propertiesMap, "_name"),
			OriginalName:     utils.PopString(propertiesMap, "_originalName"),
			FromObjectNodeID: utils.PopString(propertiesMap, "_fromObjectNodeId"),
			ToObjectNodeID:   utils.PopString(propertiesMap, "_toObjectNodeId"),
			Properties:       utils.ExtractPropertiesFromNeo4jNode(propertiesMap),
		}
		message := "Object relationship deleted successfully"
		return &model.ObjectRelationshipResponse{Success: true, Message: &message, ObjectRelationship: data}, nil
	}
	message := "Object relationship deletion failed"
	return &model.ObjectRelationshipResponse{Success: false, Message: &message, ObjectRelationship: nil}, nil
}

func (db *Neo4jDatabase) GetObjectNodeRelationship(ctx context.Context, id string) (*model.ObjectRelationshipResponse, error) {
	session := db.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	query := `MATCH () - [relationship {_id: $id}]-> () RETURN relationship`

	fmt.Println(query)

	parameters := map[string]any{
		"id": id,
	}

	result, err := session.Run(ctx, query, parameters)
	if err != nil {
		return nil, err
	}

	if result.Next(ctx) {
		record := result.Record()
		relationship, ok := record.Get("relationship")
		if !ok {
			return nil, fmt.Errorf("failed to retrieve the relationship")
		}
		neo4jRelationship, ok := relationship.(dbtype.Relationship)
		if !ok {
			return nil, fmt.Errorf("unexpected type for relationship: %T", relationship)
		}
		data := &model.ObjectRelationship{
			ID:               utils.PopString(neo4jRelationship.Props, "_id"),
			Name:             utils.PopString(neo4jRelationship.Props, "_name"),
			OriginalName:     utils.PopString(neo4jRelationship.Props, "_originalName"),
			FromObjectNodeID: utils.PopString(neo4jRelationship.Props, "_fromObjectNodeId"),
			ToObjectNodeID:   utils.PopString(neo4jRelationship.Props, "_toObjectNodeId"),
			Properties:       utils.ExtractPropertiesFromNeo4jNode(neo4jRelationship.Props),
		}
		message := "Object relationship retrieved successfully"
		return &model.ObjectRelationshipResponse{Success: true, Message: &message, ObjectRelationship: data}, nil
	} else {
		message := "Object relationship retrieval failed"
		return &model.ObjectRelationshipResponse{Success: false, Message: &message, ObjectRelationship: nil}, nil
	}
}

func (db *Neo4jDatabase) GetObjectNodeOutgoingRelationships(ctx context.Context, fromObjectNodeId string) (*model.ObjectRelationshipsResponse, error) {
	session := db.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	query := ` MATCH (fromObjectNode {_id:$fromObjectNodeId}) - [relationship] -> () RETURN relationship`

	fmt.Println(query)

	parameters := map[string]any{
		"fromObjectNodeId": fromObjectNodeId,
	}

	result, err := session.Run(ctx, query, parameters)
	if err != nil {
		return nil, err
	}

	data := []*model.ObjectRelationship{}
	for result.Next(ctx) {
		record := result.Record()
		relationship, ok := record.Get("relationship")
		if !ok {
			return nil, fmt.Errorf("failed to retrieve the relationship")
		}
		neo4jRelationship, ok := relationship.(dbtype.Relationship)
		if !ok {
			return nil, fmt.Errorf("unexpected type for relationship: %T", relationship)
		}
		data = append(data, &model.ObjectRelationship{
			ID:               utils.PopString(neo4jRelationship.Props, "_id"),
			Name:             utils.PopString(neo4jRelationship.Props, "_name"),
			OriginalName:     utils.PopString(neo4jRelationship.Props, "_originalName"),
			FromObjectNodeID: utils.PopString(neo4jRelationship.Props, "_fromObjectNodeId"),
			ToObjectNodeID:   utils.PopString(neo4jRelationship.Props, "_toObjectNodeId"),
			Properties:       utils.ExtractPropertiesFromNeo4jNode(neo4jRelationship.Props),
		})
	}
	if len(data) == 0 {
		message := "No outgoing relationships found"
		return &model.ObjectRelationshipsResponse{Success: false, Message: &message, ObjectRelationships: data}, nil
	}
	message := "Object outgoing relationships retrieved successfully"
	return &model.ObjectRelationshipsResponse{Success: true, Message: &message, ObjectRelationships: data}, nil

}

func (db *Neo4jDatabase) GetObjectNodeIncomingRelationships(ctx context.Context, toObjectNodeId string) (*model.ObjectRelationshipsResponse, error) {
	session := db.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	query := ` MATCH () - [relationship] -> (toObjectNode{_id:$toObjectNodeId}) RETURN relationship`

	fmt.Println(query)

	parameters := map[string]any{
		"toObjectNodeId": toObjectNodeId,
	}

	result, err := session.Run(ctx, query, parameters)
	if err != nil {
		return nil, err
	}

	data := []*model.ObjectRelationship{}
	for result.Next(ctx) {
		record := result.Record()
		relationship, ok := record.Get("relationship")
		if !ok {
			return nil, fmt.Errorf("failed to retrieve the relationship")
		}
		neo4jRelationship, ok := relationship.(dbtype.Relationship)
		if !ok {
			return nil, fmt.Errorf("unexpected type for relationship: %T", relationship)
		}
		data = append(data, &model.ObjectRelationship{
			ID:               utils.PopString(neo4jRelationship.Props, "_id"),
			Name:             utils.PopString(neo4jRelationship.Props, "_name"),
			OriginalName:     utils.PopString(neo4jRelationship.Props, "_originalName"),
			FromObjectNodeID: utils.PopString(neo4jRelationship.Props, "_fromObjectNodeId"),
			ToObjectNodeID:   utils.PopString(neo4jRelationship.Props, "_toObjectNodeId"),
			Properties:       utils.ExtractPropertiesFromNeo4jNode(neo4jRelationship.Props),
		})
	}
	if len(data) == 0 {
		message := "No incoming relationships found"
		return &model.ObjectRelationshipsResponse{Success: false, Message: &message, ObjectRelationships: data}, nil
	}
	message := "Object incoming relationships retrieved successfully"
	return &model.ObjectRelationshipsResponse{Success: true, Message: &message, ObjectRelationships: data}, nil

}

func (db *Neo4jDatabase) GetDomainSchemaNode(ctx context.Context, id string) (*model.DomainSchemaNodeResponse, error) {
	session := db.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	query := `MATCH (schemaDomainNode:DOMAIN_SCHEMA {_id: $id}) RETURN schemaDomainNode`

	fmt.Println(query)

	parameters := map[string]any{
		"id": id,
	}

	result, err := session.Run(ctx, query, parameters)
	if err != nil {
		return nil, err
	}

	if result.Next(ctx) {
		record := result.Record()
		schemaDomainNode, ok := record.Get("schemaDomainNode")
		if !ok {
			return nil, fmt.Errorf("failed to retrieve the schemaDomainNode")
		}
		neo4jSchemaDomainNode, ok := schemaDomainNode.(dbtype.Node)
		if !ok {
			return nil, fmt.Errorf("unexpected type for schemaDomainNode: %T", schemaDomainNode)
		}
		data := &model.DomainSchemaNode{
			ID:         utils.PopString(neo4jSchemaDomainNode.Props, "_id"),
			Name:       utils.PopString(neo4jSchemaDomainNode.Props, "_name"),
			Type:       utils.PopString(neo4jSchemaDomainNode.Props, "_type"),
			Domain:     utils.PopString(neo4jSchemaDomainNode.Props, "_domain"),
			Properties: utils.ExtractPropertiesFromNeo4jNode(neo4jSchemaDomainNode.Props),
			Labels:     neo4jSchemaDomainNode.Labels,
		}
		message := "Domain schema node retrieved successfully"
		return &model.DomainSchemaNodeResponse{Success: true, Message: &message, DomainSchemaNode: data}, nil
	}
	message := fmt.Sprintf("Domain schema node with id'%s' not found", id)
	return &model.DomainSchemaNodeResponse{Success: false, Message: &message, DomainSchemaNode: nil}, nil
}

func (db *Neo4jDatabase) GetDomainSchemaNodes(ctx context.Context) (*model.DomainSchemaNodesResponse, error) {
	session := db.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	query := `
		MATCH (schemaDomainNode:DOMAIN_SCHEMA)
		RETURN schemaDomainNode
	`

	result, err := session.Run(ctx, query, nil)
	if err != nil {
		return nil, err
	}

	data := []*model.DomainSchemaNode{}
	for result.Next(ctx) {
		record := result.Record()
		schemaDomainNode, ok := record.Get("schemaDomainNode")
		if !ok {
			return nil, fmt.Errorf("failed to retrieve the schemaDomainNode")
		}
		neo4jSchemaDomainNode, ok := schemaDomainNode.(dbtype.Node)
		if !ok {
			return nil, fmt.Errorf("unexpected type for schemaDomainNode: %T", schemaDomainNode)
		}
		data = append(data, &model.DomainSchemaNode{
			ID:         utils.PopString(neo4jSchemaDomainNode.Props, "_id"),
			Name:       utils.PopString(neo4jSchemaDomainNode.Props, "_name"),
			Type:       utils.PopString(neo4jSchemaDomainNode.Props, "_type"),
			Domain:     utils.PopString(neo4jSchemaDomainNode.Props, "_domain"),
			Properties: utils.ExtractPropertiesFromNeo4jNode(neo4jSchemaDomainNode.Props),
			Labels:     neo4jSchemaDomainNode.Labels,
		})
	}
	if len(data) == 0 {
		message := "No schema domain nodes found"
		return &model.DomainSchemaNodesResponse{Success: false, Message: &message, DomainSchemaNodes: data}, nil
	}
	message := "Schema domain nodes retrieved successfully"
	return &model.DomainSchemaNodesResponse{Success: true, Message: &message, DomainSchemaNodes: data}, nil
}

func (db *Neo4jDatabase) CreateDomainSchemaNode(ctx context.Context, domain string) (*model.DomainSchemaNodeResponse, error) {
	session := db.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	id := utils.GenerateId()
	domain = strings.Trim(domain, " ")

	query := `
		CREATE CONSTRAINT domain_schema_node_key IF NOT EXISTS
		FOR (n:DOMAIN_SCHEMA)
		REQUIRE (n._id) IS NODE KEY
		`

	_, err := session.Run(ctx, query, nil)
	if err != nil {
		return nil, err
	}

	query = `
		CREATE CONSTRAINT domain_schema_node_unique IF NOT EXISTS
		FOR (n:DOMAIN_SCHEMA)
		REQUIRE (n._name, n._type, n._domain) IS UNIQUE
	`

	_, err = session.Run(ctx, query, nil)
	if err != nil {
		return nil, err
	}

	query = `
		CREATE (schemaDomainNode:DOMAIN_SCHEMA {_id: $id, _domain: $domain, _type: "DOMAIN SCHEMA", _name: $domain})
		RETURN schemaDomainNode
	`

	fmt.Println(query)

	parameters := map[string]any{
		"id":     id,
		"domain": domain,
	}

	result, err := session.Run(ctx, query, parameters)
	if err != nil {
		return nil, err
	}

	if result.Next(ctx) {
		record := result.Record()
		schemaDomainNode, ok := record.Get("schemaDomainNode")
		if !ok {
			return nil, fmt.Errorf("failed to retrieve the schemaDomainNode")
		}
		neo4jSchemaDomainNode, ok := schemaDomainNode.(dbtype.Node)
		if !ok {
			return nil, fmt.Errorf("unexpected type for schemaDomainNode: %T", schemaDomainNode)
		}
		data := &model.DomainSchemaNode{
			ID:         utils.PopString(neo4jSchemaDomainNode.Props, "_id"),
			Name:       utils.PopString(neo4jSchemaDomainNode.Props, "_name"),
			Type:       utils.PopString(neo4jSchemaDomainNode.Props, "_type"),
			Domain:     utils.PopString(neo4jSchemaDomainNode.Props, "_domain"),
			Properties: utils.ExtractPropertiesFromNeo4jNode(neo4jSchemaDomainNode.Props),
			Labels:     neo4jSchemaDomainNode.Labels,
		}
		message := "Domain schema node created successfully"
		return &model.DomainSchemaNodeResponse{Success: true, Message: &message, DomainSchemaNode: data}, nil
	}
	message := fmt.Sprintf("Domain schema node '%s' already exists", domain)
	return &model.DomainSchemaNodeResponse{Success: false, Message: &message, DomainSchemaNode: nil}, nil
}

func (db *Neo4jDatabase) RenameDomainSchemaNode(ctx context.Context, id string, newName string) (*model.DomainSchemaNodeResponse, error) {
	session := db.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	newName = strings.TrimSpace(newName)

	query := `
	MATCH (domainSchemaNode:DOMAIN_SCHEMA {_id: $id})
	WITH domainSchemaNode, domainSchemaNode._domain as originalDomainName
	OPTIONAL MATCH (node {_domain: originalDomainName}) WHERE NOT node:DOMAIN_SCHEMA
	WITH domainSchemaNode, originalDomainName, collect(node) as nodes
	CALL {
		WITH domainSchemaNode
		SET domainSchemaNode._domain = $newName
		SET domainSchemaNode._name = $newName
	}
	WITH domainSchemaNode, originalDomainName, nodes
	FOREACH (node IN nodes | SET node._domain = $newName)
	WITH domainSchemaNode, originalDomainName,
		[node IN nodes WHERE node:TYPE_SCHEMA | node] as typeSchemaNodes,
		[node IN nodes WHERE node:RELATIONSHIP_SCHEMA | node] as relationshipSchemaNodes,
		[node IN nodes WHERE NOT node:DOMAIN_SCHEMA AND NOT node:TYPE_SCHEMA AND NOT node:RELATIONSHIP_SCHEMA | node] as objectNodes
	RETURN domainSchemaNode,
		size(objectNodes) as objectNodeCount,
		size(typeSchemaNodes) as typeSchemaNodeCount,
		size(relationshipSchemaNodes) as relationshipSchemaNodeCount,
		originalDomainName
	`
	fmt.Println(query)

	parameters := map[string]any{
		"id":      id,
		"newName": newName,
	}

	result, err := session.Run(ctx, query, parameters)
	if err != nil {
		return nil, err
	}

	objectNodeCountInt := int64(0)
	typeSchemaNodeCountInt := int64(0)
	relationshipSchemaNodeCountInt := int64(0)

	if result.Next(ctx) {
		record := result.Record()
		domainSchemaNode, ok := record.Get("domainSchemaNode")
		if !ok {
			return nil, fmt.Errorf("failed to retrieve the domainSchemaNode")
		}
		neo4jDomainSchemaNode, ok := domainSchemaNode.(dbtype.Node)
		if !ok {
			return nil, fmt.Errorf("unexpected type for domainSchemaNode: %T", domainSchemaNode)
		}
		objectNodeCount, ok := record.Get("objectNodeCount")
		if !ok {
			return nil, fmt.Errorf("failed to retrieve the objectNodeCount")
		}
		objectNodeCountInt = objectNodeCount.(int64)
		typeSchemaNodeCount, ok := record.Get("typeSchemaNodeCount")
		if !ok {
			return nil, fmt.Errorf("failed to retrieve the typeSchemaNodeCount")
		}
		typeSchemaNodeCountInt = typeSchemaNodeCount.(int64)
		relationshipSchemaNodeCount, ok := record.Get("relationshipSchemaNodeCount")
		if !ok {
			return nil, fmt.Errorf("failed to retrieve the relationshipSchemaNodeCount")
		}
		relationshipSchemaNodeCountInt = relationshipSchemaNodeCount.(int64)
		originalDomainName, ok := record.Get("originalDomainName")
		if !ok {
			return nil, fmt.Errorf("failed to retrieve the originalDomainName")
		}
		neo4jOriginalDomainName, ok := originalDomainName.(string)
		if !ok {
			return nil, fmt.Errorf("unexpected type for originalDomainName: %T", originalDomainName)
		}
		data := &model.DomainSchemaNode{
			ID:         utils.PopString(neo4jDomainSchemaNode.Props, "_id"),
			Name:       utils.PopString(neo4jDomainSchemaNode.Props, "_name"),
			Type:       utils.PopString(neo4jDomainSchemaNode.Props, "_type"),
			Domain:     utils.PopString(neo4jDomainSchemaNode.Props, "_domain"),
			Properties: utils.ExtractPropertiesFromNeo4jNode(neo4jDomainSchemaNode.Props),
			Labels:     neo4jDomainSchemaNode.Labels,
		}
		message := fmt.Sprintf("Domain schema node %s renamed successfully to %s. %d object nodes, %d type schema nodes, and %d relationship schema nodes were affected.", neo4jOriginalDomainName, newName, objectNodeCountInt, typeSchemaNodeCountInt, relationshipSchemaNodeCountInt)
		return &model.DomainSchemaNodeResponse{Success: true, Message: &message, DomainSchemaNode: data}, nil
	}
	if result.Err() != nil {
		if strings.Contains(result.Err().Error(), "already exists") {
			message := fmt.Sprintf("Domain schema node %s already exists", newName)
			return &model.DomainSchemaNodeResponse{Success: false, Message: &message, DomainSchemaNode: nil}, nil
		}
		message := fmt.Sprintf("Domain schema with id %s does not exist", id)
		return &model.DomainSchemaNodeResponse{Success: false, Message: &message, DomainSchemaNode: nil}, nil

	}
	if result.Record() == nil {
		message := fmt.Sprintf("Domain schema with id %s does not exist", id)
		return &model.DomainSchemaNodeResponse{Success: false, Message: &message, DomainSchemaNode: nil}, nil
	}
	message := fmt.Sprintf("Error renaming domain schema node with id %s", id)
	return &model.DomainSchemaNodeResponse{Success: false, Message: &message, DomainSchemaNode: nil}, nil
}

func (db *Neo4jDatabase) DeleteDomainSchemaNode(ctx context.Context, id string) (*model.DomainSchemaNodeResponse, error) {
	session := db.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	query := `
	MATCH (domainSchemaNode:DOMAIN_SCHEMA {_id: $id})
	WITH domainSchemaNode
    OPTIONAL MATCH (node {_domain: domainSchemaNode._domain}) WHERE NOT node:DOMAIN_SCHEMA
    WITH domainSchemaNode,
        collect(CASE WHEN node:TYPE_SCHEMA THEN node END) as typeNodes,
        collect(CASE WHEN node:RELATIONSHIP_SCHEMA THEN node END) as relationshipNodes,
        collect(CASE WHEN NOT node:DOMAIN_SCHEMA AND NOT node:TYPE_SCHEMA AND NOT node:RELATIONSHIP_SCHEMA THEN node END) as objectNodes
    WITH
        {
        	properties: properties(domainSchemaNode),
            labels: [label in labels(domainSchemaNode) | toString(label)]
        } as storedDomainSchema, typeNodes, relationshipNodes, objectNodes, domainSchemaNode
	WITH storedDomainSchema,
        size(typeNodes) as typeCount,
        size(relationshipNodes) as relationshipCount,
        size(objectNodes) as objectCount,
		typeNodes + relationshipNodes + objectNodes + domainSchemaNode AS allNodesToDelete
    CALL {
        WITH allNodesToDelete
        UNWIND allNodesToDelete AS nodesToDelete
        DETACH DELETE nodesToDelete
    }
    RETURN
        storedDomainSchema,
        typeCount,
        relationshipCount,
        objectCount
`

	parameters := map[string]any{
		"id": id,
	}

	fmt.Println(query)

	result, err := session.Run(ctx, query, parameters)
	if err != nil {
		message := fmt.Sprintf("Domain schema node with id %s deletion failed: Error: %s", id, err.Error())
		return &model.DomainSchemaNodeResponse{Success: false, Message: &message, DomainSchemaNode: nil}, nil
	}
	var data *model.DomainSchemaNode
	typeCountInt := int64(0)
	relationshipCountInt := int64(0)
	objectCountInt := int64(0)
	if result.Next(ctx) {
		record := result.Record()
		domainSchemaNode, ok := record.Get("storedDomainSchema")
		if !ok {
			return nil, fmt.Errorf("failed to retrieve the domainSchemaNode")
		}
		neo4jDomainSchemaNode, ok := domainSchemaNode.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("unexpected type for domainSchemaNode: %T", domainSchemaNode)
		}
		typeCount, ok := record.Get("typeCount")
		if !ok {
			return nil, fmt.Errorf("failed to retrieve the typeCount")
		}
		typeCountInt, ok = typeCount.(int64)
		if !ok {
			return nil, fmt.Errorf("failed to retrieve the typeCount")
		}
		relationshipCount, ok := record.Get("relationshipCount")
		if !ok {
			return nil, fmt.Errorf("failed to retrieve the relationshipCount")
		}
		relationshipCountInt, ok = relationshipCount.(int64)
		if !ok {
			return nil, fmt.Errorf("failed to retrieve the relationshipCount")
		}
		objectCount, ok := record.Get("objectCount")
		if !ok {
			return nil, fmt.Errorf("failed to retrieve the objectCount")
		}
		objectCountInt, ok = objectCount.(int64)
		if !ok {
			return nil, fmt.Errorf("failed to retrieve the objectCount")
		}
		labels := []string{}
		for _, label := range neo4jDomainSchemaNode["labels"].([]interface{}) {
			labels = append(labels, label.(string))
		}
		properties := make(map[string]interface{})
		for key, value := range neo4jDomainSchemaNode["properties"].(map[string]interface{}) {
			if key != "_domain" && key != "_name" && key != "_type" {
				properties[key] = value
			}
		}
		data = &model.DomainSchemaNode{
			ID:         utils.PopString(neo4jDomainSchemaNode["properties"].(map[string]interface{}), "_id"),
			Domain:     utils.PopString(neo4jDomainSchemaNode["properties"].(map[string]interface{}), "_domain"),
			Name:       utils.PopString(neo4jDomainSchemaNode["properties"].(map[string]interface{}), "_name"),
			Type:       utils.PopString(neo4jDomainSchemaNode["properties"].(map[string]interface{}), "_type"),
			Properties: utils.ExtractPropertiesFromNeo4jNode(neo4jDomainSchemaNode["properties"].(map[string]interface{})),
			Labels:     labels,
		}
	}
	if data == nil {
		message := fmt.Sprintf("Domain schema node with id %s not found", id)
		return &model.DomainSchemaNodeResponse{Success: false, Message: &message, DomainSchemaNode: nil}, nil
	}
	message := fmt.Sprintf("Domain schema node %s deleted successfully. %d type nodes, %d relationship nodes, %d object nodes deleted.", data.Name, typeCountInt, relationshipCountInt, objectCountInt)
	return &model.DomainSchemaNodeResponse{Success: true, Message: &message, DomainSchemaNode: data}, nil
}

func (db *Neo4jDatabase) CreateTypeSchemaNode(ctx context.Context, domain string, name string) (*model.TypeSchemaNodeResponse, error) {
	session := db.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	id := utils.GenerateId()
	originalName := strings.TrimSpace(name)
	domain = strings.TrimSpace(domain)
	name = utils.RemoveSpacesAndUpperCase(name)

	query := `
		CREATE CONSTRAINT type_schema_node_key IF NOT EXISTS
		FOR (n:TYPE_SCHEMA)
		REQUIRE (n._id) IS NODE KEY
		`

	_, err := session.Run(ctx, query, nil)
	if err != nil {
		return nil, err
	}

	query = `
	CREATE CONSTRAINT type_schema_node_unique IF NOT EXISTS
	FOR (n:TYPE_SCHEMA)
	REQUIRE (n._name, n._type, n._domain) IS UNIQUE
	`

	_, err = session.Run(ctx, query, nil)
	if err != nil {
		return nil, err
	}

	query = `
		CREATE (schemaTypeNode:TYPE_SCHEMA {_id: $id, _domain: $domain, _type: "TYPE SCHEMA", _name: $name, _originalName: $originalName})
		RETURN schemaTypeNode
	`
	fmt.Println(query)

	parameters := map[string]any{
		"id":           id,
		"domain":       domain,
		"name":         name,
		"originalName": originalName,
	}

	result, err := session.Run(ctx, query, parameters)
	if err != nil {
		return nil, err
	}

	if result.Next(ctx) {
		record := result.Record()
		schemaTypeNode, ok := record.Get("schemaTypeNode")
		if !ok {
			return nil, fmt.Errorf("failed to retrieve the schemaTypeNode")
		}
		neo4jSchemaTypeNode, ok := schemaTypeNode.(dbtype.Node)
		if !ok {
			return nil, fmt.Errorf("unexpected type for schemaTypeNode: %T", schemaTypeNode)
		}
		data := &model.TypeSchemaNode{
			ID:           utils.PopString(neo4jSchemaTypeNode.Props, "_id"),
			Domain:       utils.PopString(neo4jSchemaTypeNode.Props, "_domain"),
			Name:         utils.PopString(neo4jSchemaTypeNode.Props, "_name"),
			OriginalName: utils.PopString(neo4jSchemaTypeNode.Props, "_originalName"),
			Type:         utils.PopString(neo4jSchemaTypeNode.Props, "_type"),
			Properties:   utils.ExtractPropertiesFromNeo4jNode(neo4jSchemaTypeNode.Props),
			Labels:       neo4jSchemaTypeNode.Labels,
		}
		message := "Schema type node created successfully"
		return &model.TypeSchemaNodeResponse{Success: true, Message: &message, TypeSchemaNode: data}, nil
	}
	message := "Schema type node creation failed"
	return &model.TypeSchemaNodeResponse{Success: false, Message: &message, TypeSchemaNode: nil}, nil
}

func (db *Neo4jDatabase) RenameTypeSchemaNode(ctx context.Context, id string, newName string) (*model.TypeSchemaNodeResponse, error) {
	session := db.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	originalNewName := strings.TrimSpace(newName)
	newName = strings.ToUpper(originalNewName)
	newLabel := utils.RemoveSpacesAndHyphens(newName)

	// Single query to check existence, update schema node and object nodes

	query := fmt.Sprintf(`
		MATCH (existingTypeSchemaNode:TYPE_SCHEMA {_id: $id})
		WITH existingTypeSchemaNode, existingTypeSchemaNode._domain as domain, existingTypeSchemaNode._name as existingName
		OPTIONAL MATCH (duplicateTypeSchemaNode:TYPE_SCHEMA {_domain: domain, _name: $newName, _type: "TYPE SCHEMA"})
		WITH existingTypeSchemaNode, domain, existingName, duplicateTypeSchemaNode
		WHERE duplicateTypeSchemaNode IS NULL
		SET existingTypeSchemaNode._name = $newName,
			existingTypeSchemaNode._originalName = $originalNewName
		WITH existingTypeSchemaNode, domain, existingName
		OPTIONAL MATCH (objectNodes {_domain: domain, _type: existingName})
		SET objectNodes._type = $newName
		REMOVE objectNodes:`+"`"+`${existingName}`+"`"+`
		SET objectNodes:%s
		WITH existingTypeSchemaNode as typeSchemaNode, count(objectNodes) as updatedCount, existingName as previousName
		RETURN typeSchemaNode, updatedCount, previousName
	`, newLabel)

	parameters := map[string]any{
		"id":              id,
		"newName":         newName,
		"originalNewName": originalNewName,
	}

	result, err := session.Run(ctx, query, parameters)
	if err != nil {
		return nil, err
	}

	if result.Next(ctx) {
		record := result.Record()
		typeSchemaNode, ok := record.Get("typeSchemaNode")
		if !ok {
			return nil, fmt.Errorf("failed to retrieve the typeSchemaNode")
		}
		neo4jTypeSchemaNode, ok := typeSchemaNode.(dbtype.Node)
		if !ok {
			return nil, fmt.Errorf("unexpected type for typeSchemaNode: %T", typeSchemaNode)
		}
		updatedCount, ok := record.Get("updatedCount")
		if !ok {
			return nil, fmt.Errorf("failed to retrieve the updatedCount")
		}
		updatedCountInt, ok := updatedCount.(int64)
		if !ok {
			return nil, fmt.Errorf("failed to retrieve the updatedCount")
		}
		previousName, ok := record.Get("previousName")
		if !ok {
			return nil, fmt.Errorf("failed to retrieve the previousName")
		}
		previousNameString, ok := previousName.(string)
		if !ok {
			return nil, fmt.Errorf("failed to retrieve the previousName")
		}
		data := &model.TypeSchemaNode{
			ID:           utils.PopString(neo4jTypeSchemaNode.Props, "_id"),
			Domain:       utils.PopString(neo4jTypeSchemaNode.Props, "_domain"),
			Name:         utils.PopString(neo4jTypeSchemaNode.Props, "_name"),
			OriginalName: utils.PopString(neo4jTypeSchemaNode.Props, "_originalName"),
			Type:         utils.PopString(neo4jTypeSchemaNode.Props, "_type"),
			Properties:   utils.ExtractPropertiesFromNeo4jNode(neo4jTypeSchemaNode.Props),
			Labels:       neo4jTypeSchemaNode.Labels,
		}
		message := fmt.Sprintf("Schema type node renamed from %s to %s, %d object nodes updated", previousNameString, data.Name, updatedCountInt)
		return &model.TypeSchemaNodeResponse{Success: true, Message: &message, TypeSchemaNode: data}, nil
	}

	// If we get here, either the duplicate exists or the original wasn't found
	message := fmt.Sprintf("Failed to rename schema type - either %s already exists or type schema node with id %s was not found", newName, id)
	return &model.TypeSchemaNodeResponse{Success: false, Message: &message, TypeSchemaNode: nil}, nil
}

func (db *Neo4jDatabase) UpdatePropertiesOnTypeSchemaNode(ctx context.Context, id string, properties []*model.PropertyInput) (*model.TypeSchemaNodeResponse, error) {
	session := db.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	err := utils.CleanUpPropertyObjects(&properties)
	if err != nil {
		message := fmt.Sprintf("Unable to update properties. Error: %s", err.Error())
		return &model.TypeSchemaNodeResponse{Success: false, Message: &message, TypeSchemaNode: nil}, err
	}

	query := `MATCH (typeSchemaNode:TYPE_SCHEMA {_id: $id}) SET `
	query = utils.CreatePropertiesQuery(query, properties, "typeSchemaNode")
	query = strings.TrimSuffix(query, ", ")
	query += ` RETURN typeSchemaNode`

	fmt.Println(query)

	parameters := map[string]any{
		"id": id,
	}

	result, err := session.Run(ctx, query, parameters)
	if err != nil {
		return nil, err
	}

	if result.Next(ctx) {
		record := result.Record()
		typeSchemaNode, ok := record.Get("typeSchemaNode")
		if !ok {
			return nil, fmt.Errorf("failed to retrieve the typeSchemaNode")
		}
		neo4jTypeSchemaNode, ok := typeSchemaNode.(dbtype.Node)
		if !ok {
			return nil, fmt.Errorf("unexpected type for typeSchemaNode: %T", typeSchemaNode)
		}
		data := &model.TypeSchemaNode{
			ID:           utils.PopString(neo4jTypeSchemaNode.Props, "_id"),
			Domain:       utils.PopString(neo4jTypeSchemaNode.Props, "_domain"),
			Name:         utils.PopString(neo4jTypeSchemaNode.Props, "_name"),
			OriginalName: utils.PopString(neo4jTypeSchemaNode.Props, "_originalName"),
			Type:         utils.PopString(neo4jTypeSchemaNode.Props, "_type"),
			Properties:   utils.ExtractPropertiesFromNeo4jNode(neo4jTypeSchemaNode.Props),
			Labels:       neo4jTypeSchemaNode.Labels,
		}
		message := "Type schema node properties updated successfully"
		return &model.TypeSchemaNodeResponse{Success: true, Message: &message, TypeSchemaNode: data}, nil
	}
	message := "Type schema node properties update failed"
	return &model.TypeSchemaNodeResponse{Success: false, Message: &message, TypeSchemaNode: nil}, nil
}

func (db *Neo4jDatabase) DeleteTypeSchemaNode(ctx context.Context, id string) (*model.TypeSchemaNodeResponse, error) {
	session := db.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	query := `
		MATCH (typeSchemaNode:TYPE_SCHEMA {_id: $id})
		WITH typeSchemaNode
		WHERE typeSchemaNode IS NOT NULL
		OPTIONAL MATCH (objectNodes {_domain: typeSchemaNode._domain, _type: typeSchemaNode._name})
		WITH typeSchemaNode, collect(objectNodes) as objectNodesToDelete, properties(typeSchemaNode) as typeSchemaNodeProperties, labels(typeSchemaNode) as typeSchemaNodeLabels
		WITH typeSchemaNode, objectNodesToDelete, typeSchemaNodeProperties, typeSchemaNodeLabels, size(objectNodesToDelete) as objectNodesCount
		FOREACH (node IN objectNodesToDelete | DETACH DELETE node)
		DETACH DELETE typeSchemaNode
		WITH typeSchemaNode, typeSchemaNodeProperties, typeSchemaNodeLabels, objectNodesCount
		RETURN objectNodesCount, typeSchemaNodeProperties, typeSchemaNodeLabels
	`

	fmt.Println(query)

	parameters := map[string]any{
		"id": id,
	}

	result, err := session.Run(ctx, query, parameters)
	if err != nil {
		return nil, err
	}

	if result.Next(ctx) {
		record := result.Record()
		typeSchemaNodeProperties, ok := record.Get("typeSchemaNodeProperties")
		if !ok {
			return nil, fmt.Errorf("failed to retrieve the typeSchemaNodeProperties")
		}
		typeSchemaNodePropertiesMap, ok := typeSchemaNodeProperties.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("unexpected type for typeSchemaNodeProperties: %T", typeSchemaNodeProperties)
		}
		objectNodesCount, ok := record.Get("objectNodesCount")
		if !ok {
			return nil, fmt.Errorf("failed to retrieve the objectNodesCount")
		}
		objectNodesCountInt, ok := objectNodesCount.(int64)
		if !ok {
			return nil, fmt.Errorf("unexpected type for objectNodesCount: %T", objectNodesCount)
		}
		typeSchemaNodeLabels, ok := record.Get("typeSchemaNodeLabels")
		if !ok {
			return nil, fmt.Errorf("failed to retrieve the typeSchemaNodeLabels")
		}
		typeSchemaNodeLabelsSliceAny, ok := typeSchemaNodeLabels.([]interface{})
		if !ok {
			return nil, fmt.Errorf("unexpected type for typeSchemaNodeLabels: %T", typeSchemaNodeLabels)
		}
		typeSchemaNodeLabelsSliceString := make([]string, len(typeSchemaNodeLabelsSliceAny))
		for i, label := range typeSchemaNodeLabelsSliceAny {
			typeSchemaNodeLabelsSliceString[i] = label.(string)
		}
		data := &model.TypeSchemaNode{
			ID:           utils.PopString(typeSchemaNodePropertiesMap, "_id"),
			Domain:       utils.PopString(typeSchemaNodePropertiesMap, "_domain"),
			Name:         utils.PopString(typeSchemaNodePropertiesMap, "_name"),
			OriginalName: utils.PopString(typeSchemaNodePropertiesMap, "_originalName"),
			Type:         utils.PopString(typeSchemaNodePropertiesMap, "_type"),
			Properties:   utils.ExtractPropertiesFromNeo4jNode(typeSchemaNodePropertiesMap),
			Labels:       typeSchemaNodeLabelsSliceString,
		}
		message := fmt.Sprintf("Type schema node '%s' deleted successfully. %v object nodes deleted successfully", data.Name, objectNodesCountInt)
		return &model.TypeSchemaNodeResponse{Success: true, Message: &message, TypeSchemaNode: data}, nil
	}
	message := "Unable to delete type schema node"
	return &model.TypeSchemaNodeResponse{Success: false, Message: &message, TypeSchemaNode: nil}, nil

}

func (db *Neo4jDatabase) RemovePropertiesFromTypeSchemaNode(ctx context.Context, id string, properties []string) (*model.TypeSchemaNodeResponse, error) {
	session := db.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	if err := utils.CleanUpPropertyKeys(&properties); err != nil {
		return nil, err
	}

	query := `MATCH (schemaTypeNode:TYPE_SCHEMA {_id: $id}) `
	query += `OPTIONAL MATCH (objectNodes {_domain: schemaTypeNode._domain, _type: schemaTypeNode._name}) WHERE NOT objectNodes:RELATIONSHIP_SCHEMA AND NOT objectNodes:DOMAIN_SCHEMA AND NOT objectNodes:TYPE_SCHEMA SET `
	query = utils.RemovePropertiesQuery(query, properties, "schemaTypeNode")
	query = utils.RemovePropertiesQuery(query, properties, "objectNodes")
	query = strings.TrimSuffix(query, "SET ")
	query = strings.TrimSuffix(query, ", ")
	query += ` WITH schemaTypeNode, count(objectNodes) as count`
	query += ` RETURN schemaTypeNode, count`

	fmt.Println(query)

	parameters := map[string]any{
		"id": id,
	}

	result, err := session.Run(ctx, query, parameters)
	if err != nil {
		return nil, err
	}

	if result.Next(ctx) {
		record := result.Record()
		schemaTypeNode, ok := record.Get("schemaTypeNode")
		if !ok {
			return nil, fmt.Errorf("failed to retrieve the schemaTypeNode")
		}
		neo4jSchemaTypeNode, ok := schemaTypeNode.(dbtype.Node)
		if !ok {
			return nil, fmt.Errorf("unexpected type for schemaTypeNode: %T", schemaTypeNode)
		}
		count, _ := record.Get("count")
		if !ok {
			return nil, fmt.Errorf("failed to retrieve the count")
		}
		countInt, ok := count.(int64)
		if !ok {
			return nil, fmt.Errorf("unexpected type for count: %T", count)
		}

		data := &model.TypeSchemaNode{
			ID:           utils.PopString(neo4jSchemaTypeNode.Props, "_id"),
			Domain:       utils.PopString(neo4jSchemaTypeNode.Props, "_domain"),
			Name:         utils.PopString(neo4jSchemaTypeNode.Props, "_name"),
			OriginalName: utils.PopString(neo4jSchemaTypeNode.Props, "_originalName"),
			Type:         utils.PopString(neo4jSchemaTypeNode.Props, "_type"),
			Properties:   utils.ExtractPropertiesFromNeo4jNode(neo4jSchemaTypeNode.Props),
			Labels:       neo4jSchemaTypeNode.Labels,
		}
		message := fmt.Sprintf("%v properties removed from schema type node of type %s. %v object nodes updated successfully", len(properties), data.Name, countInt)
		return &model.TypeSchemaNodeResponse{Success: true, Message: &message, TypeSchemaNode: data}, nil
	}
	message := fmt.Sprintf("Unable to remove properties from schema type node %s", id)
	return &model.TypeSchemaNodeResponse{Success: false, Message: &message, TypeSchemaNode: nil}, nil
}

func (db *Neo4jDatabase) GetTypeSchemaNodes(ctx context.Context, domain *string) (*model.TypeSchemaNodesResponse, error) {
	session := db.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	query := ``
	if domain != nil {
		trimmedDomain := strings.Trim(*domain, " ")
		domain = &trimmedDomain
		query = `
		MATCH (schemaTypeNode:TYPE_SCHEMA {_domain: $domain})
		RETURN schemaTypeNode
	`
	} else {
		query = `
		MATCH (schemaTypeNode:TYPE_SCHEMA)
		RETURN schemaTypeNode
	`
	}

	fmt.Println(query)

	parameters := map[string]any{
		"domain": domain,
	}

	result, err := session.Run(ctx, query, parameters)
	if err != nil {
		return nil, err
	}

	data := []*model.TypeSchemaNode{}
	for result.Next(ctx) {
		record := result.Record()
		schemaTypeNode, ok := record.Get("schemaTypeNode")
		if !ok {
			return nil, fmt.Errorf("failed to retrieve the schemaTypeNode")
		}
		neo4jSchemaTypeNode, ok := schemaTypeNode.(dbtype.Node)
		if !ok {
			return nil, fmt.Errorf("unexpected type for schemaTypeNode: %T", schemaTypeNode)
		}
		data = append(data, &model.TypeSchemaNode{
			ID:           utils.PopString(neo4jSchemaTypeNode.Props, "_id"),
			Name:         utils.PopString(neo4jSchemaTypeNode.Props, "_name"),
			Type:         utils.PopString(neo4jSchemaTypeNode.Props, "_type"),
			Domain:       utils.PopString(neo4jSchemaTypeNode.Props, "_domain"),
			OriginalName: utils.PopString(neo4jSchemaTypeNode.Props, "_originalName"),
			Properties:   utils.ExtractPropertiesFromNeo4jNode(neo4jSchemaTypeNode.Props),
			Labels:       neo4jSchemaTypeNode.Labels,
		})
	}
	if len(data) == 0 {
		message := "No schema type nodes found"
		return &model.TypeSchemaNodesResponse{Success: false, Message: &message, TypeSchemaNodes: nil}, nil
	}
	message := "Schema type nodes retrieved successfully"
	return &model.TypeSchemaNodesResponse{Success: true, Message: &message, TypeSchemaNodes: data}, nil
}

func (db *Neo4jDatabase) GetTypeSchemaNode(ctx context.Context, id string) (*model.TypeSchemaNodeResponse, error) {
	session := db.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	query := `MATCH (schemaTypeNode:TYPE_SCHEMA {_id: $id}) RETURN schemaTypeNode`

	fmt.Println(query)

	parameters := map[string]any{
		"id": id,
	}

	result, err := session.Run(ctx, query, parameters)
	if err != nil {
		return nil, err
	}

	if result.Next(ctx) {
		record := result.Record()
		schemaTypeNode, ok := record.Get("schemaTypeNode")
		if !ok {
			return nil, fmt.Errorf("failed to retrieve the schemaTypeNode")
		}
		neo4jSchemaTypeNode, ok := schemaTypeNode.(dbtype.Node)
		if !ok {
			return nil, fmt.Errorf("unexpected type for schemaTypeNode: %T", schemaTypeNode)
		}
		data := &model.TypeSchemaNode{
			ID:           utils.PopString(neo4jSchemaTypeNode.Props, "_id"),
			Name:         utils.PopString(neo4jSchemaTypeNode.Props, "_name"),
			Type:         utils.PopString(neo4jSchemaTypeNode.Props, "_type"),
			Domain:       utils.PopString(neo4jSchemaTypeNode.Props, "_domain"),
			OriginalName: utils.PopString(neo4jSchemaTypeNode.Props, "_originalName"),
			Properties:   utils.ExtractPropertiesFromNeo4jNode(neo4jSchemaTypeNode.Props),
			Labels:       neo4jSchemaTypeNode.Labels,
		}
		message := "Schema type node retrieved successfully"
		return &model.TypeSchemaNodeResponse{Success: true, Message: &message, TypeSchemaNode: data}, nil
	}
	message := "Unable to retrieve schema type node"
	return &model.TypeSchemaNodeResponse{Success: false, Message: &message, TypeSchemaNode: nil}, nil
}

func (db *Neo4jDatabase) RenamePropertyOnTypeSchemaNode(ctx context.Context, id string, oldPropertyName string, newPropertyName string) (*model.TypeSchemaNodeResponse, error) {
	session := db.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	oldPropertyName = strings.ReplaceAll(strings.TrimSpace(strings.ToLower(oldPropertyName)), " ", "_")
	newPropertyName = strings.ReplaceAll(strings.TrimSpace(strings.ToLower(newPropertyName)), " ", "_")

	if strings.HasPrefix(oldPropertyName, "_") {
		return nil, fmt.Errorf("oldPropertyName cannot be %s", oldPropertyName)
	}

	if strings.HasPrefix(newPropertyName, "_") {
		return nil, fmt.Errorf("newPropertyName cannot be %s", newPropertyName)
	}

	query := fmt.Sprintf(`MATCH (schemaTypeNode:TYPE_SCHEMA {_id: $id}) WHERE schemaTypeNode.%s IS NULL `, newPropertyName)
	query += `OPTIONAL MATCH (objectNodes {_domain: schemaTypeNode._domain, _type: schemaTypeNode._name}) SET `
	query = utils.RenamePropertyQuery(query, oldPropertyName, newPropertyName, "schemaTypeNode")
	query = utils.RenamePropertyQuery(query, oldPropertyName, newPropertyName, "objectNodes")
	query = strings.TrimSuffix(query, ", ")
	query += ` WITH schemaTypeNode, count(objectNodes) as count`
	query += ` RETURN schemaTypeNode, count`

	fmt.Println(query)

	parameters := map[string]any{
		"id":              id,
		"oldPropertyName": oldPropertyName,
		"newPropertyName": newPropertyName,
	}

	result, err := session.Run(ctx, query, parameters)
	if err != nil {
		return nil, err
	}

	countInt := int64(0)
	if result.Next(ctx) {
		record := result.Record()
		schemaTypeNode, ok := record.Get("schemaTypeNode")
		if !ok {
			return nil, fmt.Errorf("failed to retrieve the schemaTypeNode")
		}
		neo4jSchemaTypeNode, ok := schemaTypeNode.(dbtype.Node)
		if !ok {
			return nil, fmt.Errorf("unexpected type for schemaTypeNode: %T", schemaTypeNode)
		}
		count, ok := record.Get("count")
		if !ok {
			return nil, fmt.Errorf("failed to retrieve the count")
		}
		countInt, ok = count.(int64)
		if !ok {
			return nil, fmt.Errorf("unexpected type for count: %T", count)
		}
		data := &model.TypeSchemaNode{
			ID:           utils.PopString(neo4jSchemaTypeNode.GetProperties(), "_id"),
			Name:         utils.PopString(neo4jSchemaTypeNode.GetProperties(), "_name"),
			Type:         utils.PopString(neo4jSchemaTypeNode.GetProperties(), "_type"),
			Domain:       utils.PopString(neo4jSchemaTypeNode.GetProperties(), "_domain"),
			OriginalName: utils.PopString(neo4jSchemaTypeNode.GetProperties(), "_originalName"),
			Properties:   utils.ExtractPropertiesFromNeo4jNode(neo4jSchemaTypeNode.Props),
			Labels:       neo4jSchemaTypeNode.Labels,
		}
		message := fmt.Sprintf("%s property renamed to %s on schema type node of type %s. %v object nodes updated successfully", oldPropertyName, newPropertyName, data.Name, countInt)
		return &model.TypeSchemaNodeResponse{Success: true, Message: &message, TypeSchemaNode: data}, nil
	}
	message := fmt.Sprintf("Unable to rename property '%s' on schema type node with id '%s'. Either the new property '%s' already exists or the node id '%s' is not valid.", oldPropertyName, id, newPropertyName, id)
	return &model.TypeSchemaNodeResponse{Success: false, Message: &message, TypeSchemaNode: nil}, nil
}

func (db *Neo4jDatabase) CreateRelationshipSchemaNode(ctx context.Context, name string, domain string, fromTypeSchemaNodeId string, toTypeSchemaNodeId string) (*model.RelationshipSchemaNodeResponse, error) {
	session := db.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	id := utils.GenerateId()
	domain = strings.TrimSpace(domain)
	originalName := strings.TrimSpace(name)
	name = utils.RemoveSpacesAndHyphens(strings.ToUpper(name))

	query := `
		CREATE CONSTRAINT relationship_schema_node_key IF NOT EXISTS
		FOR (n:RELATIONSHIP_SCHEMA)
		REQUIRE (n._id) IS NODE KEY
		`

	fmt.Println(query)

	_, err := session.Run(ctx, query, nil)
	if err != nil {
		message := fmt.Sprintf("Unable to create relationship schema constraint. Error: %s", err.Error())
		return &model.RelationshipSchemaNodeResponse{Success: false, Message: &message, RelationshipSchemaNode: nil}, err
	}

	query = `
		CREATE CONSTRAINT relationship_schema_node_unique IF NOT EXISTS
		FOR (n:RELATIONSHIP_SCHEMA)
		REQUIRE (n._name, n._domain, n._fromTypeSchemaNodeId, n._toTypeSchemaNodeId) IS UNIQUE
	`

	fmt.Println(query)

	_, err = session.Run(ctx, query, nil)
	if err != nil {
		message := fmt.Sprintf("Unable to create relationship schema constraint. Error: %s", err.Error())
		return &model.RelationshipSchemaNodeResponse{Success: false, Message: &message, RelationshipSchemaNode: nil}, err
	}

	query = `
		CREATE (relationshipSchemaNode:RELATIONSHIP_SCHEMA {_id: $id, _domain: $domain, _name: $name, _originalName: $originalName, _type: "RELATIONSHIP SCHEMA", _fromTypeSchemaNodeId: $fromTypeSchemaNodeId, _toTypeSchemaNodeId: $toTypeSchemaNodeId})
		RETURN relationshipSchemaNode
	`

	fmt.Println(query)

	parameters := map[string]any{
		"id":                   id,
		"domain":               domain,
		"name":                 name,
		"originalName":         originalName,
		"fromTypeSchemaNodeId": fromTypeSchemaNodeId,
		"toTypeSchemaNodeId":   toTypeSchemaNodeId,
	}

	result, err := session.Run(ctx, query, parameters)
	if err != nil {
		return nil, err
	}

	if result.Next(ctx) {
		record := result.Record()
		relationshipSchemaNode, ok := record.Get("relationshipSchemaNode")
		if !ok {
			return nil, fmt.Errorf("failed to retrieve the relationshipSchemaNode")
		}
		neo4jRelationshipSchemaNode, ok := relationshipSchemaNode.(dbtype.Node)
		if !ok {
			return nil, fmt.Errorf("unexpected type for relationshipSchemaNode: %T", relationshipSchemaNode)
		}
		data := &model.RelationshipSchemaNode{
			ID:                   utils.PopString(neo4jRelationshipSchemaNode.Props, "_id"),
			Domain:               utils.PopString(neo4jRelationshipSchemaNode.Props, "_domain"),
			Name:                 utils.PopString(neo4jRelationshipSchemaNode.Props, "_name"),
			OriginalName:         utils.PopString(neo4jRelationshipSchemaNode.Props, "_originalName"),
			Type:                 utils.PopString(neo4jRelationshipSchemaNode.Props, "_type"),
			FromTypeSchemaNodeID: utils.PopString(neo4jRelationshipSchemaNode.Props, "_fromTypeSchemaNodeId"),
			ToTypeSchemaNodeID:   utils.PopString(neo4jRelationshipSchemaNode.Props, "_toTypeSchemaNodeId"),
			Properties:           utils.ExtractPropertiesFromNeo4jNode(neo4jRelationshipSchemaNode.Props),
			Labels:               neo4jRelationshipSchemaNode.Labels,
		}
		message := "Relationship schema created successfully"
		return &model.RelationshipSchemaNodeResponse{Success: true, Message: &message, RelationshipSchemaNode: data}, nil
	}
	message := "Unable to create relationship schema"
	return &model.RelationshipSchemaNodeResponse{Success: false, Message: &message, RelationshipSchemaNode: nil}, nil
}

func (db *Neo4jDatabase) RenameRelationshipSchemaNode(ctx context.Context, id string, newName string) (*model.RelationshipSchemaNodeResponse, error) {
	session := db.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	originalNewName := strings.TrimSpace(newName)
	newName = utils.RemoveSpacesAndHyphens(strings.ToUpper(newName))

	query := fmt.Sprintf(`
    OPTIONAL MATCH (relationshipSchemaNode:RELATIONSHIP_SCHEMA {_id: $id})
    WHERE relationshipSchemaNode IS NOT NULL
    WITH relationshipSchemaNode, relationshipSchemaNode._domain as domain, relationshipSchemaNode._name as existingName
    OPTIONAL MATCH (duplicateRelationshipSchemaNode:RELATIONSHIP_SCHEMA {_domain: domain, _name: $newName})
    WITH relationshipSchemaNode, domain, existingName, duplicateRelationshipSchemaNode
    WHERE duplicateRelationshipSchemaNode IS NULL
    SET relationshipSchemaNode._name = $newName,
        relationshipSchemaNode._originalName = $originalNewName
    WITH relationshipSchemaNode, domain, existingName
    OPTIONAL MATCH (fromObjectNode)-[oldRel {_name: existingName}]->(toObjectNode)
    WITH relationshipSchemaNode, fromObjectNode, toObjectNode, oldRel, existingName
    CALL apoc.do.when(
        fromObjectNode IS NOT NULL AND toObjectNode IS NOT NULL AND oldRel IS NOT NULL,
        'CREATE (fromObjectNode)-[newRel:%s]->(toObjectNode) SET newRel = properties(oldRel), newRel._name = $newName, newRel._originalName = $originalNewName RETURN newRel',
        'RETURN null as newRel',
        {fromObjectNode: fromObjectNode, toObjectNode: toObjectNode, oldRel: oldRel, newName: $newName, originalNewName: $originalNewName}
    ) YIELD value
    WITH relationshipSchemaNode, collect(oldRel) as oldRels, count(value.newRel) as updatedCount, existingName as previousName
    FOREACH (oldRel IN oldRels | DELETE oldRel)
    RETURN relationshipSchemaNode, updatedCount, previousName
`, newName)

	fmt.Println(query)

	parameters := map[string]any{
		"id":              id,
		"newName":         newName,
		"originalNewName": originalNewName,
	}

	result, err := session.Run(ctx, query, parameters)
	if err != nil {
		return nil, err
	}

	if result.Next(ctx) {
		record := result.Record()
		relationshipSchemaNode, ok := record.Get("relationshipSchemaNode")
		if !ok {
			return nil, fmt.Errorf("failed to retrieve the relationshipSchemaNode")
		}
		neo4jRelationshipSchemaNode, ok := relationshipSchemaNode.(dbtype.Node)
		if !ok {
			return nil, fmt.Errorf("unexpected type for relationshipSchemaNode: %T", relationshipSchemaNode)
		}
		updatedCount, ok := record.Get("updatedCount")
		if !ok {
			return nil, fmt.Errorf("failed to retrieve the updatedCount")
		}
		updatedCountInt, ok := updatedCount.(int64)
		if !ok {
			return nil, fmt.Errorf("failed to retrieve the updatedCount")
		}
		previousName, ok := record.Get("previousName")
		if !ok {
			return nil, fmt.Errorf("failed to retrieve the previousName")
		}
		previousNameString, ok := previousName.(string)
		if !ok {
			return nil, fmt.Errorf("failed to retrieve the previousName")
		}
		data := &model.RelationshipSchemaNode{
			ID:                   utils.PopString(neo4jRelationshipSchemaNode.Props, "_id"),
			Name:                 utils.PopString(neo4jRelationshipSchemaNode.Props, "_name"),
			OriginalName:         utils.PopString(neo4jRelationshipSchemaNode.Props, "_originalName"),
			Domain:               utils.PopString(neo4jRelationshipSchemaNode.Props, "_domain"),
			Type:                 utils.PopString(neo4jRelationshipSchemaNode.Props, "_type"),
			FromTypeSchemaNodeID: utils.PopString(neo4jRelationshipSchemaNode.Props, "_fromTypeSchemaNodeId"),
			ToTypeSchemaNodeID:   utils.PopString(neo4jRelationshipSchemaNode.Props, "_toTypeSchemaNodeId"),
			Properties:           utils.ExtractPropertiesFromNeo4jNode(neo4jRelationshipSchemaNode.Props),
			Labels:               neo4jRelationshipSchemaNode.Labels,
		}
		message := fmt.Sprintf("%s relationship schema node renamed to %s. %v object nodes updated successfully", previousNameString, newName, updatedCountInt)
		return &model.RelationshipSchemaNodeResponse{Success: true, Message: &message, RelationshipSchemaNode: data}, nil
	}
	if result.Err() != nil {
		return nil, result.Err()
	}
	message := "Unable to rename relationship schema node"
	return &model.RelationshipSchemaNodeResponse{Success: false, Message: &message, RelationshipSchemaNode: nil}, nil
}

func (db *Neo4jDatabase) UpdatePropertiesOnRelationshipSchemaNode(ctx context.Context, id string, properties []*model.PropertyInput) (*model.RelationshipSchemaNodeResponse, error) {
	session := db.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	if err := utils.CleanUpPropertyObjects(&properties); err != nil {
		message := fmt.Sprintf("Unable to update properties. Error: %s", err.Error())
		return &model.RelationshipSchemaNodeResponse{Success: false, Message: &message, RelationshipSchemaNode: nil}, nil
	}

	query := `MATCH (relationshipSchemaNode:RELATIONSHIP_SCHEMA {_id: $id}) SET `
	query = utils.CreatePropertiesQuery(query, properties, "relationshipSchemaNode")
	query = strings.TrimSuffix(query, ", ")
	query += ` RETURN relationshipSchemaNode`

	fmt.Println(query)

	parameters := map[string]any{
		"id": id,
	}

	result, err := session.Run(ctx, query, parameters)
	if err != nil {
		return nil, err
	}

	if result.Next(ctx) {
		record := result.Record()
		relationshipSchemaNode, ok := record.Get("relationshipSchemaNode")
		if !ok {
			return nil, fmt.Errorf("failed to retrieve the relationshipSchemaNode")
		}
		neo4jRelationshipSchemaNode, ok := relationshipSchemaNode.(dbtype.Node)
		if !ok {
			return nil, fmt.Errorf("unexpected type for relationshipSchemaNode: %T", relationshipSchemaNode)
		}
		data := &model.RelationshipSchemaNode{
			ID:                   utils.PopString(neo4jRelationshipSchemaNode.Props, "_id"),
			Name:                 utils.PopString(neo4jRelationshipSchemaNode.Props, "_name"),
			OriginalName:         utils.PopString(neo4jRelationshipSchemaNode.Props, "_originalName"),
			Domain:               utils.PopString(neo4jRelationshipSchemaNode.Props, "_domain"),
			Type:                 utils.PopString(neo4jRelationshipSchemaNode.Props, "_type"),
			FromTypeSchemaNodeID: utils.PopString(neo4jRelationshipSchemaNode.Props, "_fromTypeSchemaNodeId"),
			ToTypeSchemaNodeID:   utils.PopString(neo4jRelationshipSchemaNode.Props, "_toTypeSchemaNodeId"),
			Properties:           utils.ExtractPropertiesFromNeo4jNode(neo4jRelationshipSchemaNode.Props),
			Labels:               neo4jRelationshipSchemaNode.Labels,
		}
		message := "Relationship schema node properties updated successfully"
		return &model.RelationshipSchemaNodeResponse{Success: true, Message: &message, RelationshipSchemaNode: data}, nil
	}
	if result.Err() != nil {
		return nil, result.Err()
	}
	message := "Relationship schema node properties update failed"
	return &model.RelationshipSchemaNodeResponse{Success: false, Message: &message, RelationshipSchemaNode: nil}, nil
}

func (db *Neo4jDatabase) RenamePropertyOnRelationshipSchemaNode(ctx context.Context, id string, oldPropertyName string, newPropertyName string) (*model.RelationshipSchemaNodeResponse, error) {
	session := db.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	oldPropertyName = utils.RemoveSpacesAndLowerCase(oldPropertyName)
	newPropertyName = utils.RemoveSpacesAndLowerCase(newPropertyName)

	if strings.HasPrefix(oldPropertyName, "_") {
		return nil, fmt.Errorf("oldPropertyName cannot be %s", oldPropertyName)
	}

	if strings.HasPrefix(newPropertyName, "_") {
		return nil, fmt.Errorf("newPropertyName cannot be %s", newPropertyName)
	}

	query := fmt.Sprintf(`
    MATCH (relationshipSchemaNode:RELATIONSHIP_SCHEMA {_id: $id})
    WITH relationshipSchemaNode
    WHERE relationshipSchemaNode.%s IS NOT NULL
    WITH relationshipSchemaNode, relationshipSchemaNode.%s as oldValue
    SET relationshipSchemaNode.%s = oldValue
    REMOVE relationshipSchemaNode.%s
    WITH relationshipSchemaNode
    OPTIONAL MATCH ()-[rel {_name: relationshipSchemaNode._name}]->()
    WITH relationshipSchemaNode, collect(rel) as relationships, count(rel) as updatedCount
    FOREACH (r IN relationships |
        SET r.%s = r.%s
        REMOVE r.%s
    )
    RETURN relationshipSchemaNode, updatedCount
`, oldPropertyName, oldPropertyName, newPropertyName, oldPropertyName,
		newPropertyName, oldPropertyName, oldPropertyName)

	fmt.Println(query)

	parameters := map[string]any{
		"id": id,
	}

	result, err := session.Run(ctx, query, parameters)
	if err != nil {
		return nil, err
	}

	if result.Next(ctx) {
		record := result.Record()
		relationshipSchemaNode, ok := record.Get("relationshipSchemaNode")
		if !ok {
			return nil, fmt.Errorf("failed to retrieve the relationshipSchemaNode")
		}
		neo4jRelationshipSchemaNode, ok := relationshipSchemaNode.(dbtype.Node)
		if !ok {
			return nil, fmt.Errorf("unexpected type for relationshipSchemaNode: %T", relationshipSchemaNode)
		}
		updatedCount, ok := record.Get("updatedCount")
		if !ok {
			return nil, fmt.Errorf("failed to retrieve the updatedCount")
		}
		updatedCountInt, ok := updatedCount.(int64)
		if !ok {
			return nil, fmt.Errorf("failed to retrieve the updatedCount")
		}
		data := &model.RelationshipSchemaNode{
			ID:                   utils.PopString(neo4jRelationshipSchemaNode.Props, "_id"),
			Name:                 utils.PopString(neo4jRelationshipSchemaNode.Props, "_name"),
			Domain:               utils.PopString(neo4jRelationshipSchemaNode.Props, "_domain"),
			Type:                 utils.PopString(neo4jRelationshipSchemaNode.Props, "_type"),
			FromTypeSchemaNodeID: utils.PopString(neo4jRelationshipSchemaNode.Props, "_fromTypeSchemaNodeId"),
			ToTypeSchemaNodeID:   utils.PopString(neo4jRelationshipSchemaNode.Props, "_toTypeSchemaNodeId"),
			Properties:           utils.ExtractPropertiesFromNeo4jNode(neo4jRelationshipSchemaNode.Props),
			Labels:               neo4jRelationshipSchemaNode.Labels,
		}
		message := fmt.Sprintf("%s relationship schema node property renamed to %s. %v object nodes updated successfully", oldPropertyName, newPropertyName, updatedCountInt)
		return &model.RelationshipSchemaNodeResponse{Success: true, Message: &message, RelationshipSchemaNode: data}, nil
	}
	if result.Err() != nil {
		return nil, result.Err()
	}
	message := "Unable to rename relationship schema node property"
	return &model.RelationshipSchemaNodeResponse{Success: false, Message: &message, RelationshipSchemaNode: nil}, nil
}

func (db *Neo4jDatabase) RemovePropertiesFromRelationshipSchemaNode(ctx context.Context, id string, properties []string) (*model.RelationshipSchemaNodeResponse, error) {
	session := db.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	if err := utils.CleanUpPropertyKeys(&properties); err != nil {
		message := fmt.Sprintf("Unable to remove properties. Error: %s", err.Error())
		return &model.RelationshipSchemaNodeResponse{Success: false, Message: &message, RelationshipSchemaNode: nil}, nil
	}

	query := `MATCH (relationshipSchemaNode:RELATIONSHIP_SCHEMA {_id: $id}) `
	query += `OPTIONAL MATCH ()-[rel {_name: relationshipSchemaNode._name}]->() `
	query += `WITH relationshipSchemaNode, collect(rel) as relationships, count(rel) as updatedCount `
	query += `SET `
	query = utils.RemovePropertiesQuery(query, properties, "relationshipSchemaNode")
	query = strings.TrimSuffix(query, ", ")
	query += ` WITH relationshipSchemaNode, relationships, updatedCount `
	query += `FOREACH (rel IN relationships | SET `
	query = utils.RemovePropertiesQuery(query, properties, "rel")
	query = strings.TrimSuffix(query, ", ")
	query += `) RETURN relationshipSchemaNode, updatedCount`

	fmt.Println(query)

	parameters := map[string]any{
		"id": id,
	}

	result, err := session.Run(ctx, query, parameters)
	if err != nil {
		return nil, err
	}

	if result.Next(ctx) {
		record := result.Record()
		relationshipSchemaNode, ok := record.Get("relationshipSchemaNode")
		if !ok {
			return nil, fmt.Errorf("failed to retrieve the relationshipSchemaNode")
		}
		neo4jRelationshipSchemaNode, ok := relationshipSchemaNode.(dbtype.Node)
		if !ok {
			return nil, fmt.Errorf("unexpected type for relationshipSchemaNode: %T", relationshipSchemaNode)
		}
		updatedCount, ok := record.Get("updatedCount")
		if !ok {
			return nil, fmt.Errorf("failed to retrieve the updatedCount")
		}
		updatedCountInt, ok := updatedCount.(int64)
		if !ok {
			return nil, fmt.Errorf("failed to retrieve the updatedCount")
		}
		data := &model.RelationshipSchemaNode{
			ID:                   utils.PopString(neo4jRelationshipSchemaNode.Props, "_id"),
			Name:                 utils.PopString(neo4jRelationshipSchemaNode.Props, "_name"),
			OriginalName:         utils.PopString(neo4jRelationshipSchemaNode.Props, "_originalName"),
			Domain:               utils.PopString(neo4jRelationshipSchemaNode.Props, "_domain"),
			Type:                 utils.PopString(neo4jRelationshipSchemaNode.Props, "_type"),
			FromTypeSchemaNodeID: utils.PopString(neo4jRelationshipSchemaNode.Props, "_fromTypeSchemaNodeId"),
			ToTypeSchemaNodeID:   utils.PopString(neo4jRelationshipSchemaNode.Props, "_toTypeSchemaNodeId"),
			Properties:           utils.ExtractPropertiesFromNeo4jNode(neo4jRelationshipSchemaNode.Props),
			Labels:               neo4jRelationshipSchemaNode.Labels,
		}
		message := fmt.Sprintf("%v relationship schema node properties removed successfully. %v object nodes updated successfully", len(properties), updatedCountInt)
		return &model.RelationshipSchemaNodeResponse{Success: true, Message: &message, RelationshipSchemaNode: data}, nil
	}
	if result.Err() != nil {
		return nil, result.Err()
	}
	message := "Unable to remove properties from relationship schema"
	return &model.RelationshipSchemaNodeResponse{Success: false, Message: &message, RelationshipSchemaNode: nil}, nil
}

func (db *Neo4jDatabase) DeleteRelationshipSchemaNode(ctx context.Context, id string) (*model.RelationshipSchemaNodeResponse, error) {
	session := db.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	query := `
    MATCH (relationshipSchemaNode:RELATIONSHIP_SCHEMA {_id: $id})
    OPTIONAL MATCH ()-[rel {_name: relationshipSchemaNode._name}]->()
    WITH relationshipSchemaNode, collect(rel) as relationships, properties(relationshipSchemaNode) as relationshipSchemaNodeProperties, labels(relationshipSchemaNode) as relationshipSchemaNodeLabels
    WITH relationshipSchemaNode, relationships, relationshipSchemaNodeProperties, relationshipSchemaNodeLabels, size(relationships) as relationshipsCount
    FOREACH (rel IN relationships | DELETE rel)
    DELETE relationshipSchemaNode
    WITH relationshipSchemaNodeProperties, relationshipSchemaNodeLabels, relationshipsCount
    RETURN relationshipsCount, relationshipSchemaNodeProperties, relationshipSchemaNodeLabels
`

	fmt.Println(query)

	parameters := map[string]any{
		"id": id,
	}

	result, err := session.Run(ctx, query, parameters)
	if err != nil {
		return nil, err
	}
	if result.Next(ctx) {
		record := result.Record()
		relationshipsCount, ok := record.Get("relationshipsCount")
		if !ok {
			return nil, fmt.Errorf("failed to retrieve the relationshipsCount")
		}
		relationshipsCountInt, ok := relationshipsCount.(int64)
		if !ok {
			return nil, fmt.Errorf("failed to retrieve the relationshipsCount")
		}
		relationshipSchemaNodeProperties, ok := record.Get("relationshipSchemaNodeProperties")
		if !ok {
			return nil, fmt.Errorf("failed to retrieve the relationshipSchemaNodeProperties")
		}
		relationshipSchemaNodePropertiesMap, ok := relationshipSchemaNodeProperties.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("failed to retrieve the relationshipSchemaNodeProperties")
		}
		relationshipSchemaNodeLabels, ok := record.Get("relationshipSchemaNodeLabels")
		if !ok {
			return nil, fmt.Errorf("failed to retrieve the relationshipSchemaNodeLabels")
		}
		relationshipSchemaNodeLabelsArray, ok := relationshipSchemaNodeLabels.([]interface{})
		if !ok {
			return nil, fmt.Errorf("failed to retrieve the relationshipSchemaNodeLabels")
		}
		labels := []string{}
		for _, label := range relationshipSchemaNodeLabelsArray {
			labels = append(labels, label.(string))
		}
		data := &model.RelationshipSchemaNode{
			ID:                   utils.PopString(relationshipSchemaNodePropertiesMap, "_id"),
			Name:                 utils.PopString(relationshipSchemaNodePropertiesMap, "_name"),
			OriginalName:         utils.PopString(relationshipSchemaNodePropertiesMap, "_originalName"),
			Domain:               utils.PopString(relationshipSchemaNodePropertiesMap, "_domain"),
			Type:                 utils.PopString(relationshipSchemaNodePropertiesMap, "_type"),
			FromTypeSchemaNodeID: utils.PopString(relationshipSchemaNodePropertiesMap, "_fromTypeSchemaNodeId"),
			ToTypeSchemaNodeID:   utils.PopString(relationshipSchemaNodePropertiesMap, "_toTypeSchemaNodeId"),
			Properties:           utils.ExtractPropertiesFromNeo4jNode(relationshipSchemaNodePropertiesMap),
			Labels:               labels,
		}
		message := fmt.Sprintf("Relationship schema node %s deleted successfully. %v relationships deleted successfully", data.Name, relationshipsCountInt)
		return &model.RelationshipSchemaNodeResponse{Success: true, Message: &message, RelationshipSchemaNode: data}, nil
	}
	if result.Err() != nil {
		return nil, result.Err()
	}
	message := "Unable to delete relationship schema node"
	return &model.RelationshipSchemaNodeResponse{Success: false, Message: &message, RelationshipSchemaNode: nil}, nil

}

func (db *Neo4jDatabase) GetTypeSchemaNodeRelationships(ctx context.Context, id string) (*model.RelationshipSchemaNodesResponse, error) {
	return nil, nil
}
