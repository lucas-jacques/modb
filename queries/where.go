package queries

import (
	"strings"
)

type simpleCond struct {
	right Expr
	op    string
	left  Expr
}

func newSimpleCond(left Expr, op string, right Expr) simpleCond {
	return simpleCond{left: left, op: op, right: right}
}

type logicalExpr struct {
	conds    []Expr
	operator string
}

func (w logicalExpr) Build(params ParamsSet) (string, []any) {
	var parts []string
	var values []any

	for _, cond := range w.conds {
		part, condValues := cond.Build(params)
		parts = append(parts, part)
		values = append(values, condValues...)
	}

	return "(" + strings.Join(parts, " "+w.operator+" ") + ")", values
}

func And(conds ...Expr) logicalExpr {
	return logicalExpr{conds: conds, operator: "AND"}
}

func Or(conds ...Expr) logicalExpr {
	return logicalExpr{conds: conds, operator: "OR"}
}

func EQ(left Expr, right Expr) simpleCond {
	return newSimpleCond(left, "=", right)
}

func NE(left Expr, right Expr) simpleCond {
	return newSimpleCond(left, "!=", right)
}

func GT(left Expr, right Expr) simpleCond {
	return newSimpleCond(left, ">", right)
}

func GTE(left Expr, right Expr) simpleCond {
	return newSimpleCond(left, ">=", right)
}

func LT(left Expr, right Expr) simpleCond {
	return newSimpleCond(left, "<", right)
}

func LTE(left Expr, right Expr) simpleCond {
	return newSimpleCond(left, "<=", right)
}

func (w simpleCond) Build(params ParamsSet) (string, []any) {
	leftPart, leftValues := w.left.Build(params)
	rightPart, rightValues := w.right.Build(params)

	b := strings.Builder{}
	b.WriteString("(")
	b.WriteString(leftPart)
	b.WriteString(" ")
	b.WriteString(w.op)
	b.WriteString(" ")
	b.WriteString(rightPart)
	b.WriteString(")")

	return b.String(), append(leftValues, rightValues...)
}

type inCond struct {
	col  Expr
	vals []Expr
}

func (w inCond) Build(params ParamsSet) (string, []any) {
	var rightParts []string
	var values []any

	leftPart, leftValues := w.col.Build(params)

	values = append(values, leftValues...)

	for _, val := range w.vals {
		part, valValues := val.Build(params)
		rightParts = append(rightParts, part)
		values = append(values, valValues...)
	}

	return "(" + leftPart + " IN (" + strings.Join(rightParts, ", ") + "))", values
}

func IN(expr Expr, values []Expr) inCond {
	return inCond{col: expr, vals: values}
}
