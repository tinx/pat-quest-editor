package app

import (
	"fmt"

	"github.com/tinx/pat-quest-editor/backend/internal/domain"
	"github.com/tinx/pat-quest-editor/backend/internal/ports"
)

// QuestValidatorService implements quest validation logic.
type QuestValidatorService struct {
	refData ports.ReferenceDataRepository
}

// NewQuestValidatorService creates a new quest validator.
func NewQuestValidatorService(refData ports.ReferenceDataRepository) *QuestValidatorService {
	return &QuestValidatorService{refData: refData}
}

// Validate checks a quest against all rules.
func (v *QuestValidatorService) Validate(quest *domain.Quest) *domain.ValidationResult {
	result := domain.NewValidationResult()

	v.validateUniqueNodeIDs(quest, result)
	v.validateNodeConnections(quest, result)
	v.validateEntryPoints(quest, result)
	v.validateTerminalNodes(quest, result)
	v.validatePlayerDecisionDialogs(quest, result)
	v.validateNoCycles(quest, result)
	v.validateReferences(quest, result)

	return result
}

func (v *QuestValidatorService) validateUniqueNodeIDs(quest *domain.Quest, result *domain.ValidationResult) {
	seen := make(map[int]bool)
	for _, node := range quest.QuestNodes {
		if seen[node.NodeID] {
			result.AddNodeError(node.NodeID, "duplicate NodeID")
		}
		seen[node.NodeID] = true
	}
}

func (v *QuestValidatorService) validateNodeConnections(quest *domain.Quest, result *domain.ValidationResult) {
	nodeIDs := make(map[int]bool)
	for _, node := range quest.QuestNodes {
		nodeIDs[node.NodeID] = true
	}

	// Check that all NextNodes references exist
	for _, node := range quest.QuestNodes {
		for _, nextID := range node.NextNodes {
			if !nodeIDs[nextID] {
				result.AddNodeError(node.NodeID, "NextNodes references non-existent NodeID")
			}
		}
		// Also check dialog options
		for _, opt := range node.Options {
			for _, nextID := range opt.NextNodes {
				if !nodeIDs[nextID] {
					result.AddNodeError(node.NodeID, "dialog option NextNodes references non-existent NodeID")
				}
			}
		}
	}

	// Check that non-EntryPoint nodes have incoming connections
	hasIncoming := make(map[int]bool)
	for _, node := range quest.QuestNodes {
		for _, nextID := range node.NextNodes {
			hasIncoming[nextID] = true
		}
		for _, opt := range node.Options {
			for _, nextID := range opt.NextNodes {
				hasIncoming[nextID] = true
			}
		}
	}

	for _, node := range quest.QuestNodes {
		if node.NodeType != "EntryPoint" && !hasIncoming[node.NodeID] {
			result.AddNodeError(node.NodeID, "non-EntryPoint node has no incoming connections")
		}
	}
}

func (v *QuestValidatorService) validateEntryPoints(quest *domain.Quest, result *domain.ValidationResult) {
	hasEntryPoint := false
	for _, node := range quest.QuestNodes {
		if node.NodeType == "EntryPoint" {
			hasEntryPoint = true
			break
		}
	}
	if !hasEntryPoint {
		result.AddGlobalError("quest must have at least one EntryPoint node")
	}
}

func (v *QuestValidatorService) validateTerminalNodes(quest *domain.Quest, result *domain.ValidationResult) {
	terminalActions := map[string]bool{
		"CompleteQuest": true,
		"FailQuest":     true,
		"DeclineQuest":  true,
		"PostponeQuest": true,
	}

	for _, node := range quest.QuestNodes {
		if node.NodeType != "Actions" {
			continue
		}

		terminalCount := 0
		for _, action := range node.Actions {
			if actionStr, ok := action.(string); ok {
				if terminalActions[actionStr] {
					terminalCount++
				}
			}
		}

		// Check: terminal action nodes should not have NextNodes
		if terminalCount > 0 && len(node.NextNodes) > 0 {
			result.AddNodeError(node.NodeID, "terminal action node should not have NextNodes")
		}

		// Check: at most one terminal action per node
		if terminalCount > 1 {
			result.AddNodeError(node.NodeID, "Actions node has more than one terminal action")
		}
	}
}

