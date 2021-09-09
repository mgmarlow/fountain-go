package main

import (
	"fmt"
	"io/ioutil"
)

func main() {
	file, err := ioutil.ReadFile("./fixtures/brick&steel.fountain")
	if err != nil {
		panic(err)
	}

	fileContents := string(file)
	l := NewLexer(fileContents)
	for l.token != TEndOfFile {
		fmt.Printf("%s\n", l)
		l.Next()
	}
}
