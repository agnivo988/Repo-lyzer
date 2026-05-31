package ui

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"unicode"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type FileEditModel struct {
	filePath      string
	repoOwner     string
	repoName      string
	defaultBranch string
	isOwner       bool
	width         int
	height        int
	Done          bool
	statusMsg     string
	clonePath     string
	isCloned      bool
}

func NewFileEditModel(filePath, repoFullName, defaultBranch string) FileEditModel {
	parts := strings.Split(repoFullName, "/")
	repoOwner := ""
	repoName := ""
	if len(parts) == 2 {
		repoOwner = parts[0]
		repoName = parts[1]
	}

	// Check if repo is already cloned to Desktop
	desktopPath := getDesktopPath()
	clonePath := filepath.Join(desktopPath, repoName)
	isCloned := false
	if _, err := os.Stat(filepath.Join(clonePath, ".git")); err == nil {
		isCloned = true
	}

	if defaultBranch == "" {
		defaultBranch = "main"
	}

	return FileEditModel{
		filePath:      filePath,
		repoOwner:     repoOwner,
		repoName:      repoName,
		defaultBranch: defaultBranch,
		clonePath:     clonePath,
		isCloned:      isCloned,
	}
}

func (m *FileEditModel) SetOwnership(isOwner bool) {
	m.isOwner = isOwner
}

func (m FileEditModel) Init() tea.Cmd { return nil }

func (m FileEditModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case cloneResultMsg:
		if msg.err != nil {
			m.statusMsg = fmt.Sprintf("Clone failed: %v", msg.err)
		} else {
			m.statusMsg = "✓ Repository cloned successfully!"
			m.isCloned = true
		}

	case openResultMsg:
		if msg.err != nil {
			m.statusMsg = fmt.Sprintf("Failed to open: %v", msg.err)
		} else {
			m.statusMsg = "✓ Opened in editor"
		}

	case tea.KeyMsg:
		switch msg.String() {
		case "v":
			// View file on GitHub in browser
			return m, m.openInBrowser()
		case "e", "enter":
			// Open in VS Code (new window)
			return m, m.openInVSCode()
		case "c":
			// Clone repository to Desktop
			return m, m.cloneToDesktop()
		case "o":
			// Open cloned repo in VS Code
			if m.isCloned {
				return m, m.openClonedInVSCode()
			}
		case "esc":
			m.Done = true
		}
	}

	return m, nil
}

type cloneResultMsg struct {
	err error
}

type openResultMsg struct {
	err error
}

func (m FileEditModel) View() string {
	content := TitleStyle.Render("📝 FILE VIEWER") + "\n\n"

	content += fmt.Sprintf("File: %s\n", SelectedStyle.Render(m.filePath))
	content += fmt.Sprintf("Repository: %s/%s\n\n", m.repoOwner, m.repoName)

	// Clone status
	if m.isCloned {
		content += SuccessStyle.Render("✓ Repository cloned to Desktop\n")
		content += SubtleStyle.Render(fmt.Sprintf("  Path: %s\n\n", m.clonePath))
	} else {
		content += SubtleStyle.Render("○ Repository not cloned locally\n\n")
	}

	content += "Available actions:\n"
	content += "  [v] View file on GitHub (browser)\n"
	content += "  [e] Open file in VS Code (GitHub URL)\n"
	content += "  [c] Clone repository to Desktop\n"
	if m.isCloned {
		content += SuccessStyle.Render("  [o] Open cloned repo in VS Code (new window)\n")
	} else {
		content += SubtleStyle.Render("  [o] Open cloned repo (clone first)\n")
	}

	if m.statusMsg != "" {
		content += "\n" + InputStyle.Render(m.statusMsg)
	}

	content += "\n\n" + SubtleStyle.Render("ESC back to file tree")

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Left, lipgloss.Top,
		BoxStyle.Render(content),
	)
}

// getDesktopPath returns the user's Desktop folder path
func getDesktopPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "."
	}

	switch runtime.GOOS {
	case "windows":
		return filepath.Join(home, "Desktop")
	case "darwin":
		return filepath.Join(home, "Desktop")
	default:
		// Linux - check for Desktop folder
		desktop := filepath.Join(home, "Desktop")
		if _, err := os.Stat(desktop); err == nil {
			return desktop
		}
		return home
	}
}

// openInBrowser opens the file on GitHub in the default browser
func (m FileEditModel) openInBrowser() tea.Cmd {
	return func() tea.Msg {
		githubURL, err := buildBlobURL("github.com", m.repoOwner, m.repoName, m.defaultBranch, m.filePath)
		if err != nil {
			return openResultMsg{err}
		}

		err = openExternalURL(githubURL)
		return openResultMsg{err}
	}
}

