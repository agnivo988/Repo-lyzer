package output

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func PrintLanguages(langs map[string]int) {

	fmt.Println(SectionStyle.Render("\n⛳ Language Breakdown"))

	if len(langs) == 0 {
		fmt.Println("No language data available.")
		return
	}

	total := 0
	for _, v := range langs {
		total += v
	}

	for lang, size := range langs {
		percent := float64(size) / float64(total) * 100
		bar := lipgloss.NewStyle().Foreground(lipgloss.Color("#7CFF00")).Render(strings.Repeat("🟩", int(percent/5)))

		fmt.Printf("%-10s %s %.1f%%\n", lang, bar, percent)
	}
}

type AnalysisOutput struct {
	Repo          interface{}    `json:"Repo"`
	Commits       interface{}    `json:"Commits"`
	Contributors  interface{}    `json:"Contributors"`
	FileTree      interface{}    `json:"FileTree"`
	Languages     map[string]int `json:"Languages"`
	HealthScore   int            `json:"HealthScore"`
	BusFactor     int            `json:"BusFactor"`
	BusRisk       string         `json:"BusRisk"`
	MaturityScore int            `json:"MaturityScore"`
	MaturityLevel string         `json:"MaturityLevel"`
}

func SaveAnalysisJSON(
	path string,
	repo interface{},
	commits interface{},
	contributors interface{},
	fileTree interface{},
	langs map[string]int,
	healthScore int,
	busFactor int,
	busRisk string,
	maturityScore int,
	maturityLevel string,
) error {

	data := AnalysisOutput{
		Repo:          repo,
		Commits:       commits,
		Contributors:  contributors,
		FileTree:      fileTree,
		Languages:     langs,
		HealthScore:   healthScore,
		BusFactor:     busFactor,
		BusRisk:       busRisk,
		MaturityScore: maturityScore,
		MaturityLevel: maturityLevel,
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(path, jsonData, 0644)
	if err != nil {
		return err
	}

	return nil
}
