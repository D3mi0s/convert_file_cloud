package main

import (
	"log"
	"os"

	"auth-service/handlers"
	"auth-service/models"
	"auth-service/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {
	dsn := "host=localhost user=postgres password=admin dbname=auth_service port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err := db.AutoMigrate(&models.User{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	userRepo := &repository.UserRepository{DB: db}

	authHandler := &handlers.AuthHandler{
		UserRepo: userRepo,
		JWTKey:   []byte(os.Getenv("JWT_SECRET_KEY")),
	}

	r := gin.Default()
	r.Use(CORSMiddleware())

	r.POST("/register", authHandler.Register)
	r.POST("/login", authHandler.Login)

	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
