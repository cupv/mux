package domain

// Card represents a vocabulary card entity
type Card struct {
	ID      int    `json:"id"`
	Word    string `json:"word"`
	Meaning string `json:"meaning"`
}

// CardRepository defines the interface for card storage operations
type CardRepository interface {
	GetAllCards() ([]Card, error)
	GetCardByID(id int) (*Card, error)
	CreateCard(card *Card) error
	UpdateCard(card *Card) error
	DeleteCard(id int) error
}
