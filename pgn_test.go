package chess

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

type pgnTest struct {
	PostPos *Position
	PGN     string
}

var (
	validPGNs = []pgnTest{}
)

func TestValidPGNs(t *testing.T) {
	for _, test := range validPGNs {
		game, err := decodePGN(test.PGN)
		if err != nil {
			t.Fatalf("recieved unexpected pgn error %s", err.Error())
		}
		if game.Position().String() != test.PostPos.String() {
			t.Fatalf("expected board to be \n%s\nFEN:%s\n but got \n%s\n\nFEN:%s\n",
				test.PostPos.board.Draw(), test.PostPos.String(),
				game.Position().board.Draw(), game.Position().String())
		}
	}
}

type commentTest struct {
	PGN         string
	MoveNumber  int
	CommentText string
}

var (
	commentTests = []commentTest{}
)

func TestCommentsDetection(t *testing.T) {
	for _, test := range commentTests {
		game, err := decodePGN(test.PGN)
		if err != nil {
			t.Fatal(err)
		}
		comment := strings.Join(game.Comments()[test.MoveNumber], " ")
		if comment != test.CommentText {
			t.Fatalf("expected pgn comment to be %s but got %s", test.CommentText, comment)
		}
	}
}

func TestNewGameComments(t *testing.T) {
	for _, test := range commentTests {
		pgn, err := PGN(strings.NewReader(test.PGN))
		if err != nil {
			t.Fatal(err)
		}
		game := NewGame(pgn)
		comment := strings.Join(game.Comments()[test.MoveNumber], " ")
		if comment != test.CommentText {
			t.Fatalf("expected pgn comment to be %s but got %s", test.CommentText, comment)
		}
	}
}

func TestWriteComments(t *testing.T) {
	pgn := mustParsePGN("fixtures/pgns/0005.pgn")
	game, err := decodePGN(pgn)
	if err != nil {
		t.Fatal(err)
	}
	game, err = decodePGN(game.String())
	if err != nil {
		t.Fatal(err)
	}
	if len(game.Comments()[7]) != 2 {
		t.Fatalf("expected %d comments for move 7 but got %d", 2, len(game.Comments()[7]))
	}
}

func TestScanner(t *testing.T) {
	m := map[string]int{
		"fixtures/pgns/0006.pgn": 5,
		"fixtures/pgns/0007.pgn": 5,
		"fixtures/pgns/0013.pgn": 3,
	}
	for fname, count := range m {
		f, err := os.Open(fname)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		scanner := NewScanner(f)
		games := []*Game{}
		for scanner.Scan() {
			game := scanner.Next()
			if len(game.Moves()) == 0 {
				continue
			}
			games = append(games, game)
		}
		if len(games) != count {
			t.Fatalf(fname+" expected %d games but got %d", count, len(games))
		}
	}
}

func TestScannerWithFromPosFENs(t *testing.T) {
	finalPositions := []string{
		"rnbqkbnr/pp2pppp/2p5/3p4/3PP3/5P2/PPP3PP/RNBQKBNR b KQkq - 0 3",
		"r2qkb1r/pp1n1ppp/2p2n2/4p3/2BPP1b1/2P2N2/PP4PP/RNBQ1RK1 b kq - 0 8",
		"rnbqk2r/pp2nppp/2p1p3/3p4/1b1PP3/2NB1P2/PPPB2PP/R2QK1NR b KQkq - 5 6",
		"rnbqk1nr/pp2ppbp/2p3p1/3p4/3PP3/2N1BP2/PPP3PP/R2QKBNR b KQkq - 3 5",
		"rnb1kbnr/pp3ppp/1qp5/8/3NP3/2N5/PPP3PP/R1BQKB1R b KQkq - 0 7",
	}
	fname := "fixtures/pgns/0014.pgn"
	f, err := os.Open(fname)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	scanner := NewScanner(f)
	games := []*Game{}
	for idx := 0; scanner.Scan(); {
		game := scanner.Next()
		if len(game.moves) == 0 {
			continue
		}
		finalPos := game.Position().String()
		if finalPos != finalPositions[idx] {
			t.Fatalf(fname+" game %v expected final pos %v but got %v", idx,
				finalPositions[idx], finalPos)
		}
		games = append(games, game)
		idx++
	}
	if len(games) != len(finalPositions) {
		t.Fatalf(fname+" expected %v games but got %v", len(finalPositions),
			len(games))
	}
}

func BenchmarkPGN(b *testing.B) {
	pgn := mustParsePGN("fixtures/pgns/0001.pgn")
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		opt, _ := PGN(strings.NewReader(pgn))
		NewGame(opt)
	}
}

func mustParsePGN(fname string) string {
	f, err := os.Open(fname)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	return string(b)
}
func TestGamesFromPGN(t *testing.T) {
	for _, test := range validPGNs {
		reader := strings.NewReader(test.PGN)
		games, err := GamesFromPGN(reader)
		if err != nil {
			t.Fatalf("fail to read games from valid pgn: %s", err.Error())
		}
		if len(games) != 1 {
			t.Fatalf("expected to get 1 game from pgn, got %d", len(games))
		}
	}
}
