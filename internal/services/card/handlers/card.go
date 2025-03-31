package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/cupv/mux/internal/services/card/models"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

var cards []models.Card

// InitRouter initializes the router with all card-related routes
func InitRouter() *mux.Router {
	router := mux.NewRouter()

	// Define routes for Vocabulary Card API
	router.HandleFunc("/card", GetCards).Methods("GET")
	router.HandleFunc("/card", CreateCard).Methods("POST")
	router.HandleFunc("/card/{id}", GetCard).Methods("GET")
	router.HandleFunc("/card/{id}", UpdateCard).Methods("PUT")
	router.HandleFunc("/card/{id}", DeleteCard).Methods("DELETE")

	return router
}

// GetCards retrieves the list of all cards
func GetCards(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
    // Simulate slow response (5 seconds) to demonstrate that gracefull shutdown work 
	// time.Sleep(5 * time.Second)
	time.Sleep(15 * time.Second)
	json.NewEncoder(w).Encode(cards)
}

// CreateCard creates a new card
func CreateCard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var card models.Card
	if err := json.NewDecoder(r.Body).Decode(&card); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Assign UUID and creation time
	card.ID = uuid.New().String()
	card.CreatedAt = time.Now().Unix()

	cards = append(cards, card)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(card)
}

// GetCard retrieves details of a card by ID
func GetCard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id := params["id"]

	for _, card := range cards {
		if card.ID == id {
			json.NewEncoder(w).Encode(card)
			return
		}
	}
	respondWithError(w, http.StatusNotFound, "Card not found")
}

// UpdateCard updates an existing card
func UpdateCard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id := params["id"]

	for i, card := range cards {
		if card.ID == id {
			var updatedCard models.Card
			if err := json.NewDecoder(r.Body).Decode(&updatedCard); err != nil {
				respondWithError(w, http.StatusBadRequest, "Invalid request payload")
				return
			}
			updatedCard.ID = id
			updatedCard.CreatedAt = cards[i].CreatedAt // Preserve original CreatedAt
			cards[i] = updatedCard
			json.NewEncoder(w).Encode(updatedCard)
			return
		}
	}
	respondWithError(w, http.StatusNotFound, "Card not found")
}

// DeleteCard deletes a card by ID
func DeleteCard(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	for i, card := range cards {
		if card.ID == id {
			cards = append(cards[:i], cards[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	respondWithError(w, http.StatusNotFound, "Card not found")
}

// respondWithError sends a standardized error response
func respondWithError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
