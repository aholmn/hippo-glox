package scanner

import (
	"hippo-glox/token"
	"strconv"
)

var source string
var current int
var bytes int
var tokens []token.Token
var line int = 1

func Scan(arg string) []token.Token {
	bytes = len(arg)
	source = arg
	for current < bytes {
		findToken()
	}
	return tokens
}

func findToken() {
	var lexeme string
	var c = source[current]
	pos := current

	switch c {
	case ' ':
	case '+':
		add(token.Token{token.Plus, nil, "+", line})
	case '-':
		add(token.Token{token.Minus, nil, "-", line})
	case '*':
		add(token.Token{token.Multiplication, nil, "*", line})
	case ';':
		add(token.Token{token.Semicolon, nil, ";", line})
	case '{':
		add(token.Token{token.LeftBrace, nil, "{", line})
	case '}':
		add(token.Token{token.RightBrace, nil, "}", line})
	case '(':
		add(token.Token{token.LeftParen, nil, "(", line})
	case ')':
		add(token.Token{token.RightParen, nil, ")", line})
	case ',':
		add(token.Token{token.Comma, nil, ",", line})
	case '\n':
		line++
	case '!':
		add(token.Token{token.Bang, nil, "!", line})
	case '=':
		add(token.Token{token.Equal, nil, "=", line})
	case '>':
		if !isAtEnd() && (current+1 < bytes && next() == '=') {
			add(token.Token{token.GreaterEqual, nil, ">", line})
			advance()
		} else {
			add(token.Token{token.Greater, nil, ">", line})
		}
		advance()
	case '<':
		if !isAtEnd() && (current+1 < bytes && next() == '=') {
			add(token.Token{token.LessEqual, nil, ">", line})
			advance()
		} else {
			add(token.Token{token.Less, nil, "<", line})
		}
		advance()
	case '"':
		for !isAtEnd() && (current+1 < bytes && next() != '"') {
			advance()
		}
		value := source[pos+1 : current+1]
		add(token.Token{token.String, value, value, line})
		advance()
	default:
		if isNumber(c) {
			for !isAtEnd() && ((current+1) < bytes && isNumber(next())) {
				advance()
			}
			lexeme = source[pos : current+1]
			value, _ := strconv.Atoi(lexeme)
			add(token.Token{token.Number, value, lexeme, line})
		} else if isAlpha(c) {
			for !isAtEnd() && ((current+1) < bytes && isAlphaNumeric(next())) {
				advance()
			}
			lexeme = source[pos : current+1]
			if val, ok := keywords[lexeme]; ok {
				if lexeme == "true" {
					add(token.Token{val, true, lexeme, line})
				} else if lexeme == "false" {
					add(token.Token{val, false, lexeme, line})
				} else {
					add(token.Token{val, lexeme, lexeme, line})
				}
			} else {
				add(token.Token{token.Identifier, lexeme, lexeme, line})
			}
		}

	}
	advance()
}

func advance() {
	current++
}

func add(token token.Token) {
	tokens = append(tokens, token)
}

func isAtEnd() bool {
	return current > bytes
}

func next() byte {
	return source[current+1]
}

func isNumber(c byte) bool {
	_, err := strconv.Atoi(string(c))
	return err == nil
}

func isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

func isAlphaNumeric(c byte) bool {
	return isNumber(c) || isAlpha(c)
}

var keywords = map[string]int{
	"print":    token.Print,
	"let":      token.Let,
	"if":       token.If,
	"else":     token.Else,
	"or":       token.Or,
	"and":      token.And,
	"true":     token.Boolean,
	"false":    token.Boolean,
	"while":    token.While,
	"for":      token.For,
	"function": token.Function,
	"return":   token.Return,
}
