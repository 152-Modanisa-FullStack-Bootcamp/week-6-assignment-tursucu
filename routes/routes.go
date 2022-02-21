package routes

import (
	"net/http"
	"regexp"
)

type Route struct {
	Method  string
	Regex   *regexp.Regexp
	Handler http.HandlerFunc
}
type CtxKey struct{}

func NewRoute(method, pattern string, handler http.HandlerFunc) Route {
	return Route{method, regexp.MustCompile("^" + pattern + "$"), handler}
}

func GetField(r *http.Request, index int) string {
	fields := r.Context().Value(CtxKey{}).([]string)
	return fields[index]
}
