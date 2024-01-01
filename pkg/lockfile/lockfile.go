package lockfile

import (
	"context"
	"fmt"
	"os"
	"time"
)

const (
	ExecutionTimeout = 30 * time.Second // 30 seconds
	HeartBeatTicker  = 15 * time.Second // 15 seconds
)

type Locker struct {
	lockFilePath string
	lockFile     *os.File
	ticker       *time.Ticker
}

func NewLocker(lockFilePath string) *Locker {
	return &Locker{
		lockFilePath: lockFilePath,
	}
}

func (l *Locker) Lock(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	// Try to create the lock file
	lockFile, err := os.OpenFile(l.lockFilePath, os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		if os.IsExist(err) {
			// Some other process has already locked this file
			info, err := os.Stat(l.lockFilePath)
			if err != nil {
				return err
			}
			if time.Since(info.ModTime()) > ExecutionTimeout {
				// The lock is older than execution timeout, remove it.
				os.Remove(l.lockFilePath)
				// Try to create the lock file again.
				lockFile, err = os.OpenFile(l.lockFilePath, os.O_CREATE|os.O_EXCL, 0600)
				if err != nil {
					return err
				}
			} else {
				return fmt.Errorf("another instance is already holding the lock")
			}
		} else {
			// Other error
			return err
		}
	}

	// Start a ticker to refresh the lock every heat beat amount.
	l.ticker = time.NewTicker(HeartBeatTicker)
	l.lockFile = lockFile
	go func() {
		defer l.ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				// The operation is done, stop refreshing.
				return
			case <-l.ticker.C:
				// Touch the lock file to update its modified time.
				now := time.Now()
				os.Chtimes(l.lockFilePath, now, now)
			}
		}
	}()

	return nil
}

func (l *Locker) Unlock() error {
	if l.ticker != nil {
		l.ticker.Stop()
	}
	if l.lockFile != nil {
		l.lockFile.Close()
	}
	err := os.Remove(l.lockFilePath)
	if err != nil {
		return err
	}
	return nil
}
