module github.com/apache/dubbo-go/metrics/prometheus

replace github.com/apache/dubbo-go v1.4.1 => ../../

go 1.13

require (
	github.com/apache/dubbo-go v1.4.1
	github.com/prometheus/client_golang v1.1.0
	github.com/stretchr/testify v1.5.1
)
