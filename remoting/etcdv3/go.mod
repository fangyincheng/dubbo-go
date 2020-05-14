module github.com/apache/dubbo-go/remoting/etcdv3

replace github.com/apache/dubbo-go v1.4.1 => ../../

replace github.com/coreos/bbolt v1.3.3 => go.etcd.io/bbolt v1.3.3

go 1.13

require (
	github.com/apache/dubbo-go v1.4.1
	github.com/coreos/bbolt v1.3.3 // indirect
	github.com/coreos/etcd v3.3.20+incompatible
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/go-systemd v0.0.0-20190719114852-fd7a80b32e1f // indirect
	github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/dubbogo/getty v1.3.5
	github.com/gogo/protobuf v1.3.1 // indirect
	github.com/golang/groupcache v0.0.0-20200121045136-8c9f03a8e57e // indirect
	github.com/golang/protobuf v1.4.1 // indirect
	github.com/google/btree v1.0.0 // indirect
	github.com/google/uuid v1.1.1 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.0 // indirect
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.14.5 // indirect
	github.com/jonboulle/clockwork v0.1.0 // indirect
	github.com/pkg/errors v0.8.1
	github.com/prometheus/client_golang v1.6.0 // indirect
	github.com/soheilhy/cmux v0.1.4 // indirect
	github.com/stretchr/testify v1.5.1
	github.com/tmc/grpc-websocket-proxy v0.0.0-20200427203606-3cfed13b9966 // indirect
	github.com/xiang90/probing v0.0.0-20190116061207-43a291ad63a2 // indirect
	go.etcd.io/bbolt v1.3.3 // indirect
	go.etcd.io/etcd v3.3.13+incompatible
	golang.org/x/time v0.0.0-20200416051211-89c76fbcd5d1 // indirect
	google.golang.org/genproto v0.0.0-20200513103714-09dca8ec2884 // indirect
	google.golang.org/grpc v1.29.1
	sigs.k8s.io/yaml v1.2.0 // indirect
)
