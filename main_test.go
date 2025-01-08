package chess

import (
	"fmt"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	fmt.Println("Setting up tests...")

	// Setup: Initialize resources
	validMoves = []moveTest{
		// pawn moves
		{m: &Move{s1: E2, s2: E4}, pos: unsafeFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")},
		{m: &Move{s1: A2, s2: A3}, pos: unsafeFEN("8/8/8/8/8/8/P7/5K1k w - - 0 1")},
		{m: &Move{s1: A7, s2: A6}, pos: unsafeFEN("k7/p7/8/8/8/8/8/7K b - - 0 1")},
		{m: &Move{s1: A7, s2: A5}, pos: unsafeFEN("k7/p7/8/8/8/8/8/7K b - - 0 1")},
		{m: &Move{s1: C4, s2: B5}, pos: unsafeFEN("k7/8/8/1p1p4/2P5/8/8/7K w - - 0 1")},
		{m: &Move{s1: C4, s2: D5}, pos: unsafeFEN("k7/8/8/1p1p4/2P5/8/8/7K w - - 0 1")},
		{m: &Move{s1: C4, s2: C5}, pos: unsafeFEN("k7/8/8/1p1p4/2P5/8/8/7K w - - 0 1")},
		{m: &Move{s1: C5, s2: B4}, pos: unsafeFEN("k7/8/8/2p5/1P1P4/8/8/7K b - - 0 1")},
		{m: &Move{s1: C5, s2: D4}, pos: unsafeFEN("k7/8/8/2p5/1P1P4/8/8/7K b - - 0 1")},
		{m: &Move{s1: C5, s2: C4}, pos: unsafeFEN("k7/8/8/2p5/1P1P4/8/8/7K b - - 0 1")},
		{m: &Move{s1: A4, s2: B3}, pos: unsafeFEN("2r3k1/1q1nbppp/r3p3/3pP3/pPpP4/P1Q2N2/2RN1PPP/2R4K b - b3 0 23")},
		{m: &Move{s1: A2, s2: A1, promo: Queen}, pos: unsafeFEN("K6k/8/8/8/8/8/p7/8 b - - 0 1")},
		{m: &Move{s1: E7, s2: E6}, pos: unsafeFEN("r2qkbnr/pppnpppp/8/3p4/6b1/1P3NP1/PBPPPP1P/RN1QKB1R b KQkq - 2 4")},
		// knight moves
		{m: &Move{s1: E4, s2: F6}, pos: unsafeFEN("k7/8/8/3pp3/4N3/8/5B2/7K w - - 0 1")},
		{m: &Move{s1: E4, s2: D6}, pos: unsafeFEN("k7/8/8/3pp3/4N3/8/5B2/7K w - - 0 1")},
		{m: &Move{s1: E4, s2: G5}, pos: unsafeFEN("k7/8/8/3pp3/4N3/8/5B2/7K w - - 0 1")},
		{m: &Move{s1: E4, s2: G3}, pos: unsafeFEN("k7/8/8/3pp3/4N3/8/5B2/7K w - - 0 1")},
		{m: &Move{s1: E4, s2: D2}, pos: unsafeFEN("k7/8/8/3pp3/4N3/8/5B2/7K w - - 0 1")},
		{m: &Move{s1: E4, s2: C3}, pos: unsafeFEN("k7/8/8/3pp3/4N3/8/5B2/7K w - - 0 1")},
		{m: &Move{s1: E4, s2: C5}, pos: unsafeFEN("k7/8/8/3pp3/4N3/8/5B2/7K w - - 0 1")},
		{m: &Move{s1: B8, s2: D7}, pos: unsafeFEN("rn1qkb1r/pp3ppp/2p1pn2/3p4/2PP4/2NQPN2/PP3PPP/R1B1K2R b KQkq - 0 7")},
		{m: &Move{s1: F6, s2: E4}, pos: unsafeFEN("r1b1k2r/ppp2ppp/2p2n2/4N3/4P3/2P5/PPP2PPP/R1BK3R b kq - 0 8")},
		// bishop moves
		{m: &Move{s1: E4, s2: H7}, pos: unsafeFEN("k7/8/8/3pp3/4B3/5N2/8/7K w - - 0 1")},
		{m: &Move{s1: E4, s2: D5}, pos: unsafeFEN("k7/8/8/3pp3/4B3/5N2/8/7K w - - 0 1")},
		{m: &Move{s1: E4, s2: B1}, pos: unsafeFEN("k7/8/8/3pp3/4B3/5N2/8/7K w - - 0 1")},
		// rook moves
		{m: &Move{s1: B2, s2: B4}, pos: unsafeFEN("k7/1p5b/4N3/4p3/8/8/1R6/1B5K w - - 0 1")},
		{m: &Move{s1: B2, s2: B7}, pos: unsafeFEN("k7/1p5b/4N3/4p3/8/8/1R6/1B5K w - - 0 1")},
		{m: &Move{s1: B2, s2: A2}, pos: unsafeFEN("k7/1p5b/4N3/4p3/8/8/1R6/1B5K w - - 0 1")},
		{m: &Move{s1: B2, s2: H2}, pos: unsafeFEN("k7/1p5b/4N3/4p3/8/8/1R6/1B5K w - - 0 1")},
		{m: &Move{s1: E1, s2: E8}, pos: unsafeFEN("r3r1k1/p4p1p/3p4/1p4p1/2pP4/2P2P2/PP3P1P/R3RK2 w - g6 0 22")},
		// queen moves
		{m: &Move{s1: B2, s2: E5}, pos: unsafeFEN("k7/1p5b/4N3/4p3/8/8/1Q6/1B5K w - - 0 1")},
		{m: &Move{s1: B2, s2: A1}, pos: unsafeFEN("k7/1p5b/4N3/4p3/8/8/1Q6/1B5K w - - 0 1")},
		{m: &Move{s1: B2, s2: A2}, pos: unsafeFEN("k7/1p5b/4N3/4p3/8/8/1Q6/1B5K w - - 0 1")},
		{m: &Move{s1: B2, s2: H2}, pos: unsafeFEN("k7/1p5b/4N3/4p3/8/8/1Q6/1B5K w - - 0 1")},
		{m: &Move{s1: D8, s2: D1}, pos: unsafeFEN("r1bqk2r/ppp2ppp/2p2n2/4N3/4P3/2P5/PPP2PPP/R1BQK2R b KQkq - 0 7")},
		// king moves
		{m: &Move{s1: E4, s2: E5}, pos: unsafeFEN("k4r2/8/8/8/4K3/8/8/8 w - - 0 1")},
		{m: &Move{s1: E4, s2: E3}, pos: unsafeFEN("k4r2/8/8/8/4K3/8/8/8 w - - 0 1")},
		{m: &Move{s1: E4, s2: D3}, pos: unsafeFEN("k4r2/8/8/8/4K3/8/8/8 w - - 0 1")},
		{m: &Move{s1: E4, s2: D4}, pos: unsafeFEN("k4r2/8/8/8/4K3/8/8/8 w - - 0 1")},
		{m: &Move{s1: E4, s2: D5}, pos: unsafeFEN("k4r2/8/8/8/4K3/8/8/8 w - - 0 1")},
		{m: &Move{s1: E4, s2: E5}, pos: unsafeFEN("k4r2/8/8/8/4K3/8/8/8 w - - 0 1")},
		// castling
		{m: &Move{s1: E1, s2: G1}, pos: unsafeFEN("r3k2r/8/8/8/8/8/8/R3K2R w KQkq - 0 1")},
		{m: &Move{s1: E1, s2: C1}, pos: unsafeFEN("r3k2r/8/8/8/8/8/8/R3K2R w KQkq - 0 1")},
		{m: &Move{s1: E8, s2: G8}, pos: unsafeFEN("r3k2r/8/8/8/8/8/8/R3K2R b KQkq - 0 1")},
		{m: &Move{s1: E8, s2: C8}, pos: unsafeFEN("r3k2r/8/8/8/8/8/8/R3K2R b KQkq - 0 1")},
		// king moving in front of enemy pawn http://en.lichess.org/4HXJOtpN#75
		{m: &Move{s1: F8, s2: G7}, pos: unsafeFEN("3rrk2/8/2p3P1/1p2nP1p/pP2p3/P1B1NbPB/2P2K2/5R2 b - - 1 38")},
	}

	invalidMoves = []moveTest{
		// out of turn moves
		{m: &Move{s1: E7, s2: E5}, pos: unsafeFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")},
		{m: &Move{s1: E2, s2: E4}, pos: unsafeFEN("rnbqkbnr/1ppppppp/p7/8/8/8/PPPPPPPP/RNBQKBNR b KQkq - 0 1")},
		// pawn moves
		{m: &Move{s1: E2, s2: D3}, pos: unsafeFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")},
		{m: &Move{s1: E2, s2: F3}, pos: unsafeFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")},
		{m: &Move{s1: E2, s2: E5}, pos: unsafeFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")},
		{m: &Move{s1: A2, s2: A1}, pos: unsafeFEN("k7/8/8/8/8/8/p7/7K b - - 0 1")},
		{m: &Move{s1: E6, s2: E5}, pos: unsafeFEN(`2b1r3/2k2p1B/p2np3/4B3/8/5N2/PP1K1PPP/3R4 b - - 2 1`)},
		{m: &Move{s1: H7, s2: H5}, pos: unsafeFEN(`2bqkbnr/rpppp2p/2n2p2/p5pB/5P2/4P3/PPPP2PP/RNBQK1NR b KQk - 4 6`)},
		// knight moves
		{m: &Move{s1: E4, s2: F2}, pos: unsafeFEN("k7/8/8/3pp3/4N3/8/5B2/7K w - - 0 1")},
		{m: &Move{s1: E4, s2: F3}, pos: unsafeFEN("k7/8/8/3pp3/4N3/8/5B2/7K w - - 0 1")},
		// bishop moves
		{m: &Move{s1: E4, s2: C6}, pos: unsafeFEN("k7/8/8/3pp3/4B3/5N2/8/7K w - - 0 1")},
		{m: &Move{s1: E4, s2: E5}, pos: unsafeFEN("k7/8/8/3pp3/4B3/5N2/8/7K w - - 0 1")},
		{m: &Move{s1: E4, s2: E4}, pos: unsafeFEN("k7/8/8/3pp3/4B3/5N2/8/7K w - - 0 1")},
		{m: &Move{s1: E4, s2: F3}, pos: unsafeFEN("k7/8/8/3pp3/4B3/5N2/8/7K w - - 0 1")},
		// rook moves
		{m: &Move{s1: B2, s2: B1}, pos: unsafeFEN("k7/1p5b/4N3/4p3/8/8/1R6/1B5K w - - 0 1")},
		{m: &Move{s1: B2, s2: C3}, pos: unsafeFEN("k7/1p5b/4N3/4p3/8/8/1R6/1B5K w - - 0 1")},
		{m: &Move{s1: B2, s2: B8}, pos: unsafeFEN("k7/1p5b/4N3/4p3/8/8/1R6/1B5K w - - 0 1")},
		{m: &Move{s1: B2, s2: G7}, pos: unsafeFEN("k7/1p5b/4N3/4p3/8/8/1R6/1B5K w - - 0 1")},
		// queen moves
		{m: &Move{s1: B2, s2: B1}, pos: unsafeFEN("k7/1p5b/4N3/4p3/8/8/1Q6/1B5K w - - 0 1")},
		{m: &Move{s1: B2, s2: C4}, pos: unsafeFEN("k7/1p5b/4N3/4p3/8/8/1Q6/1B5K w - - 0 1")},
		{m: &Move{s1: B2, s2: B8}, pos: unsafeFEN("k7/1p5b/4N3/4p3/8/8/1Q6/1B5K w - - 0 1")},
		{m: &Move{s1: B2, s2: G7}, pos: unsafeFEN("k7/1p5b/4N3/4p3/8/8/1Q6/1B5K w - - 0 1")},
		// king moves
		{m: &Move{s1: E4, s2: F3}, pos: unsafeFEN("k4r2/8/8/8/4K3/8/8/8 w - - 0 1")},
		{m: &Move{s1: E4, s2: F4}, pos: unsafeFEN("k4r2/8/8/8/4K3/8/8/8 w - - 0 1")},
		{m: &Move{s1: E4, s2: F5}, pos: unsafeFEN("k4r2/8/8/8/4K3/8/8/8 w - - 0 1")},
		// castleing
		{m: &Move{s1: E1, s2: B1}, pos: unsafeFEN("r3k2r/8/8/8/8/8/8/R3K2R w KQkq - 0 1")},
		{m: &Move{s1: E8, s2: B8}, pos: unsafeFEN("r3k2r/8/8/8/8/8/8/R3K2R b KQkq - 0 1")},
		{m: &Move{s1: E1, s2: C1}, pos: unsafeFEN("r3k2r/8/8/8/8/8/8/R2QK2R w KQkq - 0 1")},
		{m: &Move{s1: E1, s2: C1}, pos: unsafeFEN("2r1k2r/8/8/8/8/8/8/R3K2R w KQk - 0 1")},
		{m: &Move{s1: E1, s2: C1}, pos: unsafeFEN("3rk2r/8/8/8/8/8/8/R3K2R w KQk - 0 1")},
		{m: &Move{s1: E1, s2: G1}, pos: unsafeFEN("r3k2r/8/8/8/8/8/8/R3K2R w Qkq - 0 1")},
		{m: &Move{s1: E1, s2: C1}, pos: unsafeFEN("r3k2r/8/8/8/8/8/8/R3K2R w Kkq - 0 1")},
		// invalid promotion for non-pawn move
		{m: &Move{s1: B8, s2: D7, promo: Pawn}, pos: unsafeFEN("rn1qkb1r/pp3ppp/2p1pn2/3p4/2PP4/2NQPN2/PP3PPP/R1B1K2R b KQkq - 0 7")},
		// en passant on doubled pawn file http://en.lichess.org/TnRtrHxf#24
		{m: &Move{s1: E3, s2: F6}, pos: unsafeFEN("r1b2rk1/pp2b1pp/1qn1p3/3pPp2/1P1P4/P2BPN2/6PP/RN1Q1RK1 w - f6 0 13")},
		// can't move piece out of pin (even if checking enemy king) http://en.lichess.org/JCRBhXH7#62
		{m: &Move{s1: E1, s2: E7}, pos: unsafeFEN("4R3/1r1k2pp/p1p5/1pP5/8/8/1PP3PP/2K1Rr2 w - - 5 32")},
		// invalid one up pawn capture
		{m: &Move{s1: E6, s2: E5}, pos: unsafeFEN(`2b1r3/2k2p1B/p2np3/4B3/8/5N2/PP1K1PPP/3R4 b - - 2 1`)},
		// invalid two up pawn capture
		{m: &Move{s1: H7, s2: H5}, pos: unsafeFEN(`2bqkbnr/rpppp2p/2n2p2/p5pB/5P2/4P3/PPPP2PP/RNBQK1NR b KQk - 4 6`)},
		// invalid pawn move d5e4
		{m: &Move{s1: D5, s2: E4}, pos: unsafeFEN(`rnbqkbnr/pp2pppp/8/2pp4/3P4/4PN2/PPP2PPP/RNBQKB1R b KQkq - 0 3`)},
	}

	positionUpdates = []moveTest{
		{
			m:       &Move{s1: E2, s2: E4},
			pos:     unsafeFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"),
			postPos: unsafeFEN("rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1"),
		},
		{
			m:       &Move{s1: E1, s2: G1, tags: KingSideCastle},
			pos:     unsafeFEN("r3k2r/8/8/8/8/8/8/R3K2R w KQkq - 0 1"),
			postPos: unsafeFEN("r3k2r/8/8/8/8/8/8/R4RK1 b kq - 1 1"),
		},
		{
			m:       &Move{s1: A4, s2: B3, tags: EnPassant},
			pos:     unsafeFEN("2r3k1/1q1nbppp/r3p3/3pP3/pPpP4/P1Q2N2/2RN1PPP/2R4K b - b3 0 23"),
			postPos: unsafeFEN("2r3k1/1q1nbppp/r3p3/3pP3/11pP4/PpQ2N2/2RN1PPP/2R4K w - - 0 24"),
		},
		{
			m:       &Move{s1: E1, s2: G1, tags: KingSideCastle},
			pos:     unsafeFEN("r2qk2r/pp1n1ppp/2pbpn2/3p4/2PP4/1PNQPN2/P4PPP/R1B1K2R w KQkq - 1 9"),
			postPos: unsafeFEN("r2qk2r/pp1n1ppp/2pbpn2/3p4/2PP4/1PNQPN2/P4PPP/R1B2RK1 b kq - 2 9"),
		},
		// half move clock - knight move to f3 from starting position
		{
			m:       &Move{s1: G1, s2: F3},
			pos:     unsafeFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"),
			postPos: unsafeFEN("rnbqkbnr/pppppppp/8/8/8/5N2/PPPPPPPP/RNBQKB1R b KQkq - 1 1"),
		},
		// half move clock - king side castle
		{
			m:       &Move{s1: E1, s2: G1, tags: KingSideCastle},
			pos:     unsafeFEN("r3k2r/8/8/8/8/8/8/R3K2R w KQkq - 0 1"),
			postPos: unsafeFEN("r3k2r/8/8/8/8/8/8/R4RK1 b kq - 1 1"),
		},
		// half move clock - queen side castle
		{
			m:       &Move{s1: E1, s2: C1, tags: QueenSideCastle},
			pos:     unsafeFEN("r3k2r/ppqn1ppp/2pbpn2/3p4/2PP4/1PNQPN2/P2B1PPP/R3K2R w KQkq - 3 10"),
			postPos: unsafeFEN("r3k2r/ppqn1ppp/2pbpn2/3p4/2PP4/1PNQPN2/P2B1PPP/2KR3R b kq - 4 10"),
		},
		// half move clock - pawn push
		{
			m:       &Move{s1: E2, s2: E4},
			pos:     unsafeFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"),
			postPos: unsafeFEN("rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1"),
		},
		// half move clock - pawn capture
		{
			m:       &Move{s1: E4, s2: D5, tags: Capture},
			pos:     unsafeFEN("r1bqkbnr/ppp1pppp/2n5/3p4/4P3/5N2/PPPP1PPP/RNBQKB1R w KQkq - 2 3"),
			postPos: unsafeFEN("r1bqkbnr/ppp1pppp/2n5/3P4/8/5N2/PPPP1PPP/RNBQKB1R b KQkq - 0 3"),
		},
		// half move clock - en passant
		{
			m:       &Move{s1: E5, s2: F6, tags: EnPassant},
			pos:     unsafeFEN("r1bqkbnr/ppp1p1pp/2n5/3pPp2/8/5N2/PPPP1PPP/RNBQKB1R w KQkq f6 0 4"),
			postPos: unsafeFEN("r1bqkbnr/ppp1p1pp/2n2P2/3p4/8/5N2/PPPP1PPP/RNBQKB1R b KQkq - 0 4"),
		},
		// half move clock - piece captured by knight
		{
			m:       &Move{s1: C6, s2: D4, tags: Capture},
			pos:     unsafeFEN("r1bqkbnr/ppp1p1pp/2n5/3pPp2/3N4/8/PPPP1PPP/RNBQKB1R b KQkq - 1 4"),
			postPos: unsafeFEN("r1bqkbnr/ppp1p1pp/8/3pPp2/3n4/8/PPPP1PPP/RNBQKB1R w KQkq - 0 5"),
		},
	}

	perfResults = []perfTest{
		{pos: unsafeFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"), nodesPerDepth: []int{
			20, 400, 8902, 197281,
			// 4865609, 119060324, 3195901860, 84998978956, 2439530234167, 69352859712417
		}},
		{pos: unsafeFEN("r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1"), nodesPerDepth: []int{
			48, 2039, 97862,
			// 4085603, 193690690
		}},
		{pos: unsafeFEN("8/2p5/3p4/KP5r/1R3p1k/8/4P1P1/8 w - - 0 1"), nodesPerDepth: []int{
			14, 191, 2812, 43238, 674624,
			// 11030083, 178633661
		}},
		{pos: unsafeFEN("r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1"), nodesPerDepth: []int{
			6, 264, 9467, 422333,
			// 15833292, 706045033
		}},
		{pos: unsafeFEN("r2q1rk1/pP1p2pp/Q4n2/bbp1p3/Np6/1B3NBn/pPPP1PPP/R3K2R b KQ - 0 1"), nodesPerDepth: []int{
			6, 264, 9467, 422333,
			// 15833292, 706045033
		}},
		{pos: unsafeFEN("rnbq1k1r/pp1Pbppp/2p5/8/2B5/8/PPP1NnPP/RNBQK2R w KQ - 1 8"), nodesPerDepth: []int{
			44, 1486, 62379,
			// 2103487, 89941194
		}},
		{pos: unsafeFEN("r4rk1/1pp1qppp/p1np1n2/2b1p1B1/2B1P1b1/P1NP1N2/1PP1QPPP/R4RK1 w - - 0 10"), nodesPerDepth: []int{
			46, 2079, 89890,
			// 3894594, 164075551, 6923051137, 287188994746, 11923589843526, 490154852788714
		}},
	}

	invalidDecodeTests = []notationDecodeTest{
		{
			// opening for white
			N:    AlgebraicNotation{},
			Pos:  unsafeFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"),
			Text: "e5",
		},
		{
			// http://en.lichess.org/W91M4jms#14
			N:    AlgebraicNotation{},
			Pos:  unsafeFEN("rn1qkb1r/pp3ppp/2p1pn2/3p4/2PP4/2NQPN2/PP3PPP/R1B1K2R b KQkq - 0 7"),
			Text: "Nd7",
		},
		{
			// http://en.lichess.org/W91M4jms#17
			N:       AlgebraicNotation{},
			Pos:     unsafeFEN("r2qk2r/pp1n1ppp/2pbpn2/3p4/2PP4/1PNQPN2/P4PPP/R1B1K2R w KQkq - 1 9"),
			Text:    "O-O-O-O",
			PostPos: unsafeFEN("r2qk2r/pp1n1ppp/2pbpn2/3p4/2PP4/1PNQPN2/P4PPP/R1B2RK1 b kq - 2 9"),
		},
		{
			// http://en.lichess.org/W91M4jms#23
			N:    AlgebraicNotation{},
			Pos:  unsafeFEN("3r1rk1/pp1nqppp/2pbpn2/3p4/2PP4/1PNQPN2/PB3PPP/3RR1K1 b - - 5 12"),
			Text: "dx4",
		},
		{
			// should not assume pawn for unknown piece type "n"
			N:    AlgebraicNotation{},
			Pos:  unsafeFEN("rnbqkbnr/pppp1ppp/8/4p3/4P3/8/PPPP1PPP/RNBQKBNR w KQkq e6 0 2"),
			Text: "nf3",
		},
		{
			// disambiguation should not allow for this since it is not a capture
			N:    AlgebraicNotation{},
			Pos:  unsafeFEN("rnbqkbnr/ppp1pppp/8/3p4/3P4/8/PPP1PPPP/RNBQKBNR w KQkq - 0 2"),
			Text: "bf4",
		},
		{
			// invalid notation
			N:    AlgebraicNotation{},
			Pos:  unsafeFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"),
			Text: "g1f3",
		},
	}

	validPGNs = []pgnTest{
		{
			PostPos: unsafeFEN("4r3/6P1/2p2P1k/1p6/pP2p1R1/P1B5/2P2K2/3r4 b - - 0 45"),
			PGN:     mustParsePGN("fixtures/pgns/0001.pgn"),
		},
		{
			PostPos: unsafeFEN("4r3/6P1/2p2P1k/1p6/pP2p1R1/P1B5/2P2K2/3r4 b - - 0 45"),
			PGN:     mustParsePGN("fixtures/pgns/0002.pgn"),
		},
		{
			PostPos: unsafeFEN("2r2rk1/pp1bBpp1/2np4/2pp2p1/1bP5/1P4P1/P1QPPPBP/3R1RK1 b - - 0 3"),
			PGN:     mustParsePGN("fixtures/pgns/0003.pgn"),
		},
		{
			PostPos: unsafeFEN("r3kb1r/2qp1pp1/b1n1p2p/pp2P3/5n1B/1PPQ1N2/P1BN1PPP/R3K2R w KQkq - 1 14"),
			PGN:     mustParsePGN("fixtures/pgns/0004.pgn"),
		},
		{
			PostPos: unsafeFEN("rnbqkbnr/ppp2ppp/4p3/3p4/3PP3/8/PPP2PPP/RNBQKBNR w KQkq d6 0 3"),
			PGN:     mustParsePGN("fixtures/pgns/0008.pgn"),
		},
		{
			PostPos: unsafeFEN("r1bqkbnr/1ppp1ppp/p1n5/1B2p3/4P3/5N2/PPPP1PPP/RNBQK2R w KQkq - 0 4"),
			PGN:     mustParsePGN("fixtures/pgns/0009.pgn"),
		},
		{
			PostPos: unsafeFEN("r1bqkbnr/1ppp1ppp/p1n5/1B2p3/4P3/5N2/PPPP1PPP/RNBQK2R w KQkq - 0 4"),
			PGN:     mustParsePGN("fixtures/pgns/0010.pgn"),
		},
		{
			PostPos: unsafeFEN("8/8/6p1/4R3/6kQ/r2P1pP1/5P2/6K1 b - - 3 42"),
			PGN:     mustParsePGN("fixtures/pgns/0011.pgn"),
		},
		{
			PostPos: StartingPosition(),
			PGN:     mustParsePGN("fixtures/pgns/0012.pgn"),
		},
		{
			PostPos: unsafeFEN("rnbqkbnr/pppp1ppp/8/4p3/8/5N2/PPPPPPPP/RNBQKB1R w KQkq e6 0 2"),
			PGN:     mustParsePGN("fixtures/pgns/0015.pgn"),
		},
	}

	commentTests = []commentTest{
		{
			PGN:         mustParsePGN("fixtures/pgns/0005.pgn"),
			MoveNumber:  7,
			CommentText: `(-0.25 → 0.39) Inaccuracy. cxd4 was best. [%eval 0.39] [%clk 0:05:05]`,
		},
		{
			PGN:         mustParsePGN("fixtures/pgns/0009.pgn"),
			MoveNumber:  5,
			CommentText: `This opening is called the Ruy Lopez.`,
		},
		{
			PGN:         mustParsePGN("fixtures/pgns/0010.pgn"),
			MoveNumber:  5,
			CommentText: `This opening is called the Ruy Lopez.`,
		},
	}

	// Run the tests
	code := m.Run()

	// Teardown: Clean up resources (always executed, even if tests panic)
	fmt.Println("Tearing down tests...")

	// Exit with the appropriate code (important!)
	os.Exit(code)
}
