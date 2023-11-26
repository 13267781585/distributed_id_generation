package main

import "uuid_server/logic"

func main() {
	logic.InitPool()
	defer logic.StopPool()
}
