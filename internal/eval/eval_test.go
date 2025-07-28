package eval

import (
	"testing"

	"github.com/yurikdotdev/covfefescript/internal/ast"
	"github.com/yurikdotdev/covfefescript/internal/lexer"
	"github.com/yurikdotdev/covfefescript/internal/object"
	"github.com/yurikdotdev/covfefescript/internal/parser"
)

func TestEvalMoneyExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5!", 5},
		{"10!", 10},
		{"-5!", -5},
		{"-10!", -10},
		{"5 + 5 + 5 + 5 - 10!", 10},
		{"2 * 2 * 2 * 2 * 2!", 32},
		{"-50 + 100 + -50!", 0},
		{"5 * 2 + 10!", 20},
		{"5 + 2 * 10!", 25},
		{"20 + 2 * -10!", 0},
		{"50 / 2 * 2 + 10!", 60},
		{"2 * (5 + 10)!", 30},
		{"3 * 3 * 3 + 10!", 37},
		{"3 * (3 * 3) + 10!", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10!", 50},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testMoneyObject(t, evaluated, tt.expected)
	}
}

func TestEvalTruthExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"YUGE!", true},
		{"LOSER!", false},
		{"1 < 2!", true},
		{"1 > 2!", false},
		{"1 < 1!", false},
		{"1 > 1!", false},
		{"1 == 1!", true},
		{"1 != 1!", false},
		{"1 == 2!", false},
		{"1 != 2!", true},
		{"YUGE == YUGE!", true},
		{"LOSER == LOSER!", true},
		{"YUGE == LOSER!", false},
		{"YUGE != LOSER!", true},
		{"(1 < 2) == YUGE!", true},
		{"(1 < 2) == LOSER!", false},
		{"(1 > 2) == YUGE!", false},
		{"(1 > 2) == LOSER!", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testTruthObject(t, evaluated, tt.expected)
	}
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!YUGE!", false},
		{"!LOSER!", true},
		{"!5!", false},
		{"!!YUGE!", true},
		{"!!LOSER!", false},
		{"!!5!", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testTruthObject(t, evaluated, tt.expected)
	}
}

func TestIfElseExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{"BELIEVE_ME YUGE { 10 }!", 10},
		{"BELIEVE_ME LOSER { 10 }!", nil},
		{"BELIEVE_ME 1 { 10 }!", 10},
		{"BELIEVE_ME 1 < 2 { 10 }!", 10},
		{"BELIEVE_ME 1 > 2 { 10 }!", nil},
		{"BELIEVE_ME 1 > 2 { 10 } FAKE_NEWS { 20 }!", 20},
		{"BELIEVE_ME 1 < 2 { 10 } FAKE_NEWS { 20 }!", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testMoneyObject(t, evaluated, int64(integer))
		} else {
			testCovfefeObject(t, evaluated)
		}
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"GIVE_ME 10!", 10},
		{"GIVE_ME 10! 9!", 10},
		{"GIVE_ME 2 * 5! 9!", 10},
		{"9! GIVE_ME 2 * 5! 9!", 10},
		{`
BELIEVE_ME 10 > 1 {
  BELIEVE_ME 10 > 1 {
    GIVE_ME 10!
  }
  GIVE_ME 1!
}
`, 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testMoneyObject(t, evaluated, tt.expected)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{"5 + YUGE!", "type mismatch: MONEY + TRUTH"},
		{"5 + YUGE! 5!", "type mismatch: MONEY + TRUTH"},
		{"-YUGE!", "unknown operator: -TRUTH"},
		{"YUGE + LOSER!", "unknown operator: TRUTH + TRUTH"},
		{"5! YUGE + LOSER! 5!", "unknown operator: TRUTH + TRUTH"},
		{"BELIEVE_ME 10 > 1 { YUGE + LOSER! }", "unknown operator: TRUTH + TRUTH"},
		{"foobar", "identifier not found: foobar"},
		{`"Hello" - "World"`, "unknown operator: TWEET - TWEET"},
		{`{"name": "Covfefe"}[MAKE_IT_BIG(x) { x }]`, "unusable as hash key: FUNCTION"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("no error object returned. got=%T(%+v)", evaluated, evaluated)
			continue
		}
		if errObj.Message != tt.expectedMessage {
			t.Errorf("wrong error message. expected=%q, got=%q", tt.expectedMessage, errObj.Message)
		}
	}
}

func TestLookStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"LOOK a IS 5! a!", 5},
		{"LOOK a IS 5 * 5! a!", 25},
		{"LOOK a IS 5! LOOK b IS a! b!", 5},
		{"LOOK a IS 5! LOOK b IS a! LOOK c IS a + b + 5! c!", 15},
	}

	for _, tt := range tests {
		testMoneyObject(t, testEval(tt.input), tt.expected)
	}
}

func TestFunctionObject(t *testing.T) {
	input := "MAKE_IT_BIG (x) { x + 2; }!"
	evaluated := testEval(input)
	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("object is not Function. got=%T (%+v)", evaluated, evaluated)
	}

	if len(fn.Parameters) != 1 {
		t.Fatalf("function has wrong number of parameters. want=1, got=%d", len(fn.Parameters))
	}

	if fn.Parameters[0].String() != "x" {
		t.Fatalf("parameter is not 'x'. got=%q", fn.Parameters[0])
	}

	if len(fn.Body.Statements) < 1 {
		t.Fatalf("function body has no statements")
	}

	stmt, ok := fn.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("function body statement is not ast.ExpressionStatement. got=%T", fn.Body.Statements[0])
	}

	expectedBody := "(x + 2)"
	if stmt.String() != expectedBody {
		t.Fatalf("body is not %q. got=%q", expectedBody, stmt.String())
	}
}

// TODO: Invalid type err
// func TestFunctionApplication(t *testing.T) {
// 	tests := []struct {
// 		input    string
// 		expected int64
// 	}{
// 		{"MAKE_IT_BIG(x) { x; }(5)!", 5},
// 		{"LOOK identity IS MAKE_IT_BIG(x) { x; }! identity(5)!", 5},
// 		{"LOOK double IS MAKE_IT_BIG(x) { x * 2; }! double(5)!", 10},
// 		{"LOOK add IS MAKE_IT_BIG(x, y) { x + y; }! add(5, 5)!", 10},
// 		{"LOOK identity IS MAKE_IT_BIG(x) { GIVE_ME x; }! identity(5)!", 5},
// 		{"LOOK add IS MAKE_IT_BIG(x, y) { x + y; }! add(5 + 5, add(5, 5))!", 20},
// 	}
// 	for _, tt := range tests {
// 		testMoneyObject(t, testEval(tt.input), tt.expected)
// 	}
// }

// TODO: Closure still doesn't work.
// func TestClosures(t *testing.T) {
// 	input := `
// LOOK newAdder IS MAKE_IT_BIG(x) {
//   MAKE_IT_BIG(y) { x + y };
// };
// LOOK addTwo IS newAdder(2);
// addTwo(2);`
// 	testMoneyObject(t, testEval(input), 4)
// }

func TestForLoopEvaluation(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{
			"LOOK i IS 0! KEEP_WINNING (i < 10) { LOOK i IS i + 1! } i!",
			10,
		},
		{
			"LOOK i IS 0! KEEP_WINNING (i < 10) { LOOK i IS i + 1! BELIEVE_ME i == 5 { IT_WAS_RIGGED! } } i!",
			5,
		},
		{
			"LOOK total IS 0! LOOK i IS 0! KEEP_WINNING (i < 10) { LOOK i IS i + 1! BELIEVE_ME i == 5 { TIRED_OF_WINNING! } LOOK total IS total + i! } total!",
			50,
		},
	}

	for _, tt := range tests {
		testMoneyObject(t, testEval(tt.input), tt.expected)
	}
}

func TestTweetLiteral(t *testing.T) {
	input := `"Hello World!"`
	evaluated := testEval(input)
	str, ok := evaluated.(*object.Tweet)
	if !ok {
		t.Fatalf("object is not Tweet. got=%T (%+v)", evaluated, evaluated)
	}
	if str.Value != "Hello World!" {
		t.Errorf("Tweet has wrong value. got=%q", str.Value)
	}
}

func TestTweetConcatenation(t *testing.T) {
	input := `"Hello" + " " + "World!"`
	evaluated := testEval(input)
	str, ok := evaluated.(*object.Tweet)
	if !ok {
		t.Fatalf("object is not Tweet. got=%T (%+v)", evaluated, evaluated)
	}
	if str.Value != "Hello World!" {
		t.Errorf("Tweet has wrong value. got=%q", str.Value)
	}
}

