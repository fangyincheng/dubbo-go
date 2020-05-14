module github.com/apache/dubbo-go/registry/kubernetes

replace github.com/apache/dubbo-go v1.4.1 => ../../

replace github.com/apache/dubbo-go/remoting/kubernetes v0.0.0 => ../../remoting/kubernetes

require (
	github.com/apache/dubbo-go v1.4.1
	github.com/apache/dubbo-go/remoting/kubernetes v0.0.0
	github.com/dubbogo/getty v1.3.5
	github.com/dubbogo/gost v1.9.0
	github.com/pkg/errors v0.8.1
	github.com/stretchr/testify v1.5.1
	k8s.io/api v0.0.0-20190325185214-7544f9db76f6
	k8s.io/client-go v8.0.0+incompatible
)

go 1.13
