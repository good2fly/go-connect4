package game

import (
	"fmt"
	"slices"
)

type Cell uint8

const (
	EMPTY   Cell = ' '
	PLAYER1 Cell = 'X'
	PLAYER2 Cell = 'O'
)

// Board represents a 6x7 Connect 4 grid.
// The zero-value of Board is an empty grid ready to play.
type Board struct {
	Width        int
	Height       int
	Grid         [][]Cell
	TopEmptyRows []int
}

func NewBoard(width, height int) *Board {
	if width <= 0 {
		width = 7
	}
	if height <= 0 {
		height = 6
	}

	grid := make([][]Cell, height)
	for i := range grid {
		grid[i] = make([]Cell, width)
	}

	top := make([]int, width)

	return &Board{
		Width:        width,
		Height:       height,
		Grid:         grid,
		TopEmptyRows: top,
	}
}

func (b *Board) Clone() *Board {
	newBoard := &Board{
		Width:        b.Width,
		Height:       b.Height,
		Grid:         make([][]Cell, b.Height),
		TopEmptyRows: slices.Clone(b.TopEmptyRows),
	}
	for i := range b.Grid {
		newBoard.Grid[i] = slices.Clone(b.Grid[i])
	}
	return newBoard
}

func (b *Board) Reset() {
	for i := range b.Grid {
		for j := range b.Grid[i] {
			b.Grid[i][j] = EMPTY
		}
	}
	for i := range b.TopEmptyRows {
		b.TopEmptyRows[i] = 0
	}
}

func (b *Board) ValidMoves() []int {
	moves := []int{}
	for col := 0; col < b.Width; col++ {
		if b.TopEmptyRows[col] < b.Height {
			moves = append(moves, col)
		}
	}
	return moves
}

func (b *Board) IsValidMove(col int) bool {
	return col >= 0 && col < b.Width && b.TopEmptyRows[col] < b.Height
}

// MakeMove places a piece for the given player in the specified column.
// It returns the Move made or an error if the move is invalid.
func (b *Board) MakeMove(col int, player Cell) error {
	if col < 0 || col >= b.Width {
		return fmt.Errorf("invalid column %d", col)
	}
	row := b.TopEmptyRows[col]
	if !b.IsValidMove(col) {
		return fmt.Errorf("column is full %d", col)
	}
	b.Grid[row][col] = player
	b.TopEmptyRows[col]++
	return nil
}

func (b *Board) UndoMove(col int) error {
	if col < 0 || col >= b.Width {
		return fmt.Errorf("invalid column %d", col)
	}
	b.TopEmptyRows[col]--
	b.Grid[b.TopEmptyRows[col]][col] = EMPTY
	return nil
}

// IsWinning checks if the last move made resulted in a win for the player.
// Note, that the move was already made on the board.
// A win is defined as at least 4 consecutive pieces of the same kind in any direction (horizontal, vertical, diagonal).
func (b *Board) WasWinningMove(move int) bool {
	moveCol := move
	moveRow := b.TopEmptyRows[moveCol] - 1
	player := b.Grid[moveRow][moveCol]

	directions := [][2]int{
		{1, 0},  // Vertical
		{0, 1},  // Horizontal
		{1, 1},  // Diagonal /
		{1, -1}, // Diagonal \
	}

	// Check all directions from the last move
	for _, dir := range directions {
		count := 1
		for d := -1; d <= 1; d += 2 { // both directions of a particular kind (vertical/horizontal/diagonal)
			r, c := moveRow, moveCol
			for {
				r += d * dir[0]
				c += d * dir[1]
				if r < 0 || r >= b.Height || c < 0 || c >= b.Width || b.Grid[r][c] != player {
					break
				}
				count++
				if count >= 4 {
					return true
				}
			}
		}
	}
	return false
}

// IsFull checks if the board is completely filled (i.e., no more valid moves).
func (b *Board) IsFull() bool {
	for col := 0; col < b.Width; col++ {
		if b.TopEmptyRows[col] < b.Height {
			return false
		}
	}
	return true
}

// ParseFromString populates the board state from a string representation.
// The string should be of length Width*Height, with characters 'X', 'O', or ' ', representing
// each cell in row-major order (bottom row first).
func (b *Board) ParseFromString(boardStr string) bool {
	if len(boardStr) != b.Width*b.Height {
		return false
	}
	for i := range b.TopEmptyRows {
		b.TopEmptyRows[i] = 0
	}
	for strRow := 0; strRow < b.Height; strRow++ {
		for col := 0; col < b.Width; col++ {
			boardRow := b.Height - 1 - strRow
			piece := Cell(boardStr[strRow*b.Width+col])
			b.Grid[boardRow][col] = piece
			if piece != EMPTY && b.TopEmptyRows[col] < boardRow+1 {
				b.TopEmptyRows[col] = boardRow + 1
			}
		}
	}
	return true
}
