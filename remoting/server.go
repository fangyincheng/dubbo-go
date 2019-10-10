package remoting

type Server interface {
	Start()
	Stop()
}

type BaseServer struct {
}

func (bs *BaseServer) Send() {

}
