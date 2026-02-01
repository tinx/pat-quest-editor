package domain

// I18nString represents a localized string with multiple language versions.
type I18nString struct {
	EnUS string `yaml:"en-US" json:"en-US"`
	DeDE string `yaml:"de-DE" json:"de-DE"`
}

// Quest represents a complete quest definition.
type Quest struct {
	QuestTypeVersion int          `yaml:"QuestTypeVersion" json:"QuestTypeVersion"`
	QuestVersion     int          `yaml:"QuestVersion" json:"QuestVersion"`
	QuestID          string       `yaml:"QuestID" json:"QuestID"`
	QuestType        string       `yaml:"QuestType" json:"QuestType"`
	DisplayName      I18nString   `yaml:"DisplayName" json:"DisplayName"`
	Repeatable       string       `yaml:"Repeatable" json:"Repeatable"`
	QuestNodes       []QuestNode  `yaml:"QuestNodes" json:"QuestNodes"`
}

// QuestNode represents a node in the quest state machine.
type QuestNode struct {
	NodeID              int                  `yaml:"NodeID" json:"NodeID"`
	NodeType            string               `yaml:"NodeType" json:"NodeType"`
	NextNodes           []int                `yaml:"NextNodes,omitempty" json:"NextNodes,omitempty"`
	NextNodesIfTrue     []int                `yaml:"NextNodesIfTrue,omitempty" json:"NextNodesIfTrue,omitempty"`
	NextNodesIfFalse    []int                `yaml:"NextNodesIfFalse,omitempty" json:"NextNodesIfFalse,omitempty"`
	Conditions          []Condition          `yaml:"Conditions,omitempty" json:"Conditions,omitempty"`
	ConditionsRequired  string               `yaml:"ConditionsRequired,omitempty" json:"ConditionsRequired,omitempty"`
	ConversationPartner string               `yaml:"ConversationPartner,omitempty" json:"ConversationPartner,omitempty"`
	Speaker             string               `yaml:"Speaker,omitempty" json:"Speaker,omitempty"`
	Text                *I18nString          `yaml:"Text,omitempty" json:"Text,omitempty"`
	Options             []DialogOption       `yaml:"Options,omitempty" json:"Options,omitempty"`
	Messages            []DialogMessage      `yaml:"Messages,omitempty" json:"Messages,omitempty"`
	Actions             []Action             `yaml:"Actions,omitempty" json:"Actions,omitempty"`
}

// Condition represents a condition that can be checked.
type Condition map[string]interface{}

// Action represents an action that can be executed.
type Action interface{}

// DialogOption represents a player dialog choice.
type DialogOption struct {
	Text       I18nString  `yaml:"Text" json:"Text"`
	Conditions []Condition `yaml:"Conditions,omitempty" json:"Conditions,omitempty"`
	NextNodes  []int       `yaml:"NextNodes,omitempty" json:"NextNodes,omitempty"`
}

// DialogMessage represents a message in a dialog sequence.
type DialogMessage struct {
	Speaker string     `yaml:"Speaker" json:"Speaker"`
	Text    I18nString `yaml:"Text" json:"Text"`
}

// QuestMetadata stores editor-specific data (not part of quest YAML).
type QuestMetadata struct {
	QuestID       string                 `json:"questId"`
	NodePositions map[int]NodePosition   `json:"nodePositions"`
}

// NodePosition stores the x,y coordinates of a node on the editor canvas.
type NodePosition struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}
