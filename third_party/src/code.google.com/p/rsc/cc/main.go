// Copyright 2014 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"code.google.com/p/rsc/cc"
)

func main() {
	log.SetFlags(0)
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		cc.Read("<stdin>", os.Stdin)
	} else {
		for _, arg := range args {
			f, err := os.Open(arg)
			if err != nil {
				log.Fatal(err)
			}
			prog, err := cc.Read(arg, f)
			f.Close()
			do(prog, arg)
		}
	}
}

func do(prog *cc.Prog, file string) {
	var p Printer
	for _, decl := range prog.Decls {
		off := len(p.Bytes())
		p.Print(decl)
		if len(p.Bytes()) > off {
			p.Print(newline)
		}
	}
	os.Stdout.Write(p.Bytes())
}

// print an error; fprintf is a bad name but helps go vet.
func fprintf(span cc.Span, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Fprintf(os.Stderr, "%s:%d: %s\n", span.Start.File, span.Start.Line, msg)
}

var printed = map[interface{}]bool{}

func (p *Printer) printDecl(decl *cc.Decl) {
	if p.dup(decl) {
		return
	}

	p.Print(decl.Comments.Before)
	defer p.Print(decl.Comments.Suffix, decl.Comments.After)

	t := decl.Type
	if decl.Storage&cc.Typedef != 0 {
		if t.Kind == cc.Struct || t.Kind == cc.Union || t.Kind == cc.Union {
			if t.Tag == "" {
				t.Tag = decl.Name
			} else if decl.Name != t.Tag {
				fprintf(decl.Span, "typedef %s and tag %s do not match", decl.Name, t.Tag)
			}
			if t.Kind == cc.Enum {
				p.printEnumDecl(t)
			} else {
				p.printStructDecl(t)
			}
			return
		}
		p.Print("type ", decl.Name, " ", decl.Type)
		return
	}

	if decl.Name == "" {
		switch t.Kind {
		case cc.Struct, cc.Union:
			p.printStructDecl(t)
			return
		case cc.Enum:
			p.printEnumDecl(t)
			return
		}
		fprintf(decl.Span, "empty declaration")
		return
	}

	if t.Kind == cc.Func {
		p.printFuncDecl(decl)
		return
	}

	p.Print("var ", decl.Name, " ", decl.Type)
	if decl.Init != nil {
		p.Print(" = ", decl.Init)
	}
}

func (p *Printer) printStructDecl(t *cc.Type) {
	if p.dup(t) {
		return
	}
	if t.Kind == cc.Union {
		fprintf(t.Span, "cannot convert union")
		return
	}
	p.Print("type ", t.Tag, " struct {", indent)
	p.printStructBody(t)
	p.Print(unindent, newline, "}")
}

func (p *Printer) printStructBody(t *cc.Type) {
	for _, decl := range t.Decls {
		if decl.Name == "" {
			// Hope this is a struct definition.
			if decl.Type.Kind != cc.Struct {
				fprintf(decl.Span, "unnamed non-struct field of type %v", decl.Type)
				continue
			}
			p.printStructBody(decl.Type)
			continue
		}
		p.Print(newline, decl.Name, " ", decl.Type)
	}
}

func (p *Printer) printEnumDecl(t *cc.Type) {
	typeSuffix := ""
	if t.Tag != "" {
		typeSuffix = " " + t.Tag
		fprintf(t.Span, "cannot handle enum tags")
		return
	}
	p.Print("const (", indent)
	for i, decl := range t.Decls {
		p.Print(newline, decl.Name)
		if decl.Init == nil && i == 0 {
			if len(t.Decls) >= 2 && t.Decls[1].Init == nil {
				p.Print(typeSuffix, " = iota")
			} else {
				p.Print(typeSuffix, " = 0")
			}
		} else if decl.Init != nil {
			p.Print(typeSuffix, " = ", decl.Init.Expr)
			if i+1 < len(t.Decls) && t.Decls[i+1].Init == nil {
				p.Print(" + iota")
				if i > 0 {
					p.Print(fmt.Sprintf("-%d", i))
				}
			}
		}
	}
	p.Print(unindent, newline, ")")
}

