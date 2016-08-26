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
    "gopkg.in/validator.v2"
    "github.com/gorilla/schema"
)

var (
	verbose bool
)

func email(w http.ResponseWriter, r *http.Request) {    
    msg := new(Email)
    dec := schema.NewDecoder()
    if e := dec.Decode(msg, r.PostForm); e != nil {
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte("Sorry undecodable input."))
        fmt.Printf("ERR:email:decode: %s\n", e)
        return
    }
    if e := validator.NewValidator().Validate(msg); e != nil {
        w.WriteHeader(http.StatusUnprocessableEntity)
        w.Write([]byte("Sorry invalid input."))
        fmt.Printf("ERR:email:validate: %s\n", e)
        return
    }

    mg := mailgun.NewMailgun(os.Getenv("MAILGUN_DOMAIN"), os.Getenv("MAILGUN_APIKEY"), os.Getenv("MAILGUN_PUBLICAPIKEY"))
    message := mailgun.NewMessage("noreply@rootdev.nl", "Contact request", "From=" + msg.Email + "\n\n" + msg.Body, "rootdev@gmail.com")
    if _, idx, e := mg.Send(message); e != nil {
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte("Sorry failed sending email."))
        fmt.Printf("ERR:email:send(%d): %s\n", idx, e)
        return
    }
    w.Write([]byte("Sent email."))
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