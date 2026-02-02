# Quest Validation

This document explains how to validate the quests in the "Potions and Tinctures"
game.

## Schema

Quests are stored as YAML files. There is a JSON schema file called
"schemas/quest.json" explaining the valid strucutre of the YAML files.

## Additional Quest Logic

Not all conditions are modelled in the schema, or can be. Let's look at
the other rules and limitations that apply to quest files.

### Node Connections

 - Quests begin at one or more EntryPoint nodes. That's why EntryPoint
   nodes must not have any edges leading towards them. They are always
   the first node in any flow.
 - All node types can define where the control flow continues. They use
   properties like NextNodes, NextNodesIfTrue or NextNodesIfFalse to
   do so. However, since quests are directional acyclic graphs, it is
   not allowed to create loops. No node may direct the flow to another
   node that might then, eventually, lead back to itself.
 - Every node other than EntryPoint nodes must have at least one edge
   leading into it. Otherwise it would not be reachable.
 - Nodes must never reference the same next node more than once. No duplicate
   edges.
 - Nodes must never reference a non-existent other node as a next node.
 - Nodes must never reference themselves as a next node.

### Terminal Actions

 - Terminal actions are actions that end the quest in one way or another.
   There are three terminal actions: CompleteQuest, FailQuest and DeclineQuest.
 - Action nodes containing a terminal action are called terminal nodes.
 - Terminal nodes are the only nodes without outgoing edges.
 - Terminal nodes must not have outgoing edges.
 - Non-terminal nodes must have at least one outgoing edge. Otherwise the
   quest flow would end with unspecified behaviour.
 - Action nodes must never contain more than one terminal action.

### Journal Mechanics

 - The first Actions node in every quest flow must have a JournalEntry
   action and a QuestStageDescription action.
 - The last nodes in every quest flow must be one or more Actions nodes in
   direct succession, and at least one of them must have a JournalEntry action.

### NodeIDs

 - Every node in a quest must have a unique ID.
 - Every NodeID must be referenced at least once within a quest, except
   for NodeIDs of EntryPoint nodes.

## Higher Scope Concerns

### External References

 - Quests may reference NPCs. Quests must never reference a non-existing NPC.
 - Quests may reference world objects. Quests must never reference a
   non-existing world object type.
 - Quests may reference items. Quests must never reference a non-existing
   item type.
 - Quests may reference factions. Quests must never reference a non-existing
   faction.
 - Quests may reference resources. Quests must never reference a non-existing
   resources type.
 - Quests may reference other quests. Quests must never reference a
   non-existing quest.
 - Quests may reference known and unknown variables.
 - Quests may reference known and unknown events.

### Unique Names

 - Quest IDs must be unique accross all quests.
 - Quest display names must be unique across all quests within one language.
   (e.g. English)
 - Quest Stage Descriptions must be unique accross all quests within one
   lanuage. (e.g. English)

