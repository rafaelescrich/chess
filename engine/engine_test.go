package engine

import "testing"

func TestUndoMove(t *testing.T) {
	board := &Board{Turn: -1}
	board.PlacePiece('q', 1, 1, 8)
	board.PlacePiece('n', -1, 1, 8)
	board.PlacePiece('q', -1, 1, 8)
	board.PlacePiece('q', -1, 1, 8)
	board.Board[1].Captured = true
	board.Board[2].Captured = true
	board.Board[3].Captured = true
	m := &Move{
		Piece: 'p',
		Begin: Square{
			X: 2,
			Y: 7,
		},
		End: Square{
			X: 1,
			Y: 8,
		},
		Promotion: 'q',
		Capture:   'q',
	}
	board.UndoMove(m)
	if pos := board.Board[0].Position; pos != m.Begin {
		t.Errorf("Undone piece should have had position %s instead was on %s", m.Begin.ToString(), pos.ToString())
	}
	if captured := board.Board[1].Captured; !captured {
		t.Error("Did not uncapture correct piece")
	}
	if captured := board.Board[2].Captured; captured {
		t.Error("Undone captured piece still captured")
	}
	if captured := board.Board[3].Captured; !captured {
		t.Error("Uncaptured more than one piece")
	}
	if name := board.Board[0].Name; name != 'p' {
		t.Errorf("Undoing a promotion gave piece name %s instead of pawn", string(name))
	}
	board = &Board{Turn: -1}
	board.PlacePiece('k', 1, 7, 1)
	board.PlacePiece('r', 1, 6, 1)
	m = &Move{
		Piece: 'k',
		Begin: Square{
			X: 5,
			Y: 1,
		},
		End: Square{
			X: 7,
			Y: 1,
		},
	}
	board.UndoMove(m)
	if board.Board[1].Position.X != 8 || board.Board[1].Position.Y != 1 {
		t.Errorf("Undoing castle should have left rook at (8, 1), instead at %+v", board.Board[1].Position)
	}
	board = &Board{Turn: 1}
	board.PlacePiece('q', 1, 1, 1)
	board.PlacePiece('r', -1, 8, 8)
	board.PlacePiece('r', -1, 8, 7)
	m1 := board.Board[0].makeMoveTo(8, 8)
	m1.Capture = 'r'
	m2 := board.Board[2].makeMoveTo(8, 8)
	m2.Capture = 'q'
	board.ForceMove(m1)
	board.ForceMove(m2)
	board.UndoMove(m2)
	board.UndoMove(m1)
	for i, p := range board.Board {
		if p.Captured {
			t.Errorf("Double capture erased piece at index %d", i)
		}
	}
}

func TestForceMove(t *testing.T) {
	board := &Board{Turn: -1}
	board.PlacePiece('k', 1, 1, 1)
	m := board.Board[0].makeMoveTo(2, 2)
	board.ForceMove(m)
	if board.Board[0].Position.X != 2 || board.Board[0].Position.Y != 2 {
		t.Errorf("ForceMove didn't move the king, king should be at 2,2 instead at %+v", board.Board[0].Position)
	}
	board.PlacePiece('p', -1, 1, 2)
	m = board.Board[1].makeMoveTo(1, 1)
	m.Promotion = 'r'
	board.ForceMove(m)
	if board.Board[1].Name != 'r' {
		t.Errorf("Promotion didn't go through, promoted pawn's name is %s instead of rook", string(board.Board[1].Name))
	}
	board = &Board{Turn: 1}
	board.PlacePiece('k', 1, 5, 1)
	board.PlacePiece('r', 1, 8, 1)
	for i, _ := range board.Board {
		board.Board[i].Can_castle = true
	}
	board.ForceMove(board.Board[0].makeMoveTo(7, 1))
	if board.Board[1].Position.X != 6 || board.Board[1].Position.Y != 1 {
		t.Errorf("Forced castling left rook on %+v", board.Board[1].Position)
	}
}

