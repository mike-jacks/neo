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

// CreateSchemaNode creates a new schema node
func (db *Neo4jDatabase) CreateSchemaNode(ctx context.Context, sourceSchemaNodeName *string, createSchemaNodeInput model.CreateSchemaNodeInput) (*model.SchemaNode, error) {
	session := db.GetDriver().NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	// Ensure uniqueness constraint on SchemaNode name
	constraintQuery := "CREATE CONSTRAINT unique_schema_node_name_domain IF NOT EXISTS FOR (n:SchemaNode) REQUIRE (n.name,n.domain) IS UNIQUE"
	if err := utils.CreateConstraint(ctx, db.GetDriver(), []*string{&constraintQuery}); err != nil {
		return nil, fmt.Errorf("failed to create constraint: %w", err)
	}

	schemaNode := &model.SchemaNode{
		Name:   strings.TrimSpace(strings.ToUpper(createSchemaNodeInput.Name)),
		Domain: strings.TrimSpace(strings.ToUpper(createSchemaNodeInput.Domain)),
	}

	var query string
	parameters := map[string]interface{}{
		"name":   schemaNode.Name,
		"domain": schemaNode.Domain,
	}

	if sourceSchemaNodeName == nil {
		query = `
			CREATE (n:SchemaNode {name: $name, domain: $domain})
			RETURN n;
		`
	} else {
		query = `
			MATCH (sourceNode:SchemaNode {name: $sourceSchemaNodeName})
			CREATE (newNode:SchemaNode {name: $name, domain: $domain})
			CREATE (sourceNode)<-[:BELONGS_TO]-(newNode)
			RETURN newNode;
		`
		parameters["sourceSchemaNodeName"] = strings.TrimSpace(strings.ToUpper(*sourceSchemaNodeName))
	}

	var newSchemaNode *model.SchemaNode

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		result, err := tx.Run(ctx, query, parameters)
		if err != nil {
			return nil, err
		}

		record, err := result.Single(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to create node: %w", err)
		}

		node, ok := record.Values[0].(neo4j.Node)
		if !ok {
			return nil, fmt.Errorf("unexpected result type: %T", record.Values[0])
		}

		newSchemaNode = &model.SchemaNode{
			Name:   node.Props["name"].(string),
			Domain: node.Props["domain"].(string),
		}

		return newSchemaNode, nil
	})

	if err != nil {
		return nil, err
	}

	return newSchemaNode, nil
}

func (db *Neo4jDatabase) InsertSchemaNode(ctx context.Context, domain string, parentName string, childName string) (*model.SchemaNode, error) {
	return nil, fmt.Errorf("not implemented")
}

func (db *Neo4jDatabase) UpdateSchemaNode(ctx context.Context, domain string, name string, updateSchemaNodeInput model.UpdateSchemaNodeInput) ([]*model.SchemaNode, error) {
	domain = strings.TrimSpace(strings.ToUpper(domain))
	name = strings.TrimSpace(strings.ToUpper(name))

	session := db.GetDriver().NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	var nodes []*model.SchemaNode

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		query := `
			MATCH (root:SchemaNode {name: $origName, domain: $origDomain})
		`
		setClause := []string{}
		parameters := map[string]interface{}{
			"origName":   name,
			"origDomain": domain,
		}

		if updateSchemaNodeInput.Name != nil {
			setClause = append(setClause, "root.name = $newName")
			parameters["newName"] = strings.TrimSpace(strings.ToUpper(*updateSchemaNodeInput.Name))
		}

		if updateSchemaNodeInput.Domain != nil {
			setClause = append(setClause, "root.domain = $newDomain")
			parameters["newDomain"] = strings.TrimSpace(strings.ToUpper(*updateSchemaNodeInput.Domain))
		}

		if len(setClause) > 0 {
			query += "SET " + strings.Join(setClause, ", ")
		} else {
			return nil, fmt.Errorf("no fields to update")
		}

		query += `
			WITH root
			OPTIONAL MATCH (root)-[:HAS_PROPERTY]->(property:SchemaProperty)
			OPTIONAL MATCH (root)-[:HAS_LABEL]->(label:SchemaLabel)
			OPTIONAL MATCH (root)-[:HAS_RELATIONSHIP]->(relationship:SchemaRelationship)
			OPTIONAL MATCH (descendant:SchemaNode)-[:BELONGS_TO*]->(root)
			OPTIONAL MATCH (descendant)-[:BELONGS_TO]->(parent:SchemaNode)
			OPTIONAL MATCH (descendant)-[:HAS_PROPERTY]->(descendantProperty:SchemaProperty)
			OPTIONAL MATCH (descendant)-[:HAS_LABEL]->(descendantLabel:SchemaLabel)
			OPTIONAL MATCH (descendant)-[:HAS_RELATIONSHIP]->(descendantRelationship:SchemaRelationship)
			RETURN root, collect(property) as properties, collect(label) as labels, collect(relationship) as relationships, collect(descendant) as descendants, collect(descendantProperty) as descendantProperties, collect(descendantLabel) as descendantLabels, collect(descendantRelationship) as descendantRelationships
		`

		result, err := tx.Run(ctx, query, parameters)
		if err != nil {
			return nil, err
		}

		if result.Next(ctx) {
			record := result.Record()
			root, _ := record.Get("root")
			properties, _ := record.Get("properties")
			labels, _ := record.Get("labels")
			relationships, _ := record.Get("relationships")
			descendants, _ := record.Get("descendants")
			descendantProperties, _ := record.Get("descendantProperties")
			descendantLabels, _ := record.Get("descendantLabels")
			descendantRelationships, _ := record.Get("descendantRelationships")

			var parentName *string
			if val, ok := root.(neo4j.Node).Props["parentName"]; ok && val != nil {
				strVal := val.(string)
				parentName = &strVal
			}

			nodes = append(nodes, &model.SchemaNode{
				Name:          root.(neo4j.Node).Props["name"].(string),
				Domain:        root.(neo4j.Node).Props["domain"].(string),
				ParentName:    parentName,
				Properties:    extractProperties(properties),
				Labels:        extractLabels(labels),
				Relationships: extractRelationships(relationships),
			})

			for _, descendant := range descendants.([]interface{}) {
				descendantNode := descendant.(neo4j.Node)
				var parentName *string
				if val, ok := descendantNode.Props["parentName"]; ok && val != nil {
					strVal := val.(string)
					parentName = &strVal
				}

				node := &model.SchemaNode{
					Name:          descendantNode.Props["name"].(string),
					Domain:        descendantNode.Props["domain"].(string),
					ParentName:    parentName,
					Properties:    extractProperties(descendantProperties),
					Labels:        extractLabels(descendantLabels),
					Relationships: extractRelationships(descendantRelationships),
				}
				nodes = append(nodes, node)
			}
		}

		return nil, result.Err()
	})

	if err != nil {
		return nil, err
	}

	return nodes, nil
}

