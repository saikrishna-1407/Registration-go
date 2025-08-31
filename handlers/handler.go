package handlers

import (
	"finalreg/internal/forms"
	model "finalreg/internal/models"
	"finalreg/internal/providers"
	"finalreg/internal/services"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Handler() {
	fmt.Println(" HTTP handler initialized")

	// Register the route
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		fmt.Fprintf(w, "Hello, you've reached the server on port 8000!\n")
	})

}
func SuccessResponseHandler(c *gin.Context) {
	userID, _ := c.Get("userID")
	email, _ := c.Get("email")

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"userId":  userID,
		"email":   email,
	})
}

func RegisterUserHandler(repo providers.RepoStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input forms.UserForm

		// Bind JSON
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
			return
		}

		// Validate passwords
		if input.Password != input.ConfirmPassword {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Passwords do not match"})
			return
		}

		// Call service
		user, err := services.RegisterUserService(repo, input)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// ✅ Success response
		c.JSON(http.StatusCreated, gin.H{
			"message": "User registered successfully",
			"userId":  user.ID,
			"email":   user.Email,
		})
	}
}

func GetAllUsersHandler(repo providers.RepoStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("GET /api/users received")

		users, err := services.GetAllUsersService(repo)
		if err != nil {
			fmt.Printf("Failed to fetch users: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch users"})
			return
		}

		fmt.Printf("Found %d users\n", len(users))
		c.JSON(http.StatusOK, users)
	}
}

func GetUserByIDHandler(repo providers.RepoStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		user, err := services.GetUserByIDService(repo, uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		c.JSON(http.StatusOK, user)
	}
}

func DeleteUserHandler(repo providers.RepoStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		if err := services.DeleteUserService(repo, uint(id)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
	}
}

func GetUserCountHandler(repo providers.RepoStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		var count int64
		db := repo.GetDB()
		db.Model(&model.User{}).Count(&count)

		c.JSON(http.StatusOK, gin.H{"total_users": count})
	}
}

func UpdateUserHandler(repo providers.RepoStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		var input model.User
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		user, err := services.UpdateUserService(repo, uint(id), input)
		if err != nil {
			if err.Error() == "user not found" {
				c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "User updated successfully",
			"user":    user,
		})
	}
}
