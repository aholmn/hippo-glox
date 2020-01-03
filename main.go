package main

import (
	"errors"
	"fmt"
	"hippo-glox/interpreter"
	"hippo-glox/parser"
	"hippo-glox/scanner"
	"io/ioutil"
	"os"
)

func main() {
	bytes, err := readFile()
	if err != nil {
		fmt.Println(err)
	}

	pErrors, stmts := parser.Parse(scanner.Scan(string(bytes)))

	if len(pErrors) > 0 {
		for _, e := range pErrors {
			fmt.Println(e)
		}
		return
	}

	interpreter.Interpret(stmts)
}

func readFile() ([]byte, error) {
	if len(os.Args) < 2 {
		return nil, errors.New("No file provided")
	}
	return ioutil.ReadFile(os.Args[1])
}
