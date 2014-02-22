package queryme

type Field string

type Value interface{}

type SortOrder struct {
	Field Field
	Ascending bool
}

type Predicate interface {
	Accept(visitor PredicateVisitor)
}

type PredicateVisitor interface {
	VisitNot(operand Predicate)
	VisitAnd(operands []Predicate)
	VisitOr(operands []Predicate)
	VisitEq(field Field, operands []Value)
	VisitLt(field Field, operand Value)
	VisitLe(field Field, operand Value)
	VisitGt(field Field, operand Value)
	VisitGe(field Field, operand Value)
	VisitFts(field Field, query string)
}

type Not struct {
	Operand Predicate
}

func (p Not) Accept(visitor PredicateVisitor) {
	visitor.VisitNot(p.Operand)
}

type And []Predicate

func (p And) Accept(visitor PredicateVisitor) {
	visitor.VisitAnd(p)
}

type Or []Predicate

func (p Or) Accept(visitor PredicateVisitor) {
	visitor.VisitOr(p)
}

type Eq struct {
	Field Field
	Operands []Value
}

func (p Eq) Accept(visitor PredicateVisitor) {
	visitor.VisitEq(p.Field, p.Operands)
}

type Lt struct {
	Field Field
	Operand Value
}

func (p Lt) Accept(visitor PredicateVisitor) {
	visitor.VisitLt(p.Field, p.Operand)
}

type Le struct {
	Field Field
	Operand Value
}

func (p Le) Accept(visitor PredicateVisitor) {
	visitor.VisitLe(p.Field, p.Operand)
}

type Gt struct {
	Field Field
	Operand Value
}

func (p Gt) Accept(visitor PredicateVisitor) {
	visitor.VisitGt(p.Field, p.Operand)
}

type Ge struct {
	Field Field
	Operand Value
}

func (p Ge) Accept(visitor PredicateVisitor) {
	visitor.VisitGe(p.Field, p.Operand)
}

type Fts struct {
	Field Field
	Query string
}

func (p Fts) Accept(visitor PredicateVisitor) {
	visitor.VisitFts(p.Field, p.Query)
}

type fieldsAccumulator struct {
	Index map[Field]struct{}
	Slice []Field
}

func (acc *fieldsAccumulator) saveField(field Field) {
	if _, ok := acc.Index[field]; !ok {
		acc.Index[field] = struct{}{}
		acc.Slice = append(acc.Slice, field)
	}
}

func (acc *fieldsAccumulator) VisitNot(operand Predicate) {
	operand.Accept(acc)
}

func (acc *fieldsAccumulator) VisitAnd(operands []Predicate) {
	for _, p := range operands {
		p.Accept(acc)
	}
}

func (acc *fieldsAccumulator) VisitOr(operands []Predicate) {
	for _, p := range operands {
		p.Accept(acc)
	}
}

func (acc *fieldsAccumulator) VisitEq(field Field, operands []Value) {
	acc.saveField(field)
}

func (acc *fieldsAccumulator) VisitLt(field Field, operand Value) {
	acc.saveField(field)
}

func (acc *fieldsAccumulator) VisitLe(field Field, operand Value) {
	acc.saveField(field)
}

func (acc *fieldsAccumulator) VisitGt(field Field, operand Value) {
	acc.saveField(field)
}

func (acc *fieldsAccumulator) VisitGe(field Field, operand Value) {
	acc.saveField(field)
}

func (acc *fieldsAccumulator) VisitFts(field Field, query string) {
	acc.saveField(field)
}

func Fields(predicate Predicate) []Field {
	var acc fieldsAccumulator
	acc.Index = make(map[Field]struct{})
	acc.Slice = make([]Field, 0, 4)
	predicate.Accept(&acc)
	return acc.Slice
}
