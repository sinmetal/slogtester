package slog

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

// Entry is Stackdriver Logging Entry
type Entry struct {
	Timestamp   time.Time `json:"timestamp"`
	Severity    string    `json:"severity"`
	severity    Severity
	JSONPayload interface{} `json:"jsonPayload"`
}

// Timestamp is Stackdriver Logging Timestamp
type Timestamp struct {
	Seconds int64 `json:"seconds"`
	Nanos   int   `json:"nanos"`
}

// Log is Log Object
type Log struct {
	Entry    Entry    `json:"entry"`
	Messages []string `json:"messages"`
}

// Info is Add Log Message for Info Level
func (l *Log) Info(message string) {
	m := strings.Replace(message, "\n", "", -1)
	l.Messages = append(l.Messages, m)
	l.setSeverity(INFO)
}

// Infof is Add Log Message for Info Level
func (l *Log) Infof(format string, v ...interface{}) {
	l.Info(fmt.Sprintf(format, v...))
}

// Error is Add Log Message for Error Level
func (l *Log) Error(message string) {
	m := strings.Replace(message, "\n", "", -1)
	l.Messages = append(l.Messages, m)
	l.setSeverity(ERROR)
}

// Errorf is Add Log Message for Error Level
func (l *Log) Errorf(format string, v ...interface{}) {
	l.Error(fmt.Sprintf(format, v...))
}

// Flush is Flush to Log
func (l *Log) Flush() {
	j, err := l.flush()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("%s\n", string(j))
}

func (l *Log) setSeverity(severity Severity) {
	if l.Entry.Severity == "" {
		l.Entry.severity = severity
		l.Entry.Severity = severity.String()
		return
	}

	if l.Entry.severity < severity {
		l.Entry.severity = severity
		l.Entry.Severity = severity.String()
		return
	}
}

func (l *Log) flush() ([]byte, error) {
	b, err := json.Marshal(l.Messages)
	if err == nil {
		// l.Entry.Message = string(b)
	} else {
		return nil, err
	}

	b, err = json.Marshal(l.Entry)
	if err == nil {
		l.Messages = nil
	}
	return b, err
}

var logMap map[context.Context]*Log

func init() {
	logMap = make(map[context.Context]*Log)
}

// Info is output info level Log
func Info(ctx context.Context, message string) {
	e, ok := logMap[ctx]
	if !ok {
		e = &Log{
			Entry: Entry{
				Severity:  "INFO",
				Timestamp: time.Now(),
			},
			Messages: []string{},
		}
		go log(ctx)
	}
	e.Messages = append(e.Messages, message)
	logMap[ctx] = e
}

func log(ctx context.Context) {
	fmt.Println("log start")
	select {
	case <-ctx.Done():
		fmt.Println("log ctx.Done()")
		e, ok := logMap[ctx]
		if ok {
			fmt.Println("log ctx.Done() logMap ok;")
			encoder := json.NewEncoder(os.Stdout)
			e.Entry.JSONPayload = e.Messages
			if err := encoder.Encode(e.Entry); err != nil {
				_, err := os.Stderr.WriteString(err.Error())
				if err != nil {
					panic(err)
				}
			}
			fmt.Println("log ctx.Done() logMap ok; done.")
		}
	}
	fmt.Println("log end")
}
