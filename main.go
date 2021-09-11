package main

import (
	"fmt"
	"io/ioutil"

	"github.com/mgmarlow/fountain/lexer"
)

func main() {
	file, err := ioutil.ReadFile("./fixtures/brick&steel.fountain")
	if err != nil {
		panic(err)
	}

	fileContents := string(file)
	l := lexer.NewLexer(fileContents)
	for l.Token != lexer.TEndOfFile {
		fmt.Println(l)
		l.Next()
	}
}
