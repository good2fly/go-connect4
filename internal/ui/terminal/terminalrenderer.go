package terminal

import (
	"fmt"
	"io"
	"strings"

	"uac.tech/connect4/internal/game"
)

type TerminalRenderer struct {
	Output io.Writer
}

func (tr *TerminalRenderer) Game(g *game.Game, lastMoveRow int, lastMoveCol int) {
	var sb strings.Builder
	b := g.Board

	// Build the column headers (0 1 2...)
	sb.WriteString("\033[34m|\033[0m")
	for c := 0; c < b.Width; c++ {
		sb.WriteString(fmt.Sprintf("  %d \033[34m|\033[0m", c))
	}
	sb.WriteString("\n")

	// Print rows from Top to Bottom
	for r := b.Height - 1; r >= 0; r-- {
		sb.WriteString("\033[34m" + "+" + strings.Repeat("----+", b.Width) + "\033[0m\n")
		sb.WriteString("\033[34m|\033[0m")
		for c := 0; c < b.Width; c++ {
			cell := b.Grid[r][c]
			switch cell {
			case game.PLAYER1:
				if lastMoveRow == r && lastMoveCol == c {
					sb.WriteString("(\033[31m⬤\033[0m )\033[0m\033[34m|\033[0m") // Red
				} else {
					sb.WriteString(" \033[31m⬤\033[0m\033[34m  |\033[0m") // Red
				}
			case game.PLAYER2:
				if lastMoveRow == r && lastMoveCol == c {
					sb.WriteString("(\033[33m⬤\033[0m )\033[0m\033[34m|\033[0m") // Yellow
				} else {
					sb.WriteString(" \033[33m⬤\033[0m\033[34m  |\033[0m") // Yellow
				}
			default:
				sb.WriteString("  \033[34m·\033[0m \033[34m|\033[0m")
			}
		}
		sb.WriteString("\n")
	}

	// Bottom border
	sb.WriteString("\033[34m+" + strings.Repeat("----+", b.Width) + "\033[0m\n")
	fmt.Fprintln(tr.Output, sb.String())
}

func (tr *TerminalRenderer) Error(err error) {
	fmt.Fprintf(tr.Output, "Error: %v\n", err)
}

func (tr *TerminalRenderer) Message(msg string) {
	fmt.Fprintln(tr.Output, msg)

}

func (tr *TerminalRenderer) GameOver(g *game.Game) {
	tr.Game(g, -1, -1) // Draw the final board state one last time!
	switch g.GameState {
	case game.Player1Win:
		fmt.Fprintf(tr.Output, "🏆 %s Wins!\n", g.Player1.Name)
	case game.Player2Win:
		fmt.Fprintf(tr.Output, "🏆 %s Wins!\n", g.Player2.Name)
	case game.Draw:
		fmt.Fprintln(tr.Output, "🤝 It's a draw!")
	default:
		fmt.Fprintln(tr.Output, "Game ended.")
	}
}
