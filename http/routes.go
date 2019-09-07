package http

import (
	"net/http"
)

//Route defines an http route
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

//Routes is a collection of http routes supported by this web application
type Routes []Route
