// main.go
package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// main is the entry point of the program.
func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

type gameState int

const (
	nameInput gameState = iota
	gamePlaying
)

// model represents the state of our Tic-Tac-Toe game.
type model struct {
	board         [3][3]string // Represents the 3x3 game board
	cursorX       int          // The cursor's X position (column)
	cursorY       int          // The cursor's Y position (row)
	player        string       // The current player ("X" or "O")
	winner        string       // The winner of the game, if any
	isDraw        bool         // True if the game is a draw
	player1Name   string
	player2Name   string
	player1Score  int
	player2Score  int
	inputs        []textinput.Model
	focusIndex    int
	gameState     gameState
	winningCells  []struct{ x, y int }
}

// initialModel creates the initial state of the game.
func initialModel() model {
	m := model{
		board: [3][3]string{
			{" ", " ", " "},
			{" ", " ", " "},
			{" ", " ", " "},
		},
		cursorX:     0,
		cursorY:     0,
		player:      "X",
		winner:      "",
		isDraw:      false,
		gameState:   nameInput,
		inputs:      make([]textinput.Model, 2),
		focusIndex:  0,
		winningCells: []struct{ x, y int }{},
	}

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.Cursor.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
		t.CharLimit = 32
		t.Placeholder = fmt.Sprintf("Player %d", i+1)
		if i == 0 {
			t.Focus()
			t.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
		}
		m.inputs[i] = t
	}

	return m
}

func (m model) resetGame() model {
	m.board = [3][3]string{
		{" ", " ", " "},
		{" ", " ", " "},
		{" ", " ", " "},
	}
	m.cursorX = 0
	m.cursorY = 0
	m.player = "X"
	m.winner = ""
	m.isDraw = false
	m.winningCells = []struct{ x, y int }{}
	return m
}

// Init is called once when the program starts.
func (m model) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles incoming messages and updates the model accordingly.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	switch m.gameState {
	case nameInput:
		return updateNameInput(msg, m)
	case gamePlaying:
		return updateGamePlaying(msg, m)
	default:
		return m, nil
	}
}

func updateNameInput(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if m.focusIndex == len(m.inputs)-1 {
				m.player1Name = m.inputs[0].Value()
				m.player2Name = m.inputs[1].Value()
				m.gameState = gamePlaying
				return m, nil
			}
			m.focusIndex++
			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focusIndex {
					m.inputs[i].Focus()
					m.inputs[i].PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
				} else {
					m.inputs[i].Blur()
					m.inputs[i].PromptStyle = lipgloss.NewStyle()
				}
			}
			return m, nil
		case "up":
			if m.focusIndex > 0 {
				m.focusIndex--
			}
			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focusIndex {
					m.inputs[i].Focus()
					m.inputs[i].PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
				} else {
					m.inputs[i].Blur()
					m.inputs[i].PromptStyle = lipgloss.NewStyle()
				}
			}
		case "down":
			if m.focusIndex < len(m.inputs)-1 {
				m.focusIndex++
			}
			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focusIndex {
					m.inputs[i].Focus()
					m.inputs[i].PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
				} else {
					m.inputs[i].Blur()
					m.inputs[i].PromptStyle = lipgloss.NewStyle()
				}
			}
		}
	}

	cmd := m.updateInputs(msg)
	return m, cmd
}

func (m *model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}
	return tea.Batch(cmds...)
}

func updateGamePlaying(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+r":
			return initialModel(), nil
		case "r":
			return m.resetGame(), nil
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
		case "enter", " ":
			if m.winner != "" || m.isDraw {
				return m.resetGame(), nil
			}

			if m.board[m.cursorY][m.cursorX] == " " {
				m.board[m.cursorY][m.cursorX] = m.player
				won, cells := checkWinner(m.board, m.player)
				if won {
					m.winner = m.player
					m.winningCells = cells
					if m.player == "X" {
						m.player1Score++
					} else {
						m.player2Score++
					}
				} else if checkDraw(m.board) {
					m.isDraw = true
				} else {
					if m.player == "X" {
						m.player = "O"
					} else {
						m.player = "X"
					}
				}
			}
		}
	}
	return m, nil
}

