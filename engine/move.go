package engine

import (
	"errors"
)

// piece name + beginning and ending squares
type Move struct {
	Piece      string // Piece.Name
	Begin, End Square
	Score      float64
}

func maxInt(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func minInt(x, y int) int {
	if x > y {
		return y
	}
	return x
}

// Modifies a board in-place.
// Returns an error without modifying board if illegal move.
// Removes a captured piece entirely from board.
// Changes the turn of the board once move is successfully completed.
func (b *Board) Move(m *Move) error {
	//for readability, this should be towards the end of the file
	if m.Piece == "k" && m.Begin.X-m.End.X != 1 && m.End.X-m.Begin.X != 1 {
		err := b.castleHandler(m)
		if err == nil {
			b.Turn *= -1
			b.Lastmove = *m
		}
		return err
	}

	var piecefound bool
	var pieceindex int
	var capture bool
	var capturedpiece int
	for i, p := range b.Board {
		if m.Begin == p.Position && m.Piece == p.Name && b.Turn == p.Color {
			pieceindex = i
			piecefound = true
		} else if m.End == p.Position && p.Color == b.Turn*-1 {
			capture = true
			capturedpiece = i
		}
		if piecefound && capture {
			break
		}
	}
	if !piecefound {
		return errors.New("func Move: invalid piece")
	}
	var legal bool
	legals := b.Board[pieceindex].legalMoves(b, true)
	for _, move := range legals {
		if *m == move {
			legal = true
			b.Board[pieceindex].Position = move.End
			break
		}
	}
	if !legal {
		return errors.New("func Move: illegal move")
	}

	// en passant
	if !capture && m.Piece == "p" && (m.Begin.X-m.End.X == 1 || m.End.X-m.Begin.X == 1) {
		capture = true
		for i, p := range b.Board {
			if p.Position.X == m.End.X && p.Position.Y == m.Begin.Y {
				capturedpiece = i
				break
			}
		}
	}

	if capture {
		b.Board[capturedpiece].Position = Square{
			X: 0,
			Y: 0,
		}
	}
	b.Board[pieceindex].Can_double_move = false
	if m.Piece == "k" || m.Piece == "r" {
		b.Board[pieceindex].Can_castle = false
	}
	for i, _ := range b.Board {
		b.Board[i].Can_en_passant = false
	}
	if m.Piece == "p" && m.Begin.Y-m.End.Y == 2*-b.Board[pieceindex].Color {
		b.Board[pieceindex].Can_en_passant = true
	}
	b.Turn *= -1
	b.Lastmove = *m
	return nil
}

func (b *Board) castleHandler(m *Move) error {
	if b.isCheck(b.Turn) {
		return errors.New("func castleHandler: king is in check")
	}

	var kingindex int
	var rookindex int
	var rookfound bool
	for i, p := range b.Board {
		if m.Begin == p.Position && m.Piece == p.Name && b.Turn == p.Color {
			kingindex = i
		} else if p.Name == "r" && ((m.End.X == 7 && p.Position.X == 8) || (m.End.X == 3 && p.Position.X == 1)) {
			if b.Turn == 1 && p.Position.Y == 1 {
				rookfound = true
				rookindex = i
			} else if b.Turn == -1 && p.Position.Y == 8 {
				rookfound = true
				rookindex = i
			}
		}
		if rookfound && kingindex > 0 {
			break
		}
	}
	if !b.Board[kingindex].Can_castle {
		return errors.New("func castleHandler: king has already moved")
	}
	if !rookfound {
		return errors.New("func castleHandler: no rook in position to castle to given side")
	}
	if !b.Board[rookindex].Can_castle {
		return errors.New("func castleHandler: rook has already moved")
	}
	for i := minInt(b.Board[rookindex].Position.X, b.Board[kingindex].Position.X) + 1; i < maxInt(b.Board[rookindex].Position.X, b.Board[kingindex].Position.X); i++ {
		s := &Square{
			X: i,
			Y: b.Board[kingindex].Position.Y,
		}
		if b.occupied(s) != 0 {
			return errors.New("func castleHandler: castle path is blocked")
		}
	}
	b.Board[kingindex].Position = m.End
	if b.isCheck(b.Turn) {
		b.Board[kingindex].Position = m.Begin
		return errors.New("func castleHandler: castle places user in check")
	}
	if m.End.X == 7 {
		b.Board[rookindex].Position.X = 6
	}
	if m.End.X == 3 {
		b.Board[rookindex].Position.X = 4
	}
	return nil
}
