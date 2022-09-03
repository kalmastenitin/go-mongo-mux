package main

import (
	"log"
	"mux-mongo-api/configs"
	"mux-mongo-api/routes"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	configs.ConnectDB()

	routes.UserRoute(router)
	router.Use(mux.CORSMethodMiddleware(router))
	log.Println("Server Started Successfully!")
	log.Fatal(http.ListenAndServe(":8000", router))
}
