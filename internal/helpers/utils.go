package helpers

import (
	model "finalreg/internal/models"
	"fmt"
	"regexp"

	"gorm.io/gorm"
)

// ValidateRegex checks if a field matches a regex pattern
func ValidateRegex(field, value, pattern, errorMsg string) error {
	fmt.Printf(" Validating %s with regex...\n", field)
	matched := regexp.MustCompile(pattern).MatchString(value)
	if !matched {
		err := fmt.Errorf("%s: %s", field, errorMsg)
		fmt.Printf(" Regex failed: %v\n", err)
		return err
	}
	fmt.Printf(" %s passed regex check\n", field)
	return nil
}

// IsEmailUnique checks if email already exists in DB
func IsEmailUnique(db *gorm.DB, email string) error {
	fmt.Printf(" Checking if email '%s' is unique...\n", email)
	var count int64
	err := db.Model(&model.User{}).Where("email = ?", email).Count(&count).Error
	if err != nil {
		fmt.Printf(" DB error while checking email: %v\n", err)
		return err
	}
	if count > 0 {
		err := fmt.Errorf("email already exists")
		fmt.Printf(" %v\n", err)
		return err
	}
	fmt.Println(" Email is unique")
	return nil
}

// IsUsernameUnique checks if username already exists
func IsUsernameUnique(db *gorm.DB, username string) error {
	fmt.Printf(" Checking if username '%s' is unique...\n", username)
	var count int64
	err := db.Model(&model.User{}).Where("username = ?", username).Count(&count).Error
	if err != nil {
		fmt.Printf(" DB error while checking username: %v\n", err)
		return err
	}
	if count > 0 {
		err := fmt.Errorf("username already exists")
		fmt.Printf(" %v\n", err)
		return err
	}
	fmt.Println(" Username is unique")
	return nil
}
