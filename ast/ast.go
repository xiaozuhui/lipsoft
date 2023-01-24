package ast

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