func TestAttacking(t *testing.T) {
	board := &Board{Turn: 1}
	board.PlacePiece('k', 1, 1, 1)
	board.PlacePiece('r', 1, 2, 2)
	board.PlacePiece('p', 1, 2, 3)
	rook := board.Board[1]
	s := &Square{
		X: 4,
		Y: 2,
	}
	if !rook.Attacking(s, board) {
		t.Errorf("Rook not attacking on open line, should be attacking %+v from %+v", s, rook.Position)
	}
	s.X, s.Y = 2, 3
	if !rook.Attacking(s, board) {
		t.Errorf("Rook not attacking own piece, should be attacking %+v from %+v", s, rook.Position)
	}
	s.Y = 5
	if rook.Attacking(s, board) {
		t.Errorf("Rook attacking through own piece, should not be attacking %+v from %+v", s, rook.Position)
	}
}

func TestMakeMoveTo(t *testing.T) {
	board := &Board{Turn: 1}
	board.PlacePiece('k', 1, 1, 1)
	m := board.Board[0].makeMoveTo(2, 2)
	if m.Piece != 'k' {
		t.Error("Warped piece name from ", 'k', " to ", m.Piece)
	}
	if m.Begin != board.Board[0].Position {
		t.Errorf("Piece originated at 1,1. Current location: %+v, move begin: %+v", board.Board[0].Position, m.Begin)
	}
	if m.End.X != 2 || m.End.Y != 2 {
		t.Errorf("Incorrect ending square. Should be 2, 2, ended up at %+v", m.End)
	}
}

func TestAllLegalMoves(t *testing.T) {
	board := &Board{Turn: -1}
	board.PlacePiece('k', -1, 1, 1)
	board.PlacePiece('k', 1, 8, 8)
	board.PlacePiece('p', -1, 4, 3)
	moves := board.AllLegalMoves()
	if moveslen := len(moves); moveslen != 4 {
		t.Errorf("Too many possible moves on the board. 4 moves expected, %d moves recieved", moveslen)
	}
	for i, m1 := range moves {
		for j, m2 := range moves {
			if m2 == m1 && i != j {
				t.Error("Duplicate moves returned, ", moves)
			}
		}
	}
}

func TestCopyMove(t *testing.T) {
	move := &Move{
		Piece: 'k',
		Begin: Square{
			X: 1,
			Y: 1,
		},
		End: Square{
			X: 2,
			Y: 2,
		},
		Score: 2,
	}
	newmove := move.CopyMove()
	if !(newmove.Piece == move.Piece && newmove.Begin == move.Begin && newmove.End == move.End) {
		t.Errorf("Something went wrong copying the move, %+v was expected, %+v was returned", move, newmove)
	}
	newmove.Score = 3
	if move.Score != 2 {
		t.Error("Changing newmove changed master move")
	}
}

func TestIsOver(t *testing.T) {
	board := &Board{Turn: 1}
	board.PlacePiece('k', 1, 1, 1)
	board.PlacePiece('q', -1, 2, 2)
	board.PlacePiece('r', -1, 8, 2)
	if result := board.IsOver(); result != -2 {
		t.Errorf("Expected black wins, got a result of %d", result)
	}
	board.Board[1].Position.Y = 3
	if result := board.IsOver(); result != 1 {
		t.Errorf("Expected stalemate, got a result of %d", result)
	}
	board = &Board{Turn: -1}
	board.PlacePiece('k', 1, 1, 1)
	board.PlacePiece('k', -1, 8, 8)
	board.PlacePiece('b', 1, 6, 6)
	board.PlacePiece('r', -1, 8, 7)
	board.PlacePiece('r', -1, 7, 8)
	if over := board.IsOver(); over != 0 {
		t.Errorf("Black is in check but can block, IsOver still returned %d", over)
	}
}

