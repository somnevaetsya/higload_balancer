package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var hitsCounter = promauto.NewCounter(prometheus.CounterOpts{
	Name: "hits_counter",
	Help: "Number of hits to the server",
})

func IncrementHitsMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		hitsCounter.Inc()
		return next(c)
	}
}

func LoggingMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()
		res := next(c)
		fmt.Printf("%s %s %s %d %s\n",
			c.RealIP(),
			c.Request().Method,
			c.Request().RequestURI,
			c.Response().Status,
			time.Since(start))
		return res
	}
}

func WorkHandler(ctx echo.Context) error {
	const maxWaitingTime = 500
	time.Sleep(time.Duration(rand.Intn(maxWaitingTime)) * time.Millisecond)
	return ctx.String(http.StatusOK, "Done!\n")
}

func main() {
	go func() {
		prometheusRouter := echo.New()
		prometheusRouter.GET("/prometheus", echo.WrapHandler(promhttp.Handler()))
		log.Fatal(prometheusRouter.Start(":9090"))
	}()

	router := echo.New()
	router.Use(LoggingMiddleware)
	router.Use(IncrementHitsMiddleware)

	router.GET("/*", WorkHandler)
	log.Fatal(router.Start(":80"))
}
