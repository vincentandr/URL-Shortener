package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

var port string = ":3000"

type Router struct {
	*mux.Router
}

// Router constructor
func NewRouter() *Router{
	r := &Router{
		Router: mux.NewRouter(),
	}
	r.routes()
	return r
}

func (r *Router) routes() {
	r.HandleFunc("/", Home).Methods("GET")
	r.HandleFunc("/getProducts", GetProducts).Methods("GET")
}

func main() {
	fmt.Println("Server starting at port " + port)

	r := NewRouter()

	http.ListenAndServe(port, r)
}