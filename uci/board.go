package uci

import (
	"fmt"
)

// Internal references

type Color uint8

const (
	White Color = 0b10000000
	Black Color = 0b01000000
)

type Board struct {
	toPlay Color
}

func NewBoard(fen string) (*Board, error) {
	if fen == "" {
		return nil, fmt.Errorf("missing FEN string")
	}

	return &Board{
		toPlay: White,
	}, nil
}

func (board *Board) ToPlay() Color {
	return board.toPlay
}
