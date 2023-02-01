package mancala

import (
	"fmt"
	"mancala/lobby"
)

/*
* Mancala board looks like this:
* 2 players with 6 pits each and 1 personal pits for each player.
*
* 	Player1
* ____________________
* |	5 4 3 2 1 0	   |
* ____________________
* 0(P1 pit)	|	0(P2 pit)
* ____________________
* |	0 1 2 3 4 5  |
* ____________________
*
 */

type GameStatus string

const (
	// Game started or in progress
	Started GameStatus = "started"
	// Game is finished with an outcome
	PlayerTopWon    GameStatus = "p_top_won"
	PlayerBottomWon GameStatus = "p_bottom_won"
	Tie             GameStatus = "tie"
)

// TurnLog is a log of all turns
type TurnLog struct {
	Turn     int64
	Player   PlayerNumber
	PitIndex int
}

type Mancala struct {
	Board     *Board
	LobbySlug string

	NextPlayer PlayerNumber // 0 or 1
	Turn       int64        // count total amount of turns
	TurnsLog   []TurnLog    // log of all turns
	Status     GameStatus
}

// Start creates new Mancala game
// We don't define which player is first, it's always bottom one
func Start(lob *lobby.Lobby) *Mancala {
	return &Mancala{
		LobbySlug:  lob.Slug,
		Board:      NewBoard(),
		NextPlayer: PlayerBottom,
		Turn:       0,
		TurnsLog:   []TurnLog{},
		Status:     Started,
	}
}

func (m *Mancala) PlayTurn(player PlayerNumber, pitIndex int) error {
	// Check that move is valid(correct player, correct pit, game is in progress)
	if err := m.validateMove(player, pitIndex); err != nil {
		return err
	}
	// Log turns(analisys or show log to users?)
	m.TurnsLog = append(m.TurnsLog, TurnLog{Turn: m.Turn, Player: player, PitIndex: pitIndex})
	m.Turn += 1

	// Get stones from the pit
	stones := m.Board.pickStones(player, pitIndex)

	// Place stones in the next pits as a part of our turn
	m.placeStones(player, pitIndex+1, stones)

	// After the turn check if game finished
	if m.IsFinished() {
		m.FinalizeGame()
	}

	return nil
}

func (m *Mancala) validateMove(nextPlayer PlayerNumber, pitIndex int) error {
	if nextPlayer != m.NextPlayer {
		return fmt.Errorf("it's not your turn")
	}
	if m.Status != Started || m.IsFinished() {
		return fmt.Errorf("game is already finished")
	}

	if pitIndex < 0 || pitIndex > 5 {
		return fmt.Errorf("invalid pit index")
	}

	if m.Board.checkStones(nextPlayer, pitIndex) <= 0 {
		return fmt.Errorf("selected pit is empty")
	}

	return nil
}

func (m *Mancala) IsFinished() bool {
	return m.Board.arePitsEmpty(m.Board.Top) || m.Board.arePitsEmpty(m.Board.Bottom)
}

func (m *Mancala) FinalizeGame() {
	topScore, bottomScore := m.Board.getScores()

	switch {
	case topScore > bottomScore:
		m.Status = PlayerTopWon
	case topScore < bottomScore:
		m.Status = PlayerBottomWon
	default:
		m.Status = Tie
	}

}

func (m *Mancala) placeStones(player PlayerNumber, pitIndex, stones int) {
	// temp variable to store amount of stones in the current pit
	pitScore := 0
	// index of a currentPit we're evaluating
	currIndex := pitIndex

	// When player picked stones from his own pit
	for currIndex < amountOfPits {
		// On each pit we place one stone
		stones, pitScore = m.Board.placeStoneOnPlayerSide(player, currIndex, stones)

		// If we still have stones, we can continue placing them:
		if stones == 0 {
			break
		}
		currIndex++
	}

	switch {
	case currIndex != amountOfPits-1 && stones == 0 && pitScore == 1:
		// If we placed last stone in our personal pit, we can capture all stones from opposite pit
		stones = m.Board.pickStones(player, currIndex)

		oppositePitIndex := (currIndex - amountOfPits + 1) * -1
		stones += m.Board.pickStones(player, oppositePitIndex)

		m.Board.storeStones(player, stones)

		m.NextPlayer = nextPlayer(player)
		return
	case currIndex != amountOfPits-1 && stones == 0:
		// Otherwise just continue with next player's turn
		m.NextPlayer = nextPlayer(player)
		return
	}

	// When we finishsed our side
	// And we still have stones, store them in our personal pit
	if stones != 0 {
		stones = m.Board.storeStone(player, stones)
	}

	// If there are no more stones, we can skip opponent's turn
	if stones == 0 {
		return
	}

	// If we sitl have stones, we place them in enemy pits
	for index := 0; index < amountOfPits; index++ {
		stones, _ = m.Board.placeStoneOnEnemySide(player, index, stones)

		// If we're out of stones, turn ends
		if stones == 0 {
			m.NextPlayer = nextPlayer(player)
			return
		}
	}

	// If we sitll have stones after traversing personal and enemy pit
	// Start placing stones in our personal pit again
	m.placeStones(player, 0, stones)

}
