package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	cors "github.com/itsjamie/gin-cors"
)

// SetupRouter sets up the Gin router
func SetupRouter(srv *Service) (*gin.Engine, error) {
	fmt.Println(" Initializing Gin router...")

	router := gin.New()

	// Middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(requestid.New())
	router.Use(cors.Middleware(cors.Config{
		Origins:        "*",
		Methods:        "GET, POST",
		RequestHeaders: "Origin, Authorization, Content-Type, Content-Length",
		ExposedHeaders: "",
	}))

	// status check
	router.GET("/CheckStatus", func(c *gin.Context) {
		fmt.Println(" GET /CheckStatus received")
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Register user
	router.POST("/api/register", RegisterUserHandler(srv.Db),
		SuccessResponseHandler,
	)

	// Get all users
	router.GET("/api/users", GetAllUsersHandler(srv.Db))
	router.GET("/api/users/:id", GetUserByIDHandler(srv.Db))
	router.GET("/api/users/count", GetUserCountHandler(srv.Db))
	router.PUT("/api/users/:id", UpdateUserHandler(srv.Db))
	router.DELETE("/api/users/:id", DeleteUserHandler(srv.Db))

	return router, nil
}
