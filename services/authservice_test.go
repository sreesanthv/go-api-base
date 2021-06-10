package services

import (
	"testing"

	"github.com/jackc/pgx/v4"
	"github.com/sirupsen/logrus"
	"github.com/sreesanthv/go-api-base/database"
)

func TestIsValidPassword(t *testing.T) {
	as := NewTestAuthService()

	dbPass, _ := HashPassword("testpass")
	account := &database.AccountStore{
		Password: dbPass,
	}

	if !as.IsValidPassword(account, "testpass") {
		t.Error("Issue in matching account password", dbPass)
	}

	if as.IsValidPassword(account, "test_pass") {
		t.Error("Issue in matching account password", dbPass)
	}
}

func NewTestAuthService() *AuthService {
	log := logrus.Logger{}
	store := database.NewStore(&pgx.Conn{}, &log)
	return NewAuthService(logrus.New(), store)
}
