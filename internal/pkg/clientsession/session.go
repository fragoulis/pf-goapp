package clientsession

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"goapp/internal/pkg/watcher"

	"github.com/gorilla/websocket"
)

// Session is a single WebSocket client connection.
type Session struct {
	id         int
	connection *websocket.Conn
}

// New creates a WebSocket client session.
func New(id int) *Session {
	return &Session{
		id: id,
	}
}

// Start connects to the server and logs counter values until the session ends.
func (s *Session) Start(ctx context.Context) error {
	connection, _, err := websocket.DefaultDialer.DialContext(
		ctx,
		"ws://localhost:8080/goapp/ws",
		http.Header{
			"Origin": []string{
				// Can be replaced with the actual origin if needed or a list of allowed origins.
				"http://localhost:8080",
			},
		},
	)
	if err != nil {
		return fmt.Errorf("connection %d: connect: %w", s.id, err)
	}
	s.connection = connection
	defer connection.Close() //nolint:errcheck

	doneCh := make(chan struct{})
	defer close(doneCh)

	go func() {
		select {
		case <-ctx.Done(): // exit on ctrl+c
			s.close()
		case <-doneCh: // exit on client failure
		}
	}()

	for {
		if err := s.read(ctx); err != nil {
			return err
		}
	}
}

func (s *Session) read(ctx context.Context) error {
	_, message, err := s.connection.ReadMessage()
	if err != nil {

		if ctx.Err() != nil {
			return ctx.Err()
		}
		return fmt.Errorf("connection %d: read: %w", s.id, err)
	}

	var counter watcher.Counter
	if err := json.Unmarshal(message, &counter); err != nil {
		return fmt.Errorf("connection %d: decode response: %w", s.id, err)
	}

	log.Printf("[conn #%d] iteration: %d, value: %s", s.id, counter.Iteration, counter.Value)

	return nil
}

func (s *Session) close() {
	_ = s.connection.WriteControl(
		websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
		time.Now().Add(time.Second),
	)
	_ = s.connection.Close()
}
