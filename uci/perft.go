package uci

import (
	// Internal references
	"goche/logger"
	"os"
)

func PerftDepth(depth int, fen string) error {
	logger.Debug("perft to depth %d with FEN: %s", depth, fen)

	return nil
}

func PerftWithFen(fen string) error {
	logger.Debug("perft with FEN: %s", fen)

	return nil
}

func PerftWithFile(filename string) error {
	logger.Debug("perft from file: %s", filename)

	_, err := os.Open(filename)
	if err != nil {
		return err
	}

	return nil
}
