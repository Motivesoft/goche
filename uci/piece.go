package uci

type PieceMoveMask struct {
	PawnMoveMask     [64]uint64
	PawnSlideMask    [64]uint64
	PawnCaptureMask  [64]uint64
	KnightMoveMask   [64]uint64
	DiagonalMoveMask [64]uint64
	StraightMoveMask [64]uint64
	QueenMoveMask    [64]uint64
	KingMoveMask     [64]uint64
}

// 64-bit constant masks using this template:
//
//	0b0000000000000000000000000000000000000000000000000000000000000000
const (
	WhiteSideOfTheBoardMask       = 0b0000000000000000000000000000000011111111111111111111111111111111
	WhitePawnPromotionMask        = 0b1111111100000000000000000000000000000000000000000000000000000000
	WhitePawnSlideEligibilityMask = 0b0000000000000000000000000000000000000000000000001111111100000000

	BlackSideOfTheBoardMask       = 0b1111111111111111111111111111111100000000000000000000000000000000
	BlackPawnPromotionMask        = 0b0000000000000000000000000000000000000000000000000000000011111111
	BlackPawnSlideEligibilityMask = 0b0000000011111111000000000000000000000000000000000000000000000000
)

var PieceMoveMasks PieceMoveMask

func init() {
	for rankIndex := 0; rankIndex < 8; rankIndex++ {
		for fileIndex := 0; fileIndex < 8; fileIndex++ {
			squareIndex := rankIndex*8 + fileIndex

			_ = setIfOnBoard(&PieceMoveMasks.KnightMoveMask[squareIndex], fileIndex-2, rankIndex-1)
			_ = setIfOnBoard(&PieceMoveMasks.KnightMoveMask[squareIndex], fileIndex+2, rankIndex-1)
			_ = setIfOnBoard(&PieceMoveMasks.KnightMoveMask[squareIndex], fileIndex-2, rankIndex+1)
			_ = setIfOnBoard(&PieceMoveMasks.KnightMoveMask[squareIndex], fileIndex+2, rankIndex+1)
			_ = setIfOnBoard(&PieceMoveMasks.KnightMoveMask[squareIndex], fileIndex-1, rankIndex-2)
			_ = setIfOnBoard(&PieceMoveMasks.KnightMoveMask[squareIndex], fileIndex+1, rankIndex-2)
			_ = setIfOnBoard(&PieceMoveMasks.KnightMoveMask[squareIndex], fileIndex-1, rankIndex+2)
			_ = setIfOnBoard(&PieceMoveMasks.KnightMoveMask[squareIndex], fileIndex+1, rankIndex+2)
		}
	}
}

func setIfOnBoard(bitboard *uint64, destinationFile int, destinationRank int) bool {
	if destinationFile >= 0 && destinationFile < 8 && destinationRank >= 0 && destinationRank < 8 {
		*bitboard |= 1 << (destinationRank*8 + destinationFile)
		return true
	}

	return false
}
