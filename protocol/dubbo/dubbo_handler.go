package dubbo

import (
	"context"
	"reflect"
)

import (
	"github.com/dubbogo/getty"
	"github.com/dubbogo/hessian2"
	perrors "github.com/pkg/errors"
)

import (
	"github.com/dubbo/go-for-apache-dubbo/common"
	"github.com/dubbo/go-for-apache-dubbo/common/constant"
	"github.com/dubbo/go-for-apache-dubbo/common/logger"
	"github.com/dubbo/go-for-apache-dubbo/protocol/invocation"
)

type DubboHandler struct {
	Url common.URL
}

func NewDubboHandler(url common.URL) *DubboHandler {
	return &DubboHandler{Url: url}
}

func (dh *DubboHandler) Received(conn, message interface{}) {
	session, ok := conn.(getty.Session)
	if !ok {
		logger.Error("[DubboHandler] [Received] @conn is not getty.Session.")
		return
	}

	pkg, ok := message.(*DubboPackage)
	if !ok {
		logger.Error("[DubboHandler] [Received] @message is not DubboPackage.")
		return
	}

	err := DecodeBody(pkg)
	if err != nil {
		logger.Errorf("[DubboHandler] [DecodeBody] error: %v", err)
		return
	}
	dh.OnMessage(session, pkg)
}

func (dh *DubboHandler) OnMessage(session getty.Session, p *DubboPackage) {

	export, ok := dubboProtocol.ExporterMap().Load(dh.Url.Key())
	if !ok {
		p.Header.ResponseStatus = hessian.Response_SERVER_ERROR
		p.Body = perrors.New("No exporter!")
		reply(session, p, hessian.PackageResponse)
		return
	}
	invoker := export.(*DubboExporter).GetInvoker()
	if invoker != nil {
		result := invoker.Invoke(invocation.NewRPCInvocationForProvider(p.Service.Method, p.Body.(map[string]interface{})["args"].([]interface{}), map[string]string{
			constant.PATH_KEY: p.Service.Path,
			//attachments[constant.GROUP_KEY] = url.GetParam(constant.GROUP_KEY, "")
			constant.INTERFACE_KEY: p.Service.Interface,
			constant.VERSION_KEY:   p.Service.Version,
		}))
		if err := result.Error(); err != nil {
			p.Header.ResponseStatus = hessian.Response_SERVER_ERROR
			p.Body = err
			reply(session, p, hessian.PackageResponse)
			return
		}
		if res := result.Result(); res != nil {
			p.Header.ResponseStatus = hessian.Response_OK
			p.Body = res
			reply(session, p, hessian.PackageResponse)
			return
		}
	}

	dh.callService(p, nil)

	// not twoway
	if p.Header.Type&hessian.PackageRequest_TwoWay == 0x00 {
		return
	}
	reply(session, p, hessian.PackageResponse)
}

func (dh *DubboHandler) callService(req *DubboPackage, ctx context.Context) {

	defer func() {
		if e := recover(); e != nil {
			req.Header.ResponseStatus = hessian.Response_BAD_REQUEST
			if err, ok := e.(error); ok {
				logger.Errorf("callService panic: %#v", err)
				req.Body = e.(error)
			} else if err, ok := e.(string); ok {
				logger.Errorf("callService panic: %#v", perrors.New(err))
				req.Body = perrors.New(err)
			} else {
				logger.Errorf("callService panic: %#v", e)
				req.Body = e
			}
		}
	}()

	svcIf := req.Body.(map[string]interface{})["service"]
	if svcIf == nil {
		logger.Errorf("service not found!")
		req.Header.ResponseStatus = hessian.Response_SERVICE_NOT_FOUND
		req.Body = perrors.New("service not found")
		return
	}
	svc := svcIf.(*common.Service)
	method := svc.Method()[req.Service.Method]
	if method == nil {
		logger.Errorf("method not found!")
		req.Header.ResponseStatus = hessian.Response_SERVICE_NOT_FOUND
		req.Body = perrors.New("method not found")
		return
	}

	in := []reflect.Value{svc.Rcvr()}
	if method.CtxType() != nil {
		in = append(in, method.SuiteContext(ctx))
	}

	// prepare argv
	argv := req.Body.(map[string]interface{})["args"]
	if (len(method.ArgsType()) == 1 || len(method.ArgsType()) == 2 && method.ReplyType() == nil) && method.ArgsType()[0].String() == "[]interface {}" {
		in = append(in, reflect.ValueOf(argv))
	} else {
		for i := 0; i < len(argv.([]interface{})); i++ {
			in = append(in, reflect.ValueOf(argv.([]interface{})[i]))
		}
	}

	// prepare replyv
	var replyv reflect.Value
	if method.ReplyType() == nil {
		replyv = reflect.New(method.ArgsType()[len(method.ArgsType())-1].Elem())
		in = append(in, replyv)
	}

	returnValues := method.Method().Func.Call(in)

	var retErr interface{}
	if len(returnValues) == 1 {
		retErr = returnValues[0].Interface()
	} else {
		replyv = returnValues[0]
		retErr = returnValues[1].Interface()
	}
	if retErr != nil {
		req.Header.ResponseStatus = hessian.Response_SERVER_ERROR
		req.Body = retErr.(error)
	} else {
		req.Body = replyv.Interface()
	}
}

func reply(session getty.Session, req *DubboPackage, tp hessian.PackageType) {
	resp := &DubboPackage{
		Header: hessian.DubboHeader{
			SerialID:       req.Header.SerialID,
			Type:           tp,
			ID:             req.Header.ID,
			ResponseStatus: req.Header.ResponseStatus,
		},
	}

	if req.Header.Type&hessian.PackageRequest != 0x00 {
		resp.Body = req.Body
	} else {
		resp.Body = nil
	}

	if err := session.WritePkg(Encode(resp), WritePkg_Timeout); err != nil {
		logger.Errorf("WritePkg error: %#v, %#v", perrors.WithStack(err), req.Header)
	}
}
