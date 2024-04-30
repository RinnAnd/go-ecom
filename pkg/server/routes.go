package server

import (
	"github.com/gorilla/mux"
)

func (s *Server) RegisterRoutes() {
	router := mux.NewRouter()

	router.HandleFunc("/products", s.MakeProduct).Methods("POST")
	router.HandleFunc("/products", s.GetProducts).Methods("GET")
	router.HandleFunc("/products/{id}", s.GetProduct).Methods("GET")
	router.HandleFunc("/products/{id}", s.EditProduct).Methods("PUT")
	router.HandleFunc("/products/{id}", s.DeleteProduct).Methods("DELETE")

	s.httpServer.Handler = router
}
