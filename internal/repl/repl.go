package repl

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/peterh/liner"
	"github.com/yurikdotdev/covfefescript/internal/eval"
	"github.com/yurikdotdev/covfefescript/internal/lexer"
	"github.com/yurikdotdev/covfefescript/internal/object"
	"github.com/yurikdotdev/covfefescript/internal/parser"
)

func InitREPL() {
	line := liner.NewLiner()
	defer line.Close()

	if len(os.Args) > 1 {
		RunFile(os.Args[1])
	} else {
		fmt.Print(TRUMP_ASCII_INTRO)
		fmt.Println("\nCovfefeScript 0.1 | The Best Interpreter, Believe Me.")
		fmt.Print("Type 'CHYNA' to quit.\n\n")
		StartREPL(line, os.Stdout)
	}
}

func RunFile(filename string) {
	if !strings.HasSuffix(filename, ".covfefe") {
		fmt.Println("SAD! This is not a .covfefe file. Very sad.")
		return
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("SAD! Couldn't read file %s: %s\n", filename, err)
		return
	}

	l := lexer.New(string(data))
	p := parser.New(l)
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		printParserErrors(os.Stdout, p.Errors())
		return
	}

	env := object.NewEnvironment()
	evaluated := eval.Eval(program, env)

	// fmt.Fprintf(os.Stderr, "DEBUG: evaluated = %v, type = %s\n", evaluated, evaluated.Type())

	if evaluated != nil {
		if evaluated.Type() == object.ERROR_OBJ {
			fmt.Fprintln(os.Stderr, evaluated.Inspect())
			os.Exit(1)
		}
	}
}

func StartREPL(line *liner.State, out io.Writer) {
	env := object.NewEnvironment()

	for {
		input, err := line.Prompt(REPL_PROMPT_SIGN)
		if err != nil {
			if err == io.EOF {
				fmt.Print(TRUMP_ASCII_OUTRO)
			} else {
				fmt.Printf("SAD! A terrible error reading input: %s\n", err)
			}
			return
		}

		if input == "" {
			continue
		}

		line.AppendHistory(input)

		if input == "CHYNA" {
			fmt.Print(TRUMP_ASCII_OUTRO)
			return
		}

		// HE WHO SHALL NOT BE NAMED
		target := string([]rune{101, 112, 115, 116, 101, 105, 110})
		if strings.ToLower(input) == target {
			fmt.Print("Really bro?")
			return
		}

		l := lexer.New(input)
		p := parser.New(l)
		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		evaluated := eval.Eval(program, env)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, TRUMP_ASCII_LOSER)
	io.WriteString(out, "SAD! Unbelievable, what a loser:\n")
	for _, msg := range errors {
		io.WriteString(out, ""+msg+"\n\n")
	}
}
