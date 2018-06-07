package main

import (
	"context"
	"time"
	"fmt"

	"github.com/sinmetal/slog"
)

func main() {
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

	<- errCh
}

func workWithCancel(ctx context.Context, v string) {
	cctx, cancel := context.WithCancel(ctx)
	defer cancel()
	logger(cctx, fmt.Sprintf("Hello SLOG WithCancel. %v", v))
}

func workWithTimeout(ctx context.Context, v string) {
	cctx, cancel := context.WithTimeout(ctx, 2 * time.Second)
	defer cancel()
	logger(cctx, fmt.Sprintf("Hello SLOG WithTimeout. %v", v))
	time.Sleep(3 * time.Second)
}

func workWithDeadline(ctx context.Context, v string) {
	cctx, cancel := context.WithDeadline(ctx, time.Now().Add(2 * time.Second))
	defer cancel()
	logger(cctx, fmt.Sprintf("Hello SLOG WithDadline. %v", v))
	time.Sleep(3 * time.Second)
}

func logger(ctx context.Context, msg string) {
	slog.Info(ctx, msg)
}