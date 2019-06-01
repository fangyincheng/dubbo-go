package protocol

import (
	"github.com/dubbo/go-for-apache-dubbo/common"
)

type Handler interface {
	Received(interface{}, interface{})
}

// Base Handler
type BaseHandler struct {
	H   Handler
	Url common.URL
}

func NewBaseHandler(handler Handler, url common.URL) *BaseHandler {
	return &BaseHandler{H: handler, Url: url}
}

func (bh *BaseHandler) Received(conn, message interface{}) {
}
