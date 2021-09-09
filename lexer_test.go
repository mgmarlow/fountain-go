package main

import (
	"reflect"
	"testing"
)

func runTest(contents string) []string {
	result := []string{}

	l := NewLexer(contents)
	for l.token != TEndOfFile {
		result = append(result, l.String())
		l.Next()
	}

	return result
}

// Headings
func TestHeadingExt(t *testing.T) {
	fixture := "EXT. FOO BAR\n\nsome action text"
	got := runTest(fixture)
	expected := []string{
		"Heading: EXT. FOO BAR",
		"Text: some action text",
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("expected %v but got %v", expected, got)
	}
}

func TestHeadingInt(t *testing.T) {
	fixture := "INT. FOO BAR\n\nsome action text"
	got := runTest(fixture)
	expected := []string{
		"Heading: INT. FOO BAR",
		"Text: some action text",
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("expected %v but got %v", expected, got)
	}
}

func TestHeadingCustom(t *testing.T) {
	fixture := ".FOO BAR\n\nsome action text"
	got := runTest(fixture)
	expected := []string{
		"Heading: .FOO BAR",
		"Text: some action text",
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("expected %v but got %v", expected, got)
	}
}

// Characters
func TestCharacter(t *testing.T) {
	fixture := "BRIAN\nsome action text"
	got := runTest(fixture)
	expected := []string{
		"Character: BRIAN",
		"Text: some action text",
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("expected %v but got %v", expected, got)
	}
}

func TestCharacterFirstLast(t *testing.T) {
	fixture := "BRIAN MCDONALD\nsome action text"
	got := runTest(fixture)
	expected := []string{
		"Character: BRIAN MCDONALD",
		"Text: some action text",
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("expected %v but got %v", expected, got)
	}
}

// Pagebreak
func TestPagebreak(t *testing.T) {
	fixture := "some action text\n\n===\n\n"
	got := runTest(fixture)
	expected := []string{
		"Text: some action text",
		"Pagebreak: ===",
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("expected %v but got %v", expected, got)
	}
}

func TestPagebreakWithNoEndingNewlines(t *testing.T) {
	fixture := "some action text\n\n==="
	got := runTest(fixture)
	expected := []string{
		"Text: some action text",
		"Pagebreak: ===",
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("expected %v but got %v", expected, got)
	}
}

// Centers
func TestCenteredText(t *testing.T) {
	fixture := "> THE END <"
	got := runTest(fixture)
	expected := []string{
		"Centered Action: THE END",
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("expected %v but got %v", expected, got)
	}
}

// Transitions
func TestTransition(t *testing.T) {
	fixture := "CUT TO:"
	got := runTest(fixture)
	expected := []string{
		"Transition: CUT TO:",
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("expected %v but got %v", expected, got)
	}
}

func TestTransitionCustom(t *testing.T) {
	fixture := "> BURN TO PINK."
	got := runTest(fixture)
	expected := []string{
		"Transition: BURN TO PINK",
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("expected %v but got %v", expected, got)
	}
}

// Parentheticals
func TestParenthetical(t *testing.T) {
	fixture := "(starting the engine)"
	got := runTest(fixture)
	expected := []string{
		"Parenthetical: (starting the engine)",
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("expected %v but got %v", expected, got)
	}
}

func TestParentheticalWithCharacter(t *testing.T) {
	fixture := "BRIAN\n(starting the engine)"
	got := runTest(fixture)
	expected := []string{
		"Character: BRIAN",
		"Parenthetical: (starting the engine)",
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("expected %v but got %v", expected, got)
	}
}

func TestParentheticalWithCharacterAndAction(t *testing.T) {
	fixture := "BRIAN\n(starting the engine)\nsome action"
	got := runTest(fixture)
	expected := []string{
		"Character: BRIAN",
		"Parenthetical: (starting the engine)",
		"Text: some action",
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("expected %v but got %v", expected, got)
	}
}

// Action text
func TestAction(t *testing.T) {
	fixture := "some action"
	got := runTest(fixture)
	expected := []string{
		"Text: some action",
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("expected %v but got %v", expected, got)
	}
}

func TestActionWithTabs(t *testing.T) {
	fixture := "\tsome action"
	got := runTest(fixture)
	expected := []string{
		"Text:     some action",
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("expected %v but got %v", expected, got)
	}
}

// TODO: This doesn't follow the spec.
// https://fountain.io/syntax#section-action
func TestActionWithWhitespace(t *testing.T) {
	fixture := "   some action   "
	got := runTest(fixture)
	expected := []string{
		"Text:    some action   ",
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("expected %v but got %v", expected, got)
	}
}

func TestUnderline(t *testing.T) {
	fixture := "_some action_"
	got := runTest(fixture)
	expected := []string{
		"Underscore: _",
		"Text: some action",
		"Underscore: _",
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("expected %v but got %v", expected, got)
	}
}

func TestItalics(t *testing.T) {
	fixture := "*some action*"
	got := runTest(fixture)
	expected := []string{
		"Asterisk: *",
		"Text: some action",
		"Asterisk: *",
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("expected %v but got %v", expected, got)
	}
}

func TestBold(t *testing.T) {
	fixture := "**some action**"
	got := runTest(fixture)
	expected := []string{
		"Asterisk: *",
		"Asterisk: *",
		"Text: some action",
		"Asterisk: *",
		"Asterisk: *",
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("expected %v but got %v", expected, got)
	}
}

func TestUnderlinedBoldItalic(t *testing.T) {
	fixture := "_***some action***_"
	got := runTest(fixture)
	expected := []string{
		"Underscore: _",
		"Asterisk: *",
		"Asterisk: *",
		"Asterisk: *",
		"Text: some action",
		"Asterisk: *",
		"Asterisk: *",
		"Asterisk: *",
		"Underscore: _",
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("expected %v but got %v", expected, got)
	}
}

func TestMultilineUnderlinedBoldItalic(t *testing.T) {
	fixture := "_***some\naction***_"
	got := runTest(fixture)
	expected := []string{
		"Underscore: _",
		"Asterisk: *",
		"Asterisk: *",
		"Asterisk: *",
		"Text: some",
		"Text: action",
		"Asterisk: *",
		"Asterisk: *",
		"Asterisk: *",
		"Underscore: _",
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("expected %v but got %v", expected, got)
	}
}

// Boneyard
func TestBoneyard(t *testing.T) {
	fixture := "/* some action */"
	got := runTest(fixture)
	expected := []string{
		"Boneyard Open: /*",
		"Text: some action",
		"Boneyard End: */",
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("expected %v but got %v", expected, got)
	}
}

func TestMultilineBoneyard(t *testing.T) {
	fixture := "/* some\naction */"
	got := runTest(fixture)
	expected := []string{
		"Boneyard Open: /*",
		"Text: some",
		"Text: action",
		"Boneyard End: */",
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("expected %v but got %v", expected, got)
	}
}

// Notes
func TestNotes(t *testing.T) {
	fixture := "[[ some action ]]"
	got := runTest(fixture)
	expected := []string{
		"Note Open: [[",
		"Text: some action",
		"Note End: ]]",
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("expected %v but got %v", expected, got)
	}
}

func TestMultilineNotes(t *testing.T) {
	fixture := "[[ some\naction ]]"
	got := runTest(fixture)
	expected := []string{
		"Note Open: [[",
		"Text: some",
		"Text: action",
		"Note End: ]]",
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("expected %v but got %v", expected, got)
	}
}
