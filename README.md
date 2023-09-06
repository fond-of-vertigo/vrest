# vrest
Simple REST Client

## Design Goals 
    * No dependencies
    * Simplicity
        * Easy interface with chainable methods
        * Only one error as result
            * If err = nil, the request was successful
            * No need to inspect a response struct
    * Efficiency
        * Don't create copies of request and response bodies
    * Tracing
        * Full support for tracing requests and responses
        * Bodies can be traced if available
    * Configurability
        * Set defaults in the client instance
        * Override all defaults in the request instance
        * Many overridable functions 
    * Testability
        * Super easy to mock request execution
    * Limited scope
        * Easy handling for JSON and XML APIs
        * No support planned for non-rest features like Forms or Multipart
        * No redirect handling
        * No retry handling, yet
        
