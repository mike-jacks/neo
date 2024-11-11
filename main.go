package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/joho/godotenv"
	"github.com/mike-jacks/neo/db"
	"github.com/mike-jacks/neo/generated"
	"github.com/mike-jacks/neo/resolver"
	"github.com/rs/cors"
)

func setupGraphQLServer(db db.Database) *handler.Server {
	resolver := resolver.NewResolver(db)
	schema := generated.NewExecutableSchema(generated.Config{Resolvers: resolver})
	server := handler.NewDefaultServer(schema)

	// Add WebSocket transport without InitFunc
	server.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
	})

	return server
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found")
	}

	driver, err := db.SetupNeo4jDriver()
	if err != nil {
		log.Fatal(err)
	}
	defer driver.Close(context.Background())

	neo4jdb := &db.Neo4jDatabase{Driver: driver}

	srv := setupGraphQLServer(neo4jdb)

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	http.Handle("/graphql", corsHandler.Handler(playground.Handler("GraphQL Playground", "/query")))
	http.Handle("/query", corsHandler.Handler(srv))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	url := os.Getenv("URL")
	if url == "" {
		url = "http://localhost" + ":" + port
	}

	log.Printf("Connect to %s/graphql for GraphQL Playground", url)
	log.Printf("Connect to %s/query for GraphQL API", url)
	log.Printf("Connect to https://console.neo4j.io for Neo4j Browser Console")
	log.Fatal(http.ListenAndServe(":"+port, nil))

}
