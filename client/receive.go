package client

import "context"

// 对来自客户端的 message(request、response、notification)进行接收处理
// 对 request、notification 路由到对应的handler，对 response 传入 request 的 chan

func (client *Client) Receive(ctx context.Context, msg []byte) error {

}
