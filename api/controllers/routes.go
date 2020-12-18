package controllers

import "github.com/dmvvilela/go-poc-2/api/middlewares"

func (s *Server) initializeRoutes() {
	// Home Route
	s.Router.HandleFunc("/", middlewares.SetMiddlewareJSON(s.Home)).Methods("GET")

	// Contact routes
	s.Router.HandleFunc("/contacts", middlewares.SetMiddlewareJSON(s.GetContacts)).Methods("GET")
	s.Router.HandleFunc("/contacts", middlewares.SetMiddlewareJSON(s.CreateContact)).Methods("POST")
	s.Router.HandleFunc("/contacts/{id}", middlewares.SetMiddlewareJSON(s.GetContact)).Methods("GET")
	s.Router.HandleFunc("/contacts/{id}", middlewares.SetMiddlewareJSON(s.UpdateContact)).Methods("PUT")
	s.Router.HandleFunc("/contacts/{id}", s.DeleteContact).Methods("DELETE")
}
