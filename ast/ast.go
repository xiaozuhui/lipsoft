package ast

import "lipsoft/token"

/*
 * @Author: xiaozuhui
 * @Date: 2023-01-18 17:00:18
 * @LastEditors: xiaozuhui
 * @LastEditTime: 2023-01-18 17:03:56
 * @Description: let x = 5;
 */

type Node interface {
	// TokenLiteral 返回关联的词法单元的字面量
	TokenLiteral() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

// Program 每个ast节点的根节点
type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

// LetStatement 标识符、产生值的表达式、词法单元
type LetStatement struct {
	Token token.Token // token.LET词法单元
	Name  *Identifier // 绑定的标识符
	Value Expression  // 产生值的表达式
}

func (ls *LetStatement) statementNode() {}

func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}

// Identifier 即`let x = 5;`中的x
type Identifier struct {
	Token token.Token // token.IDENT词法单元
	Value string
}

func (i *Identifier) expressionNode() {}

func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}
