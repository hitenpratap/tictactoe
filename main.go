// main.go
package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// main is the entry point of the program.
func main() {
	// Initialize the Bubble Tea program with our model
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

// model represents the state of our Tic-Tac-Toe game.
type model struct {
	board   [3][3]string // Represents the 3x3 game board
	cursorX int          // The cursor's X position (column)
	cursorY int          // The cursor's Y position (row)
	player  string       // The current player ("X" or "O")
	winner  string       // The winner of the game, if any
	isDraw  bool         // True if the game is a draw
}

// initialModel creates the initial state of the game.
func initialModel() model {
	return model{
		board: [3][3]string{
			{" ", " ", " "},
			{" ", " ", " "},
			{" ", " ", " "},
		},
		cursorX: 0,
		cursorY: 0,
		player:  "X",
		winner:  "",
		isDraw:  false,
	}
}

// Init is called once when the program starts.
func (m model) Init() tea.Cmd {
	// We don't need to do anything here for this simple game.
	return nil
}

// Update handles incoming messages and updates the model accordingly.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Make sure we're handling a key press
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle key presses
		switch msg.String() {
		// Quit the game
		case "ctrl+c", "q":
			return m, tea.Quit

		// Movement keys
		case "up", "k":
			if m.cursorY > 0 {
				m.cursorY--
			}
		case "down", "j":
			if m.cursorY < 2 {
				m.cursorY++
			}
		case "left", "h":
			if m.cursorX > 0 {
				m.cursorX--
			}
		case "right", "l":
			if m.cursorX < 2 {
				m.cursorX++
			}

		// Place a marker
		case "enter", " ":
			// If there's a winner or a draw, enter will reset the game
			if m.winner != "" || m.isDraw {
				return initialModel(), nil
			}

			// Check if the selected cell is empty
			if m.board[m.cursorY][m.cursorX] == " " {
				// Place the player's marker
				m.board[m.cursorY][m.cursorX] = m.player

				// Check for a winner
				if checkWinner(m.board, m.player) {
					m.winner = m.player
				} else if checkDraw(m.board) {
					m.isDraw = true
				} else {
					// Switch players
					if m.player == "X" {
						m.player = "O"
					} else {
						m.player = "X"
					}
				}
			}
		// Reset the game
		case "r":
			return initialModel(), nil
		}
	}

	return m, nil
}

// View renders the UI.
func (m model) View() string {
	s := "Tic-Tac-Toe\n\n"

	// This will hold the rendered board
	var boardView string
	// This will hold the rendered rows, which will be joined vertically
	var rows []string

	// Iterate over the board rows
	for i := 0; i < 3; i++ {
		// This will hold the rendered cells for a single row
		var rowItems []string
		// Iterate over the board columns
		for j := 0; j < 3; j++ {
			cell := m.board[i][j]
			// Define the style for a single cell
			style := lipgloss.NewStyle().
				Border(lipgloss.NormalBorder(), true).
				BorderForeground(lipgloss.Color("63")).
				Width(5).
				Height(1).
				Align(lipgloss.Center, lipgloss.Center)

			// Highlight the cell under the cursor
			if m.cursorY == i && m.cursorX == j {
				style = style.Copy().BorderForeground(lipgloss.Color("205"))
			}

			var renderedCell string
			// Color the player markers
			if cell == "X" {
				renderedCell = style.Copy().Foreground(lipgloss.Color("202")).Render(cell)
			} else if cell == "O" {
				renderedCell = style.Copy().Foreground(lipgloss.Color("39")).Render(cell)
			} else {
				renderedCell = style.Render(cell)
			}
			// Add the rendered cell to the row
			rowItems = append(rowItems, renderedCell)
		}
		// Join the cells horizontally to form a row
		rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Top, rowItems...))
	}

	// Join the rows vertically to form the final board
	boardView = lipgloss.JoinVertical(lipgloss.Left, rows...)
	s += boardView

	// Display game status
	if m.winner != "" {
		s += fmt.Sprintf("\n\nPlayer %s wins! (Press Enter to play again)", m.winner)
	} else if m.isDraw {
		s += "\n\nIt's a draw! (Press Enter to play again)"
	} else {
		s += fmt.Sprintf("\n\nPlayer %s's turn", m.player)
	}

	s += "\n\nUse arrow keys or h/j/k/l to move.\n"
	s += "Press Enter or Space to place your marker.\n"
	s += "Press 'r' to reset the game.\n"
	s += "Press 'q' or 'ctrl+c' to quit.\n"

	return s
}

// checkWinner checks if the given player has won the game.
func checkWinner(board [3][3]string, player string) bool {
	// Check rows
	for i := 0; i < 3; i++ {
		if board[i][0] == player && board[i][1] == player && board[i][2] == player {
			return true
		}
	}

	// Check columns
	for i := 0; i < 3; i++ {
		if board[0][i] == player && board[1][i] == player && board[2][i] == player {
			return true
		}
	}

	// Check diagonals
	if board[0][0] == player && board[1][1] == player && board[2][2] == player {
		return true
	}
	if board[0][2] == player && board[1][1] == player && board[2][0] == player {
		return true
	}

	return false
}

// checkDraw checks if the game is a draw.
func checkDraw(board [3][3]string) bool {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if board[i][j] == " " {
				return false // There's an empty cell, so not a draw
			}
		}
	}
	return true // All cells are filled
}
