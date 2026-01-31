package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/tinx/pat-quest-editor/backend/internal/adapters/filesystem"
	httpAdapter "github.com/tinx/pat-quest-editor/backend/internal/adapters/http"
	"github.com/tinx/pat-quest-editor/backend/internal/adapters/storage"
	"github.com/tinx/pat-quest-editor/backend/internal/app"
)

// allowedDevOrigins contains origins allowed in development mode
var allowedDevOrigins = []string{
	"http://localhost:5173",  // Vite default
	"http://localhost:3000",  // Common React dev port
	"http://127.0.0.1:5173",
	"http://127.0.0.1:3000",
}

// corsMiddleware adds CORS headers for development mode (localhost only)
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		// Only allow specific localhost origins
		allowed := false
		for _, o := range allowedDevOrigins {
			if origin == o {
				allowed = true
				break
			}
		}

		if allowed {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.Header().Set("Vary", "Origin")
		}

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// validatePathSafety ensures a path doesn't escape expected boundaries using path traversal
func validatePathSafety(basePath, targetPath string) error {
	absBase, err := filepath.Abs(basePath)
	if err != nil {
		return fmt.Errorf("failed to resolve base path: %w", err)
	}
	absTarget, err := filepath.Abs(targetPath)
	if err != nil {
		return fmt.Errorf("failed to resolve target path: %w", err)
	}

	// Check for path traversal attempts
	if strings.Contains(targetPath, "..") {
		return fmt.Errorf("path contains invalid traversal sequence")
	}

	// Ensure target is within base directory
	if !strings.HasPrefix(absTarget, absBase+string(filepath.Separator)) && absTarget != absBase {
		return fmt.Errorf("path escapes base directory")
	}
	return nil
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

	// Validate database path doesn't contain traversal attempts
	if strings.Contains(*dbPath, "..") {
		log.Fatalf("Database path contains invalid traversal sequence: %s", *dbPath)
	}
	staticPath, err := filepath.Abs(*staticDir)
	if err != nil {
		log.Printf("Warning: failed to resolve static path: %v", err)
		staticPath = *staticDir
	}

	// Initialize repositories
	questRepo := filesystem.NewQuestFileRepository(questsPath)
	refDataRepo, err := filesystem.NewReferenceDataFileRepository(dataPath)
	if err != nil {
		log.Fatalf("Failed to initialize reference data repository: %v", err)
	}

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
