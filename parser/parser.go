package parser

import (
	"Interpreter/expr"
	"Interpreter/fault"
	"Interpreter/token"
	"Interpreter/tokentype"
	"errors"
	"strconv"
)

type Parser struct {
	tokens  []token.Token
	current int
}

// Constructor for parser
func NewParser(tokens []token.Token) *Parser {
	this_parser := Parser{tokens, 0}
	return &this_parser
}

// Get current token
func (p *Parser) peek() token.Token {
	return p.tokens[p.current]
}

// Check if End of file is reached
func (p *Parser) isAtEnd() bool {
	return p.peek().GetType() == tokentype.EOF
}

// Return previous token
func (p *Parser) previous() token.Token {
	return p.tokens[p.current-1]
}

// Return current token and increment
func (p *Parser) advance() token.Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

// Report if the current token matches provided token
func (p *Parser) check(tok_type tokentype.TokenType) bool {
	if p.isAtEnd() {
		return false
	}

	return p.peek().GetType() == tok_type
}

// Report if the current token matches the list of tokens
func (p *Parser) match(tok_types ...tokentype.TokenType) bool {
	for _, tok_type := range tok_types {
		if p.check(tok_type) {
			p.advance()
			return true
		}
	}

	return false
}

// Report the error
func (p *Parser) consume(tok_type tokentype.TokenType, message string) (token.Token, error) {
	if p.check(tok_type) {
		return p.advance(), nil
	}

	fault.TokenError(p.peek(), message)

	return token.Token{}, errors.New("Parser Error: issue occured while parsing token: " + tok_type.String())
}

// In the case of an error advance tokens to get fresh start
func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().GetType() == tokentype.SEMI_COLON {
			return
		}

		switch p.peek().GetType() {
		case tokentype.CLASS:
		case tokentype.FN:
		case tokentype.AUTO:
		case tokentype.FOR:
		case tokentype.IF:
		case tokentype.WHILE:
		case tokentype.PRINT:
		case tokentype.RETURN:
			return
		}

		p.advance()
	}
}

// Check primary operations and literals: L7
func (p *Parser) primary() (expr.Expr, error) {
	// Check literals w/o values
	if p.match(tokentype.FALSE, tokentype.TRUE, tokentype.NULL) {
		return expr.NewLiteral(p.previous().GetType()), nil
	}

	// check literals with values
	if p.match(tokentype.NUMBER, tokentype.STRING) {
		tok := p.previous()

		if tok.GetType() == tokentype.NUMBER {
			value, err := strconv.ParseFloat(tok.GetLiteral(), 64)

			if err != nil {
				return expr.Literal{}, err
			}

			return expr.NewNumber(value), nil
		}

		return expr.NewString(tok.GetLiteral()), nil
	}

	// check parenthesis
	if p.match(tokentype.L_PAREN) {
		// find expression
		expression, err := p.expression()

		// cascade error
		if err != nil {
			return expression, err
		}

		_, err = p.consume(tokentype.R_PAREN, "Expect ')' after expression.")

		// Report error and do not create any expression
		if err != nil {
			return nil, err
		}

		return expr.NewGrouping(expression), nil
	}

	// unknown character
	return expr.NewLiteral(p.peek().GetType()), errors.New("Parser Error: unknown literal in primary")
}

// check for unary operations: L6
func (p *Parser) unary() (expr.Expr, error) {
	if p.match(tokentype.BANG, tokentype.MINUS) {
		operator := p.previous()
		right, err := p.unary()

		// cascade error
		if err != nil {
			return right, err
		}

		return expr.NewUnary(operator, right), nil
	}

	return p.primary()
}

// Higher level arithmetic operations: L5
func (p *Parser) factor() (expr.Expr, error) {
	expression, err := p.unary()

	if err != nil {
		return expression, err
	}

	for p.match(tokentype.SLASH, tokentype.STAR) {
		operator := p.previous()
		right, err := p.unary()

		if err != nil {
			return right, err
		}

		expression = expr.NewBinary(expression, operator, right)
	}

	return expression, nil
}

// Lower priority arithmetic operations: L4
func (p *Parser) term() (expr.Expr, error) {
	expression, err := p.factor()

	// cascade error
	if err != nil {
		return expression, nil
	}

	for p.match(tokentype.MINUS, tokentype.PLUS) {
		operator := p.previous()
		right, err := p.factor()

		// cascade error
		if err != nil {
			return right, err
		}

		expression = expr.NewBinary(expression, operator, right)
	}

	return expression, nil
}

// Evaluate comparisons between item: Level3
func (p *Parser) comparison() (expr.Expr, error) {
	expression, err := p.term()

	// cascade error
	if err != nil {
		return expression, err
	}

	for p.match(tokentype.GREATER, tokentype.GREATER_EQUAL, tokentype.LESS, tokentype.LESS_EQUAL) {
		operator := p.previous()
		right, err := p.term()

		// cascade error
		if err != nil {
			return right, err
		}

		expression = expr.NewBinary(expression, operator, right)
	}

	return expression, nil
}

// Evaluate equality between expressions: Level2
func (p *Parser) equality() (expr.Expr, error) {
	expression, err := p.comparison()

	// cascade error
	if err != nil {
		return expression, err
	}

	for p.match(tokentype.NOT_EQUALS, tokentype.EQUALS) {
		operator := p.previous()
		right, err := p.comparison()

		// cascade error
		if err != nil {
			return right, err
		}

		expression = expr.NewBinary(expression, operator, right)
	}

	// no errors
	return expression, nil
}

// Entry into recursive descent parsing: Level1
func (p *Parser) expression() (expr.Expr, error) {
	return p.equality()
}

func (p *Parser) Parse() expr.Expr {
	expression, _ := p.expression()
	return expression
}