// DeleteSchemaNode is the resolver for the deleteSchemaNode field.
func (db *Neo4jDatabase) DeleteSchemaNode(ctx context.Context, domain string, name string) (bool, error) {
	domain = strings.TrimSpace(strings.ToUpper(domain))
	name = strings.TrimSpace(strings.ToUpper(name))

	session := db.GetDriver().NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	result, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		query := `
			MATCH (nodeToDelete:SchemaNode {name: $name, domain: $domain})
			OPTIONAL MATCH (nodeToDelete)<-[:BELONGS_TO]-(childNode:SchemaNode)
			OPTIONAL MATCH (nodeToDelete)-[:BELONGS_TO]->(parentNode:SchemaNode)
			WITH nodeToDelete, childNode, parentNode, COUNT(nodeToDelete) > 0 AS nodeExists
			FOREACH (_ IN CASE WHEN childNode IS NOT NULL AND parentNode IS NOT NULL THEN [1] ELSE [] END |
				CREATE (childNode)-[:BELONGS_TO]->(parentNode)
			)
			WITH nodeToDelete, nodeExists
			DETACH DELETE nodeToDelete
			RETURN nodeExists
		`
		parameters := map[string]interface{}{
			"name":   name,
			"domain": domain,
		}
		result, err := tx.Run(ctx, query, parameters)
		if err != nil {
			return false, err
		}

		record, err := result.Single(ctx)
		if err != nil {
			if err.Error() == "Result contains no more records" {
				return false, fmt.Errorf("node with name '%s' and domain '%s' does not exist", name, domain)
			}
			return false, err
		}
		nodeExists, ok := record.Get("nodeExists")
		if !ok {
			return false, fmt.Errorf("node with name '%s' and domain '%s' does not exist", name, domain)
		}
		return nodeExists.(bool), nil
	})

	if err != nil {
		return false, err
	}

	nodeExists := result.(bool)
	if !nodeExists {
		return false, fmt.Errorf("node with name '%s' and domain '%s' does not exist", name, domain)
	}

	return nodeExists, nil
}

func extractProperties(value interface{}) []*model.SchemaProperty {
	properties := []*model.SchemaProperty{}
	if value == nil {
		return properties
	}

	for _, v := range value.([]interface{}) {
		if property, ok := v.(neo4j.Node); ok {
			properties = append(properties, &model.SchemaProperty{
				Name:                 property.Props["name"].(string),
				Type:                 property.Props["type"].(string),
				Domain:               property.Props["domain"].(string),
				ParentSchemaNodeName: property.Props["parentSchemaNodeName"].(string),
			})
		}
	}

	return properties
}

func extractLabels(value interface{}) []*model.SchemaLabel {
	labels := []*model.SchemaLabel{}
	if value == nil {
		return labels
	}

	for _, v := range value.([]interface{}) {
		if label, ok := v.(neo4j.Node); ok {
			labels = append(labels, &model.SchemaLabel{
				Name:                 label.Props["name"].(string),
				Domain:               label.Props["domain"].(string),
				ParentSchemaNodeName: label.Props["parentSchemaNodeName"].(string),
			})
		}
	}

	return labels
}

func extractRelationships(value interface{}) []*model.SchemaRelationship {
	relationships := []*model.SchemaRelationship{}
	if value == nil {
		return relationships
	}

	for _, v := range value.([]interface{}) {
		if relationship, ok := v.(neo4j.Node); ok {
			relationships = append(relationships, &model.SchemaRelationship{
				Name:                 relationship.Props["name"].(string),
				Domain:               relationship.Props["domain"].(string),
				TargetSchemaNodeName: relationship.Props["targetSchemaNodeName"].(string),
				ParentSchemaNodeName: relationship.Props["parentSchemaNodeName"].(string),
			})
		}
	}

	return relationships
}
