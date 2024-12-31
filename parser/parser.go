package parser

import (
	"Interpreter/expr"
	"Interpreter/fault"
	"Interpreter/stmt"
	"Interpreter/token"
	"Interpreter/tokentype"
	"errors"
	"fmt"
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
func (p *Parser) peek() *token.Token {
	return &p.tokens[p.current]
}

// Check if End of file is reached
func (p *Parser) isAtEnd() bool {
	return p.peek().GetType() == tokentype.EOF
}

// Return previous token
func (p *Parser) previous() *token.Token {
	return &p.tokens[p.current-1]
}

// Return current token and increment
func (p *Parser) advance() *token.Token {
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
func (p *Parser) consume(tok_type tokentype.TokenType, message string) (*token.Token, error) {
	if p.check(tok_type) {
		return p.advance(), nil
	}

	fault.TokenError(*p.peek(), message)

	return &token.Token{}, errors.New("Parser Error: issue occured while parsing token: " + tok_type.String())
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
			value, ok := tok.GetLiteral().(float64)

			if !ok {
				return expr.Literal{}, nil
			}

			return expr.NewNumber(value), nil
		}

		return expr.NewString(tok.GetLiteralStr()), nil
	}

	// Get Identifier
	if p.match(tokentype.IDENTIFIER) {
		return expr.NewVar(*p.previous()), nil
	}

	// check parenthesis
	if p.match(tokentype.L_PAREN) {
		// find expression
		expression, err := p.expression()

		// cascade error
		if err != nil {
			return nil, err
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
			return nil, err
		}

		return expr.NewUnary(*operator, right), nil
	}

	return p.primary()
}

// Higher level arithmetic operations: L5
func (p *Parser) factor() (expr.Expr, error) {
	expression, err := p.unary()

	if err != nil {
		return nil, err
	}

	for p.match(tokentype.SLASH, tokentype.STAR) {
		operator := p.previous()
		right, err := p.unary()

		if err != nil {
			return nil, err
		}

		expression = expr.NewBinary(expression, *operator, right)
	}

	return expression, nil
}

// Lower priority arithmetic operations: L4
func (p *Parser) term() (expr.Expr, error) {
	expression, err := p.factor()

	// cascade error
	if err != nil {
		return nil, err
	}

	for p.match(tokentype.MINUS, tokentype.PLUS) {
		operator := p.previous()
		right, err := p.factor()

		// cascade error
		if err != nil {
			return nil, err
		}

		expression = expr.NewBinary(expression, *operator, right)
	}

	return expression, nil
}

// Evaluate comparisons between item: Level3
func (p *Parser) comparison() (expr.Expr, error) {
	expression, err := p.term()

	// cascade error
	if err != nil {
		return nil, err
	}

	for p.match(tokentype.GREATER, tokentype.GREATER_EQUAL, tokentype.LESS, tokentype.LESS_EQUAL) {
		operator := p.previous()
		right, err := p.term()

		// cascade error
		if err != nil {
			return nil, err
		}

		expression = expr.NewBinary(expression, *operator, right)
	}

	return expression, nil
}

// Evaluate equality between expressions: Level2
func (p *Parser) equality() (expr.Expr, error) {
	expression, err := p.comparison()

	// cascade error
	if err != nil {
		return nil, err
	}

	for p.match(tokentype.NOT_EQUALS, tokentype.EQUALS) {
		operator := p.previous()
		right, err := p.comparison()

		// cascade error
		if err != nil {
			return nil, err
		}

		expression = expr.NewBinary(expression, *operator, right)
	}

	// no errors
	return expression, nil
}

func (p *Parser) and() (expr.Expr, error) {
	express, err := p.equality()

	if err != nil {
		return nil, err
	}

	for p.match(tokentype.AND) {
		operator := p.previous()

		right, err := p.equality()

		if err != nil {
			return nil, err
		}

		express = expr.NewLogical(express, *operator, right)
	}

	return express, nil

}

func (p *Parser) or() (expr.Expr, error) {
	express, err := p.and()

	if err != nil {
		return nil, err
	}

	for p.match(tokentype.OR) {
		operator := p.previous()

		right, err := p.and()

		if err != nil {
			return nil, err
		}

		express = expr.NewLogical(express, *operator, right)
	}
	return express, nil
}

func (p *Parser) assignment() (expr.Expr, error) {
	express, err := p.or()

	if err != nil {
		return nil, err
	}

	if p.match(tokentype.ASSIGNMENT) {
		equals := p.previous()
		value, err := p.assignment()

		if err != nil {
			return nil, err
		}

		s, ok := express.(*expr.Var)

		if !ok {
			fault.TokenError(*equals, "Invalid assignment target.")
			fmt.Printf("Dynamic type: %T\n", express)
			return expr.Literal{}, errors.New("invalid assignment target")
		}

		name := s.GetToken()
		return expr.NewAssign(name, value), nil
	}
	return express, nil
}

// Entry into recursive descent parsing: Level1
func (p *Parser) expression() (expr.Expr, error) {
	return p.assignment()
}

