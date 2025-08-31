package services

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"finalreg/enums"

	"finalreg/internal/forms"
	model "finalreg/internal/models"
	"finalreg/internal/providers"

	"golang.org/x/crypto/bcrypt"
)

// RegisterUserHandler handles the /api/register endpoint to post the data
func RegisterUserService(repo providers.RepoStore, input forms.UserForm) (*model.User, error) {
	// Validate input
	if err := validateInput(input); err != nil {
		return nil, err
	}

	// Hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("could not hash password: %w", err)
	}

	// Create user model
	user := &model.User{
		FullName:      input.FullName,
		Email:         input.Email,
		Password:      string(hashed),
		Username:      input.Username,
		DateOfBirth:   input.DateOfBirth,
		PhoneNumber:   input.PhoneNumber,
		Gender:        input.Gender,
		Country:       input.Country,
		State:         input.State,
		PinCode:       input.PinCode,
		ReferralCode:  input.ReferralCode,
		TermsAccepted: input.TermsAccepted,
	}

	// Save user in DB
	createdUser, err := repo.CreateUser(user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return createdUser, nil
}

func UpdateUserService(repo providers.RepoStore, id uint, updatedData model.User) (*model.User, error) {
	db := repo.GetDB()

	var user model.User
	// find user
	if err := db.First(&user, id).Error; err != nil {
		return nil, errors.New("user not found")
	}
	// update allowed fields
	user.FullName = updatedData.FullName
	user.PhoneNumber = updatedData.PhoneNumber
	user.State = updatedData.State
	user.Country = updatedData.Country

	// save updates
	if err := db.Save(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// handles GET /api/users
func GetAllUsersService(repo providers.RepoStore) ([]model.User, error) {
	var users []model.User

	db := repo.GetDB()
	if err := db.Find(&users).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch users: %w", err)
	}

	return users, nil
}

func GetUserByIDService(repo providers.RepoStore, id uint) (*model.User, error) {
	var user model.User
	if err := repo.GetDB().First(&user, id).Error; err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return &user, nil
}

func DeleteUserService(repo providers.RepoStore, id uint) error {
	if err := repo.GetDB().Delete(&model.User{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

// validates all the input data
func validateInput(input forms.UserForm) error {
	fmt.Println(" Running validations...")

	// 1. Email
	if input.Email == "" {
		return fmt.Errorf("email is required")
	}
	if !regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`).MatchString(input.Email) {
		return fmt.Errorf("invalid email format")
	}
	fmt.Println(" Email valid")

	// 2. Username
	if len(input.Username) < 4 || len(input.Username) > 20 {
		return fmt.Errorf("username must be 4-20 characters")
	}
	if !regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString(input.Username) {
		return fmt.Errorf("username can only have letters, numbers, underscore")
	}
	fmt.Println(" Username valid")

	// 3. Full Name
	if len(input.FullName) < 3 || len(input.FullName) > 50 {
		return fmt.Errorf("full name must be 3-50 characters")
	}
	if !regexp.MustCompile(`^[a-zA-Z\s]+$`).MatchString(input.FullName) {
		return fmt.Errorf("full name can only have letters and spaces")
	}
	fmt.Println(" Full name valid")

	// 4. Password
	if len(input.Password) < 8 {
		return fmt.Errorf("password must be at least 8 characters")
	}
	if !regexp.MustCompile(`[A-Z]`).MatchString(input.Password) {
		return fmt.Errorf("password must have uppercase")
	}
	if !regexp.MustCompile(`[a-z]`).MatchString(input.Password) {
		return fmt.Errorf("password must have lowercase")
	}
	if !regexp.MustCompile(`[0-9]`).MatchString(input.Password) {
		return fmt.Errorf("password must have number")
	}
	if !regexp.MustCompile(`[!@#$%^&*]`).MatchString(input.Password) {
		return fmt.Errorf("password must have special character")
	}
	if input.Password != input.ConfirmPassword {
		return fmt.Errorf("passwords do not match")
	}
	fmt.Println(" Password valid")

	// 5. DOB (18+)
	if time.Since(input.DateOfBirth).Hours()/24/365 < 18 {
		return fmt.Errorf("you must be 18 or older")
	}
	fmt.Println(" Age valid")

	// 6. Phone Number
	if !regexp.MustCompile(`^\+91\d{10}$`).MatchString(input.PhoneNumber) {
		return fmt.Errorf("phone must be +91XXXXXXXXXX")
	}
	fmt.Println(" Phone valid")

	// 7. Gender
	if input.Gender != "" {
		if !stringInSlice(input.Gender, enums.ValidGenders) {
			return fmt.Errorf("gender must be male, female, other, or prefer_not_to_say")
		}
	}
	fmt.Println(" Gender valid")

	// 8. Country
	if !stringInSlice(input.Country, enums.ValidCountries) {
		return fmt.Errorf("invalid country")
	}
	fmt.Println(" Country valid")

	// 9. State & PinCode (if India)
	if input.Country == "India" {
		if input.State == "" || !stringInSlice(input.State, enums.IndianStates) {
			return fmt.Errorf("invalid Indian state")
		}
		if !regexp.MustCompile(`^\d{6}$`).MatchString(input.PinCode) {
			return fmt.Errorf("pin code must be 6 digits")
		}
		fmt.Println(" State & pinCode valid for India")
	}

	// 10. Terms
	if !input.TermsAccepted {
		return fmt.Errorf("you must accept terms")
	}
	fmt.Println(" Terms accepted")

	fmt.Println(" All validations passed!")
	return nil
}

// checks if a string is in a slice
func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
