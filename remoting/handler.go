package remoting

import (
	"github.com/apache/dubbo-go/common"
)

// Handler for remoting
type Handler interface {
	Send(url common.URL, message interface{}) error
	Received(url common.URL, message interface{}) error
}

/////////////////////////////
// base handler
/////////////////////////////

type BaseHandler struct {
}

func (bh *BaseHandler) Send(url common.URL, message interface{}) error {
	return nil
}

func (bh *BaseHandler) Received(url common.URL, message interface{}) error {
	return nil
}
