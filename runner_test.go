package router_test

import (
	"context"
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
			<-time.After(200 * time.Millisecond)
			cancel()
		}()

		err := router.Run(ctx, s)

		assert.ErrorIs(t, err, context.Canceled)
	})

	t.Run("given a server that is shutdown", func(t *testing.T) {
		s := router.NewServer(":8990", nil)

		go func() {
			s.Shutdown(context.Background())
		}()

		err := router.Run(context.Background(), s)

		assert.ErrorIs(t, err, http.ErrServerClosed)
	})
}
