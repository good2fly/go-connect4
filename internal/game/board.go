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
	width        int
	height       int
	grid         [][]Cell
	topEmptyRows []int
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
		width:        width,
		height:       height,
		grid:         grid,
		topEmptyRows: top,
	}
}

func (b *Board) Clone() *Board {
	newBoard := &Board{
		width:        b.width,
		height:       b.height,
		grid:         make([][]Cell, b.height),
		topEmptyRows: slices.Clone(b.topEmptyRows),
	}
	for i := range b.grid {
		newBoard.grid[i] = slices.Clone(b.grid[i])
	}
	return newBoard
}

func (b *Board) Reset() {
	for i := range b.grid {
		for j := range b.grid[i] {
			b.grid[i][j] = EMPTY
		}
	}
	for i := range b.topEmptyRows {
		b.topEmptyRows[i] = 0
	}
}

func (b *Board) Width() int {
	return b.width
}

func (b *Board) Height() int {
	return b.height
}

func (b *Board) CellAt(row, col int) Cell {
	return b.grid[row][col]
}

func (b *Board) TopNonEmptyRow(col int) int {
	return b.topEmptyRows[col] - 1
}

// ValidMoves returns a slice of column indices where a piece can be legally placed.
func (b *Board) ValidMoves() []int {
	moves := []int{}
	for col := 0; col < b.width; col++ {
		if b.topEmptyRows[col] < b.height {
			moves = append(moves, col)
		}
	}
	return moves
}

func (b *Board) IsValidMove(col int) bool {
	return col >= 0 && col < b.width && b.topEmptyRows[col] < b.height
}

// MakeMove places a piece for the given player in the specified column.
// It returns the Move made or an error if the move is invalid.
func (b *Board) MakeMove(col int, player Cell) error {
	if col < 0 || col >= b.width {
		return fmt.Errorf("invalid column %d", col)
	}
	row := b.topEmptyRows[col]
	if !b.IsValidMove(col) {
		return fmt.Errorf("column is full %d", col)
	}
	b.grid[row][col] = player
	b.topEmptyRows[col]++
	return nil
}

func (b *Board) UndoMove(col int) error {
	if col < 0 || col >= b.width {
		return fmt.Errorf("invalid column %d", col)
	}
	b.topEmptyRows[col]--
	b.grid[b.topEmptyRows[col]][col] = EMPTY
	return nil
}

// IsWinning checks if the last move made resulted in a win for the player.
// Note, that the move was already made on the board.
// A win is defined as at least 4 consecutive pieces of the same kind in any direction (horizontal, vertical, diagonal).
func (b *Board) WasWinningMove(move int) bool {
	moveCol := move
	moveRow := b.topEmptyRows[moveCol] - 1
	player := b.grid[moveRow][moveCol]

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
				if r < 0 || r >= b.height || c < 0 || c >= b.width || b.grid[r][c] != player {
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
	for col := 0; col < b.width; col++ {
		if b.topEmptyRows[col] < b.height {
			return false
		}
	}
	return true
}

// ParseFromString populates the board state from a string representation.
// The string should be of length width*height, with characters 'X', 'O', or ' ', representing
// each cell in row-major order (bottom row first).
func (b *Board) ParseFromString(boardStr string) bool {
	if len(boardStr) != b.width*b.height {
		return false
	}
	for i := range b.topEmptyRows {
		b.topEmptyRows[i] = 0
	}
	for strRow := 0; strRow < b.height; strRow++ {
		for col := 0; col < b.width; col++ {
			boardRow := b.height - 1 - strRow
			piece := Cell(boardStr[strRow*b.width+col])
			b.grid[boardRow][col] = piece
			if piece != EMPTY && b.topEmptyRows[col] < boardRow+1 {
				b.topEmptyRows[col] = boardRow + 1
			}
		}
	}
	return true
}
