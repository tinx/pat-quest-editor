package main

import (
	"fmt"
	"strings"
)

// ValidateCrossQuest validates rules that span multiple quests.
func ValidateCrossQuest(quests []*Quest) []ValidationError {
	var errors []ValidationError

	questIDs := buildQuestIDSet(quests)

	errors = append(errors, validateUniqueQuestIDs(quests)...)
	errors = append(errors, validateUniqueDisplayNames(quests)...)
	errors = append(errors, validateUniqueQuestStageDescriptions(quests)...)
	errors = append(errors, validateQuestReferences(quests, questIDs)...)

	return errors
}

func buildQuestIDSet(quests []*Quest) map[string]bool {
	ids := make(map[string]bool)
	for _, q := range quests {
		ids[q.QuestID] = true
	}
	return ids
}

func validateUniqueQuestIDs(quests []*Quest) []ValidationError {
	var errors []ValidationError

	seen := make(map[string]int)
	for _, q := range quests {
		seen[q.QuestID]++
	}

	for questID, count := range seen {
		if count > 1 {
			errors = append(errors, ValidationError{
				Message: fmt.Sprintf("duplicate QuestID %q found %d times", questID, count),
			})
		}
	}

	return errors
}

func validateUniqueDisplayNames(quests []*Quest) []ValidationError {
	var errors []ValidationError

	// Check en-US
	seenEnUS := make(map[string][]string)
	for _, q := range quests {
		if q.DisplayName.EnUS != "" {
			seenEnUS[q.DisplayName.EnUS] = append(seenEnUS[q.DisplayName.EnUS], q.QuestID)
		}
	}
	for name, questIDs := range seenEnUS {
		if len(questIDs) > 1 {
			errors = append(errors, ValidationError{
				Message: fmt.Sprintf("duplicate DisplayName %q (en-US) in quests: %s", name, strings.Join(questIDs, ", ")),
			})
		}
	}

	// Check de-DE
	seenDeDE := make(map[string][]string)
	for _, q := range quests {
		if q.DisplayName.DeDE != "" {
			seenDeDE[q.DisplayName.DeDE] = append(seenDeDE[q.DisplayName.DeDE], q.QuestID)
		}
	}
	for name, questIDs := range seenDeDE {
		if len(questIDs) > 1 {
			errors = append(errors, ValidationError{
				Message: fmt.Sprintf("duplicate DisplayName %q (de-DE) in quests: %s", name, strings.Join(questIDs, ", ")),
			})
		}
	}

	return errors
}

func validateUniqueQuestStageDescriptions(quests []*Quest) []ValidationError {
	var errors []ValidationError

	seenEnUS := make(map[string][]string)
	seenDeDE := make(map[string][]string)

	for _, q := range quests {
		for _, node := range q.QuestNodes {
			for _, action := range node.Actions {
				if actionMap, ok := action.(map[string]interface{}); ok {
					if qsd, ok := actionMap["QuestStageDescription"].(map[string]interface{}); ok {
						if enUS, ok := qsd["en-US"].(string); ok && enUS != "" {
							seenEnUS[enUS] = appendUnique(seenEnUS[enUS], q.QuestID)
						}
						if deDE, ok := qsd["de-DE"].(string); ok && deDE != "" {
							seenDeDE[deDE] = appendUnique(seenDeDE[deDE], q.QuestID)
						}
					}
				}
			}
		}
	}

	for desc, questIDs := range seenEnUS {
		if len(questIDs) > 1 {
			errors = append(errors, ValidationError{
				Message: fmt.Sprintf("duplicate QuestStageDescription %q (en-US) in quests: %s", desc, strings.Join(questIDs, ", ")),
			})
		}
	}

	for desc, questIDs := range seenDeDE {
		if len(questIDs) > 1 {
			errors = append(errors, ValidationError{
				Message: fmt.Sprintf("duplicate QuestStageDescription %q (de-DE) in quests: %s", desc, strings.Join(questIDs, ", ")),
			})
		}
	}

	return errors
}

func validateQuestReferences(quests []*Quest, questIDs map[string]bool) []ValidationError {
	var errors []ValidationError

	for _, q := range quests {
		for _, node := range q.QuestNodes {
			// Check conditions
			for _, cond := range node.Conditions {
				if qc, ok := cond["QuestCompleted"].(string); ok && qc != "" {
					if !questIDs[qc] {
						errors = append(errors, ValidationError{
							QuestID: q.QuestID,
							NodeID:  intPtr(node.NodeID),
							Message: fmt.Sprintf("QuestCompleted references non-existent quest %q", qc),
						})
					}
				}
			}

			// Check dialog option conditions
			for _, opt := range node.Options {
				for _, cond := range opt.Conditions {
					if qc, ok := cond["QuestCompleted"].(string); ok && qc != "" {
						if !questIDs[qc] {
							errors = append(errors, ValidationError{
								QuestID: q.QuestID,
								NodeID:  intPtr(node.NodeID),
								Message: fmt.Sprintf("QuestCompleted references non-existent quest %q", qc),
							})
						}
					}
				}
			}
		}
	}

	return errors
}

func appendUnique(slice []string, s string) []string {
	for _, item := range slice {
		if item == s {
			return slice
		}
	}
	return append(slice, s)
}
