package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rocketblend/rocketblend/internal/orchestration"
	"github.com/rocketblend/rocketblend/internal/orchestration/worker/supervisor"
	"github.com/rocketblend/rocketblend/pkg/types"
)

type (
	Config struct {
		ServerURL  string
		Hostname   string
		WorkingDir string
		Driver     types.Driver
		Logger     types.Logger
	}

	Option func(*Config)

	Worker struct {
		config     Config
		conn       *websocket.Conn
		supervisor *supervisor.Supervisor
	}
)

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

func WithDriver(d types.Driver) Option {
	return func(c *Config) {
		c.Driver = d
	}
}

func WithLogger(logger types.Logger) Option {
	return func(c *Config) {
		c.Logger = logger
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
		config:     config,
		supervisor: supervisor.New(config.Logger, config.Driver),
	}
}

func (w *Worker) Start(ctx context.Context) error {
	if err := w.supervisor.Discover(ctx, w.config.WorkingDir); err != nil {
		w.config.Logger.Error("failed to discover projects", map[string]any{"err": err})
	}

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

	projects := w.supervisor.Registered()
	capsPayload, _ := json.Marshal(orchestration.WorkerCapabilities{
		Hostname: w.config.Hostname,
		Projects: projects,
	})
	if err := w.send(orchestration.JobEventHandshake, capsPayload); err != nil {
		return err
	}

	w.config.Logger.Info("handshake sent", map[string]any{
		"hostname": w.config.Hostname,
		"projects": len(projects),
	})

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
				w.config.Logger.Info("worker shutdown requested", nil)
				return nil
			}

			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				w.config.Logger.Info("server closed connection", nil)
				return nil
			}

			w.config.Logger.Error("read error", map[string]any{"err": err})
			return err
		}

		var event orchestration.JobEvent
		if err := json.Unmarshal(message, &event); err != nil {
			continue
		}

		w.config.Logger.Info("received event", map[string]any{"type": event.Type})

		if event.Type == orchestration.JobEventJob {
			go w.handleJob(event.Payload)
		}
	}
}

func (w *Worker) send(typ orchestration.JobEventType, payload []byte) error {
	evt := orchestration.JobEvent{Type: typ, Payload: payload}
	return w.conn.WriteJSON(evt)
}

func (w *Worker) handleJob(payload []byte) {
	var wrapper struct {
		ID      string                   `json:"id"`
		Request orchestration.JobRequest `json:"request"`
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

	w.sendResult(wrapper.ID, orchestration.JobStatusCompleted, resBytes, "")
}

func (w *Worker) sendResult(id string, status orchestration.JobStatus, res json.RawMessage, errStr string) {
	result := orchestration.JobResult{
		ID:     id,
		Status: status,
		Result: res,
		Error:  errStr,
	}
	p, _ := json.Marshal(result)

	// Map JobStatus to JobEventType
	var eventType orchestration.JobEventType
	if status == orchestration.JobStatusCompleted {
		eventType = orchestration.JobEventCompleted
	} else {
		eventType = orchestration.JobEventFailed
	}

	w.send(eventType, p)
	w.config.Logger.Info("sent result", map[string]any{"id": id, "status": status})
}
