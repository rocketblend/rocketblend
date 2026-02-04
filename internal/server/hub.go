package server

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/rocketblend/rocketblend/pkg/types"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Hub maintains the set of active workers and broadcasts messages
type Hub struct {
	logger        types.Logger
	queue         *JobQueue
	workers       map[*WorkerConnection]bool // All connected workers
	capabilities  map[string][]*WorkerConnection
	projects      map[string][]*WorkerConnection
	register      chan *WorkerConnection
	unregister    chan *WorkerConnection
	jobDispatch   chan dispatchReq
	statusUpdates chan statusUpdateReq
	mu            sync.RWMutex
}

type WorkerConnection struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
	caps types.WorkerCapabilities
}

type dispatchReq struct {
	job *types.JobResult
	req types.JobRequest
}

type statusUpdateReq struct {
	jobID  string
	status types.JobStatus
	result json.RawMessage
	error  string
}

func NewHub(l types.Logger, q *JobQueue) *Hub {
	return &Hub{
		logger:        l,
		queue:         q,
		workers:       make(map[*WorkerConnection]bool),
		capabilities:  make(map[string][]*WorkerConnection),
		projects:      make(map[string][]*WorkerConnection),
		register:      make(chan *WorkerConnection),
		unregister:    make(chan *WorkerConnection),
		jobDispatch:   make(chan dispatchReq),
		statusUpdates: make(chan statusUpdateReq),
	}
}

func (h *Hub) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return

		case client := <-h.register:
			h.mu.Lock()
			h.workers[client] = true
			h.logger.Info("worker connected", map[string]interface{}{"host": client.conn.RemoteAddr().String()})
			h.mu.Unlock()

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.workers[client]; ok {
				delete(h.workers, client)
				close(client.send)
				// Remove from indexes... (simplified for now)
				h.logger.Info("worker disconnected", nil)
			}
			h.mu.Unlock()

		case update := <-h.statusUpdates:
			h.queue.Update(update.jobID, update.status, update.result, update.error)
			h.logger.Info("job updated", map[string]interface{}{
				"id":     update.jobID,
				"status": update.status,
			})

		case req := <-h.jobDispatch:
			h.mu.RLock()
			var chosen *WorkerConnection
			for w := range h.workers {
				// Just pick the first available worker for now.
				// Later we will check if the worker has the project.
				chosen = w
				break
			}
			h.mu.RUnlock()

			if chosen != nil {
				// Send Job
				payload, _ := json.Marshal(types.JobEvent{
					Type: "job",
					Payload: mustMarshal(map[string]interface{}{
						"id":      req.job.ID,
						"request": req.req, // Worker needs the full request
					}),
				})
				chosen.send <- payload

				h.queue.Update(req.job.ID, types.JobStatusAssigned, nil, "")
			} else {
				h.logger.Warn("no worker found for job", map[string]interface{}{"id": req.job.ID})
			}
		}
	}
}

func (h *Hub) Dispatch(job *types.JobResult, req types.JobRequest) {
	h.jobDispatch <- dispatchReq{job, req}
}

func (h *Hub) WorkerCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.workers)
}

func (s *Server) handleWorkerWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.config.Logger.Error("ws upgrade failed", map[string]interface{}{"err": err})
		return
	}

	client := &WorkerConnection{hub: s.hub, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client

	// Start pump loops
	go client.writePump()
	go client.readPump()
}

func (c *WorkerConnection) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			break
		}

		var event types.JobEvent
		if err := json.Unmarshal(message, &event); err != nil {
			continue
		}

		switch event.Type {
		case "handshake":
			var caps types.WorkerCapabilities
			json.Unmarshal(event.Payload, &caps)
			c.caps = caps
			// Update indexes in Hub...

		case "completed", "failed":
			var res types.JobResult // Partial result from worker
			json.Unmarshal(event.Payload, &res)
			c.hub.statusUpdates <- statusUpdateReq{
				jobID:  res.ID,
				status: res.Status,
				result: res.Result,
				error:  res.Error,
			}
		}
	}
}

func (c *WorkerConnection) writePump() {
	defer c.conn.Close()
	for {
		select {
		case msg, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.conn.WriteMessage(websocket.TextMessage, msg)
		}
	}
}

func mustMarshal(v interface{}) []byte {
	b, _ := json.Marshal(v)
	return b
}
