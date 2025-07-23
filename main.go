// @title TZ
// @version 0.1
// @description completed test task for junior go developer

// @host localhost:8080
// @BasePath /
// @query.collection.format multi

package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	config, err := LoadConfig()
	if err != nil {
		log.Fatal("Error loading config: ", err)
	}

	db, err := NewDB(fmt.Sprintf("host=db user=%s password=%s dbname=%s port=5432 sslmode=disable", config.User, config.Password, config.Name))
	if err != nil {
		log.Fatal("failed to init DB", err)
	}

	subscriptionRepository := NewSubscriptionRepository(db)
	subscriptionHandler := NewSubscriptionHandler(subscriptionRepository)

	router := mux.NewRouter()
	subscriptionHandler.RegisterRoutes(router)

	router.HandleFunc("/swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./docs/swagger.json")
	})

	router.PathPrefix("/swagger/").Handler(httpSwagger.Handler())

	log.Printf("Listening on port: %d...\n", config.Port)

	addr := fmt.Sprintf(":%d", config.Port)
	log.Fatal(http.ListenAndServe(addr, router))
}
