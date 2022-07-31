package main

import (
	"testing"
	"reflect"
)

func TestWalkNesting(t *testing.T) {
	input := []Token{
		{"text", "From what seems like only INCHES AWAY. "},
		{"underscore", "_"},
		{"text", "Steel's face FILLS the "},
		{"asterisk", "*"},
		{"text", "Leupold Mark 4"},
		{"asterisk", "*"},
		{"text", " scope"},
		{"underscore", "_"},
		{"text", "."},
	}
	got := Parse(input)

	wanted := node{}

	if !reflect.DeepEqual(got, wanted) {
		t.Error("\nExpected:", wanted, "\nGot:", got)
	}
}
