package simple

import (
	model_pkg "github.com/timbertom-gmbh/go-crud/examples/simple"
	rpc_pkg "github.com/timbertom-gmbh/go-crud/examples/simple"
)

func convertSimpleToModel(rpcObj *rpc_pkg.Simple) *model_pkg.SimpleModel {
	return &model_pkg.SimpleModel{
		Some:  rpcObj.Some,
		Basic: rpcObj.Basic,
		Value: rpcObj.Value,
	}
}
