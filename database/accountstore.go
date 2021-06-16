package database

type AccountStore struct {
	ID       int64
	Email    string
	Name     string
	Password string
}

func (s *Store) GetAccount(email string) (*AccountStore, error) {
	as := new(AccountStore)
	err := s.db.QueryRow(s.ctx, "SELECT id, email, name, password FROM accounts WHERE email = $1 AND active = TRUE", email).Scan(
		&as.ID, &as.Email, &as.Name, &as.Password,
	)

	if err != nil && s.IsSQLError(err) {
		s.logger.Error("Error in GetAccount query:", err)
		return as, err
	}

	return as, nil
}

func (s *Store) GetAccountById(id int64) (*AccountStore, error) {
	as := new(AccountStore)
	err := s.db.QueryRow(s.ctx, "SELECT id, email, name, password FROM accounts WHERE id = $1 AND active = TRUE", id).Scan(
		&as.ID, &as.Email, &as.Name, &as.Password,
	)

	if err != nil && s.IsSQLError(err) {
		s.logger.Error("Error in GetAccount query", err)
		return as, err
	}

	return as, nil
}

func (s *Store) CreateAccount(info map[string]string) (*AccountStore, error) {
	var userId int64
	sql := "INSERT INTO accounts (email, name, password) VALUES($1, $2, $3) RETURNING id"
	err := s.db.QueryRow(s.ctx, sql, info["email"], info["name"], info["password"]).Scan(&userId)
	if err != nil {
		s.logger.Error("Error creating account:", err)
		return nil, err
	}

	return &AccountStore{
		ID:       userId,
		Email:    info["email"],
		Name:     info["name"],
		Password: info["password"],
	}, nil
}

func (s *Store) SaveLoginTime(userId int64) error {
	sql := "UPDATE accounts SET last_login	= NOW() WHERE id = $1"
	_, err := s.db.Exec(s.ctx, sql, userId)
	if err != nil {
		s.logger.Error("Error setting last login time:", err)
	}

	return err
}
