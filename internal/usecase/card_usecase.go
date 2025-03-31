package usecase

import (
	"github.com/cupv/mux/internal/domain"
	"github.com/cupv/mux/internal/repository"
)

type CardUsecase interface {
	FetchCards() ([]domain.Card, error)
}

type cardUsecase struct {
	cardRepo repository.CardRepository
}

func NewCardUsecase(cardRepo repository.CardRepository) CardUsecase {
	return &cardUsecase{cardRepo}
}

func (u *cardUsecase) FetchCards() ([]domain.Card, error) {
	return u.cardRepo.GetAllCards()
}
