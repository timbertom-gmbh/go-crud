package simple

import "testing"

func TestConverter(t *testing.T) {
	rpcObj := &Simple{
		Some:  1,
		More:  56,
		Basic: "test",
		Value: "test2",
	}

	model := convertSimpleToModel(rpcObj)
	if model.Basic != rpcObj.Basic {
		t.Error("Basic field doesnt match")
	}
	if model.EvenMore != int32(rpcObj.More) {
		t.Error("Basic field doesnt match")
	}
}
