package mock

import (
	"github.com/jackc/pgx"
	"github.com/sreesanthv/go-api-base/database"
)

type StoreMock struct {
	store map[int64]database.AccountStore
	id    int64
}

func NewStore() *StoreMock {
	store := make(map[int64]database.AccountStore)
	return &StoreMock{
		store: store,
	}
}

func (s *StoreMock) CreateAccount(info map[string]string) (*database.AccountStore, error) {
	s.id++
	a := database.AccountStore{
		ID:       s.id,
		Email:    info["email"],
		Name:     info["name"],
		Password: info["password"],
	}
	s.store[a.ID] = a
	return &a, nil
}

func (s *StoreMock) GetAccount(email string) (*database.AccountStore, error) {
	var act *database.AccountStore
	err := pgx.ErrNoRows
	for _, account := range s.store {
		if email == account.Email {
			act = &account
			err = nil
			break
		}
	}
	return act, err
}

func (s *StoreMock) GetAccountById(id int64) (*database.AccountStore, error) {
	val, ok := s.store[id]
	if !ok {
		return &val, pgx.ErrNoRows
	}
	return &val, nil
}

func (s *StoreMock) SaveLoginTime(userId int64) error {
	return nil
}
