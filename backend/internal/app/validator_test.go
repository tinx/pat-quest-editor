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
			{NodeID: 1, NodeType: "Actions", Actions: []domain.Action{
				map[string]interface{}{"JournalEntry": map[string]interface{}{}},
				map[string]interface{}{"QuestStageDescription": map[string]interface{}{}},
				"CompleteQuest",
			}},
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
			{NodeID: 2, NodeType: "Actions", Actions: []domain.Action{
				map[string]interface{}{"JournalEntry": map[string]interface{}{}},
				map[string]interface{}{"QuestStageDescription": map[string]interface{}{}},
				"CompleteQuest",
			}},
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
			{NodeID: 2, NodeType: "Actions", Actions: []domain.Action{
				map[string]interface{}{"JournalEntry": map[string]interface{}{}},
				map[string]interface{}{"QuestStageDescription": map[string]interface{}{}},
				"CompleteQuest",
			}},
			{NodeID: 3, NodeType: "Actions", Actions: []domain.Action{
				map[string]interface{}{"JournalEntry": map[string]interface{}{}},
				map[string]interface{}{"QuestStageDescription": map[string]interface{}{}},
				"FailQuest",
			}},
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
			{NodeID: 2, NodeType: "Actions", Actions: []domain.Action{
				map[string]interface{}{"JournalEntry": map[string]interface{}{}},
				map[string]interface{}{"QuestStageDescription": map[string]interface{}{}},
				"CompleteQuest",
			}},
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
			{NodeID: 2, NodeType: "Actions", Actions: []domain.Action{
				map[string]interface{}{"JournalEntry": map[string]interface{}{}},
				map[string]interface{}{"QuestStageDescription": map[string]interface{}{}},
				"CompleteQuest",
			}},
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

// Helper to create a valid terminal Actions node with journal actions
func validTerminalActions(actions ...domain.Action) []domain.Action {
	result := []domain.Action{
		map[string]interface{}{"JournalEntry": map[string]interface{}{}},
		map[string]interface{}{"QuestStageDescription": map[string]interface{}{}},
	}
	return append(result, actions...)
}

// ============ Tests for new validation rules ============

func TestValidate_DuplicateEdges(t *testing.T) {
	validator := NewQuestValidatorService(&mockReferenceData{})

	quest := &domain.Quest{
		QuestID: "TestQuest",
		QuestNodes: []domain.QuestNode{
			{NodeID: 0, NodeType: "EntryPoint", NextNodes: []int{1, 1}}, // Duplicate edge
			{NodeID: 1, NodeType: "Actions", Actions: validTerminalActions("CompleteQuest")},
		},
	}

	result := validator.Validate(quest)

	if result.Valid {
		t.Error("expected invalid quest due to duplicate edges")
	}

	found := false
	for _, err := range result.Errors {
		if err.Message == "duplicate edge to NodeID 1" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected duplicate edge error, got: %v", result.Errors)
	}
}

func TestValidate_SelfReference(t *testing.T) {
	validator := NewQuestValidatorService(&mockReferenceData{})

	quest := &domain.Quest{
		QuestID: "TestQuest",
		QuestNodes: []domain.QuestNode{
			{NodeID: 0, NodeType: "EntryPoint", NextNodes: []int{1}},
			{NodeID: 1, NodeType: "Dialog", ConversationPartner: "NPC:Smith", NextNodes: []int{1}}, // Self-reference
		},
	}

	result := validator.Validate(quest)

	if result.Valid {
		t.Error("expected invalid quest due to self-reference")
	}

	found := false
	for _, err := range result.Errors {
		if err.Message == "node references itself in NextNodes" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected self-reference error, got: %v", result.Errors)
	}
}

