package interpreter

import (
	"hippo-glox/ast"
	"hippo-glox/environment"
	"hippo-glox/stdlib"
)

func Interpret(statements []ast.Stmt) {
	env := environment.Env{stdlib.Functions, nil}
	for _, stmt := range statements {
		stmt.Execute(&env)
	}
}