func (p *Printer) printFuncDecl(decl *cc.Decl) {
	if decl.Body == nil {
		// wait for definition
		return
	}
	p.Print("func ", decl.Name, "(")
	for i, arg := range decl.Type.Decls {
		if i > 0 {
			p.Print(", ")
		}
		if arg.Name == "..." {
			p.Print("args ...interface{}")
			continue
		}
		p.Print(arg.Name, " ", arg.Type)
	}
	p.Print(")")
	if !decl.Type.Base.Is(cc.Void) {
		p.Print(" ", decl.Type.Base)
	}
	p.Print(" ", decl.Body, newline)
}

// printer, to move to a new file later

type printSpecial int

const (
	indent printSpecial = iota
	unindent
	untab
	newline
)

type Printer struct {
	buf    bytes.Buffer
	indent int
	html   bool

	printed map[interface{}]bool
	suffix  []cc.Comment // suffix comments to print at next newline
}

func (p *Printer) dup(x interface{}) bool {
	if p.printed[x] {
		return true
	}
	if p.printed == nil {
		p.printed = make(map[interface{}]bool)
	}
	p.printed[x] = true
	return false
}

func (p *Printer) StartHTML() {
	p.buf.WriteString("<pre>")
	p.html = true
}

func (p *Printer) EndHTML() {
	p.buf.WriteString("</pre>")
}

func (p *Printer) Bytes() []byte {
	return p.buf.Bytes()
}

func (p *Printer) String() string {
	return p.buf.String()
}

type exprPrec struct {
	expr *cc.Expr
	prec int
}

type nestBlock struct {
	stmt *cc.Stmt
	more bool
}

var htmlEscaper = strings.NewReplacer("<", "&lt;", ">", "&gt;", "&", "&amp;")

func (p *Printer) Print(args ...interface{}) {
	for _, arg := range args {
		switch arg := arg.(type) {
		default:
			fmt.Fprintf(&p.buf, "(?%T)", arg)
		case string:
			if p.html {
				htmlEscaper.WriteString(&p.buf, arg)
			} else {
				p.buf.WriteString(arg)
			}
		case exprPrec:
			p.printExpr(arg.expr, arg.prec)
		case *cc.Expr:
			p.printExpr(arg, precLow)
		case *cc.Prefix:
			p.printPrefix(arg)
		case *cc.Init:
			p.printInit(arg)
		case *cc.Prog:
			p.printProg(arg)
		case *cc.Stmt:
			p.printStmt(arg)
		case *cc.Type:
			p.printType(arg)
		case *cc.Decl:
			p.printDecl(arg)
		case cc.Storage:
			p.Print(arg.String())
		case []cc.Comment:
			for _, com := range arg {
				p.Print(com)
			}
		case cc.Comment:
			com := arg
			if com.Suffix {
				p.suffix = append(p.suffix, com)
			} else {
				for _, line := range strings.Split(com.Text, "\n") {
					p.Print(line, newline)
				}
			}
		case nestBlock:
			if arg.stmt.Op == cc.Block {
				p.Print(" ", arg.stmt)
			} else {
				p.Print(" {", indent, newline, arg.stmt, unindent, newline, "}")
			}
		case printSpecial:
			switch arg {
			default:
				fmt.Fprintf(&p.buf, "(?special:%d)", arg)
			case indent:
				p.indent++
			case unindent:
				p.indent--
			case untab:
				b := p.buf.Bytes()
				if len(b) > 0 && b[len(b)-1] == '\t' {
					p.buf.Truncate(len(b) - 1)
				}
			case newline:
				for _, com := range p.suffix {
					p.Print(" ", com.Text)
				}
				p.suffix = p.suffix[:0]
				p.buf.WriteString("\n")
				for i := 0; i < p.indent; i++ {
					p.buf.WriteByte('\t')
				}
			}
		}
	}
}

const (
	precNone = iota
	precArrow
	precAddr
	precMul
	precAdd
	precCmp
	precAndAnd
	precOrOr
	precComma
	precLow
)

