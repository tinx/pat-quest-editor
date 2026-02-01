# Copilot Instructions for PAT Quest Editor

## Project Overview

This is a quest editor for "Potions and Tinctures" (PAT), a game in active development. The editor is a web application for visually designing quest state machines stored as YAML files.

## Build & Run Commands

```bash
# Backend
cd backend
go build ./...                    # Build all
go test ./...                     # Run all tests
go test ./internal/app -v         # Run validator tests with verbose output
go run ./cmd/server               # Start server on :8080

# Frontend
cd frontend
npm install                       # Install dependencies
npm run dev                       # Dev server with hot reload
npm run build                     # Production build to dist/

# Full stack development
# Terminal 1: cd backend && go run ./cmd/server
# Terminal 2: cd frontend && npm run dev
# Frontend proxies /api to backend via vite.config.js
```

## Architecture

### Stack
- **Backend**: GoLang with hexagonal architecture (ports/adapters pattern)
- **Frontend**: React + React Flow + Vite
- **Storage**: Quest files in YAML, editor metadata (node positions) in SQLite

### Backend Structure
```
backend/
├── cmd/server/          # Entry point, CLI flags
├── internal/
│   ├── domain/          # Core models (Quest, Item, Faction, etc.)
│   ├── ports/           # Interfaces (repositories, validators)
│   ├── adapters/
│   │   ├── filesystem/  # YAML file repositories
│   │   ├── storage/     # SQLite metadata repository
│   │   └── http/        # REST API handlers
│   └── app/             # Business logic (validator)
```

### Key Design Decisions
- No authentication - runs locally or behind authenticating proxy
- Minimize external dependencies; implement small features in-project
- All frontend assets must be self-hosted; no CDN links

## Quest Data Model

Quests are directed acyclic graphs (DAGs) with multiple possible active states.

### Node Types
- `EntryPoint` - Quest start points (only type allowed without incoming connections)
- `ConditionWatcher` - Waits for conditions to be met
- `Dialog` / `Decision` - NPC conversations with optional player choices
- `Actions` - Executes actions (rewards, quest state changes, stage updates)

### Terminal Actions
Action nodes ending a quest flow must contain exactly one of: `CompleteQuest`, `FailQuest`, `DeclineQuest`. These nodes must not have `NextNodes`.

### Validation Rules (beyond JSON schema)
- All non-EntryPoint nodes must have at least one incoming connection
- NodeIDs must be unique within a quest
- All NextNodes references must point to existing NodeIDs
- No cycles allowed (DAG only)
- References to NPCs, items, resources, and factions must exist in `data/*.yaml`

## File Locations

- Quest schema: `schemas/quest.json`
- Reference schemas: `schemas/{item,faction,resource,npc}.json`
- Quest files: `quests/`
- Reference data: `data/{items,factions,resources,npcs}.yaml`

## Conventions

### IDs
Pattern: `^[A-Z][A-Za-z0-9\.\-_:]*$` (e.g., `PAT_Demo_Quest`, `NPC:Smith`)

### i18n
All user-facing strings require both `en-US` and `de-DE` translations.

### Code Style
- Keep functions small and readable
- Unit tests mandatory for business logic
