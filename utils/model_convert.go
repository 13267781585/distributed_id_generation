package utils

import (
	"fmt"

	"uuid_server/kitex_gen/uuid/generator/server"
	"uuid_server/model"
)

func ConvertModelBoundToRpc(modelBound *model.UUIDBound) (*server.UUIDBound, error) {
	if modelBound == nil || modelBound.IntervalBound == nil {
		return nil, fmt.Errorf("bound is nil")
	}
	rpcBound := &server.UUIDBound{}
	rpcBound.Start = &modelBound.IntervalBound.Start
	rpcBound.End = &modelBound.IntervalBound.End
	return rpcBound, nil
}

func ConvertModelBoundsToRpc(modelBounds []*model.UUIDBound) ([]*server.UUIDBound, error) {
	if len(modelBounds) == 0 {
		return []*server.UUIDBound{}, nil
	}

	rpcBounds := make([]*server.UUIDBound, len(modelBounds))
	var err error
	for i, bound := range modelBounds {
		if rpcBounds[i], err = ConvertModelBoundToRpc(bound); err != nil {
			return nil, err
		}
	}
	return rpcBounds, nil
}
