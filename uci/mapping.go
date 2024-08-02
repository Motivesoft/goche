package uci

import (
	"fmt"

	"math/bits"
)

/*
func squareToIndex(square string) uint8 {
	return uint8((square[0] - 'a') + ((square[1] - '1') * 8))
}

func squareToIndex32(square string) uint32 {
	return uint32((square[0] - 'a') + ((square[1] - '1') * 8))
}

func indexToSquare(index uint8) string {
	return fmt.Sprintf("%c%c", 'a'+(index%8), '1'+(index/8))
}

func index32ToSquare(index uint32) string {
	return fmt.Sprintf("%c%c", 'a'+(index%8), '1'+(index/8))
}
*/

type NumberType interface {
	uint8 | uint32 | uint64
}

func squareToIndex[numberType NumberType](square string) numberType {
	return rankFileToIndex[numberType](square[0]-'a', square[1]-'1')
	//return numberType((square[0] - 'a') + ((square[1] - '1') * 8))
}

func rankFileToIndex[numberType NumberType](file byte, rank byte) numberType {
	return numberType(file + rank*8)
}

func indexToSquare[numberType NumberType](index numberType) string {
	return fmt.Sprintf("%c%c", 'a'+(index%8), '1'+(index/8))
}

func indexToBitboard[numberType NumberType](index numberType) uint64 {
	return 1 << index
}

/*
func bitboardToIndex[numberType NumberType](bitboard uint64) numberType {
	// TODO get this from a lookup table or a bitshift
	return 0
}
*/

func bitScanForward(index *int, mask uint64) bool {
	if mask == 0 {
		return false
	}

	*index = bits.TrailingZeros64(mask)
	return true
}

func bitScanReverse(index *int, mask uint64) bool {
	if mask == 0 {
		return false
	}

	*index = bits.Len64(mask) - 1
	return true
}