func TestOccupied(t *testing.T) {
	b := &Board{}
	b.SetUpPieces()
	whitesquare := &Square{
		X: 1,
		Y: 1,
	}
	blacksquare := &Square{
		X: 8,
		Y: 8,
	}
	emptysquare := &Square{
		X: 5,
		Y: 5,
	}
	nonsquare := &Square{
		X: 10,
		Y: 10,
	}
	if out, _ := b.Occupied(whitesquare); out != 1 {
		t.Errorf("expected 1, got %d", out)
	}
	if out, _ := b.Occupied(blacksquare); out != -1 {
		t.Errorf("expected -1, got %d", out)
	}
	if out, _ := b.Occupied(emptysquare); out != 0 {
		t.Errorf("expected 0, got %d", out)
	}
	if out, _ := b.Occupied(nonsquare); out != -2 {
		t.Errorf("expected -2, got %d", out)
	}
}

func TestIsCheck(t *testing.T) {
	board := &Board{Turn: 1}
	board.PlacePiece('k', 1, 1, 1)
	board.PlacePiece('k', -1, 8, 8)
	board.PlacePiece('r', 1, 8, 1)
	if check := board.IsCheck(1); check == true {
		t.Error("False positive when determining check")
	}
	if check := board.IsCheck(-1); check == false {
		t.Error("False negative when determining check")
	}
	if king := board.Board[0]; king.Position.X != 1 || king.Position.Y != 1 {
		t.Errorf("isCheck modified board, king moves from {X: 1, Y:1} to %+v", king.Position)
	}
}

func TestMoveIsCheck(t *testing.T) {
	board := &Board{Turn: 1}
	board.PlacePiece('k', 1, 1, 1)
	board.PlacePiece('b', 1, 2, 2)
	board.PlacePiece('q', -1, 4, 4)
	checkmove := board.Board[1].makeMoveTo(3, 1)
	if check := moveIsCheck(board, checkmove); !check {
		t.Error("Check not recognized")
	}
	okmove := board.Board[1].makeMoveTo(3, 3)
	if check := moveIsCheck(board, okmove); check {
		t.Error("False positive with ok move")
	}
	capturemove := board.Board[1].makeMoveTo(4, 4)
	if check := moveIsCheck(board, capturemove); check {
		t.Error("Capturing pinning piece with pinned piece places user in check")
	}
	board = &Board{Turn: 1}
	board.PlacePiece('k', 1, 1, 1)
	board.PlacePiece('r', -1, 8, 1)
	board.PlacePiece('b', 1, 7, 2)
	m := board.Board[2].makeMoveTo(8, 1)
	if check := moveIsCheck(board, m); check {
		t.Error("Capturing the attacking piece still places user in check")
	}
}

func TestMove(t *testing.T) {
	board := &Board{Turn: 1}
	board.PlacePiece('r', 1, 1, 1)
	board.PlacePiece('n', -1, 2, 1)
	m := board.Board[0].makeMoveTo(2, 1)
	if err := board.Move(m); err != nil {
		t.Errorf("Got an unexpected error making a legal capture: %s", err)
	}
	out := []*Piece{
		&Piece{
			Name: 'r',
			Position: Square{
				Y: 1,
				X: 2,
			},
			Color: 1,
			Directions: [][2]int{
				{1, 0},
				{-1, 0},
				{0, 1},
				{0, -1},
			},
			Infinite_direction: true,
		},
		&Piece{
			Name: 'n',
			Position: Square{
				Y: 1,
				X: 2,
			},
			Color: -1,
			Directions: [][2]int{
				{1, 2},
				{-1, 2},
				{1, -2},
				{-1, -2},
				{2, 1},
				{-2, 1},
				{2, -1},
				{-2, -1},
			},
			Captured: true,
		},
	}
	if !(len(board.Board) == len(out) && board.Board[0].Position == out[0].Position && board.Board[1].Captured) {
		t.Errorf("Expected: %+v\nGot: %+v", out, board.Board)
	}
	board.Turn = 1
	m = &Move{
		Piece: 'r',
		Begin: Square{
			Y: 8,
			X: 8,
		},
		End: Square{
			Y: 7,
			X: 8,
		},
	}
	if err := board.Move(m); err == nil {
		t.Error("Accessing an invalid piece did not return an error")
	}
	m = board.Board[0].makeMoveTo(4, 4)
	if err := board.Move(m); err == nil {
		t.Error("Attempting an illegal move did not return an error")
	}
	board = &Board{Turn: 1}
	board.PlacePiece('p', -1, 2, 5)
	board.Board[0].Can_en_passant = true
	board.PlacePiece('p', 1, 3, 5)
	m = board.Board[1].makeMoveTo(2, 6)
	if err := board.Move(m); err != nil {
		t.Errorf("En passant unexpected error: %s", err)
	}
	if !board.Board[0].Captured {
		t.Error("After en passant, captured piece not taken off board. Captured is still false.")
	}
	board = &Board{Turn: 1}
	board.PlacePiece('p', 1, 1, 7)
	m = board.Board[0].makeMoveTo(1, 8)
	m.Promotion = 'q'
	if err := board.Move(m); err != nil {
		t.Errorf("Promoting pawn raised error %s", err)
	}
	if piece := board.Board[0]; piece.Name != 'q' {
		t.Errorf("Pawn failed to promote properly, resulted in %+v", piece)
	}
}

