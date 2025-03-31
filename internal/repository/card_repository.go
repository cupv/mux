package repository

import (
    "database/sql"
	"github.com/cupv/mux/internal/domain"
)

type CardRepository interface {
    GetAllCards() ([]domain.Card, error)
}

type cardRepository struct {
    db *sql.DB
}

func NewCardRepository(db *sql.DB) CardRepository {
    return &cardRepository{db}
}

func (r *cardRepository) GetAllCards() ([]domain.Card, error) {
    rows, err := r.db.Query("SELECT id, word, meaning FROM cards")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var cards []domain.Card
    for rows.Next() {
        var card domain.Card
        if err := rows.Scan(&card.ID, &card.Word, &card.Meaning); err != nil {
            return nil, err
        }
        cards = append(cards, card)
    }
    return cards, nil
}
