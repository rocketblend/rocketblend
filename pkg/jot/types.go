package jot

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

	Downloader interface {
		Download(path string, downloadUrl string) error
	}

	Extractor interface {
		Extract(path string, extractPath string) error
	}

	// Driver is what is used to interact with the jot database. It runs
	// transactions, and provides log output
	Driver struct {
		mutex      sync.Mutex
		mutexes    map[string]*sync.Mutex
		dir        string     // the directory where jot will create the database
		log        Logger     // the logger jot will use for logging
		downloader Downloader // the downloader jot will use for downloading
		extractor  Extractor  // the extractor jot will use for extracting
	}

	// Options uses for specification of working golang-jot
	Options struct {
		Logger     // the logger jot will use (configurable)
		Downloader // the downloader jot will use (configurable)
		Extractor  // the extractor jot will use (configurable)
	}
)
