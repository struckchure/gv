package index

import "github.com/struckchure/gv"

type Page struct{}

func (Page) Rpc() []gv.RpcMethod {
	return []gv.RpcMethod{
		{
			Name: "ListPosts",
			Handler: func(request gv.RpcRequest) gv.RpcResponse {
				return gv.RpcResponse{}
			},
		},
	}
}

func (Page) Get(c any) error {
	return nil
}

func (Page) Post(c any) error {
	return nil
}

func (Page) Patch(c any) error {
	return nil
}

func (Page) Delete(c any) error {
	return nil
}
