package main

import (
	"fmt"
	"reflect"
	"testing"
)

func runTokenMatch(t *testing.T, got, wanted []Token) {
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

type test struct {
	input string
	want  []Token
}

func runTestingTable(t *testing.T, tests []test, format string) {
	t.Helper()
	for _, test := range tests {
		t.Run(fmt.Sprintf(format, test.input), func(t *testing.T) {
			got := Tokenize(test.input)
			runTokenMatch(t, got, test.want)
		})
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
	runTokenMatch(t, got, wanted)
}

func TestTokenizeSceneHeading(t *testing.T) {
	tests := []test{
		{"EXT. BRICK'S POOL - DAY", []Token{{"scene_heading", "EXT. BRICK'S POOL - DAY"}}},
		{"INT. HOUSE - DAY", []Token{{"scene_heading", "INT. HOUSE - DAY"}}},
		{".SNIPER SCOPE POV", []Token{{"scene_heading", ".SNIPER SCOPE POV"}}},
	}
	runTestingTable(t, tests, "SceneHeading %s")
}

// Leading ellipses shouldn't be interpreted as a scene heading.
func TestEllipsesNotSceneHeading(t *testing.T) {
	input := "...foo bar"
	want := []Token{{"text", "...foo bar"}}
	got := Tokenize(input)
	runTokenMatch(t, got, want)
}

func TestCharacter(t *testing.T) {
	tests := []test{
		{"STEEL", []Token{{"character", "STEEL"}}},
		{"@McCLANE", []Token{{"character", "McCLANE"}}},
	}
	runTestingTable(t, tests, "CharacterDialogue %s")
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
	runTokenMatch(t, got, want)
}

func TestParens(t *testing.T) {
	input := "(foo bar)"
	want := []Token{
		{"oparen", "("},
		{"text", "foo bar"},
		{"cparen", ")"},
	}
	got := Tokenize(input)
	runTokenMatch(t, got, want)
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
	runTokenMatch(t, got, want)
}

func TestLyric(t *testing.T) {
	input := "~Willy Wonka! Willy Wonka! The amazing chocolatier!"
	want := []Token{
		{"tilde", "~"},
		{"text", "Willy Wonka! Willy Wonka! The amazing chocolatier!"},
	}
	got := Tokenize(input)
	runTokenMatch(t, got, want)
}

func TestTransition(t *testing.T) {
	tests := []test{
		{"CUT TO:", []Token{{"transition", "CUT TO:"}}},
		{"FADE TO:", []Token{{"transition", "FADE TO:"}}},
		{"ENTER TO:", []Token{{"transition", "ENTER TO:"}}},
		{"> Burn to white.", []Token{{"transition", "Burn to white."}}},
	}
	runTestingTable(t, tests, "Transition %s")
}

func TestCenteredText(t *testing.T) {
	tests := []test{
		{">THE END<", []Token{{"centered_text", "THE END"}}},
		{"> THE END <", []Token{{"centered_text", "THE END"}}},
	}
	runTestingTable(t, tests, "CenteredText %s")
}

func TestEmphasis(t *testing.T) {
	tests := []test{
		{"*italics*", []Token{
			{"asterisk", "*"},
			{"text", "italics"},
			{"asterisk", "*"},
		}},
		{"**bold**", []Token{
			{"asterisk", "*"},
			{"asterisk", "*"},
			{"text", "bold"},
			{"asterisk", "*"},
			{"asterisk", "*"},
		}},
		{"***bold italics***", []Token{
			{"asterisk", "*"},
			{"asterisk", "*"},
			{"asterisk", "*"},
			{"text", "bold italics"},
			{"asterisk", "*"},
			{"asterisk", "*"},
			{"asterisk", "*"},
		}},
		{"_underline_", []Token{
			{"underscore", "_"},
			{"text", "underline"},
			{"underscore", "_"},
		}},
	}
	runTestingTable(t, tests, "Emphasis %s")
}

func TestNesting(t *testing.T) {
	input := "From what seems like only INCHES AWAY. _Steel's face FILLS the *Leupold Mark 4* scope_."
	want := []Token{
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
	got := Tokenize(input)
	runTokenMatch(t, got, want)
}

func TestPageBreak(t *testing.T) {
	tests := []test{
		{"===", []Token{{"page_break", "==="}}},
		{"==", []Token{{"text", "=="}}},
		{"====", []Token{{"text", "===="}}},
	}
	runTestingTable(t, tests, "PageBreak %s")
}

func TestNotes(t *testing.T) {
	tests := []test{
		{"[[This section needs work.]]", []Token{
			{"onote", "[["},
			{"text", "This section needs work."},
			{"cnote", "]]"},
		}},
		{"Foo[[bar]]", []Token{
			{"text", "Foo"},
			{"onote", "[["},
			{"text", "bar"},
			{"cnote", "]]"},
		}},
		{"Foo[bar]", []Token{
			{"text", "Foo"},
			{"obrace", "["},
			{"text", "bar"},
			{"cbrace", "]"},
		}},
	}
	runTestingTable(t, tests, "Note %s")
}

func TestBoneyard(t *testing.T) {
	tests := []test{
		{"/*This section needs work.*/", []Token{
			{"oboneyard", "/*"},
			{"text", "This section needs work."},
			{"cboneyard", "*/"},
		}},
		{"Foo/*bar*/", []Token{
			{"text", "Foo"},
			{"oboneyard", "/*"},
			{"text", "bar"},
			{"cboneyard", "*/"},
		}},
		{"Foo/bar/", []Token{
			{"text", "Foo"},
			{"forward_slash", "/"},
			{"text", "bar"},
			{"forward_slash", "/"},
		}},
	}
	runTestingTable(t, tests, "Boneyard %s")
}

func TestSection(t *testing.T) {
	tests := []test{
		{"# Act", []Token{
			{"h1", "Act"},
		}},
		{"## Sequence", []Token{
			{"h2", "Sequence"},
		}},
		{"### Scene", []Token{
			{"h3", "Scene"},
		}},
	}
	runTestingTable(t, tests, "Section %s")
}

const fixture = `INT.  WILL'S BEDROOM - NIGHT (1973)

WILL BLOOM, AGE 3, listens wide-eyed as his father EDWARD BLOOM, 40's and handsome, tells the story.  In every gesture, Edward is bigger than life, describing each detail with absolute conviction.

EDWARD
I didn't put any stock into such speculation or superstition.  All I knew was I'd been trying to catch that fish since I was a boy no bigger than you.
(closer)
And on the day you were born, that was the day I finally caught him.

EXT.  CAMPFIRE - NIGHT (1977)

A few years later, and Will sits with the other INDIAN GUIDES as Edward continues telling the story to the tribe.

EDWARD (V.O)
Now, I'd tried everything on it:  worms, lures, peanut butter, peanut butter-and-cheese.  But on that day I had a revelation:  if that fish was the ghost of a thief, the usual bait wasn't going to work.  I would have to use something he truly desired.

Edward points to his wedding band, glinting in the firelight.

LITTLE BRAVE
(confused)
Your finger?

Edward slips his ring off.`

func TestTokenizeFixture(t *testing.T) {
	want := []Token{
		{"scene_heading", "INT.  WILL'S BEDROOM - NIGHT (1973)"},
		{"newline", ""},
		{"newline", ""},
		{"text", "WILL BLOOM, AGE 3, listens wide-eyed as his father EDWARD BLOOM, 40's and handsome, tells the story.  In every gesture, Edward is bigger than life, describing each detail with absolute conviction."},
		{"newline", ""},
		{"newline", ""},
		{"character", "EDWARD"},
		{"newline", ""},
		{"text", "I didn't put any stock into such speculation or superstition.  All I knew was I'd been trying to catch that fish since I was a boy no bigger than you."},
		{"newline", ""},
		{"oparen", "("},
		{"text", "closer"},
		{"cparen", ")"},
		{"newline", ""},
		{"text", "And on the day you were born, that was the day I finally caught him."},
		{"newline", ""},
		{"newline", ""},
		{"scene_heading", "EXT.  CAMPFIRE - NIGHT (1977)"},
		{"newline", ""},
		{"newline", ""},
		{"text", "A few years later, and Will sits with the other INDIAN GUIDES as Edward continues telling the story to the tribe."},
		{"newline", ""},
		{"newline", ""},
		{"character", "EDWARD"},
		{"oparen", "("},
		{"character", "V.O"}, // TODO: I wonder if we should allow nested expressions in characters
		{"cparen", ")"},
		{"newline", ""},
		{"text", "Now, I'd tried everything on it:  worms, lures, peanut butter, peanut butter-and-cheese.  But on that day I had a revelation:  if that fish was the ghost of a thief, the usual bait wasn't going to work.  I would have to use something he truly desired."},
		{"newline", ""},
		{"newline", ""},
		{"text", "Edward points to his wedding band, glinting in the firelight."},
		{"newline", ""},
		{"newline", ""},
		{"character", "LITTLE BRAVE"},
		{"newline", ""},
		{"oparen", "("},
		{"text", "confused"},
		{"cparen", ")"},
		{"newline", ""},
		{"text", "Your finger?"},
		{"newline", ""},
		{"newline", ""},
		{"text", "Edward slips his ring off."},
	}

	got := Tokenize(fixture)
	runTokenMatch(t, got, want)
}

func BenchmarkTokenizer(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Tokenize(fixture)
	}
}
