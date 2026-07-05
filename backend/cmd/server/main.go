package main

import (
	"log"
	"net/http"
	"os"
	"social-network/backend/internal/middleware"
	"social-network/backend/internal/router"
	"social-network/backend/pkg/db/sqlite"
)

func main() {
	log.Println("[BOOTSTRAP] Commencing system initialization protocol...")

	dbPath := "./socialnetwork.db"

	// 1. Data Bedrock Layer Connection
	database, err := sqlite.Connect(dbPath)
	if err != nil {
		log.Fatalf("[CRITICAL] Database setup failed: %v", err)
	}
	defer database.Db.Close()

	if err := database.RunMigrations(); err != nil {
		log.Fatalf("[CRITICAL] Migration execution failed: %v", err)
	}

	// 2. Build Routing Topology
	mux := router.NewRouter(database.Db)

	// 3. Chain Edge Security & Visibility Framework (Starting basic)
	handlerPipeline := middleware.Logger(middleware.RateLimit(middleware.CORS(mux)))

	// 4. Bind and Listen on Network Port Boundary
	serverAddr := ":8080"
	log.Printf("[SERVER] Edge infrastructure online. Bound to host port %s", serverAddr)
	
	if err := http.ListenAndServe(serverAddr, handlerPipeline); err != nil {
		log.Fatalf("[CRITICAL] Network server collapsed: %v", err)
		os.Exit(1)
	}
}