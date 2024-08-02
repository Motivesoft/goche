package uci

import (
	"fmt"
	"goche/utility"

	"strconv"
	"strings"
	"unicode"
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

	FullMoveLSB = 0b00000000010000000000000000000000
	HalfMoveLSB = 0b00000000000000000001000000000000
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

	// We are going to assume that the FEN string is more or less well formed, not least because we
	// expect that it is being provided by another piece of software that also conforms to UCI

	currentComponent := piecePlacement
	remainder := fen

	var component string
	for {
		component, remainder = utility.SplitNextWord(remainder)

		switch currentComponent {
		case piecePlacement:
			// Piece placement starts on the 8th rank, 1st file
			reminderBit := indexToBitboard(squareToIndex[uint32]("a8"))
			bitboardBit := reminderBit

			// Parse the piece placement section of the FEN string
			for _, character := range component {
				if character == '/' {
					// Line break - move to next rank (towards LSB)
					reminderBit >>= 8
					bitboardBit = reminderBit
					continue
				}

				// Empty square(s)
				if unicode.IsDigit(character) {
					bitboardBit <<= uint64(character - '0')
					continue
				}

				// Piece color
				if unicode.IsLower(character) {
					board.blackPieces |= bitboardBit
				} else if unicode.IsUpper(character) {
					board.whitePieces |= bitboardBit
				}

				// Piece type
				switch unicode.ToUpper(character) {
				case 'P':
					board.pawns |= bitboardBit

				case 'N':
					board.knights |= bitboardBit

				case 'B':
					board.bishops |= bitboardBit

				case 'R':
					board.rooks |= bitboardBit

				case 'Q':
					board.queens |= bitboardBit

				case 'K':
					board.kings |= bitboardBit
				}

				bitboardBit <<= 1
			}

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
			if component == "-" {
				board.clearEnPassantIndex()
			} else {
				value := squareToIndex[uint32](component)
				board.setEnPassantIndex(value)
			}

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

	// TODO Remove this debug output
	board.printBoard()
	board.printMasks()

	return board, nil
}

func (board *Board) GetMoves(moveList []Move) ([]Move, error) {
	// moveList = append(moveList, ...)

	// OK, so how are we going to do this:
	// - setup directional variables for whether white or black, or copy the code for each, or something clever
	// - pseudo legal or fully legal
	// - return an array or callback per move
	// - ordered move list (captures and checks, ...)
	// - attack/defend maps (or whatever they are)

	// Basic implementation for now
	// Source and Target for player and opponent
	sourceMask := utility.If(board.gameState&WhiteMask == WhiteMask, board.whitePieces, board.blackPieces)
	targetMask := utility.If(board.gameState&WhiteMask == WhiteMask, board.blackPieces, board.whitePieces)

	// Generate all possible moves
	// moveList = append(moveList, board.generatePawnMoves(sourceMask, targetMask)...)
	moveList = board.generateKnightMoves(moveList, sourceMask, targetMask)
	// moveList = append(moveList, board.generateBishopMoves(sourceMask, targetMask)...)
	// moveList = append(moveList, board.generateRookMoves(sourceMask, targetMask)...)
	// moveList = append(moveList, board.generateQueenMoves(sourceMask, targetMask)...)
	// moveList = append(moveList, board.generateKingMoves(sourceMask, targetMask)...)

	return moveList, nil
}

func (b *Board) generateKnightMoves(moveList []Move, sourceMask uint64, _ uint64) []Move {
	pieceSet := b.knights
	pieceSet &= utility.If(b.gameState&WhiteMask == WhiteMask, b.whitePieces, b.blackPieces)

	var pieceIndex int
	var targetIndex int
	for bitScanReverse(&pieceIndex, pieceSet) {
		pieceSet ^= 1 << pieceIndex

		targetSquares := PieceMoveMasks.KnightMoveMask[pieceIndex]
		for bitScanReverse(&targetIndex, targetSquares) {
			targetSquares ^= 1 << targetIndex

			if (1<<targetIndex)&sourceMask == 0 {
				moveList = append(moveList, NewMove(uint16(pieceIndex), uint16(targetIndex)))
			}
		}
	}

	return moveList
}

func (b *Board) MakeMove(move Move) *Board {
	// Copy the board state
	backupBoard := *b
	return &backupBoard
}

func (b *Board) UnmakeMove(backupBoard *Board) {
	// Restore the board state
	b.whitePieces = backupBoard.whitePieces
	b.blackPieces = backupBoard.blackPieces
	b.pawns = backupBoard.pawns
	b.knights = backupBoard.knights
	b.bishops = backupBoard.bishops
	b.rooks = backupBoard.rooks
	b.queens = backupBoard.queens
	b.kings = backupBoard.kings
	b.gameState = backupBoard.gameState
}

func (b *Board) getFullMoveNumber() uint32 {
	return (b.gameState & FullMoveMask) >> FullMoveShift
}

func (b *Board) setFullMoveNumber(number uint32) {
	// Mask out the current value, shift the new number, and it to size and or it back into the state
	// The and-to-size step should be unnecessary, but at least means than it prevents a rogue value impacting
	// any other state bits
	b.gameState = (b.gameState & FullMoveXOR) | ((number << FullMoveShift) & FullMoveMask)
}

func (b *Board) incrementFullMoveNumber() {
	// Increment the number in situ
	b.gameState = (b.gameState & FullMoveXOR) | (((b.gameState & FullMoveMask) + FullMoveLSB) & FullMoveMask)
}

func (b *Board) getHalfMoveClock() uint32 {
	return (b.gameState & HalfMoveMask) >> HalfMoveShift
}

func (b *Board) setHalfMoveClock(number uint32) {
	// Mask out the current value, shift the new number, and it to size and or it back into the state
	// The and-to-size step should be unnecessary, but at least means than it prevents a rogue value impacting
	// any other state bits
	b.gameState = (b.gameState & HalfMoveXOR) | ((number << HalfMoveShift) & HalfMoveMask)
}

func (b *Board) incrementHalfMoveClock() {
	b.gameState = (b.gameState & HalfMoveXOR) | ((b.gameState & HalfMoveMask) + HalfMoveLSB)
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

func (b *Board) printBoard() {
	fmt.Println("  ABCDEFGH")
	fmt.Println("  --------")

	for row := 0; row < 8; row++ {
		rank := 8 - row
		fmt.Printf("%d|", rank)

		for column := 0; column < 8; column++ {
			file := 'a' + column

			bitboardBit := indexToBitboard(squareToIndex[uint32](fmt.Sprintf("%c%d", file, rank)))

			var piece string
			if b.pawns&bitboardBit != 0 {
				piece = utility.If(b.whitePieces&bitboardBit == bitboardBit, "P", "p")
			} else if b.knights&bitboardBit != 0 {
				piece = utility.If(b.whitePieces&bitboardBit == bitboardBit, "N", "n")
			} else if b.bishops&bitboardBit != 0 {
				piece = utility.If(b.whitePieces&bitboardBit == bitboardBit, "B", "b")
			} else if b.rooks&bitboardBit != 0 {
				piece = utility.If(b.whitePieces&bitboardBit == bitboardBit, "R", "r")
			} else if b.queens&bitboardBit != 0 {
				piece = utility.If(b.whitePieces&bitboardBit == bitboardBit, "Q", "q")
			} else if b.kings&bitboardBit != 0 {
				piece = utility.If(b.whitePieces&bitboardBit == bitboardBit, "K", "k")
			} else {
				piece = " "
			}

			if (rank+file)&1 == 0 {
				fmt.Printf("\033[40;1m%s\033[0m", piece)
			} else {
				fmt.Printf("\033[47;1m%s\033[0m", piece)
			}
		}
		fmt.Printf("|%d\n", rank)
	}

	fmt.Println("  --------")
	fmt.Println("  ABCDEFGH")
	fmt.Println()

	if b.gameState&WhiteMask == WhiteMask {
		fmt.Printf("White to play\n")
	} else {
		fmt.Printf("White to play\n")
	}

	fmt.Print("Castling rights:   ")
	if b.gameState&CastlingMask_WK == CastlingMask_WK {
		fmt.Printf("K")
	}
	if b.gameState&CastlingMask_WQ == CastlingMask_WQ {
		fmt.Printf("Q")
	}
	if b.gameState&CastlingMask_BK == CastlingMask_BK {
		fmt.Printf("k")
	}
	if b.gameState&CastlingMask_BQ == CastlingMask_BQ {
		fmt.Printf("q")
	}
	fmt.Println()

	fmt.Printf("En passant square: ")
	if b.getEnPassantIndex() != 0 {
		fmt.Printf("%s\n", indexToSquare[uint32](b.getEnPassantIndex()))
	} else {
		fmt.Printf("[none]\n")
	}

	fmt.Printf("Half move clock:   %d\n", b.getHalfMoveClock())
	fmt.Printf("Full move number:  %d\n", b.getFullMoveNumber())
	fmt.Println()
}

func (b *Board) printMasks() {
	fmt.Printf("Pawns:      	   %064b\n", b.pawns)
	fmt.Printf("Knights:    	   %064b\n", b.knights)
	fmt.Printf("Bishops:    	   %064b\n", b.bishops)
	fmt.Printf("Rooks:      	   %064b\n", b.rooks)
	fmt.Printf("Queens:     	   %064b\n", b.queens)
	fmt.Printf("Kings:      	   %064b\n", b.kings)
	fmt.Printf("White:      	   %064b\n", b.whitePieces)
	fmt.Printf("Black:      	   %064b\n", b.blackPieces)
	fmt.Printf("All:        	   %064b\n", b.whitePieces|b.blackPieces)
	fmt.Println()
	fmt.Printf("Game State:        %032b\n", b.gameState)
	fmt.Printf("Color To Play:     %032b\n", b.gameState&(WhiteMask|BlackMask))
	fmt.Printf("Castling Rights:   %032b\n", b.gameState&(CastlingMask_WK|CastlingMask_WQ|CastlingMask_BK|CastlingMask_BQ))
	fmt.Printf("En Passant Square: %032b\n", b.gameState&EnPassantMask)
	fmt.Printf("Half Move Clock:   %032b\n", b.gameState&HalfMoveMask)
	fmt.Printf("Full Move Number:  %032b\n", b.gameState&FullMoveMask)
	fmt.Println()
}
