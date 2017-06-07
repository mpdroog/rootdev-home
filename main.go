// main contains a super lightweight HTTP-server
// so we can let this 'daemon' work with a really
// small attach vector (keeping maintenance super low)
package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/schema"
	"gopkg.in/mailgun/mailgun-go.v1"
	"gopkg.in/validator.v2"
	"net/http"
	"os"
)

var (
	verbose       bool
	mailgunDomain string
	mailgunApi    string
	mailgunApiPub string

	decode = schema.NewDecoder()
)

func emailDecode(r *http.Request) (Email, error) {
	msg := Email{}

	if r.Body == nil {
		return msg, fmt.Errorf("Empty body")
	}
	if e := r.ParseForm(); e != nil {
		return msg, e
	}

	if e := decode.Decode(&msg, r.PostForm); e != nil {
		return msg, e
	}
	if e := validator.Validate(msg); e != nil {
		return msg, e
	}
	return msg, nil
}

func email(w http.ResponseWriter, r *http.Request) {
	msg, e := emailDecode(r)
	if e != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("Sorry invalid input."))
		fmt.Printf("ERR:email:validate: %s\n", e)
		return
	}

	mg := mailgun.NewMailgun(mailgunDomain, mailgunApi, mailgunApiPub)
	message := mailgun.NewMessage("noreply@rootdev.nl", "Contact request", "From "+msg.Email+"\n\n"+msg.Body, "rootdev@gmail.com")

	if _, idx, e := mg.Send(message); e != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Sorry failed sending email."))
		fmt.Printf("ERR:email:send(%s): %s\n", idx, e)
		return
	}
	w.Write([]byte("Sent email."))
}

func main() {
	listen := ":80"
	flag.BoolVar(&verbose, "v", false, "Verbose-mode (log more)")
	flag.StringVar(&listen, "l", ":80", "Listener port")
	flag.Parse()

	mailgunDomain = os.Getenv("MAILGUN_DOMAIN")
	mailgunApi = os.Getenv("MAILGUN_APIKEY")
	mailgunApiPub = os.Getenv("MAILGUN_PUBLICAPIKEY")

	if mailgunDomain == "" {
		fmt.Printf("ERR:main:env: Missing MAILGUN_DOMAIN\n")
		return
	}
	if mailgunApi == "" {
		fmt.Printf("ERR:main:env: Missing MAILGUN_APIKEY\n")
		return
	}
	if mailgunApiPub == "" {
		fmt.Printf("ERR:main:env: Missing MAILGUN_PUBLICAPIKEY\n")
		return
	}

	fs := http.FileServer(http.Dir("build"))
	http.Handle("/", fs)
	http.HandleFunc("/cloud/init", cloudinit)
	http.HandleFunc("/cloud/ipxe", cloudinit)
	http.HandleFunc("/action/email", email)

	if verbose {
		fmt.Printf("Listening on %s\n", listen)
	}
	if e := http.ListenAndServe(listen, nil); e != nil {
		panic(e)
	}
}
