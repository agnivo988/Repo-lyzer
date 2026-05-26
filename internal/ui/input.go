package ui

import (
	"fmt"
	"strings"

	"github.com/agnivo988/Repo-lyzer/internal/github"
	tea "github.com/charmbracelet/bubbletea"
)

type InputModel struct {
	input string
	err   error
}

func NewInputModel() InputModel {
	return InputModel{}
}

func (m InputModel) Init() tea.Cmd {
	return nil
}

func (m InputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			owner, repo, err := github.ParseGitHubURL(m.input)
			if err != nil {
				m.err = err
			} else {
				m.input = owner + "/" + repo
				m.err = nil
				return m, func() tea.Msg { return AnalyzeRepoMsg{repoName: m.input} }
			}
		case tea.KeyBackspace:
			if len(m.input) > 0 {
				m.input = m.input[:len(m.input)-1]
				m.err = nil
			}
		case tea.KeyRunes:
			m.input += string(msg.Runes)
			m.err = nil
		case tea.KeyEsc:
			m.err = nil
			return m, func() tea.Msg { return BackToMenuMsg{} }
		case tea.KeyCtrlU:
			m.input = ""
			m.err = nil
		case tea.KeyCtrlW:
			m.input = strings.TrimRight(m.input, " ")
			if idx := strings.LastIndex(m.input, " "); idx >= 0 {
				m.input = m.input[:idx+1]
			} else {
				m.input = ""
			}
			m.err = nil
		}
	}
	return m, nil
}

func (m InputModel) View() string {
	inputContent :=
		TitleStyle.Render("📥 ENTER REPOSITORY") + "\n\n" +
			InputStyle.Render("> "+m.input) + "\n\n" +
			SubtleStyle.Render("Format: owner/repo or GitHub URL  •  Press Enter to analyze")

	if m.err != nil {
		inputContent += "\n\n" + ErrorStyle.Render(fmt.Sprintf("Error: %v", m.err))
	}

	box := BoxStyle.Render(inputContent)

	return box
}

func (m *InputModel) SetInput(input string) {
	m.input = input
}

func (m *InputModel) GetInput() string {
	return m.input
}

func (m *InputModel) ClearError() {
	m.err = nil
}

func (m *InputModel) SetError(err error) {
	m.err = err
}

func sanitizeRepoInput(input string) string {
	input = strings.TrimSpace(input)
	if input == "" {
		return ""
	}
	// Handle GitHub URLs
	if strings.HasPrefix(input, "https://github.com/") {
		parts := strings.Split(input, "/")
		if len(parts) >= 5 {
			return parts[3] + "/" + parts[4]
		}
	}
	// Handle owner/repo format
	if strings.Contains(input, "/") && len(strings.Split(input, "/")) == 2 {
		return input
	}
	return ""
}
