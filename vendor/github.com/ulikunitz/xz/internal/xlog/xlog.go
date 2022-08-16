// Copyright 2014-2021 Ulrich Kunitz. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package xlog provides a simple logging package that allows to disable
// certain message categories. It defines a type, Logger, with multiple
// methods for formatting output. The package has also a predefined
// 'standard' Logger accessible through helper function Print[f|ln],
// Fatal[f|ln], Panic[f|ln], Warn[f|ln], Print[f|ln] and Debug[f|ln]
// that are easier to use then creating a Logger manually. That logger
// writes to standard error and prints the date and time of each logged
// message, which can be configured using the function SetFlags.
//
// The Fatal functions call os.Exit(1) after the message is output
// unless not suppressed by the flags. The Panic functions call panic
// after the writing the log message unless suppressed.
package xlog

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
	"time"
)

// The flags define what information is prefixed to each log entry
// generated by the Logger. The Lno* versions allow the suppression of
// specific output. The bits are or'ed together to control what will be
// printed. There is no control over the order of the items printed and
// the format. The full format is:
//
//	2009-01-23 01:23:23.123123 /a/b/c/d.go:23: message
const (
	Ldate         = 1 << iota // the date: 2009-01-23
	Ltime                     // the time: 01:23:23
	Lmicroseconds             // microsecond resolution: 01:23:23.123123
	Llongfile                 // full file name and line number: /a/b/c/d.go:23
	Lshortfile                // final file name element and line number: d.go:23
	Lnopanic                  // suppresses output from Panic[f|ln] but not the panic call
	Lnofatal                  // suppresses output from Fatal[f|ln] but not the exit
	Lnowarn                   // suppresses output from Warn[f|ln]
	Lnoprint                  // suppresses output from Print[f|ln]
	Lnodebug                  // suppresses output from Debug[f|ln]
	// initial values for the standard logger
	Lstdflags = Ldate | Ltime | Lnodebug
)

// A Logger represents an active logging object that generates lines of
// output to an io.Writer. Each logging operation if not suppressed
// makes a single call to the Writer's Write method. A Logger can be
// used simultaneously from multiple goroutines; it guarantees to
// serialize access to the Writer.
type Logger struct {
	mu sync.Mutex // ensures atomic writes; and protects the following
	// fields
	prefix string    // prefix to write at beginning of each line
	flag   int       // properties
	out    io.Writer // destination for output
	buf    []byte    // for accumulating text to write
}

// New creates a new Logger. The out argument sets the destination to
// which the log output will be written. The prefix appears at the
// beginning of each log line. The flag argument defines the logging
// properties.
func New(out io.Writer, prefix string, flag int) *Logger {
	return &Logger{out: out, prefix: prefix, flag: flag}
}

// std is the standard logger used by the package scope functions.
var std = New(os.Stderr, "", Lstdflags)

// itoa converts the integer to ASCII. A negative widths will avoid
// zero-padding. The function supports only non-negative integers.
func itoa(buf *[]byte, i int, wid int) {
	var u = uint(i)
	if u == 0 && wid <= 1 {
		*buf = append(*buf, '0')
		return
	}
	var b [32]byte
	bp := len(b)
	for ; u > 0 || wid > 0; u /= 10 {
		bp--
		wid--
		b[bp] = byte(u%10) + '0'
	}
	*buf = append(*buf, b[bp:]...)
}

// formatHeader puts the header into the buf field of the buffer.
func (l *Logger) formatHeader(t time.Time, file string, line int) {
	l.buf = append(l.buf, l.prefix...)
	if l.flag&(Ldate|Ltime|Lmicroseconds) != 0 {
		if l.flag&Ldate != 0 {
			year, month, day := t.Date()
			itoa(&l.buf, year, 4)
			l.buf = append(l.buf, '-')
			itoa(&l.buf, int(month), 2)
			l.buf = append(l.buf, '-')
			itoa(&l.buf, day, 2)
			l.buf = append(l.buf, ' ')
		}
		if l.flag&(Ltime|Lmicroseconds) != 0 {
			hour, min, sec := t.Clock()
			itoa(&l.buf, hour, 2)
			l.buf = append(l.buf, ':')
			itoa(&l.buf, min, 2)
			l.buf = append(l.buf, ':')
			itoa(&l.buf, sec, 2)
			if l.flag&Lmicroseconds != 0 {
				l.buf = append(l.buf, '.')
				itoa(&l.buf, t.Nanosecond()/1e3, 6)
			}
			l.buf = append(l.buf, ' ')
		}
	}
	if l.flag&(Lshortfile|Llongfile) != 0 {
		if l.flag&Lshortfile != 0 {
			short := file
			for i := len(file) - 1; i > 0; i-- {
				if file[i] == '/' {
					short = file[i+1:]
					break
				}
			}
			file = short
		}
		l.buf = append(l.buf, file...)
		l.buf = append(l.buf, ':')
		itoa(&l.buf, line, -1)
		l.buf = append(l.buf, ": "...)
	}
}

