package router_test

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/maddiesch/go-router"
	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	t.Run("given a cancel context", func(t *testing.T) {
		s := router.NewServer(":8990", nil)

		ctx, cancel := context.WithCancel(context.Background())

		go func() {
			<-time.After(50 * time.Millisecond)
			cancel()
		}()

		err := router.Run(ctx, s)

		assert.ErrorIs(t, err, context.Canceled)
	})

	t.Run("given a server that is shutdown", func(t *testing.T) {
		s := router.NewServer(":8990", nil)

		go func() {
			<-time.After(50 * time.Millisecond)
			s.Shutdown(context.Background())
		}()

		err := router.Run(context.Background(), s)

		assert.ErrorIs(t, err, http.ErrServerClosed)
	})
}

func TestRunServer(t *testing.T) {
	s := router.NewServer(":8990", nil)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		<-time.After(50 * time.Millisecond)
		cancel()
	}()

	stopErr := errors.New("stop return error")

	err := router.RunServer(ctx, router.RunInput{
		Server: s,
		Run: func(_ context.Context, s *http.Server) error {
			return s.ListenAndServe()
		},
		Stop: func(ctx context.Context, s *http.Server) error {
			s.Shutdown(ctx)
			return stopErr
		},
	})

	assert.ErrorIs(t, err, stopErr)
}
