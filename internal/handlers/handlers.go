package handlers

import (
	"fmt"
	"strings"

	"github.com/emarifer/gofiber-realtime-chat-htmx/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"golang.org/x/crypto/bcrypt"
)

/********** Auth Handlers **********/

type authService interface {
	CreateUser(user, pass string) error
	CheckUser(username string) (services.User, error)
}

func newAuthHandler(us authService, ss *session.Store) *authHandler {
	return &authHandler{userServices: us, store: ss}
}

type authHandler struct {
	userServices authService
	store        *session.Store
}

func (ah *authHandler) homeHandler(c *fiber.Ctx) error {
	user, _ := c.Locals(username_key).(string)
	errMsg, succMsg := getMessages(c)

	return c.Render("index", fiber.Map{
		"Title":   fmt.Sprintf("Home (%s's chat)", user),
		"errMsg":  errMsg,
		"succMsg": succMsg,
	})
}

func (ah *authHandler) registerHandler(c *fiber.Ctx) error {
	errMsg, succMsg := getMessages(c)

	return c.Render("register", fiber.Map{
		"Title":   "Register",
		"errMsg":  errMsg,
		"succMsg": succMsg,
	})
}

func (ah *authHandler) signupHandler(c *fiber.Ctx) error {
	user := strings.Trim(c.FormValue("username"), " ")
	pass := strings.Trim(c.FormValue("password"), " ")

	// simple server-side validation...
	if user == "" || pass == "" {
		return setFlash(
			c, "error", []byte("fields cannot be empty"),
		).Redirect("/register?error=fields cannot be empty-400")
	}

	if err := ah.userServices.CreateUser(user, pass); err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return setFlash(
				c, "error", []byte("user already exists"),
			).Redirect("/register?error=user already exists-403")
		}

		return fiber.NewError(
			fiber.StatusInternalServerError,
			"failed to save user to db: database temporarily out of service",
		)
	}

	return setFlash(
		c, "success", []byte("user registered!!"),
	).Redirect("/login?success=user registered-201")
}

func (ah *authHandler) loginHandler(c *fiber.Ctx) error {
	errMsg, succMsg := getMessages(c)

	return c.Render("login", fiber.Map{
		"Title":   "Login",
		"errMsg":  errMsg,
		"succMsg": succMsg,
	})
}

func (ah *authHandler) signinHandler(c *fiber.Ctx) error {
	next := c.Query("next", "/")

	// obtaining the time zone from the POST request of the login form
	tzone := ""
	if len(c.GetReqHeaders()["X-Timezone"]) != 0 {
		tzone = c.GetReqHeaders()["X-Timezone"][0]
		// fmt.Println("Tzone:", tzone)
	}

	username := strings.Trim(c.FormValue("username"), " ")
	pass := strings.Trim(c.FormValue("password"), " ")
	// simple server-side validation...
	if username == "" || pass == "" {
		return setFlash(
			c, "error", []byte("fields cannot be empty"),
		).Redirect("/login?error=fields cannot be empty-400")
	}

	user, err := ah.userServices.CheckUser(username)
	if err != nil {
		if strings.Contains(err.Error(), "no such table: users") {
			return fiber.NewError(
				fiber.StatusInternalServerError,
				"failed to save user to db: database temporarily out of service",
			)
		}

		return setFlash(
			c, "error", []byte("invalid credentials"),
		).Redirect("/login?error=invalid credentials-401")
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.Password), []byte(pass),
	); err != nil {
		return setFlash(
			c, "error", []byte("invalid credentials"),
		).Redirect("/login?error=invalid credentials-401")
	}

	sess, err := ah.store.Get(c)
	if err != nil {
		// return c.Redirect("/login?error=failed to fetch session 500")

		return fiber.NewError(
			fiber.StatusInternalServerError,
			"failed to fetch session: database temporarily out of service",
		)
	}

	if sess.Fresh() {
		sess.Set(username_key, user.Username)
		sess.Set(tzone_key, tzone)

		if err := sess.Save(); err != nil {
			// return c.Redirect("/login?error=failed to save session 500")

			return fiber.NewError(
				fiber.StatusInternalServerError,
				"failed to fetch session: database temporarily out of service",
			)
		}

		// return c.Redirect("/")
		return setFlash(
			c, "success", []byte("you have successfully logged in!!"),
		).Redirect("/")
	}

	return c.Redirect(next)
}

func (ah *authHandler) logoutHandler(c *fiber.Ctx) error {
	user, _ := c.Locals(username_key).(string)
	sess, err := ah.store.Get(c)
	if err != nil {

		return c.Redirect("/login")
	}

	sess.Delete(user)
	sess.Destroy()

	return setFlash(
		c, "success", []byte("you have successfully logged out!!"),
	).Redirect("/login")
}
