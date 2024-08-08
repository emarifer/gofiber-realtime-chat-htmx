package handlers

import (
	"bytes"
	"encoding/json"
	"html/template"
	"log"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func newWebsocketHandler(sm *manager) *websocketHandler {
	return &websocketHandler{stateManager: sm}
}

type websocketHandler struct {
	stateManager *manager
}

func (wh *websocketHandler) chatHandler() fiber.Handler {

	return websocket.New(func(c *websocket.Conn) {
		user, _ := c.Locals(username_key).(string)
		tz, _ := c.Locals(tzone_key).(string)
		defer func() {
			wh.stateManager.remove(user)
			c.Close()
			log.Println("remove user:", user)
		}()
		log.Println("logged in user:", user)

		var (
			mt  int
			msg []byte
			err error
		)

		wh.stateManager.add(user, c)

		// Display all messages when connecting to websocket
		for _, m := range wh.stateManager.getMessages() {
			if err = wh.stateManager.getConnection(user).WriteMessage(
				websocket.TextMessage,
				getMessageTemplate(user, m.Username, tz, m),
			); err != nil {
				log.Println("write:", err)
				break
			}
		}
		// for _, ms := range mg.GetMessages() {
		// 	fmt.Printf("Messages: %#v\n", ms)
		// }

		for {
			if mt, msg, err = c.ReadMessage(); err != nil {
				log.Println("read:", err)
				break
			}
			// The field of the `Message` struct that stores
			// the message text has to match the name of
			// the form input on the `Home` page (index.tmpl file)
			var m Message
			if err := json.Unmarshal(msg, &m); err != nil {
				log.Println(err)
			}
			m.Username = user
			m.time = time.Now().UTC()
			log.Printf("Decoded MSG: %#v\n", m)
			wh.stateManager.addMessage(m)

			// broadcast the message to all connected users
			for u, conn := range wh.stateManager.getConnectedUsers() {
				if err = conn.WriteMessage(
					mt, getMessageTemplate(u, user, tz, &m),
				); err != nil {
					log.Println("write:", err)
					break
				}
			}
		}
	})
}

// getMessageTemplate generates the template for
// the message and returns a byte portion representing the template,
// which will be sent over the websocket.
func getMessageTemplate(u, currentUser, tz string, m *Message) []byte {
	var itsme bool
	if u == currentUser {
		itsme = true
	}

	tmpl, err := template.ParseFiles("web/views/message.tmpl")
	if err != nil {
		log.Fatalf("template parsing: %s", err)
	}

	data := fiber.Map{
		"ItsMe":    itsme,
		"Username": currentUser,
		"Text":     m.Msg,
		"Time":     convertTime(tz, m.time),
	}

	// Render the template with the message as data.
	var renderedMessage bytes.Buffer
	err = tmpl.Execute(&renderedMessage, data)
	if err != nil {
		log.Fatalf("template execution: %s", err)
	}

	return renderedMessage.Bytes()
}

// convertTime is a convenience function that
// returns the appropriately formatted time
// given the UTC datetime and time zone.
func convertTime(tz string, t time.Time) string {
	loc, _ := time.LoadLocation(tz)

	return t.In(loc).Format("15:04 -0700")
}
