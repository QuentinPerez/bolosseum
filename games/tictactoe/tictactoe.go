package tictactoe

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/moul/bolosseum/bots"
	"github.com/moul/bolosseum/games"
)

var pieces = []string{"X", "O"}

type TictactoeGame struct {
	games.BotsBasedGame

	board map[string]string `json:"board",omitempty`
}

func NewGame() (*TictactoeGame, error) {
	game := TictactoeGame{}
	game.Bots = make([]bots.Bot, 0)
	game.board = make(map[string]string, 9)

	for y := 0; y < 3; y++ {
		for x := 0; x < 3; x++ {
			game.board[fmt.Sprintf("%d-%d", y, x)] = ""
		}
	}

	return &game, nil
}

func (g *TictactoeGame) CheckArgs(args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("You need to specify 2 bots")
	}
	return nil
}

func (g *TictactoeGame) checkBoard() (bots.Bot, error) {
	// check if the board is invalid
	if len(g.board) != 9 {
		return nil, fmt.Errorf("Invalid board: %d cases", len(g.board))
	}

	// check if there is a winner
	for idx, piece := range pieces {

		// check for horizontal match
		for y := 0; y < 3; y++ {
			valid := true
			for x := 0; x < 3; x++ {
				if g.board[fmt.Sprintf("%d-%d", y, x)] != piece {
					valid = false
					break
				}
			}
			if valid {
				return g.Bots[idx], nil
			}
		}

		// check for vertical match
		for x := 0; x < 3; x++ {
			valid := true
			for y := 0; y < 3; y++ {
				if g.board[fmt.Sprintf("%d-%d", y, x)] != piece {
					valid = false
					break
				}
			}
			if valid {
				return g.Bots[idx], nil
			}
		}

		// diagonals
		if g.board["0-0"] == piece && g.board["1-1"] == piece && g.board["2-2"] == piece {
			return g.Bots[idx], nil
		}
		if g.board["0-2"] == piece && g.board["1-1"] == piece && g.board["2-0"] == piece {
			return g.Bots[idx], nil
		}
	}

	return nil, nil
}

func (g *TictactoeGame) Run(gameID string) error {
	if err := bots.InitTurnBasedBots(g.Bots, g.Name(), gameID); err != nil {
		return err
	}

	// play
	for turn := 0; turn < 9; turn++ {
		idx := turn % 2
		bot := g.Bots[idx]
		piece := pieces[idx]

		reply, err := bot.SendMessage(bots.QuestionMessage{
			GameID:      gameID,
			Game:        g.Name(),
			Action:      "play-turn",
			Board:       g.board,
			You:         piece,
			PlayerIndex: idx,
		})
		if err != nil {
			return err
		}

		g.board[reply.Play.(string)] = piece

		// check board
		winner, err := g.checkBoard()
		if err != nil {
			return err
		}
		if winner != nil {
			logrus.Warnf("Player %d (%s) won", idx, winner.Name())
			return nil
		}
	}

	logrus.Warnf("Draw")
	return nil
}

func (g *TictactoeGame) Name() string {
	return "tictactoe"
}