// View renders the UI.
func (m model) View() string {
	if m.gameState == nameInput {
		return viewNameInput(m)
	}
	return viewGamePlaying(m)
}

func viewNameInput(m model) string {
	var b strings.Builder
	b.WriteString("Enter Player Names\n\n")
	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}
	b.WriteString("\n\nPress Enter to continue.")
	return b.String()
}

func viewGamePlaying(m model) string {
	s := "Tic-Tac-Toe\n\n"
	s += fmt.Sprintf("Score: %s (X) %d - %d %s (O)\n\n", m.player1Name, m.player1Score, m.player2Score, m.player2Name)

	var boardView string
	var rows []string

	for i := 0; i < 3; i++ {
		var rowItems []string
		for j := 0; j < 3; j++ {
			cell := m.board[i][j]
			style := lipgloss.NewStyle().
				Border(lipgloss.NormalBorder(), true).
				BorderForeground(lipgloss.Color("63")).
				Width(5).
				Height(1).
				Align(lipgloss.Center, lipgloss.Center)

			if m.cursorY == i && m.cursorX == j {
				style = style.Copy().BorderForeground(lipgloss.Color("205"))
			}

			for _, winningCell := range m.winningCells {
				if winningCell.x == j && winningCell.y == i {
					style = style.Copy().Foreground(lipgloss.Color("196"))
				}
			}

			var renderedCell string
			if cell == "X" {
				renderedCell = style.Copy().Foreground(lipgloss.Color("202")).Render(cell)
			} else if cell == "O" {
				renderedCell = style.Copy().Foreground(lipgloss.Color("39")).Render(cell)
			} else {
				renderedCell = style.Render(cell)
			}
			rowItems = append(rowItems, renderedCell)
		}
		rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Top, rowItems...))
	}

	boardView = lipgloss.JoinVertical(lipgloss.Left, rows...)
	s += boardView

	if m.winner != "" {
		winnerName := m.player1Name
		if m.winner == "O" {
			winnerName = m.player2Name
		}
		s += fmt.Sprintf("\n\n%s wins! (Press Enter to play again)", winnerName)
	} else if m.isDraw {
		s += "\n\nIt's a draw! (Press Enter to play again)"
	} else {
		playerName := m.player1Name
		if m.player == "O" {
			playerName = m.player2Name
		}
		s += fmt.Sprintf("\n\n%s's turn (%s)", playerName, m.player)
	}

	s += "\n\nUse arrow keys or h/j/k/l to move.\n"
	s += "Press Enter or Space to place your marker.\n"
	s += "Press 'r' to reset the game.\n"
	s += "Press 'ctrl+r' to reset scores and names.\n"
	s += "Press 'q' or 'ctrl+c' to quit.\n"

	return s
}

// checkWinner checks if the given player has won the game.
func checkWinner(board [3][3]string, player string) (bool, []struct{ x, y int }) {
	// Check rows
	for i := 0; i < 3; i++ {
		if board[i][0] == player && board[i][1] == player && board[i][2] == player {
			return true, []struct{ x, y int }{{0, i}, {1, i}, {2, i}}
		}
	}

	// Check columns
	for i := 0; i < 3; i++ {
		if board[0][i] == player && board[1][i] == player && board[2][i] == player {
			return true, []struct{ x, y int }{{i, 0}, {i, 1}, {i, 2}}
		}
	}

	// Check diagonals
	if board[0][0] == player && board[1][1] == player && board[2][2] == player {
		return true, []struct{ x, y int }{{0, 0}, {1, 1}, {2, 2}}
	}
	if board[0][2] == player && board[1][1] == player && board[2][0] == player {
		return true, []struct{ x, y int }{{2, 0}, {1, 1}, {0, 2}}
	}

	return false, nil
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
