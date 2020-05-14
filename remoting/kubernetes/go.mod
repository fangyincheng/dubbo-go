module github.com/apache/dubbo-go/remoting/kubernetes

replace github.com/apache/dubbo-go v1.4.1 => ../../

go 1.13

require (
	github.com/apache/dubbo-go v1.4.1
	github.com/gogo/protobuf v1.3.1 // indirect
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b // indirect
	github.com/google/btree v1.0.0 // indirect
	github.com/gregjones/httpcache v0.0.0-20190611155906-901d90724c79 // indirect
	github.com/peterbourgon/diskv v2.0.1+incompatible // indirect
	github.com/pkg/errors v0.8.1
	github.com/stretchr/testify v1.5.1
	golang.org/x/time v0.0.0-20200416051211-89c76fbcd5d1 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	k8s.io/api v0.0.0-20190325185214-7544f9db76f6
	k8s.io/apimachinery v0.0.0-20190223001710-c182ff3b9841
	k8s.io/client-go v8.0.0+incompatible
	k8s.io/kube-openapi v0.0.0-20191107075043-30be4d16710a // indirect
)
