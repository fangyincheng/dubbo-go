package remoting

import (
	"github.com/apache/dubbo-go/common"
)

type Transporter interface {
	Bind(url common.URL, handler Handler) (Server, error)
	Connect(url common.URL, handler Handler) (Client, error)
}
