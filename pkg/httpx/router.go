package httpx

import (
	"net/http"
)

type (
	httpPath  map[string]http.HandlerFunc
	httpRoute map[string]httpPath
)

type Router struct {
	rules httpRoute
}

func NewRouter() *Router {
	return &Router{
		rules: make(httpRoute),
	}
}

func (r *Router) Handle(method string, path string, handler http.HandlerFunc) {
	_, exist := r.rules[path]
	if !exist {
		r.rules[path] = make(httpPath)
	}
	r.rules[path][method] = handler
}

func (r *Router) findHandler(path string, method string) (http.HandlerFunc, bool, bool) {
	_, exist := r.rules[path]
	handler, methodExist := r.rules[path][method]
	return handler, methodExist, exist
}

func (r *Router) ServeHTTP(w http.ResponseWriter, request *http.Request) {
	handler, methodExist, exist := r.findHandler(request.URL.Path, request.Method)

	if !exist {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if !methodExist {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	handler(w, request)
}
