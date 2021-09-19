package main

import (
	"io/ioutil"

	"github.com/mgmarlow/fountain/emitters"
	"github.com/mgmarlow/fountain/parser"
)

func main() {
	file, err := ioutil.ReadFile("./fixtures/brick&steel.fountain")
	if err != nil {
		panic(err)
	}

	fileContents := string(file)
	p := parser.NewParser(fileContents)
	j := p.Emit(emitters.NewJsonEmitter())
	_ = ioutil.WriteFile("test.json", j, 0644)
}
