package vrest

type TraceMaker interface {
	NewTrace(req *Request) Trace
}

type Trace interface {
	OnAfterRequest(req *Request)
	End()
}

type NopTraceMaker struct{}
type NopTrace struct{}

var nopTrace *NopTrace = &NopTrace{}

func (*NopTraceMaker) NewTrace(req *Request) Trace {
	return nopTrace
}

func (*NopTrace) OnAfterRequest(req *Request) {
}

func (*NopTrace) End() {
}
