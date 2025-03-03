package main

import (
	"github.com/gin-contrib/cors"
	"github.com/joho/godotenv"
	"log"
	"user_crud/controller"
	"user_crud/repository"
	"user_crud/service"
)

func loadEnvVariable() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	loadEnvVariable()
	db := repository.PostgresConnect()
	defer db.Close()
	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository)
	r := controller.NewUserController(userService)
	r.Use(cors.Default())
	err := r.Run("0.0.0.0:8080")
	if err != nil {
		return
	}
}
