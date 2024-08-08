package app

import (
	"database/sql"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

type App struct {
	App   *fiber.App
	DB    *sql.DB
	Store *session.Store
}

func NewApp(a *fiber.App, d *sql.DB, s *session.Store) *App {
	return &App{App: a, DB: d, Store: s}
}

func (a *App) Start() {
	a.App.Static("/static", "web/static", fiber.Static{Compress: true})

	// 404 Handler
	a.App.Use(func(c *fiber.Ctx) error {
		return fiber.NewError(fiber.StatusNotFound, "page not found")
	})

	log.Fatalln(a.App.Listen(":8080"))
}
