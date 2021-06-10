package services

import (
	"github.com/sirupsen/logrus"
	"github.com/sreesanthv/go-api-base/database"
)

type AuthService struct {
	logger *logrus.Logger
	store  *database.Store
	redis  *database.Redis
}

func NewAuthService(log *logrus.Logger, store *database.Store, redis *database.Redis) *AuthService {
	return &AuthService{
		logger: log,
		store:  store,
		redis:  redis,
	}
}

func (s *AuthService) GetUser(email string) *database.AccountStore {
	user, _ := s.store.GetUser(email)
	return user
}

// validate password entered  - login
func (s *AuthService) IsValidPassword(act *database.AccountStore, password string) bool {
	return CheckPasswordHash(password, act.Password)
}
