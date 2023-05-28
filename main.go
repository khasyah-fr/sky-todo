package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"

	"github.com/khasyah-fr/sky-todo/routes"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("failed to load environment variables: %v", err)
	}

	e := echo.New()

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	routes.ConfigureActivityRoutes(e)
	routes.ConfigureTodoRoutes(e)

	port := "3030"
	go e.Logger.Fatal(e.Start(":" + port))
}
