package main

import (
	"net/http"
	"strings"
)

type Handler interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
}

type Router struct {
	handlers map[string]map[string]HandlerFunc
}

func (r *Router) HandleFunc(method, pattern string, h HandlerFunc) {
	m, ok := r.handlers[method]
	if !ok {
		m = make(map[string]HandlerFunc)
		r.handlers[method] = m
	}

	m[pattern] = h
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	for pattern, handler := range r.handlers[req.Method] {
		if ok, params := match(pattern, req.URL.Path); ok {

			// Context 생성
			ctx := Context{
				Params:         make(map[string]interface{}),
				ResponseWriter: w,
				Request:        req,
			}

			for k, v := range params {
				ctx.Params[k] = v
			}

			handler(&ctx)
			return
		}
	}

	http.NotFound(w, req)
}

func (r *Router) handler() HandlerFunc {
	return func(c *Context) {
		for pattern, handler := range r.handlers[c.Request.Method] {
			if ok, params := match(pattern, c.Request.URL.Path); ok {
				for k, v := range params {
					c.Params[k] = v
				}
				handler(c)
				return
			}
		}

		http.NotFound(c.ResponseWriter, c.Request)
		return
	}
}

func match(pattern, path string) (bool, map[string]string) {
	if pattern == path {
		return true, nil
	}

	patterns := strings.Split(pattern, "/")
	paths := strings.Split(path, "/")

	if len(patterns) != len(paths) {
		return false, nil
	}

	params := make(map[string]string)

	for i, p := range patterns {

		switch {
		case p == paths[i]:
			continue
		case len(p) > 0 && p[0] == ':':
			params[p[1:]] = paths[i]
		default:
			return false, nil
		}
	}

	return true, params
}
