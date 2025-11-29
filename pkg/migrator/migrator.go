package migrator

import (
	"fmt"
	"log"

	"ice/config"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(cfg config.MySQLConfig) error {
	dsn := fmt.Sprintf("mysql://%s:%s@tcp(%s:%s)/%s?multiStatements=true",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)

	m, err := migrate.New(
		"file://internal/migration/mysql",
		dsn,
	)
	if err != nil {
		return fmt.Errorf("migration init failed: %w", err)
	}
	defer func() {
		if sourceErr, dbErr := m.Close(); sourceErr != nil || dbErr != nil {
			log.Printf("migration close error: source=%v, db=%v", sourceErr, dbErr)
		}
	}()

	if err := m.Up(); err != nil && err.Error() != "no change" {
		return fmt.Errorf("migration up failed: %w", err)
	}

	log.Println("migration complete")
	return nil
}
