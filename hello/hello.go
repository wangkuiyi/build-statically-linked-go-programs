package main

import (
    "fmt"
    "net/http"
    "log"
)


func main() {
    http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
	    fmt.Fprintf(w, "Hello!")
    })

    if e := http.ListenAndServe(":9090", nil); e != nil {
        log.Fatal("ListenAndServe: ", e)
    }
}
