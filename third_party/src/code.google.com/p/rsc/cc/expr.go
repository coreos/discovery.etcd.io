// Copyright 2013 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cc

import "fmt"

// An Expr is a parsed C expression.
type Expr struct {
	SyntaxInfo
	Op    ExprOp   // operator
	Left  *Expr    // left (or only) operand
	Right *Expr    // right operand
	List  []*Expr  // operand list, for Comma, Cond, Call
	Text  string   // name or literal, for Name, Number, Goto, Arrow, Dot
	Texts []string // list of literals, for String
	Type  *Type    // type operand, for SizeofType, Offsetof, Cast, CastInit, VaArg
	Init  *Init    // initializer, for CastInit

	// derived information
	XDecl *Decl
	XType *Type // expression type, derived
}

func (x *Expr) String() string {
	var p Printer
	p.printExpr(x, precLow)
	return p.String()
}

type ExprOp int

const (
	_          ExprOp = iota
	Add               // Left + Right
	AddEq             // Left += Right
	Addr              // &Left
	And               // Left & Right
	AndAnd            // Left && Right
	AndEq             // Left &= Right
	Arrow             // Left->Text
	Call              // Left(List)
	Cast              // (Type)Left
	CastInit          // (Type){Init}
	Comma             // x, y, z; List = {x, y, z}
	Cond              // x ? y : z; List = {x, y, z}
	Div               // Left / Right
	DivEq             // Left /= Right
	Dot               // Left.Name
	Eq                // Left = Right
	EqEq              // Left == Right
	Gt                // Left > Right
	GtEq              // Left >= Right
	Index             // Left[Right]
	Indir             // *Left
	Lsh               // Left << Right
	LshEq             // Left <<= Right
	Lt                // Left < Right
	LtEq              // Left <= Right
	Minus             // -Left
	Mod               // Left % Right
	ModEq             // Left %= Right
	Mul               // Left * Right
	MulEq             // Left *= Right
	Name              // Text (function, variable, or enum name)
	Not               // !Left
	NotEq             // Left != Right
	Number            // Text (numeric or chraracter constant)
	Offsetof          // offsetof(Type, Left)
	Or                // Left | Right
	OrEq              // Left |= Right
	OrOr              // Left || Right
	Paren             // (Left)
	Plus              //  +Left
	PostDec           // Left--
	PostInc           // Left++
	PreDec            // --Left
	PreInc            // ++Left
	Rsh               // Left >> Right
	RshEq             // Left >>= Right
	SizeofExpr        // sizeof(Left)
	SizeofType        // sizeof(Type)
	String            // Text (quoted string literal)
	Sub               // Left - Right
	SubEq             // Left -= Right
	Twid              // ~Left
	VaArg             // va_arg(Left, Type)
	Xor               // Left ^ Right
	XorEq             // Left ^= Right
)

var exprOpString = []string{
	Add:        "Add",
	AddEq:      "AddEq",
	Addr:       "Addr",
	And:        "And",
	AndAnd:     "AndAnd",
	AndEq:      "AndEq",
	Arrow:      "Arrow",
	Call:       "Call",
	Cast:       "Cast",
	CastInit:   "CastInit",
	Comma:      "Comma",
	Cond:       "Cond",
	Div:        "Div",
	DivEq:      "DivEq",
	Dot:        "Dot",
	Eq:         "Eq",
	EqEq:       "EqEq",
	Gt:         "Gt",
	GtEq:       "GtEq",
	Index:      "Index",
	Indir:      "Indir",
	Lsh:        "Lsh",
	LshEq:      "LshEq",
	Lt:         "Lt",
	LtEq:       "LtEq",
	Minus:      "Minus",
	Mod:        "Mod",
	ModEq:      "ModEq",
	Mul:        "Mul",
	MulEq:      "MulEq",
	Name:       "Name",
	Not:        "Not",
	NotEq:      "NotEq",
	Number:     "Number",
	Offsetof:   "Offsetof",
	Or:         "Or",
	OrEq:       "OrEq",
	OrOr:       "OrOr",
	Paren:      "Paren",
	Plus:       "Plus",
	PostDec:    "PostDec",
	PostInc:    "PostInc",
	PreDec:     "PreDec",
	PreInc:     "PreInc",
	Rsh:        "Rsh",
	RshEq:      "RshEq",
	SizeofExpr: "SizeofExpr",
	SizeofType: "SizeofType",
	String:     "String",
	Sub:        "Sub",
	SubEq:      "SubEq",
	Twid:       "Twid",
	VaArg:      "VaArg",
	Xor:        "Xor",
	XorEq:      "XorEq",
}

func (op ExprOp) String() string {
	if 0 <= int(op) && int(op) <= len(exprOpString) {
		return exprOpString[op]
	}
	return fmt.Sprintf("ExprOp(%d)", op)
}

// Prefix is an initializer prefix.
type Prefix struct {
	Span  Span
	Dot   string // .Dot =
	Index *Expr  // [Index] =
}

// Init is an initializer expression.
type Init struct {
	SyntaxInfo
	Prefix []*Prefix // list of prefixes
	Expr   *Expr     // Expr
	Braced []*Init   // {Braced}
}
