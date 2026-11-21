package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/mohamedhany/tcp-port-scanner-go/internal/output"
	"github.com/mohamedhany/tcp-port-scanner-go/internal/ports"
	"github.com/mohamedhany/tcp-port-scanner-go/internal/scanner"
)

const version = "0.1.0"

func main() {
	var host, portExpression string
	var timeout time.Duration
	var workers int
	var jsonOutput, showVersion bool

	flag.StringVar(&host, "host", "", "target hostname or IP address")
	flag.StringVar(&portExpression, "ports", "", "ports to scan, e.g. 22,80,443,8000-8010")
	flag.DurationVar(&timeout, "timeout", 500*time.Millisecond, "timeout per TCP connection")
	flag.IntVar(&workers, "workers", 1, "maximum concurrent workers")
	flag.BoolVar(&jsonOutput, "json", false, "write results as JSON")
	flag.BoolVar(&showVersion, "version", false, "print version")
	flag.Parse()

	if showVersion {
		fmt.Println(version)
		return
	}
	if strings.TrimSpace(host) == "" || strings.TrimSpace(portExpression) == "" {
		fmt.Fprintln(os.Stderr, "error: -host and -ports are required")
		flag.Usage()
		os.Exit(2)
	}
	if timeout <= 0 {
		fmt.Fprintln(os.Stderr, "error: -timeout must be greater than zero")
		os.Exit(2)
	}
	if workers <= 0 {
		fmt.Fprintln(os.Stderr, "error: -workers must be greater than zero")
		os.Exit(2)
	}

	portsToScan, err := ports.ParseExpression(portExpression)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(2)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	fmt.Fprintf(os.Stderr, "Scanning %s ports %s (timeout %s, workers %d)\n", host, portExpression, timeout, workers)
	results, err := scanner.ScanConcurrent(ctx, host, portsToScan, timeout, workers)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(3)
	}
	if ctx.Err() != nil {
		fmt.Fprintln(os.Stderr, "scan cancelled")
		os.Exit(4)
	}

	if jsonOutput {
		err = output.WriteJSON(os.Stdout, results)
	} else {
		err = output.WriteHuman(os.Stdout, results)
	}
	if err != nil && !errors.Is(err, os.ErrClosed) {
		fmt.Fprintf(os.Stderr, "error writing output: %s\n", err)
		os.Exit(3)
	}
}
