package game

import (
	"strings"
	"testing"
)

type WantCell struct {
	r, c int
	val  Cell
}

func TestParseFromString(t *testing.T) {
	tests := []struct {
		name        string
		boardStr    string
		wantTopRows []int
		wantCells   []WantCell
	}{
		{
			name:        "Empty Board",
			boardStr:    "       " + "       " + "       " + "       " + "       " + "       ",
			wantTopRows: []int{0, 0, 0, 0, 0, 0, 0},
			wantCells: []WantCell{
				{r: 0, c: 0, val: EMPTY},
			},
		},
		{
			name: "Full Board",
			boardStr: "" +
				"XXXXXXX" +
				"OOOOOOO" +
				"XXXXXXX" +
				"OOOOOOO" +
				"XXXXXXX" +
				"OOOOOOO",
			wantTopRows: []int{6, 6, 6, 6, 6, 6, 6},
			wantCells: []WantCell{
				{r: 0, c: 0, val: PLAYER2},
				{r: 1, c: 0, val: PLAYER1},
				{r: 5, c: 6, val: PLAYER1},
			},
		},
		{
			name:        "Single piece at bottom left",
			boardStr:    "       " + "       " + "       " + "       " + "       " + "X      ",
			wantTopRows: []int{1, 0, 0, 0, 0, 0, 0},
			wantCells: []WantCell{
				{r: 0, c: 0, val: PLAYER1},
			},
		},
		{
			name: "Full column 0",
			// 'X' at indices 0, 7, 14, 21, 28, 35 (Column 0)
			boardStr: "X      " +
				"X      " +
				"X      " +
				"X      " +
				"X      " +
				"X      ",
			wantTopRows: []int{6, 0, 0, 0, 0, 0, 0},
			wantCells: []WantCell{
				{r: 5, c: 0, val: PLAYER1},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := NewBoard(7, 6)
			if ok := b.ParseFromString(tt.boardStr); !ok {
				t.Fatalf("ParseFromString failed")
			}

			// Check the gravity cache
			for i, got := range b.topEmptyRows {
				if got != tt.wantTopRows[i] {
					t.Errorf("Col %d: TopEmptyRows = %d, want %d", i, got, tt.wantTopRows[i])
				}
			}

			// Check a list of specific critical cells
			for _, wantCell := range tt.wantCells {
				if b.grid[wantCell.r][wantCell.c] != wantCell.val {
					t.Errorf("Cell [%d][%d] = %c, want %c",
						wantCell.r, wantCell.c, b.grid[wantCell.r][wantCell.c], wantCell.val)
				}
			}
		})
	}
}

func TestIsFull(t *testing.T) {
	tests := []struct {
		name     string
		boardStr string
		want     bool
	}{
		{
			name:     "Empty board is not full",
			boardStr: strings.Repeat(" ", 42),
			want:     false,
		},
		{
			name:     "Board with one slot left is not full",
			boardStr: " OXXXXO" + "XXXXOOO" + "OOOXXXX" + "XXXXOOO" + "OOOXXXX" + "XXXXOOO",
			want:     false,
		},
		{
			name:     "Completely full board",
			boardStr: "OXXXXOO" + "XXXXOOO" + "OOOXXXX" + "XXXXOOO" + "OOOXXXX" + "XXXXOOO",
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := NewBoard(7, 6)
			b.ParseFromString(tt.boardStr)
			if got := b.IsFull(); got != tt.want {
				t.Errorf("IsFull() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsWinningMove(t *testing.T) {
	tests := []struct {
		name     string
		boardStr string
		lastMove int
		want     bool
	}{
		{
			name:     "Horizontal Win",
			boardStr: "       " + "       " + "       " + "       " + "       " + "XXXX   ",
			lastMove: 3,
			want:     true,
		},
		{
			name:     "Vertical Win",
			boardStr: "O      " + "O      " + "O      " + "O      " + "       " + "       ",
			lastMove: 0, // Topmost O
			want:     true,
		},
		{
			name:     "Diagonal Win (Forward Slash)",
			boardStr: "   X   " + "  X    " + " X     " + "X      " + "       " + "       ",
			lastMove: 3,
			want:     true,
		},
		{
			name:     "No Win (3 in a row)",
			boardStr: "       " + "       " + "       " + "       " + "       " + "XXX    ",
			lastMove: 2,
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := NewBoard(7, 6)
			b.ParseFromString(tt.boardStr)
			if got := b.WasWinningMove(tt.lastMove); got != tt.want {
				t.Errorf("%s: IsWinningMove() = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}
