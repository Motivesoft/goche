package uci

import "fmt"

// Move requires:
// - from square (6 bits for index)
// - to square (6 bits for index)
// - promotion piece (2 bits for knight, bishop, rook, queen)
// - 2 bits spare, assuming we don't need them for either promotion, castling, enpassant, check, capture, checkmate...
// TODO Start with 16 bits and then extend later?

type Move uint16

func NewPromotionMove(from, to uint16, promotionPiece uint16) Move {
	return Move(from | (to << 6) | (promotionPiece << 12))
}

func NewMove(from, to uint16) Move {
	return Move(from | (to << 6))
}

func (m Move) From() uint8 {
	return uint8(m & 0b111111)
}

func (m Move) To() uint8 {
	return uint8((m >> 6) & 0b111111)
}

func (m Move) PrintMove() {
	from := m.From()
	to := m.To()
	fmt.Printf("%016b %c%c%c%c\n", m, 'a'+from%8, '1'+from/8, 'a'+to%8, '1'+to/8)
}
