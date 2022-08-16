package scribble

import (
	"sync"
)

type (
	// Logger is a generic logger interface
	Logger interface {
		Debug(string ...interface{})
		Info(string ...interface{})
		Print(string ...interface{})
		Warn(string ...interface{})
		Warning(string ...interface{})
		Error(string ...interface{})
		Fatal(string ...interface{})
		Panic(string ...interface{})
	}

	// Driver is what is used to interact with the scribble database. It runs
	// transactions, and provides log output
	Driver struct {
		mutex   sync.Mutex
		mutexes map[string]*sync.Mutex
		dir     string // the directory where scribble will create the database
		log     Logger // the logger scribble will log to
	}

	// Options uses for specification of working golang-scribble
	Options struct {
		Logger // the logger scribble will use (configurable)
	}
)
