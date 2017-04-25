// Copyright 2016 Qiang Xue. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package routing

import (
	"github.com/acsellers/inflections"
	"reflect"
	"strings"
)

// RouteGroup represents a group of routes that share the same path prefix.
type RouteGroup struct {
	prefix   string
	router   *Router
	handlers []Handler
}

// newRouteGroup creates a new RouteGroup with the given path prefix, router, and handlers.
func newRouteGroup(prefix string, router *Router, handlers []Handler) *RouteGroup {
	return &RouteGroup{
		prefix:   prefix,
		router:   router,
		handlers: handlers,
	}
}

// Get adds a GET route to the router with the given route path and handlers.
func (r *RouteGroup) Get(path string, handlers ...Handler) *Route {
	return newRoute(path, r).Get(handlers...)
}

// Post adds a POST route to the router with the given route path and handlers.
func (r *RouteGroup) Post(path string, handlers ...Handler) *Route {
	return newRoute(path, r).Post(handlers...)
}

// Put adds a PUT route to the router with the given route path and handlers.
func (r *RouteGroup) Put(path string, handlers ...Handler) *Route {
	return newRoute(path, r).Put(handlers...)
}

// Patch adds a PATCH route to the router with the given route path and handlers.
func (r *RouteGroup) Patch(path string, handlers ...Handler) *Route {
	return newRoute(path, r).Patch(handlers...)
}

// Delete adds a DELETE route to the router with the given route path and handlers.
func (r *RouteGroup) Delete(path string, handlers ...Handler) *Route {
	return newRoute(path, r).Delete(handlers...)
}

// Connect adds a CONNECT route to the router with the given route path and handlers.
func (r *RouteGroup) Connect(path string, handlers ...Handler) *Route {
	return newRoute(path, r).Connect(handlers...)
}

// Head adds a HEAD route to the router with the given route path and handlers.
func (r *RouteGroup) Head(path string, handlers ...Handler) *Route {
	return newRoute(path, r).Head(handlers...)
}

// Options adds an OPTIONS route to the router with the given route path and handlers.
func (r *RouteGroup) Options(path string, handlers ...Handler) *Route {
	return newRoute(path, r).Options(handlers...)
}

// Trace adds a TRACE route to the router with the given route path and handlers.
func (r *RouteGroup) Trace(path string, handlers ...Handler) *Route {
	return newRoute(path, r).Trace(handlers...)
}

// Any adds a route with the given route, handlers, and the HTTP methods as listed in routing.Methods.
func (r *RouteGroup) Any(path string, handlers ...Handler) *Route {
	route := newRoute(path, r)
	for _, method := range Methods {
		route.add(method, handlers)
	}
	return route
}

// RegisterController Registers a controller (i.e. struct) with functions
// of the type: 'func(*routing.Context) error' and skips all other functions
func (r *RouteGroup) RegisterController(controller interface{}) {
	t := reflect.TypeOf(controller)
	for i := 0; i < t.NumMethod(); i++ {
		method := t.Method(i)
		// Check for func(*routing.Context) error
		if method.Type.NumOut() == 1 && method.Type.Out(0).Name() == "error" &&
			method.Type.NumIn() == 2 && method.Type.In(1).String() == "*routing.Context" {
			uri := "/" + inflections.Underscore(t.Name())
			uri += "/" + inflections.Underscore(method.Name)
			mt, found := t.MethodByName(method.Name + "Params")
			if found {
				r := mt.Func.Call([]reflect.Value{reflect.ValueOf(controller)})
				if len(r) == 1 {
					uri += "/" + r[1].Interface().(string)
				}
			}
			r.Any(uri, func(ctx *Context) error {
				var params []reflect.Value
				params = append(params, reflect.ValueOf(controller))
				params = append(params, reflect.ValueOf(ctx))
				errVal := method.Func.Call(params)
				i := errVal[0].Interface()
				if i != nil {
					return i.(error)
				}
				return nil
			})
		}
	}
}

// To adds a route to the router with the given HTTP methods, route path, and handlers.
// Multiple HTTP methods should be separated by commas (without any surrounding spaces).
func (r *RouteGroup) To(methods, path string, handlers ...Handler) *Route {
	route := newRoute(path, r)
	for _, method := range strings.Split(methods, ",") {
		route.add(method, handlers)
	}
	return route
}

// Group creates a RouteGroup with the given route path prefix and handlers.
// The new group will combine the existing path prefix with the new one.
// If no handler is provided, the new group will inherit the handlers registered
// with the current group.
func (r *RouteGroup) Group(prefix string, handlers ...Handler) *RouteGroup {
	if len(handlers) == 0 {
		handlers = make([]Handler, len(r.handlers))
		copy(handlers, r.handlers)
	}
	return newRouteGroup(r.prefix+prefix, r.router, handlers)
}

// Use registers one or multiple handlers to the current route group.
// These handlers will be shared by all routes belong to this group and its subgroups.
func (r *RouteGroup) Use(handlers ...Handler) {
	r.handlers = append(r.handlers, handlers...)
}
