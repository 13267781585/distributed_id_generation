package tools

import (
	"testing"
	"uuid_server/model"
)

func TestHandler(t *testing.T) {
	list := NewLinkedList(EnableLock)
	list.MPush(&model.Bound{1, 100}, &model.Bound{333, 444}, &model.Bound{444, 445})

	count := int64(222)
	bounds := make([]interface{}, 0)
	for count > 0 {
		node := list.GetFirstNode()
		if node == nil {
			return
		}
		val, _ := node.data.(*model.Bound)
		interval := val.End - val.Start
		if count >= interval {
			list.Pop()
			bounds = append(bounds, val)
		} else {
			bounds = append(bounds, &model.Bound{Start: val.Start, End: val.Start + interval})
			val.Start = val.Start + interval
		}

		count = count - interval
	}
	t.Log(bounds)
}
