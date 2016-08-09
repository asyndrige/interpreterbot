package main

import (
	"log"
	"net/http"
)

// Decorator function for restricting http verb-methods.
func httpRestrict(h http.HandlerFunc, verb []string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		for _, v := range verb {
			if r.Method != v {
				http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
				return
			}
			h(w, r)
		}
	}
}

// Dummy error checking function.
func checkErr(err error) {
	if err != nil {
		log.Println(err)
	}
}
