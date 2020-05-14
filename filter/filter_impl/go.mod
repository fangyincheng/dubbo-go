module github.com/apache/dubbo-go/filter/filter_impl

replace github.com/apache/dubbo-go v1.4.1 => ../../

go 1.13

require (
	github.com/afex/hystrix-go v0.0.0-20180502004556-fa1af6a1f4f5
	github.com/apache/dubbo-go v1.4.1
	github.com/apache/dubbo-go-hessian2 v1.5.0
	github.com/golang/mock v1.3.1
	github.com/mitchellh/mapstructure v1.3.0
	github.com/modern-go/concurrent v0.0.0-20180228061459-e0a39a4cb421
	github.com/opentracing/opentracing-go v1.1.0
	github.com/pkg/errors v0.8.1
	github.com/smartystreets/goconvey v1.6.4 // indirect
	github.com/stretchr/testify v1.5.1
	gopkg.in/yaml.v2 v2.2.2
)
