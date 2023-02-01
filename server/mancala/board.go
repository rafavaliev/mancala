package mancala

const amountOfPits = 6

type Pits [amountOfPits]int

func NewPits() Pits {
	return Pits{6, 6, 6, 6, 6, 6}
}

type Board struct {
	Top            Pits
	TopPlayerStore PlayerStore

	Bottom            Pits
	BottomPlayerStore PlayerStore
}

func NewBoard() *Board {
	return &Board{
		Top:               NewPits(),
		TopPlayerStore:    0,
		Bottom:            NewPits(),
		BottomPlayerStore: 0,
	}
}

func (b *Board) arePitsEmpty(pits Pits) bool {
	for _, pit := range pits {
		if pit != 0 {
			return false
		}
	}
	return true
}

func (b *Board) checkStones(player PlayerNumber, index int) int {
	switch player {
	case PlayerBottom:
		return b.Bottom[index]
	case PlayerTop:
		return b.Top[index]
	}
	return 0
}

func (b *Board) pickStones(player PlayerNumber, index int) int {
	var stones int
	switch player {
	case PlayerBottom:
		stones = b.Bottom[index]
		b.Bottom[index] = 0
	case PlayerTop:
		stones = b.Top[index]
		b.Top[index] = 0
	}
	return stones
}

func (b *Board) pickStonesEnemySide(player PlayerNumber, index int) int {
	pits := b.Bottom
	if player == PlayerBottom {
		pits = b.Top
	}
	stones := pits[index]
	pits[index] = 0
	return stones
}

func (b *Board) getScores() (int, int) {
	topScore := int(b.TopPlayerStore)
	for _, pitValues := range b.Top {
		topScore += pitValues
	}

	bottomScore := int(b.BottomPlayerStore)
	for _, pitValues := range b.Bottom {
		bottomScore += pitValues
	}
	return topScore, bottomScore
}

// When we traversed all pits, we can store a stone in a player's store
func (b *Board) storeStone(player PlayerNumber, stones int) int {
	switch player {
	case PlayerBottom:
		b.BottomPlayerStore++
	case PlayerTop:
		b.TopPlayerStore++
	}
	return stones - 1
}

func (b *Board) storeStones(player PlayerNumber, stones int) int {
	switch player {
	case PlayerBottom:
		b.BottomPlayerStore += PlayerStore(stones)
	case PlayerTop:
		b.TopPlayerStore += PlayerStore(stones)
	}
	return 0
}

// Places a stone on the player side of the board, returns the amount of stones left and the amount of stones in the pit
func (b *Board) placeStoneOnPlayerSide(player PlayerNumber, pitIndex int, stones int) (int, int) {
	pitScore := 0
	switch {
	case player == PlayerBottom:
		b.Bottom[pitIndex] += 1
		pitScore = b.Bottom[pitIndex]
	case player == PlayerTop:
		b.Top[pitIndex] += 1
		pitScore = b.Top[pitIndex]
	}
	return stones - 1, pitScore
}

// Places a stone on the enemy side of the board, returns the amount of stones left and the amount of stones in the pit
func (b *Board) placeStoneOnEnemySide(player PlayerNumber, pitIndex int, stones int) (int, int) {
	pitScore := 0
	switch {
	case player == PlayerBottom:
		b.Top[pitIndex]++
		pitScore = b.Top[pitIndex]
	case player == PlayerTop:
		b.Bottom[pitIndex]++
		pitScore = b.Bottom[pitIndex]
	}
	return stones - 1, pitScore
}
