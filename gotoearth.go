package gotoearth

import (
	"fmt"
	"strings"
)

// Event is a downtoearth event.
type Event struct {
	Body        map[string]interface{} `json:"body"`
	Path        map[string]string      `json:"path"`
	Querystring map[string]string      `json:"querystring"`
	Route       string                 `json:"route"`
	Stage       string                 `json:"stage"`
}

// Handler is used to call a method delegate based on an event.
type Handler interface {
	Handle(evt interface{}) (interface{}, error)
}

// Router is for delegating handling methods.
type Router struct {
	// Handlers are types that satisfy the Handler interface.
	// This is public so you can set it directly rather than using SetHandler.
	Handlers map[string]Handler
}

// RouteParts returns the HTTP verb and route path for the Event.
func RouteParts(r string) map[string]string {
	parts := strings.SplitN(r, ":", 2)
	return map[string]string{
		"verb": parts[0],
		"path": parts[1],
	}
}

// Parts returns the HTTP verb and route path for the Event.
func (e Event) Parts() map[string]string {
	return RouteParts(e.Route)
}

// Route routes the given event to the correct delegate method based upon route.
func (r Router) Route(route string, evt interface{}) (interface{}, error) {
	if route, ok := r.Handlers[route]; ok {
		return route.Handle(evt)
	}
	return "", fmt.Errorf("%s: no matching route", route)
}

// SimpleRoute routes the event to the correct delegate method.
// This will expect a gotoearth.Event and pass it in full.
func (r Router) SimpleRoute(evt Event) (interface{}, error) {
	if route, ok := r.Handlers[evt.Route]; ok {
		return route.Handle(evt)
	}
	return "", fmt.Errorf("%s: no matching route", evt.Route)
}

// SetHandler adds a Handler to the Router.
// This is probably superfluous. However, there may be need for fancy controls.
func (r *Router) SetHandler(route string, handler Handler) {
	if r.Handlers == nil {
		r.Handlers = map[string]Handler{}
	}
	r.Handlers[route] = handler
}
