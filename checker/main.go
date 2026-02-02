package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	questsPath := flag.String("quests", "./quests", "Path to quests directory")
	dataPath := flag.String("data", "./data", "Path to reference data directory")
	quiet := flag.Bool("quiet", false, "Only output errors, no summary")
	flag.Parse()

	exitCode := run(*questsPath, *dataPath, *quiet)
	os.Exit(exitCode)
}

func run(questsPath, dataPath string, quiet bool) int {
	// Load reference data
	refData, err := LoadReferenceData(dataPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return 2
	}

	// Load all quests
	quests, loadErrors := LoadQuests(questsPath)

	// Print load errors
	for _, err := range loadErrors {
		fmt.Printf("[LOAD ERROR]: %v\n", err)
	}

	if len(quests) == 0 && len(loadErrors) == 0 {
		if !quiet {
			fmt.Println("No quests found.")
		}
		return 0
	}

	// Run single-quest validation
	var singleErrors []ValidationError
	for _, quest := range quests {
		singleErrors = append(singleErrors, ValidateQuest(quest, refData)...)
	}

	// Run cross-quest validation
	crossErrors := ValidateCrossQuest(quests)

	// Print all errors
	allErrors := append(singleErrors, crossErrors...)
	for _, verr := range allErrors {
		fmt.Println(formatError(verr))
	}

	// Summary
	totalErrors := len(loadErrors) + len(allErrors)
	if !quiet {
		fmt.Println(strings.Repeat("-", 40))
		fmt.Printf("Checked %d quests, found %d issues.\n", len(quests), totalErrors)
	}

	if totalErrors > 0 {
		return 1
	}
	return 0
}

func formatError(err ValidationError) string {
	if err.QuestID != "" {
		if err.NodeID != nil {
			return fmt.Sprintf("[%s] Node %d: %s", err.QuestID, *err.NodeID, err.Message)
		}
		return fmt.Sprintf("[%s]: %s", err.QuestID, err.Message)
	}
	return fmt.Sprintf("[CROSS-QUEST]: %s", err.Message)
}
