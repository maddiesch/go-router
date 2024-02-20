package router

import (
	"context"
	"net/http"
)

func Run(ctx context.Context, server *http.Server) error {
	errChan := make(chan error, 1)

	go func() {
		errChan <- server.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		if err := server.Shutdown(context.TODO()); err != nil {
			return err
		}

		return context.Cause(ctx)
	case err := <-errChan:
		return err
	}
}
