package usecase

import (
	"github.com/cupv/mux/internal/domain"
	"github.com/cupv/mux/internal/repository"
)

type CreateCardItem struct {
	Word    string
	Meaning string
}

type CardUsecase interface {
	FetchCards() ([]domain.Card, error)
	Create(item CreateCardItem) (int64, error)
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

func (u *cardUsecase) Create(item CreateCardItem) (int64, error) {
	return u.cardRepo.Add(repository.AddCardItem{
		Word:    item.Word,
		Meaning: item.Meaning,
	})
}
