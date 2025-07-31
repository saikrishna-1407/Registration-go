package services

import (
	"fmt"
	"net/http"
	"regexp"
	"time"

	"finalreg/enums"
	model "finalreg/internal/models"
	"finalreg/internal/providers"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// RegisterUserHandler handles the /api/register endpoint to post the data
func RegisterUserHandler(repo providers.RepoStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("\n\n REGISTER Starting in services.go ")

		// This will Parse JSON
		var input map[string]interface{}
		fmt.Println(" Parsing JSON...")
		if err := c.ShouldBindJSON(&input); err != nil {
			fmt.Printf(" JSON error: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}
		fmt.Printf(" Parsed input: %+v\n", input)

		// This will Validate the data
		if err := validateInput(input); err != nil {
			fmt.Printf(" Validation failed: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		fmt.Println(" All validations passed!")

		// This will Hash the password
		password := input["password"].(string)
		fmt.Println(" Hashing password...")
		hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			fmt.Printf(" Hashing failed: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not hash password"})
			return
		}
		fmt.Println(" Password hashed")

		// This will parse the  DOB
		dobStr := input["dateOfBirth"].(string)
		dob, err := time.Parse("2006-01-02", dobStr)
		if err != nil {
			fmt.Printf(" Invalid date format: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
			return
		}

		// This will Create user model
		user := &model.User{
			FullName:      input["fullName"].(string),
			Email:         input["email"].(string),
			Password:      string(hashed),
			Username:      input["username"].(string),
			DateOfBirth:   dob,
			PhoneNumber:   input["phoneNumber"].(string),
			Gender:        input["gender"].(string),
			Country:       input["country"].(string),
			State:         input["state"].(string),
			PinCode:       input["pinCode"].(string),
			ReferralCode:  input["referralCode"].(string),
			TermsAccepted: input["termsAccepted"].(bool),
		}

		// This will Save data to database
		fmt.Println(" Saving user to database...")
		if err := repo.CreateUser(user); err != nil {
			fmt.Printf(" Failed to save user: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}
		fmt.Printf(" User saved! ID: %d, Email: %s\n", user.ID, user.Email)

		// This will Send success message back
		c.JSON(http.StatusCreated, gin.H{
			"message": "User registered successfully",
			"userId":  user.ID,
			"email":   user.Email,
		})
	}
}

// GetAllUsersHandler handles GET /api/users to get the data
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

// This will validate all the Input Data
func validateInput(input map[string]interface{}) error {
	fmt.Println(" Running validations...")

	// 1. Email
	email, _ := input["email"].(string)
	if email == "" {
		return fmt.Errorf("email is required")
	}
	if !regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`).MatchString(email) {
		return fmt.Errorf("invalid email format")
	}
	fmt.Println(" Email valid")

	// 2. Username
	username, _ := input["username"].(string)
	if len(username) < 4 || len(username) > 20 {
		return fmt.Errorf("username must be 4-20 characters")
	}
	if !regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString(username) {
		return fmt.Errorf("username can only have letters, numbers, underscore")
	}
	fmt.Println(" Username valid")

	// 3. Full Name
	fullName, _ := input["fullName"].(string)
	if len(fullName) < 3 || len(fullName) > 50 {
		return fmt.Errorf("full name must be 3-50 characters")
	}
	if !regexp.MustCompile(`^[a-zA-Z\s]+$`).MatchString(fullName) {
		return fmt.Errorf("full name can only have letters and spaces")
	}
	fmt.Println(" Full name valid")

	// 4. Password
	password, _ := input["password"].(string)
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters")
	}
	if !regexp.MustCompile(`[A-Z]`).MatchString(password) {
		return fmt.Errorf("password must have uppercase")
	}
	if !regexp.MustCompile(`[a-z]`).MatchString(password) {
		return fmt.Errorf("password must have lowercase")
	}
	if !regexp.MustCompile(`[0-9]`).MatchString(password) {
		return fmt.Errorf("password must have number")
	}
	if !regexp.MustCompile(`[!@#$%^&*]`).MatchString(password) {
		return fmt.Errorf("password must have special character")
	}
	if password != input["confirmPassword"].(string) {
		return fmt.Errorf("passwords do not match")
	}
	fmt.Println(" Password valid")

	// 5. DOB (18+)
	dobStr, _ := input["dateOfBirth"].(string)
	dob, err := time.Parse("2006-01-02", dobStr)
	if err != nil {
		return fmt.Errorf("invalid date format")
	}
	if time.Since(dob).Hours()/24/365 < 18 {
		return fmt.Errorf("you must be 18 or older")
	}
	fmt.Println(" Age valid")

	// 6. Phone Number
	phone, _ := input["phoneNumber"].(string)
	if !regexp.MustCompile(`^\+91\d{10}$`).MatchString(phone) {
		return fmt.Errorf("phone must be +91XXXXXXXXXX")
	}
	fmt.Println(" Phone valid")

	// 7. Gender
	gender, _ := input["gender"].(string)
	if gender != "" {
		valid := stringInSlice(gender, enums.ValidGenders)
		if !valid {
			return fmt.Errorf("gender must be male, female, other, or prefer_not_to_say")
		}
	}
	fmt.Println(" Gender valid")

	// 8. Country
	country, _ := input["country"].(string)
	if !stringInSlice(country, enums.ValidCountries) {
		return fmt.Errorf("invalid country")
	}
	fmt.Println(" Country valid")

	// 9. State & PinCode (if India)
	if country == "India" {
		state, _ := input["state"].(string)
		if state == "" || !stringInSlice(state, enums.IndianStates) {
			return fmt.Errorf("invalid Indian state")
		}
		pinCode, _ := input["pinCode"].(string)
		if !regexp.MustCompile(`^\d{6}$`).MatchString(pinCode) {
			return fmt.Errorf("pin code must be 6 digits")
		}
		fmt.Println(" State & pinCode valid for India")
	}

	// 10. Terms
	if !input["termsAccepted"].(bool) {
		return fmt.Errorf("you must accept terms")
	}
	fmt.Println(" Terms accepted")

	fmt.Println(" All validations passed!")
	return nil
}

// stringInSlice checks if a string is in a slice
func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
