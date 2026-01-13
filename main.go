package main

import (
	"fmt"
	"net/http"
)

func main() {
	sm := http.NewServeMux()
	sm.HandleFunc("/", handlerIndex)
	http.ListenAndServe(":8080", sm)
}

func handlerIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "hello world index")
}
