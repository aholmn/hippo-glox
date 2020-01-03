package parser

import (
	"fmt"
	"hippo-glox/ast"
	"hippo-glox/token"
)

var current int
var tokens []token.Token
var statements []ast.Stmt
var pErrors []string

func Parse(arg []token.Token) ([]string, []ast.Stmt) {
	tokens = arg
	for !isAtEnd() {
		statements = append(statements, declaration())
	}
	return pErrors, statements
}

func declaration() ast.Stmt {
	if match(token.Function) {
		return function()
	}
	if match(token.Let) {
		return varDeclaration()
	}
	return statement()
}

func varDeclaration() ast.LetStmt {
	name := eat(token.Identifier, "Expect variable name")
	if match(token.Equal) {
		initializer := expression()
		eat(token.Semicolon, "Expect ';' after variable declaration")
		return ast.LetStmt{initializer, name}
	} else {
		eat(token.Semicolon, "Expect ';' after variable declaration")
		return ast.LetStmt{nil, name}
	}
}

func function() ast.FunctionStmt {
	name := eat(token.Identifier, "Expect function name")
	eat(token.LeftParen, "Expect '(' after function name")
	var params []token.Token

	if !check(token.RightParen) {
		params = append(params, eat(token.Identifier, "Expect parameter name."))
		for match(token.Comma) {
			params = append(params, eat(token.Identifier, "Expect parameter name."))
		}
	}
	eat(token.RightParen, "Expect ')' after parameters.")
	eat(token.LeftBrace, "Expect '{' before function body")
	body := blockStatement()
	return ast.FunctionStmt{name, params, body}
}

func statement() ast.Stmt {
	if match(token.For) {
		return forStatement()
	}
	if match(token.If) {
		return ifStatement()
	}
	if match(token.Print) {
		return printStatement()
	}
	if match(token.Return) {
		return returnStatement()
	}
	if match(token.While) {
		return whileStatement()
	}
	if match(token.LeftBrace) {
		return blockStatement()
	}

	return expressionStatement()
}

func forStatement() ast.Stmt {
	eat(token.LeftParen, "Expect '(' after 'for'.")

	var initializer ast.Stmt
	var condition ast.Expr
	var increment ast.Expr
	var body ast.Stmt

	if match(token.Semicolon) {
		initializer = nil
	} else if match(token.Let) {
		initializer = varDeclaration()
	} else {
		initializer = expressionStatement()
	}
	if !check(token.Semicolon) {
		condition = expression()
	} else {
		condition = nil

	}
	eat(token.Semicolon, "Expect ';' after loop condition")

	if !check(token.RightParen) {
		increment = expression()
	} else {
		increment = nil
	}
	
	eat(token.RightParen, "Expect ')' after for clauses.")
	eat(token.LeftBrace, "Expect '{' after for clauses.")

	body = blockStatement()
	if increment != nil {
		body = ast.BlockStmt{[]ast.Stmt{body, ast.ExprStmt{increment}}}
	}
	if condition == nil {
		condition = ast.Literal{true}
	}

	body = ast.WhileStmt{condition, body}
	if initializer != nil {
		body = ast.BlockStmt{[]ast.Stmt{initializer, body}}
	}

	return body
}

func whileStatement() ast.WhileStmt {
	eat(token.LeftParen, "Expect '(' after 'while'.")
	condition := expression()
	eat(token.RightParen, "Expect ')' after condition.")
	eat(token.LeftBrace, "Expect '{' after parenthesis.")
	body := blockStatement()
	return ast.WhileStmt{condition, body}
}

func ifStatement() ast.IfStmt {
	eat(token.LeftParen, "Expect '(' after 'if'.")
	condition := expression()
	eat(token.RightParen, "Expect ')' after if condition")

	eat(token.LeftBrace, "Expect '{' after if condition")
	thenBranch := blockStatement()

	var elseBranch ast.Stmt
	if match(token.Else) {
		eat(token.LeftBrace, "Expect '{' after if condition")
		elseBranch = blockStatement()
	}
	return ast.IfStmt{condition, thenBranch, elseBranch}
}

func blockStatement() ast.BlockStmt {
	var blocks []ast.Stmt

	for !check(token.RightBrace) && !isAtEnd() {
		blocks = append(blocks, declaration())
	}
	eat(token.RightBrace, "Expect '}' after block.")
	return ast.BlockStmt{blocks}
}

