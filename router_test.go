package router_test

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/maddiesch/go-router"
	"github.com/maddiesch/go-router/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRouter(t *testing.T) {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})))

	t.Run("given a basic router with sub routers", func(t *testing.T) {
		r := router.New()

		r.Use(middleware.Logger(slog.LevelInfo))

		r.HandleFunc(http.MethodGet, "/", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		sub := r.Sub("/api")
		sub.HandleFunc(http.MethodGet, "/posts/:id", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("id:" + r.PathValue("id")))
		})

		s := httptest.NewServer(r)
		t.Cleanup(s.Close)

		t.Run("requesting a root path", func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, s.URL+"/", nil)

			resp, err := http.DefaultClient.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusOK, resp.StatusCode)
		})

		t.Run("requesting sub-router path", func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, s.URL+"/api/posts/foobar", nil)

			resp, err := http.DefaultClient.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusOK, resp.StatusCode)

			content, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			assert.Equal(t, "id:foobar", string(content))
		})
	})
}
