package workloader

import (
	"context"
	"sync"
)

const (
	Sequential ExecutionMode = iota
	Concurrent
	ControlledConcurrency
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

func Run(ctx context.Context, opts *RunOpts) error {
	switch opts.Mode {
	case Sequential:
		return runSequentially(ctx, opts.Tasks)
	case Concurrent, ControlledConcurrency:
		if opts.MaxConcurrency == 0 {
			return runConcurrently(ctx, opts.Tasks)
		}
		return runWithControlledConcurrency(ctx, opts.Tasks, opts.MaxConcurrency)
	case WorkerPool:
		// For WorkerPool, maxConcurrency represents the number of workers
		return runWithWorkerPool(ctx, opts.Tasks, opts.MaxConcurrency)
	default:
		return runConcurrently(ctx, opts.Tasks)
	}
}

func runSequentially(ctx context.Context, tasks []Task) error {
	for _, task := range tasks {
		if err := task(ctx); err != nil {
			return err
		}
	}
	return nil
}

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

	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

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

	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

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

	// Collect errors, if any
	for err := range errChan {
		if err != nil {
			return err
		}
	}
	return nil
}
