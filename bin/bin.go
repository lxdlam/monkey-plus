package bin

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/lxdlam/monkey-plus/evaluator"
	"github.com/lxdlam/monkey-plus/lexer"
	"github.com/lxdlam/monkey-plus/object"
	"github.com/lxdlam/monkey-plus/parser"
	"io"
	"log"
	"os"
	"strings"
)

func Run(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	var codes bytes.Buffer

	for {
		scanned := scanner.Scan()

		if !scanned {
			break
		}
		codes.WriteString(scanner.Text() + "\n")
	}

	l := lexer.New(codes.String())
	p := parser.New(l)

	program := p.ParseProgram()

	errorLength := len(p.Errors())
	if errorLength != 0 {
		var errors []string
		for idx, msg := range p.Errors() {
			errors = append(errors, fmt.Sprintf(" parse error[%d/%d]: %s", idx+1, errorLength, msg))
		}
		_, err := io.WriteString(out, strings.Join(errors, "\n"))
		if err != nil {
			log.Fatalf(err.Error())
		}
		return
	}

	evaluated := evaluator.Eval(program, env)
	if evaluated != nil && evaluated.Inspect() != "null" {
		_, err := io.WriteString(out, evaluated.Inspect())
		if err != nil {
			log.Fatalf(err.Error())
		}
	}
}

func RunFile(path string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf(err.Error())
	}

	defer file.Close()

	Run(file, os.Stdout)
}

func RunCode(code string) {
	Run(strings.NewReader(code), os.Stdout)
}
