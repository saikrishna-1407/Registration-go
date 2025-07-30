package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println(" Middleware: ErrorHandler started")

		defer func() {
			if r := recover(); r != nil {
				fmt.Printf(" PANIC RECOVERED: %v\n", r)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong!"})
			}
		}()

		fmt.Println(" Proceeding to next handler...")
		c.Next()

	}
}
