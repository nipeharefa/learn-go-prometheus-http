package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

func lama(ctx context.Context) error {
	return nil
}

func serang(count int) {
	iChan := make(chan int, count)

	// _ = make([]int, 100000)

	var wg sync.WaitGroup

	wg.Add(count)

	for i := 0; i < count; i++ {
		go func(ii chan int, i int, wg *sync.WaitGroup) {
			ii <- i*rand.Intn(4) + 1
			wg.Done()
		}(iChan, i, &wg)
	}
	wg.Wait()

	time.Sleep(1 * time.Second)

	close(iChan)
	for i := range iChan {
		fmt.Println(i)
	}
}

func main() {

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Method("GET", "/metrics", promhttp.Handler())
	r.Get("/naik", func(w http.ResponseWriter, r *http.Request) {
		serang(9000)
		w.Write([]byte("welcome"))
	})
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {

		ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
		defer cancel()

		errChan := make(chan error, 1)

		go func(e chan error, ctx context.Context) {
			e <- lama(ctx)
		}(errChan, ctx)

		select {
		case <-ctx.Done():
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("timeout"))
			return
		case e := <-errChan:
			if e != nil {
				w.WriteHeader(http.StatusBadRequest)
			}
		}
		w.Write([]byte("welcome"))
	})

	logger.Info("Running on :4123")
	http.ListenAndServe("0.0.0.0:4123", r)
}
