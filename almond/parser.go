package almond

import (
	"errors"
	"fmt"
)

// PARSER DESCRIPTION
type Parser struct {
	tokens  []Token
	current int
}

// Ctor
func NewParser(tokens []Token) *Parser {
	thisParser := Parser{tokens, 0}
	return &thisParser
}

// RECURSIVE DESCENT

// ------ Entry
func (p *Parser) Parse() []Stmt {
	var statements []Stmt

	for !p.isAtEnd() {
		statements = append(statements, p.declarationStmt())
	}

	return statements
}

// ------ Statement handling

// get the expression from statment
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

// evaluate block or scoped statement
func (p *Parser) blockStmt() ([]Stmt, error) {
	var statements []Stmt

	for !p.check(R_BRACE) && !p.isAtEnd() {
		statements = append(statements, p.declarationStmt())
	}

	_, err := p.consume(R_BRACE, "Expected '}' after block.")

	if err != nil {
		return nil, err
	}
	return statements, nil
}

// evaluate print statment
func (p *Parser) printStmt() (Stmt, error) {
	value, err := p.expression()

	if err != nil {
		return nil, err
	}

	p.consume(SEMI_COLON, "Expect ';' after value.")
	return *NewPrintStmt(value), nil
}

// evaluate whole statment
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