func (p *Parser) expressionStmt() (stmt.Stmt, error) {
	value, err := p.expression()

	if err != nil {
		return nil, err
	}

	_, err = p.consume(tokentype.SEMI_COLON, "Expect ';' after value.")

	if err != nil {
		return nil, err
	}
	return *stmt.NewExpression(value), nil
}

func (p *Parser) printStmt() (stmt.Stmt, error) {
	value, err := p.expression()

	if err != nil {
		return nil, err
	}

	p.consume(tokentype.SEMI_COLON, "Expect ';' after value.")
	return *stmt.NewPrint(value), nil
}

func (p *Parser) block() ([]stmt.Stmt, error) {
	var statements []stmt.Stmt

	for !p.check(tokentype.R_BRACE) && !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}

	_, err := p.consume(tokentype.R_BRACE, "Expected '}' after block.")

	if err != nil {
		return nil, err
	}
	return statements, nil
}

func (p *Parser) ifStmt() (stmt.Stmt, error) {
	p.consume(tokentype.L_PAREN, "Expected '(' after 'if'")

	condition, err := p.expression()

	if err != nil {
		return nil, err
	}

	p.consume(tokentype.R_PAREN, "Expected ')' after if condition")

	then_branch, err := p.statement()

	if err != nil {
		return nil, err
	}

	var else_branch stmt.Stmt
	else_branch = nil

	if p.match(tokentype.ELSE) {
		else_branch, err = p.statement()

		if err != nil {
			return nil, err
		}
	}

	return *stmt.NewIf(condition, then_branch, else_branch), nil
}

func (p *Parser) whileStmt() (stmt.Stmt, error) {
	_, err := p.consume(tokentype.L_PAREN, "Expect '(' after while.")

	if err != nil {
		return nil, err
	}

	condition, err := p.expression()

	if err != nil {
		return nil, err
	}

	_, err = p.consume(tokentype.R_PAREN, "Expect ')' after condition")

	if err != nil {
		return nil, err
	}

	body, err := p.statement()

	if err != nil {
		return nil, err
	}

	return stmt.NewWhile(condition, body), nil

}

func (p *Parser) forStmt() (stmt.Stmt, error) {
	_, err := p.consume(tokentype.L_PAREN, "Expect '(' after 'for'.")

	if err != nil {
		return nil, err
	}

	var initializer stmt.Stmt

	if p.match(tokentype.SEMI_COLON) {
		initializer = nil
	} else if p.match(tokentype.AUTO) {
		initializer, err = p.varDeclaration()

		if err != nil {
			return nil, err
		}
	} else {
		initializer, err = p.expressionStmt()

		if err != nil {
			return nil, err
		}
	}

	var condition expr.Expr

	if !p.check(tokentype.SEMI_COLON) {
		condition, err = p.expression()

		if err != nil {
			return nil, err
		}
	}
	_, err = p.consume(tokentype.SEMI_COLON, "Expect ';' after loop condition")

	if err != nil {
		return nil, err
	}

	var increment expr.Expr
	if !p.check(tokentype.R_PAREN) {
		increment, err = p.expression()

		if err != nil {
			return nil, err
		}
	}

	_, err = p.consume(tokentype.R_PAREN, "Expect ')' after the for clause.")

	if err != nil {
		return nil, err
	}

	body, err := p.statement()

	if err != nil {
		return nil, err
	}

	if increment != nil {
		body = stmt.NewBlock([]stmt.Stmt{body, stmt.NewExpression(increment)})
	}

	if condition == nil {
		condition = expr.NewLiteral(tokentype.TRUE)
	}

	body = stmt.NewWhile(condition, body)

	if initializer != nil {
		body = stmt.NewBlock([]stmt.Stmt{initializer, body})
	}

	return body, nil
}

func (p *Parser) statement() (stmt.Stmt, error) {
	if p.match(tokentype.FOR) {
		return p.forStmt()
	}
	if p.match(tokentype.IF) {
		return p.ifStmt()
	}
	if p.match(tokentype.PRINT) {
		return p.printStmt()
	}

	if p.match(tokentype.WHILE) {
		return p.whileStmt()
	}

	if p.match(tokentype.L_BRACE) {
		value, err := p.block()

		if err != nil {
			return nil, err
		}

		return stmt.NewBlock(value), nil
	}
	return p.expressionStmt()

}

func (p *Parser) varDeclaration() (stmt.Stmt, error) {
	// TODO: handle error
	name, err := p.consume(tokentype.IDENTIFIER, "Expected variable name.")

	if err != nil {
		return nil, err
	}

	var initializer expr.Expr

	if p.match(tokentype.ASSIGNMENT) {
		initializer, err = p.expression()

		if err != nil {
			return nil, err
		}
	}

	_, err = p.consume(tokentype.SEMI_COLON, "Expected ';' after variable declaration.")

	if err != nil {
		return nil, err
	}

	return stmt.NewVar(*name, initializer), nil

}

func (p *Parser) declaration() stmt.Stmt {
	var statement stmt.Stmt
	var err error

	if p.match(tokentype.AUTO) {
		statement, err = p.varDeclaration()
	} else {
		statement, err = p.statement()
	}

	if err != nil {
		p.synchronize()
		return nil
	}

	return statement
}

func (p *Parser) Parse() []stmt.Stmt {
	var statements []stmt.Stmt

	for !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}

	return statements
}
