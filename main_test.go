// main_test.go
package main

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// TestInitialModel verifies that the game starts with the correct default state.
func TestInitialModel(t *testing.T) {
	m := initialModel()

	if m.player != "X" {
		t.Errorf("Expected player to be 'X', but got '%s'", m.player)
	}
	if m.cursorX != 0 || m.cursorY != 0 {
		t.Errorf("Expected cursor at (0, 0), but got (%d, %d)", m.cursorX, m.cursorY)
	}
	if m.winner != "" {
		t.Errorf("Expected no winner at start, but got '%s'", m.winner)
	}
	if m.isDraw {
		t.Error("Expected isDraw to be false at start, but got true")
	}
	for i, row := range m.board {
		for j, cell := range row {
			if cell != " " {
				t.Errorf("Expected board at (%d, %d) to be empty, but got '%s'", i, j, cell)
			}
		}
	}
}

// TestCheckWinner covers all winning scenarios and some non-winning ones.
func TestCheckWinner(t *testing.T) {
	// A map of test cases, with a descriptive name for each.
	testCases := map[string]struct {
		board  [3][3]string
		player string
		want   bool
	}{
		"Row 1 win for X": {
			board:  [3][3]string{{"X", "X", "X"}, {" ", "O", " "}, {"O", " ", " "}},
			player: "X",
			want:   true,
		},
		"Row 2 win for X": {
			board:  [3][3]string{{" ", "O", " "}, {"X", "X", "X"}, {"O", " ", " "}},
			player: "X",
			want:   true,
		},
		"Row 3 win for X": {
			board:  [3][3]string{{" ", "O", " "}, {"O", " ", " "}, {"X", "X", "X"}},
			player: "X",
			want:   true,
		},
		"Column 1 win for O": {
			board:  [3][3]string{{"O", "X", "X"}, {"O", "X", " "}, {"O", " ", " "}},
			player: "O",
			want:   true,
		},
		"Column 2 win for O": {
			board:  [3][3]string{{"X", "O", "X"}, {" ", "O", " "}, {"X", "O", " "}},
			player: "O",
			want:   true,
		},
		"Column 3 win for O": {
			board:  [3][3]string{{"X", "X", "O"}, {" ", " ", "O"}, {"X", " ", "O"}},
			player: "O",
			want:   true,
		},
		"Diagonal 1 win for X": {
			board:  [3][3]string{{"X", "O", "O"}, {" ", "X", " "}, {" ", " ", "X"}},
			player: "X",
			want:   true,
		},
		"Diagonal 2 win for O": {
			board:  [3][3]string{{"X", "X", "O"}, {" ", "O", " "}, {"O", " ", "X"}},
			player: "O",
			want:   true,
		},
		"No winner": {
			board:  [3][3]string{{"X", "O", "X"}, {"O", "X", "O"}, {"O", "X", " "}},
			player: "X",
			want:   false,
		},
		"Empty board": {
			board:  [3][3]string{{" ", " ", " "}, {" ", " ", " "}, {" ", " ", " "}},
			player: "O",
			want:   false,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			got := checkWinner(tc.board, tc.player)
			if got != tc.want {
				t.Errorf("checkWinner() = %v, want %v", got, tc.want)
			}
		})
	}
}

// TestCheckDraw checks for draw and non-draw conditions.
func TestCheckDraw(t *testing.T) {
	testCases := map[string]struct {
		board [3][3]string
		want  bool
	}{
		"Full board is a draw": {
			board: [3][3]string{{"X", "O", "X"}, {"O", "X", "X"}, {"O", "X", "O"}},
			want:  true,
		},
		"Full board with a winner is still considered a full board": {
			board: [3][3]string{{"X", "O", "X"}, {"O", "X", "O"}, {"X", "X", "X"}},
			want:  true,
		},
		"Board with empty space is not a draw": {
			board: [3][3]string{{"X", "O", "X"}, {"O", " ", "X"}, {"O", "X", "O"}},
			want:  false,
		},
		"Empty board is not a draw": {
			board: [3][3]string{{" ", " ", " "}, {" ", " ", " "}, {" ", " ", " "}},
			want:  false,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			got := checkDraw(tc.board)
			if got != tc.want {
				t.Errorf("checkDraw() = %v, want %v", got, tc.want)
			}
		})
	}
}

