# Quest Notes

## Specification

There is a JSON schema for quest YAML files. The schema can be found in
this repository at 'schemas/quest.json'. This makes it possible to validate
quest files for syntactic correctness, and limited semantic correctness.
However, there are additional rules and constraints that are not modeled
in the JSON schema, but must be verified in code.

## Rules beyond the Schema

 - The only node type allowed to not have any connections coming into
   it are 'EntryPoint' nodes. Every other node type must have at least
   one connection leading into it.
 - Flows must end in a node of type Actions, and one of the actions in
   this node must be one of:
    - CompleteQuest
    - FailQuest
    - DeclineQuest
   Actions nodes with any of these actions must not have a "NextNodes"
   field.
 - Actions nodes can only have one of these actions, never more
   than one of them:
    - CompleteQuest
    - FailQuest
    - DeclineQuest
 - Flows must never be cyclic. This is a directed acyclic graph.
   No loops allowed.
 - NodeID value must be unique within the quest.
 - The NextNodes field is an array of values indicating which quest
   node or nodes will continue the quest from the current node. Therefore,
   all NodeIDs given in this array must exist in the quest YAML file.
 - Quests must not reference unknown conversation partners or speakers.
 - Quests must not reference unknown items.
 - Quests must not reference unknown resources.
 - Quests must not reference unknown factions.
 - Quests may reference arbitrary variables.

