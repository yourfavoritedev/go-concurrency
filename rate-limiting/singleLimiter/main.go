package main

import (
	"context"
	"log"
	"os"
	"sync"

	"golang.org/x/time/rate"
)

type apiConnection struct {
	rateLimiter *rate.Limiter
}

func (a *apiConnection) readFile(ctx context.Context) error {
	if err := a.rateLimiter.Wait(ctx); err != nil {
		return err
	}
	return nil
}

func (a *apiConnection) resolveAddress(ctx context.Context) error {
	if err := a.rateLimiter.Wait(ctx); err != nil {
		return err
	}
	return nil
}

func main() {
	defer log.Printf("Done.")

	open := func() *apiConnection {
		return &apiConnection{
			// 1 event per second with maximum depth of 1 token per bucket (only 1 request can be processed per second)
			rateLimiter: rate.NewLimiter(rate.Limit(1), 1),
		}
	}

	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime | log.LUTC)

	apiConnection := open()
	var wg sync.WaitGroup
	wg.Add(20)

	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			// wait on the rate limiter to have enough access tokens for us to complete our request
			err := apiConnection.readFile(context.Background())
			if err != nil {
				log.Printf("Cannot ReadFile: %v", err)
			}
			log.Printf("ReadFile")
		}()
	}

	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			err := apiConnection.resolveAddress(context.Background())
			if err != nil {
				log.Printf("Cannot ResolveAddress: %v", err)
			}
			log.Printf("ResolveAddress")
		}()
	}

	wg.Wait()

}
