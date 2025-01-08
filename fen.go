package chess

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// Decodes FEN notation into a GameState.  An error is returned
// if there is a parsing error or if the FEN is illegal. FEN
// notation format: rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1
// for 960 mode it can decode both shredder-fen and x-fen wuthout knowing
// which one it is
func decodeFEN(fen string, isNineSixty bool) (*Position, error) {
	fen = strings.TrimSpace(fen)
	parts := strings.Split(fen, " ")
	if len(parts) != 6 {
		return nil, fmt.Errorf("chess: fen invalid notation %s must have 6 sections", fen)
	}
	b, err := fenBoard(parts[0])
	if err != nil {
		return nil, err
	}
	turn, ok := fenTurnMap[parts[1]]
	if !ok {
		return nil, fmt.Errorf("chess: fen invalid turn %s", parts[1])
	}
	rights, err := formCastleRights(parts[2], isNineSixty, b)
	if err != nil {
		return nil, err
	}
	sq, err := formEnPassant(parts[3], b, turn)
	if err != nil {
		return nil, err
	}
	halfMoveClock, err := strconv.Atoi(parts[4])
	if err != nil || halfMoveClock < 0 {
		return nil, fmt.Errorf("chess: fen invalid half move clock %s", parts[4])
	}
	moveCount, err := strconv.Atoi(parts[5])
	if err != nil || moveCount < 1 {
		return nil, fmt.Errorf("chess: fen invalid move count %s", parts[5])
	}

	pos := &Position{
		board:           b,
		turn:            turn,
		castleRights:    rights,
		enPassantSquare: sq,
		halfMoveClock:   halfMoveClock,
		moveCount:       moveCount,
	}

	// Make sure the player in next turn cannot capture opponent's king.
	cp := pos.copy()
	cp.turn = cp.turn.Other()
	if isInCheck(cp) {
		return nil, fmt.Errorf("chess: fen illegal %s , king can be captured in next move", fen)
	}

	return pos, nil
}

// generates board from fen format: rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR
func fenBoard(boardStr string) (*Board, error) {
	// make sure boardStr contains only one K and one k.
	if strings.Count(boardStr, "K") != 1 || strings.Count(boardStr, "k") != 1 {
		return nil, fmt.Errorf("chess: fen illegal board %s , both black and white should have one king each", boardStr)
	}
	rankStrs := strings.Split(boardStr, "/")
	if len(rankStrs) != 8 {
		return nil, fmt.Errorf("chess: fen invalid board %s", boardStr)
	}
	m := map[Square]Piece{}
	for i, rankStr := range rankStrs {
		rank := Rank(7 - i)
		if (rank == Rank1 || rank == Rank8) && (strings.ContainsRune(rankStr, 'p') || strings.ContainsRune(rankStr, 'P')) {
			return nil, fmt.Errorf("chess: fen illegal board %s , pawns cannot be on first or last rank", boardStr)
		}
		fileMap, err := fenFormRank(rankStr)
		if err != nil {
			return nil, err
		}
		for file, piece := range fileMap {
			m[NewSquare(file, rank)] = piece
		}
	}
	return NewBoard(m), nil
}

func fenFormRank(rankStr string) (map[File]Piece, error) {
	count := 0
	m := map[File]Piece{}
	err := fmt.Errorf("chess: fen invalid rank %s", rankStr)
	for _, r := range rankStr {
		c := fmt.Sprintf("%c", r)
		piece := fenPieceMap[c]
		if piece == NoPiece {
			skip, err := strconv.Atoi(c)
			if err != nil {
				return nil, err
			}
			count += skip
			continue
		}
		m[File(count)] = piece
		count++
	}
	if count != 8 {
		return nil, err
	}
	return m, nil
}

