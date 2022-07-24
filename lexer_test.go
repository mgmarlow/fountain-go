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

func TestTextWithNewlines(t *testing.T) {
	input := `They drink long and well from the beers.

And then there's a long beat.`
	got := Tokenize(input)
	wanted := []Token{
		{"text", "They drink long and well from the beers."},
		{"newline", ""},
		{"newline", ""},
		{"text", "And then there's a long beat."},
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
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Scene Heading %s", test.input), func(t *testing.T) {
			got := Tokenize(test.input)
			testTokenMatch(t, got, test.wanted)
		})
	}
}

// Leading ellipses shouldn't be interpreted as a scene heading.
func TestEllipsesNotSceneHeading(t *testing.T) {
	input := "...foo bar"
	want := []Token{{"text", "...foo bar"}}
	got := Tokenize(input)
	testTokenMatch(t, got, want)
}

func TestCharacter(t *testing.T) {
	tests := []struct {
		input  string
		wanted []Token
	}{
		{"STEEL", []Token{{"character", "STEEL"}}},
		// TODO:
		// {"HANS (on the radio)", []Token{{"character", "HANS (on the radio)"}}},
		{"@McCLANE", []Token{{"character", "McCLANE"}}},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Character Dialogue %s", test.input), func(t *testing.T) {
			got := Tokenize(test.input)
			testTokenMatch(t, got, test.wanted)
		})
	}
}
