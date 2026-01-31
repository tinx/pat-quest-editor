package filesystem

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/tinx/pat-quest-editor/backend/internal/domain"
)

// ReferenceDataFileRepository implements ReferenceDataRepository using the filesystem.
type ReferenceDataFileRepository struct {
	itemsPath     string
	factionsPath  string
	resourcesPath string
	npcsPath      string
}

// NewReferenceDataFileRepository creates a new filesystem-based reference data repository.
func NewReferenceDataFileRepository(dataPath string) *ReferenceDataFileRepository {
	return &ReferenceDataFileRepository{
		itemsPath:     dataPath + "/items.yaml",
		factionsPath:  dataPath + "/factions.yaml",
		resourcesPath: dataPath + "/resources.yaml",
		npcsPath:      dataPath + "/npcs.yaml",
	}
}

// ListItems returns all defined items.
func (r *ReferenceDataFileRepository) ListItems() ([]domain.Item, error) {
	var items []domain.Item
	if err := r.loadYAMLFile(r.itemsPath, &items); err != nil {
		return nil, fmt.Errorf("failed to load items: %w", err)
	}
	return items, nil
}

// ListFactions returns all defined factions.
func (r *ReferenceDataFileRepository) ListFactions() ([]domain.Faction, error) {
	var factions []domain.Faction
	if err := r.loadYAMLFile(r.factionsPath, &factions); err != nil {
		return nil, fmt.Errorf("failed to load factions: %w", err)
	}
	return factions, nil
}

// ListResources returns all defined resources.
func (r *ReferenceDataFileRepository) ListResources() ([]domain.Resource, error) {
	var resources []domain.Resource
	if err := r.loadYAMLFile(r.resourcesPath, &resources); err != nil {
		return nil, fmt.Errorf("failed to load resources: %w", err)
	}
	return resources, nil
}

// ListNPCs returns all defined NPCs.
func (r *ReferenceDataFileRepository) ListNPCs() ([]domain.NPC, error) {
	var npcs []domain.NPC
	if err := r.loadYAMLFile(r.npcsPath, &npcs); err != nil {
		return nil, fmt.Errorf("failed to load NPCs: %w", err)
	}
	return npcs, nil
}

// GetItem retrieves an item by ID.
func (r *ReferenceDataFileRepository) GetItem(itemID string) (*domain.Item, error) {
	items, err := r.ListItems()
	if err != nil {
		return nil, err
	}
	for _, item := range items {
		if item.ItemID == itemID {
			return &item, nil
		}
	}
	return nil, fmt.Errorf("item not found: %s", itemID)
}

// GetFaction retrieves a faction by ID.
func (r *ReferenceDataFileRepository) GetFaction(factionID string) (*domain.Faction, error) {
	factions, err := r.ListFactions()
	if err != nil {
		return nil, err
	}
	for _, faction := range factions {
		if faction.FactionID == factionID {
			return &faction, nil
		}
	}
	return nil, fmt.Errorf("faction not found: %s", factionID)
}

// GetResource retrieves a resource by ID.
func (r *ReferenceDataFileRepository) GetResource(resourceID string) (*domain.Resource, error) {
	resources, err := r.ListResources()
	if err != nil {
		return nil, err
	}
	for _, resource := range resources {
		if resource.ResourceID == resourceID {
			return &resource, nil
		}
	}
	return nil, fmt.Errorf("resource not found: %s", resourceID)
}

// GetNPC retrieves an NPC by ID.
func (r *ReferenceDataFileRepository) GetNPC(npcID string) (*domain.NPC, error) {
	npcs, err := r.ListNPCs()
	if err != nil {
		return nil, err
	}
	for _, npc := range npcs {
		if npc.NPCID == npcID {
			return &npc, nil
		}
	}
	return nil, fmt.Errorf("NPC not found: %s", npcID)
}

func (r *ReferenceDataFileRepository) loadYAMLFile(path string, v interface{}) error {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // Return empty if file doesn't exist
		}
		return err
	}
	return yaml.Unmarshal(data, v)
}