// openInVSCode opens the GitHub file URL in VS Code
func (m FileEditModel) openInVSCode() tea.Cmd {
	return func() tea.Msg {
		// Use vscode.dev to open the file in browser-based VS Code
		vscodeURL, err := buildVSCodeBlobURL(m.repoOwner, m.repoName, m.defaultBranch, m.filePath)
		if err != nil {
			return openResultMsg{err}
		}

		err = openExternalURL(vscodeURL)
		return openResultMsg{err}
	}
}

// buildBlobURL creates a safe GitHub blob URL for a repository file.
func buildBlobURL(host, owner, repo, branch, filePath string) (string, error) {
	owner = strings.TrimSpace(owner)
	repo = strings.TrimSpace(repo)
	branch = strings.TrimSpace(branch)
	if owner == "" || repo == "" {
		return "", fmt.Errorf("invalid repository reference")
	}
	if branch == "" {
		branch = "main"
	}

	safePath, err := sanitizeRepoPath(filePath)
	if err != nil {
		return "", err
	}

	base := &url.URL{Scheme: "https", Host: host}
	base.Path = path.Join(owner, repo, "blob", branch, safePath)
	return base.String(), nil
}

// buildVSCodeBlobURL creates a safe vscode.dev blob URL for a repository file.
func buildVSCodeBlobURL(owner, repo, branch, filePath string) (string, error) {
	owner = strings.TrimSpace(owner)
	repo = strings.TrimSpace(repo)
	branch = strings.TrimSpace(branch)
	if owner == "" || repo == "" {
		return "", fmt.Errorf("invalid repository reference")
	}
	if branch == "" {
		branch = "main"
	}

	safePath, err := sanitizeRepoPath(filePath)
	if err != nil {
		return "", err
	}

	base := &url.URL{Scheme: "https", Host: "vscode.dev"}
	base.Path = path.Join("github", owner, repo, "blob", branch, safePath)
	return base.String(), nil
}

// sanitizeRepoPath normalizes a repository path and validates control characters.
func sanitizeRepoPath(filePath string) (string, error) {
	cleaned := strings.TrimSpace(strings.TrimPrefix(filePath, "/"))
	if cleaned == "" {
		return "", fmt.Errorf("invalid file path")
	}

	parts := strings.Split(cleaned, "/")
	normalized := make([]string, 0, len(parts))
	for _, p := range parts {
		if p == "" {
			continue
		}

		for _, r := range p {
			if unicode.IsControl(r) {
				return "", fmt.Errorf("invalid file path")
			}
		}

		normalized = append(normalized, p)
	}

	if len(normalized) == 0 {
		return "", fmt.Errorf("invalid file path")
	}

	return strings.Join(normalized, "/"), nil
}

// openExternalURL opens a URL using the OS handler without invoking a shell.
func openExternalURL(targetURL string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", targetURL)
	case "darwin":
		cmd = exec.Command("open", targetURL)
	default:
		cmd = exec.Command("xdg-open", targetURL)
	}

	return cmd.Start()
}

// cloneToDesktop clones the repository to the Desktop folder
func (m FileEditModel) cloneToDesktop() tea.Cmd {
	return func() tea.Msg {
		desktopPath := getDesktopPath()
		clonePath := filepath.Join(desktopPath, m.repoName)

		// Check if already exists
		if _, err := os.Stat(clonePath); err == nil {
			return cloneResultMsg{fmt.Errorf("folder already exists: %s", clonePath)}
		}

		// Clone the repository
		repoURL := fmt.Sprintf("https://github.com/%s/%s.git", m.repoOwner, m.repoName)
		cmd := exec.Command("git", "clone", repoURL, clonePath)

		err := cmd.Run()
		if err != nil {
			return cloneResultMsg{err}
		}

		// Open file manager to show the cloned folder
		openFileManager(clonePath)

		return cloneResultMsg{nil}
	}
}

// openClonedInVSCode opens the cloned repository in a new VS Code window
func (m FileEditModel) openClonedInVSCode() tea.Cmd {
	return func() tea.Msg {
		// Open VS Code in a new window with the cloned repo
		// -n flag opens a new window
		filePath := filepath.Join(m.clonePath, strings.TrimPrefix(m.filePath, "/"))

		// First try to open VS Code with the specific file
		cmd := exec.Command("code", "-n", m.clonePath, "-g", filePath)
		err := cmd.Start()

		if err != nil {
			// Fallback: just open the folder
			cmd = exec.Command("code", "-n", m.clonePath)
			err = cmd.Start()
		}

		return openResultMsg{err}
	}
}
