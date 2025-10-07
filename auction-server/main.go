// main.go
package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/micahli/fl-auction/auction-server/graph"
	"github.com/micahli/fl-auction/auction-server/internal/service"
	"github.com/micahli/fl-auction/auction-server/internal/store"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/websocket"
	"github.com/rs/cors"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// Initialize the data store
	auctionStore := store.NewAuctionStore()

	// Initialize the service layer
	auctionService := service.NewAuctionService(auctionStore)

	// Create the GraphQL resolver
	resolver := graph.NewResolver(auctionService, auctionStore)

	// Create the GraphQL server with the generated schema
	srv := handler.New(graph.NewExecutableSchema(graph.Config{
		Resolvers: resolver,
	}))

	// Configure HTTP transports
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	// Configure WebSocket transport for subscriptions
	srv.AddTransport(&transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// In production, validate origin properly
				return true
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	})

	// Add GraphQL extensions
	srv.Use(extension.Introspection{})

	// Configure CORS for frontend access
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:5173", "http://localhost:8080"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		Debug:            false,
	})

	// Setup HTTP routes
	http.Handle("/", playground.Handler("GraphQL Playground", "/query"))
	http.Handle("/query", corsHandler.Handler(srv))

	// Start the server
	log.Printf("🚀 Server starting on http://localhost:%s", port)
	log.Printf("📊 GraphQL Playground: http://localhost:%s/", port)
	log.Printf("🔌 GraphQL Endpoint: http://localhost:%s/query", port)
	log.Printf("⚡ WebSocket Endpoint: ws://localhost:%s/query", port)
	log.Printf("\n📝 Try these queries in the playground:\n")
	log.Printf("   - Create auction: mutation { createAuction(startingBid: 100, duration: 30, extendedBidding: true) { id status } }\n")
	log.Printf("   - Place bid: mutation { placeBid(userId: \"user123\", amount: 150) { id amount } }\n")
	log.Printf("   - Subscribe: subscription { auctionEvents { type auction { currentBid currentWinner timeRemaining } } }\n")

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