var opPrec = []int{
	cc.Add:        precAdd,
	cc.AddEq:      precLow,
	cc.Addr:       precAddr,
	cc.And:        precMul,
	cc.AndAnd:     precAndAnd,
	cc.AndEq:      precLow,
	cc.Arrow:      precArrow,
	cc.Call:       precArrow,
	cc.Cast:       precAddr,
	cc.CastInit:   precAddr,
	cc.Comma:      precComma,
	cc.Cond:       precComma,
	cc.Div:        precMul,
	cc.DivEq:      precLow,
	cc.Dot:        precArrow,
	cc.Eq:         precLow,
	cc.EqEq:       precCmp,
	cc.Gt:         precCmp,
	cc.GtEq:       precCmp,
	cc.Index:      precArrow,
	cc.Indir:      precAddr,
	cc.Lsh:        precMul,
	cc.LshEq:      precLow,
	cc.Lt:         precCmp,
	cc.LtEq:       precCmp,
	cc.Minus:      precAddr,
	cc.Mod:        precMul,
	cc.ModEq:      precLow,
	cc.Mul:        precMul,
	cc.MulEq:      precLow,
	cc.Name:       precNone,
	cc.Not:        precAddr,
	cc.NotEq:      precCmp,
	cc.Number:     precNone,
	cc.Offsetof:   precAddr,
	cc.Or:         precAdd,
	cc.OrEq:       precLow,
	cc.OrOr:       precOrOr,
	cc.Paren:      precLow,
	cc.Plus:       precAddr,
	cc.PostDec:    precAddr,
	cc.PostInc:    precAddr,
	cc.PreDec:     precAddr,
	cc.PreInc:     precAddr,
	cc.Rsh:        precMul,
	cc.RshEq:      precLow,
	cc.SizeofExpr: precAddr,
	cc.SizeofType: precAddr,
	cc.String:     precNone,
	cc.Sub:        precAdd,
	cc.SubEq:      precLow,
	cc.Twid:       precAddr,
	cc.VaArg:      precAddr,
	cc.Xor:        precAdd,
	cc.XorEq:      precLow,
}

var opStr = []string{
	cc.Add:        "+",
	cc.AddEq:      "+=",
	cc.Addr:       "&",
	cc.And:        "&",
	cc.AndAnd:     "&&",
	cc.AndEq:      "&=",
	cc.Div:        "/",
	cc.DivEq:      "/=",
	cc.Eq:         "=",
	cc.EqEq:       "==",
	cc.Gt:         ">",
	cc.GtEq:       ">=",
	cc.Indir:      "*",
	cc.Lsh:        "<<",
	cc.LshEq:      "<<=",
	cc.Lt:         "<",
	cc.LtEq:       "<=",
	cc.Minus:      "-",
	cc.Mod:        "%",
	cc.ModEq:      "%=",
	cc.Mul:        "*",
	cc.MulEq:      "*=",
	cc.Not:        "!",
	cc.NotEq:      "!=",
	cc.Or:         "|",
	cc.OrEq:       "|=",
	cc.OrOr:       "||",
	cc.Plus:       "+",
	cc.PreDec:     "--",
	cc.PreInc:     "++",
	cc.Rsh:        ">>",
	cc.RshEq:      ">>=",
	cc.Sub:        "-",
	cc.SubEq:      "-=",
	cc.Twid:       "^",
	cc.Xor:        "^",
	cc.XorEq:      "^=",
	cc.SizeofExpr: "sizeof ",
}

