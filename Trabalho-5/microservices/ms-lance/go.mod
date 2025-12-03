module ms-lance

go 1.25.0

require (
	common v0.0.0-00010101000000-000000000000
	github.com/rabbitmq/amqp091-go v1.10.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	google.golang.org/grpc/cmd/protoc-gen-go-grpc v1.6.0 // indirect
	google.golang.org/protobuf v1.36.10 // indirect
)

replace common => ../common
