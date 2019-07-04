package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

var serverName string = ""

func computeBomb(cycles int) {
	for i := 0; i < cycles; i++ {
		randF1 := rand.Float64()
		randF2 := rand.Float64()
		randF2 = (randF1 * randF2) - randF2
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	msg := "Hello, I'm your webserver today!!!\nMy name is, %s.\n"
	query := r.URL.Query()
	sleepTime := 0
	sleepTimeStr := ""
	cycles := 0
	cyclesStr := ""

	if serverName == "" {
		msg = fmt.Sprintf(msg, "uhhhh...yeah, I forgot.")
	} else {
		msg = fmt.Sprintf(msg, serverName)
	}

	if sleepTimeStr = query.Get("sleepTime"); sleepTimeStr != "" {
		sleepTimeC, err := strconv.Atoi(sleepTimeStr)
		sleepTime = sleepTimeC
		if sleepTime < 0 || err != nil {
			log.Print("sleepTime Err non-null")
			sleepTime = 0
		}
	}

	if cyclesStr = query.Get("cycles"); cyclesStr != "" {
		cyclesC, err := strconv.Atoi(cyclesStr)
		cycles = cyclesC
		if cycles < 0 || err != nil {
			cycles = 0
		}
	}

	computeBomb(cycles)
	time.Sleep(time.Duration(sleepTime) * time.Millisecond)
	msg += "I've slept " + strconv.Itoa(sleepTime) + " milliseconds.\n"
	msg += "I've computed " + strconv.Itoa(cycles) + " times.\n"
	msg += "Goodbye!\n"
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
	serverName = getServerName(os.Args[1:])
	http.HandleFunc("/request", handler)
	log.Fatal(http.ListenAndServe(":80", nil))
}
