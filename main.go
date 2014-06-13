package main

import (
	"flag"
	"fmt"
	"github.com/ActiveState/tail"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

var port = flag.Int("http-port", 24987, "Http Port for claudius")

func Ping(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprint(w, "ping")
}

func Log(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	containerId := ps.ByName("containerId")
	fmt.Fprintf(w, "Logs for %s", containerId)
}

func Tail(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	containerId := ps.ByName("containerId")
	logFilePath := fmt.Sprintf("/var/lib/docker/containers/%s/%s-json.log", containerId, containerId)
	hj := w.(http.Hijacker)
	conn, bufrw, err := hj.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer conn.Close()
	t, err := tail.TailFile(logFilePath, tail.Config{Follow: true})
	for line := range t.Lines {
		bufrw.WriteString(fmt.Sprintf("%s\n", line.Text))
		bufrw.Flush()
	}
}

func main() {
	flag.Parse()
	httpPort := fmt.Sprintf(":%d", *port)
	router := httprouter.New()
	router.GET("/ping", Ping)
	router.GET("/logs/:containerId", Log)
	router.GET("/tail/:containerId", Tail)
	log.Fatal(http.ListenAndServe(httpPort, router))
}
