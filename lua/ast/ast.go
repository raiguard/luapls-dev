package ast

import (
	"fmt"
	"strings"

	"github.com/raiguard/luapls/lua/token"
)

type Node interface {
	String() string
}

type Block struct {
	Node
	Statements []Statement
}

func (b *Block) String() string {
	var out string
	for _, stmt := range b.Statements {
		out += stmt.String() + "\n"
	}
	return strings.TrimSpace(out)
}

type Statement interface {
	Node
	statementNode()
}

type AssignmentStatement struct {
	Token token.Token
	Vars  []Identifier
	Exps  []Expression
}

func (as *AssignmentStatement) statementNode() {}
func (as *AssignmentStatement) String() string {
	return fmt.Sprintf("%s = %s", nodeListToString(as.Vars), nodeListToString(as.Exps))
}

type BreakStatement token.Token

func (bs *BreakStatement) statementNode() {}
func (bs *BreakStatement) String() string {
	return bs.Literal
}

type DoStatement struct {
	Token token.Token
	Body  Block
}

func (ds *DoStatement) statementNode() {}
func (ds *DoStatement) String() string {
	return fmt.Sprintf("%s\n%s\nend", ds.Token.Literal, ds.Body.String())
}

type ForStatement struct {
	Var   Identifier
	Start Expression
	End   Expression
	Step  *Expression // Optional
	Body  Block
}

func (fs *ForStatement) statementNode() {}
func (fs *ForStatement) String() string {
	if fs.Step != nil {
		return fmt.Sprintf(
			"for %s = %s, %s, %s do\n%s\nend",
			fs.Var.String(),
			fs.Start.String(),
			fs.End.String(),
			(*fs.Step).String(),
			fs.Body.String(),
		)
	} else {
		return fmt.Sprintf(
			"for %s = %s, %s do\n%s\nend",
			fs.Var.String(),
			fs.Start.String(),
			fs.End.String(),
			fs.Body.String(),
		)
	}
}

type ForInStatement struct {
	Vars []Identifier
	Exps []Expression
	Body Block
}

func (fs *ForInStatement) statementNode() {}
func (fs *ForInStatement) String() string {
	return fmt.Sprintf("for %s in %s do\n%s\nend", nodeListToString(fs.Vars), nodeListToString(fs.Exps), fs.Body.String())
}

type FunctionStatement struct {
	Name    Identifier
	Params  []Identifier
	Body    Block
	IsLocal bool
}

func (fs *FunctionStatement) statementNode() {}
func (fs *FunctionStatement) String() string {
	localStr := ""
	if fs.IsLocal {
		localStr = "local "
	}
	return fmt.Sprintf(
		"%sfunction %s(%s)\n%s\nend",
		localStr,
		fs.Name.String(),
		nodeListToString(fs.Params),
		fs.Body.String(),
	)
}

type GotoStatement struct {
	Token token.Token
	Label Identifier
}

func (gs *GotoStatement) statementNode() {}
func (gs *GotoStatement) String() string {
	return fmt.Sprintf("%s %s", gs.Token.Literal, gs.Label.String())
}

type IfStatement struct {
	Token   token.Token
	Clauses []IfClause
}

func (is *IfStatement) statementNode() {}
func (is *IfStatement) String() string {
	return fmt.Sprintf("%send", nodeListToString(is.Clauses))
}

type IfClause struct {
	Condition Expression
	Body      Block
}

func (ic IfClause) statementNode() {}
func (ic IfClause) String() string {
	return fmt.Sprintf("if %s then\n%s\n", ic.Condition.String(), ic.Body.String())
}

type LabelStatement struct {
	Token token.Token
	Label Identifier
}

func (ls *LabelStatement) statementNode() {}
func (ls *LabelStatement) String() string {
	return fmt.Sprintf("::%s::", ls.Label.String())
}

type LocalStatement struct {
	Token token.Token
	Names []Identifier
	Exps  []Expression
}

func (ls *LocalStatement) statementNode() {}
func (ls *LocalStatement) String() string {
	return fmt.Sprintf("%s %s = %s", ls.Token.Literal, nodeListToString(ls.Names), nodeListToString(ls.Exps))
}

type RepeatStatement struct {
	Token     token.Token
	Body      Block
	Condition Expression
}

func (rs *RepeatStatement) statementNode() {}
func (rs *RepeatStatement) String() string {
	return fmt.Sprintf("%s\n%s\nuntil %s", rs.Token.Literal, rs.Body.String(), rs.Condition.String())
}

type ReturnStatement struct {
	Exps []Expression
}

func (rs *ReturnStatement) statementNode() {}
func (rs *ReturnStatement) String() string {
	return fmt.Sprintf("return %s", nodeListToString(rs.Exps))
}

type WhileStatement struct {
	Token     token.Token
	Condition Expression
	Body      Block
}

func (ws *WhileStatement) statementNode() {}
func (ws *WhileStatement) String() string {
	return fmt.Sprintf("%s %s do\n%s\nend", ws.Token.Literal, ws.Condition.String(), ws.Body.String())
}

type Expression interface {
	Node
	expressionNode()
}

type BinaryExpression struct {
	Token token.Token
	Left  Expression
	Right Expression
}

func (ie *BinaryExpression) expressionNode() {}
func (ie *BinaryExpression) String() string {
	return fmt.Sprintf("(%s %s %s)", ie.Left.String(), ie.Token.Literal, ie.Right.String())
}

type Identifier token.Token

func (i Identifier) expressionNode() {}
func (i Identifier) String() string  { return i.Literal }

type NumberLiteral struct {
	Token token.Token
	Value float64
}

func (nl *NumberLiteral) expressionNode() {}
func (nl *NumberLiteral) String() string  { return nl.Token.Literal }

type UnaryExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (pe *UnaryExpression) expressionNode() {}
func (pe *UnaryExpression) String() string {
	return fmt.Sprintf("(%s%s)", pe.Operator, pe.Right.String())
}

type StringLiteral struct {
	Token token.Token
	Value string // Without quotes
}

func (sl *StringLiteral) expressionNode() {}
func (sl *StringLiteral) String() string  { return sl.Token.Literal }

func nodeListToString[T Node](nodes []T) string {
	items := []string{}
	for _, node := range nodes {
		items = append(items, node.String())
	}
	return strings.Join(items, ", ")
}
