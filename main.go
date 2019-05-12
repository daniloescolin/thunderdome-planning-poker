package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gobuffalo/packr"
	"github.com/gorilla/mux"
)

func main() {
	var listenPort = fmt.Sprintf(":%s", GetEnv("PORT", "8080"))
	go h.run()

	// box := packr.New("webui", "./dist")
	box := packr.NewBox("./dist")
	staticHandler := http.FileServer(box)

	router := mux.NewRouter()
	router.PathPrefix("/css/").Handler(staticHandler)
	router.PathPrefix("/js/").Handler(staticHandler)
	router.PathPrefix("/img/").Handler(staticHandler)
	router.HandleFunc("/api/warrior", RecruitWarriorHandler).Methods("POST")
	router.HandleFunc("/api/battle", CreateBattleHandler).Methods("POST")
	router.HandleFunc("/api/battle/{id}", GetBattleHandler)
	router.HandleFunc("/api/arena/{id}", serveWs)
	router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = "/"
		staticHandler.ServeHTTP(w, r)
	})

	srv := &http.Server{
		Handler: router,
		Addr:    listenPort,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("Access the WebUI via 127.0.0.1" + listenPort)

	log.Fatal(srv.ListenAndServe())
}