package main

import "net/http"

func GoHandler(w http.ResponseWriter, r *http.Request) {
	message := r.FormValue("message")

	if message == "gorules" {
		w.Write([]byte("RIGHT"))
	} else {
		http.Error(w, "Error", http.StatusNotFound)
	}
}
