package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"goapp/internal/app/client"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lmsgprefix | log.Lshortfile)
}

func main() {
	connections := flag.Int("n", 1, "number of parallel connections")
	flag.Parse()

	if *connections < 1 {
		log.Fatal("-n must be greater than zero")
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := client.Start(ctx, *connections); err != nil {
		log.Fatal(err)
	}
}
