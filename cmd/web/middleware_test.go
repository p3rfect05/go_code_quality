package main

import (
	"net/http"
	"testing"
)

func TestNoSurf(t *testing.T) {
	var myH myHandler

	h := NoSurf(&myH)

	switch v := h.(type) {
	case http.Handler:
	//do nothing
	default:
		t.Error("Type is not http.Handler but %T", v)
	}
}

func TestSessionLoad(t *testing.T) {
	var myH myHandler

	h := SessionLoad(&myH)

	switch v := h.(type) {
	case http.Handler:
	//do nothing
	default:
		t.Error("Type is not http.Handler but %T", v)
	}
}