func (p *Printer) printExpr(x *cc.Expr, prec int) {
	if x == nil {
		return
	}
	if p.html {
		fmt.Fprintf(&p.buf, "<span title='%s type %v'>", x.Op, x.XType)
		defer fmt.Fprintf(&p.buf, "</span>")
	}

	p.Print(x.Comments.Before)
	defer p.Print(x.Comments.Suffix, x.Comments.After)

	var newPrec int
	if 0 <= int(x.Op) && int(x.Op) < len(opPrec) {
		newPrec = opPrec[x.Op]
	}
	if prec < newPrec {
		p.Print("(")
		defer p.Print(")")
	}
	prec = newPrec

	var str string
	if 0 <= int(x.Op) && int(x.Op) < len(opStr) {
		str = opStr[x.Op]
	}
	if str != "" {
		if x.Right != nil {
			// binary operator
			// left associative
			p.Print(exprPrec{x.Left, prec}, " ", str, " ", exprPrec{x.Right, prec - 1})
		} else {
			// unary operator
			if (x.Op == cc.Plus || x.Op == cc.Minus || x.Op == cc.Addr) && x.Left.Op == x.Op ||
				x.Op == cc.Plus && x.Left.Op == cc.PreInc ||
				x.Op == cc.Minus && x.Left.Op == cc.PreDec {
				prec-- // force parenthesization +(+x) not ++x
			}
			p.Print(str, exprPrec{x.Left, prec})
		}
		return
	}

	// special cases
	switch x.Op {
	default:
		panic(fmt.Sprintf("printExpr missing case for %v", x.Op))

	case cc.Arrow:
		p.Print(exprPrec{x.Left, prec}, ".", x.Text)

	case cc.Call:
		p.Print(exprPrec{x.Left, precAddr}, "(")
		for i, y := range x.List {
			if i > 0 {
				p.Print(", ")
			}
			p.printExpr(y, precComma)
		}
		p.Print(")")

	case cc.Cast:
		p.Print("(", x.Type, ")", exprPrec{x.Left, prec})

	case cc.CastInit:
		p.Print("(", x.Type, ")", x.Init)

	case cc.Comma:
		for i, y := range x.List {
			if i > 0 {
				p.Print(", ")
			}
			p.printExpr(y, prec-1)
		}

	case cc.Cond:
		p.Print(exprPrec{x.List[0], prec - 1}, " ? ", exprPrec{x.List[1], prec}, " : ", exprPrec{x.List[2], prec})

	case cc.Dot:
		p.Print(exprPrec{x.Left, prec}, ".", x.Text)

	case cc.Index:
		p.Print(exprPrec{x.Left, prec}, "[", exprPrec{x.Right, precLow}, "]")

	case cc.Name, cc.Number:
		p.Print(x.Text)

	case cc.String:
		for i, str := range x.Texts {
			if i > 0 {
				p.Print(" ")
			}
			p.Print(str)
		}

	case cc.Offsetof:
		p.Print("offsetof(", x.Type, ", ", exprPrec{x.Left, precComma}, ")")

	case cc.Paren:
		p.Print("(", exprPrec{x.Left, prec}, ")")

	case cc.PostDec:
		p.Print(exprPrec{x.Left, prec}, "--")

	case cc.PostInc:
		p.Print(exprPrec{x.Left, prec}, "++")

	case cc.SizeofType:
		p.Print("sizeof(", x.Type, ")")

	case cc.VaArg:
		p.Print("va_arg(", exprPrec{x.Left, precComma}, ", ", x.Type, ")")
	}
}

func (p *Printer) printPrefix(x *cc.Prefix) {
	if x.Dot != "" {
		p.Print(".", x.Dot)
	} else {
		p.Print("[", x.Index, "]")
	}
}

func (p *Printer) printInit(x *cc.Init) {
	p.Print(x.Comments.Before)
	defer p.Print(x.Comments.Suffix, x.Comments.After)

	if len(x.Prefix) > 0 {
		for _, pre := range x.Prefix {
			p.Print(pre)
		}
		p.Print(" = ")
	}
	if x.Expr != nil {
		p.printExpr(x.Expr, precComma)
	} else {
		nl := len(x.Braced) > 0 && x.Braced[0].Span.Start.Line != x.Braced[len(x.Braced)-1].Span.End.Line
		p.Print("{")
		if nl {
			p.Print(indent)
		}
		for i, y := range x.Braced {
			if i > 0 {
				p.Print(",")
			}
			if nl {
				p.Print(newline)
			} else if i > 0 {
				p.Print(" ")
			}
			p.Print(y)
		}
		if nl {
			p.Print(unindent, newline)
		}
		p.Print("}")
	}
}

func (p *Printer) printProg(x *cc.Prog) {
	p.Print(x.Comments.Before)
	defer p.Print(x.Comments.Suffix, x.Comments.After)

	for _, decl := range x.Decls {
		p.Print(decl, newline)
	}
}

