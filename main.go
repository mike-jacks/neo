package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/joho/godotenv"
	"github.com/mike-jacks/neo/graph/generated"
	"github.com/mike-jacks/neo/graph/resolver"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func setupGraphQLServer(driver neo4j.DriverWithContext) *handler.Server {
	resolver := resolver.NewResolver(driver)
	schema := generated.NewExecutableSchema(generated.Config{Resolvers: resolver})
	server := handler.NewDefaultServer(schema)
	return server
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	driver, err := setupNeo4jDriver()
	if err != nil {
		log.Fatal(err)
	}
	defer driver.Close(context.Background())

	srv := setupGraphQLServer(driver)

	http.Handle("/playground", playground.Handler("GraphQL Playground", "/query"))
	http.Handle("/query", srv)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Connect to http://localhost:%s/graphiql for GraphiQL", port)
	log.Printf("Connect to http://localhost:%s/playground for GraphQL Playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func setupNeo4jDriver() (neo4j.DriverWithContext, error) {
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

// ... existing code ...
