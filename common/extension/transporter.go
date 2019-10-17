package extension

import "github.com/apache/dubbo-go/remoting"

var (
	transporters = make(map[string]func() remoting.Transporter)
)

func SetTransporter(name string, v func() remoting.Transporter) {
	transporters[name] = v
}

func GetTransporter(name string) remoting.Transporter {
	if transporters[name] == nil {
		panic("transporter for " + name + " is not existing, make sure you have import the package.")
	}
	return transporters[name]()
}
