package migrations

import (
	"fmt"

	"github.com/mehmetymw/debezium-postgres-es/infrastructure/persistence/models"
	"gorm.io/gorm"
)

// RunMigrations runs all database migrations
func RunMigrations(db *gorm.DB) error {
	fmt.Println("Running database migrations...")

	// Auto migrate the models
	if err := db.AutoMigrate(&models.Order{}); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	fmt.Println("Database migration completed")
	return nil
}
