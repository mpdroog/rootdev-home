// main contains a super lightweight HTTP-server
// so we can let this 'daemon' work with a really
// small attach vector (keeping maintenance super low)
package main

import (
	"net/http"
	"fmt"
	"flag"
	"github.com/unrolled/secure"
)

var (
	verbose bool
)

func email(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func main() {
	listen := ":8022"
	flag.BoolVar(&verbose, "v", false, "Verbose-mode (log more)")
	flag.Parse()

	secureMiddleware := secure.New(secure.Options{
        AllowedHosts:          []string{"rootdev.nl"},
        STSSeconds:            315360000,
        STSIncludeSubdomains:  true,
        STSPreload:            true,
        FrameDeny:             true,
        ContentTypeNosniff:    true,
        BrowserXssFilter:      true,
        ContentSecurityPolicy: "default-src 'self'",
    })

	fs := http.FileServer(http.Dir("static"))
  	http.Handle("/", secureMiddleware.Handler(fs))
    http.HandleFunc("/action/email", secureMiddleware.Handler(email))

    if verbose {
	    fmt.Printf("Listening on %s\n", listen)
	}
    if e := http.ListenAndServeTLS(listen, "tls/server.crt", "tls/server.key", app); e != nil {
    	panic(e)
    }
}