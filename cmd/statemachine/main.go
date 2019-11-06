package main

import (
	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
	"github.com/olivernadj/post-proc/internal/api"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)


type Specification struct {
	TTProcess   int `envconfig:"TT_PROCESS" required:"false" default:"10"`
	TTDelete	int `envconfig:"TT_DELETE" required:"false" default:"10"`
}

func main() {
	var s Specification
	err := envconfig.Process("", &s)
	if err != nil {
		log.Fatal(err.Error())
	}

	//check for action timeouts in every 5 seconds
	ticker := time.NewTicker(5 * time.Second)
	// trap Ctrl+C and call cancel
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	defer func() {
		log.Printf("signal stop")
		signal.Stop(c)
	}()
	go func() {
		for {
			select {
			case <-ticker.C:
				//handle events
				api.HandleProcess(s.TTDelete)
				api.HandleDelete(s.TTProcess)
			case <-c:
				log.Printf("gracefully stop statemachinge")
				ticker.Stop()
				return
			}
		}
	}()

	// http for metrics only
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/metrics", metricsHandler)
	log.Fatal(http.ListenAndServe(":8082", r))
}

func metricsHandler (w http.ResponseWriter, r *http.Request) {
	p := promhttp.Handler()
	p.ServeHTTP(w, r)
}
