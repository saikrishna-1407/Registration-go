package store

import (
	model "finalreg/internal/models"
	"fmt"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	fmt.Println(" Running migrations...")
	err := db.AutoMigrate(&model.User{})
	if err != nil {
		fmt.Printf(" Migration failed: %v\n", err)
		return err
	}
	fmt.Println(" Migration completed: User table ready")
	return nil
}
