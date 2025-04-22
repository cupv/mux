package http

import (
	"encoding/json"
	"net/http"

	"github.com/cupv/mux/internal/usecase"
)

type CreateCardDto struct {
	Word    string `json:"word"`
	Meaning string `json:"meaning"`
}

type CardHandler struct {
	usecase usecase.CardUsecase
}

func NewCardHandler(u usecase.CardUsecase) *CardHandler {
	return &CardHandler{u}
}

func (h *CardHandler) GetCards(w http.ResponseWriter, r *http.Request) {
	cards, err := h.usecase.FetchCards()
	if err != nil {
		http.Error(w, "Failed to retrieve cards", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cards)
}

func (h *CardHandler) Create(w http.ResponseWriter, r *http.Request) {

	var dto CreateCardDto

	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	cardId, err := h.usecase.Create(usecase.CreateCardItem{
		Word:    dto.Word,
		Meaning: dto.Meaning,
	})
	if err != nil {
		http.Error(w, "Failed to retrieve cards", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cardId)
}
