package migrate

import (
	"fmt"

	"github.com/go-pg/migrations"
)

const bootstrapAdminAccount = `
INSERT INTO accounts (id, email, name, active, roles)
VALUES (DEFAULT, 'admin@local.io', 'Admin User', true, '{admin}')
`

const bootstrapUserAccount = `
INSERT INTO accounts (id, email, name, active)
VALUES (DEFAULT, 'user@local.io', 'User', true)`

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
