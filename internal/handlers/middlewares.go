package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func (ah *authHandler) auth(c *fiber.Ctx) error {
	next := c.Path()
	s, err := ah.store.Get(c)
	if err != nil {
		// return c.Redirect("/login")
		return fiber.NewError(
			fiber.StatusInternalServerError,
			"failed to fetch session: database temporarily out of service",
		)
	}

	if len(s.Keys()) > 0 {
		user := s.Get(username_key)
		if user == nil {
			return c.Redirect(fmt.Sprintf("/login?next=%s", next))
		}

		tzone, _ := s.Get(tzone_key).(string)
		c.Locals(username_key, user.(string))
		c.Locals(tzone_key, tzone)
	} else {
		return c.Redirect(fmt.Sprintf("/login?next=%s", next))
	}

	return c.Next()
}

func (ah *authHandler) flag(c *fiber.Ctx) error {
	s, err := ah.store.Get(c)
	if err != nil {
		return fiber.NewError(
			fiber.StatusInternalServerError,
			"failed to fetch session: database temporarily out of service",
		)
	}

	user := s.Get(username_key)

	if user == nil {
		c.Locals(from_protected_key, false)

		return c.Next()
	}

	c.Locals(from_protected_key, true)

	return c.Next()
}
