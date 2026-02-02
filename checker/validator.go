package main

import "fmt"

// ValidateQuest validates a single quest and returns any errors found.
func ValidateQuest(quest *Quest, refData *ReferenceData) []ValidationError {
	var errors []ValidationError

	errors = append(errors, validateUniqueNodeIDs(quest)...)
	errors = append(errors, validateNodeConnections(quest)...)
	errors = append(errors, validateEntryPoints(quest)...)
	errors = append(errors, validateTerminalNodes(quest)...)
	errors = append(errors, validateOutgoingEdges(quest)...)
	errors = append(errors, validateNoCycles(quest)...)
	errors = append(errors, validateReferences(quest, refData)...)

	return errors
}

func validateUniqueNodeIDs(quest *Quest) []ValidationError {
	var errors []ValidationError
	seen := make(map[int]bool)

	for _, node := range quest.QuestNodes {
		if seen[node.NodeID] {
			errors = append(errors, ValidationError{
				QuestID: quest.QuestID,
				NodeID:  intPtr(node.NodeID),
				Message: "duplicate NodeID",
			})
		}
		seen[node.NodeID] = true
	}

	return errors
}

func validateNodeConnections(quest *Quest) []ValidationError {
	var errors []ValidationError

	nodeIDs := make(map[int]bool)
	for _, node := range quest.QuestNodes {
		nodeIDs[node.NodeID] = true
	}

	for _, node := range quest.QuestNodes {
		// Check NextNodes
		seen := make(map[int]bool)
		for _, nextID := range node.NextNodes {
			if nextID == node.NodeID {
				errors = append(errors, ValidationError{
					QuestID: quest.QuestID,
					NodeID:  intPtr(node.NodeID),
					Message: "node references itself in NextNodes",
				})
			}
			if seen[nextID] {
				errors = append(errors, ValidationError{
					QuestID: quest.QuestID,
					NodeID:  intPtr(node.NodeID),
					Message: fmt.Sprintf("duplicate edge to node %d", nextID),
				})
			}
			seen[nextID] = true
			if !nodeIDs[nextID] {
				errors = append(errors, ValidationError{
					QuestID: quest.QuestID,
					NodeID:  intPtr(node.NodeID),
					Message: fmt.Sprintf("NextNodes references non-existent node %d", nextID),
				})
			}
		}

		// Check NextNodesIfTrue
		for _, nextID := range node.NextNodesIfTrue {
			if nextID == node.NodeID {
				errors = append(errors, ValidationError{
					QuestID: quest.QuestID,
					NodeID:  intPtr(node.NodeID),
					Message: "node references itself in NextNodesIfTrue",
				})
			}
			if !nodeIDs[nextID] {
				errors = append(errors, ValidationError{
					QuestID: quest.QuestID,
					NodeID:  intPtr(node.NodeID),
					Message: fmt.Sprintf("NextNodesIfTrue references non-existent node %d", nextID),
				})
			}
		}

		// Check NextNodesIfFalse
		for _, nextID := range node.NextNodesIfFalse {
			if nextID == node.NodeID {
				errors = append(errors, ValidationError{
					QuestID: quest.QuestID,
					NodeID:  intPtr(node.NodeID),
					Message: "node references itself in NextNodesIfFalse",
				})
			}
			if !nodeIDs[nextID] {
				errors = append(errors, ValidationError{
					QuestID: quest.QuestID,
					NodeID:  intPtr(node.NodeID),
					Message: fmt.Sprintf("NextNodesIfFalse references non-existent node %d", nextID),
				})
			}
		}

		// Check dialog option NextNodes
		for _, opt := range node.Options {
			for _, nextID := range opt.NextNodes {
				if nextID == node.NodeID {
					errors = append(errors, ValidationError{
						QuestID: quest.QuestID,
						NodeID:  intPtr(node.NodeID),
						Message: "dialog option references itself",
					})
				}
				if !nodeIDs[nextID] {
					errors = append(errors, ValidationError{
						QuestID: quest.QuestID,
						NodeID:  intPtr(node.NodeID),
						Message: fmt.Sprintf("dialog option references non-existent node %d", nextID),
					})
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
			errors = append(errors, ValidationError{
				QuestID: quest.QuestID,
				NodeID:  intPtr(node.NodeID),
				Message: "non-EntryPoint node has no incoming connections",
			})
		}
	}

	return errors
}

func validateEntryPoints(quest *Quest) []ValidationError {
	var errors []ValidationError

	hasEntryPoint := false
	for _, node := range quest.QuestNodes {
		if node.NodeType == "EntryPoint" {
			hasEntryPoint = true
			break
		}
	}

	if !hasEntryPoint {
		errors = append(errors, ValidationError{
			QuestID: quest.QuestID,
			Message: "quest must have at least one EntryPoint node",
		})
	}

	return errors
}

func validateTerminalNodes(quest *Quest) []ValidationError {
	var errors []ValidationError

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

		if terminalCount > 0 && len(node.NextNodes) > 0 {
			errors = append(errors, ValidationError{
				QuestID: quest.QuestID,
				NodeID:  intPtr(node.NodeID),
				Message: "terminal action node must not have NextNodes",
			})
		}

		if terminalCount > 1 {
			errors = append(errors, ValidationError{
				QuestID: quest.QuestID,
				NodeID:  intPtr(node.NodeID),
				Message: "Actions node has more than one terminal action",
			})
		}
	}

	return errors
}

func validateOutgoingEdges(quest *Quest) []ValidationError {
	var errors []ValidationError

	terminalActions := map[string]bool{
		"CompleteQuest": true,
		"FailQuest":     true,
		"DeclineQuest":  true,
	}

	for _, node := range quest.QuestNodes {
		outgoingCount := len(node.NextNodes) + len(node.NextNodesIfTrue) + len(node.NextNodesIfFalse)
		for _, opt := range node.Options {
			outgoingCount += len(opt.NextNodes)
		}

		// Check if this is a terminal node
		isTerminal := false
		if node.NodeType == "Actions" {
			for _, action := range node.Actions {
				if actionStr, ok := action.(string); ok {
					if terminalActions[actionStr] {
						isTerminal = true
						break
					}
				}
			}
		}

		// Non-terminal nodes must have outgoing edges (except Decision which uses Options)
		if !isTerminal && outgoingCount == 0 && node.NodeType != "Decision" {
			errors = append(errors, ValidationError{
				QuestID: quest.QuestID,
				NodeID:  intPtr(node.NodeID),
				Message: "non-terminal node has no outgoing edges",
			})
		}
	}

	return errors
}

func validateNoCycles(quest *Quest) []ValidationError {
	var errors []ValidationError

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
	const (
		white = 0 // not visited
		gray  = 1 // visiting (in current path)
		black = 2 // visited
	)

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
				return true
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
		errors = append(errors, ValidationError{
			QuestID: quest.QuestID,
			Message: "quest contains a cycle",
		})
	}

	return errors
}

