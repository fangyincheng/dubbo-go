module github.com/apache/dubbo-go/registry/consul

require (
	github.com/apache/dubbo-go v1.4.1
	github.com/dubbogo/gost v1.9.0
	github.com/hashicorp/consul v1.5.3
	github.com/hashicorp/consul/api v1.1.0
	github.com/pkg/errors v0.8.1
	github.com/stretchr/testify v1.5.1
)

replace github.com/apache/dubbo-go v1.4.1 => ../../

go 1.13