func formCastleRights(castleStr string, isNineSixty bool, board *Board) (*CastleRights, error) {
	cr := &CastleRights{
		nineSixtyMode:         isNineSixty,
		aSideRookStartingFile: "",
		hSideRookStartingFile: "",
		whiteKingSideCastle:   false,
		whiteQueenSideCastle:  false,
		blackKingSideCastle:   false,
		blackQueenSideCastle:  false,
	}
	if castleStr == "-" {
		return cr, nil
	}
	if len(castleStr) < 1 || len(castleStr) > 4 {
		return cr, fmt.Errorf("chess: fen invalid castle rights %s , length should be 1 to 4", castleStr)
	}
	if hasDuplicateCharacters(castleStr) {
		return cr, fmt.Errorf("chess: fen invalid castle rights %s, has duplicate characters", castleStr)
	}
	if !isNineSixty { // normal mode
		if !hasOnlyKQkq(castleStr) {
			return cr, fmt.Errorf("chess: fen invalid castle rights %s , normal mode should only have KQkq", castleStr)
		}
		if strings.ContainsRune(castleStr, 'K') {
			if board.Piece(E1) != WhiteKing || board.Piece(H1) != WhiteRook {
				return cr, fmt.Errorf("chess: fen illegal castle rights %s , normal white kingside missing pieces", castleStr)
			}
			cr.whiteKingSideCastle = true
		}
		if strings.ContainsRune(castleStr, 'Q') {
			if board.Piece(E1) != WhiteKing || board.Piece(A1) != WhiteRook {
				return cr, fmt.Errorf("chess: fen illegal castle rights %s, normal white queenside missing pieces", castleStr)
			}
			cr.whiteQueenSideCastle = true
		}
		if strings.ContainsRune(castleStr, 'k') {
			if board.Piece(E8) != BlackKing || board.Piece(H8) != BlackRook {
				return cr, fmt.Errorf("chess: fen illegal castle rights %s , normal black kingsie missing pieces", castleStr)
			}
			cr.blackKingSideCastle = true
		}
		if strings.ContainsRune(castleStr, 'q') {
			if board.Piece(E8) != BlackKing || board.Piece(A8) != BlackRook {
				return cr, fmt.Errorf("chess: fen illegal castle rights %s , normal black queenside missing pieces", castleStr)
			}
			cr.blackQueenSideCastle = true
		}
		return cr, nil
	} else { // 960 mode, handles both Shredder-FEN and X-FEN
		kingStartingFile := ""
		for _, c := range castleStr {
			if unicode.IsUpper(c) { // white castle
				if board.whiteKingSq.Rank() != Rank1 {
					return cr, fmt.Errorf("chess: fen illegal castle rights %s , white king should be on rank 1 for castle", castleStr)
				}
				if board.whiteKingSq.File() == FileA || board.whiteKingSq.File() == FileH {
					return cr, fmt.Errorf("chess: fen illegal castle rights %s , white king cant be on file A or H for castle in 960", castleStr)
				}
				if kingStartingFile != "" && kingStartingFile != board.whiteKingSq.File().String() {
					return cr, fmt.Errorf("chess: fen illegal castle rights %s , white and black kings must be on same file for both to have castle rights", castleStr)
				}
				kingStartingFile = board.whiteKingSq.File().String()
				if c == 'K' {
					if cr.whiteKingSideCastle {
						return cr, fmt.Errorf("chess: fen invalid castle rights %s , white king side castle info provided more than once", castleStr)
					}
					for sq := H1; sq > board.whiteKingSq; sq-- {
						if board.Piece(sq) == WhiteRook {
							if cr.hSideRookStartingFile != "" && cr.hSideRookStartingFile != sq.File().String() {
								return cr, fmt.Errorf("chess: fen invalid castle rights %s , rook starting king side file missmatch", castleStr)
							}
							cr.hSideRookStartingFile = sq.File().String()
							cr.whiteKingSideCastle = true
							break
						}
					}
					if !cr.whiteKingSideCastle {
						return cr, fmt.Errorf("chess: fen invalid castle rights %s , no white kingside rook found", castleStr)
					}
				} else if c == 'Q' {
					if cr.whiteQueenSideCastle {
						return cr, fmt.Errorf("chess: fen invalid castle rights %s , white queen side castle info provided more than once", castleStr)
					}
					for sq := A1; sq < board.whiteKingSq; sq++ {
						if board.Piece(sq) == WhiteRook {
							if cr.aSideRookStartingFile != "" && cr.aSideRookStartingFile != sq.File().String() {
								return cr, fmt.Errorf("chess: fen invalid castle rights %s , rook starting queen side file missmatch", castleStr)
							}
							cr.aSideRookStartingFile = sq.File().String()
							cr.whiteQueenSideCastle = true
							break
						}
					}
					if !cr.whiteQueenSideCastle {
						return cr, fmt.Errorf("chess: fen invalid castle rights %s , no white queenside rook found", castleStr)
					}
				} else if c >= 'A' && c <= 'H' {
					if runeToFileMap[c] > board.whiteKingSq.File() { // king side
						if cr.whiteKingSideCastle {
							return cr, fmt.Errorf("chess: fen invalid castle rights %s , white king side castle info provided more than once", castleStr)
						}
						if board.Piece(NewSquare(runeToFileMap[c], Rank1)) != WhiteRook {
							return cr, fmt.Errorf("chess: fen invalid castle rights %s , no white kingside rook found", castleStr)
						}
						if cr.hSideRookStartingFile != "" && cr.hSideRookStartingFile != runeToFileMap[c].String() {
							return cr, fmt.Errorf("chess: fen invalid castle rights %s , rook starting king side file missmatch", castleStr)
						}
						cr.hSideRookStartingFile = runeToFileMap[c].String()
						cr.whiteKingSideCastle = true
					} else if runeToFileMap[c] < board.whiteKingSq.File() { // queen side
						if cr.whiteQueenSideCastle {
							return cr, fmt.Errorf("chess: fen invalid castle rights %s , white queen side castle info provided more than once", castleStr)
						}
						if board.Piece(NewSquare(runeToFileMap[c], Rank1)) != WhiteRook {
							return cr, fmt.Errorf("chess: fen invalid castle rights %s , no white queenside rook found", castleStr)
						}
						if cr.aSideRookStartingFile != "" && cr.aSideRookStartingFile != runeToFileMap[c].String() {
							return cr, fmt.Errorf("chess: fen invalid castle rights %s , rook starting queen side file missmatch", castleStr)
						}
						cr.aSideRookStartingFile = runeToFileMap[c].String()
						cr.whiteQueenSideCastle = true
					} else {
						return cr, fmt.Errorf("chess: fen illegal castle rights %s , rook can't be on king", castleStr)
					}
				} else {
					return cr, fmt.Errorf("chess: fen invalid castle rights %s , unknown characters", castleStr)
				}
			} else if unicode.IsLower(c) { // black castle
				if board.blackKingSq.Rank() != Rank8 {
					return cr, fmt.Errorf("chess: fen illegal castle rights %s , black king should be on rank 8 for castle", castleStr)
				}
				if board.blackKingSq.File() == FileA || board.blackKingSq.File() == FileH {
					return cr, fmt.Errorf("chess: fen illegal castle rights %s , black king cant be on file A or H for castle in 960", castleStr)
				}
				if kingStartingFile != "" && kingStartingFile != board.blackKingSq.File().String() {
					return cr, fmt.Errorf("chess: fen illegal castle rights %s , white and black kings must be on same file for both to have castle rights", castleStr)
				}
				kingStartingFile = board.blackKingSq.File().String()
				if c == 'k' {
					if cr.blackKingSideCastle {
						return cr, fmt.Errorf("chess: fen invalid castle rights %s , black king side castle info provided more than once", castleStr)
					}
					for sq := H8; sq > board.blackKingSq; sq-- {
						if board.Piece(sq) == BlackRook {
							if cr.hSideRookStartingFile != "" && cr.hSideRookStartingFile != sq.File().String() {
								return cr, fmt.Errorf("chess: fen invalid castle rights %s , rook starting king side file missmatch", castleStr)
							}
							cr.hSideRookStartingFile = sq.File().String()
							cr.blackKingSideCastle = true
							break
						}
					}
					if !cr.blackKingSideCastle {
						return cr, fmt.Errorf("chess: fen invalid castle rights %s , no black kingside rook found", castleStr)
					}
				} else if c == 'q' {
					if cr.blackQueenSideCastle {
						return cr, fmt.Errorf("chess: fen invalid castle rights %s , black queen side castle info provided more than once", castleStr)
					}
					for sq := A8; sq < board.blackKingSq; sq++ {
						if board.Piece(sq) == BlackRook {
							if cr.aSideRookStartingFile != "" && cr.aSideRookStartingFile != sq.File().String() {
								return cr, fmt.Errorf("chess: fen invalid castle rights %s , rook starting queen side file missmatch", castleStr)
							}
							cr.aSideRookStartingFile = sq.File().String()
							cr.blackQueenSideCastle = true
							break
						}
					}
					if !cr.blackQueenSideCastle {
						return cr, fmt.Errorf("chess: fen invalid castle rights %s , no black queenside rook found", castleStr)
					}
				} else if c >= 'a' && c <= 'h' {
					if runeToFileMap[c] > board.blackKingSq.File() { // king side
						if cr.blackKingSideCastle {
							return cr, fmt.Errorf("chess: fen invalid castle rights %s , black king side castle info provided more than once", castleStr)
						}
						if board.Piece(NewSquare(runeToFileMap[c], Rank8)) != BlackRook {
							return cr, fmt.Errorf("chess: fen invalid castle rights %s , no black kingside rook found", castleStr)
						}
						if cr.hSideRookStartingFile != "" && cr.hSideRookStartingFile != runeToFileMap[c].String() {
							return cr, fmt.Errorf("chess: fen invalid castle rights %s , rook starting king side file missmatch", castleStr)
						}
						cr.hSideRookStartingFile = runeToFileMap[c].String()
						cr.blackKingSideCastle = true
					} else if runeToFileMap[c] < board.blackKingSq.File() { // queen side
						if cr.blackQueenSideCastle {
							return cr, fmt.Errorf("chess: fen invalid castle rights %s , black queen side castle info provided more than once", castleStr)
						}
						if board.Piece(NewSquare(runeToFileMap[c], Rank8)) != BlackRook {
							return cr, fmt.Errorf("chess: fen invalid castle rights %s , no black queenside rook found", castleStr)
						}
						if cr.aSideRookStartingFile != "" && cr.aSideRookStartingFile != runeToFileMap[c].String() {
							return cr, fmt.Errorf("chess: fen invalid castle rights %s , rook starting queen side file missmatch", castleStr)
						}
						cr.aSideRookStartingFile = runeToFileMap[c].String()
						cr.blackQueenSideCastle = true
					} else {
						return cr, fmt.Errorf("chess: fen illegal castle rights %s , rook can't be on king", castleStr)
					}
				} else {
					return cr, fmt.Errorf("chess: fen invalid castle rights %s , unknown character", castleStr)
				}
			} else {
				return cr, fmt.Errorf("chess: fen invalid castle rights %s , unknown character, niether lowercase nor uppercase", castleStr)
			}
		}
		return cr, nil
	}
}

