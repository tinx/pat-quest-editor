package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// LoadQuests loads all quest files from the given directory.
func LoadQuests(questsPath string) ([]*Quest, []error) {
	var quests []*Quest
	var errors []error

	err := filepath.Walk(questsPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".yaml") && !strings.HasSuffix(path, ".yml") {
			return nil
		}

		quest, err := loadQuestFile(path)
		if err != nil {
			errors = append(errors, fmt.Errorf("failed to load %s: %w", path, err))
			return nil
		}
		quests = append(quests, quest)
		return nil
	})

	if err != nil {
		errors = append(errors, fmt.Errorf("failed to walk quests directory: %w", err))
	}

	return quests, errors
}

func loadQuestFile(path string) (*Quest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var quest Quest
	if err := yaml.Unmarshal(data, &quest); err != nil {
		return nil, err
	}

	return &quest, nil
}

// LoadReferenceData loads all reference data from the given directory.
func LoadReferenceData(dataPath string) (*ReferenceData, error) {
	refData := &ReferenceData{
		NPCs:      make(map[string]bool),
		Items:     make(map[string]bool),
		Factions:  make(map[string]bool),
		Resources: make(map[string]bool),
		Objects:   make(map[string]bool),
	}

	// Load NPCs
	npcs, err := loadYAMLList[NPC](filepath.Join(dataPath, "npcs.yaml"))
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to load NPCs: %w", err)
	}
	for _, npc := range npcs {
		refData.NPCs[npc.NPCID] = true
	}

	// Load Items
	items, err := loadYAMLList[Item](filepath.Join(dataPath, "items.yaml"))
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to load items: %w", err)
	}
	for _, item := range items {
		refData.Items[item.ItemID] = true
	}

	// Load Factions
	factions, err := loadYAMLList[Faction](filepath.Join(dataPath, "factions.yaml"))
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to load factions: %w", err)
	}
	for _, faction := range factions {
		refData.Factions[faction.FactionID] = true
	}

	// Load Resources
	resources, err := loadYAMLList[Resource](filepath.Join(dataPath, "resources.yaml"))
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to load resources: %w", err)
	}
	for _, resource := range resources {
		refData.Resources[resource.ResourceID] = true
	}

	// Load Objects
	objects, err := loadYAMLList[Object](filepath.Join(dataPath, "objects.yaml"))
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to load objects: %w", err)
	}
	for _, obj := range objects {
		refData.Objects[obj.ObjectID] = true
	}

	return refData, nil
}

func loadYAMLList[T any](path string) ([]T, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var items []T
	if err := yaml.Unmarshal(data, &items); err != nil {
		return nil, err
	}

	return items, nil
}
