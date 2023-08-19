/**
User: cr-mao
Date: 2023/8/18 15:21
Email: crmao@qq.com
Desc: node_client.go
*/
package gate

// NewNodeClient 新建节点客户端
func (t *Transporter) NewNodeClient(ep *endpoint.Endpoint) (transport.NodeClient, error) {
	t.once.Do(func() {
		t.builder = client.NewBuilder(&t.opts.client)
	})

	cc, err := t.builder.Build(ep.Target())
	if err != nil {
		return nil, err
	}

	return node.NewClient(cc), nil
}
