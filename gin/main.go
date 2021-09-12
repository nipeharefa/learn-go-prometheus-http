package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	// "github.com/prometheus/client_golang/prometheus/promhttp"
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
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	// r.GET(gin)

	h := gin.WrapH(promhttp.Handler())

	r.GET("/metrics", h)

	r.GET("/naik", func(c *gin.Context) {
		_, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)

		defer cancel()

		serang(9000)
		c.String(http.StatusOK, "welcome")
	})

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.Run("0.0.0.0:4123") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
