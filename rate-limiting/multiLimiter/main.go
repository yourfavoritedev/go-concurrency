package main

import (
	"context"
	"log"
	"os"
	"sort"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type rateLimiter interface {
	Wait(context.Context) error
	Limit() rate.Limit
}

type multiLimiter struct {
	limiters []rateLimiter
}

type apiConnection struct {
	rateLimiter rateLimiter
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

// multiLimiter needs to integrate Wait and Limit in order to be used in rateLimiter interface
func (l *multiLimiter) Wait(ctx context.Context) error {
	for _, l := range l.limiters {
		if err := l.Wait(ctx); err != nil {
			return err
		}
	}
	return nil
}

// will return the most restrictive limit in the slice
func (l *multiLimiter) Limit() rate.Limit {
	return l.limiters[0].Limit()
}

func main() {
	defer log.Printf("Done.")

	generateMultiLimiter := func(limiters []rateLimiter) *multiLimiter {
		byLimit := func(i, j int) bool {
			return limiters[i].Limit() < limiters[j].Limit()
		}

		sort.Slice(limiters, byLimit)
		return &multiLimiter{limiters: limiters}
	}

	per := func(eventCount int, duration time.Duration) rate.Limit {
		return rate.Every(duration / time.Duration(eventCount))
	}

	open := func() *apiConnection {
		secondLimit := rate.NewLimiter(per(2, time.Second), 1)
		minuteLimit := rate.NewLimiter(per(10, time.Minute), 10)

		limiters := []rateLimiter{secondLimit, minuteLimit}

		return &apiConnection{
			rateLimiter: generateMultiLimiter(limiters),
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
