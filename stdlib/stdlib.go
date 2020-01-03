package stdlib

import(
	"hippo-glox/environment"
	"hippo-glox/function"
)

var sum = (func(env *environment.Env, args []interface{}) interface{} {
	return args[0].(int) + args[1].(int)
})

var Functions = map[string]interface{}{
	"sum":   function.Function{2, nil, sum},
}