func TestValidate_SelfReferenceInConditionBranch(t *testing.T) {
	validator := NewQuestValidatorService(&mockReferenceData{})

	quest := &domain.Quest{
		QuestID: "TestQuest",
		QuestNodes: []domain.QuestNode{
			{NodeID: 0, NodeType: "EntryPoint", NextNodes: []int{1}},
			{NodeID: 1, NodeType: "ConditionBranch",
				Conditions:      []domain.Condition{{"QuestCompleted": "SomeQuest"}},
				NextNodesIfTrue: []int{1}, // Self-reference
			},
		},
	}

	result := validator.Validate(quest)

	if result.Valid {
		t.Error("expected invalid quest due to self-reference in ConditionBranch")
	}

	found := false
	for _, err := range result.Errors {
		if err.Message == "node references itself in NextNodesIfTrue" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected self-reference error, got: %v", result.Errors)
	}
}

func TestValidate_NonTerminalActionsWithoutNextNodes(t *testing.T) {
	validator := NewQuestValidatorService(&mockReferenceData{})

	quest := &domain.Quest{
		QuestID: "TestQuest",
		QuestNodes: []domain.QuestNode{
			{NodeID: 0, NodeType: "EntryPoint", NextNodes: []int{1}},
			{NodeID: 1, NodeType: "Actions", Actions: []domain.Action{
				map[string]interface{}{"JournalEntry": map[string]interface{}{}},
			}}, // No terminal action and no NextNodes
		},
	}

	result := validator.Validate(quest)

	if result.Valid {
		t.Error("expected invalid quest due to non-terminal Actions without NextNodes")
	}

	found := false
	for _, err := range result.Errors {
		if err.Message == "non-terminal Actions node must have NextNodes (quest flow ends with unspecified behaviour)" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected unspecified behaviour error, got: %v", result.Errors)
	}
}

func TestValidate_UnreferencedNodeID(t *testing.T) {
	validator := NewQuestValidatorService(&mockReferenceData{})

	quest := &domain.Quest{
		QuestID: "TestQuest",
		QuestNodes: []domain.QuestNode{
			{NodeID: 0, NodeType: "EntryPoint", NextNodes: []int{1}},
			{NodeID: 1, NodeType: "Actions", Actions: validTerminalActions("CompleteQuest")},
			{NodeID: 2, NodeType: "Actions", Actions: validTerminalActions("FailQuest")}, // Never referenced
		},
	}

	result := validator.Validate(quest)

	if result.Valid {
		t.Error("expected invalid quest due to unreferenced NodeID")
	}

	found := false
	for _, err := range result.Errors {
		if err.Message == "NodeID is never referenced by any other node" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected unreferenced NodeID error, got: %v", result.Errors)
	}
}

func TestValidate_EntryPointNodeID_NotRequireReference(t *testing.T) {
	validator := NewQuestValidatorService(&mockReferenceData{})

	quest := &domain.Quest{
		QuestID: "TestQuest",
		QuestNodes: []domain.QuestNode{
			{NodeID: 0, NodeType: "EntryPoint", NextNodes: []int{1}}, // EntryPoints don't need to be referenced
			{NodeID: 1, NodeType: "Actions", Actions: validTerminalActions("CompleteQuest")},
		},
	}

	result := validator.Validate(quest)

	if !result.Valid {
		t.Errorf("expected valid quest (EntryPoint doesn't need reference), got errors: %v", result.Errors)
	}
}

func TestValidate_JournalAtFlowStart_Missing(t *testing.T) {
	validator := NewQuestValidatorService(&mockReferenceData{})

	quest := &domain.Quest{
		QuestID: "TestQuest",
		QuestNodes: []domain.QuestNode{
			{NodeID: 0, NodeType: "EntryPoint", NextNodes: []int{1}},
			{NodeID: 1, NodeType: "Actions", Actions: []domain.Action{
				map[string]interface{}{"JournalEntry": map[string]interface{}{}},
				"CompleteQuest",
			}}, // Missing QuestStageDescription
		},
	}

	result := validator.Validate(quest)

	if result.Valid {
		t.Error("expected invalid quest due to missing QuestStageDescription at flow start")
	}

	found := false
	for _, err := range result.Errors {
		if err.Message == "first Actions node in flow must have QuestStageDescription action" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected missing QuestStageDescription error, got: %v", result.Errors)
	}
}

