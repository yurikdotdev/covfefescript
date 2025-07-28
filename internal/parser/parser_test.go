package parser

import (
	"github.com/yurikdotdev/covfefescript/internal/ast"
	"github.com/yurikdotdev/covfefescript/internal/lexer"
	"fmt"
	"testing"
)

func TestLookStatements(t *testing.T) {
	input := `
LOOK x IS 5!
LOOK y IS YUGE!
LOOK foobar IS y!
`
	program := parseAndCheckErrors(t, input, 0)
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testLookStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

func TestReturnStatements(t *testing.T) {
	input := `
GIVE_ME 5!
GIVE_ME YUGE!
GIVE_ME foobar!
`
	program := parseAndCheckErrors(t, input, 0)
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
	}

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt not *ast.ReturnStatement. got=%T", stmt)
			continue
		}
		if returnStmt.TokenLiteral() != "GIVE_ME" {
			t.Errorf("returnStmt.TokenLiteral not 'GIVE_ME', got %q", returnStmt.TokenLiteral())
		}
	}
}

func TestForLoopStatement(t *testing.T) {
	input := `KEEP_WINNING (x < 10) { LOOK x IS x + 1! }`

	program := parseAndCheckErrors(t, input, 0)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ForLoopStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ForLoopStatement. got=%T",
			program.Statements[0])
	}

	if !testInfixExpression(t, stmt.Condition, "x", "<", 10) {
		return
	}

	if len(stmt.Body.Statements) != 1 {
		t.Errorf("loop body does not contain 1 statement. got=%d",
			len(stmt.Body.Statements))
	}

	bodyStmt, ok := stmt.Body.Statements[0].(*ast.LookStatement)
	if !ok {
		t.Fatalf("statement in body is not *ast.LookStatement. got=%T",
			stmt.Body.Statements[0])
	}

	if !testLookStatement(t, bodyStmt, "x") {
		return
	}
}

func TestBreakStatement(t *testing.T) {
	input := `IT_WAS_RIGGED!`

	program := parseAndCheckErrors(t, input, 0)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.BreakStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.BreakStatement. got=%T",
			program.Statements[0])
	}

	if stmt.TokenLiteral() != "IT_WAS_RIGGED" {
		t.Errorf("stmt.TokenLiteral not 'IT_WAS_RIGGED', got %q",
			stmt.TokenLiteral())
	}
}

func TestContinueStatement(t *testing.T) {
	input := `TIRED_OF_WINNING!`

	program := parseAndCheckErrors(t, input, 0)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ContinueStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ContinueStatement. got=%T",
			program.Statements[0])
	}

	if stmt.TokenLiteral() != "TIRED_OF_WINNING" {
		t.Errorf("stmt.TokenLiteral not 'TIRED_OF_WINNING', got %q",
			stmt.TokenLiteral())
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar!"
	program := parseAndCheckErrors(t, input, 0)
	testSingleExpressionStatement(t, program)
}

func TestMoneyLiteralExpression(t *testing.T) {
	input := "5!"
	program := parseAndCheckErrors(t, input, 0)
	stmt := testSingleExpressionStatement(t, program)
	literal, ok := stmt.Expression.(*ast.MoneyLiteral)
	if !ok {
		t.Fatalf("exp not *ast.MoneyLiteral. got=%T", stmt.Expression)
	}
	if literal.Value != 5 {
		t.Errorf("literal.Value not %d. got=%d", 5, literal.Value)
	}
}

