package router

import (
	"context"
	"net/http"
	"time"
)

// Run runs the HTTP server using the provided context and server configuration.
// It returns an error if the server fails to start or encounters an error during shutdown.
func Run(ctx context.Context, server *http.Server) error {
	return RunServer(ctx, RunInput{
		Server: server,
		Run: func(_ context.Context, server *http.Server) error {
			return server.ListenAndServe()
		},
		Stop: func(ctx context.Context, server *http.Server) error {
			return server.Shutdown(ctx)
		},
	})
}

type RunInput struct {
	Server *http.Server                              // The HTTP Server to run
	Run    func(context.Context, *http.Server) error // The function called to run the server
	Stop   func(context.Context, *http.Server) error // The function called if the context is canceled to stop the server
}

// RunServer runs the server using the provided context and input parameters.
// It starts a goroutine to execute the Run method of the input object, which accepts the passed server.
// It waits for the server to complete or for the context to be canceled.
// If the context is canceled, it stops the server using the Stop method of the input object.
// Returns an error if the server fails to start or if the context is canceled.
func RunServer(ctx context.Context, in RunInput) error {
	errChan := make(chan error, 1)

	go func() {
		errChan <- in.Run(ctx, in.Server)
	}()

	select {
	case <-ctx.Done():
		if err := in.Stop(context.TODO(), in.Server); err != nil {
			return err
		}

		return context.Cause(ctx)
	case err := <-errChan:
		return err
	}
}

// NewServer creates a new HTTP server with the specified address and handler.
// The server has a read timeout of 60 seconds, a write timeout of 60 seconds,
// and an idle timeout of 5 minutes.
func NewServer(addr string, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  5 * time.Minute,
	}
}
