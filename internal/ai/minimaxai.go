package ai

import (
	"fmt"
	"math"

	"uac.tech/connect4/internal/game"
)

type TurnType int

const (
	Maximizing TurnType = iota
	Minimizing
)

type MinimaxAI struct {
	Depth int
	Self  game.Cell
	Other game.Cell
	Debug bool
}

const WIN_SCORE = 1_000_000

func (ai *MinimaxAI) GetMove(board *game.Board) int {

	var shadowBoard = board.Clone()                                      // safety copy
	bestMove, score := ai.minimax(shadowBoard, ai.Depth, Minimizing, -1) // AI is always minimizing for itself
	if ai.Debug {
		fmt.Printf(">>> Best move = %d with score = %d and depth = %d\n", bestMove, score, ai.Depth)
	}
	return bestMove
}

func (ai *MinimaxAI) minimax(board *game.Board, depth int, turn TurnType, lastMove int) (bestMove int, score int) {

	// 1. Check for a Win IMMEDIATELY
	if lastMove != -1 && board.WasWinningMove(lastMove) {
		// If it was a win, who made it?
		// If it's currently the Maximizing turn, the Minimizer (AI) just moved and won.
		if turn == Maximizing {
			return -1, -WIN_SCORE - depth // Minimizer won (prefer faster wins)
		}
		return -1, WIN_SCORE + depth // Maximizer won (prefer faster wins)
	}

	if depth == 0 {
		return -1, ai.evaluate(board, lastMove, turn)
	}

	validMoves := board.ValidMoves()
	bestMove = -1
	if len(validMoves) == 0 {
		return bestMove, 0
	}
	var bestScore int
	var cell game.Cell

	if turn == Minimizing {
		bestScore = math.MaxInt
		cell = ai.Self
	} else {
		bestScore = math.MinInt
		cell = ai.Other
	}

	for _, move := range board.ValidMoves() {
		err := board.MakeMove(move, cell)
		if err != nil {
			fmt.Printf("AI failed to make move: %v", err)
			break
		}
		_, score := ai.minimax(board, depth-1, switchTurn(turn), move)
		if ai.Debug {
			fmt.Printf("Depth %d: Move %d in a %v turn resulted in score %d (so far best move = %d, best score = %d)\n", depth, move, turn, score, bestMove, bestScore)
		}
		if turn == Minimizing {
			if score < bestScore {
				bestScore = score
				bestMove = move
			}
		} else {
			if score > bestScore {
				bestScore = score
				bestMove = move
			}
		}
		board.UndoMove(move)
	}
	return bestMove, bestScore
}

func switchTurn(turn TurnType) TurnType {
	if turn == Minimizing {
		return Maximizing
	}
	return Minimizing
}

func (ai *MinimaxAI) evaluate(board *game.Board, lastMove int, turnType TurnType) int {
	score := 0
	midCol := board.Width >> 1

	for col := 0; col < board.Width; col++ {
		// Weight the center column more (e.g., center = 3, edges = 0)
		dist := midCol - Abs(midCol-col)
		for row := 0; row < board.Height; row++ {
			if board.Grid[row][col] == ai.Self {
				score -= dist // AI wants a negative score
			} else if board.Grid[row][col] == ai.Other {
				score += dist // Human wants a positive score
			}
		}
	}
	return score
}

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
