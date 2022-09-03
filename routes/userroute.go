package routes

import (
	"encoding/json"
	"mux-mongo-api/controllers"
	"mux-mongo-api/helpers"
	"mux-mongo-api/responses"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")
		if len(authHeader) != 2 {
			w.WriteHeader(http.StatusUnauthorized)
			response := responses.UserResponse{Status: http.StatusUnauthorized, Message: "error", Data: map[string]interface{}{"data": "Invalid Token"}}
			json.NewEncoder(w).Encode(response)
			return
		}
		tokenState := helpers.ValidateAccessToken(authHeader[1])
		if !tokenState {
			w.WriteHeader(http.StatusUnauthorized)
			response := responses.UserResponse{Status: http.StatusUnauthorized, Message: "error", Data: map[string]interface{}{"data": "Token Expired"}}
			json.NewEncoder(w).Encode(response)
			return
		}
		return
	})
}

func UserRoute(router *mux.Router) {
	router.HandleFunc("/user/register", controllers.Register).Methods("POST")
	router.HandleFunc("/user/{userId}", controllers.GetUser).Methods("GET")
	router.HandleFunc("/user/", controllers.GetAllUser).Methods("GET")
	router.Handle("/user/{userId}", middleware(http.HandlerFunc(controllers.DeleteUser))).Methods("DELETE")
	// router.HandleFunc("/user/{userId}", controllers.DeleteUser).Methods("DELETE")
	router.HandleFunc("/user/activate/{userId}", controllers.ActivateUser).Methods("POST")
	router.HandleFunc("/user/login", controllers.LoginUser).Methods("POST")
}
