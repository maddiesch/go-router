package router_test

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
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

		r.Use(middleware.RequestID(), middleware.Logger(slog.LevelInfo))

		r.HandleFunc(http.MethodGet, "/", func(w http.ResponseWriter, r *http.Request) {
			spew.Dump(r.Context())
			w.WriteHeader(http.StatusOK)
		})

		r.Handle(http.MethodGet, "/health-check", router.HealthCheckHandler())

		r.HandleFunc(http.MethodGet, "/flush", func(w http.ResponseWriter, r *http.Request) {
			rc := http.NewResponseController(w)
			if err := rc.SetWriteDeadline(time.Now().Add(1 * time.Second)); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if err := rc.Flush(); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		})

		r.Sub("/api", func(sub *router.Router) {
			sub.Use(middleware.NoCache())
			sub.HandleFunc(http.MethodGet, "/posts/:id", func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("id:" + r.PathValue("id")))
			})
		})

		customRequestID := middleware.RequestID(func() string {
			return "custom-id"
		})

		r.Handle(http.MethodGet, "/custom-request-id", customRequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})))

		r.HandleFunc(http.MethodGet, "/panic", func(w http.ResponseWriter, r *http.Request) {
			panic("testing panic recovery")
		})

		s := httptest.NewServer(r)
		t.Cleanup(s.Close)

		t.Run("requesting a health check", func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, s.URL+"/health-check", nil)

			resp, err := http.DefaultClient.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusOK, resp.StatusCode)

			content, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			assert.Equal(t, "OK", string(content))
		})

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

		t.Run("requesting a panic route", func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, s.URL+"/panic", nil)

			resp, err := http.DefaultClient.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		})

		t.Run("getting a custom request id", func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, s.URL+"/custom-request-id", nil)

			resp, err := http.DefaultClient.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusOK, resp.StatusCode)
			assert.Equal(t, "custom-id", resp.Header.Get("X-Request-ID"))
		})

		t.Run("ensure handler response writer still supports Flush", func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, s.URL+"/flush", nil)

			resp, err := http.DefaultClient.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusOK, resp.StatusCode)
		})
	})
}
