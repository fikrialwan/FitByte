package main

import (
	"log"
	"os"

	"github.com/fikrialwan/FitByte/config"
	"github.com/fikrialwan/FitByte/internal/controller"
	"github.com/fikrialwan/FitByte/internal/repository"
	"github.com/fikrialwan/FitByte/internal/routes"
	"github.com/fikrialwan/FitByte/internal/service"
	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()

	registerRoutesAndInjectDependency(server)

	run(server)
}

func run(server *gin.Engine) {
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	host := os.Getenv("APP_ENV")
	var serve string
	if host == "develop" {
		serve = "0.0.0.0:" + port
	} else {
		serve = ":" + port
	}

	if err := server.Run(serve); err != nil {
		log.Fatalf("error running server: %v", serve)
	}
}

func registerRoutesAndInjectDependency(server *gin.Engine) {
	db := config.InitDb()

	userRepository := repository.NewUserRepository(db)

	jwtService := service.NewJwtService()
	userService := service.NewUserService(userRepository, jwtService)

	userController := controller.NewUserController(userService)

	// registerRoutes
	routes.RegisterUserRoutes(server, userController, jwtService)
}
