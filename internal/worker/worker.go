package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rocketblend/rocketblend/pkg/types"
)

type (
	Config struct {
		ServerURL  string
		Hostname   string
		WorkingDir string
		Logger     types.Logger
	}

	Option func(*Config)

	Worker struct {
		config Config
		conn   *websocket.Conn
	}
)

func WithLogger(logger types.Logger) Option {
	return func(c *Config) {
		c.Logger = logger
	}
}

func WithServerURL(url string) Option {
	return func(c *Config) {
		c.ServerURL = url
	}
}

func WithWorkingDir(dir string) Option {
	return func(c *Config) {
		c.WorkingDir = dir
	}
}

func WithHostname(name string) Option {
	return func(c *Config) {
		c.Hostname = name
	}
}

func New(opts ...Option) *Worker {
	hostname, _ := os.Hostname()
	config := Config{
		ServerURL: "http://localhost:8080",
		Hostname:  hostname,
	}

	for _, opt := range opts {
		opt(&config)
	}

	return &Worker{
		config: config,
	}
}

func (w *Worker) Start(ctx context.Context) error {
	url := fmt.Sprintf("%s/ws/worker", w.config.ServerURL)
	if w.config.ServerURL[:4] == "http" {
		url = "ws" + w.config.ServerURL[4:] + "/ws/worker"
	}

	w.config.Logger.Info("connecting to server", map[string]any{"url": url})
	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return err
	}
	w.conn = c
	defer c.Close()

	capsPayload, _ := json.Marshal(types.WorkerCapabilities{
		Hostname: w.config.Hostname,
	})
	if err := w.send("handshake", capsPayload); err != nil {
		return err
	}

	w.config.Logger.Info("worker connected", nil)

	go func() {
		<-ctx.Done()
		w.config.Logger.Info("worker shutting down", nil)
		c.Close()
	}()

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}

			w.config.Logger.Error("read error", map[string]any{"err": err})
			return err
		}

		var event types.JobEvent
		if err := json.Unmarshal(message, &event); err != nil {
			continue
		}

		w.config.Logger.Info("received event", map[string]any{"type": event.Type})

		if event.Type == "job" {
			go w.handleJob(event.Payload)
		}
	}
}

func (w *Worker) send(typ string, payload []byte) error {
	evt := types.JobEvent{Type: typ, Payload: payload}
	return w.conn.WriteJSON(evt)
}

func (w *Worker) handleJob(payload []byte) {
	var wrapper struct {
		ID      string           `json:"id"`
		Request types.JobRequest `json:"request"`
	}
	json.Unmarshal(payload, &wrapper)

	w.config.Logger.Info("processing job", map[string]any{
		"id":   wrapper.ID,
		"task": wrapper.Request.Task,
	})

	// Simulate work
	time.Sleep(2 * time.Second)

	res := map[string]string{"message": "Processed successfully by basic worker"}
	resBytes, _ := json.Marshal(res)

	w.sendResult(wrapper.ID, types.JobStatusCompleted, resBytes, "")
}

func (w *Worker) sendResult(id string, status types.JobStatus, res json.RawMessage, errStr string) {
	result := types.JobResult{
		ID:     id,
		Status: status,
		Result: res,
		Error:  errStr,
	}
	p, _ := json.Marshal(result)
	w.send(string(status), p)
	w.config.Logger.Info("sent result", map[string]any{"id": id, "status": status})
}
