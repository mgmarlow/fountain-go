package lexer

import (
	"io/ioutil"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func getTokens(contents string) []string {
	result := []string{}

	l := NewLexer(contents)
	for l.Token != TEndOfFile {
		result = append(result, l.String())
		l.Next()
	}

	return result
}

func runTest(t *testing.T, text string, expected []string) {
	t.Helper()
	got := getTokens(text)

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("expected %v but got %v", expected, got)
	}
}

func TestAgainstFixture(t *testing.T) {
	fixtureFile := filepath.Join("../", "fixtures", "brick&steel.fountain")
	file, err := ioutil.ReadFile(fixtureFile)
	if err != nil {
		t.Errorf("unable to read fixture")
	}

	got := []string{}
	fileContents := string(file)
	l := NewLexer(fileContents)
	for l.Token != TEndOfFile {
		got = append(got, l.String())
		l.Next()
	}

	raw := `<slugline value="EXT. BRICK'S PATIO - DAY">
<text value="A gorgeous day.  The sun is shining.  But BRICK BRADDOCK, retired police detective, is sitting quietly, contemplating -- something.">
<text value="The SCREEN DOOR slides open and DICK STEEL, his former partner and fellow retiree, emerges with two cold beers.">
<dialogue value="STEEL">
<text value="Beer's ready!">
<dialogue value="BRICK">
<text value="Are they cold?">
<dialogue value="STEEL">
<text value="Does a bear crap in the woods?">
<text value="Steel sits.  They laugh at the dumb joke.">
<dialogue value="STEEL">
<paren_open>
<text value="beer raised">
<paren_close>
<text value="To retirement.">
<dialogue value="BRICK">
<text value="To retirement.">
<text value="They drink long and well from the beers.">
<text value="And then there's a long beat.">
<text value="Longer than is funny.">
<text value="Long enough to be depressing.">
<text value="The men look at each other.">
<dialogue value="STEEL">
<text value="Screw retirement.">
<dialogue value="BRICK ^">
<text value="Screw retirement.">
<transition value="SMASH CUT TO:">
<slugline value="INT. TRAILER HOME - DAY">
<text value="This is the home of THE BOY BAND, AKA DAN and JACK.  They too are drinking beer, and counting the take from their last smash-and-grab.  Money, drugs, and ridiculous props are strewn about the table.">
<dialogue value="JACK">
<paren_open>
<text value="in Vietnamese, subtitled">
<paren_close>
<asterisk>
<text value="Did you know Brick and Steel are retired?">
<asterisk>
<dialogue value="DAN">
<text value="Then let's retire them.">
<underscore>
<text value="Permanently">
<underscore>
<slugline value=".">
<text value="Jack begins to argue vociferously in Vietnamese ">
<paren_open>
<dialogue value="?">
<paren_close>
<text value=", But mercifully we...">
<transition value="CUT TO:">
<slugline value="EXT. BRICK'S POOL - DAY">
<text value="Steel, in the middle of a heated phone call:">
<dialogue value="STEEL">
<text value="They're coming out of the woodwork!">
<paren_open>
<text value="pause">
<paren_close>
<text value="No, everybody we've put away!">
<paren_open>
<text value="pause">
<paren_close>
<text value="Point Blank Sniper?">
<slugline value=".SNIPER SCOPE POV">
<text value="From what seems like only INCHES AWAY.  ">
<underscore>
<text value="Steel's face FILLS the ">
<asterisk>
<text value="Leupold Mark 4">
<asterisk>
<text value=" scope">
<underscore>
<slugline value=".">
<dialogue value="STEEL">
<text value="The man's a myth!">
<text value="Steel turns and looks straight into the cross-hairs.">
<dialogue value="STEEL">
<paren_open>
<text value="oh crap">
<paren_close>
<text value="Hello...">
<transition value="CUT TO:">
<slugline value=".OPENING TITLES">
<centered_text value="> BRICK BRADDOCK <">
<centered_text value="> & DICK STEEL IN <">
<centered_text value="> BRICK & STEEL <">
<centered_text value="> FULL RETIRED <">
<transition value="SMASH CUT TO:">
<slugline value="EXT. WOODEN SHACK - DAY">
<text value="COGNITO, the criminal mastermind, is SLAMMED against the wall.">
<dialogue value="COGNITO">
<text value="Woah woah woah, Brick and Steel!">
<text value="Sure enough, it's Brick and Steel, roughing up their favorite usual suspect.">
<dialogue value="COGNITO">
<text value="What is it you want with me, DICK?">
<text value="Steel SMACKS him.">
<dialogue value="STEEL">
<text value="Who's coming after us?">
<dialogue value="COGNITO">
<text value="Everyone's coming after you mate!  Scorpio, The Boy Band, Sparrow, Point Blank Sniper...">
<text value="As he rattles off the long list, Brick and Steel share a look.  This is going to be BAD.">
<transition value="CUT TO:">
<slugline value="INT. GARAGE - DAY">
<text value="BRICK and STEEL get into Mom's PORSCHE, Steel at the wheel.  They pause for a beat, the gravity of the situation catching up with them.">
<dialogue value="BRICK">
<text value="This is everybody we've ever put away.">
<dialogue value="STEEL">
<paren_open>
<text value="starting the engine">
<paren_close>
<text value="So much for retirement!">
<text value="They speed off.  To destiny!">
<transition value="CUT TO:">
<slugline value="EXT. PALATIAL MANSION - DAY">
<text value="An EXTREMELY HANDSOME MAN drinks a beer.  Shirtless, unfortunately.">
<text value="His minion approaches offscreen:">
<dialogue value="MINION">
<text value="We found Brick and Steel!">
<dialogue value="HANDSOME MAN">
<text value="I want them dead.  DEAD!">
<text value="Beer flies.">
<transition value="> BURN TO PINK.">
<centered_text value="> THE END <">`

	expected := strings.Split(raw, "\n")

	for i, v := range got {
		if i >= len(expected) {
			t.Errorf("ran out of expected tokens")
			return
		}

		if v != expected[i] {
			t.Errorf("expected %v but got %v", v, expected[i])
		}
	}
}

