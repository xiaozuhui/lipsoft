package ast

/*
 * @Author: xiaozuhui
 * @Date: 2023-01-18 17:00:18
 * @LastEditors: xiaozuhui
 * @LastEditTime: 2023-01-18 17:03:56
 * @Description:
 */

type Node interface {
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
