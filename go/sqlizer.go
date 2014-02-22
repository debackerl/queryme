package queryme

import (
	"bytes"
	"fmt"
	"strings"
)

func ToSql(predicate Predicate, allowedFields []Field) (sql string, values []interface{}) {
	b := newQueryBuilder(allowedFields)
	predicate.Accept(b)
	sql = b.Sql()
	values = b.Values()
	return
}

type queryBuilder struct {
	VariablePlaceHolder string
	EscapedIdentifierGenerator func(string)string

	allowedFields map[Field]struct{}
	sql *bytes.Buffer
	values []interface{}
}

func newQueryBuilder(allowedFields []Field) *queryBuilder {
	allowedFieldsIndex := make(map[Field]struct{})

	for _, f := range allowedFields {
		if _, ok := allowedFieldsIndex[f]; !ok {
			allowedFieldsIndex[f] = struct{}{}
		}
	}

	b := &queryBuilder{
		"?",
		func(id string) string { return "`" + strings.Replace(id, "`", "``", -1) + "`" },
		allowedFieldsIndex,
		bytes.NewBuffer(make([]byte, 0, 64)),
		make([]interface{}, 0, 8)}

	return b
}

func (b *queryBuilder) Sql() string {
	return b.sql.String()
}

func (b *queryBuilder) Values() []interface{} {
	return b.values
}

func (b *queryBuilder) AppendSql(sql string) {
	b.sql.WriteString(sql)
}

func (b *queryBuilder) AppendId(id Field) {
	if _, ok := b.allowedFields[id]; !ok {
		panic(fmt.Errorf("unauthorized field accessed: %q", id))
	}

	b.sql.WriteString(b.EscapedIdentifierGenerator(string(id)))
}

func (b *queryBuilder) AppendValue(value interface{}) {
	b.values = append(b.values, value)
}

func (b *queryBuilder) VisitNot(operand Predicate) {
	b.AppendSql("(NOT ")
	operand.Accept(b)
	b.AppendSql(")")
}

func (b *queryBuilder) VisitPredicates(sqlOperator string, operands []Predicate) {
	b.AppendSql("(")
	for i, p := range operands {
		if i > 0 {
			b.AppendSql(" ")
			b.AppendSql(sqlOperator)
			b.AppendSql(" ")
		}
		p.Accept(b)
	}
	b.AppendSql(")")
}

func (b *queryBuilder) VisitAnd(operands []Predicate) {
	b.VisitPredicates("AND", operands)
}

func (b *queryBuilder) VisitOr(operands []Predicate) {
	b.VisitPredicates("OR", operands)
}

func (b *queryBuilder) VisitEq(field Field, operands []Value) {
	switch len(operands) {
		case 0:
			b.AppendSql("false")
		case 1:
			b.AppendId(field)
			b.AppendSql("=")
			b.AppendSql(b.VariablePlaceHolder)
			b.AppendValue(operands[0])
		default:
			b.AppendId(field)
			b.AppendSql(" IN (")
			for i, op := range operands {
				if i > 0 {
					b.AppendSql(",")
				}
				b.AppendSql(b.VariablePlaceHolder)
				b.AppendValue(op)
			}
			b.AppendSql(")")
	}
}

func (b *queryBuilder) VisitLt(field Field, operand Value) {
	b.AppendId(field)
	b.AppendSql("<")
	b.AppendSql(b.VariablePlaceHolder)
	b.AppendValue(operand)
}

func (b *queryBuilder) VisitLe(field Field, operand Value) {
	b.AppendId(field)
	b.AppendSql("<=")
	b.AppendSql(b.VariablePlaceHolder)
	b.AppendValue(operand)
}

func (b *queryBuilder) VisitGt(field Field, operand Value) {
	b.AppendId(field)
	b.AppendSql(">")
	b.AppendSql(b.VariablePlaceHolder)
	b.AppendValue(operand)
}

func (b *queryBuilder) VisitGe(field Field, operand Value) {
	b.AppendId(field)
	b.AppendSql(">=")
	b.AppendSql(b.VariablePlaceHolder)
	b.AppendValue(operand)
}

func (b *queryBuilder) VisitFts(field Field, query string) {
	b.AppendSql("MATCH (")
	b.AppendId(field)
	b.AppendSql(") AGAINST (")
	b.AppendSql(b.VariablePlaceHolder)
	b.AppendSql(")")
	b.AppendValue(query)
}
