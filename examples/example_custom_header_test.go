package main_test

import (
	"fmt"

	"github.com/fond-of-vertigo/vrest"
)

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
