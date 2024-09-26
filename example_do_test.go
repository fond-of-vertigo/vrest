package vrest_test

import (
	"fmt"
	"net/http"
	"time"

	"github.com/fond-of-vertigo/vrest"
)

type testBody struct {
	Text   string
	Number int
}

func ExampleMockDoer() {
	client := vrest.NewWithClient(&http.Client{
		Timeout: 10 * time.Second,
	})

	client.Overridable.Do = vrest.MockDoer(testBody{Text: "text", Number: 456}, nil)

	respBody := testBody{}

	err := client.NewRequest().
		SetResponseBody(&respBody).
		DoGet("https://jsonplaceholder.typicode.com/todos/1")
	if err != nil {
		panic(err)
	}

	fmt.Println("Response:", respBody)
	// Output: Response: {text 456}
}
