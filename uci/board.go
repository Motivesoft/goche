package uci

import (
	"fmt"
	"goche/utility"
	"strconv"
	"strings"
)

// Internal references

const (
	FullMoveMask    uint32 = 0b11111111110000000000000000000000
	HalfMoveMask    uint32 = 0b00000000001111111111000000000000
	EnPassantMask   uint32 = 0b00000000000000000000111111000000
	CastlingMask_WK uint32 = 0b00000000000000000000000000100000
	CastlingMask_WQ uint32 = 0b00000000000000000000000000010000
	CastlingMask_BK uint32 = 0b00000000000000000000000000001000
	CastlingMask_BQ uint32 = 0b00000000000000000000000000000100
	BlackMask       uint32 = 0b00000000000000000000000000000010
	WhiteMask       uint32 = 0b00000000000000000000000000000001
)

const (
	FullMoveXOR  uint32 = 0b00000000001111111111111111111111
	HalfMoveXOR  uint32 = 0b11111111110000000000111111111111
	EnPassantXOR uint32 = 0b11111111111111111111000000111111
)

const (
	FullMoveShift  = 22
	HalfMoveShift  = 12
	EnPassantShift = 6
)

type Board struct {
	// Partly made up of elements recognisable from FEN and the rest is transient
	blackPieces uint64
	whitePieces uint64
	pawns       uint64
	knights     uint64
	bishops     uint64
	rooks       uint64
	queens      uint64
	kings       uint64
	gameState   uint32

	// gameState is designed to fit into a 32-bit uint, from lsb to msb:
	//  2 bits  whose turn it is (01 = white, 10 = black. Two bits in case this buys some other advantage)
	//  4 bits castling			 (1111 for KQkq)
	//  6 bits enpassant		 (1<<x == uint64 square, 0 for none as that is not a valid enpassant square)
	// 10 bits halfmove			 (enough for 1024, which should be enough)
	// 10 bits fullmove			 (enough for 1024, which should be enough)
}

func NewBoard(fen string) (*Board, error) {
	if fen == "" {
		return nil, fmt.Errorf("missing FEN string")
	}

	// Create the object and then populate it from the FEN string
	board := &Board{}

	type FenComponent int

	const (
		piecePlacement FenComponent = iota
		activeColor
		castlingRights
		enPassantSquare
		halfMoveClock
		fullMoveNumber
	)

	// We are going to assume that the FEN string is well formed, not least because we can reasonably
	// assume that it is being called by another piece of software that also conforms to UCI

	currentComponent := piecePlacement
	remainder := fen

	for {
		component, remainder := utility.SplitNextWord(remainder)
		switch currentComponent {
		case piecePlacement:
		case activeColor:
			board.gameState |= utility.If(strings.Contains(component, "w"), WhiteMask, BlackMask)

		case castlingRights:
			if strings.Contains(component, "K") {
				board.gameState |= CastlingMask_WK
			}
			if strings.Contains(component, "Q") {
				board.gameState |= CastlingMask_WQ
			}
			if strings.Contains(component, "k") {
				board.gameState |= CastlingMask_BK
			}
			if strings.Contains(component, "q") {
				board.gameState |= CastlingMask_BQ
			}

		case enPassantSquare:
			value := squareToIndex[uint32](component)
			board.setEnPassantIndex(value)

		case halfMoveClock:
			value, _ := strconv.Atoi(component)
			board.setHalfMoveClock(uint32(value))

		case fullMoveNumber:
			value, _ := strconv.Atoi(component)
			board.setFullMoveNumber(uint32(value))
		}

		// Done?
		if currentComponent == fullMoveNumber {
			break
		}

		// Step onto the next item
		currentComponent++
		if remainder == "" {
			return nil, fmt.Errorf("unexpected end of FEN string")
		}
	}

	return board, nil
}

func (board *Board) GetMoves(moveList []uint) ([]uint, error) {
	// moveList = append(moveList, ...)

	return moveList, nil
}

func (b *Board) MakeMove(move uint) string {
	return ""
}

func (b *Board) UnmakeMove(undo string) {
}

func (b *Board) getFullMoveNumber() uint32 {
	return (b.gameState % FullMoveMask) >> FullMoveShift
}

func (b *Board) setFullMoveNumber(number uint32) {
	b.gameState = (b.gameState & FullMoveXOR) | (number << FullMoveShift)
}

func (b *Board) incrementFullMoveNumber() {
	number := (b.gameState % FullMoveMask) >> FullMoveShift
	number++
	b.gameState = (b.gameState & FullMoveXOR) | (number << FullMoveShift)
}

func (b *Board) getHalfMoveClock() uint32 {
	return (b.gameState % HalfMoveMask) >> HalfMoveShift
}

func (b *Board) setHalfMoveClock(number uint32) {
	b.gameState = (b.gameState & HalfMoveXOR) | (number << HalfMoveShift)
}

func (b *Board) incrementHalfMoveClock() {
	number := (b.gameState % HalfMoveMask) >> HalfMoveShift
	number++
	b.gameState = (b.gameState & HalfMoveXOR) | (number << HalfMoveShift)
}

func (b *Board) getEnPassantIndex() uint32 {
	return (b.gameState & EnPassantMask) >> EnPassantShift
}

func (b *Board) setEnPassantIndex(index uint32) {
	b.gameState = (b.gameState & EnPassantXOR) | (index << EnPassantShift)
}

func (b *Board) clearEnPassantIndex() {
	b.gameState = (b.gameState & EnPassantXOR)
}
