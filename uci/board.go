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

func (board *Board) GetMoves(moveList []uint) ([]uint, error) {
	moveList = append(moveList, 1, 2, 3, 4, 5, 6, 7, 8)
	moveList = append(moveList, 1, 2, 3, 4, 5, 6, 7, 8)
	moveList = append(moveList, 1, 2, 3, 4, 5, 6, 7, 8)
	moveList = append(moveList, 1, 2, 3, 4, 5, 6, 7, 8)
	moveList = append(moveList, 1, 2, 3, 4, 5, 6, 7, 8)
	moveList = append(moveList, 1, 2, 3, 4, 5, 6, 7, 8)
	moveList = append(moveList, 1, 2, 3, 4, 5, 6, 7, 8)
	moveList = append(moveList, 1, 2, 3, 4, 5, 6, 7, 8)
	moveList = append(moveList, 1, 2, 3, 4, 5, 6, 7, 8)
	moveList = append(moveList, 1, 2, 3, 4, 5, 6, 7, 8)
	moveList = append(moveList, 1, 2, 3, 4, 5, 6, 7, 8)
	moveList = append(moveList, 1, 2, 3, 4, 5, 6, 7, 8)
	moveList = append(moveList, 1, 2, 3, 4, 5, 6, 7, 8)
	moveList = append(moveList, 1, 2, 3, 4, 5, 6, 7, 8)
	moveList = append(moveList, 1, 2, 3, 4, 5, 6, 7, 8)
	moveList = append(moveList, 1, 2, 3, 4, 5, 6, 7, 8)
	moveList = append(moveList, 1, 2, 3, 4, 5, 6, 7, 8)
	moveList = append(moveList, 1, 2, 3, 4, 5, 6, 7, 8)
	moveList = append(moveList, 1, 2, 3, 4, 5, 6, 7, 8)
	moveList = append(moveList, 1, 2, 3, 4, 5, 6, 7, 8)
	moveList = append(moveList, 1, 2, 3, 4, 5, 6, 7, 8)
	moveList = append(moveList, 1, 2, 3, 4, 5, 6, 7, 8)
	moveList = append(moveList, 1, 2, 3, 4, 5, 6, 7, 8)
	moveList = append(moveList, 1, 2, 3, 4, 5, 6, 7, 8)
	moveList = append(moveList, 1, 2, 3, 4, 5, 6, 7, 8)
	moveList = append(moveList, 1, 2, 3, 4, 5, 6, 7, 8)
	moveList = append(moveList, 1, 2, 3, 4, 5, 6, 7, 8)
	moveList = append(moveList, 1, 2, 3, 4, 5, 6, 7, 8)
	moveList = append(moveList, 1, 2, 3, 4, 5, 6, 7, 8)
	moveList = append(moveList, 1, 2, 3, 4, 5, 6, 7, 8)
	moveList = append(moveList, 1, 2, 3, 4, 5, 6, 7, 8)
	moveList = append(moveList, 1, 2, 3, 4, 5, 6, 7)
	return moveList, nil
}

func (b *Board) MakeMove(move uint) string {
	return ""
}

func (b *Board) UnmakeMove(undo string) {
}
