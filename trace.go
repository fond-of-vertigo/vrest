package vrest

type TraceMaker interface {
	NewTrace(req *Request) Trace
}

type Trace interface {
	OnAfterRequest(req *Request)
	End()
}