func TestValidate_JournalAtFlowStart_MissingJournalEntry(t *testing.T) {
	validator := NewQuestValidatorService(&mockReferenceData{})

	quest := &domain.Quest{
		QuestID: "TestQuest",
		QuestNodes: []domain.QuestNode{
			{NodeID: 0, NodeType: "EntryPoint", NextNodes: []int{1}},
			{NodeID: 1, NodeType: "Actions", Actions: []domain.Action{
				map[string]interface{}{"QuestStageDescription": map[string]interface{}{}},
				"CompleteQuest",
			}}, // Missing JournalEntry
		},
	}

	result := validator.Validate(quest)

	if result.Valid {
		t.Error("expected invalid quest due to missing JournalEntry at flow start")
	}

	found := false
	for _, err := range result.Errors {
		if err.Message == "first Actions node in flow must have JournalEntry action" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected missing JournalEntry error, got: %v", result.Errors)
	}
}

func TestValidate_JournalAtFlowEnd_Missing(t *testing.T) {
	validator := NewQuestValidatorService(&mockReferenceData{})

	quest := &domain.Quest{
		QuestID: "TestQuest",
		QuestNodes: []domain.QuestNode{
			{NodeID: 0, NodeType: "EntryPoint", NextNodes: []int{1}},
			{NodeID: 1, NodeType: "Actions", Actions: []domain.Action{
				map[string]interface{}{"JournalEntry": map[string]interface{}{}},
				map[string]interface{}{"QuestStageDescription": map[string]interface{}{}},
			}, NextNodes: []int{2}},
			{NodeID: 2, NodeType: "Dialog", ConversationPartner: "NPC:Smith", Messages: []domain.DialogMessage{
				{Speaker: "NPC:Smith"},
			}, NextNodes: []int{3}},
			{NodeID: 3, NodeType: "Actions", Actions: []domain.Action{"CompleteQuest"}}, // No JournalEntry after last non-Actions node
		},
	}

	result := validator.Validate(quest)

	if result.Valid {
		t.Error("expected invalid quest due to missing JournalEntry at flow end")
	}

	found := false
	for _, err := range result.Errors {
		if err.Message == "terminal Actions chain must contain a JournalEntry action" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected missing JournalEntry in terminal chain error, got: %v", result.Errors)
	}
}

func TestValidate_JournalAtFlowEnd_InChain(t *testing.T) {
	validator := NewQuestValidatorService(&mockReferenceData{})

	quest := &domain.Quest{
		QuestID: "TestQuest",
		QuestNodes: []domain.QuestNode{
			{NodeID: 0, NodeType: "EntryPoint", NextNodes: []int{1}},
			{NodeID: 1, NodeType: "Actions", Actions: []domain.Action{
				map[string]interface{}{"JournalEntry": map[string]interface{}{}},
				map[string]interface{}{"QuestStageDescription": map[string]interface{}{}},
			}, NextNodes: []int{2}},
			{NodeID: 2, NodeType: "Actions", Actions: []domain.Action{
				map[string]interface{}{"JournalEntry": map[string]interface{}{}}, // JournalEntry in preceding Actions node
			}, NextNodes: []int{3}},
			{NodeID: 3, NodeType: "Actions", Actions: []domain.Action{"CompleteQuest"}}, // Terminal without JournalEntry is OK
		},
	}

	result := validator.Validate(quest)

	if !result.Valid {
		t.Errorf("expected valid quest (JournalEntry in Actions chain), got errors: %v", result.Errors)
	}
}

