package controllers

import middlewares "github.com/AnhHoangQuach/go-intern-spores/middlewares"

func (s *Server) initializeRoutes() {
	s.Router.HandleFunc("/users/signup", middlewares.SetMiddlewareJSON(s.signup)).Methods("POST")
	s.Router.HandleFunc("/home", middlewares.SetMiddlewareJSON(s.home)).Methods("GET")
}