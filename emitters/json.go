package emitters

import (
	"encoding/json"

	"github.com/mgmarlow/fountain/parser"
)

type JsonEmitter struct{}

func NewJsonEmitter() JsonEmitter {
	return JsonEmitter{}
}

func (emitter JsonEmitter) Emit(root parser.Composite) []byte {
	j, err := json.Marshal(root)

	if err != nil {
		// TODO: Better error handling
		panic(err)
	}

	return j
}