func TestHeadingExt(t *testing.T) {
	runTest(t, "EXT. FOO BAR", []string{
		`<slugline value="EXT. FOO BAR">`,
	})
}

func TestHeadingInt(t *testing.T) {
	runTest(t, "INT. FOO BAR", []string{
		`<slugline value="INT. FOO BAR">`,
	})
}

func TestCustomHeading(t *testing.T) {
	runTest(t, ".FOO BAR", []string{
		`<slugline value=".FOO BAR">`,
	})
}

func TestPagebreak(t *testing.T) {
	runTest(t, "some action text\n\n===\n\n", []string{
		`<text value="some action text">`,
		`<equals>`,
		`<equals>`,
		`<equals>`,
	})
}

func TestCenteredText(t *testing.T) {
	runTest(t, "> THE END <", []string{
		`<centered_text value="> THE END <">`,
	})
}

func TestParenthetical(t *testing.T) {
	runTest(t, "(starting the engine)", []string{
		`<paren_open>`,
		`<text value="starting the engine">`,
		`<paren_close>`,
	})
}

func TestActionWithWhitespace(t *testing.T) {
	runTest(t, "   some action   ", []string{
		`<text value="   some action   ">`,
	})
}

func TestActionWithSpecialCharacters(t *testing.T) {
	runTest(t, "some action/other stuff", []string{
		`<text value="some action/other stuff">`,
	})
}

func TestUnderline(t *testing.T) {
	runTest(t, "_some action_", []string{
		`<underscore>`,
		`<text value="some action">`,
		`<underscore>`,
	})
}

func TestItalics(t *testing.T) {
	runTest(t, "*some action*", []string{
		`<asterisk>`,
		`<text value="some action">`,
		`<asterisk>`,
	})
}

func TestBold(t *testing.T) {
	runTest(t, "**some action**", []string{
		`<asterisk>`,
		`<asterisk>`,
		`<text value="some action">`,
		`<asterisk>`,
		`<asterisk>`,
	})
}

func TestBoneyard(t *testing.T) {
	runTest(t, "/* some action */", []string{
		`<boneyard_open>`,
		`<text value=" some action ">`,
		`<boneyard_end>`,
	})
}

func TestBoneyardWithinText(t *testing.T) {
	runTest(t, "Some more action /* some action */ more stuff", []string{
		`<text value="Some more action ">`,
		`<boneyard_open>`,
		`<text value=" some action ">`,
		`<boneyard_end>`,
		`<text value=" more stuff">`,
	})
}

func TestNotes(t *testing.T) {
	runTest(t, "[[ some action ]]", []string{
		`<note_open>`,
		`<text value=" some action ">`,
		`<note_close>`,
	})
}

func TestNotesWithinText(t *testing.T) {
	runTest(t, "Some more action [[ some action ]] some other action", []string{
		`<text value="Some more action ">`,
		`<note_open>`,
		`<text value=" some action ">`,
		`<note_close>`,
		`<text value=" some other action">`,
	})
}

func TestDialogue(t *testing.T) {
	runTest(t, "FOO BAR\nI'm a cool dude.", []string{
		`<dialogue value="FOO BAR">`,
		`<text value="I'm a cool dude.">`,
	})
}

func TestTransition(t *testing.T) {
	runTest(t, "CUT TO:", []string{
		`<transition value="CUT TO:">`,
	})
}

func TestBeginsWithMatchTruthy(t *testing.T) {
	got := beginsWith("some string", "some")
	expected := true

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("expected %v but got %v", expected, got)
	}
}

func TestBeginsWithFalsy(t *testing.T) {
	got := beginsWith("some string", "foo")
	expected := false

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("expected %v but got %v", expected, got)
	}
}

func TestEndsWithMatchTruthy(t *testing.T) {
	got := endsWith("CUT TO:", "TO:")
	expected := true

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("expected %v but got %v", expected, got)
	}
}

func TestEndsWithFalsy(t *testing.T) {
	got := endsWith("CUT TO:", "foo")
	expected := false

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("expected %v but got %v", expected, got)
	}
}
