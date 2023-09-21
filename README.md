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
