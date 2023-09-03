package parser

import (
	"fmt"

	"github.com/raiguard/luapls/lua/ast"
	"github.com/raiguard/luapls/lua/token"
)

func (p *Parser) parseStatement() ast.Statement {
	var stat ast.Statement
	switch p.curToken.Type {
	case token.BREAK:
		stat = p.parseBreakStatement()
	case token.GOTO:
		stat = p.parseGotoStatement()
	case token.IDENT:
		stat = p.parseAssignmentStatement()
	case token.IF:
		stat = p.parseIfStatement()
	case token.LOCAL:
		stat = p.parseLocalStatement()
	case token.REPEAT:
		stat = p.parseRepeatStatement()
	case token.WHILE:
		stat = p.parseWhileStatement()
	default:
		p.errors = append(p.errors, "Unexpected <exp>")
		return nil
	}
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stat
}

func (p *Parser) parseAssignmentStatement() *ast.AssignmentStatement {
	stmt := &ast.AssignmentStatement{
		Token: p.curToken,
	}
	stmt.Vars = parseNodeList(p, p.parseIdentifier)
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}
	p.nextToken()
	stmt.Exps = p.parseExpressionList()

	return stmt
}

func (p *Parser) parseBreakStatement() *ast.BreakStatement {
	stmt := ast.BreakStatement(p.curToken)
	return &stmt
}

func (p *Parser) parseGotoStatement() *ast.GotoStatement {
	stmt := ast.GotoStatement{Token: p.curToken}
	if !p.expectPeek(token.IDENT) {
		return nil
	}
	stmt.Label = *p.parseIdentifier()
	return &stmt
}

func (p *Parser) parseIfStatement() *ast.IfStatement {
	stmt := &ast.IfStatement{Token: p.curToken}

	p.nextToken()

	stmt.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.THEN) {
		return nil
	}

	p.nextToken()

	block := p.ParseBlock()

	if block == nil {
		p.errors = append(p.errors, "Failed to parse block")
		return nil
	}

	stmt.Block = *block

	if !p.curTokenIs(token.END) {
		p.errors = append(p.errors, fmt.Sprintf("Expected 'end', got %s", p.curToken.Literal))
		return nil
	}

	return stmt
}

func (p *Parser) parseLocalStatement() *ast.LocalStatement {
	stmt := &ast.LocalStatement{
		Token: p.curToken,
	}
	if !p.expectPeek(token.IDENT) {
		return nil
	}
	stmt.Names = parseNodeList(p, p.parseIdentifier)
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}
	p.nextToken()
	stmt.Exps = p.parseExpressionList()

	return stmt
}

func (p *Parser) parseRepeatStatement() *ast.RepeatStatement {
	stmt := &ast.RepeatStatement{
		Token: p.curToken,
	}
	p.nextToken()
	stmt.Block = *p.ParseBlock()
	if !p.curTokenIs(token.UNTIL) {
		return nil
	}
	p.nextToken()
	stmt.Condition = p.parseExpression(LOWEST)
	return stmt
}

func (p *Parser) parseWhileStatement() *ast.WhileStatement {
	stmt := &ast.WhileStatement{
		Token: p.curToken,
	}
	p.nextToken()
	stmt.Condition = p.parseExpression(LOWEST)
	if !p.expectPeek(token.DO) {
		return nil
	}
	p.nextToken()
	stmt.Block = *p.ParseBlock()
	return stmt
}

func parseNodeList[T ast.Node](p *Parser, parseFunc func() *T) []T {
	values := []T{}
	val := parseFunc()
	if val == nil {
		return values
	}
	values = append(values, *val)

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		val = parseFunc()
		if val == nil {
			break
		}
		values = append(values, *val)
	}

	return values
}
