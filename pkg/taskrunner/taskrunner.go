package taskrunner

import (
	"context"
	"errors"
	"sync"
)

const (
	Sequential ExecutionMode = iota
	Concurrent
	WorkerPool
)

type (
	ExecutionMode int

	RunOpts struct {
		Tasks          []Task
		Mode           ExecutionMode
		MaxConcurrency int
	}

	Task func(ctx context.Context) error
)

var (
	ErrNegativeConcurrency = errors.New("max concurrency must be greater than or equal to 0")
	ErrNoTasks             = errors.New("no tasks to execute")
	ErrInvalidMode         = errors.New("invalid execution mode")
)

// Run executes the tasks based on the provided options.
func Run(ctx context.Context, opts *RunOpts) error {
	if err := validateOpts(opts); err != nil {
		return err
	}

	if len(opts.Tasks) == 1 {
		return opts.Tasks[0](ctx)
	}

	switch opts.Mode {
	case Sequential:
		return runSequentially(ctx, opts.Tasks)
	case Concurrent:
		if opts.MaxConcurrency > 0 {
			return runWithControlledConcurrency(ctx, opts.Tasks, opts.MaxConcurrency)
		}

		return runConcurrently(ctx, opts.Tasks)
	case WorkerPool:
		return runWithWorkerPool(ctx, opts.Tasks, opts.MaxConcurrency)
	default:
		return ErrInvalidMode
	}
}

// validateOpts checks the options for any configuration errors before execution.
func validateOpts(opts *RunOpts) error {
	if opts.MaxConcurrency < 0 {
		return ErrNegativeConcurrency
	}

	if len(opts.Tasks) == 0 {
		return ErrNoTasks
	}

	return nil
}

// runSequentially executes the tasks sequentially.
func runSequentially(ctx context.Context, tasks []Task) error {
	for _, task := range tasks {
		if err := task(ctx); err != nil {
			return err
		}
	}

	return nil
}

// runConcurrently executes the tasks concurrently.
func runConcurrently(ctx context.Context, tasks []Task) error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(tasks))

	for _, task := range tasks {
		wg.Add(1)
		go func(t Task) {
			defer wg.Done()
			if err := t(ctx); err != nil {
				errChan <- err
			}
		}(task)
	}

	wg.Wait()
	close(errChan)

	return collectErrors(errChan)
}

// runWithControlledConcurrency executes the tasks concurrently with controlled concurrency.
func runWithControlledConcurrency(ctx context.Context, tasks []Task, maxConcurrency int) error {
	sem := make(chan struct{}, maxConcurrency)
	var wg sync.WaitGroup
	errChan := make(chan error, len(tasks))

	for _, task := range tasks {
		sem <- struct{}{} // Acquire a token
		wg.Add(1)
		go func(t Task) {
			defer wg.Done()
			defer func() { <-sem }() // Release the token
			if err := t(ctx); err != nil {
				errChan <- err
			}
		}(task)
	}

	wg.Wait()
	close(errChan)

	return collectErrors(errChan)
}

// worker is a helper function that executes tasks from a channel.
func worker(ctx context.Context, tasks <-chan Task, wg *sync.WaitGroup, errChan chan<- error) {
	defer wg.Done()
	for task := range tasks {
		select {
		case <-ctx.Done():
			errChan <- ctx.Err()
			return
		default:
			if err := task(ctx); err != nil {
				errChan <- err
			}
		}
	}
}

// runWithWorkerPool executes the tasks using a worker pool.
func runWithWorkerPool(ctx context.Context, tasks []Task, numWorkers int) error {
	tasksChan := make(chan Task, len(tasks))
	errChan := make(chan error, len(tasks))
	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(ctx, tasksChan, &wg, errChan)
	}

	// Send tasks to the channel
	for _, task := range tasks {
		tasksChan <- task
	}
	close(tasksChan) // No more tasks are coming, close the channel

	// Wait for all workers to finish
	wg.Wait()
	close(errChan)

	return collectErrors(errChan)
}

// collectErrors aggregates all errors from the error channel.
func collectErrors(errChan chan error) error {
	var err error
	for e := range errChan {
		if e != nil {
			err = e
		}
	}

	return err
}