func TestLegalMoves(t *testing.T) {
	board := &Board{Turn: 1}
	board.PlacePiece('r', 1, 2, 1)
	board.PlacePiece('p', 1, 2, 2)
	board.PlacePiece('n', -1, 5, 1)
	board.PlacePiece('p', 1, 1, 3)
	board.PlacePiece('p', -1, 3, 3)
	rookmoves := make([]*Move, 0)
	for x := 1; x <= 5; x++ {
		if x != 2 {
			m := board.Board[0].makeMoveTo(x, 1)
			rookmoves = append(rookmoves, m)
		}
	}
	rooklegalmoves := board.Board[0].legalMoves(board, false)
	if len(rooklegalmoves) != len(rookmoves) {
		t.Errorf("Size of rook legal moves do not match, %d generated manually vs %d generated automatically", len(rookmoves), len(rooklegalmoves))
	}
	pawnmoves := make([]*Move, 0)
	m := board.Board[1].makeMoveTo(2, 3)
	pawnmoves = append(pawnmoves, m)
	m = board.Board[1].makeMoveTo(3, 3)
	pawnmoves = append(pawnmoves, m)
	m = board.Board[1].makeMoveTo(2, 4)
	pawnmoves = append(pawnmoves, m)
	pawnlegalmoves := board.Board[1].legalMoves(board, false)
	for i, m := range pawnmoves {
		if m.ToString() != pawnlegalmoves[i].ToString() {
			t.Errorf("Pawn legal moves failure on move %s when should have been %s", m.ToString(), pawnlegalmoves[i].ToString())
		}
	}
	board.PlacePiece('p', 1, 6, 6)
	board.Board[len(board.Board)-1].Captured = true
	if moves := board.Board[len(board.Board)-1].legalMoves(board, false); len(moves) != 0 {
		t.Error("Captured piece has legal moves")
	}
	board = &Board{Turn: 1}
	board.PlacePiece('p', -1, 2, 5)
	board.Board[0].Can_en_passant = true
	board.PlacePiece('p', 1, 3, 5)
	if numlegalmoves := len(board.Board[1].legalMoves(board, false)); numlegalmoves != 2 {
		t.Error("En passant not recognized as legal move")
	}
	board = &Board{Turn: 1}
	board.PlacePiece('p', 1, 1, 7)
	if numlegalmoves := len(board.Board[0].legalMoves(board, false)); numlegalmoves == 1 {
		t.Error("Only one legal move recognized for promoting pawn")
	}
	board = &Board{Turn: 1}
	board.PlacePiece('k', 1, 1, 1)
	if numlegalmoves := len(board.Board[0].legalMoves(board, true)); numlegalmoves != 3 {
		t.Errorf("%d moves generated for king in corner", numlegalmoves)
	}
	board = &Board{Turn: 1}
	board.PlacePiece('k', 1, 5, 1)
	board.PlacePiece('r', 1, 1, 1)
	board.PlacePiece('r', 1, 8, 1)
	for i, _ := range board.Board {
		board.Board[i].Can_castle = true
	}
	castles := 0
	for _, m := range board.Board[0].legalMoves(board, false) {
		if m.End.X == 3 || m.End.X == 7 {
			castles += 1
		}
	}
	if castles != 2 {
		t.Errorf("The wrong amount of valid castles were found. Expected 2, got %d", castles)
	}
	board = &Board{Turn: 1}
	board.PlacePiece('p', 1, 2, 2)
	board.PlacePiece('p', -1, 2, 3)
	if numlegalmoves := len(board.Board[0].legalMoves(board, true)); numlegalmoves != 0 {
		t.Errorf("Blocked pawn still had %d legal move(s)", numlegalmoves)
	}
	board = &Board{Turn: 1}
	board.PlacePiece('b', 1, 1, 1)
	board.PlacePiece('n', -1, 3, 3)
	if m := board.Board[0].legalMoves(board, true)[1]; m.Capture != 'n' {
		t.Errorf("Bishop capturing knight had capture %s", string(m.Capture))
	}
	board.PlacePiece('p', 1, 1, 2)
	board.Turn = -1
	for _, m := range board.Board[1].legalMoves(board, true) {
		if m.End.X == 1 && m.End.Y == 2 {
			if m.Capture == 'p' {
				break
			} else {
				t.Errorf("Knight capturing pawn gave capture %s", string(m.Capture))
				break
			}
		}
	}
	board.Board[2].Captured = true
	for _, m := range board.Board[1].legalMoves(board, true) {
		if m.End.X == 1 && m.End.Y == 2 {
			if m.Capture != 0 {
				break
				t.Errorf("Capturing a previously captured piece gave capture %s", string(m.Capture))
			} else {
				break
			}
		}
	}
	board.Board[2].Captured = false
	board.Turn = 1
	board.PlacePiece('q', -1, 2, 3)
	for _, m := range board.Board[2].legalMoves(board, true) {
		if m.End.X == 2 && m.End.Y == 3 {
			if m.Capture == 'q' {
				break
			} else {
				t.Errorf("Pawn capturing queen gave capture %s", string(m.Capture))
			}
		}
	}
}