func hasDuplicateCharacters(s string) bool {
	charMap := make(map[rune]bool)
	for _, char := range s {
		if charMap[char] {
			return true // Duplicate found
		}
		charMap[char] = true
	}
	return false
}
func hasOnlyKQkq(s string) bool {
	for _, char := range s {
		if !strings.ContainsRune("KQkq", char) {
			return false
		}
	}
	return true
}
func hasSomeKQkq(s string) bool {
	for _, char := range s {
		if strings.ContainsRune("KQkq", char) {
			return true
		}
	}
	return false
}

func formEnPassant(enPassant string, board *Board, turn Color) (Square, error) {
	if enPassant == "-" {
		return NoSquare, nil
	}
	sq := strToSquareMap[enPassant]
	if sq == NoSquare || !(sq.Rank() == Rank3 || sq.Rank() == Rank6) {
		return NoSquare, fmt.Errorf("chess: fen invalid En Passant square %s", enPassant)
	}

	if (sq.Rank() == Rank3 && turn != Black) || (sq.Rank() == Rank6 && turn != White) {
		return NoSquare, fmt.Errorf("chess: fen invalid En Passant square %s, not possible for given turn", enPassant)
	}

	if (sq.Rank() == Rank3 && board.Piece(NewSquare(sq.File(), Rank4)) != WhitePawn) || (sq.Rank() == Rank6 && board.Piece(NewSquare(sq.File(), Rank5)) != BlackPawn) {
		return NoSquare, fmt.Errorf("chess: fen invalid En Passant square %s, corresponding pawn not present", enPassant)
	}

	return sq, nil
}

var (
	fenPieceMap = map[string]Piece{
		"K": WhiteKing,
		"Q": WhiteQueen,
		"R": WhiteRook,
		"B": WhiteBishop,
		"N": WhiteKnight,
		"P": WhitePawn,
		"k": BlackKing,
		"q": BlackQueen,
		"r": BlackRook,
		"b": BlackBishop,
		"n": BlackKnight,
		"p": BlackPawn,
	}

	fenTurnMap = map[string]Color{
		"w": White,
		"b": Black,
	}
)
