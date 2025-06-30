package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// Vulnerable: Hardcoded credentials for demonstration
const dbUser = "admin"
const dbPassword = "supersecret" // SonarQube will flag this

// hello handles the /hello endpoint
func hello(w http.ResponseWriter, req *http.Request) {
	span := tracer.StartSpan("web.request", tracer.ResourceName("/hello"))
	defer span.Finish()

	// Simulate logging credentials (another vulnerability)
	log.Printf("Connecting to DB with user=%s and password=%s\n", dbUser, dbPassword)

	rand.Seed(time.Now().UnixNano())

	if rand.Float32() > 0.80 {
		err := fmt.Errorf("500 - Something bad happened!")
		log.Println("500 - error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		span.Finish(tracer.WithError(err))
		return
	}

	log.Println("200 - Success")
	_, err := fmt.Fprintf(w, `{"route":"/hello","response":"hello world main branch"}`)
	if err != nil {
		log.Printf("Failed to write response: %v\n", err)
	}
}

// empty handles the / endpoint
func empty(w http.ResponseWriter, req *http.Request) {
	span := tracer.StartSpan("web.request", tracer.ResourceName("/"))
	defer span.Finish()

	rand.Seed(time.Now().UnixNano())

	if rand.Float32() > 0.98 {
		err := fmt.Errorf("500 - Something bad happened!")
		log.Println("500 - error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		span.Finish(tracer.WithError(err))
		return
	}

	log.Println("200 - Success")
	_, err := fmt.Fprintf(w, `{"route":"/","branch":"main"}`)
	if err != nil {
		log.Printf("Failed to write response: %v\n", err)
	}
}

func main() {
	tracer.Start()
	defer tracer.Stop()

	r := mux.NewRouter()
	r.HandleFunc("/", empty).Methods("GET")
	r.HandleFunc("/hello", hello).Methods("GET")

	log.Println("Starting server on :3000")
	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Fatalf("Server failed to start: %v\n", err)
	}
}
