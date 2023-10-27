package main

// May be provided automatically by your editor, so they may be removed if you hit save right after typing this bit.
import (
	"fmt"
	"net/http"
)

func returnResponse(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<div id=\"response\">You said: %s!</div>", r.Header.Get("Hx-Prompt"))
}

func main() {
	
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	    http.ServeFile(w, r, "index.html")
	})

	http.HandleFunc("/submit", returnResponse)
	fmt.Println("Starting the server on :3000...")
	http.ListenAndServe("127.0.0.1:3000", nil)
}
