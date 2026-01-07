package main

import (
	"flag"
	"os"

	"uac.tech/connect4/internal/ai"
	"uac.tech/connect4/internal/game"
	"uac.tech/connect4/internal/ui/terminal"
)

func main() {

	// Define flags (Pointer approach)
	// flag.Int(name, default value, description)
	aiDepth := flag.Int("depth", 4, "Depth for the minimax AI (1-8)")
	playerName := flag.String("name", "Alice", "Human player name")
	aiName := flag.String("ainame", "Bob", "AI player name")
	isAiFirst := flag.Bool("aifirst", false, "Whether the AI plays first")
	isDebug := flag.Bool("debug", false, "Enable debug logging")

	// You MUST call this to actually parse the os.Args
	flag.Parse()

	player1 := game.Player{
		Name:  *playerName,
		Input: &terminal.HumanCLI{Input: os.Stdin, Output: os.Stdout},
		Cell:  game.PLAYER1,
	}
	// player2 := game.Player{
	// 	Name:  "Bob",
	// 	Input: &terminal.HumanCLI{Input: os.Stdin, Output: os.Stdout},
	// 	Cell:  game.PLAYER2,
	// }
	var self, other game.Cell
	if *isAiFirst {
		self = game.PLAYER1
		other = game.PLAYER2
	} else {
		self = game.PLAYER2
		other = game.PLAYER1
	}
	player2 := game.Player{
		Name:  *aiName,
		Input: &ai.MinimaxAI{Depth: *aiDepth, Self: self, Other: other, Debug: *isDebug},
		Cell:  game.PLAYER2,
	}
	renderer := &terminal.TerminalRenderer{Output: os.Stdout}

	game := game.NewGame(7, 6, &player2, &player1, renderer)
	game.Play()

}
