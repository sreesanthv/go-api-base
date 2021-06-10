// Package database implements postgres connection and queries.
package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/spf13/viper"
)

// DBConn returns a postgres connection pool.
func DBConn() (*pgx.Conn, error) {

	connUrl := fmt.Sprintf("postgres://%s:%s@%s/%s", viper.GetString("db_user"), viper.GetString("db_password"), viper.GetString("db_addr"), viper.GetString("db_database"))
	db, err := pgx.Connect(context.Background(), connUrl)
	if err != nil {
		return nil, err
	}

	return db, nil
}