func TestValidate_MultipleEntryPoints_EachNeedsJournal(t *testing.T) {
	validator := NewQuestValidatorService(&mockReferenceData{})

	quest := &domain.Quest{
		QuestID: "TestQuest",
		QuestNodes: []domain.QuestNode{
			{NodeID: 0, NodeType: "EntryPoint", NextNodes: []int{2}},
			{NodeID: 1, NodeType: "EntryPoint", NextNodes: []int{3}},
			{NodeID: 2, NodeType: "Actions", Actions: validTerminalActions("CompleteQuest")},
			{NodeID: 3, NodeType: "Actions", Actions: []domain.Action{"FailQuest"}}, // Missing journal actions
		},
	}

	result := validator.Validate(quest)

	if result.Valid {
		t.Error("expected invalid quest: second EntryPoint flow missing journal actions")
	}

	// Should have errors about missing JournalEntry and QuestStageDescription for node 3
	foundJournal := false
	foundStage := false
	for _, err := range result.Errors {
		if *err.NodeID == 3 {
			if err.Message == "first Actions node in flow must have JournalEntry action" {
				foundJournal = true
			}
			if err.Message == "first Actions node in flow must have QuestStageDescription action" {
				foundStage = true
			}
		}
	}
	if !foundJournal || !foundStage {
		t.Errorf("expected journal errors for node 3, got: %v", result.Errors)
	}
}

func TestValidate_DuplicateEdgesInDecisionOptions(t *testing.T) {
	validator := NewQuestValidatorService(&mockReferenceData{})

	quest := &domain.Quest{
		QuestID: "TestQuest",
		QuestNodes: []domain.QuestNode{
			{NodeID: 0, NodeType: "EntryPoint", NextNodes: []int{1}},
			{NodeID: 1, NodeType: "Decision",
				ConversationPartner: "NPC:Smith",
				Options: []domain.DialogOption{
					{Text: domain.I18nString{EnUS: "Option 1"}, NextNodes: []int{2}},
					{Text: domain.I18nString{EnUS: "Option 2"}, NextNodes: []int{2}}, // Same target, but different options - this is OK
				},
			},
			{NodeID: 2, NodeType: "Actions", Actions: validTerminalActions("CompleteQuest")},
		},
	}

	result := validator.Validate(quest)

	// This should actually be valid - different options can lead to the same node
	// The duplicate check is for the same node referencing the same target multiple times
	if !result.Valid {
		t.Errorf("expected valid quest (different options to same node is OK), got errors: %v", result.Errors)
	}
}

func TestValidate_DuplicateEdgesWithinSameOption(t *testing.T) {
	validator := NewQuestValidatorService(&mockReferenceData{})

	quest := &domain.Quest{
		QuestID: "TestQuest",
		QuestNodes: []domain.QuestNode{
			{NodeID: 0, NodeType: "EntryPoint", NextNodes: []int{1}},
			{NodeID: 1, NodeType: "Decision",
				ConversationPartner: "NPC:Smith",
				Options: []domain.DialogOption{
					{Text: domain.I18nString{EnUS: "Option 1"}, NextNodes: []int{2, 2}}, // Duplicate within same option
				},
			},
			{NodeID: 2, NodeType: "Actions", Actions: validTerminalActions("CompleteQuest")},
		},
	}

	result := validator.Validate(quest)

	if result.Valid {
		t.Error("expected invalid quest due to duplicate edges within decision option")
	}

	found := false
	for _, err := range result.Errors {
		if err.Message == "duplicate edge to NodeID 2" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected duplicate edge error, got: %v", result.Errors)
	}
}

func TestValidate_SelfReferenceInDialogOption(t *testing.T) {
	validator := NewQuestValidatorService(&mockReferenceData{})

	quest := &domain.Quest{
		QuestID: "TestQuest",
		QuestNodes: []domain.QuestNode{
			{NodeID: 0, NodeType: "EntryPoint", NextNodes: []int{1}},
			{NodeID: 1, NodeType: "Decision",
				ConversationPartner: "NPC:Smith",
				Options: []domain.DialogOption{
					{Text: domain.I18nString{EnUS: "Option 1"}, NextNodes: []int{1}}, // Self-reference
				},
			},
		},
	}

	result := validator.Validate(quest)

	if result.Valid {
		t.Error("expected invalid quest due to self-reference in dialog option")
	}

	found := false
	for _, err := range result.Errors {
		if err.Message == "node references itself in option 1 NextNodes" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected self-reference error, got: %v", result.Errors)
	}
}
