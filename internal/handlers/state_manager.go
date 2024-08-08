package handlers

import (
	"sync"
	"time"

	"github.com/gofiber/contrib/websocket"
)

type Message struct {
	Msg      string    `json:"msg"`
	Username string    `json:"-"`
	time     time.Time `json:"-"`
}

// manager is the global state of the application for managing websockets
type manager struct {
	mutex *sync.Mutex

	connections map[string]*websocket.Conn
	messages    []*Message
}

func newManager() *manager {
	return &manager{
		connections: map[string]*websocket.Conn{},
		mutex:       &sync.Mutex{},
	}
}

// add adds websocket connection to connection manager
func (m *manager) add(username string, c *websocket.Conn) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.connections[username] = c
}

// addMessage adds message to the `messages` property of the connection manager
func (m *manager) addMessage(msg Message) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.messages = append(m.messages, &msg)
}

// remove removes websocket connection from the connection manager
func (m *manager) remove(username string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	delete(m.connections, username)
}

// getConnection retrieve user websocket connection
func (m *manager) getConnection(username string) *websocket.Conn {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	return m.connections[username]
}

// getConnectedUsers get the username/connection map
func (m *manager) getConnectedUsers() map[string]*websocket.Conn {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	return m.connections
}

// getMessages get the list of messages
func (m *manager) getMessages() []*Message {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	return m.messages
}
