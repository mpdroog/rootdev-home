// main contains a super lightweight HTTP-server
// so we can let this 'daemon' work with a really
// small attach vector (keeping maintenance super low)
package main

import (
    "os"
	"net/http"
	"fmt"
	"flag"
    "gopkg.in/mailgun/mailgun-go.v1"
)

var (
	verbose bool
)

func email(w http.ResponseWriter, r *http.Request) {    
    mg := mailgun.NewMailgun(os.Getenv("MAILGUN_DOMAIN"), os.Getenv("MAILGUN_APIKEY"), os.Getenv("MAILGUN_PUBLICAPIKEY"))
    message := mailgun.NewMessage("noreply@rootdev.nl", "Contact request", "Hello from Mailgun Go!", "rootdev@gmail.com")

    mg.Send(message)
}

func main() {
	listen := ":8022"
	flag.BoolVar(&verbose, "v", false, "Verbose-mode (log more)")
	flag.Parse()

	fs := http.FileServer(http.Dir("build"))
  	http.Handle("/", fs)
    http.HandleFunc("/action/email", email)

    if verbose {
	    fmt.Printf("Listening on %s\n", listen)
	}
    if e := http.ListenAndServe(listen, nil); e != nil {
    	panic(e)
    }
}