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
		{"@McCLANE", []Token{{"character", "McCLANE"}}},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Character Dialogue %s", test.input), func(t *testing.T) {
			got := Tokenize(test.input)
			testTokenMatch(t, got, test.wanted)
		})
	}
}

func TestCharacterWithParenthetical(t *testing.T) {
	input := "HANS (on the radio)"
	want := []Token{
		{"character", "HANS"},
		{"oparen", "("},
		{"text", "on the radio"},
		{"cparen", ")"},
	}
	got := Tokenize(input)
	testTokenMatch(t, got, want)
}

func TestParens(t *testing.T) {
	input := "(foo bar)"
	want := []Token{
		{"oparen", "("},
		{"text", "foo bar"},
		{"cparen", ")"},
	}
	got := Tokenize(input)
	testTokenMatch(t, got, want)
}

func TestDualDialogue(t *testing.T) {
	input := `BRICK
Screw retirement.

STEEL ^
Screw retirement.`
	want := []Token{
		{"character", "BRICK"},
		{"newline", ""},
		{"text", "Screw retirement."},
		{"newline", ""},
		{"newline", ""},
		{"character", "STEEL"},
		{"caret", "^"},
		{"newline", ""},
		{"text", "Screw retirement."},
	}
	got := Tokenize(input)
	testTokenMatch(t, got, want)
}

func TestLyric(t *testing.T) {
	input := "~Willy Wonka! Willy Wonka! The amazing chocolatier!"
	want := []Token{
		{"tilde", "~"},
		{"text", "Willy Wonka! Willy Wonka! The amazing chocolatier!"},
	}
	got := Tokenize(input)
	testTokenMatch(t, got, want)
}

func TestTransition(t *testing.T) {
	tests := []struct {
		input  string
		wanted []Token
	}{
		{"CUT TO:", []Token{{"transition", "CUT TO:"}}},
		{"FADE TO:", []Token{{"transition", "FADE TO:"}}},
		{"ENTER TO:", []Token{{"transition", "ENTER TO:"}}},
		{"> Burn to white.", []Token{{"transition", "Burn to white."}}},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Transition %s", test.input), func(t *testing.T) {
			got := Tokenize(test.input)
			testTokenMatch(t, got, test.wanted)
		})
	}
}
