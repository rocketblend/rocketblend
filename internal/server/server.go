package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/rocketblend/rocketblend/pkg/types"
	"github.com/rs/xid"
)

type (
	Config struct {
		Port   int
		Host   string
		Logger types.Logger
	}

	Option func(*Config)

	Server struct {
		config Config
		router *mux.Router
		hub    *Hub
		queue  *JobQueue
	}
)

func WithLogger(logger types.Logger) Option {
	return func(c *Config) {
		c.Logger = logger
	}
}

func WithHost(host string) Option {
	return func(c *Config) {
		c.Host = host
	}
}

func WithPort(port int) Option {
	return func(c *Config) {
		c.Port = port
	}
}

func New(opts ...Option) *Server {
	config := Config{
		Port: 8080,
		Host: "0.0.0.0",
	}

	for _, opt := range opts {
		opt(&config)
	}

	s := &Server{
		config: config,
		router: mux.NewRouter(),
		queue:  NewJobQueue(),
	}
	s.hub = NewHub(config.Logger, s.queue)
	s.routes()
	return s
}

func (s *Server) Start(ctx context.Context) error {
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
	s.config.Logger.Info("listening", map[string]interface{}{"address": addr})

	srv := &http.Server{
		Addr:    addr,
		Handler: s.router,
	}

	go s.hub.Run(ctx)

	go func() {
		<-ctx.Done()
		s.config.Logger.Info("server shutting down", nil)
		srv.Shutdown(context.Background())
	}()

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (s *Server) routes() {
	api := s.router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/jobs", s.handlePostJob).Methods("POST")
	api.HandleFunc("/jobs/{id}", s.handleGetJob).Methods("GET")
	api.HandleFunc("/system/info", s.handleSystemInfo).Methods("GET") // Metadata endpoint

	s.router.HandleFunc("/ws/worker", s.handleWorkerWS)
}

func (s *Server) handlePostJob(w http.ResponseWriter, r *http.Request) {
	var req types.JobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	jobID := xid.New().String()
	job := &types.JobResult{
		ID:        jobID,
		Status:    types.JobStatusQueued,
		CreatedAt: time.Now(),
	}

	// Add to queue
	s.queue.Add(job, req)

	// Notify Hub to dispatch
	s.hub.Dispatch(job, req)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"job_id": jobID})
}

func (s *Server) handleGetJob(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	job, exists := s.queue.Get(id)
	if !exists {
		http.Error(w, "Job not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(job)
}

func (s *Server) handleSystemInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	info := map[string]interface{}{
		"workers": s.hub.WorkerCount(),
	}
	json.NewEncoder(w).Encode(info)
}

// In-Memory Queue (simplified)
type JobQueue struct {
	jobs sync.Map // map[string]*types.JobResult
	reqs sync.Map // map[string]types.JobRequest
}

func NewJobQueue() *JobQueue {
	return &JobQueue{}
}

func (q *JobQueue) Add(job *types.JobResult, req types.JobRequest) {
	q.jobs.Store(job.ID, job)
	q.reqs.Store(job.ID, req)
}

func (q *JobQueue) Get(id string) (*types.JobResult, bool) {
	val, ok := q.jobs.Load(id)
	if !ok {
		return nil, false
	}
	return val.(*types.JobResult), true
}

func (q *JobQueue) Update(id string, status types.JobStatus, result json.RawMessage, err string) {
	val, ok := q.jobs.Load(id)
	if !ok {
		return
	}
	j := val.(*types.JobResult)
	j.Status = status
	j.Result = result
	j.Error = err
	if status == types.JobStatusProcessing {
		now := time.Now()
		j.StartedAt = &now
	}
	if status == types.JobStatusCompleted || status == types.JobStatusFailed {
		now := time.Now()
		j.EndedAt = &now
	}
	q.jobs.Store(id, j)
}
