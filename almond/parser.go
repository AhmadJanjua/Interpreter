package almond

import (
	"errors"
	"fmt"
)

type Parser struct {
	tokens  []Token
	current int
}

// Constructor for parser
func NewParser(tokens []Token) *Parser {
	thisParser := Parser{tokens, 0}
	return &thisParser
}

// Get current token
func (p *Parser) peek() *Token {
	return &p.tokens[p.current]
}

// Check if End of file is reached
func (p *Parser) isAtEnd() bool {
	return p.peek().GetType() == EOF
}

// Return previous token
func (p *Parser) previous() *Token {
	return &p.tokens[p.current-1]
}

// Return current token and increment
func (p *Parser) advance() *Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

// Report if the current token matches provided token
func (p *Parser) check(tokType TokenType) bool {
	if p.isAtEnd() {
		return false
	}

	return p.peek().GetType() == tokType
}

// Report if the current token matches the list of tokens
func (p *Parser) match(tokTypes ...TokenType) bool {
	for _, tokType := range tokTypes {
		if p.check(tokType) {
			p.advance()
			return true
		}
	}

	return false
}

// Report the error
func (p *Parser) consume(tokType TokenType, message string) (*Token, error) {
	if p.check(tokType) {
		return p.advance(), nil
	}

	TokenError(*p.peek(), message)

	return &Token{}, errors.New("Parser Error: issue occured while parsing token: " + tokType.String())
}

// In the case of an error advance tokens to get fresh start
func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().GetType() == SEMI_COLON {
			return
		}

		switch p.peek().GetType() {
		case CLASS:
		case FN:
		case AUTO:
		case FOR:
		case IF:
		case WHILE:
		case PRINT:
		case RETURN:
			return
		}

		p.advance()
	}
}

// Check primary operations and literals: L7
func (p *Parser) primary() (Expr, error) {
	// Check literals w/o values
	if p.match(FALSE, TRUE, NULL) {
		return NewLiteral(p.previous().GetType()), nil
	}

	// check literals with values
	if p.match(NUMBER, STRING) {
		tok := p.previous()

		if tok.GetType() == NUMBER {
			value, ok := tok.GetLiteral().(float64)

			if !ok {
				return Literal{}, nil
			}

			return NewNumber(value), nil
		}

		return NewString(tok.GetLiteralStr()), nil
	}

	// Get Identifier
	if p.match(IDENTIFIER) {
		return NewVarExpr(*p.previous()), nil
	}

	// check parenthesis
	if p.match(L_PAREN) {
		// find expression
		expression, err := p.expression()

		// cascade error
		if err != nil {
			return nil, err
		}

		_, err = p.consume(R_PAREN, "Expect ')' after expression.")

		// Report error and do not create any expression
		if err != nil {
			return nil, err
		}

		return NewGroupingExpr(expression), nil
	}

	// unknown character
	return NewLiteral(p.peek().GetType()), errors.New("Parser Error: unknown literal in primary")
}

// check for unary operations: L6
func (p *Parser) unary() (Expr, error) {
	if p.match(BANG, MINUS) {
		operator := p.previous()
		right, err := p.unary()

		// cascade error
		if err != nil {
			return nil, err
		}

		return NewUnaryExpr(*operator, right), nil
	}

	return p.primary()
}

// Higher level arithmetic operations: L5
func (p *Parser) factor() (Expr, error) {
	expression, err := p.unary()

	if err != nil {
		return nil, err
	}

	for p.match(SLASH, STAR) {
		operator := p.previous()
		right, err := p.unary()

		if err != nil {
			return nil, err
		}

		expression = NewBinaryExpr(expression, *operator, right)
	}

	return expression, nil
}

// Lower priority arithmetic operations: L4
func (p *Parser) term() (Expr, error) {
	expression, err := p.factor()

	// cascade error
	if err != nil {
		return nil, err
	}

	for p.match(MINUS, PLUS) {
		operator := p.previous()
		right, err := p.factor()

		// cascade error
		if err != nil {
			return nil, err
		}

		expression = NewBinaryExpr(expression, *operator, right)
	}

	return expression, nil
}

// Evaluate comparisons between item: Level3
func (p *Parser) comparison() (Expr, error) {
	expression, err := p.term()

	// cascade error
	if err != nil {
		return nil, err
	}

	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		operator := p.previous()
		right, err := p.term()

		// cascade error
		if err != nil {
			return nil, err
		}

		expression = NewBinaryExpr(expression, *operator, right)
	}

	return expression, nil
}

// Evaluate equality between expressions: Level2
func (p *Parser) equality() (Expr, error) {
	expression, err := p.comparison()

	// cascade error
	if err != nil {
		return nil, err
	}

	for p.match(NOT_EQUALS, EQUALS) {
		operator := p.previous()
		right, err := p.comparison()

		// cascade error
		if err != nil {
			return nil, err
		}

		expression = NewBinaryExpr(expression, *operator, right)
	}

	// no errors
	return expression, nil
}

func (p *Parser) and() (Expr, error) {
	express, err := p.equality()

	if err != nil {
		return nil, err
	}

	for p.match(AND) {
		operator := p.previous()

		right, err := p.equality()

		if err != nil {
			return nil, err
		}

		express = NewLogicalExpr(express, *operator, right)
	}

	return express, nil

}

