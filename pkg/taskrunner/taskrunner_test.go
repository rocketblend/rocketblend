package taskrunner_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/rocketblend/rocketblend/pkg/taskrunner"
)

func TestRunSequential(t *testing.T) {
	tasks := []taskrunner.Task[int]{
		func(ctx context.Context) (int, error) { time.Sleep(10 * time.Millisecond); return 1, nil },
		func(ctx context.Context) (int, error) { return 2, nil },
		func(ctx context.Context) (int, error) { return 3, nil },
	}

	opts := &taskrunner.RunOpts[int]{
		Tasks: tasks,
		Mode:  taskrunner.Sequential,
	}

	results, err := taskrunner.Run(context.Background(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []int{1, 2, 3}
	for i, result := range results {
		if result != expected[i] {
			t.Errorf("expected result %d, got %d", expected[i], result)
		}
	}
}

func TestRunConcurrent(t *testing.T) {
	tasks := []taskrunner.Task[int]{
		func(ctx context.Context) (int, error) { time.Sleep(10 * time.Millisecond); return 1, nil },
		func(ctx context.Context) (int, error) { return 2, nil },
		func(ctx context.Context) (int, error) { return 3, nil },
	}

	opts := &taskrunner.RunOpts[int]{
		Tasks: tasks,
		Mode:  taskrunner.Concurrent,
	}

	results, err := taskrunner.Run(context.Background(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []int{1, 2, 3}
	for i, result := range results {
		if result != expected[i] {
			t.Errorf("expected result %d, got %d", expected[i], result)
		}
	}
}

func TestRunWithControlledConcurrency(t *testing.T) {
	tasks := []taskrunner.Task[int]{
		func(ctx context.Context) (int, error) { time.Sleep(10 * time.Millisecond); return 1, nil },
		func(ctx context.Context) (int, error) { return 2, nil },
		func(ctx context.Context) (int, error) { return 3, nil },
	}

	opts := &taskrunner.RunOpts[int]{
		Tasks:          tasks,
		Mode:           taskrunner.Concurrent,
		MaxConcurrency: 2,
	}

	results, err := taskrunner.Run(context.Background(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []int{1, 2, 3}
	for i, result := range results {
		if result != expected[i] {
			t.Errorf("expected result %d, got %d", expected[i], result)
		}
	}
}

func TestRunWithWorkerPool(t *testing.T) {
	tasks := []taskrunner.Task[int]{
		func(ctx context.Context) (int, error) { time.Sleep(10 * time.Millisecond); return 1, nil },
		func(ctx context.Context) (int, error) { return 2, nil },
		func(ctx context.Context) (int, error) { return 3, nil },
	}

	opts := &taskrunner.RunOpts[int]{
		Tasks:          tasks,
		Mode:           taskrunner.WorkerPool,
		MaxConcurrency: 2,
	}

	results, err := taskrunner.Run(context.Background(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []int{1, 2, 3}
	for i, result := range results {
		if result != expected[i] {
			t.Errorf("expected result %d, got %d", expected[i], result)
		}
	}
}

func TestRunWithErrors(t *testing.T) {
	expectedErr := errors.New("task error")
	tasks := []taskrunner.Task[int]{
		func(ctx context.Context) (int, error) { return 1, nil },
		func(ctx context.Context) (int, error) { return 0, expectedErr },
		func(ctx context.Context) (int, error) { return 3, nil },
	}

	opts := &taskrunner.RunOpts[int]{
		Tasks: tasks,
		Mode:  taskrunner.Concurrent,
	}

	results, err := taskrunner.Run(context.Background(), opts)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error %v, got %v", expectedErr, err)
	}
	if results != nil {
		t.Fatalf("expected nil results, got %v", results)
	}
}

func TestRunWithCancellation(t *testing.T) {
	tasks := []taskrunner.Task[int]{
		func(ctx context.Context) (int, error) { time.Sleep(100 * time.Millisecond); return 1, nil },
		func(ctx context.Context) (int, error) { return 2, nil },
		func(ctx context.Context) (int, error) { return 3, nil },
	}

	opts := &taskrunner.RunOpts[int]{
		Tasks: tasks,
		Mode:  taskrunner.Concurrent,
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	results, err := taskrunner.Run(ctx, opts)
	if err == nil {
		t.Fatalf("expected error due to cancellation, got nil")
	}

	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled error, got %v", err)
	}

	if results != nil {
		t.Fatalf("expected nil results, got %v", results)
	}
}

func TestRunNoTasks(t *testing.T) {
	opts := &taskrunner.RunOpts[int]{
		Tasks: []taskrunner.Task[int]{},
		Mode:  taskrunner.Concurrent,
	}

	results, err := taskrunner.Run(context.Background(), opts)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if !errors.Is(err, taskrunner.ErrNoTasks) {
		t.Fatalf("expected ErrNoTasks error, got %v", err)
	}

	if results != nil {
		t.Fatalf("expected nil results, got %v", results)
	}
}

func TestRunNegativeConcurrency(t *testing.T) {
	tasks := []taskrunner.Task[int]{
		func(ctx context.Context) (int, error) { return 1, nil },
	}

	opts := &taskrunner.RunOpts[int]{
		Tasks:          tasks,
		Mode:           taskrunner.Concurrent,
		MaxConcurrency: -1,
	}

	results, err := taskrunner.Run(context.Background(), opts)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if !errors.Is(err, taskrunner.ErrNegativeConcurrency) {
		t.Fatalf("expected ErrNegativeConcurrency error, got %v", err)
	}

	if results != nil {
		t.Fatalf("expected nil results, got %v", results)
	}
}

func TestRunInvalidMode(t *testing.T) {
	tasks := []taskrunner.Task[int]{
		func(ctx context.Context) (int, error) { return 1, nil },
	}

	opts := &taskrunner.RunOpts[int]{
		Tasks: tasks,
		Mode:  taskrunner.ExecutionMode(999), // Invalid mode
	}

	results, err := taskrunner.Run(context.Background(), opts)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if !errors.Is(err, taskrunner.ErrInvalidMode) {
		t.Fatalf("expected ErrInvalidMode error, got %v", err)
	}

	if results != nil {
		t.Fatalf("expected nil results, got %v", results)
	}
}
