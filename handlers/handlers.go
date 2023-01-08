package handlers

import (
	"fmt"
	"github.com/bernsblack/fiddleware/examples"
	"github.com/gorilla/mux"
	"net/http"
)

func Ping(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "pong")
	if err != nil {
		panic(err)
	}
}

func GetCircle(w http.ResponseWriter, r *http.Request) {
	c := examples.Circle{Radius: 234.3243}
	// response with json of c
	_, _ = fmt.Fprintf(w, "%+v", c)
}

func PingWithId(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	_, _ = fmt.Fprintf(w, fmt.Sprintf("%+v", map[string]interface{}{"id": id}))
}

func HandlerWithoutMiddleware(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, "pong 3")
}
