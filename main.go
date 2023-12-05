package main

import (
	"fmt"
	"net"
	"uuid_server/kitex_gen/uuid/generator/server/uuidgeneratorserver"
	"uuid_server/logic"

	"github.com/cloudwego/kitex/server"
)

func main() {
	logic.InitPool()
	defer logic.StopPool()

	svr := uuidgeneratorserver.NewServer(&UUIDGeneratorServerImpl{}, server.WithServiceAddr(&net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 8888}))
	if err := svr.Run(); err != nil {
		fmt.Printf("server err:%v", err)
	}
}
