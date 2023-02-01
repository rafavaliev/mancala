package mancala

import (
	"mancala/lobby"
	"reflect"
	"testing"
)

func TestMancala_PlayTurn(t *testing.T) {

	l := &lobby.Lobby{Slug: "hello"}

	tests := []struct {
		name      string
		currState *Mancala
		wantState *Mancala
		player    PlayerNumber
		pitIndex  int
		wantErr   bool
	}{
		{
			name:      "validation: wrong player's turn",
			currState: Start(l),
			player:    PlayerTop,
			pitIndex:  0,
			wantErr:   true,
		},
		{
			name:      "validation: wrong pit index",
			currState: Start(l),
			player:    PlayerBottom,
			pitIndex:  55,
			wantErr:   true,
		},
		{
			name: "validation: empty pit",
			currState: &Mancala{
				Board: &Board{
					Top:    NewPits(),
					Bottom: [6]int{0, 6, 6, 6, 6, 6},
				},
				NextPlayer: PlayerBottom,
				Status:     Started,
			},
			player:   PlayerBottom,
			pitIndex: 0,
			wantErr:  true,
		},
		{
			name: "validation: mancala game is finished",
			currState: &Mancala{
				Board: &Board{
					Top:    NewPits(),
					Bottom: [6]int{0, 0, 0, 0, 0, 0},
				},
				NextPlayer: PlayerBottom,
				Status:     Started,
			},
			player:   PlayerBottom,
			pitIndex: 0,
			wantErr:  true,
		},
		{
			name:      "bottom player turn gets another turn",
			currState: Start(l),
			wantState: func() *Mancala {
				m := Start(l)
				m.Board.Bottom = [6]int{0, 7, 7, 7, 7, 7}
				m.Board.BottomPlayerStore = 1
				m.NextPlayer = PlayerBottom
				return m
			}(),
			player:   PlayerBottom,
			pitIndex: 0,
		},
		{
			name: "bottom player turn doesn't affect top player's store ",
			currState: func() *Mancala {
				m := Start(l)
				m.Board.Bottom = [6]int{2, 2, 0, 0, 0, 8}
				return m
			}(),
			wantState: func() *Mancala {
				m := Start(l)
				m.Board.Bottom = [6]int{3, 2, 0, 0, 0, 0}
				m.Board.BottomPlayerStore = 1
				m.Board.Top = [6]int{7, 7, 7, 7, 7, 7}
				m.Board.TopPlayerStore = 0
				m.NextPlayer = PlayerTop
				return m
			}(),
			player:   PlayerBottom,
			pitIndex: 5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.currState.PlayTurn(tt.player, tt.pitIndex); (err != nil) != tt.wantErr {
				t.Log(err)
				t.Errorf("PlayTurn() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr {
				return
			}

			if !reflect.DeepEqual(tt.currState.Board.Bottom, tt.wantState.Board.Bottom) {
				t.Errorf("PlayTurn().Bottom got = %v, want %v", tt.currState.Board.Bottom, tt.wantState.Board.Bottom)
			}
			if !reflect.DeepEqual(tt.currState.Board.Top, tt.wantState.Board.Top) {
				t.Errorf("PlayTurn().Top got = %v, want %v", tt.currState.Board.Top, tt.wantState.Board.Top)
			}
			if tt.currState.Board.TopPlayerStore != tt.wantState.Board.TopPlayerStore {
				t.Errorf("PlayTurn().TopPlayerStore got = %v, want %v", tt.currState.Board, tt.wantState.Board)
			}
			if tt.currState.Board.BottomPlayerStore != tt.wantState.Board.BottomPlayerStore {
				t.Errorf("PlayTurn().BottomPlayerStore got = %v, want %v", tt.currState.Board, tt.wantState.Board)
			}
			if tt.currState.NextPlayer != tt.wantState.NextPlayer {
				t.Errorf("PlayTurn().NextPlayer got = %v, want %v", tt.currState.NextPlayer, tt.wantState.NextPlayer)
			}
		})
	}
}
