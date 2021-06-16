package router

import (
	"crud-backend/middleware"

	"github.com/gorilla/mux"
)

func Router() *mux.Router {

	router := mux.NewRouter()

	router.HandleFunc("/api/bbm/{id}", middleware.GetBbm).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/bbm", middleware.GetAllBbm).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/newbbm", middleware.CreateBbm).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/bbm/{id}", middleware.UpdateBbm).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/bbm/{id}", middleware.DeleteBbm).Methods("DELETE", "OPTIONS")

	return router
}