func (l *Logger) output(calldepth int, now time.Time, s string) error {
	var file string
	var line int
	if l.flag&(Lshortfile|Llongfile) != 0 {
		l.mu.Unlock()
		var ok bool
		_, file, line, ok = runtime.Caller(calldepth)
		if !ok {
			file = "???"
			line = 0
		}
		l.mu.Lock()
	}
	l.buf = l.buf[:0]
	l.formatHeader(now, file, line)
	l.buf = append(l.buf, s...)
	if len(s) == 0 || s[len(s)-1] != '\n' {
		l.buf = append(l.buf, '\n')
	}
	_, err := l.out.Write(l.buf)
	return err
}

// Output writes the string s with the header controlled by the flags to
// the l.out writer. A newline will be appended if s doesn't end in a
// newline. Calldepth is used to recover the PC, although all current
// calls of Output use the call depth 2. Access to the function is serialized.
func (l *Logger) Output(calldepth, noflag int, v ...interface{}) error {
	now := time.Now()
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.flag&noflag != 0 {
		return nil
	}
	s := fmt.Sprint(v...)
	return l.output(calldepth+1, now, s)
}

// Outputf works like output but formats the output like Printf.
func (l *Logger) Outputf(calldepth int, noflag int, format string, v ...interface{}) error {
	now := time.Now()
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.flag&noflag != 0 {
		return nil
	}
	s := fmt.Sprintf(format, v...)
	return l.output(calldepth+1, now, s)
}

// Outputln works like output but formats the output like Println.
func (l *Logger) Outputln(calldepth int, noflag int, v ...interface{}) error {
	now := time.Now()
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.flag&noflag != 0 {
		return nil
	}
	s := fmt.Sprintln(v...)
	return l.output(calldepth+1, now, s)
}

// Panic prints the message like Print and calls panic. The printing
// might be suppressed by the flag Lnopanic.
func (l *Logger) Panic(v ...interface{}) {
	l.Output(2, Lnopanic, v...)
	s := fmt.Sprint(v...)
	panic(s)
}

// Panic prints the message like Print and calls panic. The printing
// might be suppressed by the flag Lnopanic.
func Panic(v ...interface{}) {
	std.Output(2, Lnopanic, v...)
	s := fmt.Sprint(v...)
	panic(s)
}

// Panicf prints the message like Printf and calls panic. The printing
// might be suppressed by the flag Lnopanic.
func (l *Logger) Panicf(format string, v ...interface{}) {
	l.Outputf(2, Lnopanic, format, v...)
	s := fmt.Sprintf(format, v...)
	panic(s)
}

// Panicf prints the message like Printf and calls panic. The printing
// might be suppressed by the flag Lnopanic.
func Panicf(format string, v ...interface{}) {
	std.Outputf(2, Lnopanic, format, v...)
	s := fmt.Sprintf(format, v...)
	panic(s)
}

// Panicln prints the message like Println and calls panic. The printing
// might be suppressed by the flag Lnopanic.
func (l *Logger) Panicln(v ...interface{}) {
	l.Outputln(2, Lnopanic, v...)
	s := fmt.Sprintln(v...)
	panic(s)
}

// Panicln prints the message like Println and calls panic. The printing
// might be suppressed by the flag Lnopanic.
func Panicln(v ...interface{}) {
	std.Outputln(2, Lnopanic, v...)
	s := fmt.Sprintln(v...)
	panic(s)
}

// Fatal prints the message like Print and calls os.Exit(1). The
// printing might be suppressed by the flag Lnofatal.
func (l *Logger) Fatal(v ...interface{}) {
	l.Output(2, Lnofatal, v...)
	os.Exit(1)
}

// Fatal prints the message like Print and calls os.Exit(1). The
// printing might be suppressed by the flag Lnofatal.
func Fatal(v ...interface{}) {
	std.Output(2, Lnofatal, v...)
	os.Exit(1)
}

// Fatalf prints the message like Printf and calls os.Exit(1). The
// printing might be suppressed by the flag Lnofatal.
func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.Outputf(2, Lnofatal, format, v...)
	os.Exit(1)
}

// Fatalf prints the message like Printf and calls os.Exit(1). The
// printing might be suppressed by the flag Lnofatal.
func Fatalf(format string, v ...interface{}) {
	std.Outputf(2, Lnofatal, format, v...)
	os.Exit(1)
}

// Fatalln prints the message like Println and calls os.Exit(1). The
// printing might be suppressed by the flag Lnofatal.
func (l *Logger) Fatalln(format string, v ...interface{}) {
	l.Outputln(2, Lnofatal, v...)
	os.Exit(1)
}

