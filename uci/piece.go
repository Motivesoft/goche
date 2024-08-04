package uci

type PieceMoveMask struct {
	WhitePawnSlideMask   [64]uint64
	WhitePawnCaptureMask [64]uint64
	BlackPawnSlideMask   [64]uint64
	BlackPawnCaptureMask [64]uint64
	KnightMoveMask       [64]uint64
	DiagonalMoveMask     [64]uint64
	StraightMoveMask     [64]uint64
	QueenMoveMask        [64]uint64
	KingMoveMask         [64]uint64
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

			// Knight moves
			_ = setIfOnBoard(&PieceMoveMasks.KnightMoveMask[squareIndex], fileIndex-2, rankIndex-1)
			_ = setIfOnBoard(&PieceMoveMasks.KnightMoveMask[squareIndex], fileIndex+2, rankIndex-1)
			_ = setIfOnBoard(&PieceMoveMasks.KnightMoveMask[squareIndex], fileIndex-2, rankIndex+1)
			_ = setIfOnBoard(&PieceMoveMasks.KnightMoveMask[squareIndex], fileIndex+2, rankIndex+1)
			_ = setIfOnBoard(&PieceMoveMasks.KnightMoveMask[squareIndex], fileIndex-1, rankIndex-2)
			_ = setIfOnBoard(&PieceMoveMasks.KnightMoveMask[squareIndex], fileIndex+1, rankIndex-2)
			_ = setIfOnBoard(&PieceMoveMasks.KnightMoveMask[squareIndex], fileIndex-1, rankIndex+2)
			_ = setIfOnBoard(&PieceMoveMasks.KnightMoveMask[squareIndex], fileIndex+1, rankIndex+2)

			// Directional masks forBishop/Rook/Queen moves
			for r := -1; r <= 1; r++ {
				for f := -1; f <= 1; f++ {
					// Straight or diagonal? In either case, go as far as we can and then break
					if f == 0 || r == 0 {
						for d := 1; d < 8; d++ {
							if !setIfOnBoard(&PieceMoveMasks.StraightMoveMask[squareIndex], fileIndex+f*d, rankIndex+r*d) {
								break
							}
						}
					} else {
						for d := 1; d < 8; d++ {
							if !setIfOnBoard(&PieceMoveMasks.DiagonalMoveMask[squareIndex], fileIndex+f*d, rankIndex+r*d) {
								break
							}
						}
					}
				}
			}

			// King moves
			for r := -1; r <= 1; r++ {
				for f := -1; f <= 1; f++ {
					if f == 0 && r == 0 {
						continue
					}
					_ = setIfOnBoard(&PieceMoveMasks.KingMoveMask[squareIndex], fileIndex+f, rankIndex+r)
				}
			}

			// Pawn moves - set for both colors
			if rankIndex > 0 {
				// White
				if rankIndex < 7 {
					_ = setIfOnBoard(&PieceMoveMasks.WhitePawnSlideMask[squareIndex], fileIndex, rankIndex+1)
				}

				if rankIndex == 1 {
					_ = setIfOnBoard(&PieceMoveMasks.WhitePawnSlideMask[squareIndex], fileIndex, rankIndex+2)
				}

				_ = setIfOnBoard(&PieceMoveMasks.WhitePawnCaptureMask[squareIndex], fileIndex-1, rankIndex+1)
				_ = setIfOnBoard(&PieceMoveMasks.WhitePawnCaptureMask[squareIndex], fileIndex+1, rankIndex+1)

				// Black
				if rankIndex > 0 {
					_ = setIfOnBoard(&PieceMoveMasks.BlackPawnSlideMask[squareIndex], fileIndex, rankIndex-1)
				}

				if rankIndex == 6 {
					_ = setIfOnBoard(&PieceMoveMasks.BlackPawnSlideMask[squareIndex], fileIndex, rankIndex-2)
				}

				_ = setIfOnBoard(&PieceMoveMasks.BlackPawnCaptureMask[squareIndex], fileIndex-1, rankIndex+1)
				_ = setIfOnBoard(&PieceMoveMasks.BlackPawnCaptureMask[squareIndex], fileIndex+1, rankIndex+1)
			}
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
