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

	raw := `<text value="Title:">
<underscore>
<asterisk>
<asterisk>
<text value="BRICK & STEEL">
<asterisk>
<asterisk>
<underscore>
<underscore>
<asterisk>
<asterisk>
<text value="FULL RETIRED">
<asterisk>
<asterisk>
<underscore>
<text value="Credit: Written by">
<text value="Author: Stu Maschwitz">
<text value="Source: Story by KTM">
<text value="Draft date: 1/27/2012">
<text value="Contact:">
<text value="Next Level Productions">
<text value="1588 Mission Dr.">
<text value="Solvang, CA 93463">
<text value="EXT. BRICK'S PATIO - DAY">
<text value="A gorgeous day.  The sun is shining.  But BRICK BRADDOCK, retired police detective, is sitting quietly, contemplating -- something.">
<text value="The SCREEN DOOR slides open and DICK STEEL, his former partner and fellow retiree, emerges with two cold beers.">
<text value="STEEL">
<text value="Beer's ready!">
<text value="BRICK">
<text value="Are they cold?">
<text value="STEEL">
<text value="Does a bear crap in the woods?">
<text value="Steel sits.  They laugh at the dumb joke.">
<text value="STEEL">
<paren_open>
<text value="beer raised">
<paren_close>
<text value="To retirement.">
<text value="BRICK">
<text value="To retirement.">
<text value="They drink long and well from the beers.">
<text value="And then there's a long beat.">
<text value="Longer than is funny.">
<text value="Long enough to be depressing.">
<text value="The men look at each other.">
<text value="STEEL">
<text value="Screw retirement.">
<text value="BRICK ^">
<text value="Screw retirement.">
<text value="SMASH CUT TO:">
<text value="INT. TRAILER HOME - DAY">
<text value="This is the home of THE BOY BAND, AKA DAN and JACK.  They too are drinking beer, and counting the take from their last smash-and-grab.  Money, drugs, and ridiculous props are strewn about the table.">
<text value="JACK">
<paren_open>
<text value="in Vietnamese, subtitled">
<paren_close>
<asterisk>
<text value="Did you know Brick and Steel are retired?">
<asterisk>
<text value="DAN">
<text value="Then let's retire them.">
<underscore>
<text value="Permanently">
<underscore>
<text value=".">
<text value="Jack begins to argue vociferously in Vietnamese ">
<paren_open>
<text value="?">
<paren_close>
<text value=", But mercifully we...">
<text value="CUT TO:">
<text value="EXT. BRICK'S POOL - DAY">
<text value="Steel, in the middle of a heated phone call:">
<text value="STEEL">
<text value="They're coming out of the woodwork!">
<paren_open>
<text value="pause">
<paren_close>
<text value="No, everybody we've put away!">
<paren_open>
<text value="pause">
<paren_close>
<text value="Point Blank Sniper?">
<text value=".SNIPER SCOPE POV">
<text value="From what seems like only INCHES AWAY.  ">
<underscore>
<text value="Steel's face FILLS the ">
<asterisk>
<text value="Leupold Mark 4">
<asterisk>
<text value=" scope">
<underscore>
<text value=".">
<text value="STEEL">
<text value="The man's a myth!">
<text value="Steel turns and looks straight into the cross-hairs.">
<text value="STEEL">
<paren_open>
<text value="oh crap">
<paren_close>
<text value="Hello...">
<text value="CUT TO:">
<text value=".OPENING TITLES">
<caret_close>
<text value=" BRICK BRADDOCK ">
<caret_open>
<caret_close>
<text value=" & DICK STEEL IN ">
<caret_open>
<caret_close>
<text value=" BRICK & STEEL ">
<caret_open>
<caret_close>
<text value=" FULL RETIRED ">
<caret_open>
<text value="SMASH CUT TO:">
<text value="EXT. WOODEN SHACK - DAY">
<text value="COGNITO, the criminal mastermind, is SLAMMED against the wall.">
<text value="COGNITO">
<text value="Woah woah woah, Brick and Steel!">
<text value="Sure enough, it's Brick and Steel, roughing up their favorite usual suspect.">
<text value="COGNITO">
<text value="What is it you want with me, DICK?">
<text value="Steel SMACKS him.">
<text value="STEEL">
<text value="Who's coming after us?">
<text value="COGNITO">
<text value="Everyone's coming after you mate!  Scorpio, The Boy Band, Sparrow, Point Blank Sniper...">
<text value="As he rattles off the long list, Brick and Steel share a look.  This is going to be BAD.">
<text value="CUT TO:">
<text value="INT. GARAGE - DAY">
<text value="BRICK and STEEL get into Mom's PORSCHE, Steel at the wheel.  They pause for a beat, the gravity of the situation catching up with them.">
<text value="BRICK">
<text value="This is everybody we've ever put away.">
<text value="STEEL">
<paren_open>
<text value="starting the engine">
<paren_close>
<text value="So much for retirement!">
<text value="They speed off.  To destiny!">
<text value="CUT TO:">
<text value="EXT. PALATIAL MANSION - DAY">
<text value="An EXTREMELY HANDSOME MAN drinks a beer.  Shirtless, unfortunately.">
<text value="His minion approaches offscreen:">
<text value="MINION">
<text value="We found Brick and Steel!">
<text value="HANDSOME MAN">
<text value="I want them dead.  DEAD!">
<text value="Beer flies.">
<caret_close>
<text value=" BURN TO PINK.">
<caret_close>
<text value=" THE END ">
<caret_open>`

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

// Headings
func TestHeadingExt(t *testing.T) {
	runTest(t, "EXT. FOO BAR\n\nsome action text", []string{
		`<text value="EXT. FOO BAR">`,
		`<text value="some action text">`,
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
		`<caret_close>`,
		`<text value=" THE END ">`,
		`<caret_open>`,
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
