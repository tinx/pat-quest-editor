package filesystem

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/tinx/pat-quest-editor/backend/internal/domain"
)

// ReferenceDataFileRepository implements ReferenceDataRepository using the filesystem.
type ReferenceDataFileRepository struct {
	basePath      string
	itemsPath     string
	factionsPath  string
	resourcesPath string
	npcsPath      string
	objectsPath   string
}

// NewReferenceDataFileRepository creates a new filesystem-based reference data repository.
// Returns error if any constructed path would escape the base directory.
func NewReferenceDataFileRepository(dataPath string) (*ReferenceDataFileRepository, error) {
	absBase, err := filepath.Abs(dataPath)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve data path: %w", err)
	}

	// Construct paths using filepath.Join for safety
	itemsPath := filepath.Join(absBase, "items.yaml")
	factionsPath := filepath.Join(absBase, "factions.yaml")
	resourcesPath := filepath.Join(absBase, "resources.yaml")
	npcsPath := filepath.Join(absBase, "npcs.yaml")
	objectsPath := filepath.Join(absBase, "objects.yaml")

	// Validate all paths are within base directory
	for name, path := range map[string]string{
		"items":     itemsPath,
		"factions":  factionsPath,
		"resources": resourcesPath,
		"npcs":      npcsPath,
		"objects":   objectsPath,
	} {
		if err := validatePathWithinBase(absBase, path); err != nil {
			return nil, fmt.Errorf("invalid %s path: %w", name, err)
		}
	}

	return &ReferenceDataFileRepository{
		basePath:      absBase,
		itemsPath:     itemsPath,
		factionsPath:  factionsPath,
		resourcesPath: resourcesPath,
		npcsPath:      npcsPath,
		objectsPath:   objectsPath,
	}, nil
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

// ListObjects returns all defined world objects.
func (r *ReferenceDataFileRepository) ListObjects() ([]domain.Object, error) {
	var objects []domain.Object
	if err := r.loadYAMLFile(r.objectsPath, &objects); err != nil {
		return nil, fmt.Errorf("failed to load objects: %w", err)
	}
	return objects, nil
}

// GetObject retrieves a world object by ID.
func (r *ReferenceDataFileRepository) GetObject(objectID string) (*domain.Object, error) {
	objects, err := r.ListObjects()
	if err != nil {
		return nil, err
	}
	for _, obj := range objects {
		if obj.ObjectID == objectID {
			return &obj, nil
		}
	}
	return nil, fmt.Errorf("object not found: %s", objectID)
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
