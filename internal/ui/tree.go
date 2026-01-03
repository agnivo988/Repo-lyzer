package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// FileNode represents a file or directory in the repository
type FileNode struct {
	Name     string
	Type     string // "file" or "dir"
	Path     string
	Size     int64
	Children []*FileNode
	Expanded bool
}

// TreeModel represents the file tree view
type TreeModel struct {
	root        *FileNode
	cursor      int
	visibleList []*FileNode
	width       int
	height      int
	Done        bool
	SelectedPath string
}

func NewTreeModel(root *FileNode) TreeModel {
	if root == nil {
		root = &FileNode{
			Name:     "repository",
			Type:     "dir",
			Path:     "/",
			Children: []*FileNode{},
		}
	}

	m := TreeModel{
		root: root,
	}
	m.updateVisibleList()
	return m
}

func (m *TreeModel) updateVisibleList() {
	m.visibleList = []*FileNode{}
	m.addVisibleNodes(m.root, 0)
}

func (m *TreeModel) addVisibleNodes(node *FileNode, depth int) {
	if node == m.root {
		m.visibleList = append(m.visibleList, node)
	} else {
		m.visibleList = append(m.visibleList, node)
	}

	if node.Expanded && len(node.Children) > 0 {
		for _, child := range node.Children {
			m.addVisibleNodes(child, depth+1)
		}
	}
}

func (m TreeModel) Init() tea.Cmd { return nil }

func (m TreeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.visibleList)-1 {
				m.cursor++
			}
		case "right", "l":
			if m.cursor < len(m.visibleList) {
				node := m.visibleList[m.cursor]
				if node.Type == "dir" && len(node.Children) > 0 {
					node.Expanded = true
					m.updateVisibleList()
				}
			}
		case "left", "h":
			if m.cursor < len(m.visibleList) {
				node := m.visibleList[m.cursor]
				if node.Type == "dir" && node.Expanded {
					node.Expanded = false
					m.updateVisibleList()
				}
			}
		case "enter":
			if m.cursor < len(m.visibleList) {
				m.SelectedPath = m.visibleList[m.cursor].Path
				m.Done = true
			}
		case "esc":
			m.Done = true
		}
	}

	return m, nil
}

func (m TreeModel) View() string {
	if m.width == 0 || m.height == 0 {
		return "Initializing..."
	}

	content := TitleStyle.Render("📁 REPOSITORY FILE TREE") + "\n\n"

	// Display visible nodes
	startIdx := m.cursor - (m.height - 5) / 2
	if startIdx < 0 {
		startIdx = 0
	}
	endIdx := startIdx + (m.height - 5)
	if endIdx > len(m.visibleList) {
		endIdx = len(m.visibleList)
	}

	for i := startIdx; i < endIdx; i++ {
		node := m.visibleList[i]
		indent := m.getIndent(node)

		icon := "📄"
		if node.Type == "dir" {
			icon = "📁"
			if node.Expanded && len(node.Children) > 0 {
				icon = "📂"
			}
		}

		prefix := "  "
		style := NormalStyle
		if i == m.cursor {
			prefix = "▶ "
			style = SelectedStyle
		}

		line := fmt.Sprintf("%s%s%s %s", prefix, indent, icon, node.Name)
		content += style.Render(line) + "\n"
	}

	footer := SubtleStyle.Render("↑↓ navigate • ← → expand/collapse • Enter select • ESC back")
	content += "\n" + footer

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Left, lipgloss.Top,
		BoxStyle.Render(content),
	)
}

func (m TreeModel) getIndent(node *FileNode) string {
	depth := m.getNodeDepth(m.root, node)
	indent := ""
	for i := 0; i < depth; i++ {
		indent += "  "
	}
	return indent
}

func (m TreeModel) getNodeDepth(parent *FileNode, target *FileNode) int {
	if parent == target {
		return 0
	}

	for _, child := range parent.Children {
		if child == target {
			return 1
		}
		depth := m.getNodeDepth(child, target)
		if depth >= 0 {
			return depth + 1
		}
	}
	return -1
}

// BuildFileTree creates a file tree from repository content
func BuildFileTree(fileCount int, topFiles []string) *FileNode {
	root := &FileNode{
		Name:     "repository",
		Type:     "dir",
		Path:     "/",
		Children: []*FileNode{},
	}

	// Add directories
	dirs := []string{"src", "test", "docs", "config", "scripts", "build"}
	for _, dir := range dirs {
		root.Children = append(root.Children, &FileNode{
			Name:     dir,
			Type:     "dir",
			Path:     "/" + dir,
			Children: []*FileNode{},
		})
	}

	// Add sample files in root
	sampleFiles := []string{"README.md", "LICENSE", "go.mod", "go.sum", ".gitignore"}
	for _, file := range sampleFiles {
		root.Children = append(root.Children, &FileNode{
			Name: file,
			Type: "file",
			Path: "/" + file,
			Size: 1024,
		})
	}

	return root
}
