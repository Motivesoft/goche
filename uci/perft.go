package uci

import (
	// Internal references
	"fmt"
	"goche/logger"
	"goche/utility"
	"os"
	"strconv"
	"strings"
)

// Perform a search to a depth with a FEN string and display the results
func PerftDepth(depth int, fen string, divide bool) error {
	if fen == "" {
		return fmt.Errorf("missing FEN string")
	}

	logger.Debug("perft to depth %d with FEN: %s", depth, fen)

	result, err := perftRun(depth, fen, divide)
	if err != nil {
		return fmt.Errorf("run failed: %w", err)
	}

	fmt.Printf("  Depth: %d. Actual: %d\n", depth, result)

	return nil
}

// Process a FEN string that contains expected results (error if not)
func PerftWithFen(fen string, divide bool) error {
	if fen == "" {
		return fmt.Errorf("missing FEN string")
	}

	logger.Debug("perft with FEN: %s", fen)

	return perftFen(fen, divide)
}

// Read and process a file of FEN strings, each of which contain expected results (error if not)
func PerftWithFile(filename string, divide bool) error {
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

func perftRun(depth int, fen string, divide bool) (int, error) {
	//return 0, fmt.Errorf("perft not implemented")
	return 0, nil
}

func perftFen(fen string, divide bool) error {

	// FEN format is expected to be one of:
	// - fen;Ddepth expected-at-depth;Ddepth expected-at-depth
	// - fen,expected-at-depth-1,expected-at-depth-2,expected-at-depth-3,...
	//
	// e.g:
	// - rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1;D1 20;D2 400;...
	// - rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1,20,400,...

	type expectedResults struct {
		depth, expected int
	}

	//	var expected []expectedResults

	split := strings.Split(fen, ";")
	if len(split) > 1 {
		for i := 1; i < len(split); i++ {
			depth, count, err := getDepthAndExpected(split[i])
			if err != nil {
				return fmt.Errorf("badly formatted expected results: %s", fen)
			}

			//			expected = append(expected, expectedResults{depth, count})
			logger.Debug("perft to depth %d with FEN: %s", depth, fen)

			result, err := perftRun(depth, fen, divide)
			if err != nil {
				return fmt.Errorf("run failed: %w", err)
			}

			fmt.Printf("  Depth: %3d. Expected: %12d. Actual: %12d. %s\n", depth, count, result, utility.If(count == result, "PASSED", "FAILED"))
		}
	} else {
		split = strings.Split(fen, ",")
		if len(split) > 1 {
			for i := 1; i < len(split); i++ {
				count, err := strconv.Atoi(split[i])
				if err != nil {
					return fmt.Errorf("badly formatted expected results: %s", fen)
				}

				logger.Debug("perft to depth %d with FEN: %s", i, fen)

				result, err := perftRun(i, fen, divide)
				if err != nil {
					return fmt.Errorf("run failed: %w", err)
				}

				fmt.Printf("  Depth: %3d. Expected: %12d. Actual: %12d. %s\n", i, count, result, utility.If(count == result, "PASSED", "FAILED"))
			}
		} else {
			return fmt.Errorf("missing expected results: %s", fen)
		}
	}

	return nil
}

func performRun(depth int, fen string) (int, error) {

}

func getDepthAndExpected(expected string) (int, int, error) {
	if !strings.HasPrefix(expected, "D") {
		return 0, 0, fmt.Errorf("expected results must start with 'D'")
	}

	numbers := strings.Split(expected[1:], " ")
	if len(numbers) != 2 {
		return 0, 0, fmt.Errorf("expected results require depth and expected move count values")
	}

	depth, err := strconv.Atoi(numbers[0])
	if err != nil {
		return 0, 0, fmt.Errorf("depth must be a number")
	}

	count, err := strconv.Atoi(numbers[1])
	if err != nil {
		return 0, 0, fmt.Errorf("expected move count must be a number")
	}

	return depth, count, nil
}
