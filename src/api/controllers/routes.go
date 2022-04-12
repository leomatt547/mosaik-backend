package controllers

import "gitlab.informatika.org/if3250_2022_37_mosaik/mosaik-backend/src/api/middlewares"

func (s *Server) initializeRoutes() {

	// Home Route
	s.Router.HandleFunc("/", middlewares.SetMiddlewareHTML(s.Home)).Methods("GET")

	// Login Route
	s.Router.HandleFunc("/login", middlewares.SetMiddlewareJSON(s.Login)).Methods("POST")

	//Parents routes
	s.Router.HandleFunc("/resetpassword", middlewares.SetMiddlewareJSON(s.SendMailParent)).Methods("POST")
	s.Router.HandleFunc("/parents", middlewares.SetMiddlewareJSON(s.CreateParent)).Methods("POST")
	s.Router.HandleFunc("/parents/password/{id}", middlewares.SetMiddlewareJSON(s.UpdateParentPassword)).Methods("POST")
	s.Router.HandleFunc("/parents", middlewares.SetMiddlewareJSON(s.GetParents)).Methods("GET")
	s.Router.HandleFunc("/parents/{id}", middlewares.SetMiddlewareJSON(s.GetParent)).Methods("GET")
	s.Router.HandleFunc("/parents/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateParentProfile))).Methods("PUT")
	s.Router.HandleFunc("/parents/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteParent)).Methods("DELETE")
	s.Router.HandleFunc("/parents/newpassword/{id}", middlewares.SetMiddlewareJSON(s.ParentNewPassword)).Methods("POST")

	//Child routes
	// s.Router.HandleFunc("/resetpassword", middlewares.SetMiddlewareJSON(s.SendMailChild)).Methods("POST")
	s.Router.HandleFunc("/childs", middlewares.SetMiddlewareJSON(s.CreateChild)).Methods("POST")
	// s.Router.HandleFunc("/childs/password/{id}", middlewares.SetMiddlewareJSON(s.UpdateChildPassword)).Methods("POST")
	s.Router.HandleFunc("/childs/password/{id}", middlewares.SetMiddlewareJSON(s.ChildNewPassword)).Methods("POST")
	s.Router.HandleFunc("/childs", middlewares.SetMiddlewareJSON(s.GetChilds)).Methods("GET")
	s.Router.HandleFunc("/childs/{id}", middlewares.SetMiddlewareJSON(s.GetChild)).Methods("GET")
	s.Router.HandleFunc("/childs/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateChildProfile))).Methods("PUT")
	s.Router.HandleFunc("/childs/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteChild)).Methods("DELETE")

	//Parent Visit routes
	s.Router.HandleFunc("/parentvisits", middlewares.SetMiddlewareJSON(s.CreateParentVisit)).Methods("POST")
	s.Router.HandleFunc("/parentvisits", middlewares.SetMiddlewareJSON(s.GetParentVisits)).Methods("GET")
	s.Router.HandleFunc("/parentvisits/{id}", middlewares.SetMiddlewareJSON(s.GetParentVisit)).Methods("GET")
	s.Router.HandleFunc("/parentvisits/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteParentVisit)).Methods("DELETE")

	//Child Visit routes
	s.Router.HandleFunc("/childvisits", middlewares.SetMiddlewareJSON(s.CreateChildVisit)).Methods("POST")
	s.Router.HandleFunc("/childvisits", middlewares.SetMiddlewareJSON(s.GetChildVisits)).Methods("GET")
	s.Router.HandleFunc("/childvisits/{id}", middlewares.SetMiddlewareJSON(s.GetChildVisit)).Methods("GET")
	s.Router.HandleFunc("/childvisits/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteChildVisit)).Methods("DELETE")

	//Parent Download routes
	s.Router.HandleFunc("/parentdownloads", middlewares.SetMiddlewareJSON(s.CreateParentDownload)).Methods("POST")
	s.Router.HandleFunc("/parentdownloads", middlewares.SetMiddlewareJSON(s.GetParentDownloads)).Methods("GET")
	s.Router.HandleFunc("/parentdownloads/{id}", middlewares.SetMiddlewareJSON(s.GetParentDownload)).Methods("GET")
	s.Router.HandleFunc("/parentdownloads/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteParentDownload)).Methods("DELETE")

	//Child Download routes
	s.Router.HandleFunc("/childdownloads", middlewares.SetMiddlewareJSON(s.CreateChildDownload)).Methods("POST")
	s.Router.HandleFunc("/childdownloads", middlewares.SetMiddlewareJSON(s.GetChildDownloads)).Methods("GET")
	s.Router.HandleFunc("/childdownloads/{id}", middlewares.SetMiddlewareJSON(s.GetChildDownload)).Methods("GET")
	s.Router.HandleFunc("/childdownloads/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteChildDownload)).Methods("DELETE")

	//Url routes
	s.Router.HandleFunc("/urls", middlewares.SetMiddlewareJSON(s.CreateUrl)).Methods("POST")
	s.Router.HandleFunc("/urls", middlewares.SetMiddlewareJSON(s.GetUrls)).Methods("GET")
	s.Router.HandleFunc("/urls/{id}", middlewares.SetMiddlewareJSON(s.GetUrl)).Methods("GET")

	//Web Checker
	s.Router.HandleFunc("/nsfw", middlewares.SetMiddlewareJSON(s.SavedSearchChecker)).Methods("POST")
	s.Router.HandleFunc("/nsfw/images", middlewares.SetMiddlewareJSON(s.GetImagesInLink)).Methods("POST")
	s.Router.HandleFunc("/nsfw", middlewares.SetMiddlewareJSON(s.GetNSFWUrls)).Methods("GET")

	//Blacklist
	s.Router.HandleFunc("/blacklist", middlewares.SetMiddlewareJSON(s.CreateBlackList)).Methods("POST")
	s.Router.HandleFunc("/blacklist", middlewares.SetMiddlewareJSON(s.GetBlacklist)).Methods("GET")
	s.Router.HandleFunc("/blacklist/{id}", middlewares.SetMiddlewareJSON(s.GetBlacklistByID)).Methods("GET")
	s.Router.HandleFunc("/blacklist/check", middlewares.SetMiddlewareJSON(s.BlacklistChecker)).Methods("POST")
	s.Router.HandleFunc("/blacklist/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteBlacklist)).Methods("DELETE")

	//Whitelist
	s.Router.HandleFunc("/whitelist", middlewares.SetMiddlewareJSON(s.CreateWhitelist)).Methods("POST")
	s.Router.HandleFunc("/whitelist", middlewares.SetMiddlewareJSON(s.GetWhitelist)).Methods("GET")
	s.Router.HandleFunc("/whitelist/{id}", middlewares.SetMiddlewareJSON(s.GetWhitelistByID)).Methods("GET")
	s.Router.HandleFunc("/whitelist/check", middlewares.SetMiddlewareJSON(s.WhitelistChecker)).Methods("POST")
	s.Router.HandleFunc("/whitelist/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteWhitelist)).Methods("DELETE")
}
