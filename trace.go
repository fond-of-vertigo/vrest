package vrest

// TraceMaker defines an interface for handling traces of HTTP requests.
// The interface is designed with Open Telemetry in mind.
// vrest creates a new trace for each request.
// A trace is only created, if the request could be built successfully.
//
// A TraceMaker/Trace has full access to all Request data, but it should
// not modify anything. You can access the BodyBytes of the request and
// response if they are available.
// If the Request uses a Reader, BodyBytes will be nil.
type TraceMaker interface {
	// NewTrace is called just before a request is about to be executed
	// by the HTTP client. NewTrace is NOT called, if a request could not
	// be built, for example because the request body could be marshaled.
	NewTrace(req *Request) Trace
}

// Trace is created for each new request.
type Trace interface {
	// OnAfterRequest is called just after a request was executed.
	// It is called any time, even if the request/response was not
	// successful.
	OnAfterRequest(req *Request)

	// End is called with defer, just after the trace has been created.
	// It can be used to end/close a trace.
	End()
}
