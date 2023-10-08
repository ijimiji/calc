package calculator

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"math/big"
)

const prec = 10

func New() *Calculator {
	return &Calculator{}
}

type Calculator struct {
}

func (c *Calculator) Eval(expr string) (string, error) {
	fs := token.NewFileSet()
	tr, _ := parser.ParseExpr(expr)
	ast.Print(fs, tr)

	evaled := c.eval(tr).(*ast.BasicLit)

	return evaled.Value, nil
}

func (c *Calculator) eval(node ast.Expr) ast.Expr {
	switch node.(type) {
	case *ast.BinaryExpr:
		return c.evalBinary(node.(*ast.BinaryExpr))
	case *ast.BasicLit:
		return node
	case *ast.UnaryExpr:
		return c.negative(node.(*ast.UnaryExpr))
	}
	return nil
}

func (c *Calculator) negative(node *ast.UnaryExpr) *ast.BasicLit {
	return node.X.(*ast.BasicLit)
}

func (c *Calculator) evalBinary(node *ast.BinaryExpr) *ast.BasicLit {
	a, _ := new(big.Float).SetPrec(prec).SetString(c.eval(node.X).(*ast.BasicLit).Value)
	b, _ := new(big.Float).SetPrec(prec).SetString(c.eval(node.Y).(*ast.BasicLit).Value)

	switch node.Op {
	case token.ADD:
		res := &ast.BasicLit{Value: new(big.Float).Add(a, b).SetPrec(prec).String(), Kind: token.FLOAT}
		fmt.Printf("%s %s %s = %s\n", a.String(), node.Op.String(), b.String(), res.Value)
		return res
	case token.SUB:
		res := &ast.BasicLit{Value: new(big.Float).Sub(a, b).String(), Kind: token.FLOAT}
		fmt.Printf("%s %s %s = %s\n", a.String(), node.Op.String(), b.String(), res.Value)
		return res
	case token.MUL:
		res := &ast.BasicLit{Value: new(big.Float).Mul(a, b).String(), Kind: token.FLOAT}
		fmt.Printf("%s %s %s = %s\n", a.String(), node.Op.String(), b.String(), res.Value)
		return res
	case token.QUO:
		res := &ast.BasicLit{Value: new(big.Float).Quo(a, b).String(), Kind: token.FLOAT}
		fmt.Printf("%s %s %s = %s\n", a.String(), node.Op.String(), b.String(), res.Value)
		return res
	}
	return nil
}
