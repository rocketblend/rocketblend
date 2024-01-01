package lockfile

import (
	"context"
	"errors"
	"os"
	"sync"
	"time"

	"github.com/flowshot-io/x/pkg/logger"
)

const (
	ExecutionTimeout = 30 * time.Second // 30 seconds
	HeartBeatTicker  = 15 * time.Second // 15 seconds
)

type Locker struct {
	logger       logger.Logger
	lockFilePath string
	lockFile     *os.File
	ticker       *time.Ticker
	mutex        sync.Mutex
}

type (
	Options struct {
		Path   string
		Logger logger.Logger
	}

	Option func(*Options)
)

func WithPath(path string) Option {
	return func(o *Options) {
		o.Path = path
	}
}

func WithLogger(logger logger.Logger) Option {
	return func(o *Options) {
		o.Logger = logger
	}
}

// New creates a lock and returns a cancel function to release it.
func New(ctx context.Context, opts ...Option) (cancelFunc func(), err error) {
	options := &Options{
		Logger: logger.NoOp(),
	}

	for _, opt := range opts {
		opt(options)
	}

	locker := &Locker{
		logger:       options.Logger,
		lockFilePath: options.Path,
	}

	ctx, cancel := context.WithCancel(ctx)
	if err := locker.lock(ctx); err != nil {
		cancel()
		return nil, err
	}

	// Return a function to cancel the context and release the lock.
	return func() {
		cancel()
	}, nil
}

// TODO: Look into using flock instead of a lock file.
func (l *Locker) lock(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	l.mutex.Lock()
	defer l.mutex.Unlock()

	lockFile, err := l.tryLockFile()
	if err != nil {
		return err
	}

	l.ticker = time.NewTicker(HeartBeatTicker)
	l.lockFile = lockFile

	go l.manageLock(ctx)

	return nil
}

func (l *Locker) tryLockFile() (*os.File, error) {
	lockFile, err := os.OpenFile(l.lockFilePath, os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		if os.IsExist(err) {
			if err := l.handleExistingLock(); err != nil {
				return nil, err
			}

			return os.OpenFile(l.lockFilePath, os.O_CREATE|os.O_EXCL, 0600)
		}

		return nil, err
	}

	return lockFile, nil
}

func (l *Locker) handleExistingLock() error {
	info, err := os.Stat(l.lockFilePath)
	if err != nil {
		return err
	}

	if time.Since(info.ModTime()) > ExecutionTimeout {
		return os.Remove(l.lockFilePath)
	}

	return errors.New("another instance is already holding the lock")
}

func (l *Locker) manageLock(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			l.unlock()
			return
		case <-l.ticker.C:
			// Touch the lock file to update its modified time
			now := time.Now()
			if chtimesErr := os.Chtimes(l.lockFilePath, now, now); chtimesErr != nil {
				l.logger.Error("error updating lock file timestamp", map[string]interface{}{"error": chtimesErr})
			}
		}
	}
}

func (l *Locker) unlock() {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	// Stop the ticker only if it's not already stopped
	if l.ticker != nil {
		l.ticker.Stop()
		l.ticker = nil
	}

	// Close the lock file only if it's not already closed
	if l.lockFile != nil {
		err := l.lockFile.Close()
		if err != nil {
			l.logger.Error("error closing lock file", map[string]interface{}{"error": err})
		}
		l.lockFile = nil
	}

	// Remove the lock file only if it exists
	if _, err := os.Stat(l.lockFilePath); err == nil {
		if err := os.Remove(l.lockFilePath); err != nil {
			l.logger.Error("error removing lock file", map[string]interface{}{"error": err})
		}
	}
}
