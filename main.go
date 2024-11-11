package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/mike-jacks/neo/db"
	"github.com/mike-jacks/neo/generated"
	"github.com/mike-jacks/neo/resolver"
	"github.com/rs/cors"
)

func setupGraphQLServer(db db.Database) *handler.Server {
	resolver := resolver.NewResolver(db)
	schema := generated.NewExecutableSchema(generated.Config{Resolvers: resolver})
	server := handler.New(schema)

	allowedOrigins := map[string]bool{
		"https://neo-frontend-v2.vercel.app": true,
		"http://localhost:5173":              true,
		"http://localhost":                   true,
		"http://172.0.0.1":                   true,
	}

	// Add WebSocket transport without InitFunc
	server.AddTransport(transport.Websocket{
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				origin := strings.TrimRight(r.Header.Get("Origin"), "/")
				log.Printf("WebSocket Origin: '%s'", origin)

				// Debug allowed origins
				for allowed := range allowedOrigins {
					log.Printf("Comparing with allowed origin: '%s'", allowed)
				}

				// Compare directly with map
				if !allowedOrigins[origin] {
					log.Printf("WebSocket connection rejected from origin: '%s'", origin)
					return false
				}

				log.Printf("WebSocket connection accepted from origin: '%s'", origin)
				return true
			},
		},
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
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{"*"},
	})

	loggingHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request received: Method: %s, Path: %s", r.Method, r.URL.Path)
		if websocket.IsWebSocketUpgrade(r) {
			log.Printf("WebSocket Upgrade Detected. Origin: %s", r.Header.Get("Origin"))
			srv.ServeHTTP(w, r)
			return
		} else {
			log.Printf("Non-WebSocket Request Detected. Origin: %s", r.Header.Get("Origin"))
			corsHandler.Handler(srv).ServeHTTP(w, r)
		}
	})

	http.Handle("/graphql", corsHandler.Handler(playground.Handler("GraphQL Playground", "/query")))
	http.Handle("/query", loggingHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	url := os.Getenv("URL")

	var websocketUrl string
	if url == "" {
		url = "http://localhost" + ":" + port
		websocketUrl = "ws://localhost" + ":" + port
	} else {
		websocketUrl = strings.Replace(url, "https://", "wss://", 1)
	}

	log.Printf("Connect to %s/graphql for GraphQL Playground", url)
	log.Printf("GraphQL WebSocket endpoint: %s/query", websocketUrl)
	log.Printf("Connect to %s/query for GraphQL API", url)
	log.Printf("Connect to https://console.neo4j.io for Neo4j Browser Console")
	log.Fatal(http.ListenAndServe(":"+port, nil))

}
