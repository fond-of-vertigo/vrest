package vrest

import "io"

type TraceMaker interface {
	New(req *Request) Trace
}

type Trace interface {
	io.Closer
	OnAfterRequest(req *Request)
}

type NopTraceMaker struct{}
type NopTrace struct{}

func (*NopTraceMaker) New(req *Request) Trace {
	return &NopTrace{}
}

func (*NopTrace) Close() error {
	return nil
}
func (*NopTrace) OnAfterRequest(req *Request) {
}
