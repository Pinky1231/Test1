package main

import (
	"fmt"
	"net/http"
	"slices"
)

func handler(w http.ResponseWriter, r *http.Request) {
	vals := r.URL.Query()
	fmt.Println(vals)
	if slices.Contains(vals["name"], "Oleg") {
		w.Write([]byte("Hiii OOleg"))
		return
	} else {
		w.Write([]byte(vals["name"][0]))
		return
	}

	// w.Write([]byte("Вы на главной странице!"))

}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/ss", handler)
	http.ListenAndServe(":8080", mux)

}
