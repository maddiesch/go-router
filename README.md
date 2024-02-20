# Router

A simple wrapper for [httprouter](https://github.com/julienschmidt/httprouter) that supports middleware & sub-routers. Additionally, it uses the new [Go 1.22 http.Request feature to set Path values on the http.Request](https://pkg.go.dev/net/http#Request.PathValue).
