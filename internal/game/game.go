package game

import "fmt"

type PlayerInput interface {
	GetMove(board *Board) int
}

type Player struct {
	Name  string
	Input PlayerInput
	Cell  Cell
}

type GameStateType int

const (
	Playing GameStateType = iota
	Draw
	Player1Win
	Player2Win
)

type Game struct {
	GameState GameStateType
	Board     *Board
	Player1   Player
	Player2   Player
	renderer  GameRenderer
}

type GameRenderer interface {
	Game(g *Game, lastMoveRow int, lastMoveCol int)
	Error(err error)
	Message(msg string)
	GameOver(g *Game)
}

func NewGame(width, height int, player1 *Player, player2 *Player, renderer GameRenderer) *Game {
	return &Game{
		GameState: Playing,
		Board:     NewBoard(width, height),
		Player1:   *player1,
		Player2:   *player2,
		renderer:  renderer,
	}
}

func (g *Game) Play() {
	// for rendering last move
	lastMovCol := -1
	lastMoveRow := -1
	currentPlayer := &g.Player1
	for {
		g.renderer.Message(fmt.Sprintf("\n%s's turn\n", currentPlayer.Name))
		g.renderer.Game(g, lastMoveRow, lastMovCol)
		moveCol := currentPlayer.Input.GetMove(g.Board)
		err := g.Board.MakeMove(moveCol, currentPlayer.Cell)
		if err != nil {
			g.renderer.Error(err)
			continue
		}
		lastMovCol = moveCol
		lastMoveRow = g.Board.TopNonEmptyRow(moveCol)
		if g.Board.WasWinningMove(moveCol) {
			// Handle win condition
			if currentPlayer == &g.Player1 {
				g.GameState = Player1Win
			} else {
				g.GameState = Player2Win
			}
			g.renderer.GameOver(g)
			break
		}
		if g.Board.IsFull() {
			g.GameState = Draw
			g.renderer.GameOver(g)
			break
		}
		// Switch players
		if currentPlayer == &g.Player1 {
			currentPlayer = &g.Player2
		} else {
			currentPlayer = &g.Player1
		}
	}
}
