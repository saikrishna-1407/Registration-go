package handlers

import (
	"fmt"
	"net/http"

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
