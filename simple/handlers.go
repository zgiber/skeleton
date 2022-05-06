package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type handler interface {
	Handle(method, url string, handler http.HandlerFunc)
	GetGreeting(http.ResponseWriter, *http.Request)
	GetGreetingWithNamedParam(http.ResponseWriter, *http.Request)
}

func initRoutes(h handler) {
	h.Handle(http.MethodGet, "/greeting", h.GetGreeting)
	h.Handle(http.MethodGet, "/greeting/:name", h.GetGreetingWithNamedParam)
}

func (t *transport) GetGreeting(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	greeting := t.service.Greeting(r.Context(), name)
	w.Write([]byte(greeting))
}

func (t *transport) GetGreetingWithNamedParam(w http.ResponseWriter, r *http.Request) {
	name := httprouter.ParamsFromContext(r.Context()).ByName("name")
	greeting := t.service.Greeting(r.Context(), name)
	w.Write([]byte(greeting))
}
