module github.com/fond-of-vertigo/vrest/examples

go 1.23.0

replace github.com/fond-of-vertigo/vrest => ../../vrest

require (
	github.com/fond-of-vertigo/vrest v0.0.0-00010101000000-000000000000
	go.opentelemetry.io/otel/trace v1.17.0
)

require (
	github.com/go-logr/logr v1.2.4 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	go.opentelemetry.io/otel v1.17.0
	go.opentelemetry.io/otel/metric v1.17.0 // indirect
)