// Fatalln prints the message like Println and calls os.Exit(1). The
// printing might be suppressed by the flag Lnofatal.
func Fatalln(format string, v ...interface{}) {
	std.Outputln(2, Lnofatal, v...)
	os.Exit(1)
}

// Warn prints the message like Print. The printing might be suppressed
// by the flag Lnowarn.
func (l *Logger) Warn(v ...interface{}) {
	l.Output(2, Lnowarn, v...)
}

// Warn prints the message like Print. The printing might be suppressed
// by the flag Lnowarn.
func Warn(v ...interface{}) {
	std.Output(2, Lnowarn, v...)
}

// Warnf prints the message like Printf. The printing might be suppressed
// by the flag Lnowarn.
func (l *Logger) Warnf(format string, v ...interface{}) {
	l.Outputf(2, Lnowarn, format, v...)
}

// Warnf prints the message like Printf. The printing might be suppressed
// by the flag Lnowarn.
func Warnf(format string, v ...interface{}) {
	std.Outputf(2, Lnowarn, format, v...)
}

// Warnln prints the message like Println. The printing might be suppressed
// by the flag Lnowarn.
func (l *Logger) Warnln(v ...interface{}) {
	l.Outputln(2, Lnowarn, v...)
}

// Warnln prints the message like Println. The printing might be suppressed
// by the flag Lnowarn.
func Warnln(v ...interface{}) {
	std.Outputln(2, Lnowarn, v...)
}

// Print prints the message like fmt.Print. The printing might be suppressed
// by the flag Lnoprint.
func (l *Logger) Print(v ...interface{}) {
	l.Output(2, Lnoprint, v...)
}

// Print prints the message like fmt.Print. The printing might be suppressed
// by the flag Lnoprint.
func Print(v ...interface{}) {
	std.Output(2, Lnoprint, v...)
}

// Printf prints the message like fmt.Printf. The printing might be suppressed
// by the flag Lnoprint.
func (l *Logger) Printf(format string, v ...interface{}) {
	l.Outputf(2, Lnoprint, format, v...)
}

// Printf prints the message like fmt.Printf. The printing might be suppressed
// by the flag Lnoprint.
func Printf(format string, v ...interface{}) {
	std.Outputf(2, Lnoprint, format, v...)
}

// Println prints the message like fmt.Println. The printing might be
// suppressed by the flag Lnoprint.
func (l *Logger) Println(v ...interface{}) {
	l.Outputln(2, Lnoprint, v...)
}

// Println prints the message like fmt.Println. The printing might be
// suppressed by the flag Lnoprint.
func Println(v ...interface{}) {
	std.Outputln(2, Lnoprint, v...)
}

// Debug prints the message like Print. The printing might be suppressed
// by the flag Lnodebug.
func (l *Logger) Debug(v ...interface{}) {
	l.Output(2, Lnodebug, v...)
}

// Debug prints the message like Print. The printing might be suppressed
// by the flag Lnodebug.
func Debug(v ...interface{}) {
	std.Output(2, Lnodebug, v...)
}

// Debugf prints the message like Printf. The printing might be suppressed
// by the flag Lnodebug.
func (l *Logger) Debugf(format string, v ...interface{}) {
	l.Outputf(2, Lnodebug, format, v...)
}

// Debugf prints the message like Printf. The printing might be suppressed
// by the flag Lnodebug.
func Debugf(format string, v ...interface{}) {
	std.Outputf(2, Lnodebug, format, v...)
}

// Debugln prints the message like Println. The printing might be suppressed
// by the flag Lnodebug.
func (l *Logger) Debugln(v ...interface{}) {
	l.Outputln(2, Lnodebug, v...)
}

// Debugln prints the message like Println. The printing might be suppressed
// by the flag Lnodebug.
func Debugln(v ...interface{}) {
	std.Outputln(2, Lnodebug, v...)
}

// Flags returns the current flags used by the logger.
func (l *Logger) Flags() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.flag
}

// Flags returns the current flags used by the standard logger.
func Flags() int {
	return std.Flags()
}

// SetFlags sets the flags of the logger.
func (l *Logger) SetFlags(flag int) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.flag = flag
}

// SetFlags sets the flags for the standard logger.
func SetFlags(flag int) {
	std.SetFlags(flag)
}

// Prefix returns the prefix used by the logger.
func (l *Logger) Prefix() string {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.prefix
}

// Prefix returns the prefix used by the standard logger of the package.
func Prefix() string {
	return std.Prefix()
}

// SetPrefix sets the prefix for the logger.
func (l *Logger) SetPrefix(prefix string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.prefix = prefix
}

// SetPrefix sets the prefix of the standard logger of the package.
func SetPrefix(prefix string) {
	std.SetPrefix(prefix)
}

// SetOutput sets the output of the logger.
func (l *Logger) SetOutput(w io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.out = w
}

// SetOutput sets the output for the standard logger of the package.
func SetOutput(w io.Writer) {
	std.SetOutput(w)
}
