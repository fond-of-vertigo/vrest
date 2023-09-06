package vrest

type TraceMaker interface {
	New(req *Request) Trace
}

type Trace interface {
	OnAfterRequest(req *Request)
	End()
}

type NopTraceMaker struct{}
type NopTrace struct{}

var nopTrace *NopTrace = &NopTrace{}

func (*NopTraceMaker) New(req *Request) Trace {
	return nopTrace
}

func (*NopTrace) OnAfterRequest(req *Request) {
}

func (*NopTrace) End() {
}