func (p *Parser) or() (Expr, error) {
	express, err := p.and()

	if err != nil {
		return nil, err
	}

	for p.match(OR) {
		operator := p.previous()

		right, err := p.and()

		if err != nil {
			return nil, err
		}

		express = NewLogicalExpr(express, *operator, right)
	}
	return express, nil
}

func (p *Parser) assignment() (Expr, error) {
	express, err := p.or()

	if err != nil {
		return nil, err
	}

	if p.match(ASSIGNMENT) {
		equals := p.previous()
		value, err := p.assignment()

		if err != nil {
			return nil, err
		}

		s, ok := express.(*VarExpr)

		if !ok {
			TokenError(*equals, "Invalid assignment target.")
			fmt.Printf("Dynamic type: %T\n", express)
			return Literal{}, errors.New("invalid assignment target")
		}

		name := s.GetToken()
		return NewAssignExpr(name, value), nil
	}
	return express, nil
}

// Entry into recursive descent parsing: Level1
func (p *Parser) expression() (Expr, error) {
	return p.assignment()
}

func (p *Parser) expressionStmt() (Stmt, error) {
	value, err := p.expression()

	if err != nil {
		return nil, err
	}

	_, err = p.consume(SEMI_COLON, "Expect ';' after value.")

	if err != nil {
		return nil, err
	}
	return *NewExprStmt(value), nil
}

func (p *Parser) printStmt() (Stmt, error) {
	value, err := p.expression()

	if err != nil {
		return nil, err
	}

	p.consume(SEMI_COLON, "Expect ';' after value.")
	return *NewPrintStmt(value), nil
}

func (p *Parser) block() ([]Stmt, error) {
	var statements []Stmt

	for !p.check(R_BRACE) && !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}

	_, err := p.consume(R_BRACE, "Expected '}' after block.")

	if err != nil {
		return nil, err
	}
	return statements, nil
}

func (p *Parser) ifStmt() (Stmt, error) {
	p.consume(L_PAREN, "Expected '(' after 'if'")

	condition, err := p.expression()

	if err != nil {
		return nil, err
	}

	p.consume(R_PAREN, "Expected ')' after if condition")

	thenBranch, err := p.statement()

	if err != nil {
		return nil, err
	}

	var elseBranch Stmt
	elseBranch = nil

	if p.match(ELSE) {
		elseBranch, err = p.statement()

		if err != nil {
			return nil, err
		}
	}

	return *NewIfStmt(condition, thenBranch, elseBranch), nil
}

func (p *Parser) whileStmt() (Stmt, error) {
	_, err := p.consume(L_PAREN, "Expect '(' after while.")

	if err != nil {
		return nil, err
	}

	condition, err := p.expression()

	if err != nil {
		return nil, err
	}

	_, err = p.consume(R_PAREN, "Expect ')' after condition")

	if err != nil {
		return nil, err
	}

	body, err := p.statement()

	if err != nil {
		return nil, err
	}

	return NewWhileStmt(condition, body), nil

}

func (p *Parser) forStmt() (Stmt, error) {
	_, err := p.consume(L_PAREN, "Expect '(' after 'for'.")

	if err != nil {
		return nil, err
	}

	var initializer Stmt

	if p.match(SEMI_COLON) {
		initializer = nil
	} else if p.match(AUTO) {
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

	var condition Expr

	if !p.check(SEMI_COLON) {
		condition, err = p.expression()

		if err != nil {
			return nil, err
		}
	}
	_, err = p.consume(SEMI_COLON, "Expect ';' after loop condition")

	if err != nil {
		return nil, err
	}

	var increment Expr
	if !p.check(R_PAREN) {
		increment, err = p.expression()

		if err != nil {
			return nil, err
		}
	}

	_, err = p.consume(R_PAREN, "Expect ')' after the for clause.")

	if err != nil {
		return nil, err
	}

	body, err := p.statement()

	if err != nil {
		return nil, err
	}

	if increment != nil {
		body = NewBlockStmt([]Stmt{body, NewExprStmt(increment)})
	}

	if condition == nil {
		condition = NewLiteral(TRUE)
	}

	body = NewWhileStmt(condition, body)

	if initializer != nil {
		body = NewBlockStmt([]Stmt{initializer, body})
	}

	return body, nil
}

func (p *Parser) statement() (Stmt, error) {
	if p.match(FOR) {
		return p.forStmt()
	}
	if p.match(IF) {
		return p.ifStmt()
	}
	if p.match(PRINT) {
		return p.printStmt()
	}

	if p.match(WHILE) {
		return p.whileStmt()
	}

	if p.match(L_BRACE) {
		value, err := p.block()

		if err != nil {
			return nil, err
		}

		return NewBlockStmt(value), nil
	}
	return p.expressionStmt()

}

func (p *Parser) varDeclaration() (Stmt, error) {
	// TODO: handle error
	name, err := p.consume(IDENTIFIER, "Expected variable name.")

	if err != nil {
		return nil, err
	}

	var initializer Expr

	if p.match(ASSIGNMENT) {
		initializer, err = p.expression()

		if err != nil {
			return nil, err
		}
	}

	_, err = p.consume(SEMI_COLON, "Expected ';' after variable declaration.")

	if err != nil {
		return nil, err
	}

	return NewVarStmt(*name, initializer), nil

}

func (p *Parser) declaration() Stmt {
	var statement Stmt
	var err error

	if p.match(AUTO) {
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

func (p *Parser) Parse() []Stmt {
	var statements []Stmt

	for !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}

	return statements
}
