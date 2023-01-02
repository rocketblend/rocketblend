package jot

import (
	"os"
	"sync"
)

// stat checks for dir, if path isn't a directory check to see if it's a file
func stat(path string) (fi os.FileInfo, err error) {
	if fi, err = os.Stat(path); os.IsNotExist(err) {
		fi, err = os.Stat(path)
	}
	return
}

// getOrCreateMutex creates a new reference specific mutex any time a reference
// is being modfied to avoid unsafe operations
func (d *Driver) getOrCreateMutex(reference string) *sync.Mutex {
	// create mutex
	d.mutex.Lock()
	defer d.mutex.Unlock()

	// if the mutex doesn't exist make it
	m, ok := d.mutexes[reference]
	if !ok {
		m = &sync.Mutex{}
		d.mutexes[reference] = m
	}

	return m
}
