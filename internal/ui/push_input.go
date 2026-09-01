package ui

import (
	"fmt"
	"os"
	"strings"

	"github.com/agnivo988/Repo-lyzer/internal/pathutil"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type PushInputModel struct {
	Step      int // 0 = local path, 1 = repo name
	LocalPath string
	RepoName  string
	Err       error
}

func NewPushInputModel() PushInputModel {
	return PushInputModel{
		Step: 0,
	}
}

type PushRepoMsg struct {
	LocalPath string
	RepoName  string
}

func (m PushInputModel) Update(msg tea.Msg) (PushInputModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.Step == 0 {
				if m.LocalPath != "" {
					resolvedPath, err := pathutil.CleanAndResolvePath(m.LocalPath)
					if err != nil {
						m.Err = err
					} else {
						// Validate local folder exists
						if stat, err := os.Stat(resolvedPath); err != nil {
							m.Err = fmt.Errorf("local directory does not exist: %s", resolvedPath)
						} else if !stat.IsDir() {
							m.Err = fmt.Errorf("path is not a directory: %s", resolvedPath)
						} else {
							m.Err = nil
							m.LocalPath = resolvedPath
							m.Step = 1
						}
					}
				} else {
					m.Err = fmt.Errorf("local directory path cannot be empty")
				}
			} else {
				if m.RepoName != "" {
					// Validate remote repo format (owner/repo)
					parts := strings.Split(m.RepoName, "/")
					if len(parts) == 2 && parts[0] != "" && parts[1] != "" {
						m.Err = nil
						return m, func() tea.Msg {
							return PushRepoMsg{
								LocalPath: m.LocalPath,
								RepoName:  m.RepoName,
							}
						}
					} else {
						m.Err = fmt.Errorf("invalid format: must be owner/repo")
					}
				} else {
					m.Err = fmt.Errorf("repository name cannot be empty")
				}
			}
		case tea.KeyBackspace:
			if m.Step == 0 {
				if len(m.LocalPath) > 0 {
					runes := []rune(m.LocalPath)
					m.LocalPath = string(runes[:len(runes)-1])
					m.Err = nil
				}
			} else {
				if len(m.RepoName) > 0 {
					runes := []rune(m.RepoName)
					m.RepoName = string(runes[:len(runes)-1])
					m.Err = nil
				}
			}
		case tea.KeyRunes:
			if m.Step == 0 {
				m.LocalPath += string(msg.Runes)
				m.Err = nil
			} else {
				m.RepoName += string(msg.Runes)
				m.Err = nil
			}
		case tea.KeyEsc:
			if m.Step == 1 {
				m.Step = 0
				m.Err = nil
			} else {
				return m, func() tea.Msg { return BackToMenuMsg{} }
			}
		case tea.KeyCtrlU:
			if m.Step == 0 {
				m.LocalPath = ""
			} else {
				m.RepoName = ""
			}
			m.Err = nil
		}
	}
	return m, nil
}

func (m PushInputModel) View(width, height int) string {
	header := TitleStyle.Render("📤 PUSH TO GITHUB")

	var prompt string
	var currentInput string
	if m.Step == 0 {
		prompt = "Enter absolute local folder path to push:\n(e.g., /Users/username/Desktop/my-project)"
		currentInput = m.LocalPath
	} else {
		prompt = fmt.Sprintf("Local Folder: %s\n\nEnter remote GitHub repository (owner/repo):", m.LocalPath)
		currentInput = m.RepoName
	}

	inputContent := fmt.Sprintf(
		"%s\n\n> %s█\n\n"+
			"This folder will be prepared and pushed to GitHub.",
		prompt,
		currentInput,
	)

	var errMsg string
	if m.Err != nil {
		errMsg = "\n" + ErrorStyle.Render(m.Err.Error())
	}

	footer := SubtleStyle.Render("Enter: continue • ESC: back • Ctrl+U: clear")

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
