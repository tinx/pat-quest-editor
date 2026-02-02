package main

// I18nString represents a localized string with multiple language versions.
type I18nString struct {
	EnUS string `yaml:"en-US"`
	DeDE string `yaml:"de-DE"`
}

// Quest represents a complete quest definition.
type Quest struct {
	QuestTypeVersion int         `yaml:"QuestTypeVersion"`
	QuestVersion     int         `yaml:"QuestVersion"`
	QuestID          string      `yaml:"QuestID"`
	QuestType        string      `yaml:"QuestType"`
	DisplayName      I18nString  `yaml:"DisplayName"`
	Repeatable       string      `yaml:"Repeatable"`
	QuestNodes       []QuestNode `yaml:"QuestNodes"`
}

// QuestNode represents a node in the quest state machine.
type QuestNode struct {
	NodeID              int                      `yaml:"NodeID"`
	NodeType            string                   `yaml:"NodeType"`
	NextNodes           []int                    `yaml:"NextNodes,omitempty"`
	NextNodesIfTrue     []int                    `yaml:"NextNodesIfTrue,omitempty"`
	NextNodesIfFalse    []int                    `yaml:"NextNodesIfFalse,omitempty"`
	Conditions          []map[string]interface{} `yaml:"Conditions,omitempty"`
	ConditionsRequired  string                   `yaml:"ConditionsRequired,omitempty"`
	ConversationPartner string                   `yaml:"ConversationPartner,omitempty"`
	Speaker             string                   `yaml:"Speaker,omitempty"`
	Text                *I18nString              `yaml:"Text,omitempty"`
	Options             []DialogOption           `yaml:"Options,omitempty"`
	Messages            []DialogMessage          `yaml:"Messages,omitempty"`
	Actions             []interface{}            `yaml:"Actions,omitempty"`
}

// DialogOption represents a player dialog choice.
type DialogOption struct {
	Text       I18nString               `yaml:"Text"`
	Conditions []map[string]interface{} `yaml:"Conditions,omitempty"`
	NextNodes  []int                    `yaml:"NextNodes,omitempty"`
}

// DialogMessage represents a message in a dialog sequence.
type DialogMessage struct {
	Speaker string     `yaml:"Speaker"`
	Text    I18nString `yaml:"Text"`
}

// NPC represents a non-player character.
type NPC struct {
	NPCID string `yaml:"NPCID"`
}

// Item represents an item type.
type Item struct {
	ItemID string `yaml:"ItemID"`
}

// Faction represents a faction.
type Faction struct {
	FactionID string `yaml:"FactionID"`
}

// Resource represents a resource type.
type Resource struct {
	ResourceID string `yaml:"ResourceID"`
}

// Object represents a world object.
type Object struct {
	ObjectID string `yaml:"ObjectID"`
}

// ReferenceData holds all reference data for validation.
type ReferenceData struct {
	NPCs      map[string]bool
	Items     map[string]bool
	Factions  map[string]bool
	Resources map[string]bool
	Objects   map[string]bool
}

// ValidationError represents a single validation issue.
type ValidationError struct {
	QuestID string
	NodeID  *int
	Message string
}
