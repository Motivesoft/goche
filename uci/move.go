package uci

// Move requires:
// - from square (6 bits for index)
// - to square (6 bits for index)
// - promotion piece (2 bits for knight, bishop, rook, queen)
// - 2 bits spare, assuming we don't need them for either promotion, castling, enpassant, check, capture, checkmate...
// TODO Start with 16 bits and then extend later?

type Move uint16
