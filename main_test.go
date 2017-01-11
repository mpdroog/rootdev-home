package main

import (
	"bytes"
	"net/http"
	"net/url"
	"testing"
)

func TestEmail(t *testing.T) {
	data := url.Values{"email": {"derp@derp.com"}, "body": {"Hello world"}}

	req, e := http.NewRequest("POST", "/action/email", bytes.NewBufferString(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if e != nil {
		t.Fatal(e)
	}

	email, e := emailDecode(req)
	if e != nil {
		t.Fatal(e)
	}
	if email.Email != "derp@derp.com" {
		t.Fatalf("Email not expected, found=" + email.Email)
	}
	if email.Body != "Hello world" {
		t.Fatalf("Body not expected, found=" + email.Body)
	}
}
