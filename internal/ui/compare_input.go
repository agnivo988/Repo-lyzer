package ui

import (
	"fmt"

	"github.com/agnivo988/Repo-lyzer/internal/github"
	tea "github.com/charmbracelet/bubbletea"
)

type CompareReposMsg struct {
	Repo1 string
	Repo2 string
}

type CompareInputModel struct {
	step   int
	repo1  string
	repo2  string
	cursor int
	err    error
}

func NewCompareInputModel() CompareInputModel {
	return CompareInputModel{
		step:   0,
		repo1:  "",
		repo2:  "",
		cursor: 0,
	}
}

func (m CompareInputModel) Init() tea.Cmd {
	return nil
}

func (m CompareInputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if m.step == 0 {
				owner, repo, err := github.ParseGitHubURL(m.repo1)
				if err != nil {
					m.err = err
				} else {
					m.repo1 = owner + "/" + repo
					m.err = nil
					m.step = 1
				}
			} else if m.step == 1 {
				owner, repo, err := github.ParseGitHubURL(m.repo2)
				if err != nil {
					m.err = err
				} else {
					m.repo2 = owner + "/" + repo
					m.err = nil
					return m, func() tea.Msg {
						return CompareReposMsg{Repo1: m.repo1, Repo2: m.repo2}
					}
				}
			}
		case "esc":
			return m, func() tea.Msg {
				return BackToMenuMsg{}
			}
		case "backspace":
			if m.step == 0 && len(m.repo1) > 0 {
				m.repo1 = m.repo1[:len(m.repo1)-1]
			} else if m.step == 1 && len(m.repo2) > 0 {
				m.repo2 = m.repo2[:len(m.repo2)-1]
			}
			m.err = nil
		case "ctrl+u":
			if m.step == 0 {
				m.repo1 = ""
			} else {
				m.repo2 = ""
			}
			m.err = nil
		default:
			if len(msg.String()) == 1 {
				if m.step == 0 {
					m.repo1 += msg.String()
				} else {
					m.repo2 += msg.String()
				}
				m.err = nil
			}
		}
	}
	return m, nil
}

func (m CompareInputModel) View() string {
	var currentInput string
	var prompt string

	if m.step == 0 {
		prompt = "📥 ENTER FIRST REPOSITORY"
		currentInput = m.repo1
	} else {
		prompt = "📥 ENTER SECOND REPOSITORY"
		currentInput = m.repo2
	}

	inputContent := TitleStyle.Render(prompt) + "\n\n"

	if m.step == 1 {
		inputContent += SubtleStyle.Render("First: "+m.repo1) + "\n\n"
	}

	inputContent += InputStyle.Render("> "+currentInput) + "\n\n"
	inputContent += SubtleStyle.Render("Format: owner/repo  •  Press Enter to continue  •  ESC to go back")

	if m.err != nil {
		inputContent += "\n\n" + ErrorStyle.Render(fmt.Sprintf("Error: %v", m.err))
	}

	box := BoxStyle.Render(inputContent)

	return box
}
