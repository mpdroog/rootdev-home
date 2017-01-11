package main

import (
	"testing"
)

func TestValidateEmail(t *testing.T) {
	// validateEmail(v interface{}, param string) error
	valid := []string{
		"pw.droog@quicknet.nl",
		"rootdev@gmail.com",
		"mpdroog@icloud.com",
		"rootdev!derp@gmail.com",
	}
	for _, email := range valid {
		if e := validateEmail(email, ""); e != nil {
			t.Errorf("Email should match: %s", email)
		}
	}

	invalid := []string{
		"mp",
		"mp.droog",
		"a@b",
	}
	for _, email := range invalid {
		if e := validateEmail(email, ""); e == nil {
			t.Errorf("Email should NOT match: %s", email)
		}
	}
}
