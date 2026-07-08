package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"parser/internal/config"
	"parser/internal/handler"
	"parser/internal/repository"
	"parser/internal/router"
	"parser/internal/server"
	"parser/internal/middleware"
	"time"

	"github.com/labstack/echo/v4"
)

const DefaultContextTimeout = 30

func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(
			"failed to load config: ",
			err,
		)
	}

	srv, err := server.New(cfg)
	if err != nil {
		log.Fatal(
			"failed to initialize server: ",
			err,
		)
	}

	repo := repository.New(
		srv.DB.Pool,
		cfg,
	)

	h := handler.New(
		repo,
	)

	e := echo.New()

	e.Use(
		middleware.CORS(cfg),
		middleware.Logger(),
	)

	router.Register(
		e,
		h,
	)

	srv.SetupHTTPServer(e)

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
	)

	go func() {

		if err := srv.Start(); err != nil &&
			!errors.Is(
				err,
				http.ErrServerClosed,
			) {

			log.Fatal(
				"failed to start server: ",
				err,
			)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel :=
		context.WithTimeout(
			context.Background(),
			DefaultContextTimeout*time.Second,
		)

	defer cancel()
	defer stop()

	if err := srv.Shutdown(
		shutdownCtx,
	); err != nil {

		log.Fatal(
			"server forced to shutdown: ",
			err,
		)
	}

	log.Println(
		"server exited properly",
	)
}