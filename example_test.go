package router_test

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/maddiesch/go-router"
	"github.com/maddiesch/go-router/middleware"
)

func ExampleRouter() {
	r := router.New()

	// Register the HealthCheckHandler to the "/health-check" route, before the middleware.
	// This will remove the middleware overhead from the HealthCheck endpoint
	r.Handle(http.MethodGet, "/health-check", router.HealthCheckHandler())

	r.Use(middleware.RequestID(), middleware.Logger(slog.LevelInfo))

	// All routes will be prefixed with "/api"
	r.Sub("/api", func(sub *router.Router) {
		sub.Use(middleware.NoCache())

		// Full path /api/posts/:id
		sub.HandleFunc("GET", "/posts/:id", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("id:" + r.PathValue("id")))
		})
	})

	s := router.NewServer(":3000", r)

	go func() {
		<-time.After(100 * time.Millisecond)
		s.Shutdown(context.Background())
	}()

	err := router.Run(context.Background(), s)

	fmt.Println(err)
	// Output: http: Server closed
}
