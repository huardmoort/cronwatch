package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/example/cronwatch/internal/config"
	"github.com/example/cronwatch/internal/heartbeat"
	"github.com/example/cronwatch/internal/notify"
	"github.com/example/cronwatch/internal/reporter"
	"github.com/example/cronwatch/internal/store"
	"github.com/example/cronwatch/internal/watcher"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	configPath := flag.String("config", "cronwatch.yaml", "path to config file")
	reportOnly := flag.Bool("report", false, "print status report and exit")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	st, err := store.New(cfg.StorePath)
	if err != nil {
		return fmt.Errorf("opening store: %w", err)
	}

	if *reportOnly {
		r := reporter.New(cfg, st)
		if err := r.Print(os.Stdout); err != nil {
			return fmt.Errorf("printing report: %w", err)
		}
		return nil
	}

	notifier := notify.New(cfg.Webhook)
	w := watcher.New(cfg, st, notifier)

	mux := heartbeat.New(cfg, st)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	log.Printf("cronwatch starting — monitoring %d job(s)", len(cfg.Jobs))

	go func() {
		if err := mux.ListenAndServe(ctx); err != nil {
			log.Printf("heartbeat server stopped: %v", err)
		}
	}()

	watcher.RunLoop(ctx, w, cfg.Interval)
	log.Println("cronwatch stopped")
	return nil
}
