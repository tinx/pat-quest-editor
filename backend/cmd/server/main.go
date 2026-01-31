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

// corsMiddleware adds CORS headers for development mode
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	// Command line flags
	addr := flag.String("addr", ":8080", "HTTP server address")
	questsDir := flag.String("quests", "../quests", "Path to quests directory")
	dataDir := flag.String("data", "../data", "Path to reference data directory")
	dbPath := flag.String("db", "editor.db", "Path to SQLite database")
	staticDir := flag.String("static", "../frontend/dist", "Path to frontend static files")
	devMode := flag.Bool("dev", false, "Enable development mode (CORS headers)")
	flag.Parse()

	// Resolve paths relative to working directory
	questsPath, err := filepath.Abs(*questsDir)
	if err != nil {
		log.Printf("Warning: failed to resolve quests path: %v", err)
		questsPath = *questsDir
	}
	dataPath, err := filepath.Abs(*dataDir)
	if err != nil {
		log.Printf("Warning: failed to resolve data path: %v", err)
		dataPath = *dataDir
	}
	dbPathAbs, err := filepath.Abs(*dbPath)
	if err != nil {
		log.Printf("Warning: failed to resolve database path: %v", err)
		dbPathAbs = *dbPath
	}
	staticPath, err := filepath.Abs(*staticDir)
	if err != nil {
		log.Printf("Warning: failed to resolve static path: %v", err)
		staticPath = *staticDir
	}

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

	// Wrap with CORS middleware in dev mode
	var finalHandler http.Handler = mux
	if *devMode {
		log.Printf("Development mode enabled (CORS headers active)")
		finalHandler = corsMiddleware(mux)
	}

	// Start server
	log.Printf("Starting server on %s", *addr)
	log.Printf("Quest files: %s", questsPath)
	log.Printf("Reference data: %s", dataPath)
	log.Printf("Database: %s", dbPathAbs)

	if err := http.ListenAndServe(*addr, finalHandler); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
