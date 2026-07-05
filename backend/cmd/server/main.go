package main

import (
	"log"
	"os"
	"social-network/backend/pkg/db/sqlite"
)

func main() {
	log.Println("[BOOTSTRAP] Initializing database engine validation scaffold...")

	dbPath := "./socialnetwork.db"

	// 1. Establish data link layer
	database, err := sqlite.Connect(dbPath)
	if err != nil {
		log.Fatalf("[CRITICAL] Database connection initialization aborted: %v", err)
	}
	defer database.Db.Close()

	// 2. Schema synchronization run
	if err := database.RunMigrations(); err != nil {
		log.Fatalf("[CRITICAL] Schema migration execution aborted: %v", err)
	}

	log.Println("[SUCCESS] Phase 1 verification complete. Bedrock layers structurally secure.")
	os.Exit(0)
}