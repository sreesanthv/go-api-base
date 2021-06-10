package database

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/sirupsen/logrus"
)

type Store struct {
	db     *pgx.Conn
	logger *logrus.Logger
	ctx    context.Context
}

func NewStore(db *pgx.Conn, logger *logrus.Logger) *Store {
	return &Store{
		db:     db,
		logger: logger,
		ctx:    context.Background(),
	}
}

func (s *Store) IsSQLError(err error) bool {
	er := false
	if err != pgx.ErrNoRows {
		er = true
	}
	return er
}
