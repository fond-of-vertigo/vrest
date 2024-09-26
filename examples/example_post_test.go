package main_test

import (
	"fmt"
	"net/url"

	"github.com/fond-of-vertigo/vrest"
)

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
