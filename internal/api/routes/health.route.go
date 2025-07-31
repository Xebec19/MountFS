package routes

import (
	"github.com/Xebec19/reimagined-lamp/internal/api/controllers"
	"github.com/gorilla/mux"
)

func createHealthRoutes(r *mux.Router) {

	s := r.PathPrefix("/health").Subrouter()

	s.HandleFunc("/test", controllers.TestHealth)
}
