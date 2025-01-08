package chess

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
)

type validNotationTest struct {
	Pos1        *Position
	Pos2        *Position
	AlgText     string
	LongAlgText string
	UCIText     string
	Description string
}

func TestValidDecoding(t *testing.T) {
	f, err := os.Open("fixtures/valid_notation_tests.json")
	if err != nil {
		t.Fatal(err)
		return
	}

	validTests := []validNotationTest{}
	if err := json.NewDecoder(f).Decode(&validTests); err != nil {
		t.Fatal(err)
		return
	}

	for _, test := range validTests {
		for i, n := range []Notation{AlgebraicNotation{}, LongAlgebraicNotation{}, UCINotation{}} {
			var moveText string
			switch i {
			case 0:
				moveText = test.AlgText
			case 1:
				moveText = test.LongAlgText
			case 2:
				moveText = test.UCIText
			}
			m, err := n.Decode(test.Pos1, moveText)
			if err != nil {
				movesStrList := []string{}
				for _, m := range test.Pos1.ValidMoves() {
					s := n.Encode(test.Pos1, m)
					movesStrList = append(movesStrList, s)
				}
				t.Fatalf("starting from board \n%s\n expected move to be valid error - %s %s\n", test.Pos1.board.Draw(), err, strings.Join(movesStrList, ","))
			}
			postPos := test.Pos1.Update(m)
			if test.Pos2.String() != postPos.String() {
				t.Errorf("starting from board \n%s%s\n after move %s\n expected board to be %s\n%s\n but was %s\n%s\n",
					test.Pos1.String(),
					test.Pos1.board.Draw(), m.String(), test.Pos2.String(),
					test.Pos2.board.Draw(), postPos.String(), postPos.board.Draw())
			}
		}

	}
}

type notationDecodeTest struct {
	N       Notation
	Pos     *Position
	Text    string
	PostPos *Position
}

var (
	invalidDecodeTests = []notationDecodeTest{}
)

func TestInvalidDecoding(t *testing.T) {
	for _, test := range invalidDecodeTests {
		if _, err := test.N.Decode(test.Pos, test.Text); err == nil {
			t.Errorf("starting from board\n%s\n expected move notation %s to be invalid", test.Pos.board.Draw(), test.Text)
		}
	}
}
