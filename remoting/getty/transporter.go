package getty

import (
	"github.com/apache/dubbo-go/common"
	"github.com/apache/dubbo-go/common/extension"
	"github.com/apache/dubbo-go/remoting"
)

const (
	GETTY = "getty"
)

func init() {
	extension.SetTransporter(GETTY, GetGettyTransporter)
}

type GettyTransporter struct{}

func (gt *GettyTransporter) Bind(url common.URL, handler remoting.Handler) (remoting.Server, error) {
	return nil, nil
}

func (gt *GettyTransporter) Connect(url common.URL, handler remoting.Handler) (remoting.Client, error) {
	return nil, nil
}

func GetGettyTransporter() remoting.Transporter {
	return &GettyTransporter{}
}
