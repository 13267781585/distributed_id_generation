package handler

import (
	"context"
	"fmt"
	
	server0 "uuid_server/kitex_gen/uuid/generator/server"
	"uuid_server/model"
	"uuid_server/service"
	"uuid_server/utils"
)

type GetUUIDBoundsHandler struct {
	service *service.UUIDService

	req  *server0.GetUUIDBoundsRequest
	resp *server0.GetUUIDBoundsResponse
	ctx  context.Context

	modelBounds []*model.UUIDBound
}

func NewGetUUIDBoundsHandler(c context.Context, req *server0.GetUUIDBoundsRequest, resp *server0.GetUUIDBoundsResponse) *GetUUIDBoundsHandler {
	return &GetUUIDBoundsHandler{
		ctx:     c,
		req:     req,
		resp:    resp,
		service: service.NewUUIDService(),
	}
}

func (h *GetUUIDBoundsHandler) Handle() {
	for _, handler := range []func() error{
		h.getUUIDBounds,
		h.packResp,
	} {
		if err := handler(); err != nil {
			fmt.Printf("GetUUIDBoundsHandler err:%v /n", err)
			h.makeResp(err)
			return
		}
	}
	h.makeResp(nil)
}

func (h *GetUUIDBoundsHandler) makeResp(err error) {
	base := &server0.Base{}
	if err == nil {
		base.Code = utils.Int64Ptr(0)
		base.Message = utils.StringPtr("")
	} else {
		base.Code = utils.Int64Ptr(-1)
		base.Message = utils.StringPtr(err.Error())
	}
	h.resp.Base = base
}

func (h *GetUUIDBoundsHandler) getUUIDBounds() error {
	var err error
	h.modelBounds, err = h.service.GetUUIDBounds(h.req.GetBizCode(), h.req.GetCount())
	if err != nil {
		return fmt.Errorf("GetUUIDBounds err:%v", err)
	}
	return nil
}

func (h *GetUUIDBoundsHandler) packResp() error {
	rpcBounds, err := utils.ConvertModelBoundsToRpc(h.modelBounds)
	if err != nil {
		return err
	}
	h.resp.UuidBounds = rpcBounds
	return nil
}
