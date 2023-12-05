package main

import (
	"context"
	"uuid_server/handler"
	server0 "uuid_server/kitex_gen/uuid/generator/server"
)

// UUIDGeneratorServerImpl implements the last service interface defined in the IDL.
type UUIDGeneratorServerImpl struct{}

// GetUUIDBounds implements the UUIDGeneratorServerImpl interface.
func (s *UUIDGeneratorServerImpl) GetUUIDBounds(ctx context.Context, req *server0.GetUUIDBoundsRequest) (resp *server0.GetUUIDBoundsResponse, err error) {
	resp = &server0.GetUUIDBoundsResponse{}
	handler.NewGetUUIDBoundsHandler(ctx, req, resp).Handle()
	return
}
