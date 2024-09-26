package vrest_test

import (
	"fmt"
	"net/url"

	"github.com/fond-of-vertigo/vrest"
)

func ExampleRequest_DoGet() {
	client := vrest.New()
	respBody := make(map[string]interface{})

	err := client.NewRequest().
		SetResponseBody(&respBody).
		DoGet("https://jsonplaceholder.typicode.com/todos/1")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Response:", respBody)
	// Output: Response: map[completed:false id:1 title:delectus aut autem userId:1]
}

func ExampleRequest_DoPost() {
	client := vrest.New()
	respBody := make(map[string]interface{})

	body := url.Values{}
	body.Set("title", "foo")
	body.Set("body", "bar")
	body.Set("userId", "1")

	err := client.NewRequest().
		SetResponseBody(&respBody).
		SetContentTypeJSON().
		SetBody(body).
		DoPost("https://jsonplaceholder.typicode.com/posts")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Response:", respBody)
	// Output: Response: map[body:[bar] id:101 title:[foo] userId:[1]]
}

func ExampleRequest_SetBasicAuth() {
	client := vrest.New()
	respBody := make(map[string]interface{})

	err := client.NewRequest().
		SetResponseBody(&respBody).
		SetBasicAuth("username", "password").
		DoGet("https://jsonplaceholder.typicode.com/todos/1")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Response:", respBody)
	// Output: Response: map[completed:false id:1 title:delectus aut autem userId:1]
}

func ExampleRequest_SetHeader() {
	client := vrest.New()
	respBody := make(map[string]interface{})

	err := client.NewRequest().
		SetResponseBody(&respBody).
		SetHeader("my-header", "my-value").
		DoGet("https://jsonplaceholder.typicode.com/todos/1")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Response:", respBody)
	// Output: Response: map[completed:false id:1 title:delectus aut autem userId:1]
}
