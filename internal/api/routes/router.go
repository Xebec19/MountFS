package routes

import "github.com/gorilla/mux"

func NewRouter() *mux.Router {

	r := mux.NewRouter()

	createHealthRoutes(r)

	return r
}
