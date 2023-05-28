package jot

import (
	"sync"

	"github.com/rocketblend/rocketblend/pkg/jot/downloader"
	"github.com/rocketblend/rocketblend/pkg/jot/extractor"
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

	// Driver is what is used to interact with the jot database. It runs
	// transactions, and provides log output
	Driver struct {
		mutex      sync.Mutex
		mutexes    map[string]*sync.Mutex
		dir        string
		log        Logger
		downloader downloader.Downloader
		extractor  extractor.Extractor
	}

	// Options uses for specification of working golang-jot
	Options struct {
		Logger
		downloader.Downloader
		extractor.Extractor
	}
)
