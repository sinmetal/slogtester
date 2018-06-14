package main

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/profiler"
	"github.com/sinmetal/slog"
)

func main() {
	if err := profiler.Start(profiler.Config{
		Service:        "slogtester",
		ServiceVersion: "1.0.0",
	}); err != nil {
		panic(err)
	}

	errCh := make(chan error)

	go func() {
		for {
			ctx := context.Background()
			workWithCancel(ctx, "1")
			workWithCancel(ctx, "2")
			time.Sleep(10 * time.Second)
		}
	}()

	go func() {
		for {
			ctx := context.Background()
			workWithTimeout(ctx, "1")
			workWithTimeout(ctx, "2")
			time.Sleep(10 * time.Second)
		}
	}()

	go func() {
		for {
			ctx := context.Background()
			workWithDeadline(ctx, "1")
			workWithDeadline(ctx, "2")
			time.Sleep(10 * time.Second)
		}
	}()

	<-errCh
}

func workWithCancel(ctx context.Context, v string) {
	ctx = slog.WithLog(ctx)
	defer slog.Flush(ctx)
	slog.SetLogName(ctx, "WithCancel")
	cctx, cancel := context.WithCancel(ctx)
	defer cancel()

	s := struct {
		Key   string
		Value string
	}{
		Key:   "message",
		Value: v,
	}
	logger(cctx, fmt.Sprintf("Hello SLOG WithCancel. %v", v))
	logger(cctx, s)
}

func workWithTimeout(ctx context.Context, v string) {
	ctx = slog.WithLog(ctx)
	defer slog.Flush(ctx)
	slog.SetLogName(ctx, "WithTimeout")
	cctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	s := struct {
		Key   string
		Value string
	}{
		Key:   "message",
		Value: v,
	}
	logger(cctx, fmt.Sprintf("Hello SLOG WithTimeout. %v", v))
	logger(cctx, s)
	time.Sleep(3 * time.Second)
}

func workWithDeadline(ctx context.Context, v string) {
	ctx = slog.WithLog(ctx)
	defer slog.Flush(ctx)
	slog.SetLogName(ctx, "WithDeadline")
	cctx, cancel := context.WithDeadline(ctx, time.Now().Add(2*time.Second))
	defer cancel()

	s := struct {
		Key   string
		Value string
	}{
		Key:   "message",
		Value: v,
	}
	logger(cctx, fmt.Sprintf("Hello SLOG WithDadline. %v", v))
	logger(cctx, s)
	time.Sleep(3 * time.Second)
}

func logger(ctx context.Context, body interface{}) {
	slog.Info(ctx, body)
}
