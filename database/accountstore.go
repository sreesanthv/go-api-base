package database

import "fmt"

type AccountStore struct {
	ID       int64
	Email    string
	Name     string
	Password string
}

func (s *Store) GetUser(email string) (*AccountStore, error) {
	as := new(AccountStore)
	err := s.db.QueryRow(s.ctx, "SELECT id, email, name, password FROM accounts WHERE email = $1", email).Scan(
		&as.ID, &as.Email, &as.Name, &as.Password,
	)

	if err != nil && s.IsSQLError(err) {
		fmt.Println("Error in GetUser query", err)
		return as, err
	}

	return as, nil
}

func (s *Store) GetUserById(id int64) (*AccountStore, error) {
	as := new(AccountStore)
	err := s.db.QueryRow(s.ctx, "SELECT id, email, name, password FROM accounts WHERE id = $1", id).Scan(
		&as.ID, &as.Email, &as.Name, &as.Password,
	)

	if err != nil && s.IsSQLError(err) {
		fmt.Println("Error in GetUser query", err)
		return as, err
	}

	return as, nil
}
