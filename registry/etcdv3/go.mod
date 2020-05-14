module github.com/apache/dubbo-go/registry/etcdv3

replace github.com/apache/dubbo-go v1.4.1 => ../../

replace github.com/apache/dubbo-go/remoting/etcdv3 v0.0.0 => ../../remoting/etcdv3

replace github.com/coreos/bbolt v1.3.3 => go.etcd.io/bbolt v1.3.3

go 1.13

require (
	github.com/apache/dubbo-go v1.4.1
	github.com/apache/dubbo-go/remoting/etcdv3 v0.0.0
	github.com/dubbogo/getty v1.3.5
	github.com/pkg/errors v0.8.1
	github.com/stretchr/testify v1.5.1
	go.etcd.io/etcd v3.3.13+incompatible
)
