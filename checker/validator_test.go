package main

import "testing"

func TestValidateUniqueNodeIDs(t *testing.T) {
	quest := &Quest{
		QuestID: "TestQuest",
		QuestNodes: []QuestNode{
			{NodeID: 0, NodeType: "EntryPoint"},
			{NodeID: 0, NodeType: "Actions"}, // Duplicate
		},
	}

	errors := validateUniqueNodeIDs(quest)

	if len(errors) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errors))
	}
	if errors[0].Message != "duplicate NodeID" {
		t.Errorf("unexpected message: %s", errors[0].Message)
	}
}

func TestValidateNodeConnections_SelfReference(t *testing.T) {
	quest := &Quest{
		QuestID: "TestQuest",
		QuestNodes: []QuestNode{
			{NodeID: 0, NodeType: "EntryPoint", NextNodes: []int{0}}, // Self-reference
		},
	}

	errors := validateNodeConnections(quest)

	found := false
	for _, err := range errors {
		if err.Message == "node references itself in NextNodes" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected self-reference error")
	}
}

func TestValidateNodeConnections_DuplicateEdge(t *testing.T) {
	quest := &Quest{
		QuestID: "TestQuest",
		QuestNodes: []QuestNode{
			{NodeID: 0, NodeType: "EntryPoint", NextNodes: []int{1, 1}}, // Duplicate
			{NodeID: 1, NodeType: "Actions", Actions: []interface{}{"CompleteQuest"}},
		},
	}

	errors := validateNodeConnections(quest)

	found := false
	for _, err := range errors {
		if err.Message == "duplicate edge to node 1" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected duplicate edge error")
	}
}

func TestValidateNodeConnections_NonExistent(t *testing.T) {
	quest := &Quest{
		QuestID: "TestQuest",
		QuestNodes: []QuestNode{
			{NodeID: 0, NodeType: "EntryPoint", NextNodes: []int{99}},
		},
	}

	errors := validateNodeConnections(quest)

	found := false
	for _, err := range errors {
		if err.Message == "NextNodes references non-existent node 99" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected non-existent node error")
	}
}

func TestValidateEntryPoints_Missing(t *testing.T) {
	quest := &Quest{
		QuestID: "TestQuest",
		QuestNodes: []QuestNode{
			{NodeID: 0, NodeType: "Actions"},
		},
	}

	errors := validateEntryPoints(quest)

	if len(errors) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errors))
	}
	if errors[0].Message != "quest must have at least one EntryPoint node" {
		t.Errorf("unexpected message: %s", errors[0].Message)
	}
}

func TestValidateTerminalNodes_WithNextNodes(t *testing.T) {
	quest := &Quest{
		QuestID: "TestQuest",
		QuestNodes: []QuestNode{
			{NodeID: 0, NodeType: "EntryPoint", NextNodes: []int{1}},
			{NodeID: 1, NodeType: "Actions", Actions: []interface{}{"CompleteQuest"}, NextNodes: []int{2}},
			{NodeID: 2, NodeType: "Actions", Actions: []interface{}{"FailQuest"}},
		},
	}

	errors := validateTerminalNodes(quest)

	found := false
	for _, err := range errors {
		if err.Message == "terminal action node must not have NextNodes" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected terminal node with NextNodes error")
	}
}

func TestValidateTerminalNodes_MultipleTerminal(t *testing.T) {
	quest := &Quest{
		QuestID: "TestQuest",
		QuestNodes: []QuestNode{
			{NodeID: 0, NodeType: "EntryPoint", NextNodes: []int{1}},
			{NodeID: 1, NodeType: "Actions", Actions: []interface{}{"CompleteQuest", "FailQuest"}},
		},
	}

	errors := validateTerminalNodes(quest)

	found := false
	for _, err := range errors {
		if err.Message == "Actions node has more than one terminal action" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected multiple terminal actions error")
	}
}

func TestValidateNoCycles(t *testing.T) {
	quest := &Quest{
		QuestID: "TestQuest",
		QuestNodes: []QuestNode{
			{NodeID: 0, NodeType: "EntryPoint", NextNodes: []int{1}},
			{NodeID: 1, NodeType: "Dialog", NextNodes: []int{2}},
			{NodeID: 2, NodeType: "Dialog", NextNodes: []int{1}}, // Cycle
		},
	}

	errors := validateNoCycles(quest)

	if len(errors) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errors))
	}
	if errors[0].Message != "quest contains a cycle" {
		t.Errorf("unexpected message: %s", errors[0].Message)
	}
}

