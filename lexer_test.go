package main

import (
	"fmt"
	"reflect"
	"testing"
)

func testTokenMatch(t *testing.T, got, wanted []Token) {
	t.Helper()

	// Ending newline assumed so we don't need to repeat it in the tests.
	wanted = append(wanted, Token{
		kind:  "newline",
		value: "",
	})

	if !reflect.DeepEqual(got, wanted) {
		t.Error("\nExpected:", wanted, "\nGot:", got)
	}
}

func TestActionWithNewlines(t *testing.T) {
	input := `They drink long and well from the beers.

And then there's a long beat.`
	got := Tokenize(input)
	wanted := []Token{
		{"action", "They drink long and well from the beers."},
		{"newline", ""},
		{"newline", ""},
		{"action", "And then there's a long beat."},
	}
	testTokenMatch(t, got, wanted)
}

func TestTokenizeSceneHeading(t *testing.T) {
	tests := []struct {
		input  string
		wanted []Token
	}{
		{"EXT. BRICK'S POOL - DAY", []Token{{"scene_heading", "EXT. BRICK'S POOL - DAY"}}},
		{"INT. HOUSE - DAY", []Token{{"scene_heading", "INT. HOUSE - DAY"}}},
		{".SNIPER SCOPE POV", []Token{{"scene_heading", ".SNIPER SCOPE POV"}}},
		// Not interpreted as a scene heading
		{"...foo bar", []Token{{"action", "...foo bar"}}},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Scene Heading %s", test.input), func(t *testing.T) {
			got := Tokenize(test.input)
			testTokenMatch(t, got, test.wanted)
		})
	}
}