// evaluate for statement
func (p *Parser) forStmt() (Stmt, error) {
	_, err := p.consume(L_PAREN, "Expect '(' after 'for'.")

	if err != nil {
		return nil, err
	}

	var initializer Stmt

	if p.match(SEMI_COLON) {
		initializer = nil
	} else if p.match(AUTO) {
		initializer, err = p.varStmt()

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

// evaluate if statement
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

// direct to correct statement
func (p *Parser) statement() (Stmt, error) {
	if p.match(IF) {
		return p.ifStmt()
	}
	if p.match(FOR) {
		return p.forStmt()
	}
	if p.match(WHILE) {
		return p.whileStmt()
	}
	if p.match(PRINT) {
		return p.printStmt()
	}

	if p.match(L_BRACE) {
		value, err := p.blockStmt()

		if err != nil {
			return nil, err
		}

		return NewBlockStmt(value), nil
	}
	return p.expressionStmt()

}

// assign function to identifier
func (p *Parser) fnStmt(kind string) (Stmt, error) {
	// Get name
	name, err := p.consume(IDENTIFIER, "Expected "+kind+" name")

	if err != nil {
		return nil, err
	}

	// Filter first parenthesis
	_, err = p.consume(L_PAREN, "Expected '(' after "+kind+" name")

	if err != nil {
		return nil, err
	}

	var params []Token

	// Get arguments
	if !p.check(R_PAREN) {
		for ok := true; ok; ok = p.match(COMMA) {
			if len(params) >= 255 {
				TokenError(*p.peek(), "cannot have more than 255 args")
			}

			param, err := p.consume(IDENTIFIER, "expected parameter name")

			if err != nil {
				return nil, err
			}

			params = append(params, *param)
		}
	}

	_, err = p.consume(R_PAREN, "Expected ')' after args")

	if err != nil {
		return nil, err
	}

	_, err = p.consume(L_BRACE, "Expected '{' before body")

	if err != nil {
		return nil, err
	}

	body, err := p.blockStmt()

	if err != nil {
		return nil, err
	}

	return *NewFnStmt(*name, params, body), nil
}

// assign value to identifier
func (p *Parser) varStmt() (Stmt, error) {
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

// statement parsing entry point + handles error
func (p *Parser) declarationStmt() Stmt {
	var statement Stmt
	var err error

	if p.match(FN) {
		statement, err = p.fnStmt("function")
	} else if p.match(AUTO) {
		statement, err = p.varStmt()
	} else {
		statement, err = p.statement()
	}

	if err != nil {
		p.synchronize()
		return nil
	}

	return statement
}

// ------ Expression handling

// evaluate expression to literal
func (p *Parser) primaryExpr() (Expr, error) {
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

// helper function to deal with calls
func (p *Parser) finishCall(callee Expr) (Expr, error) {
	var arguments []Expr

	if !p.check(R_PAREN) {
		for ok := true; ok; ok = p.match(COMMA) {
			expr, err := p.expression()

			if err != nil {
				return expr, err
			}

			// limit max arguments
			if len(arguments) >= 255 {
				return nil, errors.New("cannot exceed more than 255 arguments")
			}

			arguments = append(arguments, expr)
		}
	}

	paren, err := p.consume(R_PAREN, "expected ')' at the end of a call")

	if err != nil {
		return nil, err
	}

	return NewCallExpr(callee, *paren, arguments), nil
}

// evaluates to a call expression
func (p *Parser) callExpr() (Expr, error) {
	expr, err := p.primaryExpr()

	if err != nil {
		return expr, err
	}

	// check if this is a call
	for {
		if p.match(L_PAREN) {
			// get all the arguments
			expr, err = p.finishCall(expr)

			if err != nil {
				return expr, err
			}
		} else {
			break
		}
	}

	return expr, nil
}

// evaluate to a unary expression
func (p *Parser) unaryExpr() (Expr, error) {
	if p.match(BANG, MINUS) {
		operator := p.previous()
		right, err := p.unaryExpr()

		// cascade error
		if err != nil {
			return nil, err
		}

		return NewUnaryExpr(*operator, right), nil
	}

	return p.callExpr()
}

// evaluate to binary expression
func (p *Parser) factorExpr() (Expr, error) {
	expression, err := p.unaryExpr()

	if err != nil {
		return nil, err
	}

	for p.match(SLASH, STAR) {
		operator := p.previous()
		right, err := p.unaryExpr()

		if err != nil {
			return nil, err
		}

		expression = NewBinaryExpr(expression, *operator, right)
	}

	return expression, nil
}

// evaluate to binary expression
func (p *Parser) termExpr() (Expr, error) {
	expression, err := p.factorExpr()

	// cascade error
	if err != nil {
		return nil, err
	}

	for p.match(MINUS, PLUS) {
		operator := p.previous()
		right, err := p.factorExpr()

		// cascade error
		if err != nil {
			return nil, err
		}

		expression = NewBinaryExpr(expression, *operator, right)
	}

	return expression, nil
}

// evaluate to binary expression
func (p *Parser) comparisonExpr() (Expr, error) {
	expression, err := p.termExpr()

	// cascade error
	if err != nil {
		return nil, err
	}

	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		operator := p.previous()
		right, err := p.termExpr()

		// cascade error
		if err != nil {
			return nil, err
		}

		expression = NewBinaryExpr(expression, *operator, right)
	}

	return expression, nil
}

// evaluate to binary expression
func (p *Parser) equalityExpr() (Expr, error) {
	expression, err := p.comparisonExpr()

	// cascade error
	if err != nil {
		return nil, err
	}

	for p.match(NOT_EQUALS, EQUALS) {
		operator := p.previous()
		right, err := p.comparisonExpr()

		// cascade error
		if err != nil {
			return nil, err
		}

		expression = NewBinaryExpr(expression, *operator, right)
	}

	// no errors
	return expression, nil
}

// evaluate to logical expression
func (p *Parser) andExpr() (Expr, error) {
	express, err := p.equalityExpr()

	if err != nil {
		return nil, err
	}

	for p.match(AND) {
		operator := p.previous()

		right, err := p.equalityExpr()

		if err != nil {
			return nil, err
		}

		express = NewLogicalExpr(express, *operator, right)
	}

	return express, nil

}

// evaluate to logical expression
func (p *Parser) orExpr() (Expr, error) {
	express, err := p.andExpr()

	if err != nil {
		return nil, err
	}

	for p.match(OR) {
		operator := p.previous()

		right, err := p.andExpr()

		if err != nil {
			return nil, err
		}

		express = NewLogicalExpr(express, *operator, right)
	}
	return express, nil
}

// evaluate to assignment expression
func (p *Parser) assignmentExpr() (Expr, error) {
	express, err := p.orExpr()

	if err != nil {
		return nil, err
	}

	if p.match(ASSIGNMENT) {
		equals := p.previous()
		value, err := p.assignmentExpr()

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

// expression evaluation entry
func (p *Parser) expression() (Expr, error) {
	return p.assignmentExpr()
}

// HELPER FUNCTIONS

// ----- Bool helpers

// checks if last token is reached
func (p *Parser) isAtEnd() bool {
	return p.peek().GetType() == EOF
}

// match current token with supplied token
func (p *Parser) check(tokType TokenType) bool {
	if p.isAtEnd() {
		return false
	}

	return p.peek().GetType() == tokType
}

// match and advance tokens
func (p *Parser) match(tokTypes ...TokenType) bool {
	for _, tokType := range tokTypes {
		if p.check(tokType) {
			p.advance()
			return true
		}
	}

	return false
}

// ----- Token helpers

// check current token
func (p *Parser) peek() *Token {
	return &p.tokens[p.current]
}

// check previous token
func (p *Parser) previous() *Token {
	return &p.tokens[p.current-1]
}

// return current token and go to next
func (p *Parser) advance() *Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

// ----- Token + Error helpers

// return current token, advance and report error
func (p *Parser) consume(tokType TokenType, message string) (*Token, error) {
	if p.check(tokType) {
		return p.advance(), nil
	}

	TokenError(*p.peek(), message)

	return &Token{}, errors.New("Parser Error: issue occured while parsing token: " + tokType.String())
}

// ----- Error helpers

// advance to next non-erroneous state
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
