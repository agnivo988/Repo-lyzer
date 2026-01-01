package cmd

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/term"
)

var logo = `
 ██████╗ ███████╗██████╗  ██████╗      ██╗     ██╗   ██╗ ███████╗  ███████╗██████╗ 
 ██╔══██╗██╔════╝██╔══██╗██╔═══██╗     ██║     ╚██╗ ██╔╝ ╚════██║  ██╔════╝██╔══██╗
 ██████╔╝█████╗  ██████╔╝██║   ██║█████╗██║      ╚████╔╝     ██╔╝  █████╗  ██████╔╝
 ██╔══██╗██╔══╝  ██╔═══╝ ██║   ██║╚════╝██║       ╚██╔╝     ██╔╝   ██╔══╝  ██╔══██╗
 ██║  ██║███████╗██║     ╚██████╔╝     ███████╗   ██║      ██╔╝     ███████╗██║  ██║   
 ╚═╝  ╚═╝╚══════╝╚═╝      ╚═════╝      ╚══════╝   ╚═╝     ███████╗ ╚══════╝╚═╝  ╚═╝     
`


var (
	titleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#00E5FF"))

	boxStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#7D56F4")).
		Padding(1, 4)

	selectedStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00FF87")).
		Bold(true)

	normalStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF"))

	inputStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFD700")).
		Bold(true)
)


type menuModel struct {
	step     int
	cursor   int
	choices  []string
	input    string
	selected int
}

func initialMenu() menuModel {
	return menuModel{
		choices: []string{
			"Analyze a repository",
			"Compare repositories",
			"Exit",
		},
	}
}

func (m menuModel) Init() tea.Cmd { return nil }

func (m menuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c", "q":
			os.Exit(0)

		// Navigation
		case "up":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		// Selection
		case "enter":
	if m.step == 0 {
		// EXIT
		if m.cursor == 2 {
			os.Exit(0)
		}

		m.selected = m.cursor
		m.step = 1
		m.input = ""
		return m, nil
	}

	if m.step == 1 {
		m.step = 2 // move to execution phase
		return m, tea.Quit
	}


		// Input typing
		default:
			if m.step == 1 {
				if msg.String() == "backspace" {
					if len(m.input) > 0 {
						m.input = m.input[:len(m.input)-1]
					}
				} else {
					m.input += msg.String()
				}
			}
		}
	}

	return m, nil
}

func (m menuModel) View() string {
	// Clear screen every render
	clear := "\033[H\033[2J"

	// ---------- MENU SCREEN ----------
	if m.step == 0 {
		content := titleStyle.Render(logo) + "\n\n"

		for i, choice := range m.choices {
			cursor := "  "
			style := normalStyle

			if m.cursor == i {
				cursor = "▶ "
				style = selectedStyle
			}

			content += fmt.Sprintf("%s%s\n", cursor, style.Render(choice))
		}

		content += "\n" + lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888")).
			Render("↑ ↓ navigate • Enter select • q quit")

		box := boxStyle.Render(content)

		w, h, err := term.GetSize(os.Stdout.Fd())
           if err != nil {
	       w, h = 80, 24 // fallback
               } 

              return clear + lipgloss.Place(
	          w,
	         h,
	       lipgloss.Center,
	      lipgloss.Center,
	     box,
         )

	}

	// ---------- INPUT SCREEN ----------
	inputContent :=
		titleStyle.Render("📥 ENTER REPOSITORY") + "\n\n" +
			inputStyle.Render("> "+m.input) + "\n\n" +
			lipgloss.NewStyle().
				Foreground(lipgloss.Color("#888888")).
				Render("Format: owner/repo  •  Press Enter to run")

	box := boxStyle.Render(inputContent)

	w, h, err := term.GetSize(os.Stdout.Fd())
if err != nil {
	w, h = 80, 24 // fallback
}

return clear + lipgloss.Place(
	w,
	h,
	lipgloss.Center,
	lipgloss.Center,
	box,
)

}


func RunMenu() {
	for {
		p := tea.NewProgram(initialMenu())
		model, err := p.Run()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		m := model.(menuModel)

		switch m.selected {

		case 0: // Analyze
			rootCmd.SetArgs([]string{"analyze", m.input})
			rootCmd.Execute()

		case 1: // Compare
			parts := strings.Split(m.input, " ")
			if len(parts) != 2 {
				fmt.Println("Error!! Use: owner1/repo1 owner2/repo2")
				continue
			}
			rootCmd.SetArgs([]string{"compare", parts[0], parts[1]})
			fmt.Print("\033[H\033[2J")
			rootCmd.Execute()
		}

		// Pause before returning to menu
		fmt.Println("\nPress Enter to return to menu...")
		fmt.Scanln()
	}
}

