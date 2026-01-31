package domain

// Item represents an item type in the game.
type Item struct {
	ItemID      string     `yaml:"ItemID" json:"ItemID"`
	DisplayName I18nString `yaml:"DisplayName" json:"DisplayName"`
	Description I18nString `yaml:"Description,omitempty" json:"Description,omitempty"`
	Stackable   bool       `yaml:"Stackable,omitempty" json:"Stackable,omitempty"`
	MaxStack    int        `yaml:"MaxStack,omitempty" json:"MaxStack,omitempty"`
	Category    string     `yaml:"Category" json:"Category"`
}

// Faction represents a faction in the game.
type Faction struct {
	FactionID       string     `yaml:"FactionID" json:"FactionID"`
	DisplayName     I18nString `yaml:"DisplayName" json:"DisplayName"`
	Description     I18nString `yaml:"Description,omitempty" json:"Description,omitempty"`
	FactionType     string     `yaml:"FactionType" json:"FactionType"`
	MaxLevel        int        `yaml:"MaxLevel,omitempty" json:"MaxLevel,omitempty"`
	InitialStanding int        `yaml:"InitialStanding,omitempty" json:"InitialStanding,omitempty"`
}

// Resource represents a world resource in the game.
type Resource struct {
	ResourceID  string     `yaml:"ResourceID" json:"ResourceID"`
	DisplayName I18nString `yaml:"DisplayName" json:"DisplayName"`
	Description I18nString `yaml:"Description,omitempty" json:"Description,omitempty"`
	Category    string     `yaml:"Category" json:"Category"`
}

// NPC represents an NPC (conversation partner/speaker) in the game.
type NPC struct {
	NPCID       string      `yaml:"NPCID" json:"NPCID"`
	DisplayName I18nString  `yaml:"DisplayName" json:"DisplayName"`
	Title       *I18nString `yaml:"Title,omitempty" json:"Title,omitempty"`
	Description *I18nString `yaml:"Description,omitempty" json:"Description,omitempty"`
	Location    string      `yaml:"Location,omitempty" json:"Location,omitempty"`
	FactionID   string      `yaml:"FactionID,omitempty" json:"FactionID,omitempty"`
}
