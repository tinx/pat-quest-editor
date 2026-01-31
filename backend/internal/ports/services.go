package ports

import "github.com/tinx/pat-quest-editor/backend/internal/domain"

// QuestValidator defines the interface for quest validation.
type QuestValidator interface {
	// Validate checks a quest against all rules and returns validation results.
	Validate(quest *domain.Quest) *domain.ValidationResult
}
