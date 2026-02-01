package ports

import "github.com/tinx/pat-quest-editor/backend/internal/domain"

// QuestRepository defines operations for quest file storage.
type QuestRepository interface {
	// List returns all quest IDs available in the repository.
	List() ([]string, error)
	
	// Get retrieves a quest by its ID.
	Get(questID string) (*domain.Quest, error)
	
	// Save persists a quest to storage.
	Save(quest *domain.Quest) error
	
	// Delete removes a quest from storage.
	Delete(questID string) error
	
	// Exists checks if a quest with the given ID exists.
	Exists(questID string) (bool, error)
}

// ReferenceDataRepository defines operations for reference data (items, factions, etc.).
type ReferenceDataRepository interface {
	// ListItems returns all defined items.
	ListItems() ([]domain.Item, error)
	
	// ListFactions returns all defined factions.
	ListFactions() ([]domain.Faction, error)
	
	// ListResources returns all defined resources.
	ListResources() ([]domain.Resource, error)
	
	// ListNPCs returns all defined NPCs.
	ListNPCs() ([]domain.NPC, error)
	
	// ListObjects returns all defined world objects.
	ListObjects() ([]domain.Object, error)
	
	// GetItem retrieves an item by ID.
	GetItem(itemID string) (*domain.Item, error)
	
	// GetFaction retrieves a faction by ID.
	GetFaction(factionID string) (*domain.Faction, error)
	
	// GetResource retrieves a resource by ID.
	GetResource(resourceID string) (*domain.Resource, error)
	
	// GetNPC retrieves an NPC by ID.
	GetNPC(npcID string) (*domain.NPC, error)
	
	// GetObject retrieves a world object by ID.
	GetObject(objectID string) (*domain.Object, error)
}

// MetadataRepository defines operations for editor metadata storage.
type MetadataRepository interface {
	// GetQuestMetadata retrieves editor metadata for a quest.
	GetQuestMetadata(questID string) (*domain.QuestMetadata, error)
	
	// SaveQuestMetadata persists editor metadata for a quest.
	SaveQuestMetadata(metadata *domain.QuestMetadata) error
	
	// DeleteQuestMetadata removes editor metadata for a quest.
	DeleteQuestMetadata(questID string) error
}
