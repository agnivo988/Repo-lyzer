package ui

import (
	"fmt"
	"strings"

	"github.com/agnivo988/Repo-lyzer/internal/github"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type CloneInputModel struct {
	input string
	err   error
}

func NewCloneInputModel() CloneInputModel {
	return CloneInputModel{}
}

func (m CloneInputModel) Update(msg tea.Msg) (CloneInputModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.input != "" {
				owner, repo, err := github.ParseGitHubURL(m.input)
				if err != nil {
					m.err = err
				} else {
					m.err = nil
					return m, func() tea.Msg { return CloneRepoMsg{repoName: owner + "/" + repo} }
				}
			}
		case tea.KeyBackspace:
			if len(m.input) > 0 {
				m.input = m.input[:len(m.input)-1]
				m.err = nil // Clear error on input change
			}
		case tea.KeyRunes:
			m.input += string(msg.Runes)
			m.err = nil // Clear error on input change
		case tea.KeyEsc:
			return m, func() tea.Msg { return BackToMenuMsg{} }
		case tea.KeyCtrlU:
			m.input = ""
			m.err = nil // Clear error on clear
		}
	}
	return m, nil
}

func (m CloneInputModel) View(width, height int) string {
	header := TitleStyle.Render("📥 CLONE REPOSITORY")

	inputContent := fmt.Sprintf(
		"Enter repository to clone (owner/repo):\n\n> %s█\n\n"+
			"The repository will be cloned to your Desktop folder.",
		m.input,
	)

	var errMsg string
	if m.err != nil {
		errMsg = "\n" + ErrorStyle.Render(m.err.Error())
	}

	footer := SubtleStyle.Render("Enter: clone • ESC: back • Ctrl+U: clear")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		BoxStyle.Render(inputContent),
		errMsg,
		footer,
	)

	if width == 0 || height == 0 {
		return ""
	}

	return lipgloss.Place(
		width,
		height,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)
}

type CloneRepoMsg struct {
	repoName string
}
