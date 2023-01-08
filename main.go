package main

import (
	"fmt"
	"github.com/bernsblack/fiddleware/fiddleware"
	"github.com/bernsblack/fiddleware/handlers"
	"github.com/bernsblack/fiddleware/util"
	"github.com/gorilla/mux"
	"net/http"
	"runtime"
)

func main() {
	runtimeVersion := runtime.Version()
	fmt.Printf("runtime.Version() => %s\n", runtimeVersion)

	logSettings := map[string]struct{}{
		"/ping":      {},
		"/ping/{id}": {},
		"/circle":    {},
	}
	fmt.Printf("%+v", logSettings)

	logger := util.NewLogger()
	recorder := fiddleware.NewInMemoryRecorder(logger)

	innerRouter := mux.NewRouter()
	router := fiddleware.Wrap(innerRouter, logger, recorder, logSettings)

	// setup http server with ping endpoint
	router.HandleFunc("/ping", handlers.Ping)
	router.HandleFunc("/ping/{id}", handlers.PingWithId)
	router.HandleFunc("/circle", handlers.GetCircle).Methods(http.MethodGet)
	router.HandleFunc("/handler-without-middleware", handlers.HandlerWithoutMiddleware)

	// start http server
	err := http.ListenAndServe(":7789", router)
	if err != nil {
		panic(err)
	}

	fmt.Printf("test => %s", "test")
}
