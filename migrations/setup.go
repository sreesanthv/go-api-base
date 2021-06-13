package migrations

import (
	"log"

	"github.com/go-pg/migrations"
	"github.com/go-pg/pg"
	"github.com/spf13/viper"
)

// Migrate runs go-pg migrations
func Migrate(args []string) {
	db, err := DBConn()
	if err != nil {
		log.Fatal(err)
	}

	err = db.RunInTransaction(func(tx *pg.Tx) error {
		oldVersion, newVersion, err := migrations.Run(tx, args...)
		if err != nil {
			return err
		}
		if newVersion != oldVersion {
			log.Printf("migrated from version %d to %d\n", oldVersion, newVersion)
		} else {
			log.Printf("version is %d\n", oldVersion)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

}

// Reset runs reverts all migrations to version 0 and then applies all migrations to latest
func Reset() {
	db, err := DBConn()
	if err != nil {
		log.Fatal(err)
	}

	version, err := migrations.Version(db)
	if err != nil {
		log.Fatal(err)
	}

	err = db.RunInTransaction(func(tx *pg.Tx) error {
		for version != 0 {
			oldVersion, newVersion, err := migrations.Run(tx, "down")
			if err != nil {
				return err
			}
			log.Printf("migrated from version %d to %d\n", oldVersion, newVersion)
			version = newVersion
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}

// DB connection
func DBConn() (*pg.DB, error) {
	db := pg.Connect(&pg.Options{
		Network:  viper.GetString("db_network"),
		Addr:     viper.GetString("db_addr"),
		User:     viper.GetString("db_user"),
		Password: viper.GetString("db_password"),
		Database: viper.GetString("db_database"),
	})
	return db, nil
}