func TestStringLiteralExpression(t *testing.T) {
	input := `"hello world"!`
	program := parseAndCheckErrors(t, input, 0)
	stmt := testSingleExpressionStatement(t, program)
	literal, ok := stmt.Expression.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("exp not *ast.StringLiteral. got=%T", stmt.Expression)
	}
	if literal.Value != "hello world" {
		t.Errorf("literal.Value not %q. got=%q", "hello world", literal.Value)
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		value    any
	}{
		{"!5!", "!", 5},
		{"-15!", "-", 15},
		{"!YUGE!", "!", true},
		{"!LOSER!", "!", false},
	}

	for _, tt := range prefixTests {
		program := parseAndCheckErrors(t, tt.input, 0)
		stmt := testSingleExpressionStatement(t, program)
		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt is not ast.PrefixExpression. got=%T", stmt.Expression)
		}
		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s", tt.operator, exp.Operator)
		}
		if !testLiteralExpression(t, exp.Right, tt.value) {
			return
		}
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  any
		operator   string
		rightValue any
	}{
		{"5 + 5!", 5, "+", 5},
		{"5 - 5!", 5, "-", 5},
		{"5 * 5!", 5, "*", 5},
		{"5 / 5!", 5, "/", 5},
		{"5 > 5!", 5, ">", 5},
		{"5 < 5!", 5, "<", 5},
		{"5 == 5!", 5, "==", 5},
		{"5 != 5!", 5, "!=", 5},
		{"YUGE == YUGE!", true, "==", true},
		{"YUGE != LOSER!", true, "!=", false},
		{"LOSER == LOSER!", false, "==", false},
	}

	for _, tt := range infixTests {
		program := parseAndCheckErrors(t, tt.input, 0)
		stmt := testSingleExpressionStatement(t, program)
		testInfixExpression(t, stmt.Expression, tt.leftValue, tt.operator, tt.rightValue)
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"-a * b", "((-a) * b)"},
		{"!-a", "(!(-a))"},
		{"a + b + c", "((a + b) + c)"},
		{"a + b - c", "((a + b) - c)"},
		{"a * b * c", "((a * b) * c)"},
		{"a * b / c", "((a * b) / c)"},
		{"a + b / c", "(a + (b / c))"},
		{"a + b * c + d / e - f", "(((a + (b * c)) + (d / e)) - f)"},
		{"3 + 4! -5 * 5", "(3 + 4)((-5) * 5)"},
		{"5 > 4 == 3 < 4", "((5 > 4) == (3 < 4))"},
		{"3 + 4 * 5 == 3 * 1 + 4 * 5", "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))"},
		{"1 + (2 + 3) + 4", "((1 + (2 + 3)) + 4)"},
		{"(5 + 5) * 2", "((5 + 5) * 2)"},
		{"2 / (5 + 5)", "(2 / (5 + 5))"},
		{"-(5 + 5)", "(-(5 + 5))"},
		{"!(YUGE == YUGE)", "(!(YUGE == YUGE))"},
	}

	for _, tt := range tests {
		program := parseAndCheckErrors(t, tt.input, 0)
		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func TestIfExpression(t *testing.T) {
	input := `BELIEVE_ME x < y { x }`

	program := parseAndCheckErrors(t, input, 0)
	stmt := testSingleExpressionStatement(t, program)

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T", stmt.Expression)
	}

	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d\n", len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T", exp.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if exp.Alternative != nil {
		t.Errorf("exp.Alternative.Statements was not nil. got=%+v", exp.Alternative)
	}
}

func TestIfElseExpression(t *testing.T) {
	input := `BELIEVE_ME x < y { x } FAKE_NEWS { y }`

	program := parseAndCheckErrors(t, input, 0)
	stmt := testSingleExpressionStatement(t, program)

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T", stmt.Expression)
	}

	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d\n", len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T", exp.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if exp.Alternative == nil {
		t.Fatalf("exp.Alternative.Statements was nil.")
	}

	if len(exp.Alternative.Statements) != 1 {
		t.Errorf("exp.Alternative.Statements does not contain 1 statements. got=%d\n", len(exp.Alternative.Statements))
	}

	alternative, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T", exp.Alternative.Statements[0])
	}

	if !testIdentifier(t, alternative.Expression, "y") {
		return
	}
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `MAKE_IT_BIG add(x, y) { x + y! }`

	program := parseAndCheckErrors(t, input, 0)
	stmt := testSingleExpressionStatement(t, program)

	function, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.FunctionLiteral. got=%T", stmt.Expression)
	}

	if len(function.Parameters) != 2 {
		t.Fatalf("function literal parameters wrong. want 2, got=%d\n", len(function.Parameters))
	}

	testLiteralExpression(t, function.Parameters[0], "x")
	testLiteralExpression(t, function.Parameters[1], "y")

	if len(function.Body.Statements) != 1 {
		t.Fatalf("function.Body.Statements has not 1 statements. got=%d\n", len(function.Body.Statements))
	}

	bodyStmt, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("function body stmt is not ast.ExpressionStatement. got=%T", function.Body.Statements[0])
	}

	testInfixExpression(t, bodyStmt.Expression, "x", "+", "y")
}

func TestCallExpressionParsing(t *testing.T) {
	input := "add(1, 2 * 3, 4 + 5)!"

	program := parseAndCheckErrors(t, input, 0)
	stmt := testSingleExpressionStatement(t, program)

	exp, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.CallExpression. got=%T", stmt.Expression)
	}

	if !testIdentifier(t, exp.Function, "add") {
		return
	}

	if len(exp.Arguments) != 3 {
		t.Fatalf("wrong length of arguments. got=%d", len(exp.Arguments))
	}

	testLiteralExpression(t, exp.Arguments[0], 1)
	testInfixExpression(t, exp.Arguments[1], 2, "*", 3)
	testInfixExpression(t, exp.Arguments[2], 4, "+", 5)
}

func TestParsingArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"

	program := parseAndCheckErrors(t, input, 0)
	stmt := testSingleExpressionStatement(t, program)
	array, ok := stmt.Expression.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("exp not ast.ArrayLiteral. got=%T", stmt.Expression)
	}
	if len(array.Elements) != 3 {
		t.Fatalf("len(array.Elements) not 3. got=%d", len(array.Elements))
	}
	testMoneyLiteral(t, array.Elements[0], 1)
	testInfixExpression(t, array.Elements[1], 2, "*", 2)
	testInfixExpression(t, array.Elements[2], 3, "+", 3)
}

