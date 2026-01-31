package storage

import (
	"database/sql"
	"encoding/json"
	"fmt"

	_ "github.com/mattn/go-sqlite3"

	"github.com/tinx/pat-quest-editor/backend/internal/domain"
)

// SQLiteMetadataRepository implements MetadataRepository using SQLite.
type SQLiteMetadataRepository struct {
	db *sql.DB
}

// NewSQLiteMetadataRepository creates a new SQLite-based metadata repository.
func NewSQLiteMetadataRepository(dbPath string) (*SQLiteMetadataRepository, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool for SQLite (single writer, multiple readers)
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	// Verify connection is working
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	repo := &SQLiteMetadataRepository{db: db}
	if err := repo.initSchema(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return repo, nil
}

// Close closes the database connection.
func (r *SQLiteMetadataRepository) Close() error {
	return r.db.Close()
}

func (r *SQLiteMetadataRepository) initSchema() error {
	schema := `
		CREATE TABLE IF NOT EXISTS quest_metadata (
			quest_id TEXT PRIMARY KEY,
			node_positions TEXT NOT NULL
		);
	`
	_, err := r.db.Exec(schema)
	return err
}

// GetQuestMetadata retrieves editor metadata for a quest.
func (r *SQLiteMetadataRepository) GetQuestMetadata(questID string) (*domain.QuestMetadata, error) {
	var positionsJSON string
	err := r.db.QueryRow(
		"SELECT node_positions FROM quest_metadata WHERE quest_id = ?",
		questID,
	).Scan(&positionsJSON)

	if err == sql.ErrNoRows {
		// Return empty metadata if none exists
		return &domain.QuestMetadata{
			QuestID:       questID,
			NodePositions: make(map[int]domain.NodePosition),
		}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get metadata: %w", err)
	}

	var positions map[int]domain.NodePosition
	if err := json.Unmarshal([]byte(positionsJSON), &positions); err != nil {
		return nil, fmt.Errorf("failed to parse positions: %w", err)
	}

	return &domain.QuestMetadata{
		QuestID:       questID,
		NodePositions: positions,
	}, nil
}

// SaveQuestMetadata persists editor metadata for a quest.
func (r *SQLiteMetadataRepository) SaveQuestMetadata(metadata *domain.QuestMetadata) error {
	positionsJSON, err := json.Marshal(metadata.NodePositions)
	if err != nil {
		return fmt.Errorf("failed to marshal positions: %w", err)
	}

	_, err = r.db.Exec(`
		INSERT INTO quest_metadata (quest_id, node_positions)
		VALUES (?, ?)
		ON CONFLICT(quest_id) DO UPDATE SET node_positions = excluded.node_positions
	`, metadata.QuestID, string(positionsJSON))

	if err != nil {
		return fmt.Errorf("failed to save metadata: %w", err)
	}

	return nil
}

// DeleteQuestMetadata removes editor metadata for a quest.
func (r *SQLiteMetadataRepository) DeleteQuestMetadata(questID string) error {
	_, err := r.db.Exec("DELETE FROM quest_metadata WHERE quest_id = ?", questID)
	if err != nil {
		return fmt.Errorf("failed to delete metadata: %w", err)
	}
	return nil
}