// TestUpdatePlayerMove tests the core game logic of making a move.
func TestUpdatePlayerMove(t *testing.T) {
	m := initialModel()
	var updatedModel tea.Model

	// Player X makes a move at (0, 0)
	updatedModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updatedModel.(model)

	if m.board[0][0] != "X" {
		t.Errorf("Expected board at (0,0) to be 'X', but got '%s'", m.board[0][0])
	}
	if m.player != "O" {
		t.Errorf("Expected player to switch to 'O', but got '%s'", m.player)
	}

	// Try to move on the same spot (should not work)
	updatedModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updatedModel.(model)
	if m.player != "O" {
		t.Errorf("Player should remain 'O' after invalid move, but got '%s'", m.player)
	}

	// Player O makes a move
	m.cursorX = 1
	m.cursorY = 1
	updatedModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updatedModel.(model)

	if m.board[1][1] != "O" {
		t.Errorf("Expected board at (1,1) to be 'O', but got '%s'", m.board[1][1])
	}
	if m.player != "X" {
		t.Errorf("Expected player to switch back to 'X', but got '%s'", m.player)
	}
}

// TestUpdateCursorMovement tests that the cursor moves correctly.
func TestUpdateCursorMovement(t *testing.T) {
	m := initialModel()
	var updatedModel tea.Model

	// Move down
	updatedModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown})
	m = updatedModel.(model)
	if m.cursorY != 1 {
		t.Errorf("Expected cursorY to be 1 after moving down, got %d", m.cursorY)
	}

	// Move right
	updatedModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyRight})
	m = updatedModel.(model)
	if m.cursorX != 1 {
		t.Errorf("Expected cursorX to be 1 after moving right, got %d", m.cursorX)
	}

	// Move up
	updatedModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyUp})
	m = updatedModel.(model)
	if m.cursorY != 0 {
		t.Errorf("Expected cursorY to be 0 after moving up, got %d", m.cursorY)
	}

	// Move left
	updatedModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyLeft})
	m = updatedModel.(model)
	if m.cursorX != 0 {
		t.Errorf("Expected cursorX to be 0 after moving left, got %d", m.cursorX)
	}

	// Test boundaries (should not move past edge)
	updatedModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyUp})
	m = updatedModel.(model)
	if m.cursorY != 0 {
		t.Errorf("Expected cursorY to stay at 0 at top boundary, got %d", m.cursorY)
	}
}

// TestUpdateReset tests the game reset functionality.
func TestUpdateReset(t *testing.T) {
	m := initialModel()
	var updatedModel tea.Model

	// Make a move
	updatedModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updatedModel.(model)

	// Reset the game
	updatedModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})
	m = updatedModel.(model)

	// Check if the model is back to its initial state
	initial := initialModel()
	if m.player != initial.player {
		t.Errorf("Expected player to be '%s' after reset, but got '%s'", initial.player, m.player)
	}
	if m.board != initial.board {
		t.Errorf("Expected board to be empty after reset, but it was not")
	}
}

// TestUpdateQuit tests that the quit command works.
func TestUpdateQuit(t *testing.T) {
	m := initialModel()
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyCtrlC})

	if cmd == nil {
		t.Error("Expected a tea.Quit command, but got nil")
	}
	// We can't directly check if the command is tea.Quit,
	// but we can check if it's not nil. A more robust way would
	// require a custom command type or inspecting the command's function.
}

// TestUpdateWinCondition checks that the game correctly identifies a winner.
func TestUpdateWinCondition(t *testing.T) {
	m := initialModel()
	m.board = [3][3]string{
		{"X", "X", " "},
		{"O", "O", " "},
		{" ", " ", " "},
	}
	m.cursorX = 2
	m.cursorY = 0
	m.player = "X"

	updatedModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updatedModel.(model)

	if m.winner != "X" {
		t.Errorf("Expected winner to be 'X', but got '%s'", m.winner)
	}
}

// TestUpdateDrawCondition checks that the game correctly identifies a draw.
func TestUpdateDrawCondition(t *testing.T) {
	m := initialModel()
	m.board = [3][3]string{
		{"X", "O", "X"},
		{"X", "O", "O"},
		{"O", "X", " "},
	}
	m.cursorX = 2
	m.cursorY = 2
	m.player = "X"

	updatedModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updatedModel.(model)

	if !m.isDraw {
		t.Error("Expected isDraw to be true, but got false")
	}
	if m.winner != "" {
		t.Errorf("Expected no winner in a draw, but got '%s'", m.winner)
	}
}

