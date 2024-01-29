package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"

	"github.com/gorilla/mux"
	"github.com/newrelic/go-agent/v3/newrelic"
	"go.elastic.co/apm/module/apmgorilla"
)

func hello(w http.ResponseWriter, req *http.Request) {
	span := tracer.StartSpan("web.request", tracer.ResourceName("/"))
	defer span.Finish()
	rand.Seed(time.Now().UnixNano())
	if rand.Float32() > 0.80 {
		fmt.Println("500 - error Something bad happened!")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - error Something bad happened!"))
		span.Finish(tracer.WithError(fmt.Errorf("500 - Something bad happened!")))
	} else {
		fmt.Println("200 - Success")
		fmt.Fprintf(w, `{"route":"/hello","response":"hello world main branch"}`)
	}
}

func empty(w http.ResponseWriter, req *http.Request) {
	span := tracer.StartSpan("web.request", tracer.ResourceName("/hello"))
	defer span.Finish()
	rand.Seed(time.Now().UnixNano())
	if rand.Float32() > 0.98 {
		fmt.Println("500 - error Something bad happened!")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - error Something bad happened!"))
		span.Finish(tracer.WithError(fmt.Errorf("500 - Something bad happened!")))
	} else {
		fmt.Println("200 - Success")
		fmt.Fprintf(w, `{"route":"/","branch":"main"}`)
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", empty)
	r.HandleFunc("/hello", hello)
	fmt.Print(http.ListenAndServe(":3000", r))
}
