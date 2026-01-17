package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/agnivo988/Repo-lyzer/internal/github"
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
	root         *FileNode
	cursor       int
	visibleList  []*FileNode
	width        int
	height       int
	Done         bool
	SelectedPath string
	searchInput  textinput.Model
	searchMode   bool
	searchQuery  string
	allNodes     []*FileNode // Store all nodes for search filtering
}

// NewTreeModel creates a new tree model for displaying the repository file structure.
// It builds the file tree from analysis results if provided, otherwise creates an empty tree.
// Parameters:
//   - result: Pointer to AnalysisResult containing file tree data, can be nil
// Returns the initialized TreeModel with the file tree populated.
func NewTreeModel(result *AnalysisResult) TreeModel {
	var root *FileNode
	if result != nil {
		root = BuildFileTree(*result)
	} else {
		root = &FileNode{
			Name:     "repository",
			Type:     "dir",
			Path:     "/",
			Children: []*FileNode{},
		}
	}

	// Initialize search input
	ti := textinput.New()
	ti.Placeholder = "Search files..."
	ti.CharLimit = 50
	ti.Width = 30

	m := TreeModel{
		root:        root,
		searchInput: ti,
		searchMode:  false,
	}
	m.updateAllNodes()
	m.updateVisibleList()
	return m
}

func (m *TreeModel) updateAllNodes() {
	m.allNodes = []*FileNode{}
	m.addAllNodes(m.root)
}

func (m *TreeModel) addAllNodes(node *FileNode) {
	m.allNodes = append(m.allNodes, node)
	for _, child := range node.Children {
		m.addAllNodes(child)
	}
}

func (m *TreeModel) updateVisibleList() {
	if m.searchMode && m.searchQuery != "" {
		m.updateFilteredList()
	} else {
		m.visibleList = []*FileNode{}
		m.addVisibleNodes(m.root, 0)
	}
}

func (m *TreeModel) updateFilteredList() {
	m.visibleList = []*FileNode{}
	query := strings.ToLower(m.searchQuery)
	
	for _, node := range m.allNodes {
		if strings.Contains(strings.ToLower(node.Name), query) {
			m.visibleList = append(m.visibleList, node)
		}
	}
}

func (m *TreeModel) addVisibleNodes(node *FileNode, depth int) {
	m.visibleList = append(m.visibleList, node)

	if node.Expanded && len(node.Children) > 0 {
		for _, child := range node.Children {
			m.addVisibleNodes(child, depth+1)
		}
	}
}

func (m TreeModel) Init() tea.Cmd { return nil }

func (m TreeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		if m.searchMode {
			switch msg.String() {
			case "esc":
				m.searchMode = false
				m.searchQuery = ""
				m.searchInput.Reset()
				m.updateVisibleList()
				m.cursor = 0
			case "enter":
				m.searchMode = false
			default:
				m.searchInput, cmd = m.searchInput.Update(msg)
				m.searchQuery = m.searchInput.Value()
				m.updateVisibleList()
				m.cursor = 0
			}
		} else {
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
					node := m.visibleList[m.cursor]
					if node.Type == "file" {
						m.SelectedPath = node.Path
						m.Done = true
					}
				}
			case "/":
				m.searchMode = true
				m.searchInput.Focus()
				cmd = m.searchInput.Cursor.BlinkCmd()
			case "esc":
				m.Done = true
			}
		}
	}

	return m, cmd
}

func (m TreeModel) View() string {
	if m.width == 0 || m.height == 0 {
		return "Initializing..."
	}

	content := TitleStyle.Render("📁 REPOSITORY FILE TREE") + "\n\n"

	// Show search input if in search mode
	if m.searchMode {
		searchView := m.searchInput.View()
		content += "🔍 " + searchView + "\n\n"
	} else {
		content += SubtleStyle.Render("Press / to search files") + "\n\n"
	}

	// Display visible nodes
	startIdx := m.cursor - (m.height-8)/2
	if startIdx < 0 {
		startIdx = 0
	}
	endIdx := startIdx + (m.height - 8)
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

	footer := SubtleStyle.Render("↑↓ navigate • ← → expand/collapse • / search • Enter edit file • ESC back")
	content += "\n" + footer

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Left, lipgloss.Top,
		BoxStyle.Render(content),
	)
}

func (m TreeModel) getIndent(node *FileNode) string {
	if m.searchMode && m.searchQuery != "" {
		// For search results, show the full path as context
		pathParts := strings.Split(strings.Trim(node.Path, "/"), "/")
		if len(pathParts) > 1 {
			parentPath := strings.Join(pathParts[:len(pathParts)-1], "/")
			return SubtleStyle.Render("└─ " + parentPath + "/")
		}
		return ""
	}
	
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
func BuildFileTree(result AnalysisResult) *FileNode {
	repoName := "repository"
	if result.Repo != nil {
		repoName = result.Repo.Name
	}
	root := &FileNode{
		Name:     repoName,
		Type:     "dir",
		Path:     "/",
		Children: []*FileNode{},
	}

	// Build tree from actual FileTree data
	for _, entry := range result.FileTree {
		addEntryToTree(root, entry)
	}

	return root
}

// addEntryToTree recursively adds a TreeEntry to the FileNode tree
func addEntryToTree(root *FileNode, entry github.TreeEntry) {
	parts := strings.Split(strings.Trim(entry.Path, "/"), "/")
	current := root

	for i, part := range parts {
		isLast := i == len(parts)-1
		found := false

		// Check if node already exists
		for _, child := range current.Children {
			if child.Name == part {
				current = child
				found = true
				break
			}
		}

		// Create new node if not found
		if !found {
			nodeType := "file"
			if !isLast || entry.Type == "tree" {
				nodeType = "dir"
			}
			newNode := &FileNode{
				Name:     part,
				Type:     nodeType,
				Path:     "/" + strings.Join(parts[:i+1], "/"),
				Children: []*FileNode{},
			}
			if isLast {
				newNode.Size = int64(entry.Size)
			}
			current.Children = append(current.Children, newNode)
			current = newNode
		}
	}
}
