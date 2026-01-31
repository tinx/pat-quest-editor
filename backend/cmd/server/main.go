package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/tinx/pat-quest-editor/backend/internal/adapters/filesystem"
	httpAdapter "github.com/tinx/pat-quest-editor/backend/internal/adapters/http"
	"github.com/tinx/pat-quest-editor/backend/internal/adapters/storage"
	"github.com/tinx/pat-quest-editor/backend/internal/app"
)

func main() {
	// Command line flags
	addr := flag.String("addr", ":8080", "HTTP server address")
	questsDir := flag.String("quests", "../quests", "Path to quests directory")
	dataDir := flag.String("data", "../data", "Path to reference data directory")
	dbPath := flag.String("db", "editor.db", "Path to SQLite database")
	staticDir := flag.String("static", "../frontend/dist", "Path to frontend static files")
	flag.Parse()

	// Resolve paths relative to working directory
	questsPath, _ := filepath.Abs(*questsDir)
	dataPath, _ := filepath.Abs(*dataDir)
	dbPathAbs, _ := filepath.Abs(*dbPath)
	staticPath, _ := filepath.Abs(*staticDir)

	// Initialize repositories
	questRepo := filesystem.NewQuestFileRepository(questsPath)
	refDataRepo := filesystem.NewReferenceDataFileRepository(dataPath)

	metadataRepo, err := storage.NewSQLiteMetadataRepository(dbPathAbs)
	if err != nil {
		log.Fatalf("Failed to initialize metadata repository: %v", err)
	}
	defer metadataRepo.Close()

	// Initialize services
	validator := app.NewQuestValidatorService(refDataRepo)

	// Initialize HTTP handler
	handler := httpAdapter.NewHandler(questRepo, refDataRepo, metadataRepo, validator)

	// Set up routes
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	// Serve static files for frontend
	if _, err := os.Stat(staticPath); err == nil {
		fileServer := http.FileServer(http.Dir(staticPath))
		mux.Handle("/", fileServer)
		log.Printf("Serving static files from %s", staticPath)
	} else {
		log.Printf("Static directory not found: %s (frontend will not be served)", staticPath)
	}

	// Start server
	log.Printf("Starting server on %s", *addr)
	log.Printf("Quest files: %s", questsPath)
	log.Printf("Reference data: %s", dataPath)
	log.Printf("Database: %s", dbPathAbs)

	if err := http.ListenAndServe(*addr, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
