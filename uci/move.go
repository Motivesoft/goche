package uci

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
