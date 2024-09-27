package vrest

import (
	"errors"
	"testing"
)

func TestRequest_setResponseBody_NoPointer(t *testing.T) {
	client := New()

	body := []byte{}
	err := client.NewRequest().
		SetResponseBody(body).
		DoGet("/test")
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, ErrInvalidRequest) {
		t.Fatalf("unexpected error: %v", err)
	}
}
