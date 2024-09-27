package expr

type ExprType int

const (
	PRIMARY ExprType = iota
	POSTFIX
	UNARY
	MULTIPLICATIVE
	ADDITIVE
	SHIFT
	RELATIONAL
	EQUALITY
	BITWISE_AND
	BITWISE_XOR
	BITWISE_OR
	LOGICAL_AND
	LOGICAL_OR
	CONDITIONAL
	ASSIGNMENT
)

type Expr struct {
	expr_type ExprType

	primary_expr        *PrimaryExpr
	postfix_expr        *PostfixExpr
	unary_expr          *UnaryExpr
	multiplicative_expr *MultiplicativeExpr
	additive_expr       *AdditiveExpr
	shift_expr          *ShiftExpr
	relational_expr     *RelationalExpr
	equality_expr       *EqualityExpr
	bitwise_and_expr    *BitwiseAndExpr
	bitwise_xor_expr    *BitwiseXorExpr
	bitwise_or_expr     *BitwiseOrExpr
	logical_and_expr    *LogicalAndExpr
	logical_or_expr     *LogicalOrExpr
	conditional_expr    *ConditionalExpr
	assignment_expr     *AssignmentExpr
}

func (this *Expr) InitPrimaryExpr(primary_expr *PrimaryExpr) {
	this.expr_type = PRIMARY

	this.primary_expr = primary_expr
}

func (this *Expr) InitPostfixExpr(postfix_expr *PostfixExpr) {
	this.expr_type = POSTFIX

	this.postfix_expr = postfix_expr
}

func (this *Expr) InitUnaryExpr(unary_expr *UnaryExpr) {
	this.expr_type = UNARY

	this.unary_expr = unary_expr
}

func (this *Expr) InitMultiplicativeExpr(multiplicative_expr *MultiplicativeExpr) {
	this.expr_type = MULTIPLICATIVE

	this.multiplicative_expr = multiplicative_expr
}

func (this *Expr) InitAdditiveExpr(additive_expr *AdditiveExpr) {
	this.expr_type = ADDITIVE

	this.additive_expr = additive_expr
}

func (this *Expr) InitShiftExpr(shift_expr *ShiftExpr) {
	this.expr_type = SHIFT

	this.shift_expr = shift_expr
}

func (this *Expr) InitRelationalExpr(relational_expr *RelationalExpr) {
	this.expr_type = RELATIONAL

	this.relational_expr = relational_expr
}

func (this *Expr) InitEqualityExpr(equality_expr *EqualityExpr) {
	this.expr_type = EQUALITY

	this.equality_expr = equality_expr
}

func (this *Expr) InitBitwiseAndExpr(bitwise_and_expr *BitwiseAndExpr) {
	this.expr_type = BITWISE_AND

	this.bitwise_and_expr = bitwise_and_expr
}

func (this *Expr) InitBitwiseXorExpr(bitwise_xor_expr *BitwiseXorExpr) {
	this.expr_type = BITWISE_XOR

	this.bitwise_xor_expr = bitwise_xor_expr
}

func (this *Expr) InitBitwiseOrExpr(bitwise_or_expr *BitwiseOrExpr) {
	this.expr_type = BITWISE_OR

	this.bitwise_or_expr = bitwise_or_expr
}

func (this *Expr) InitLogicalAndExpr(logical_and_expr *LogicalAndExpr) {
	this.expr_type = LOGICAL_AND

	this.logical_and_expr = logical_and_expr
}

func (this *Expr) InitLogicalOrExpr(logical_or_expr *LogicalOrExpr) {
	this.expr_type = LOGICAL_OR

	this.logical_or_expr = logical_or_expr
}

func (this *Expr) InitConditionalExpr(conditional_expr *ConditionalExpr) {
	this.expr_type = CONDITIONAL

	this.conditional_expr = conditional_expr
}

func (this *Expr) InitAssignmentExpr(assignment_expr *AssignmentExpr) {
	this.expr_type = ASSIGNMENT

	this.assignment_expr = assignment_expr
}

func (this *Expr) ExprType() ExprType {
	return this.expr_type
}

func (this *Expr) PrimaryExpr() *PrimaryExpr {
	return this.primary_expr
}

func (this *Expr) PostfixExpr() *PostfixExpr {
	return this.postfix_expr
}

func (this *Expr) UnaryExpr() *UnaryExpr {
	return this.unary_expr
}

func (this *Expr) MultiplicativeExpr() *MultiplicativeExpr {
	return this.multiplicative_expr
}

func (this *Expr) AdditiveExpr() *AdditiveExpr {
	return this.additive_expr
}

func (this *Expr) ShiftExpr() *ShiftExpr {
	return this.shift_expr
}

func (this *Expr) RelationalExpr() *RelationalExpr {
	return this.relational_expr
}

func (this *Expr) BitwiseAndExpr() *BitwiseAndExpr {
	return this.bitwise_and_expr
}

func (this *Expr) BitwiseXorExpr() *BitwiseXorExpr {
	return this.bitwise_xor_expr
}

func (this *Expr) BitwiseOrExpr() *BitwiseOrExpr {
	return this.bitwise_or_expr
}

func (this *Expr) LogicalAndExpr() *LogicalAndExpr {
	return this.logical_and_expr
}

func (this *Expr) LogicalOrExpr() *LogicalOrExpr {
	return this.logical_or_expr
}

func (this *Expr) EqualityExpr() *EqualityExpr {
	return this.equality_expr
}

func (this *Expr) ConditionalExpr() *ConditionalExpr {
	return this.conditional_expr
}

func (this *Expr) AssignmentExpr() *AssignmentExpr {
	return this.assignment_expr
}
