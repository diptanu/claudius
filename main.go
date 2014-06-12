package main

import(
    "fmt"
    "github.com/julienschmidt/httprouter"
    "net/http"
    "log"
)


func Ping(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    fmt.Fprint(w, "ping")
}

func Log(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    fmt.Fprintf(w, "Logs for %s", ps.ByName("containerId"))
}


func main() {
    router := httprouter.New()
    router.GET("/ping", Ping)
    router.GET("/logs/:containerId", Log)
    log.Fatal(http.ListenAndServe(":8080", router))
}