func TestCanCastle(t *testing.T) {
	board := &Board{Turn: 1}
	board.PlacePiece('k', 1, 5, 1)
	board.Board[0].Can_castle = true
	board.PlacePiece('r', 1, 8, 1)
	board.Board[1].Can_castle = true
	board.PlacePiece('b', 1, 6, 1)
	if board.can_castle(8) {
		t.Error("Castle allowed through blocking piece")
	}
	board.Board[2].Color = -1
	board.Board[2].Position.Y = 2
	if board.can_castle(8) {
		t.Error("Castle allowed when king in check")
	}
	board.Board[2].Position.X = 5
	board.Board[2].Position.Y = 3
	if board.can_castle(8) {
		t.Error("Castle allowed when king placed in check")
	}
	board.Board[2].Color = 1
	board.Board[0].Can_castle = false
	if board.can_castle(8) {
		t.Error("Castle allowed after king moved")
	}
	board.Board[0].Can_castle = true
	board.Board[1].Can_castle = false
	if board.can_castle(8) {
		t.Error("Castle allowed after rook move")
	}
	board.Board[1].Can_castle = true
	board.Board[1].Position.Y = 2
	if board.can_castle(8) {
		t.Error("Castle allowed when rook out of position")
	}
	board.Board[1].Position.Y = 1
	if !board.can_castle(8) {
		t.Error("Error when making a legal castle")
	}
}

func TestToFen(t *testing.T) {
	board := &Board{Turn: 1}
	board.SetUpPieces()
	if fen := board.ToFen(); fen != "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w" {
		t.Errorf("Initial position expected fen:\n%s\nInstead got:\n%s\n", "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w", fen)
	}
	m := &Move{
		Piece: 'p',
		Begin: Square{
			X: 5,
			Y: 2,
		},
		End: Square{
			X: 5,
			Y: 4,
		},
	}
	board.ForceMove(m)
	if fen := board.ToFen(); fen != "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b" {
		t.Errorf("After 1.e4 expected fen:\n%s\nInstead got:\n%s\n", "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b", fen)
	}
}
