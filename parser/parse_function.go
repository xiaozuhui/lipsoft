package parser

import (
	"fmt"
	"lipsoft/ast"
	"lipsoft/token"
	"strconv"
)

type (
	prefixParseFn func() ast.Expression               // 前缀表达式
	infixParseFn  func(ast.Expression) ast.Expression // 中缀表达式
)

// parseIdentifier 与token.IDENT关联的将词法单元和其字面量返回节点的函数
func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

// parseIntegerLiteral 与token.INT关联的解析函数
func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}
	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}
	lit.Value = value
	return lit
}

// parsePrefixExpression 与!和-关联的解析函数
func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}
	p.nextToken() // 需要将词法单元指针前移
	expression.Right = p.parseExpression(PREFIX)
	return expression
}

// parseInfixExpression 解析中缀表达式
func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}
	precedence := p.curPrecedence() // 获取当前词法单元的优先级
	p.nextToken()                   // 当前的
	expression.Right = p.parseExpression(precedence)
	return expression
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()
	exp := p.parseExpression(LOWEST)
	if !p.exceptPeek(token.RPAREN) {
		return nil
	}
	return exp
}
