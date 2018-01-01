package simple

import (
	"net/http"
	"net/http/pprof"
)

func homepageHandler(w http.ResponseWriter, r *http.Request) {
	msg := "Welcome to Golang world!"
	w.Write([]byte(msg))
}

func main() {
	r := http.NewServeMux()
	r.HandleFunc("/", homepageHandler)

	// Register pprof handlers
	r.HandleFunc("/debug/pprof/", pprof.Index)
	r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	r.HandleFunc("/debug/pprof/profile", pprof.Profile)
	r.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	r.HandleFunc("/debug/pprof/trace", pprof.Trace)

	http.ListenAndServe(":8080", r)
}
