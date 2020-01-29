package repl

import (
	"bufio"
	"fmt"
	"github.com/lxdlam/monkey-plus/evaluator"
	"github.com/lxdlam/monkey-plus/lexer"
	"github.com/lxdlam/monkey-plus/object"
	"github.com/lxdlam/monkey-plus/parser"
	"io"
	"log"
)

const PROMPT = ">> "
const MONKEY_FACE = `            __,__
   .--.  .-"     "-.  .--.
  / .. \/  .-. .-.  \/ .. \
 | |  '|  /   Y   \  |'  | |
 | \   \  \ 0 | 0 /  /   / |
  \ '- ,\.-"""""""-./, -' /
   ''-' /_   ^ ^   _\ '-''
       |  \._   _./  |
       \   \ '~' /   /
        '._ '-=-' _.'
           '-----'
`

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()

		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			_, err := io.WriteString(out, evaluated.Inspect())
			if err != nil {
				log.Fatalf(err.Error())
			}

			_, err = io.WriteString(out, "\n")
			if err != nil {
				log.Fatalf(err.Error())
			}
		}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		_, err := io.WriteString(out, MONKEY_FACE)
		if err != nil {
			log.Fatalf(err.Error())
		}

		_, err = io.WriteString(out, "Woops! We ran into some monkey business here!\n")
		if err != nil {
			log.Fatalf(err.Error())
		}

		_, err = io.WriteString(out, "\t"+msg+"\n")
		if err != nil {
			log.Fatalf(err.Error())
		}
	}
}
