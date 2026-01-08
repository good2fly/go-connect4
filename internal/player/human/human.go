package terminal

import (
	"bufio"
	"fmt"
	"io"
	"strconv"

	"uac.tech/connect4/internal/game"
)

type HumanCLI struct {
	Input  io.Reader
	Output io.Writer
}

// GetMove satisfies the PlayerInput interface
func (h *HumanCLI) GetMove(b *game.Board) int {
	scanner := bufio.NewScanner(h.Input)
	for {
		fmt.Fprint(h.Output, "Select column: ")
		if !scanner.Scan() {
			return -1
		}

		val, err := strconv.Atoi(scanner.Text())
		if err != nil || !b.IsValidMove(val) {
			fmt.Fprintln(h.Output, "Invalid move!")
			continue
		}
		return val
	}
}
