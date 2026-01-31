package http

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/tinx/pat-quest-editor/backend/internal/domain"
	"github.com/tinx/pat-quest-editor/backend/internal/ports"
)

// maxRequestBodySize limits request body to 10MB to prevent memory exhaustion attacks
const maxRequestBodySize = 10 * 1024 * 1024

// Handler provides HTTP handlers for the quest editor API.
type Handler struct {
	quests    ports.QuestRepository
	refData   ports.ReferenceDataRepository
	metadata  ports.MetadataRepository
	validator ports.QuestValidator
}

// NewHandler creates a new HTTP handler.
func NewHandler(
	quests ports.QuestRepository,
	refData ports.ReferenceDataRepository,
	metadata ports.MetadataRepository,
	validator ports.QuestValidator,
) *Handler {
	return &Handler{
		quests:    quests,
		refData:   refData,
		metadata:  metadata,
		validator: validator,
	}
}

// RegisterRoutes registers all API routes on the given mux.
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	// Quest endpoints
	mux.HandleFunc("/api/quests", h.handleQuests)
	mux.HandleFunc("/api/quests/", h.handleQuest)

	// Reference data endpoints
	mux.HandleFunc("/api/items", h.handleItems)
	mux.HandleFunc("/api/factions", h.handleFactions)
	mux.HandleFunc("/api/resources", h.handleResources)
	mux.HandleFunc("/api/npcs", h.handleNPCs)

	// Metadata endpoints
	mux.HandleFunc("/api/metadata/", h.handleMetadata)

	// Validation endpoint
	mux.HandleFunc("/api/validate", h.handleValidate)
}

func (h *Handler) handleQuests(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.listQuests(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) handleQuest(w http.ResponseWriter, r *http.Request) {
	questID := strings.TrimPrefix(r.URL.Path, "/api/quests/")
	if questID == "" {
		http.Error(w, "quest ID required", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getQuest(w, r, questID)
	case http.MethodPut:
		h.saveQuest(w, r, questID)
	case http.MethodDelete:
		h.deleteQuest(w, r, questID)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) listQuests(w http.ResponseWriter, r *http.Request) {
	questIDs, err := h.quests.List()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.writeJSON(w, questIDs)
}

func (h *Handler) getQuest(w http.ResponseWriter, r *http.Request, questID string) {
	quest, err := h.quests.Get(questID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Include metadata in response
	metadata, _ := h.metadata.GetQuestMetadata(questID)

	response := struct {
		Quest    *domain.Quest         `json:"quest"`
		Metadata *domain.QuestMetadata `json:"metadata,omitempty"`
	}{
		Quest:    quest,
		Metadata: metadata,
	}

	h.writeJSON(w, response)
}

func (h *Handler) saveQuest(w http.ResponseWriter, r *http.Request, questID string) {
	var request struct {
		Quest    domain.Quest          `json:"quest"`
		Metadata *domain.QuestMetadata `json:"metadata,omitempty"`
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodySize)
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Ensure quest ID in URL matches body
	if request.Quest.QuestID != questID {
		http.Error(w, "quest ID mismatch", http.StatusBadRequest)
		return
	}

	// Validate the quest
	validationResult := h.validator.Validate(&request.Quest)

	// Save quest even if invalid (allows work-in-progress saves)
	if err := h.quests.Save(&request.Quest); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Save metadata if provided
	if request.Metadata != nil {
		request.Metadata.QuestID = questID
		if err := h.metadata.SaveQuestMetadata(request.Metadata); err != nil {
			log.Printf("Warning: failed to save metadata for quest %s: %v", questID, err)
		}
	}

	h.writeJSON(w, validationResult)
}

func (h *Handler) deleteQuest(w http.ResponseWriter, r *http.Request, questID string) {
	if err := h.quests.Delete(questID); err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Also delete metadata
	h.metadata.DeleteQuestMetadata(questID)

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) handleItems(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	items, err := h.refData.ListItems()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.writeJSON(w, items)
}

func (h *Handler) handleFactions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	factions, err := h.refData.ListFactions()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.writeJSON(w, factions)
}

func (h *Handler) handleResources(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	resources, err := h.refData.ListResources()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.writeJSON(w, resources)
}

func (h *Handler) handleNPCs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	npcs, err := h.refData.ListNPCs()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.writeJSON(w, npcs)
}

func (h *Handler) handleMetadata(w http.ResponseWriter, r *http.Request) {
	questID := strings.TrimPrefix(r.URL.Path, "/api/metadata/")
	if questID == "" {
		http.Error(w, "quest ID required", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		metadata, err := h.metadata.GetQuestMetadata(questID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		h.writeJSON(w, metadata)

	case http.MethodPut:
		var metadata domain.QuestMetadata
		r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodySize)
		if err := json.NewDecoder(r.Body).Decode(&metadata); err != nil {
			http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
			return
		}
		metadata.QuestID = questID
		if err := h.metadata.SaveQuestMetadata(&metadata); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) handleValidate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var quest domain.Quest
	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodySize)
	if err := json.NewDecoder(r.Body).Decode(&quest); err != nil {
		http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	result := h.validator.Validate(&quest)
	h.writeJSON(w, result)
}

func (h *Handler) writeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
