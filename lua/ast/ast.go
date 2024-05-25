package ast

import (
	"github.com/raiguard/luapls/lua/token"
)

type Node interface {
	Pos() token.Pos
	End() token.Pos
}

func Range(n Node) token.Range {
	return token.Range{Start: n.Pos(), End: n.End()}
}

type Block struct {
	Stmts    []Statement
	StartPos token.Pos `json:"-"`
	EndPos   token.Pos `json:"-"`
}

func (b *Block) Pos() token.Pos {
	return b.StartPos
}
func (b *Block) End() token.Pos {
	return b.EndPos
}

type TableField struct {
	Key      Expression
	Value    Expression
	StartPos token.Pos `json:"-"` // Needed in case of bracketed keys
}

func (tf *TableField) Pos() token.Pos {
	return tf.StartPos
}
func (tf *TableField) End() token.Pos {
	return tf.Value.End()
}

type Invalid struct {
	Exps []Expression `json:",omitempty"`
	// OR
	Token *token.Token `json:",omitempty"`
}

func (i *Invalid) expressionNode() {}
func (i *Invalid) statementNode()  {}
func (i *Invalid) Pos() token.Pos {
	if i.Token != nil {
		return i.Token.Pos
	}
	return i.Exps[0].Pos()
}
func (i *Invalid) End() token.Pos {
	if i.Token != nil {
		return i.Token.End()
	}
	return i.Exps[len(i.Exps)-1].End()
}

// A Leaf node has no children and is interactable in the editor.
type LeafNode interface {
	leaf()
}
