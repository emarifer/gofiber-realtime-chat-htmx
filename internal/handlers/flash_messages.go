package handlers

import (
	"encoding/base64"
	"time"

	"github.com/gofiber/fiber/v2"
)

// SetFlash sets a cookie with the flash message (base64 encoded)
// which is made available for the next request.
func setFlash(c *fiber.Ctx, name string, value []byte) *fiber.Ctx {
	// Create cookie
	cookie := new(fiber.Cookie)
	cookie.Name = name
	cookie.Value = encode(value)

	c.Cookie(cookie)

	return c
}

// getFlash tries to retrieve the message (as a []byte)
// set in a cookie by a handler in the previous request,
// decodes it, and overrides the cookie so that
// the message can only be read once.
func getFlash(c *fiber.Ctx, name string) ([]byte, error) {
	value, err := decode(c.Cookies(name, ""))
	if err != nil {
		return nil, err
	}

	c.Cookie(&fiber.Cookie{
		Name: name,
		Path: "/",
		// Set expiry date to the past
		MaxAge:  -1,
		Expires: time.Unix(1, 0),
	})

	// dc := &http.Cookie{
	// 	Name:    name,
	// 	Path:    "/",
	// 	MaxAge:  -1,
	// 	Expires: time.Unix(1, 0),
	// }

	return value, nil
}

// GetMessages is a convenience function that uses the getFlash function,
// ignoring the errors it may return and transforming
// the byte slices of the error/success messages into strings.
// If the slices are <nil>, return their respective empty strings.
func getMessages(c *fiber.Ctx) (string, string) {
	fmErr, _ := getFlash(c, "error")
	fmSucc, _ := getFlash(c, "success")

	return string(fmErr), string(fmSucc)
}

// -------------------------

func encode(src []byte) string {
	return base64.URLEncoding.EncodeToString(src)
}

func decode(src string) ([]byte, error) {
	return base64.URLEncoding.DecodeString(src)
}

/* SIMPLE FLASH MESSAGES IN GO:
https://www.alexedwards.net/blog/simple-flash-messages-in-golang
*/
