package ast

import (
	"fmt"
	"hippo-glox/environment"
	"hippo-glox/function"
	"hippo-glox/token"
)

type Expr interface {
	Eval(env *environment.Env) interface{}
}

type Stmt interface {
	Execute(env *environment.Env) interface{}
}

type Binary struct {
	Left     Expr
	Operator token.Token
	Right    Expr
}

type Literal struct {
	Value interface{}
}

type Unary struct {
	Operator token.Token
	Right    Expr
}

type Grouping struct {
	Expression Expr
}

type Assign struct {
	Name  token.Token
	Value Expr
}

type Variable struct {
	Name token.Token
}

type Logical struct {
	Left     Expr
	Operator token.Token
	Right    Expr
}

type Call struct {
	Callee    Expr
	Paren     token.Token
	Arguments []Expr
}

type PrintStmt struct {
	Expr Expr
}

type ExprStmt struct {
	Expr Expr
}

type LetStmt struct {
	Initializer Expr
	Name        token.Token
}

type BlockStmt struct {
	Statements []Stmt
}

type IfStmt struct {
	Condition  Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

type WhileStmt struct {
	Condition Expr
	Body      Stmt
}

type FunctionStmt struct {
	Name   token.Token
	Params []token.Token
	Body   BlockStmt
}

type ReturnStmt struct {
	Keyword token.Token
	Value   Expr
}

func (s ReturnStmt) Execute(env *environment.Env) interface{} {
	if s.Value != nil {
		panic(s.Value.Eval(env))
	}
	panic(nil)
}

func (s FunctionStmt) Execute(env *environment.Env) interface{} {
	var fun = (func(env *environment.Env, args []interface{}) interface{} {
		store := make(map[string]interface{})
		inner := environment.Env{Store: store, Outer: env}
		for i, p := range s.Params {
			inner.Store[p.Lexeme] = args[i]
		}
		s.Body.Execute(&inner)
		return nil
	})
	env.Store[s.Name.Lexeme] = function.Function{0, env, fun}
	return nil
}

func (s WhileStmt) Execute(env *environment.Env) interface{} {
	for isTruthy(s.Condition.Eval(env)) {
		s.Body.Execute(env)
	}
	return nil
}

func (s IfStmt) Execute(env *environment.Env) interface{} {
	if isTruthy(s.Condition.Eval(env)) {
		s.ThenBranch.Execute(env)
	} else if s.ElseBranch != nil {
		s.ElseBranch.Execute(env)
	}
	return nil
}

func (s BlockStmt) Execute(env *environment.Env) interface{} {
	store := make(map[string]interface{})
	inner := environment.Env{Store: store, Outer: env}

	for _, i := range s.Statements {
		i.Execute(&inner)

	}
	return nil
}

func (s LetStmt) Execute(env *environment.Env) interface{} {
	if s.Initializer != nil {
		value := s.Initializer.Eval(env)
		env.Store[s.Name.Lexeme] = value
	} else {
		env.Store[s.Name.Lexeme] = nil

	}
	return nil
}

func (s ExprStmt) Execute(env *environment.Env) interface{} {
	return s.Expr.Eval(env)
}

func (s PrintStmt) Execute(env *environment.Env) interface{} {
	fmt.Println(s.Expr.Eval(env))
	return nil
}

func (e Grouping) Eval(env *environment.Env) interface{} {
	return e.Expression.Eval(env)
}

func (e Binary) Eval(env *environment.Env) interface{} {
	left := e.Left.Eval(env)
	right := e.Right.Eval(env)
	switch e.Operator.Type {
	case token.Greater:
		return left.(int) > right.(int)
	case token.GreaterEqual:
		return left.(int) >= right.(int)
	case token.Less:
		return left.(int) < right.(int)
	case token.LessEqual:
		return left.(int) <= right.(int)
	case token.Plus:
		return left.(int) + right.(int)
	case token.Minus:
		return left.(int) - right.(int)
	case token.Multiplication:
		return left.(int) * right.(int)
	case token.Division:
		return left.(int) / right.(int)
	}
	panic(fmt.Sprintln("No match for token: ", e.Operator.Type))
}

func (e Variable) Eval(env *environment.Env) interface{} {
	return get(e.Name.Lexeme, env)
}

func (e Assign) Eval(env *environment.Env) interface{} {
	value := e.Value.Eval(env)
	assign(e.Name.Lexeme, value, env)
	return nil
}

func (e Logical) Eval(env *environment.Env) interface{} {
	left := e.Left.Eval(env)
	if e.Operator.Type == token.Or {
		if isTruthy(left) {
			return left
		}
	} else if e.Operator.Type == token.And {
		if !isTruthy(left) {
			return left
		}
	}
	return e.Right.Eval(env)
}

func (e Call) Eval(env *environment.Env) interface{} {
	callee := e.Callee.Eval(env)
	var args []interface{}
	for _, v := range e.Arguments {
		args = append(args, v.Eval(env))
	}
	switch f := callee.(type) {
	case function.Function:
		return call(f, env, args)
	default:
		return nil
	}
}
func (e Unary) Eval(env *environment.Env) interface{} {
	right := e.Right.Eval(env)

	switch e.Operator.Type {
	case token.Minus:
		return -right.(int)
	case token.Bang:
		return !isTruthy(right)
	}
	panic(fmt.Sprintln("No match for token: ", e.Operator.Type))
}

func (e Literal) Eval(env *environment.Env) interface{} {
	return e.Value
}

func call(f function.Function, env *environment.Env, args []interface{}) (res interface{}) {
	defer func() {
		if err := recover(); err != nil {
			res = err
		}
	}()
	res = f.Call(f.Closure, args)
	return res
}

func get(variable string, env *environment.Env) interface{} {
	if val, ok := env.Store[variable]; ok {
		return val
	}
	if env.Outer != nil {
		return get(variable, env.Outer)
	}
	err := fmt.Errorf("Undefined variable: ", variable)
	panic(err)
}

func assign(variable string, value interface{}, env *environment.Env) {
	if _, ok := env.Store[variable]; ok {
		env.Store[variable] = value
		return
	}
	if env.Outer != nil {
		assign(variable, value, env.Outer)
		return
	}
	panic(fmt.Sprintf("Undefined variable: ", variable))
}

func isTruthy(value interface{}) bool {
	if value == nil {
		return false
	}
	switch v := value.(type) {
	case bool:
		return v
	default:
		return true
	}
}

func isEqual(a interface{}, b interface{}) bool {
	if a == nil && b == nil {
		return true
	} else if a == nil || b == nil {
		return false
	}
	return a == b
}
