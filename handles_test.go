package main

import (
	"testing"
)

func TestBadWordReplacement(t *testing.T) {
	post := "Hello, this is a Fornax post, kerfuffle these Sharbert words."

	got := badWordReplacement(post)
	result := "Hello, this is a **** post, **** these **** words."

	if got != result {
		t.Errorf("got: %s -- wanted: %s", got, result)
	}
}