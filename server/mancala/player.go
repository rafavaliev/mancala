package mancala

type PlayerStore int

type PlayerNumber int

const (
	PlayerBottom PlayerNumber = 0
	PlayerTop    PlayerNumber = 1
)

func nextPlayer(player PlayerNumber) PlayerNumber {
	if player == PlayerBottom {
		return PlayerTop
	}
	return PlayerBottom
}
