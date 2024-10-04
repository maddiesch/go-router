package router

import (
	"log/slog"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func New(config ...RouterConfig) *Router {
	r := &Router{
		router: httprouter.New(),
		path:   "/",
	}

	r.router.RedirectTrailingSlash = true
	r.router.HandleMethodNotAllowed = false
	r.router.RedirectFixedPath = true
	r.router.PanicHandler = func(w http.ResponseWriter, r *http.Request, i any) {
		slog.ErrorContext(r.Context(), "Panic Recovered", slog.Any("panic", i))

		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	r.Configure(config...)

	return r
}

type RouterConfig func(*httprouter.Router)

type Router struct {
	router      *httprouter.Router
	path        string
	middlewares []func(http.Handler) http.Handler
}

func (r *Router) Configure(config ...RouterConfig) {
	for _, config := range config {
		config(r.router)
	}
}

func (r *Router) SubRouter(path string) *Router {
	return &Router{
		router:      r.router,
		path:        r.fullPath(path),
		middlewares: r.middlewares,
	}
}

func (r *Router) Sub(path string, fn func(*Router)) {
	sub := r.SubRouter(path)
	fn(sub)
}

func (r *Router) fullPath(path string) string {
	return httprouter.CleanPath(r.path + "/" + path)
}

func (r *Router) Use(m ...func(http.Handler) http.Handler) {
	r.middlewares = append(r.middlewares, m...)
}

func (r *Router) HandleFunc(method, path string, h func(w http.ResponseWriter, r *http.Request)) {
	r.Handle(method, path, http.HandlerFunc(h))
}

func (r *Router) Handle(method, path string, h http.Handler) {
	for i := len(r.middlewares) - 1; i >= 0; i-- {
		h = r.middlewares[i](h)
	}

	r.router.Handle(method, r.fullPath(path), func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		for _, p := range p {
			r.SetPathValue(p.Key, p.Value)
		}

		h.ServeHTTP(w, r)
	})
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.router.ServeHTTP(w, req)
}

var (
	_ http.Handler = (*Router)(nil)
)
