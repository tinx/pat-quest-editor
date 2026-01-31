package filesystem

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/tinx/pat-quest-editor/backend/internal/domain"
)

// QuestFileRepository implements QuestRepository using the filesystem.
type QuestFileRepository struct {
	basePath string
}

// NewQuestFileRepository creates a new filesystem-based quest repository.
func NewQuestFileRepository(basePath string) *QuestFileRepository {
	return &QuestFileRepository{basePath: basePath}
}

// List returns all quest IDs available in the repository.
func (r *QuestFileRepository) List() ([]string, error) {
	var questIDs []string

	err := filepath.Walk(r.basePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && (strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml")) {
			quest, err := r.loadQuestFile(path)
			if err != nil {
				return nil // Skip files that can't be parsed
			}
			questIDs = append(questIDs, quest.QuestID)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to list quests: %w", err)
	}

	return questIDs, nil
}

// Get retrieves a quest by its ID.
func (r *QuestFileRepository) Get(questID string) (*domain.Quest, error) {
	path, err := r.findQuestFile(questID)
	if err != nil {
		return nil, err
	}
	return r.loadQuestFile(path)
}

// Save persists a quest to storage.
func (r *QuestFileRepository) Save(quest *domain.Quest) error {
	// Try to find existing file, otherwise create new one
	path, err := r.findQuestFile(quest.QuestID)
	if err != nil {
		// Create new file with sanitized quest ID as filename
		filename := sanitizeFilename(quest.QuestID) + ".yaml"
		path = filepath.Join(r.basePath, filename)
	}

	data, err := yaml.Marshal(quest)
	if err != nil {
		return fmt.Errorf("failed to marshal quest: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write quest file: %w", err)
	}

	return nil
}

// Delete removes a quest from storage.
func (r *QuestFileRepository) Delete(questID string) error {
	path, err := r.findQuestFile(questID)
	if err != nil {
		return err
	}
	return os.Remove(path)
}

// Exists checks if a quest with the given ID exists.
func (r *QuestFileRepository) Exists(questID string) (bool, error) {
	_, err := r.findQuestFile(questID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (r *QuestFileRepository) findQuestFile(questID string) (string, error) {
	var foundPath string

	err := filepath.Walk(r.basePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && (strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml")) {
			quest, err := r.loadQuestFile(path)
			if err != nil {
				return nil
			}
			if quest.QuestID == questID {
				foundPath = path
				return filepath.SkipAll
			}
		}
		return nil
	})

	if err != nil {
		return "", fmt.Errorf("error searching for quest: %w", err)
	}

	if foundPath == "" {
		return "", fmt.Errorf("quest not found: %s", questID)
	}

	return foundPath, nil
}

func (r *QuestFileRepository) loadQuestFile(path string) (*domain.Quest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read quest file: %w", err)
	}

	var quest domain.Quest
	if err := yaml.Unmarshal(data, &quest); err != nil {
		return nil, fmt.Errorf("failed to parse quest file: %w", err)
	}

	return &quest, nil
}

func sanitizeFilename(s string) string {
	// Replace characters that are problematic in filenames
	replacer := strings.NewReplacer(
		":", "_",
		"/", "_",
		"\\", "_",
		" ", "_",
	)
	return replacer.Replace(s)
}