// TestUpdateResetAfterWin tests that the game resets after a win.
func TestUpdateResetAfterWin(t *testing.T) {
	m := initialModel()
	m.winner = "X"

	updatedModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updatedModel.(model)

	initial := initialModel()
	if m.player != initial.player {
		t.Errorf("Expected player to be '%s' after reset, but got '%s'", initial.player, m.player)
	}
	if m.board != initial.board {
		t.Errorf("Expected board to be empty after reset, but it was not")
	}
	if m.winner != "" {
		t.Errorf("Expected winner to be empty after reset, but got '%s'", m.winner)
	}
}

// TestUpdateResetAfterDraw tests that the game resets after a draw.
func TestUpdateResetAfterDraw(t *testing.T) {
	m := initialModel()
	m.isDraw = true

	updatedModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updatedModel.(model)

	initial := initialModel()
	if m.player != initial.player {
		t.Errorf("Expected player to be '%s' after reset, but got '%s'", initial.player, m.player)
	}
	if m.board != initial.board {
		t.Errorf("Expected board to be empty after reset, but it was not")
	}
	if m.isDraw {
		t.Error("Expected isDraw to be false after reset, but got true")
	}
}

// TestUpdateAlternateKeys tests alternate key bindings for controls.
func TestUpdateAlternateKeys(t *testing.T) {
	m := initialModel()
	var updatedModel tea.Model

	// 'k' for up
	m.cursorY = 1
	updatedModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	m = updatedModel.(model)
	if m.cursorY != 0 {
		t.Errorf("Expected cursorY to be 0 after pressing 'k', got %d", m.cursorY)
	}

	// 'j' for down
	updatedModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	m = updatedModel.(model)
	if m.cursorY != 1 {
		t.Errorf("Expected cursorY to be 1 after pressing 'j', got %d", m.cursorY)
	}

	// 'h' for left
	m.cursorX = 1
	updatedModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}})
	m = updatedModel.(model)
	if m.cursorX != 0 {
		t.Errorf("Expected cursorX to be 0 after pressing 'h', got %d", m.cursorX)
	}

	// 'l' for right
	updatedModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}})
	m = updatedModel.(model)
	if m.cursorX != 1 {
		t.Errorf("Expected cursorX to be 1 after pressing 'l', got %d", m.cursorX)
	}

	// ' ' (space) for placing a marker
	m.cursorX = 2
	m.cursorY = 2
	updatedModel, _ = m.Update(tea.KeyMsg{Type: tea.KeySpace})
	m = updatedModel.(model)
	if m.board[2][2] != "X" {
		t.Errorf("Expected board at (2,2) to be 'X' after pressing space, but got '%s'", m.board[2][2])
	}

	// 'q' for quit
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	if cmd == nil {
		t.Error("Expected a tea.Quit command after pressing 'q', but got nil")
	}
}

// TestView checks the rendered output of the model's View function.
func TestView(t *testing.T) {
	t.Run("Initial view", func(t *testing.T) {
		m := initialModel()
		view := m.View()

		if !contains(view, "Tic-Tac-Toe") {
			t.Errorf("View does not contain 'Tic-Tac-Toe'")
		}
		if !contains(view, "Player X's turn") {
			t.Errorf("View does not contain 'Player X's turn'")
		}
	})

	t.Run("Win view", func(t *testing.T) {
		m := initialModel()
		m.winner = "O"
		view := m.View()

		if !contains(view, "Player O wins!") {
			t.Errorf("View does not contain 'Player O wins!'")
		}
	})

	t.Run("Draw view", func(t *testing.T) {
		m := initialModel()
		m.isDraw = true
		view := m.View()

		if !contains(view, "It's a draw!") {
			t.Errorf("View does not contain 'It's a draw!'")
		}
	})

	t.Run("View with markers", func(t *testing.T) {
		m := initialModel()
		m.board[0][0] = "X"
		m.board[1][1] = "O"
		view := m.View()

		// A simple check to see if the markers are in the output.
		// A more robust test could check the exact rendered output,
		// but that would be brittle and dependent on lipgloss rendering.
		if !contains(view, "X") {
			t.Errorf("View does not contain marker 'X'")
		}
		if !contains(view, "O") {
			t.Errorf("View does not contain marker 'O'")
		}
	})
}

// contains is a helper function to check if a string contains a substring.
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
