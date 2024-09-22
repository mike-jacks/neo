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

	// GraphiQL handler
	graphiqlHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(`
			<!DOCTYPE html>
			<html>
			<head>
				<title>GraphiQL</title>
				<link href="https://unpkg.com/graphiql/graphiql.min.css" rel="stylesheet" />
			</head>
			<body style="margin: 0;">
				<div id="graphiql" style="height: 100vh;"></div>
				<script
					crossorigin
					src="https://unpkg.com/react/umd/react.production.min.js"
				></script>
				<script
					crossorigin
					src="https://unpkg.com/react-dom/umd/react-dom.production.min.js"
				></script>
				<script
					crossorigin
					src="https://unpkg.com/graphiql/graphiql.min.js"
				></script>
				<script>
					const graphQLFetcher = graphQLParams =>
						fetch('/query', {
							method: 'post',
							headers: {
								'Content-Type': 'application/json',
							},
							body: JSON.stringify(graphQLParams),
						}).then(response => response.json());

					ReactDOM.render(
						React.createElement(GraphiQL, { fetcher: graphQLFetcher }),
						document.getElementById('graphiql'),
					);
				</script>
			</body>
			</html>
		`))
	})

	http.Handle("/graphiql", graphiqlHandler)
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
