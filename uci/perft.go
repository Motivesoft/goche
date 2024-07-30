package uci

import (
	// Internal references
	"fmt"
	"goche/logger"
	"os"
)

func PerftDepth(depth int, fen string) error {
	if fen == "" {
		return fmt.Errorf("missing FEN string")
	}

	logger.Debug("perft to depth %d with FEN: %s", depth, fen)

	return nil
}

func PerftWithFen(fen string) error {
	if fen == "" {
		return fmt.Errorf("missing FEN string")
	}

	logger.Debug("perft with FEN: %s", fen)

	return nil
}

func PerftWithFile(filename string) error {
	if filename == "" {
		return fmt.Errorf("missing filename")
	}

	logger.Debug("perft from file: %s", filename)

	_, err := os.Open(filename)
	if err != nil {
		return err
	}

	return nil
}
