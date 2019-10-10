package remoting

type Client interface {
	Send()
}

type BaseClient struct {
}

func (bc *BaseClient) Send() {

}
