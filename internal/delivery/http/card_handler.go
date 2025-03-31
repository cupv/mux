package http

import (
	"encoding/json"
	"net/http"

	"github.com/cupv/mux/internal/usecase"
)

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
