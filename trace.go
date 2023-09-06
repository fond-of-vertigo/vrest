package vrest

import "io"

type Ender interface {
	End()
}

type TraceMaker interface {
	New(req *Request) Trace
}

type Trace interface {
	io.Closer
	OnAfterRequest(req *Request)
}

type NopTraceMaker struct{}

func (*NopTraceMaker) New(req *Request) Trace {
	return &NopTrace{}
}

type NopTrace struct{}

func (*NopTrace) Close() error {
	return nil
}

func (*NopTrace) OnAfterRequest(req *Request) {
}
