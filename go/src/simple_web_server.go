package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

var msg string = "Hello, I'm your webserver today!!! My name is: %s.\n"

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	io.WriteString(w, msg)
}

func getServerName(args []string) string {
	if len(args) > 0 {
		return args[0]
	}
	return ""
}

func main() {
	serverName := getServerName(os.Args[1:])
	if serverName == "" {
		msg = fmt.Sprintf(msg, "uhhhh...yeah, I forgot.")
	} else {
		msg = fmt.Sprintf(msg, serverName)
	}
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":80", nil))
}
