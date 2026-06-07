package main

import (
	"log"

	"parser/internal/db"
	"parser/internal/handler"
	"parser/internal/middleware"
	"parser/internal/repository"
	"parser/internal/router"

	"github.com/labstack/echo/v4"
)

func main() {

	database, err := db.New()
	if err != nil {
		log.Fatal(err)
	}
	defer database.Pool.Close()

	repo := repository.New(
		database.Pool,
	)

	h := handler.New(
		repo,
	)

	e := echo.New()

	e.Use(
		middleware.Logger(),
	)

	e.Use(
		middleware.CORS(),
	)

	router.Register(
		e,
		h,
	)

	if err := e.Start(":5555"); err != nil {
		log.Fatal(err)
	}
}
