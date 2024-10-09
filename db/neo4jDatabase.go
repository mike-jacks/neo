package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/mike-jacks/neo/model"
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
	for i, label := range labels {
		labels[i] = strings.Trim(strings.ToUpper(label), " ")
	}
	for i := range properties {
		properties[i].Key = strings.Trim(strings.ToLower(properties[i].Key), " ")
	}

	query := fmt.Sprintf(`
		CREATE CONSTRAINT IF NOT EXISTS
		FOR (n:_%v)
		REQUIRE (n.name, n.type, n.domain) IS NODE KEY
		`, domain)
	_, err := session.Run(ctx, query, nil)
	if err != nil {
		return nil, err
	}

	query = fmt.Sprintf(`
		CREATE CONSTRAINT IF NOT EXISTS
		FOR (n:%v)
		REQUIRE (n.name, n.type, n.domain) IS NODE KEY
		`, typeArg)

	_, err = session.Run(ctx, query, nil)
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

	query = fmt.Sprintf("CREATE (o:_%v:%v", domain, typeArg)
	for _, label := range labels {
		query += fmt.Sprintf(":%v", label)
	}
	query += " {name: $name, type: $typeArg, domain: $domain "
	for _, property := range properties {
		if property.Type == "STRING" {
			query += fmt.Sprintf(", %v: \"%v\"", property.Key, property.Value)
		} else {
			query += fmt.Sprintf(", %v: %v", property.Key, property.Value)
		}
	}
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

		data := map[string]interface{}{
			"labels":     neo4jNode.Labels,
			"properties": nodeProperties,
		}
		message := "Object node created successfully"
		return &model.Response{Success: true, Message: &message, Data: data}, nil
	}

	return nil, err
}
