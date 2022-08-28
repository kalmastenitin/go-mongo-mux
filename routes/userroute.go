package routes

import (
	"mux-mongo-api/controllers"

	"github.com/gorilla/mux"
)

func UserRoute(router *mux.Router) {
	router.HandleFunc("/user/register", controllers.Register).Methods("POST")
	router.HandleFunc("/user/{userId}", controllers.GetUser).Methods("GET")
	router.HandleFunc("/user/", controllers.GetAllUser).Methods("GET")
	router.HandleFunc("/user/{userId}", controllers.DeleteUser).Methods("DELETE")
	router.HandleFunc("/user/activate/{userId}", controllers.ActivateUser).Methods("POST")
	router.HandleFunc("/user/login", controllers.LoginUser).Methods("POST")
}