func TestParsingIndexExpressions(t *testing.T) {
	input := "myArray[1 + 1]"

	program := parseAndCheckErrors(t, input, 0)
	stmt := testSingleExpressionStatement(t, program)
	indexExp, ok := stmt.Expression.(*ast.IndexExpression)
	if !ok {
		t.Fatalf("exp not *ast.IndexExpression. got=%T", stmt.Expression)
	}
	if !testIdentifier(t, indexExp.Left, "myArray") {
		return
	}
	if !testInfixExpression(t, indexExp.Index, 1, "+", 1) {
		return
	}
}

func TestParsingHashLiterals(t *testing.T) {
	input := `{"one": 1, "two": 2, "three": 3}`
	program := parseAndCheckErrors(t, input, 0)
	stmt := testSingleExpressionStatement(t, program)
	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("exp is not ast.HashLiteral. got=%T", stmt.Expression)
	}
	if len(hash.Pairs) != 3 {
		t.Errorf("hash.Pairs has wrong length. got=%d", len(hash.Pairs))
	}
	expected := map[string]int64{
		"one":   1,
		"two":   2,
		"three": 3,
	}
	for key, value := range hash.Pairs {
		literal, ok := key.(*ast.StringLiteral)
		if !ok {
			t.Errorf("key is not ast.StringLiteral. got=%T", key)
		}
		expectedValue := expected[literal.String()]
		testMoneyLiteral(t, value, expectedValue)
	}
}

func TestLookStatementErrors(t *testing.T) {
	input := `
LOOK x 5!
LOOK y IS!
LOOK = 10!
`
	parseAndCheckErrors(t, input, 3)
}

func parseAndCheckErrors(t *testing.T, input string, expectedErrors int) *ast.Program {
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	errors := p.Errors()
	if len(errors) != expectedErrors {
		t.Errorf("parser has wrong number of errors. expected=%d, got=%d", expectedErrors, len(errors))
		for _, msg := range errors {
			t.Errorf("parser error: %q", msg)
		}
		t.FailNow()
	}
	return program
}

func testLookStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "LOOK" {
		t.Errorf("s.TokenLiteral not 'LOOK'. got=%q", s.TokenLiteral())
		return false
	}
	lookStmt, ok := s.(*ast.LookStatement)
	if !ok {
		t.Errorf("s not *ast.LookStatement. got=%T", s)
		return false
	}
	if lookStmt.Name.Value != name {
		t.Errorf("lookStmt.Name.Value not '%s'. got=%s", name, lookStmt.Name.Value)
		return false
	}
	if lookStmt.Name.TokenLiteral() != name {
		t.Errorf("lookStmt.Name.TokenLiteral() not '%s'. got=%s", name, lookStmt.Name.TokenLiteral())
		return false
	}
	return true
}

func testSingleExpressionStatement(t *testing.T, program *ast.Program) *ast.ExpressionStatement {
	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}
	return stmt
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected any) bool {
	switch v := expected.(type) {
	case int:
		return testMoneyLiteral(t, exp, int64(v))
	case int64:
		return testMoneyLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)
	}
	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}

func testInfixExpression(t *testing.T, exp ast.Expression, left any, operator string, right any) bool {
	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not ast.InfixExpression. got=%T(%s)", exp, exp)
		return false
	}
	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}
	if opExp.Operator != operator {
		t.Errorf("exp.Operator is not '%s'. got=%q", operator, opExp.Operator)
		return false
	}
	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}
	return true
}

func testMoneyLiteral(t *testing.T, il ast.Expression, value int64) bool {
	money, ok := il.(*ast.MoneyLiteral)
	if !ok {
		t.Errorf("il not *ast.MoneyLiteral. got=%T", il)
		return false
	}
	if money.Value != value {
		t.Errorf("money.Value not %d. got=%d", value, money.Value)
		return false
	}
	if money.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("money.TokenLiteral not %d. got=%s", value, money.TokenLiteral())
		return false
	}
	return true
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp not *ast.Identifier. got=%T", exp)
		return false
	}
	if ident.Value != value {
		t.Errorf("ident.Value not %s. got=%s", value, ident.Value)
		return false
	}
	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral not %s. got=%s", value, ident.TokenLiteral())
		return false
	}
	return true
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
	bo, ok := exp.(*ast.Boolean)
	if !ok {
		t.Errorf("exp not *ast.Boolean. got=%T", exp)
		return false
	}
	if bo.Value != value {
		t.Errorf("bo.Value not %t. got=%t", value, bo.Value)
		return false
	}
	expectedLiteral := "YUGE"
	if !value {
		expectedLiteral = "LOSER"
	}
	if bo.TokenLiteral() != expectedLiteral {
		t.Errorf("bo.TokenLiteral not %s. got=%s", expectedLiteral, bo.TokenLiteral())
		return false
	}
	return true
}
