package slog

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Entry is Stackdriver Logging Entry
type Entry struct {
	InsertID         string            `json:"insertId"`
	Severity         string            `json:"severity"`
	Labels           map[string]string `json:"labels"`
	LogName          string            `json:"logName"`
	ReceiveTimestamp time.Time         `json:"receiveTimestamp"`
	Resource         MonitoredResource `json:"resource"`
	JSONPayload      interface{}       `json:"jsonPayload"`
	Timestamp        time.Time         `json:"timestamp"`
}

// MonitoredResource is Log Resource
// https://cloud.google.com/logging/docs/reference/v2/rest/v2/MonitoredResource
type MonitoredResource struct {
	Type   string            `json:"type"`
	Labels map[string]string `json:"labels"`
}

// Log is Log Object
type Log struct {
	Entry    Entry `json:"entry"`
	Messages []string
}

type contextLogKey struct{}

// WithLog is context.ValueにLogを入れたものを返す
// Log周期開始時に利用する
func WithLog(ctx context.Context) context.Context {
	labels := map[string]string{"hoge": "fuga"}
	l := &Log{
		Entry: Entry{
			InsertID:         time.Now().String(),
			Labels:           labels,
			LogName:          "projects/metal-tile-dev1/logs/slog",
			ReceiveTimestamp: time.Now(),
			Resource: MonitoredResource{
				Type:   "slog",
				Labels: labels,
			},
			Severity:  "WARNING",
			Timestamp: time.Now(),
		},
	}

	return context.WithValue(ctx, contextLogKey{}, l)
}

// Info is output info level Log
func Info(ctx context.Context, message string) {
	l, ok := ctx.Value(contextLogKey{}).(*Log)
	if !ok {
		panic(fmt.Sprintf("not contain log. message = %s", message))
	}
	l.Messages = append(l.Messages, message)
}

// Flush is ログを出力する
func Flush(ctx context.Context) {
	l, ok := ctx.Value(contextLogKey{}).(*Log)
	if ok {
		encoder := json.NewEncoder(os.Stdout)
		l.Entry.JSONPayload = l.Messages
		if err := encoder.Encode(l.Entry); err != nil {
			_, err := os.Stdout.WriteString(err.Error())
			if err != nil {
				panic(err)
			}
		}
	}
}
