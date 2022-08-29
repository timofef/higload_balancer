package main

import (
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func main() {
	// Configuring prometheus metrics
	hits := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "hits",
	})
	if err := prometheus.Register(hits); err != nil {
		log.Fatal(err)
	}
	if err := prometheus.Register(prometheus.NewBuildInfoCollector()); err != nil {
		log.Fatal(err)
	}

	// Starting prometheus on different port
	go func() {
		metricsRouter := echo.New()
		metricsRouter.GET("/metrics", echo.WrapHandler(promhttp.Handler()))
		metricsRouter.Start(":5555")
	}()

	// Starting main app
	router := echo.New()
	router.GET("/app",
		func(ctx echo.Context) error {
			defer hits.Inc()
			// Sleep for random time interval <= 50ms
			duration := time.Duration(rand.Intn(50)) * time.Millisecond
			time.Sleep(duration)
			return ctx.String(http.StatusOK, "Hello world")
		},
	)
	router.Start(":8080")
}
