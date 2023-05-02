package main

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/n-vh/twitch-renames/internal/config"
	"github.com/n-vh/twitch-renames/internal/database"
	"github.com/n-vh/twitch-renames/internal/router"
)

func main() {
	cfg := config.Get()
	database.Connect(&database.Config{
		Host:     cfg.PgHost,
		Port:     cfg.PgPort,
		User:     cfg.PgUser,
		Password: cfg.PgPassword,
		Database: cfg.PgDatabase,
		Schema:   cfg.PgSchema,
	})

	app := fiber.New()

	// app.Use(timingMiddleware)
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))

	app.Get("/autocomplete/:username", router.AutoComplete)
	app.Get("/search/:username", router.Search)
	app.Get("/user/:id", router.User)

	app.Listen(":" + cfg.ServerPort)
}

func timingMiddleware(ctx *fiber.Ctx) error {
	start := time.Now()

	err := ctx.Next()

	ip := ctx.Context().RemoteIP().String()
	timing := int(time.Since(start).Milliseconds())

	log.Printf("Request took %dms, IP: %s", timing, ip)

	return err
}
