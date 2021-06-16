package interfaces

import "github.com/sreesanthv/go-api-base/database"

type Store interface {
	// account related
	GetAccount(email string) (*database.AccountStore, error)
	GetAccountById(id int64) (*database.AccountStore, error)
	CreateAccount(info map[string]string) (*database.AccountStore, error)
	SaveLoginTime(userId int64) error
}
