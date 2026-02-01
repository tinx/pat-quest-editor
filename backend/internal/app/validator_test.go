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
func (m *mockReferenceData) ListObjects() ([]domain.Object, error) {
	return []domain.Object{
		{ObjectID: "Object:Lorry"},
		{ObjectID: "Object:StableLamp"},
	}, nil
}
func (m *mockReferenceData) GetObject(objectID string) (*domain.Object, error) { return nil, nil }

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

func TestValidate_ConditionBranch_Valid(t *testing.T) {
	validator := NewQuestValidatorService(&mockReferenceData{})

	quest := &domain.Quest{
		QuestID: "TestQuest",
		QuestNodes: []domain.QuestNode{
			{NodeID: 0, NodeType: "EntryPoint", NextNodes: []int{1}},
			{NodeID: 1, NodeType: "ConditionBranch",
				Conditions:       []domain.Condition{{"QuestCompleted": "SomeQuest"}},
				NextNodesIfTrue:  []int{2},
				NextNodesIfFalse: []int{3},
			},
			{NodeID: 2, NodeType: "Actions", Actions: []domain.Action{"CompleteQuest"}},
			{NodeID: 3, NodeType: "Actions", Actions: []domain.Action{"FailQuest"}},
		},
	}

	result := validator.Validate(quest)

	if !result.Valid {
		t.Errorf("expected valid quest with ConditionBranch, got errors: %v", result.Errors)
	}
}

func TestValidate_ConditionBranch_OnlyTrueBranch(t *testing.T) {
	validator := NewQuestValidatorService(&mockReferenceData{})

	quest := &domain.Quest{
		QuestID: "TestQuest",
		QuestNodes: []domain.QuestNode{
			{NodeID: 0, NodeType: "EntryPoint", NextNodes: []int{1}},
			{NodeID: 1, NodeType: "ConditionBranch",
				Conditions:      []domain.Condition{{"QuestCompleted": "SomeQuest"}},
				NextNodesIfTrue: []int{2},
			},
			{NodeID: 2, NodeType: "Actions", Actions: []domain.Action{"CompleteQuest"}},
		},
	}

	result := validator.Validate(quest)

	if !result.Valid {
		t.Errorf("expected valid quest with only true branch, got errors: %v", result.Errors)
	}
}

func TestValidate_ConditionBranch_OnlyFalseBranch(t *testing.T) {
	validator := NewQuestValidatorService(&mockReferenceData{})

	quest := &domain.Quest{
		QuestID: "TestQuest",
		QuestNodes: []domain.QuestNode{
			{NodeID: 0, NodeType: "EntryPoint", NextNodes: []int{1}},
			{NodeID: 1, NodeType: "ConditionBranch",
				Conditions:       []domain.Condition{{"QuestCompleted": "SomeQuest"}},
				NextNodesIfFalse: []int{2},
			},
			{NodeID: 2, NodeType: "Actions", Actions: []domain.Action{"CompleteQuest"}},
		},
	}

	result := validator.Validate(quest)

	if !result.Valid {
		t.Errorf("expected valid quest with only false branch, got errors: %v", result.Errors)
	}
}

func TestValidate_ConditionBranch_NoBranches(t *testing.T) {
	validator := NewQuestValidatorService(&mockReferenceData{})

	quest := &domain.Quest{
		QuestID: "TestQuest",
		QuestNodes: []domain.QuestNode{
			{NodeID: 0, NodeType: "EntryPoint", NextNodes: []int{1}},
			{NodeID: 1, NodeType: "ConditionBranch",
				Conditions: []domain.Condition{{"QuestCompleted": "SomeQuest"}},
			},
		},
	}

	result := validator.Validate(quest)

	if result.Valid {
		t.Error("expected invalid quest: ConditionBranch with no branches")
	}

	found := false
	for _, err := range result.Errors {
		if err.Message == "ConditionBranch must have at least one of NextNodesIfTrue or NextNodesIfFalse" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected error about missing branches, got: %v", result.Errors)
	}
}

