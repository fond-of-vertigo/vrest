package main_test

import (
	"fmt"

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
