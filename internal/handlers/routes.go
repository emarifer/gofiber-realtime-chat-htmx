package handlers

import (
	"errors"
	"fmt"

	"github.com/emarifer/gofiber-realtime-chat-htmx/internal/app"
	"github.com/emarifer/gofiber-realtime-chat-htmx/internal/services"

	"github.com/gofiber/fiber/v2"
)

const (
	username_key       string = "username"
	from_protected_key string = "fromProtected"
	is_error_key       string = "isError"
	code_key           string = "errorCode"
	tzone_key          string = "time_zone"
)

func SetupRoutes(a *app.App) {
	// Dependency injection
	us := services.NewUserService(services.User{}, a.DB)
	ah := newAuthHandler(us, a.Store)

	sm := newManager()
	wh := newWebsocketHandler(sm)

	a.App.Get("/", ah.auth, ah.homeHandler)
	a.App.Get("/register", ah.flag, ah.registerHandler)
	a.App.Post("/register", ah.signupHandler)
	a.App.Get("/login", ah.flag, ah.loginHandler)
	a.App.Post("/login", ah.signinHandler)
	a.App.Post("/logout", ah.auth, ah.logoutHandler)

	a.App.Get("/ws", ah.auth, wh.chatHandler())
}

// CustomErrorHandler does centralized error handling.
func CustomErrorHandler(c *fiber.Ctx, err error) error {
	c.Locals(is_error_key, true)

	// Status code defaults to 500
	code := fiber.StatusInternalServerError

	// Retrieve the custom status code if it's a *fiber.Error
	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
	}

	c.Locals(code_key, code)

	data := fiber.Map{"Title": fmt.Sprintf("Error %d", code)}

	// Send custom error page
	err = c.Status(code).Render("errors", data)
	if err != nil {
		// In case the Render fails
		return c.
			Status(fiber.StatusInternalServerError).
			SendString("Internal Server Error")
	}

	// Return from handler
	return nil
}
