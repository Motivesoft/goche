package uci

import (
	// Internal references
	"bufio"
	"fmt"
	"goche/logger"
	"goche/utility"
	"os"
	"strconv"
	"strings"
	"time"
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

	fmt.Printf("  Depth: %3d. Actual: %12d\n", depth, result)

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

	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	// Create a scanner
	scanner := bufio.NewScanner(file)

	// Read the file line by line
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comment lines - but print them as they might be in the file for formatting purposes
		if line == "" || strings.HasPrefix(line, "#") {
			fmt.Println(line)
			continue
		}

		perftFen(line, divide)
	}

	// Check for any errors during scanning
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	return nil
}

func perftFen(fenWithResults string, divide bool) error {

	// FEN format is expected to be one of:
	// - fen;Ddepth expected-at-depth;Ddepth expected-at-depth
	// - fen,expected-at-depth-1,expected-at-depth-2,expected-at-depth-3,...
	//
	// e.g:
	// - rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1;D1 20;D2 400;...
	// - rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1,20,400,...

	type expectedResults struct {
		depth     int
		moveCount int
	}

	var expected []expectedResults
	var fen string

	// Extract the expected results from the fen string
	split := strings.Split(fenWithResults, ";")
	if len(split) > 1 {
		fen = split[0]
		for i := 1; i < len(split); i++ {
			depth, count, err := getDepthAndExpected(split[i])
			if err != nil {
				return fmt.Errorf("badly formatted expected results: %s", fenWithResults)
			}

			expected = append(expected, expectedResults{depth, count})
		}
	} else {
		split = strings.Split(fenWithResults, ",")
		if len(split) > 1 {
			fen = split[0]
			for depth := 1; depth < len(split); depth++ {
				count, err := strconv.Atoi(split[depth])
				if err != nil {
					return fmt.Errorf("badly formatted expected results: %s", fenWithResults)
				}

				expected = append(expected, expectedResults{depth, count})
			}
		} else {
			return fmt.Errorf("missing expected results: %s", fenWithResults)
		}
	}

	// Now run the actual test
	fmt.Printf("FEN: %s\n", fen)
	for i := 0; i < len(expected); i++ {
		depth := expected[i].depth
		count := expected[i].moveCount

		logger.Debug("perft to depth %d with FEN: %s", depth, fen)

		result, err := perftRun(depth, fen, divide)
		if err != nil {
			return fmt.Errorf("run failed: %w", err)
		}

		fmt.Printf("  Depth: %3d. Expected: %12d. Actual: %12d. %s\n", depth, count, result, utility.If(count == result, "PASSED", "FAILED"))
	}

	return nil
}

// Parse 'Ddepth expected' (e.g. 'D2 1000') into depth and expected move count and returns an error on failure
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

func perftRun(depth int, fen string, divide bool) (int, error) {
	board, err := NewBoard(fen)
	if err != nil {
		return 0, fmt.Errorf("failed to create board: %w", err)
	}

	start := time.Now()

	var nodes int
	if divide {
		nodes, err = search(board, depth, divide)
	} else {
		nodes, err = search(board, depth, divide)
	}

	if err != nil {
		return 0, fmt.Errorf("move search failed: %w", err)
	}

	elapsed := time.Since(start)
	fmt.Printf("  Depth: %3d. Nodes: %3d. Time: %s\n", depth, nodes, elapsed)

	return nodes, nil
}

func search(board *Board, depth int, divide bool) (int, error) {
	nodes := 0

	// Always return 1 at the root - but I guess this could/should actually be
	// a check of a pseudo-level moves if our engine return such things?
	if depth == 0 {
		return 1, nil
	}

	// Create an array for possible moves and allocate it to the maximum size necessary
	moveList := make([]Move, 0, 256)

	// Generate all possible moves
	moveList, err := board.GetMoves(moveList)
	if err != nil {
		return 0, fmt.Errorf("move generation failed: %w", err)
	}

	// Recurse with each move
	for i := 0; i < len(moveList); i++ {
		move := moveList[i]

		// Make the move, search the new position, unmake the move - repeat

		undo := board.MakeMove(move)

		moveNodes, err := search(board, depth-1, false)
		if err != nil {
			return 0, err
		}

		nodes += moveNodes

		// Extra reporting if requested
		if divide {
			fmt.Printf("  %3d : %d : %p\n", move, moveNodes, board)
		}

		board.UnmakeMove(undo)
	}

	return nodes, nil
}
