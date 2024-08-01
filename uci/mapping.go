package uci

import "fmt"

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
	uint8 | uint32
}

func squareToIndex[numberType NumberType](square string) numberType {
	return numberType((square[0] - 'a') + ((square[1] - '1') * 8))
}

func indexToSquare[numberType NumberType](index numberType) string {
	return fmt.Sprintf("%c%c", 'a'+(index%8), '1'+(index/8))
}
