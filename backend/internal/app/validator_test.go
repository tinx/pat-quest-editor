package app

import (
	"testing"

	"github.com/tinx/pat-quest-editor/backend/internal/domain"
)

// mockReferenceData implements ports.ReferenceDataRepository for testing.
type mockReferenceData struct{}

func (m *mockReferenceData) ListItems() ([]domain.Item, error) {
	return []domain.Item{{ItemID: "PackOfNails"}}, nil
}

func (m *mockReferenceData) ListFactions() ([]domain.Faction, error) {
	return []domain.Faction{
		{FactionID: "NPC:Smith"},
		{FactionID: "NPC:Carpenter"},
		{FactionID: "Town"},
	}, nil
}

func (m *mockReferenceData) ListResources() ([]domain.Resource, error) {
	return []domain.Resource{
		{ResourceID: "Coal"},
		{ResourceID: "IronOre"},
	}, nil
}

func (m *mockReferenceData) ListNPCs() ([]domain.NPC, error) {
	return []domain.NPC{
		{NPCID: "NPC:Smith"},
		{NPCID: "NPC:Carpenter"},
	}, nil
}

func (m *mockReferenceData) GetItem(itemID string) (*domain.Item, error)       { return nil, nil }
func (m *mockReferenceData) GetFaction(factionID string) (*domain.Faction, error) { return nil, nil }
func (m *mockReferenceData) GetResource(resourceID string) (*domain.Resource, error) { return nil, nil }
func (m *mockReferenceData) GetNPC(npcID string) (*domain.NPC, error)           { return nil, nil }

func TestValidate_ValidQuest(t *testing.T) {
	validator := NewQuestValidatorService(&mockReferenceData{})

	quest := &domain.Quest{
		QuestID: "TestQuest",
		QuestNodes: []domain.QuestNode{
			{NodeID: 0, NodeType: "EntryPoint", NextNodes: []int{1}},
			{NodeID: 1, NodeType: "Actions", Actions: []domain.Action{"CompleteQuest"}},
		},
	}

	result := validator.Validate(quest)

	if !result.Valid {
		t.Errorf("expected valid quest, got errors: %v", result.Errors)
	}
}

func TestValidate_DuplicateNodeIDs(t *testing.T) {
	validator := NewQuestValidatorService(&mockReferenceData{})

	quest := &domain.Quest{
		QuestID: "TestQuest",
		QuestNodes: []domain.QuestNode{
			{NodeID: 0, NodeType: "EntryPoint", NextNodes: []int{1}},
			{NodeID: 0, NodeType: "Actions", Actions: []domain.Action{"CompleteQuest"}},
		},
	}

	result := validator.Validate(quest)

	if result.Valid {
		t.Error("expected invalid quest due to duplicate NodeIDs")
	}

	found := false
	for _, err := range result.Errors {
		if err.Message == "duplicate NodeID" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected duplicate NodeID error")
	}
}

func TestValidate_NoEntryPoint(t *testing.T) {
	validator := NewQuestValidatorService(&mockReferenceData{})

	quest := &domain.Quest{
		QuestID: "TestQuest",
		QuestNodes: []domain.QuestNode{
			{NodeID: 0, NodeType: "Actions", Actions: []domain.Action{"CompleteQuest"}},
		},
	}

	result := validator.Validate(quest)

	if result.Valid {
		t.Error("expected invalid quest due to missing EntryPoint")
	}
}

func TestValidate_CycleDetection(t *testing.T) {
	validator := NewQuestValidatorService(&mockReferenceData{})

	quest := &domain.Quest{
		QuestID: "TestQuest",
		QuestNodes: []domain.QuestNode{
			{NodeID: 0, NodeType: "EntryPoint", NextNodes: []int{1}},
			{NodeID: 1, NodeType: "Dialog", NextNodes: []int{2}, ConversationPartner: "NPC:Smith", Messages: []domain.DialogMessage{{Speaker: "NPC:Smith"}}},
			{NodeID: 2, NodeType: "Dialog", NextNodes: []int{1}, ConversationPartner: "NPC:Smith", Messages: []domain.DialogMessage{{Speaker: "NPC:Smith"}}}, // Creates cycle
		},
	}

	result := validator.Validate(quest)

	if result.Valid {
		t.Error("expected invalid quest due to cycle")
	}

	found := false
	for _, err := range result.Errors {
		if err.Message == "quest contains a cycle (loops are not allowed)" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected cycle detection error")
	}
}