func (v *QuestValidatorService) validatePlayerDecisionDialogs(quest *domain.Quest, result *domain.ValidationResult) {
	for _, node := range quest.QuestNodes {
		if node.NodeType != "PlayerDecisionDialog" {
			continue
		}

		// PlayerDecisionDialog must not have top-level NextNodes
		if len(node.NextNodes) > 0 {
			result.AddNodeError(node.NodeID, "PlayerDecisionDialog must not have top-level NextNodes; use NextNodes in each option instead")
		}

		// Every option must have NextNodes
		for i, opt := range node.Options {
			if len(opt.NextNodes) == 0 {
				result.AddNodeError(node.NodeID, fmt.Sprintf("option %d must have NextNodes", i+1))
			}
		}
	}
}

func (v *QuestValidatorService) validateNoCycles(quest *domain.Quest, result *domain.ValidationResult) {
	// Build adjacency list
	adj := make(map[int][]int)
	for _, node := range quest.QuestNodes {
		adj[node.NodeID] = append(adj[node.NodeID], node.NextNodes...)
		for _, opt := range node.Options {
			adj[node.NodeID] = append(adj[node.NodeID], opt.NextNodes...)
		}
	}

	// DFS-based cycle detection
	white := 0 // not visited
	gray := 1  // visiting (in current path)
	black := 2 // visited

	color := make(map[int]int)
	for _, node := range quest.QuestNodes {
		color[node.NodeID] = white
	}

	var hasCycle bool
	var dfs func(nodeID int) bool
	dfs = func(nodeID int) bool {
		color[nodeID] = gray
		for _, nextID := range adj[nodeID] {
			if color[nextID] == gray {
				return true // Back edge found - cycle
			}
			if color[nextID] == white {
				if dfs(nextID) {
					return true
				}
			}
		}
		color[nodeID] = black
		return false
	}

	for _, node := range quest.QuestNodes {
		if color[node.NodeID] == white {
			if dfs(node.NodeID) {
				hasCycle = true
				break
			}
		}
	}

	if hasCycle {
		result.AddGlobalError("quest contains a cycle (loops are not allowed)")
	}
}

func (v *QuestValidatorService) validateReferences(quest *domain.Quest, result *domain.ValidationResult) {
	// Load reference data
	npcs, _ := v.refData.ListNPCs()
	items, _ := v.refData.ListItems()
	factions, _ := v.refData.ListFactions()
	resources, _ := v.refData.ListResources()

	npcIDs := make(map[string]bool)
	for _, npc := range npcs {
		npcIDs[npc.NPCID] = true
	}

	itemIDs := make(map[string]bool)
	for _, item := range items {
		itemIDs[item.ItemID] = true
	}

	factionIDs := make(map[string]bool)
	for _, faction := range factions {
		factionIDs[faction.FactionID] = true
	}

	resourceIDs := make(map[string]bool)
	for _, resource := range resources {
		resourceIDs[resource.ResourceID] = true
	}

	for _, node := range quest.QuestNodes {
		// Check conversation partners and speakers
		if node.ConversationPartner != "" && !npcIDs[node.ConversationPartner] {
			result.AddNodeError(node.NodeID, "unknown conversation partner: "+node.ConversationPartner)
		}
		if node.Speaker != "" && !npcIDs[node.Speaker] {
			result.AddNodeError(node.NodeID, "unknown speaker: "+node.Speaker)
		}

		// Check message speakers (allow "Player" as a valid speaker)
		for _, msg := range node.Messages {
			if msg.Speaker != "Player" && !npcIDs[msg.Speaker] {
				result.AddNodeError(node.NodeID, "unknown speaker in message: "+msg.Speaker)
			}
		}
	}
}
