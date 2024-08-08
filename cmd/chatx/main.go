package main

import (
	"time"

	"github.com/emarifer/gofiber-realtime-chat-htmx/internal/app"
	"github.com/emarifer/gofiber-realtime-chat-htmx/internal/db"
	"github.com/emarifer/gofiber-realtime-chat-htmx/internal/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/template/html/v2"
)

func main() {
	engine := html.New("web/views", ".tmpl")
	ap := fiber.New(fiber.Config{
		AppName:           "ChatX",
		Views:             engine,
		ViewsLayout:       "layouts/main",
		ErrorHandler:      handlers.CustomErrorHandler,
		PassLocalsToViews: true,
	})
	db, storage := db.InitDb()
	ss := session.New(session.Config{
		Storage:    storage,
		Expiration: time.Minute * 60,
	})

	// Dependency injection
	newApp := app.NewApp(ap, db, ss)

	newApp.App.Use(logger.New())

	handlers.SetupRoutes(newApp)

	newApp.Start()
}
