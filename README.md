# vrest
Simple REST client inspired by the great [resty](https://github.com/go-resty/resty/releases/tag/v2.7.0) lib.

## Why another rest lib?
We really like the API of resty, but there are a few reasons why we created our own lib:

 1. **No dependencies**: vrest has no dependencies.

 1. **Efficient memory usage**: We don't want a rest lib to copy body bytes. The lib user
    should be in control. vrest gives you access to the body bytes (if available).

 1. **HTTP abstraction**: We want to abstract the HTTP layer away, so there is only
    one error return value, no HTTP response handling required.

 1. **Limited scope**: This lib only covers default rest APIs with JSON and XML.
    We don't plan to add support for handling HTML related stuff like form data
    or multipart file uploads.

 1. **Improved logging/tracing**: Tracing requests was designed with Open Telemetry
    in mind. There is full access to all request/response data.

 1. **No data hiding**: We provide full access to all internal data in callbacks.
    You can break things, if you want, so please look at the examples how to
    perform advanced tasks like overwriting vrest functions.

 1. **Configurability**: Many configurable settings and replaceable functions.
    Most settings can be set in client (default for all requests) and each
    request to override client defaults.

 1. **Easy testing**: It is easy to replace functions in vrest. There are also some
    helper functions to mock responses.

## Installation
Installation

To install vRest, run:

```bash
go get github.com/fond-of-vertigo/vrest
```

## Usage

### Simple GET request
```go
package main_test

import (
    "fmt"
    "github.com/fond-of-vertigo/vrest"
)

func ExampleGETRequest() {
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
```

### Simple POST request
```go
package main_test

import (
	"fmt"
	"net/url"

	"github.com/fond-of-vertigo/vrest"
)

func ExamplePostRequest() {
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
```

### Use Basic Auth
```go
package main_test

import (
	"fmt"

	"github.com/fond-of-vertigo/vrest"
)

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
```

### Set custom header
```go
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
```

### Customize client
```go
package main_test

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
```

### Overriding vrest functions
We're providing a way to override vrest functions. This might be useful for testing or if you want to change the behavior of vrest.
Through the `Overridable` struct in the client, you can replace the functions you want to override.

```go
package main_test

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

```