func (p *Printer) printStmt(x *cc.Stmt) {
	if len(x.Labels) > 0 {
		p.Print(untab, unindent, x.Comments.Before, indent, "\t")
		for _, lab := range x.Labels {
			p.Print(untab, unindent, lab.Comments.Before, indent, "\t")
			p.Print(untab)
			switch {
			case lab.Name != "":
				p.Print(lab.Name)
			case lab.Expr != nil:
				p.Print("case ", lab.Expr)
			default:
				p.Print("default")
			}
			p.Print(":", lab.Comments.Suffix, newline)
		}
	} else {
		p.Print(x.Comments.Before)
	}
	defer p.Print(x.Comments.Suffix, x.Comments.After)

	switch x.Op {
	case cc.ARGBEGIN:
		p.Print("ARGBEGIN{", indent, newline, x.Body, unindent, newline, "}ARGEND")

	case cc.Block:
		p.Print("{", indent)
		for _, b := range x.Block {
			p.Print(newline, b)
		}
		p.Print(unindent, newline, "}")

	case cc.Break:
		p.Print("break")

	case cc.Continue:
		p.Print("continue")

	case cc.Do:
		p.Print("do", nestBlock{x.Body, true}, " while(", x.Expr, ");")

	case cc.Empty:
		// ok

	case cc.For:
		p.Print("for ", x.Pre, "; ", x.Expr, "; ", x.Post, nestBlock{x.Body, false})

	case cc.If:
		p.Print("if ", x.Expr, nestBlock{x.Body, x.Else != nil})
		if x.Else != nil {
			p.Print(" else", nestBlock{x.Else, false})
		}

	case cc.Goto:
		p.Print("goto ", x.Text)

	case cc.Return:
		if x.Expr == nil {
			p.Print("return")
		} else {
			p.Print("return ", x.Expr)
		}

	case cc.StmtDecl:
		p.Print(x.Decl)

	case cc.StmtExpr:
		p.Print(x.Expr)

	case cc.Switch:
		p.Print("switch ", x.Expr, nestBlock{x.Body, false})

	case cc.While:
		p.Print("for ", x.Expr, nestBlock{x.Body, false})
	}
}

func (p *Printer) printType(t *cc.Type) {
	// Shouldn't happen but handle in case it does.
	p.Print(t.Comments.Before)
	defer p.Print(t.Comments.Suffix, t.Comments.After)

	if t == cc.BoolType {
		p.Print("bool")
		return
	}
	if typemap[t.Kind] != "" {
		p.Print(typemap[t.Kind])
		return
	}

	switch t.Kind {
	default:
		p.Print(t.String()) // hope for the best

	case cc.TypedefType:
		if typemap[t.Base.Kind] != "" && strings.ToLower(t.Name) == t.Name {
			p.Print(typemap[t.Base.Kind])
			return
		}
		p.Print(t.Name)

	case cc.Ptr:
		if t.Base.Is(cc.Func) {
			p.Print(t.Base)
			return
		}
		p.Print("*", t.Base)

	case cc.Func:
		p.Print("func(")
		for i, arg := range t.Decls {
			if i > 0 {
				p.Print(", ")
			}
			if arg.Name == "..." {
				p.Print("...interface{}")
				continue
			}
			p.Print(arg.Type)
		}
		p.Print(")")
		if !t.Base.Is(cc.Void) {
			p.Print(" ", t.Base)
		}

	case cc.Array:
		if t.Width == nil {
			p.Print("[]", t.Base) // TODO
			return
		}
		p.Print("[", t.Width, "]", t.Base)
	}
}

var typemap = map[cc.TypeKind]string{
	cc.Char:      "int8",
	cc.Uchar:     "uint8",
	cc.Short:     "int16",
	cc.Ushort:    "uint16",
	cc.Int:       "int",
	cc.Uint:      "uint",
	cc.Long:      "int32",
	cc.Ulong:     "uint32",
	cc.Longlong:  "int64",
	cc.Ulonglong: "uint64",
	cc.Float:     "float32",
	cc.Double:    "float64",
}

func (p *Printer) oldprintDecl(x *cc.Decl) {
	if x.Storage != 0 {
		p.Print(x.Storage, " ")
	}
	if x.Type == nil {
		p.Print(x.Name)
	} else {
		name := x.Name
		if x.Type.Kind == cc.Func && x.Body != nil {
			name = "\n" + name
		}
		p.Print(cc.TypedName{x.Type, name})
		if x.Name == "" {
			switch x.Type.Kind {
			case cc.Struct, cc.Union, cc.Enum:
				p.Print(" {", indent)
				for _, decl := range x.Type.Decls {
					p.Print(newline, decl)
				}
				p.Print(unindent, newline, "}")
			}
		}
	}
	if x.Init != nil {
		p.Print(" = ", x.Init)
	}
	if x.Body != nil {
		p.Print(newline, x.Body)
	}
}
