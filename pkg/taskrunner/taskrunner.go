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

	RunOpts[T any] struct {
		Tasks          []Task[T]
		Mode           ExecutionMode
		MaxConcurrency int
	}

	Task[T any] func(ctx context.Context) (T, error)

	indexedResult[T any] struct {
		Index  int
		Result T
	}
)

var (
	ErrNegativeConcurrency = errors.New("max concurrency must be greater than or equal to 0")
	ErrNoTasks             = errors.New("no tasks to execute")
	ErrInvalidMode         = errors.New("invalid execution mode")
)

// Run executes the tasks based on the provided options.
func Run[T any](ctx context.Context, opts *RunOpts[T]) ([]T, error) {
	if err := validateOpts(opts); err != nil {
		return nil, err
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
		return nil, ErrInvalidMode
	}
}

// validateOpts checks the options for any configuration errors before execution.
func validateOpts[T any](opts *RunOpts[T]) error {
	if opts.MaxConcurrency < 0 {
		return ErrNegativeConcurrency
	}

	if len(opts.Tasks) == 0 {
		return ErrNoTasks
	}

	return nil
}

// runSequentially executes the tasks sequentially.
func runSequentially[T any](ctx context.Context, tasks []Task[T]) ([]T, error) {
	results := make([]T, 0, len(tasks))
	for _, task := range tasks {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			result, err := task(ctx)
			if err != nil {
				return nil, err
			}
			results = append(results, result)
		}
	}

	return results, nil
}

// runConcurrently executes the tasks concurrently.
func runConcurrently[T any](ctx context.Context, tasks []Task[T]) ([]T, error) {
	var wg sync.WaitGroup
	errChan := make(chan error, len(tasks))
	resultChan := make(chan indexedResult[T], len(tasks))

	for i, task := range tasks {
		wg.Add(1)
		go func(index int, t Task[T]) {
			defer wg.Done()
			select {
			case <-ctx.Done():
				errChan <- ctx.Err()
			default:
				result, err := t(ctx)
				if err != nil {
					errChan <- err
				} else {
					resultChan <- indexedResult[T]{index, result}
				}
			}
		}(i, task)
	}

	wg.Wait()
	close(errChan)
	close(resultChan)

	if err := collectErrors(errChan); err != nil {
		return nil, err
	}

	return collectResults(resultChan, len(tasks)), nil
}

// runWithControlledConcurrency executes the tasks concurrently with controlled concurrency.
func runWithControlledConcurrency[T any](ctx context.Context, tasks []Task[T], maxConcurrency int) ([]T, error) {
	sem := make(chan struct{}, maxConcurrency)
	var wg sync.WaitGroup
	errChan := make(chan error, len(tasks))
	resultChan := make(chan indexedResult[T], len(tasks))

	for i, task := range tasks {
		sem <- struct{}{} // Acquire a token
		wg.Add(1)
		go func(index int, t Task[T]) {
			defer wg.Done()
			defer func() { <-sem }() // Release the token
			select {
			case <-ctx.Done():
				errChan <- ctx.Err()
			default:
				result, err := t(ctx)
				if err != nil {
					errChan <- err
				} else {
					resultChan <- indexedResult[T]{index, result}
				}
			}
		}(i, task)
	}

	wg.Wait()
	close(errChan)
	close(resultChan)

	if err := collectErrors(errChan); err != nil {
		return nil, err
	}

	return collectResults(resultChan, len(tasks)), nil
}

// runWithWorkerPool executes the tasks using a worker pool.
func runWithWorkerPool[T any](ctx context.Context, tasks []Task[T], numWorkers int) ([]T, error) {
	tasksChan := make(chan struct {
		Index int
		Task  Task[T]
	}, len(tasks))
	errChan := make(chan error, len(tasks))
	resultChan := make(chan indexedResult[T], len(tasks))
	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(ctx, tasksChan, &wg, errChan, resultChan)
	}

	// Send tasks to the channel
	for i, task := range tasks {
		tasksChan <- struct {
			Index int
			Task  Task[T]
		}{i, task}
	}
	close(tasksChan) // No more tasks are coming, close the channel

	// Wait for all workers to finish
	wg.Wait()
	close(errChan)
	close(resultChan)

	if err := collectErrors(errChan); err != nil {
		return nil, err
	}

	return collectResults(resultChan, len(tasks)), nil
}

// worker processes tasks from the tasks channel.
func worker[T any](ctx context.Context, tasks <-chan struct {
	Index int
	Task  Task[T]
}, wg *sync.WaitGroup, errChan chan<- error, resultChan chan<- indexedResult[T]) {
	defer wg.Done()
	for task := range tasks {
		select {
		case <-ctx.Done():
			errChan <- ctx.Err()
			return
		default:
			result, err := task.Task(ctx)
			if err != nil {
				errChan <- err
			} else {
				resultChan <- indexedResult[T]{task.Index, result}
			}
		}
	}
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

// collectResults aggregates results from a channel into a slice and orders them.
func collectResults[T any](resultChan <-chan indexedResult[T], taskCount int) []T {
	results := make([]T, taskCount)
	for result := range resultChan {
		results[result.Index] = result.Result
	}

	return results
}