func TestValidateNoCycles_Valid(t *testing.T) {
	quest := &Quest{
		QuestID: "TestQuest",
		QuestNodes: []QuestNode{
			{NodeID: 0, NodeType: "EntryPoint", NextNodes: []int{1}},
			{NodeID: 1, NodeType: "Actions", Actions: []interface{}{"CompleteQuest"}},
		},
	}

	errors := validateNoCycles(quest)

	if len(errors) != 0 {
		t.Errorf("expected no errors, got %v", errors)
	}
}

func TestValidateQuest_Valid(t *testing.T) {
	quest := &Quest{
		QuestID: "TestQuest",
		QuestNodes: []QuestNode{
			{NodeID: 0, NodeType: "EntryPoint", NextNodes: []int{1}},
			{NodeID: 1, NodeType: "Actions", Actions: []interface{}{"CompleteQuest"}},
		},
	}
	refData := &ReferenceData{
		NPCs:      make(map[string]bool),
		Items:     make(map[string]bool),
		Factions:  make(map[string]bool),
		Resources: make(map[string]bool),
		Objects:   make(map[string]bool),
	}

	errors := ValidateQuest(quest, refData)

	if len(errors) != 0 {
		t.Errorf("expected no errors, got %v", errors)
	}
}

// Cross-quest validation tests

func TestValidateUniqueQuestIDs(t *testing.T) {
	quests := []*Quest{
		{QuestID: "Quest1"},
		{QuestID: "Quest2"},
		{QuestID: "Quest1"}, // Duplicate
	}

	errors := validateUniqueQuestIDs(quests)

	if len(errors) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errors))
	}
	if errors[0].Message != `duplicate QuestID "Quest1" found 2 times` {
		t.Errorf("unexpected message: %s", errors[0].Message)
	}
}

func TestValidateUniqueDisplayNames(t *testing.T) {
	quests := []*Quest{
		{QuestID: "Quest1", DisplayName: I18nString{EnUS: "My Quest", DeDE: "Meine Quest"}},
		{QuestID: "Quest2", DisplayName: I18nString{EnUS: "My Quest", DeDE: "Andere Quest"}},
	}

	errors := validateUniqueDisplayNames(quests)

	found := false
	for _, err := range errors {
		if err.Message == `duplicate DisplayName "My Quest" (en-US) in quests: Quest1, Quest2` {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected duplicate DisplayName error, got %v", errors)
	}
}

func TestValidateQuestReferences_Invalid(t *testing.T) {
	quests := []*Quest{
		{QuestID: "Quest1"},
		{
			QuestID: "Quest2",
			QuestNodes: []QuestNode{
				{
					NodeID:     1,
					NodeType:   "ConditionWatcher",
					Conditions: []map[string]interface{}{{"QuestCompleted": "NonExistent"}},
				},
			},
		},
	}

	questIDs := buildQuestIDSet(quests)
	errors := validateQuestReferences(quests, questIDs)

	if len(errors) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errors))
	}
	if errors[0].Message != `QuestCompleted references non-existent quest "NonExistent"` {
		t.Errorf("unexpected message: %s", errors[0].Message)
	}
}

func TestValidateQuestReferences_Valid(t *testing.T) {
	quests := []*Quest{
		{QuestID: "Quest1"},
		{
			QuestID: "Quest2",
			QuestNodes: []QuestNode{
				{
					NodeID:     1,
					NodeType:   "ConditionWatcher",
					Conditions: []map[string]interface{}{{"QuestCompleted": "Quest1"}},
				},
			},
		},
	}

	questIDs := buildQuestIDSet(quests)
	errors := validateQuestReferences(quests, questIDs)

	if len(errors) != 0 {
		t.Errorf("expected no errors, got %v", errors)
	}
}

func TestFormatError(t *testing.T) {
	nodeID := 5

	tests := []struct {
		name     string
		err      ValidationError
		expected string
	}{
		{
			name:     "quest with node",
			err:      ValidationError{QuestID: "Quest1", NodeID: &nodeID, Message: "some error"},
			expected: "[Quest1] Node 5: some error",
		},
		{
			name:     "quest without node",
			err:      ValidationError{QuestID: "Quest1", Message: "some error"},
			expected: "[Quest1]: some error",
		},
		{
			name:     "cross-quest error",
			err:      ValidationError{Message: "duplicate found"},
			expected: "[CROSS-QUEST]: duplicate found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatError(tt.err); got != tt.expected {
				t.Errorf("got %q, want %q", got, tt.expected)
			}
		})
	}
}
