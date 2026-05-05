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
	configPath := flag.String("config", "cronwatch.yaml", "path to config file")
	reportOnly := flag.Bool("report", false, "print status report and exit")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading config: %v\n", err)
		os.Exit(1)
	}

	st, err := store.New(cfg.StorePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error opening store: %v\n", err)
		os.Exit(1)
	}

	if *reportOnly {
		r := reporter.New(cfg, st)
		if err := r.Print(os.Stdout); err != nil {
			fmt.Fprintf(os.Stderr, "error printing report: %v\n", err)
			os.Exit(1)
		}
		return
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
}
