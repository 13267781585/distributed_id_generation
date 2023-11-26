package handler

import (
	"uuid_server/model"
	"uuid_server/service"
)

type UUIDHandler struct {
	service *service.UUIDService
}

func NewUUIDHandler() *UUIDHandler {
	return &UUIDHandler{
		service: service.NewUUIDService(),
	}
}

func (h *UUIDHandler) GetUUIDBounds(biz int64, count int64) ([]*model.UUIDBound, bool) {
	bounds, err := h.service.GetUUIDBounds(biz, count)
	if err != nil {
		return nil, false
	}
	return bounds, true
}