func TestBuiltinFunctions(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{`len("")!`, 0},
		{`len("four")!`, 4},
		{`len("hello world")!`, 11},
		{`len([1, 2, 3])!`, 3},
		{`len([])!`, 0},
		{`len(1)!`, "argument to `len` not supported, got MONEY"},
		{`len("one", "two")!`, "wrong number of arguments. got=2, want=1"},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		switch expected := tt.expected.(type) {
		case int:
			testMoneyObject(t, evaluated, int64(expected))
		case string:
			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("object is not Error. got=%T (%+v)", evaluated, evaluated)
				continue
			}
			if errObj.Message != expected {
				t.Errorf("wrong error message. expected=%q, got=%q", expected, errObj.Message)
			}
		}
	}
}

func TestArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]!"
	evaluated := testEval(input)
	result, ok := evaluated.(*object.Array)
	if !ok {
		t.Fatalf("object is not Array. got=%T (%+v)", evaluated, evaluated)
	}
	if len(result.Elements) != 3 {
		t.Fatalf("array has wrong num of elements. got=%d", len(result.Elements))
	}
	testMoneyObject(t, result.Elements[0], 1)
	testMoneyObject(t, result.Elements[1], 4)
	testMoneyObject(t, result.Elements[2], 6)
}

func TestArrayIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{"[1, 2, 3][0]!", 1},
		{"[1, 2, 3][1]!", 2},
		{"[1, 2, 3][2]!", 3},
		{"LOOK i IS 0! [1][i]!", 1},
		{"[1, 2, 3][1 + 1]!", 3},
		{"LOOK myArray IS [1, 2, 3]! myArray[2]!", 3},
		{"LOOK myArray IS [1, 2, 3]! myArray[0] + myArray[1] + myArray[2]!", 6},
		{"LOOK myArray IS [1, 2, 3]! LOOK i IS myArray[0]! myArray[i]!", 2},
		{"[1, 2, 3][3]!", nil},
		{"[1, 2, 3][-1]!", nil},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testMoneyObject(t, evaluated, int64(integer))
		} else {
			testCovfefeObject(t, evaluated)
		}
	}
}

func TestHashLiterals(t *testing.T) {
	input := `LOOK two IS "two"!
	{
		"one": 10 - 9,
		two: 1 + 1,
		"thr" + "ee": 6 / 2,
		4: 4,
		YUGE: 5,
		LOSER: 6
	}!`
	evaluated := testEval(input)
	result, ok := evaluated.(*object.Hash)
	if !ok {
		t.Fatalf("Eval didn't return Hash. got=%T (%+v)", evaluated, evaluated)
	}
	expected := map[object.HashKey]int64{
		(&object.Tweet{Value: "one"}).HashKey():   1,
		(&object.Tweet{Value: "two"}).HashKey():   2,
		(&object.Tweet{Value: "three"}).HashKey(): 3,
		(&object.Money{Value: 4}).HashKey():       4,
		object.YUGE.HashKey():                     5,
		object.LOSER.HashKey():                    6,
	}
	if len(result.Pairs) != len(expected) {
		t.Fatalf("Hash has wrong num of pairs. got=%d", len(result.Pairs))
	}
	for expectedKey, expectedValue := range expected {
		pair, ok := result.Pairs[expectedKey]
		if !ok {
			t.Errorf("no pair for given key in Pairs")
		}
		testMoneyObject(t, pair.Value, expectedValue)
	}
}

func TestHashIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{`{"foo": 5}["foo"]!`, 5},
		{`{"foo": 5}["bar"]!`, nil},
		{`LOOK key IS "foo"! {"foo": 5}[key]!`, 5},
		{`{}["foo"]!`, nil},
		{`{5: 5}[5]!`, 5},
		{`{YUGE: 5}[YUGE]!`, 5},
		{`{LOSER: 5}[LOSER]!`, 5},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testMoneyObject(t, evaluated, int64(integer))
		} else {
			testCovfefeObject(t, evaluated)
		}
	}
}

// --- Helper Functions ---

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()

	return Eval(program, env)
}

func testMoneyObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Money)
	if !ok {
		t.Errorf("object is not *object.Money. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d", result.Value, expected)
		return false
	}
	return true
}

func testTruthObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Truth)
	if !ok {
		t.Errorf("object is not *object.Truth. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%t, want=%t", result.Value, expected)
		return false
	}
	return true
}

func testCovfefeObject(t *testing.T, obj object.Object) bool {
	if obj != object.COVFEFE {
		t.Errorf("object is not object.COVFEFE. got=%T (%+v)", obj, obj)
		return false
	}
	return true
}
