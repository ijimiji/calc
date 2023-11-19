package calculator

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"

	"github.com/shopspring/decimal"
)

const prec = 10

func New() *Calculator {
	return &Calculator{}
}

type Calculator struct {
}

func (c *Calculator) Eval(expr string) (string, error) {
	expr = strings.Join(strings.Fields(expr), "")
	expr = strings.ReplaceAll(expr, ",", ".")
	expr = "(" + expr + ")"
	fs := token.NewFileSet()
	tr, _ := parser.ParseExpr(expr)
	ast.Print(fs, tr)

	e, err := c.eval(tr)
	if err != nil {
		return "", err
	}
	if e == nil {
		return "", fmt.Errorf("invalid expr")
	}

	return e.(*ast.BasicLit).Value, nil
}

func (c *Calculator) RoundMath(number string) (string, error) {
	num, err := decimal.NewFromString(number)
	if err != nil {
		return "", err
	}

	return num.Round(0).String(), nil
}

func (c *Calculator) RoundSimple(number string) (string, error) {
	num, err := decimal.NewFromString(number)
	if err != nil {
		return "", err
	}
	_ = num

	return fmt.Sprint(num.IntPart()), nil
}

func (c *Calculator) RoundAccounting(number string) (string, error) {
	num, err := decimal.NewFromString(number)
	if err != nil {
		return "", err
	}
	_ = num

	return num.RoundBank(0).String(), nil
}

func (c *Calculator) eval(node ast.Expr) (ast.Expr, error) {
	switch node.(type) {
	case *ast.BinaryExpr:
		return c.evalBinary(node.(*ast.BinaryExpr))
	case *ast.BasicLit:
		return node, nil
	case *ast.UnaryExpr:
		return c.negative(node.(*ast.UnaryExpr)), nil
	case *ast.ParenExpr:
		return c.eval(node.(*ast.ParenExpr).X)
	}
	return nil, nil
}

func (c *Calculator) negative(node *ast.UnaryExpr) *ast.BasicLit {
	val := node.X.(*ast.BasicLit)
	val.Value = "-" + val.Value
	return val
}

func (c *Calculator) formatDecimal(d decimal.Decimal) string {
	return d.StringFixedBank(prec)
}

func (c *Calculator) evalBinary(node *ast.BinaryExpr) (*ast.BasicLit, error) {
	if node == nil {
		return nil, fmt.Errorf("nil node")
	}
	x, err := c.eval(node.X)
	if err != nil {
		return nil, fmt.Errorf("can't eval x: %v", err)
	}
	a, err := decimal.NewFromString(x.(*ast.BasicLit).Value)
	if err != nil {
		return nil, fmt.Errorf("can't parse number: %v", err)
	}

	y, err := c.eval(node.Y)
	if err != nil {
		return nil, fmt.Errorf("can't eval x: %v", err)
	}
	b, err := decimal.NewFromString(y.(*ast.BasicLit).Value)
	if err != nil {
		return nil, fmt.Errorf("can't parse number: %v", err)
	}

	switch node.Op {
	case token.ADD:
		res := &ast.BasicLit{Value: c.formatDecimal(a.Add(b)), Kind: token.FLOAT}
		fmt.Printf("%s %s %s = %s\n", a.String(), node.Op.String(), b.String(), res.Value)
		return res, nil
	case token.SUB:
		res := &ast.BasicLit{Value: c.formatDecimal(a.Sub(b)), Kind: token.FLOAT}
		fmt.Printf("%s %s %s = %s\n", a.String(), node.Op.String(), b.String(), res.Value)
		return res, nil
	case token.MUL:
		res := &ast.BasicLit{Value: c.formatDecimal(a.Mul(b)), Kind: token.FLOAT}
		fmt.Printf("%s %s %s = %s\n", a.String(), node.Op.String(), b.String(), res.Value)
		return res, nil
	case token.QUO:
		if b.Equal(decimal.NewFromInt(0)) {
			return nil, fmt.Errorf("zero division encountered")
		}
		res := &ast.BasicLit{Value: c.formatDecimal(a.DivRound(b, 6)), Kind: token.FLOAT}
		fmt.Printf("%s %s %s = %s\n", a.String(), node.Op.String(), b.String(), res.Value)
		return res, nil
	}
	return nil, fmt.Errorf("unknow op encountered")
}

func (c *Calculator) Format(s string) string {
	parts := strings.Split(s, ".")
	intPart := parts[0]
	result := make([]rune, 0)
	rIntPart := []rune(intPart)
	count := 0

	// Process integer part
	for i := len(intPart) - 1; i >= 0; i-- {
		count++
		result = append([]rune{rIntPart[i]}, result...)
		if count%3 == 0 && i != 0 {
			result = append([]rune{' '}, result...)
		}
	}

	// Process decimal part if exists
	if len(parts) > 1 {
		decimalPart := parts[1]
		result = append(result, []rune("."+decimalPart)...)
	}

	return string(result)
}
