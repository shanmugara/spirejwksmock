package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"spirejwksmock/jwksmock"
	"syscall"
	"time"
)

func main() {
	// set default logger to stdout without timestamp
	log.SetOutput(os.Stdout)
	log.SetFlags(0)

	log.Println("Starting spirejwksmock")

	// create a context that is cancelled on SIGINT/SIGTERM
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// use a ticker instead of sleep so we can listen for shutdown
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Shutting down: context cancelled")
			return
		case <-ticker.C:
			// per-iteration timeout to avoid hanging indefinitely
			itCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
			bundles, myjwtsvid, err := jwksmock.FetchMyJWT(itCtx)
			cancel()
			if err != nil {
				log.Println("Error fetching JWT:", err)
				// don't exit; wait for next tick (could implement backoff here)
				continue
			}
			log.Println("Fetched JWT Bundle:", bundles)
			for _, b := range bundles.Bundles() {
				bb, err := b.Marshal()
				if err != nil {
					log.Println("Error marshaling bundle:", err)
					continue
				}
				log.Println("Bundle:", string(bb))
			}
			log.Println("Fetched My JWT SVID:", myjwtsvid.Marshal())
		}
	}
}
