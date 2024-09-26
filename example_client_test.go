package vrest_test

import (
	"fmt"
	"net/http"
	"time"

	"github.com/fond-of-vertigo/vrest"
)

func ExampleNewWithClient() {
	client := vrest.NewWithClient(&http.Client{
		Timeout: 10 * time.Second,
	})

	respBody := make(map[string]interface{})

	err := client.NewRequest().
		SetResponseBody(&respBody).
		DoGet("https://jsonplaceholder.typicode.com/todos/1")
	if err != nil {
		panic(err)
	}

	fmt.Println("Response:", respBody)
	// Output: Response: map[completed:false id:1 title:delectus aut autem userId:1]
}
