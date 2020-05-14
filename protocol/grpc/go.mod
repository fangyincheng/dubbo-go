module github.com/apache/dubbo-go/protocol/grpc

require (
	github.com/apache/dubbo-go v1.4.1
	github.com/apache/dubbo-go-hessian2 v1.5.0
	github.com/golang/protobuf v1.4.0
	github.com/grpc-ecosystem/grpc-opentracing v0.0.0-20180507213350-8e809c8a8645
	github.com/kr/pretty v0.1.0 // indirect
	github.com/opentracing/opentracing-go v1.1.0
	github.com/pkg/errors v0.8.1
	github.com/stretchr/testify v1.5.1
	google.golang.org/grpc v1.22.1
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
)

replace github.com/apache/dubbo-go v1.4.1 => ../../

go 1.13
