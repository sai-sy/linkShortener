package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
		fmt.Fprint(w, "Hellooo World!\n")
	})
	
	http.ListenAndServe(":8080", nil)
}
