package function

import (
	"hippo-glox/environment"
)

type Function struct {
	Arity   int
	Closure *environment.Env
	Call    func(env *environment.Env, args []interface{}) interface{}
}
