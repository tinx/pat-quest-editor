# Potions and Tinctures Quest Editor

## Context

"Potions and Tinctures", or "PAT" for short, is a game that is currently
in active development. The game contains a number of quests that the player
may or may not accept and play through.

Instead of hard-coding the quest data, stages, dialogs and other parameters
into the game code, it is better to abstract the quest content from the
game logic.

To achieve that, quests can be defined in YAML files. A JSON schema defines
which options are available for quests. And then there is one YAML file
per quest to specifiy what that particular quest looks like.

## Supported Quest Features

For PAT, we implemented these features:

 - Quests that can start in multiple different ways.
 - Quests that are only available under certain conditions, such as
   - Faction standing
   - Completing other quests first
   - World status conditions (such as "It's summer")
   - etc.
 - Repeatable quests and partially repeatable quests
 - Quests that require making a decision, which alters the outcome
 - Quests that can progress in multiple ways, for example:
   - Solve 3 out of 5 tasks, it's your choice which ones
   - Make a choice: deliver the item, or tell the guards about it
 - Quests that give multiple sub-tasks which can be handled in parallel
 - Quest rewards
 - Punishment for failure
 - Quests that require specific tasks to be completed n times
 - Per-quest journal entries, such as status updates and post-quest summaries
 - Quests with steps that can be accomplished before the quest was accepted
 - Quests involving quest items, which can't be traded or sold
 - Dialog options that are only available under certain conditions
 - Quest progression options that are only available under certain conditions
 - Quest versions, enabling migrations from older to newer states in savegames

Our data models must support storing and restoring all quest states, so
that we can enable the player to save and load the game state.

## Quest Nodes

PAT Quests are designed via a graph of nodes, which are connected to show the
flow of control. There are several types of nodes, such as dialog nodes,
action nodes, nodes waiting for certain things to happen, and so on.

Any node can potentially be reached in multiple ways, and can also potentially
lead on to multiple nodes further down.

That means quests are state machines, but can potentially have multiple active
states. (for example, a quest that can start in two different ways will be in
two active states right from the start)

## Quest Editor

The editor is implemented as a web application with a server backend. This
makes it possible to edit quests from any browser, without installing local
software. Quest files are kept on the server, as there is no need for them
to be downloaded or uploaded.

In the editor, quest nodes can be created, placed and connected. Node types
and node settings influence how the quest will play out. The editor is
mostly visual, in order to give a low barrier of entry for designing quests.

Ultimately, one quest is edited at a time, and the result is reflected in
the quest YAML file. This file can then be checked into version control such
as git.

## Quest Checker CLI

The checker is a standalone command-line utility that validates all quest files.
It performs both single-quest validation and cross-quest checks.

### Building

```bash
cd checker
go build .
```

### Usage

```bash
./checker -quests ../quests -data ../data
```

Options:
- `-quests` - Path to quests directory (default: `./quests`)
- `-data` - Path to reference data directory (default: `./data`)
- `-quiet` - Only output errors, no summary

Exit codes:
- `0` - All quests valid
- `1` - Validation issues found
- `2` - Fatal error (e.g., can't read files)

### Validation Rules

Single-quest:
- Unique NodeIDs within quest
- Valid node connections (no non-existent, self-referencing, or duplicate edges)
- Non-EntryPoint nodes have incoming connections
- At least one EntryPoint exists
- No cycles (DAG only)
- Terminal nodes have no outgoing edges
- Non-terminal nodes have outgoing edges
- References to NPCs, items, factions, resources, objects exist

Cross-quest:
- Unique QuestIDs across all quests
- Unique DisplayNames per language
- Unique QuestStageDescriptions per language
- QuestCompleted conditions reference existing quests

