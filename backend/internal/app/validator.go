package app

import (
	"fmt"
	"log"

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
	v.validateDecisions(quest, result)
	v.validateConditionBranches(quest, result)
	v.validateNoCycles(quest, result)
	v.validateReferences(quest, result)
	v.validateNoUnreferencedNodes(quest, result)
	v.validateJournalAtFlowStart(quest, result)
	v.validateJournalAtFlowEnd(quest, result)

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

	// Helper to check a single NextNodes list for duplicates, self-references, and invalid references
	checkNextNodesList := func(nodeID int, list []int, listName string) {
		seen := make(map[int]bool)
		for _, nextID := range list {
			if !nodeIDs[nextID] {
				result.AddNodeError(nodeID, fmt.Sprintf("%s references non-existent NodeID", listName))
			}
			if nextID == nodeID {
				result.AddNodeError(nodeID, fmt.Sprintf("node references itself in %s", listName))
			}
			if seen[nextID] {
				result.AddNodeError(nodeID, fmt.Sprintf("duplicate edge to NodeID %d", nextID))
			}
			seen[nextID] = true
		}
	}

	// Check that all NextNodes references exist, no duplicates, and no self-references
	for _, node := range quest.QuestNodes {
		checkNextNodesList(node.NodeID, node.NextNodes, "NextNodes")
		checkNextNodesList(node.NodeID, node.NextNodesIfTrue, "NextNodesIfTrue")
		checkNextNodesList(node.NodeID, node.NextNodesIfFalse, "NextNodesIfFalse")

		// Check dialog options - each option's NextNodes is checked independently
		for i, opt := range node.Options {
			checkNextNodesList(node.NodeID, opt.NextNodes, fmt.Sprintf("option %d NextNodes", i+1))
		}
	}

	// Check that non-EntryPoint nodes have incoming connections
	hasIncoming := make(map[int]bool)
	for _, node := range quest.QuestNodes {
		for _, nextID := range node.NextNodes {
			hasIncoming[nextID] = true
		}
		for _, nextID := range node.NextNodesIfTrue {
			hasIncoming[nextID] = true
		}
		for _, nextID := range node.NextNodesIfFalse {
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

		isTerminal := terminalCount > 0
		hasOutgoingEdges := len(node.NextNodes) > 0

		// Check: terminal action nodes should not have NextNodes
		if isTerminal && hasOutgoingEdges {
			result.AddNodeError(node.NodeID, "terminal action node should not have NextNodes")
		}

		// Check: non-terminal action nodes must have NextNodes
		if !isTerminal && !hasOutgoingEdges {
			result.AddNodeError(node.NodeID, "non-terminal Actions node must have NextNodes (quest flow ends with unspecified behaviour)")
		}

		// Check: at most one terminal action per node
		if terminalCount > 1 {
			result.AddNodeError(node.NodeID, "Actions node has more than one terminal action")
		}
	}
}

func (v *QuestValidatorService) validateDecisions(quest *domain.Quest, result *domain.ValidationResult) {
	for _, node := range quest.QuestNodes {
		if node.NodeType != "Decision" {
			continue
		}

		// Decision must not have top-level NextNodes
		if len(node.NextNodes) > 0 {
			result.AddNodeError(node.NodeID, "Decision must not have top-level NextNodes; use NextNodes in each option instead")
		}

		// Every option must have NextNodes
		for i, opt := range node.Options {
			if len(opt.NextNodes) == 0 {
				result.AddNodeError(node.NodeID, fmt.Sprintf("option %d must have NextNodes", i+1))
			}
		}
	}
}

func (v *QuestValidatorService) validateConditionBranches(quest *domain.Quest, result *domain.ValidationResult) {
	for _, node := range quest.QuestNodes {
		if node.NodeType != "ConditionBranch" {
			continue
		}

		// ConditionBranch must not have top-level NextNodes
		if len(node.NextNodes) > 0 {
			result.AddNodeError(node.NodeID, "ConditionBranch must not have top-level NextNodes; use NextNodesIfTrue and NextNodesIfFalse instead")
		}

		// ConditionBranch must have at least one condition
		if len(node.Conditions) == 0 {
			result.AddNodeError(node.NodeID, "ConditionBranch must have at least one condition")
		}

		// At least one of NextNodesIfTrue or NextNodesIfFalse must be non-empty
		if len(node.NextNodesIfTrue) == 0 && len(node.NextNodesIfFalse) == 0 {
			result.AddNodeError(node.NodeID, "ConditionBranch must have at least one of NextNodesIfTrue or NextNodesIfFalse")
		}
	}
}

func (v *QuestValidatorService) validateNoCycles(quest *domain.Quest, result *domain.ValidationResult) {
	// Build adjacency list
	adj := make(map[int][]int)
	for _, node := range quest.QuestNodes {
		adj[node.NodeID] = append(adj[node.NodeID], node.NextNodes...)
		adj[node.NodeID] = append(adj[node.NodeID], node.NextNodesIfTrue...)
		adj[node.NodeID] = append(adj[node.NodeID], node.NextNodesIfFalse...)
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
	npcs, err := v.refData.ListNPCs()
	if err != nil {
		log.Printf("Warning: failed to load NPCs for validation: %v", err)
	}
	items, err := v.refData.ListItems()
	if err != nil {
		log.Printf("Warning: failed to load items for validation: %v", err)
	}
	factions, err := v.refData.ListFactions()
	if err != nil {
		log.Printf("Warning: failed to load factions for validation: %v", err)
	}
	resources, err := v.refData.ListResources()
	if err != nil {
		log.Printf("Warning: failed to load resources for validation: %v", err)
	}
	objects, err := v.refData.ListObjects()
	if err != nil {
		log.Printf("Warning: failed to load objects for validation: %v", err)
	}

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

	objectIDs := make(map[string]bool)
	for _, obj := range objects {
		objectIDs[obj.ObjectID] = true
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

		// Check conditions
		for _, cond := range node.Conditions {
			// Check ResourceAvailability conditions
			if ra, ok := cond["ResourceAvailability"].(map[string]interface{}); ok {
				if resource, ok := ra["Resource"].(string); ok && resource != "" {
					if !resourceIDs[resource] {
						result.AddNodeError(node.NodeID, "unknown resource in ResourceAvailability: "+resource)
					}
				}
			}

			// Check ItemUsedOnObject conditions
			if iuo, ok := cond["ItemUsedOnObject"].(map[string]interface{}); ok {
				if item, ok := iuo["Item"].(string); ok && item != "" {
					if !itemIDs[item] {
						result.AddNodeError(node.NodeID, "unknown item in ItemUsedOnObject: "+item)
					}
				}
				if obj, ok := iuo["Object"].(string); ok && obj != "" {
					if !objectIDs[obj] {
						result.AddNodeError(node.NodeID, "unknown object in ItemUsedOnObject: "+obj)
					}
				}
			}

			// Check ItemUsedOnNPC conditions
			if iun, ok := cond["ItemUsedOnNPC"].(map[string]interface{}); ok {
				if item, ok := iun["Item"].(string); ok && item != "" {
					if !itemIDs[item] {
						result.AddNodeError(node.NodeID, "unknown item in ItemUsedOnNPC: "+item)
					}
				}
				if npc, ok := iun["NPC"].(string); ok && npc != "" {
					if !npcIDs[npc] {
						result.AddNodeError(node.NodeID, "unknown NPC in ItemUsedOnNPC: "+npc)
					}
				}
			}
		}
	}
}

func (v *QuestValidatorService) validateNoUnreferencedNodes(quest *domain.Quest, result *domain.ValidationResult) {
	// Collect all referenced NodeIDs (nodes that appear in any NextNodes list)
	referenced := make(map[int]bool)
	for _, node := range quest.QuestNodes {
		for _, nextID := range node.NextNodes {
			referenced[nextID] = true
		}
		for _, nextID := range node.NextNodesIfTrue {
			referenced[nextID] = true
		}
		for _, nextID := range node.NextNodesIfFalse {
			referenced[nextID] = true
		}
		for _, opt := range node.Options {
			for _, nextID := range opt.NextNodes {
				referenced[nextID] = true
			}
		}
	}

	// Check that every non-EntryPoint node is referenced
	for _, node := range quest.QuestNodes {
		if node.NodeType == "EntryPoint" {
			continue
		}
		if !referenced[node.NodeID] {
			result.AddNodeError(node.NodeID, "NodeID is never referenced by any other node")
		}
	}
}

// getOutgoingEdges returns all nodes that a given node can transition to.
func getOutgoingEdges(node *domain.QuestNode) []int {
	edges := make([]int, 0)
	edges = append(edges, node.NextNodes...)
	edges = append(edges, node.NextNodesIfTrue...)
	edges = append(edges, node.NextNodesIfFalse...)
	for _, opt := range node.Options {
		edges = append(edges, opt.NextNodes...)
	}
	return edges
}

// nodeHasAction checks if an Actions node contains a specific action type.
func nodeHasAction(node *domain.QuestNode, actionName string) bool {
	for _, action := range node.Actions {
		if actionStr, ok := action.(string); ok && actionStr == actionName {
			return true
		}
		if actionMap, ok := action.(map[string]interface{}); ok {
			if _, exists := actionMap[actionName]; exists {
				return true
			}
		}
	}
	return false
}

// isTerminalNode checks if a node is a terminal Actions node.
func isTerminalNode(node *domain.QuestNode) bool {
	if node.NodeType != "Actions" {
		return false
	}
	return nodeHasAction(node, "CompleteQuest") ||
		nodeHasAction(node, "FailQuest") ||
		nodeHasAction(node, "DeclineQuest")
}

func (v *QuestValidatorService) validateJournalAtFlowStart(quest *domain.Quest, result *domain.ValidationResult) {
	// Build node lookup map
	nodeByID := make(map[int]*domain.QuestNode)
	for i := range quest.QuestNodes {
		nodeByID[quest.QuestNodes[i].NodeID] = &quest.QuestNodes[i]
	}

	// For each EntryPoint, find the first Actions node in the flow
	for _, node := range quest.QuestNodes {
		if node.NodeType != "EntryPoint" {
			continue
		}

		// BFS to find first Actions node
		visited := make(map[int]bool)
		queue := getOutgoingEdges(&node)
		var firstActionsNode *domain.QuestNode

		for len(queue) > 0 && firstActionsNode == nil {
			currentID := queue[0]
			queue = queue[1:]

			if visited[currentID] {
				continue
			}
			visited[currentID] = true

			current, exists := nodeByID[currentID]
			if !exists {
				continue
			}

			if current.NodeType == "Actions" {
				firstActionsNode = current
				break
			}

			// Add outgoing edges to queue
			queue = append(queue, getOutgoingEdges(current)...)
		}

		if firstActionsNode == nil {
			result.AddNodeError(node.NodeID, "EntryPoint flow has no Actions node")
			continue
		}

		// Check that first Actions node has both JournalEntry and QuestStageDescription
		hasJournalEntry := nodeHasAction(firstActionsNode, "JournalEntry")
		hasQuestStageDescription := nodeHasAction(firstActionsNode, "QuestStageDescription")

		if !hasJournalEntry {
			result.AddNodeError(firstActionsNode.NodeID, "first Actions node in flow must have JournalEntry action")
		}
		if !hasQuestStageDescription {
			result.AddNodeError(firstActionsNode.NodeID, "first Actions node in flow must have QuestStageDescription action")
		}
	}
}

func (v *QuestValidatorService) validateJournalAtFlowEnd(quest *domain.Quest, result *domain.ValidationResult) {
	// Build node lookup map and incoming edges map
	nodeByID := make(map[int]*domain.QuestNode)
	incomingFrom := make(map[int][]int) // nodeID -> list of nodes that point to it
	for i := range quest.QuestNodes {
		node := &quest.QuestNodes[i]
		nodeByID[node.NodeID] = node
		for _, nextID := range getOutgoingEdges(node) {
			incomingFrom[nextID] = append(incomingFrom[nextID], node.NodeID)
		}
	}

	// Find all terminal nodes and check their Actions chains
	for _, node := range quest.QuestNodes {
		if !isTerminalNode(&node) {
			continue
		}

		// Trace backwards to find the Actions chain
		// The chain includes the terminal node and all consecutive Actions nodes before it
		actionsChain := []*domain.QuestNode{&node}
		visited := make(map[int]bool)
		visited[node.NodeID] = true

		// BFS backwards, but only follow Actions nodes
		toCheck := []int{node.NodeID}
		for len(toCheck) > 0 {
			currentID := toCheck[0]
			toCheck = toCheck[1:]

			for _, prevID := range incomingFrom[currentID] {
				if visited[prevID] {
					continue
				}
				prevNode, exists := nodeByID[prevID]
				if !exists {
					continue
				}
				// Only include Actions nodes in the chain
				if prevNode.NodeType == "Actions" {
					visited[prevID] = true
					actionsChain = append(actionsChain, prevNode)
					toCheck = append(toCheck, prevID)
				}
			}
		}

		// Check if any node in the chain has JournalEntry
		hasJournalEntry := false
		for _, chainNode := range actionsChain {
			if nodeHasAction(chainNode, "JournalEntry") {
				hasJournalEntry = true
				break
			}
		}

		if !hasJournalEntry {
			result.AddNodeError(node.NodeID, "terminal Actions chain must contain a JournalEntry action")
		}
	}
}
