package chess

import (
	"log"
	"testing"
)

type moveTest struct {
	pos     *Position
	m       *Move
	postPos *Position
}

var (
	validMoves = []moveTest{}

	invalidMoves = []moveTest{}

	positionUpdates = []moveTest{}
)

func unsafeFEN(s string) *Position {
	pos, err := decodeFEN(s, false)
	if err != nil {
		log.Printf("move_test.go: unsafeFEN(): error for fen %s\n", s)
		log.Fatal(err)
	}
	return pos
}

func TestValidMoves(t *testing.T) {
	for _, mt := range validMoves {
		if !moveIsValid(mt.pos, mt.m, false) {
			log.Println(mt.pos.String())
			log.Println(mt.pos.board.Draw())
			log.Println(mt.pos.ValidMoves())
			log.Println("In Check:", squaresAreAttacked(mt.pos, mt.pos.board.whiteKingSq))
			// log.Println("In Check:", mt.pos.inCheck())
			mt.pos.turn = mt.pos.turn.Other()
			t.Fatalf("expected move %s to be valid", mt.m)
		}
	}
}

func TestInvalidMoves(t *testing.T) {
	for _, mt := range invalidMoves {
		if moveIsValid(mt.pos, mt.m, false) {
			log.Println(mt.pos.String())
			log.Println(mt.pos.board.Draw())
			t.Fatalf("expected move %s to be invalid", mt.m)
		}
	}
}

func TestPositionUpdates(t *testing.T) {
	for _, mt := range positionUpdates {
		if !moveIsValid(mt.pos, mt.m, true) {
			log.Println(mt.pos.String())
			log.Println(mt.pos.board.Draw())
			log.Println(mt.pos.ValidMoves())
			t.Fatalf("expected move %s %v to be valid", mt.m, mt.m.tags)
		}

		postPos := mt.pos.Update(mt.m)
		if postPos.String() != mt.postPos.String() {
			t.Fatalf("starting from board \n%s%s\n after move %s\n expected board to be %s\n%s\n but was %s\n%s\n",
				mt.pos.String(),
				mt.pos.board.Draw(),
				mt.m.String(),
				mt.postPos.String(),
				mt.postPos.board.Draw(),
				postPos.String(),
				postPos.board.Draw(),
			)
		}
	}
}

type perfTest struct {
	pos           *Position
	nodesPerDepth []int
}

/* https://www.chessprogramming.org/Perft_Results */
var perfResults = []perfTest{}

func TestPerfResults(t *testing.T) {
	for _, perf := range perfResults {
		countMoves(t, perf.pos, []*Position{perf.pos}, perf.nodesPerDepth, len(perf.nodesPerDepth))
	}
}

func countMoves(t *testing.T, originalPosition *Position, positions []*Position, nodesPerDepth []int, maxDepth int) {
	if len(nodesPerDepth) == 0 {
		return
	}
	depth := maxDepth - len(nodesPerDepth) + 1
	expNodes := nodesPerDepth[0]
	newPositions := make([]*Position, 0)
	for _, pos := range positions {
		for _, move := range pos.ValidMoves() {
			newPos := pos.Update(move)
			newPositions = append(newPositions, newPos)
		}
	}
	gotNodes := len(newPositions)
	if expNodes != gotNodes {
		t.Errorf("Depth: %d Expected: %d Got: %d", depth, expNodes, gotNodes)
		t.Log("##############################")
		t.Log("# Original position info")
		t.Log("###")
		t.Log(originalPosition.String())
		t.Log(originalPosition.board.Draw())
		t.Log("##############################")
		t.Log("# Details in JSONL (http://jsonlines.org)")
		t.Log("###")
		for _, pos := range positions {
			t.Logf(`{"position": "%s", "moves": %d}`, pos.String(), len(pos.ValidMoves()))
		}
	}
	countMoves(t, originalPosition, newPositions, nodesPerDepth[1:], maxDepth)
}

func BenchmarkValidMoves(b *testing.B) {
	pos := unsafeFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		pos.ValidMoves()
		pos.validMoves = nil
	}
}

func moveIsValid(pos *Position, m *Move, useTags bool) bool {
	for _, move := range pos.ValidMoves() {
		if move.s1 == m.s1 && move.s2 == m.s2 && move.promo == m.promo {
			if useTags {
				if m.tags != move.tags {
					return false
				}
			}
			return true
		}
	}
	return false
}
