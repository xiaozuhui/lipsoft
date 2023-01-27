package parser

import (
	"fmt"
	"lipsoft/ast"
	"lipsoft/lexer"
	"lipsoft/token"
)

type Parser struct {
	lexer *lexer.Lexer // 词法分析器实例

	curToken  token.Token // 当前词法单元指针
	peekToken token.Token // 下一个词法单元指针

	errors []string // 错误处理

	prefixParseFns map[token.TokenType]prefixParseFn // 前缀表达式的映射
	infixParseFns  map[token.TokenType]infixParseFn  // 中缀表达式的映射
}

// registerPrefix 注册前缀表达式
func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

// registerInfix 注册中缀表达式
func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func New(l *lexer.Lexer) *Parser {
	p := Parser{lexer: l, errors: []string{}}

	p.nextToken() // 读取两个词法单元
	p.nextToken()

	// 后缀对应
	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)  // 前缀表达式 !
	p.registerPrefix(token.MINUS, p.parsePrefixExpression) // 前缀表达式 -
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)

	// 中缀对应
	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)     // 加
	p.registerInfix(token.MINUS, p.parseInfixExpression)    // 减
	p.registerInfix(token.SLASH, p.parseInfixExpression)    // 除
	p.registerInfix(token.ASTERISK, p.parseInfixExpression) // 乘
	p.registerInfix(token.EQ, p.parseInfixExpression)       // 等于
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)   // 不等于
	p.registerInfix(token.LT, p.parseInfixExpression)       // 大于
	p.registerInfix(token.GT, p.parseInfixExpression)       // 小于

	return &p
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead",
		t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{} // 构造ast根节点
	program.Statements = []ast.Statement{}

	for !p.curTokenIs(token.EOF) { // 遍历输入的词法单元
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken() // 同时前移curToken和peekToken
	}
	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}
	p.nextToken()

	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}

	// 如果下一个此法单元不是标识符，比如变量，或是字符串字面量，则结束
	if !p.exceptPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	// 如果下一个词法单元不是=
	if !p.exceptPeek(token.ASSIGN) {
		return nil
	}

	// TODO 跳过对表达式的操作，直到遇到分号;
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return t == p.curToken.Type
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return t == p.peekToken.Type
}

// exceptPeek 判断下一个token是否是期望词法单元，如果是，则后移词法单元指针
func (p *Parser) exceptPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t) // 错误
	return false
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)
	// 如果后续词法单元是分号，那么将当前词法单元指针指向分号；分号是可缺省的
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

// parseExpression 解析表达式，使用普拉特解析
func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()
	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() { // 判断是否到分号，并且下一个词法单元优先级是否大于期望优先级
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}
	return leftExp
}
