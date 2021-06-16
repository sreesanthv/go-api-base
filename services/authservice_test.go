package services

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/sreesanthv/go-api-base/database/mock"
)

func NewTestAuthService() *AuthService {
	return NewAuthService(logrus.New(), mock.NewStore(), mock.NewRedis())
}

func TestGetAccount(t *testing.T) {
	as := NewTestAuthService()
	as.CreateAccount(map[string]string{
		"name":     "Bob",
		"email":    "bob@fakemail.go",
		"password": "goutpassword",
	})

	user, err := as.Store.GetAccount("bob@fakemail.go")
	if err != nil {
		t.Error("Failed to get account from store", err)
	}

	if user.Password == "goutpassword" {
		t.Error("Password saved without hashing.")
	}

	if user.ID == 0 {
		t.Errorf("Invalid ID for the account")
	}
}

func TestGetAccountById(t *testing.T) {
	as := NewTestAuthService()
	addTestAccounts(as)

	user, err := as.Store.GetAccountById(2)
	if err != nil {
		t.Error("Failed to get account from store", err)
	}

	if user.ID == 0 {
		t.Errorf("Invalid ID for the account")
	}
}

func TestIsValidPassword(t *testing.T) {
	as := NewTestAuthService()
	addTestAccounts(as)

	act := as.GetAccountById(2)
	if as.IsValidPassword(act, "marthautpassword") == false {
		t.Error("Failed to check hash with original password")
	}
	if as.IsValidPassword(act, "goutpassword") == true {
		t.Error("Failed to check hash with fake password")
	}
}

func TestCreateToken(t *testing.T) {
	as := NewTestAuthService()
	addTestAccounts(as)
	act := as.GetAccountById(1)

	td, err := as.CreateToken(act)
	if err != nil {
		t.Error("Failed to create token:", err)
	}

	if td.AccessToken == td.RefreshToken {
		t.Error("Access token and refresh token is same")
	}

	if td.AccessUuid == td.RefreshUuid {
		t.Error("Access uuid and refresh uuid is same")
	}

	err = as.PersistToken(act.ID, td)
	if err != nil {
		t.Error("Error is persiting token:", err)
	}

	tdParsed, err := as.ParseToken(td.AccessToken, TokenTypeAccess)
	if err != nil {
		t.Error("Error is parsing access token:", err)
	}

	if tdParsed.UserId == 0 || tdParsed.UserId != act.ID {
		t.Errorf("Mismatch in parsed User ID. Expected: %d, Got: %d", act.ID, tdParsed.UserId)
	}

	// parse access token as refresh token
	tdParsed, err = as.ParseToken(td.AccessToken, TokenTypeRefresh)
	if err == nil {
		t.Error("Access token parsed as refresh token:", err)
	}

	as.DropToken(td.AccessToken)

	tdParsed, err = as.ParseToken(td.RefreshToken, TokenTypeRefresh)
	if err != nil {
		t.Error("Error is parsing refresh token:", err)
	}
}

func addTestAccounts(as *AuthService) {
	as.CreateAccount(map[string]string{
		"name":     "Bob",
		"email":    "bob@fakemail.go",
		"password": "goutpassword",
	})
	as.CreateAccount(map[string]string{
		"name":     "Martha",
		"email":    "martha@fakemail.go",
		"password": "marthautpassword",
	})
	as.CreateAccount(map[string]string{
		"name":     "Raju",
		"email":    "raju@fakemail.go",
		"password": "raju4u",
	})
}
