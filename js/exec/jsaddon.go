package exec

import "github.com/fastschema/qjs"

type JSAddonObj interface {
	GetName() string
	GetJSObj(ctx *qjs.Context) (*qjs.Value, error)
}
