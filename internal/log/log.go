package log

import (
	apexlog "github.com/apex/log"
	"github.com/apex/log/handlers/logfmt"
	"os"
)

// l is our logger singleton
var l apexlog.Logger

//Load logger singleton for main module
func Load(level apexlog.Level) {
	l.Level = level
	l.Handler = logfmt.New(os.Stderr)
}

//Get returns global logger
func Get() *apexlog.Logger {
	return &l
}