func validateReferences(quest *Quest, refData *ReferenceData) []ValidationError {
	var errors []ValidationError

	for _, node := range quest.QuestNodes {
		// Check conversation partners and speakers
		if node.ConversationPartner != "" && !refData.NPCs[node.ConversationPartner] {
			errors = append(errors, ValidationError{
				QuestID: quest.QuestID,
				NodeID:  intPtr(node.NodeID),
				Message: fmt.Sprintf("unknown NPC: %s", node.ConversationPartner),
			})
		}
		if node.Speaker != "" && !refData.NPCs[node.Speaker] {
			errors = append(errors, ValidationError{
				QuestID: quest.QuestID,
				NodeID:  intPtr(node.NodeID),
				Message: fmt.Sprintf("unknown speaker: %s", node.Speaker),
			})
		}

		// Check message speakers
		for _, msg := range node.Messages {
			if msg.Speaker != "Player" && !refData.NPCs[msg.Speaker] {
				errors = append(errors, ValidationError{
					QuestID: quest.QuestID,
					NodeID:  intPtr(node.NodeID),
					Message: fmt.Sprintf("unknown speaker in message: %s", msg.Speaker),
				})
			}
		}

		// Check conditions
		errors = append(errors, validateConditionReferences(quest.QuestID, node.NodeID, node.Conditions, refData)...)

		// Check dialog option conditions
		for _, opt := range node.Options {
			errors = append(errors, validateConditionReferences(quest.QuestID, node.NodeID, opt.Conditions, refData)...)
		}

		// Check action references
		errors = append(errors, validateActionReferences(quest.QuestID, node.NodeID, node.Actions, refData)...)
	}

	return errors
}