func TestValidate_TerminalNodeWithNextNodes(t *testing.T) {
	validator := NewQuestValidatorService(&mockReferenceData{})

	quest := &domain.Quest{
		QuestID: "TestQuest",
		QuestNodes: []domain.QuestNode{
			{NodeID: 0, NodeType: "EntryPoint", NextNodes: []int{1}},
			{NodeID: 1, NodeType: "Actions", Actions: []domain.Action{"CompleteQuest"}, NextNodes: []int{2}},
			{NodeID: 2, NodeType: "Actions", Actions: []domain.Action{"FailQuest"}},
		},
	}

	result := validator.Validate(quest)

	if result.Valid {
		t.Error("expected invalid quest due to terminal node with NextNodes")
	}
}

func TestValidate_UnknownNPC(t *testing.T) {
	validator := NewQuestValidatorService(&mockReferenceData{})

	quest := &domain.Quest{
		QuestID: "TestQuest",
		QuestNodes: []domain.QuestNode{
			{NodeID: 0, NodeType: "EntryPoint", NextNodes: []int{1}},
			{NodeID: 1, NodeType: "Dialog", ConversationPartner: "NPC:Unknown", Messages: []domain.DialogMessage{{Speaker: "NPC:Unknown"}}, NextNodes: []int{2}},
			{NodeID: 2, NodeType: "Actions", Actions: []domain.Action{"CompleteQuest"}},
		},
	}

	result := validator.Validate(quest)

	if result.Valid {
		t.Error("expected invalid quest due to unknown NPC")
	}
}

func TestValidate_PlayerSpeakerIsValid(t *testing.T) {
	validator := NewQuestValidatorService(&mockReferenceData{})

	quest := &domain.Quest{
		QuestID: "TestQuest",
		QuestNodes: []domain.QuestNode{
			{NodeID: 0, NodeType: "EntryPoint", NextNodes: []int{1}},
			{NodeID: 1, NodeType: "Dialog", ConversationPartner: "NPC:Smith", Messages: []domain.DialogMessage{
				{Speaker: "NPC:Smith", Text: domain.I18nString{EnUS: "Hello!"}},
				{Speaker: "Player", Text: domain.I18nString{EnUS: "Hi there!"}},
			}, NextNodes: []int{2}},
			{NodeID: 2, NodeType: "Actions", Actions: []domain.Action{"CompleteQuest"}},
		},
	}

	result := validator.Validate(quest)

	if !result.Valid {
		t.Errorf("expected valid quest with Player speaker, got errors: %v", result.Errors)
	}
}

func TestValidate_EmptyQuestProgressor(t *testing.T) {
	validator := NewQuestValidatorService(&mockReferenceData{})

	quest := &domain.Quest{
		QuestID: "TestQuest",
		QuestNodes: []domain.QuestNode{
			{NodeID: 0, NodeType: "EntryPoint", NextNodes: []int{1}},
			{NodeID: 1, NodeType: "QuestProgress", QuestProgressors: []string{"NPC:Smith", ""}, NextNodes: []int{2}},
			{NodeID: 2, NodeType: "Actions", Actions: []domain.Action{"CompleteQuest"}},
		},
	}

	result := validator.Validate(quest)

	if result.Valid {
		t.Error("expected invalid quest due to empty quest progressor")
	}

	found := false
	for _, err := range result.Errors {
		if err.Message == "quest progressor 2 is empty" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected empty quest progressor error, got: %v", result.Errors)
	}
}

func TestValidate_UnknownQuestProgressor(t *testing.T) {
	validator := NewQuestValidatorService(&mockReferenceData{})

	quest := &domain.Quest{
		QuestID: "TestQuest",
		QuestNodes: []domain.QuestNode{
			{NodeID: 0, NodeType: "EntryPoint", NextNodes: []int{1}},
			{NodeID: 1, NodeType: "QuestProgress", QuestProgressors: []string{"NPC:Unknown"}, NextNodes: []int{2}},
			{NodeID: 2, NodeType: "Actions", Actions: []domain.Action{"CompleteQuest"}},
		},
	}

	result := validator.Validate(quest)

	if result.Valid {
		t.Error("expected invalid quest due to unknown quest progressor")
	}
}