func TestValidate_ConditionBranch_NoConditions(t *testing.T) {
	validator := NewQuestValidatorService(&mockReferenceData{})

	quest := &domain.Quest{
		QuestID: "TestQuest",
		QuestNodes: []domain.QuestNode{
			{NodeID: 0, NodeType: "EntryPoint", NextNodes: []int{1}},
			{NodeID: 1, NodeType: "ConditionBranch",
				NextNodesIfTrue: []int{2},
			},
			{NodeID: 2, NodeType: "Actions", Actions: []domain.Action{"CompleteQuest"}},
		},
	}

	result := validator.Validate(quest)

	if result.Valid {
		t.Error("expected invalid quest: ConditionBranch with no conditions")
	}

	found := false
	for _, err := range result.Errors {
		if err.Message == "ConditionBranch must have at least one condition" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected error about missing conditions, got: %v", result.Errors)
	}
}

func TestValidate_ConditionBranch_NoTopLevelNextNodes(t *testing.T) {
	validator := NewQuestValidatorService(&mockReferenceData{})

	quest := &domain.Quest{
		QuestID: "TestQuest",
		QuestNodes: []domain.QuestNode{
			{NodeID: 0, NodeType: "EntryPoint", NextNodes: []int{1}},
			{NodeID: 1, NodeType: "ConditionBranch",
				Conditions:      []domain.Condition{{"QuestCompleted": "SomeQuest"}},
				NextNodes:       []int{2}, // Should not use top-level NextNodes
				NextNodesIfTrue: []int{2},
			},
			{NodeID: 2, NodeType: "Actions", Actions: []domain.Action{"CompleteQuest"}},
		},
	}

	result := validator.Validate(quest)

	if result.Valid {
		t.Error("expected invalid quest: ConditionBranch with top-level NextNodes")
	}

	found := false
	for _, err := range result.Errors {
		if err.Message == "ConditionBranch must not have top-level NextNodes; use NextNodesIfTrue and NextNodesIfFalse instead" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected error about top-level NextNodes, got: %v", result.Errors)
	}
}

func TestValidate_ConditionBranch_InvalidReference(t *testing.T) {
	validator := NewQuestValidatorService(&mockReferenceData{})

	quest := &domain.Quest{
		QuestID: "TestQuest",
		QuestNodes: []domain.QuestNode{
			{NodeID: 0, NodeType: "EntryPoint", NextNodes: []int{1}},
			{NodeID: 1, NodeType: "ConditionBranch",
				Conditions:      []domain.Condition{{"QuestCompleted": "SomeQuest"}},
				NextNodesIfTrue: []int{99}, // Non-existent node
			},
		},
	}

	result := validator.Validate(quest)

	if result.Valid {
		t.Error("expected invalid quest: ConditionBranch referencing non-existent node")
	}

	found := false
	for _, err := range result.Errors {
		if err.Message == "NextNodesIfTrue references non-existent NodeID" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected error about non-existent NodeID, got: %v", result.Errors)
	}
}

func TestValidate_ConditionBranch_CycleDetection(t *testing.T) {
	validator := NewQuestValidatorService(&mockReferenceData{})

	quest := &domain.Quest{
		QuestID: "TestQuest",
		QuestNodes: []domain.QuestNode{
			{NodeID: 0, NodeType: "EntryPoint", NextNodes: []int{1}},
			{NodeID: 1, NodeType: "ConditionBranch",
				Conditions:      []domain.Condition{{"QuestCompleted": "SomeQuest"}},
				NextNodesIfTrue: []int{1}, // Cycle back to self
			},
		},
	}

	result := validator.Validate(quest)

	if result.Valid {
		t.Error("expected invalid quest: ConditionBranch with cycle")
	}

	found := false
	for _, err := range result.Errors {
		if err.Message == "quest contains a cycle (loops are not allowed)" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected cycle detection error, got: %v", result.Errors)
	}
}