func validateConditionReferences(questID string, nodeID int, conditions []map[string]interface{}, refData *ReferenceData) []ValidationError {
	var errors []ValidationError

	for _, cond := range conditions {
		// ResourceAvailability
		if ra, ok := cond["ResourceAvailability"].(map[string]interface{}); ok {
			if resource, ok := ra["Resource"].(string); ok && resource != "" {
				if !refData.Resources[resource] {
					errors = append(errors, ValidationError{
						QuestID: questID,
						NodeID:  intPtr(nodeID),
						Message: fmt.Sprintf("unknown resource: %s", resource),
					})
				}
			}
		}

		// ItemUsedOnObject
		if iuo, ok := cond["ItemUsedOnObject"].(map[string]interface{}); ok {
			if item, ok := iuo["Item"].(string); ok && item != "" {
				if !refData.Items[item] {
					errors = append(errors, ValidationError{
						QuestID: questID,
						NodeID:  intPtr(nodeID),
						Message: fmt.Sprintf("unknown item: %s", item),
					})
				}
			}
			if obj, ok := iuo["Object"].(string); ok && obj != "" {
				if !refData.Objects[obj] {
					errors = append(errors, ValidationError{
						QuestID: questID,
						NodeID:  intPtr(nodeID),
						Message: fmt.Sprintf("unknown object: %s", obj),
					})
				}
			}
		}

		// ItemUsedOnNPC
		if iun, ok := cond["ItemUsedOnNPC"].(map[string]interface{}); ok {
			if item, ok := iun["Item"].(string); ok && item != "" {
				if !refData.Items[item] {
					errors = append(errors, ValidationError{
						QuestID: questID,
						NodeID:  intPtr(nodeID),
						Message: fmt.Sprintf("unknown item: %s", item),
					})
				}
			}
			if npc, ok := iun["NPC"].(string); ok && npc != "" {
				if !refData.NPCs[npc] {
					errors = append(errors, ValidationError{
						QuestID: questID,
						NodeID:  intPtr(nodeID),
						Message: fmt.Sprintf("unknown NPC: %s", npc),
					})
				}
			}
		}

		// FactionStanding
		if fs, ok := cond["FactionStanding"].(map[string]interface{}); ok {
			if faction, ok := fs["Faction"].(string); ok && faction != "" {
				if !refData.Factions[faction] {
					errors = append(errors, ValidationError{
						QuestID: questID,
						NodeID:  intPtr(nodeID),
						Message: fmt.Sprintf("unknown faction: %s", faction),
					})
				}
			}
		}

		// Inventory
		if inv, ok := cond["Inventory"].([]interface{}); ok {
			for _, item := range inv {
				if itemMap, ok := item.(map[string]interface{}); ok {
					if itemType, ok := itemMap["Type"].(string); ok && itemType != "" {
						if !refData.Items[itemType] {
							errors = append(errors, ValidationError{
								QuestID: questID,
								NodeID:  intPtr(nodeID),
								Message: fmt.Sprintf("unknown item type: %s", itemType),
							})
						}
					}
				}
			}
		}

		// ItemLost
		if itemLost, ok := cond["ItemLost"].(string); ok && itemLost != "" {
			if !refData.Items[itemLost] {
				errors = append(errors, ValidationError{
					QuestID: questID,
					NodeID:  intPtr(nodeID),
					Message: fmt.Sprintf("unknown item: %s", itemLost),
				})
			}
		}
	}

	return errors
}

func validateActionReferences(questID string, nodeID int, actions []interface{}, refData *ReferenceData) []ValidationError {
	var errors []ValidationError

	for _, action := range actions {
		actionMap, ok := action.(map[string]interface{})
		if !ok {
			continue
		}

		// ItemsGained
		if ig, ok := actionMap["ItemsGained"].([]interface{}); ok {
			for _, item := range ig {
				if itemMap, ok := item.(map[string]interface{}); ok {
					if itemType, ok := itemMap["Type"].(string); ok && itemType != "" {
						if !refData.Items[itemType] {
							errors = append(errors, ValidationError{
								QuestID: questID,
								NodeID:  intPtr(nodeID),
								Message: fmt.Sprintf("unknown item type in ItemsGained: %s", itemType),
							})
						}
					}
				}
			}
		}

		// ItemsLost
		if il, ok := actionMap["ItemsLost"].([]interface{}); ok {
			for _, item := range il {
				if itemMap, ok := item.(map[string]interface{}); ok {
					if itemType, ok := itemMap["Type"].(string); ok && itemType != "" {
						if !refData.Items[itemType] {
							errors = append(errors, ValidationError{
								QuestID: questID,
								NodeID:  intPtr(nodeID),
								Message: fmt.Sprintf("unknown item type in ItemsLost: %s", itemType),
							})
						}
					}
				}
			}
		}

		// FactionStanding
		if fs, ok := actionMap["FactionStanding"].(map[string]interface{}); ok {
			if faction, ok := fs["Faction"].(string); ok && faction != "" {
				if !refData.Factions[faction] {
					errors = append(errors, ValidationError{
						QuestID: questID,
						NodeID:  intPtr(nodeID),
						Message: fmt.Sprintf("unknown faction in FactionStanding: %s", faction),
					})
				}
			}
		}
	}

	return errors
}

func intPtr(i int) *int {
	return &i
}
