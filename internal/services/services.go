package services

import (
	"fmt"
	"net/http"
	"regexp"
	"time"

	"finalreg/enums"

	"finalreg/internal/forms"
	model "finalreg/internal/models"
	"finalreg/internal/providers"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// RegisterUserHandler handles the /api/register endpoint to post the data
func RegisterUserHandler(repo providers.RepoStore, Name string) gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("\n\n REGISTER Starting in services.go ")

		// Parse JSON into struct
		var input forms.UserForm
		fmt.Println(" Parsing JSON...")
		if err := c.ShouldBindJSON(&input); err != nil {
			fmt.Printf(" JSON error: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		fmt.Printf(" Parsed input: %+v\n", input)

		// Validate the data
		if err := validateInput(input); err != nil {
			fmt.Printf(" Validation failed: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		fmt.Println(" All validations passed!")

		// Hash password
		fmt.Println(" Hashing password...")
		hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			fmt.Printf(" Hashing failed: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not hash password"})
			return
		}
		fmt.Println(" Password hashed")

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

		// Save to DB
		fmt.Println(" Saving user to database...")
		if err := repo.CreateUser(user); err != nil {
			fmt.Printf(" Failed to save user: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}
		fmt.Printf(" User saved! ID: %d, Email: %s\n", user.ID, user.Email)

		// Success response
		c.Set("userID", user.ID)
		c.Set("email", user.Email)
		c.Next()
	}
}

// handles GET /api/users
func GetAllUsersHandler(repo providers.RepoStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println(" GET /api/users received in services.go")

		db := repo.GetDB()
		var users []model.User

		if err := db.Find(&users).Error; err != nil {
			fmt.Printf(" Failed to fetch users: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch users"})
			return
		}

		fmt.Printf(" Found %d users\n", len(users))
		c.JSON(http.StatusOK, users)
	}
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
