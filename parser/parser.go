package parser

import (
	"fmt"
	"lipsoft/ast"
	"lipsoft/lexer"
	"lipsoft/token"
)

type Parser struct {
	lexer *lexer.Lexer // 词法分析器实例

	curToken  token.Token // 当前词法单元
	peekToken token.Token // 下一个词法单元

	errors []string // 错误处理
}

func New(l *lexer.Lexer) *Parser {
	p := Parser{lexer: l, errors: []string{}}

	p.nextToken() // 读取两个词法单元
	p.nextToken()

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
		return nil
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
