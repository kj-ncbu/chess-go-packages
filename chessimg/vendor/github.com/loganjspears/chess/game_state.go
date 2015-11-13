package chess

import (
	"fmt"
	"strings"
)

// Side represents a side of the board.
type Side int

const (
	// KingSide is the right side of the board from white's perspective.
	KingSide Side = iota + 1
	// QueenSide is the left side of the board from white's perspective.
	QueenSide
)

// CastleRights holds the state of both sides castling abilities.
type CastleRights string

// CanCastle returns true if the given color and side combination
// can castle, otherwise returns false.
func (cr CastleRights) CanCastle(c Color, side Side) bool {
	char := "k"
	if side == QueenSide {
		char = "q"
	}
	if c == White {
		char = strings.ToUpper(char)
	}
	return strings.Contains(string(cr), char)
}

// String implements the fmt.Stringer interface and returns
// a FEN compatible string.  Ex. KQq
func (cr CastleRights) String() string {
	return string(cr)
}

// GameState represents the state of the game without regaurd
// to its outcome.  GameState is translatable to FEN notation.
type GameState struct {
	board           Board
	turn            Color
	castleRights    CastleRights
	enPassantSquare *Square
	halfMoveClock   int
	moveCount       int
}

// Board returns the gamestate's board.
func (gs *GameState) Board() Board {
	return gs.board.copy()
}

// Turn returns the color to move next.
func (gs *GameState) Turn() Color {
	return gs.turn
}

// CastleRights returns the castling rights of the state.
func (gs *GameState) CastleRights() CastleRights {
	return gs.castleRights
}

// String implements the fmt.Stringer interface and returns a
// string with the format: rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1
func (gs *GameState) String() string {
	b := gs.board.String()
	t := gs.turn.String()
	c := gs.castleRights.String()
	sq := "-"
	if gs.enPassantSquare != nil {
		sq = gs.enPassantSquare.String()
	}
	return fmt.Sprintf("%s %s %s %s %d %d", b, t, c, sq, gs.halfMoveClock, gs.moveCount)
}

// MarshalText implements the encoding.TextMarshaler interface and
// encodes the gamestate's FEN.
func (gs *GameState) MarshalText() (text []byte, err error) {
	return []byte(gs.String()), nil
}

// UnmarshalText implements the encoding.TextUnarshaler interface and
// assumes the data is in the FEN format.
func (gs *GameState) UnmarshalText(text []byte) error {
	state, err := decodeFEN(string(text))
	if err != nil {
		return err
	}
	gs.board = state.board
	gs.turn = state.turn
	gs.castleRights = state.castleRights
	gs.enPassantSquare = state.enPassantSquare
	gs.halfMoveClock = state.halfMoveClock
	gs.moveCount = state.moveCount
	return nil
}

func (gs *GameState) getOutcome() (Outcome, Method) {
	// should only happen in unit tests
	if gs.board.kingSquare(gs.turn) == nil {
		return NoOutcome, NoMethod
	}
	inCheck := gs.board.inCheck(gs.turn)
	hasMove := len(gs.validMoves()) > 0
	if !inCheck && !hasMove {
		return Draw, Stalemate
	} else if inCheck && !hasMove {
		switch gs.turn {
		case White:
			return BlackWon, Checkmate
		case Black:
			return WhiteWon, Checkmate
		}
	}
	return NoOutcome, NoMethod
}

func (gs *GameState) validMoves() []*Move {
	moves := []*Move{}
	for _, s1 := range gs.board.squaresForColor(gs.turn) {
		p := gs.board.piece(s1)
		// TODO s2 can only be an empty or enemy square
		for _, s2 := range allSquares {
			couldPromo := p.Type() == Pawn && (s2.rank == R1 || s2.rank == R8)
			if couldPromo {
				for _, pt := range PieceTypes() {
					if pt.isPromotable() {
						m := &Move{s1: s1, s2: s2, state: gs, promo: pt}
						if m.isValid() {
							moves = append(moves, m)
						}
					}
				}
			} else {
				m := &Move{s1: s1, s2: s2, state: &GameState{board: gs.board, turn: gs.turn}}
				if m.isValid() {
					moves = append(moves, m)
				}
			}

		}
	}
	return moves
}

func (gs *GameState) samePosition(g *GameState) bool {
	return gs.board.String() == g.board.String() &&
		gs.turn == g.turn &&
		gs.castleRights.String() == g.castleRights.String() &&
		gs.enPassantSquare == g.enPassantSquare
}