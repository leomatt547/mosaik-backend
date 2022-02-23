package controllers

import "gitlab.informatika.org/if3250_2022_37_mosaik/mosaik-backend/src/api/middlewares"

func (s *Server) initializeRoutes() {

	// Home Route
	s.Router.HandleFunc("/", middlewares.SetMiddlewareJSON(s.Home)).Methods("GET")

	// Login Route
	s.Router.HandleFunc("/login", middlewares.SetMiddlewareJSON(s.Login)).Methods("POST")

	//Parents routes
	s.Router.HandleFunc("/parents", middlewares.SetMiddlewareJSON(s.CreateParent)).Methods("POST")
	s.Router.HandleFunc("/parents", middlewares.SetMiddlewareJSON(s.GetParents)).Methods("GET")
	s.Router.HandleFunc("/parents/{id}", middlewares.SetMiddlewareJSON(s.GetParent)).Methods("GET")
	s.Router.HandleFunc("/parents/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateParent))).Methods("PUT")
	s.Router.HandleFunc("/parents/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteParent)).Methods("DELETE")

	//Child routes
	s.Router.HandleFunc("/childs", middlewares.SetMiddlewareJSON(s.CreateChild)).Methods("POST")
	s.Router.HandleFunc("/childs", middlewares.SetMiddlewareJSON(s.GetChilds)).Methods("GET")
	s.Router.HandleFunc("/childs/{id}", middlewares.SetMiddlewareJSON(s.GetChild)).Methods("GET")
	s.Router.HandleFunc("/childs/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateChild))).Methods("PUT")
	s.Router.HandleFunc("/childs/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteChild)).Methods("DELETE")

	//Child Visit routes
	s.Router.HandleFunc("/childvisits", middlewares.SetMiddlewareJSON(s.CreateChildVisit)).Methods("POST")
	s.Router.HandleFunc("/childvisits", middlewares.SetMiddlewareJSON(s.GetChildVisits)).Methods("GET")
	s.Router.HandleFunc("/childvisits/{id}", middlewares.SetMiddlewareJSON(s.GetChildVisit)).Methods("GET")
	s.Router.HandleFunc("/childvisits/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteChildVisit)).Methods("DELETE")

	//Parent Visit routes
	s.Router.HandleFunc("/parentvisits", middlewares.SetMiddlewareJSON(s.CreateParentVisit)).Methods("POST")
	s.Router.HandleFunc("/parentvisits", middlewares.SetMiddlewareJSON(s.GetParentVisits)).Methods("GET")
	s.Router.HandleFunc("/parentvisits/{id}", middlewares.SetMiddlewareJSON(s.GetParentVisit)).Methods("GET")
	s.Router.HandleFunc("/parentvisits/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteParentVisit)).Methods("DELETE")

	//Url routes
	s.Router.HandleFunc("/urls", middlewares.SetMiddlewareJSON(s.CreateUrl)).Methods("POST")
	s.Router.HandleFunc("/urls", middlewares.SetMiddlewareJSON(s.GetUrls)).Methods("GET")
	s.Router.HandleFunc("/urls/{id}", middlewares.SetMiddlewareJSON(s.GetUrl)).Methods("GET")
}
