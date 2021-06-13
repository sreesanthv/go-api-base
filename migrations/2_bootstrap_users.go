package migrations

import (
	"fmt"

	"github.com/go-pg/migrations"
)

const bootstrapAdminAccount = `
INSERT INTO accounts (id, email, name, active, roles, password)
VALUES (DEFAULT, 'admin@local.io', 'Admin User', true, '{admin}', '$2a$14$qa6K8ZwSK0.lQEYIPpfrW.ib0rhYsaAG1NoqZHhBcgOiZXdJ6LkbK')
`

const bootstrapUserAccount = `
INSERT INTO accounts (id, email, name, active, password)
VALUES (DEFAULT, 'user@local.io', 'User', true, '$2a$14$qa6K8ZwSK0.lQEYIPpfrW.ib0rhYsaAG1NoqZHhBcgOiZXdJ6LkbK')`

func init() {
	up := []string{
		bootstrapAdminAccount,
		bootstrapUserAccount,
	}

	down := []string{
		`TRUNCATE accounts CASCADE`,
	}

	migrations.Register(func(db migrations.DB) error {
		for _, q := range up {
			fmt.Println("up 2_bootstrap_users")
			_, err := db.Exec(q)
			if err != nil {
				return err
			}
		}
		return nil
	}, func(db migrations.DB) error {
		fmt.Println("down 2_bootstrap_users")
		for _, q := range down {
			_, err := db.Exec(q)
			if err != nil {
				return err
			}
		}
		return nil
	})
}