func printStatement() ast.PrintStmt {
	expr := expression()
	eat(token.Semicolon, "Missing semicolon")
	return ast.PrintStmt{Expr: expr}
}

func returnStatement() ast.ReturnStmt {
	keyword := previous()
	var value ast.Expr
	if (!check(token.Semicolon)) {
		value = expression()
	}
	eat(token.Semicolon, "Expect ';' after return value")
	return ast.ReturnStmt{keyword, value}
}

func expressionStatement() ast.ExprStmt {
	expr := expression()
	eat(token.Semicolon, "Missing semicolon")
	return ast.ExprStmt{Expr: expr}
}

func expression() ast.Expr {
	return assignment()
}

func assignment() ast.Expr {
	expr := or()
	if match(token.Equal) {
		equals := previous()
		value := assignment()
		switch v := expr.(type) {
		case ast.Variable:
			name := v.Name
			return ast.Assign{name, value}
		}
		panic(fmt.Sprintf("Invalid assignment target: ", equals))
	}
	return expr
}

func or() ast.Expr {
	expr := and()

	for match(token.Or) {
		operator := previous()
		right := and()
		return ast.Logical{expr, operator, right}
	}
	return expr
}

func and() ast.Expr {
	expr := equality()

	for match(token.And) {
		operator := previous()
		right := equality()
		return ast.Logical{expr, operator, right}
	}
	return expr
}

func equality() ast.Expr {
	expr := comparison()
	for match(token.BangEqual, token.EqualEqual) {
		operator := previous()
		right := comparison()
		return ast.Binary{expr, operator, right}
	}
	return expr
}

func comparison() ast.Expr {
	expr := addition()
	for match(token.Greater, token.GreaterEqual, token.Less, token.LessEqual) {
		operator := previous()
		right := addition()
		return ast.Binary{expr, operator, right}
	}
	return expr
}

func addition() ast.Expr {
	expr := multiplication()
	for match(token.Plus, token.Minus) {
		left := expr
		operator := previous()
		right := multiplication()
		expr = ast.Binary{left, operator, right}
	}
	return expr
}

func multiplication() ast.Expr {
	expr := unary()
	for match(token.Multiplication, token.Division) {
		left := expr
		operator := previous()
		right := multiplication()
		expr = ast.Binary{left, operator, right}
	}
	return expr
}

func unary() ast.Expr {
	for match(token.Minus, token.Bang) {
		operator := previous()
		right := unary()
		return ast.Unary{operator, right}
	}
	return call()
}

func call() ast.Expr {
	expr := primary()
	for true {
		if match(token.LeftParen) {
			expr = finishCall(expr)
			
		} else {
			break;
		}
	}
	return expr
}

func primary() ast.Expr {

	if match(token.Number) {
		return ast.Literal{previous().Value}
	}

	if match(token.String) {
		return ast.Literal{previous().Value}
	}

	if match(token.Boolean) {
		return ast.Literal{previous().Value}
	}

	if match(token.Identifier) {
		return ast.Variable{previous()}
	}

	if match(token.LeftParen) {
		expr := expression()
		eat(token.RightParen, "Expect ')' after expression.")
		return ast.Grouping{expr}
	}
	panic("cannot parse expression")
}

func finishCall(callee ast.Expr) ast.Expr {
	var arguments []ast.Expr
	if !check(token.RightParen) {
		
		arguments = append(arguments, expression())
		for match(token.Comma) {
			arguments = append(arguments, expression())
		}
	}
	paren := eat(token.RightParen, "Expect ')' after arguments.")
	return ast.Call{callee, paren, arguments}
}

func eat(tokenType int, msg string) token.Token {

	if !check(tokenType) {
		msg = fmt.Sprintf("%v at line %v", msg, getLine())
		pErrors = append(pErrors, msg)
		advance()
		return token.Token{Type: tokenType}
	}
	advance()
	return tokens[current-1]
}

func match(types ...int) bool {
	for _, t := range types {
		if !isAtEnd() && peek() == t {
			advance()
			return true
		}
	}
	return false
}

func peek() int {
	return tokens[current].Type
}

func advance() {
	current++
}

func isAtEnd() bool {
	return current >= len(tokens)
}

func previous() token.Token {
	return tokens[current-1]
}

func getLine() int {
	return tokens[current-1].Line
}

func check(t int) bool {
	if isAtEnd() {
		return false
	}
	return peek() == t
}
