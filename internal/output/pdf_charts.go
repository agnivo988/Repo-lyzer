package output

import (
	"bytes"
	"fmt"
	"sort"
	"time"

	"github.com/agnivo988/Repo-lyzer/internal/github"
	"github.com/wcharczuk/go-chart/v2"
	"github.com/wcharczuk/go-chart/v2/drawing"
)

// GenerateCommitActivityChart creates a line chart of commit activity over time
func GenerateCommitActivityChart(commits []github.Commit) ([]byte, error) {
	if len(commits) == 0 {
		return nil, fmt.Errorf("no commits to chart")
	}

	// Group commits by date
	commitsByDate := make(map[string]int)
	for _, commit := range commits {
		dateStr := commit.Commit.Author.Date.Format("2006-01-02")
		commitsByDate[dateStr]++
	}

	// Sort dates
	dates := make([]string, 0, len(commitsByDate))
	for date := range commitsByDate {
		dates = append(dates, date)
	}
	sort.Strings(dates)

	// Prepare chart data
	xValues := make([]time.Time, len(dates))
	yValues := make([]float64, len(dates))

	for i, dateStr := range dates {
		date, _ := time.Parse("2006-01-02", dateStr)
		xValues[i] = date
		yValues[i] = float64(commitsByDate[dateStr])
	}

	// Create chart
	graph := chart.Chart{
		Width:  800,
		Height: 400,
		Background: chart.Style{
			Padding: chart.Box{Top: 20, Left: 50, Right: 20, Bottom: 40},
		},
		XAxis: chart.XAxis{
			Name:           "Date",
			ValueFormatter: chart.TimeValueFormatterWithFormat("Jan 02"),
		},
		YAxis: chart.YAxis{
			Name: "Commits",
		},
		Series: []chart.Series{
			chart.TimeSeries{
				XValues: xValues,
				YValues: yValues,
				Style: chart.Style{
					StrokeColor: drawing.ColorFromHex("1E40AF"),
					FillColor:   drawing.ColorFromHex("1E40AF").WithAlpha(64),
				},
			},
		},
	}

	buffer := bytes.NewBuffer([]byte{})
	err := graph.Render(chart.PNG, buffer)
	return buffer.Bytes(), err
}

// GenerateLanguageDistributionChart creates a pie chart of language distribution
func GenerateLanguageDistributionChart(languages map[string]int) ([]byte, error) {
	if len(languages) == 0 {
		return nil, fmt.Errorf("no languages to chart")
	}

	// Calculate total
	total := 0
	for _, bytes := range languages {
		total += bytes
	}

	// Convert to chart values
	var values []chart.Value
	for lang, bytes := range languages {
		pct := float64(bytes) / float64(total) * 100
		if pct >= 1.0 { // Only show languages >= 1%
			values = append(values, chart.Value{
				Label: fmt.Sprintf("%s (%.1f%%)", lang, pct),
				Value: float64(bytes),
			})
		}
	}

	// Sort by value descending
	sort.Slice(values, func(i, j int) bool {
		return values[i].Value > values[j].Value
	})

	pie := chart.PieChart{
		Width:  800,
		Height: 400,
		Values: values,
	}

	buffer := bytes.NewBuffer([]byte{})
	err := pie.Render(chart.PNG, buffer)
	return buffer.Bytes(), err
}

// GenerateContributorBarChart creates a bar chart of top contributors
func GenerateContributorBarChart(contributors []github.Contributor) ([]byte, error) {
	if len(contributors) == 0 {
		return nil, fmt.Errorf("no contributors to chart")
	}

	// Take top 10 contributors
	maxContribs := 10
	if len(contributors) < maxContribs {
		maxContribs = len(contributors)
	}

	bars := make([]chart.Value, maxContribs)
	for i := 0; i < maxContribs; i++ {
		bars[i] = chart.Value{
			Label: contributors[i].Login,
			Value: float64(contributors[i].Commits),
		}
	}

	graph := chart.BarChart{
		Width:  800,
		Height: 400,
		Background: chart.Style{
			Padding: chart.Box{Top: 20, Left: 50, Right: 20, Bottom: 40},
		},
		YAxis: chart.YAxis{
			Name: "Commits",
		},
		Bars: bars,
	}

	buffer := bytes.NewBuffer([]byte{})
	err := graph.Render(chart.PNG, buffer)
	return buffer.Bytes(), err
}

// GenerateHealthScoreGauge creates a visual gauge for health score
func GenerateHealthScoreGauge(score int) ([]byte, error) {
	// Create a simple bar representation of the health score
	bars := []chart.Value{
		{Label: fmt.Sprintf("Health: %d/100", score), Value: float64(score)},
	}

	graph := chart.BarChart{
		Width:  600,
		Height: 200,
		Background: chart.Style{
			Padding: chart.Box{Top: 20, Left: 50, Right: 20, Bottom: 40},
		},
		YAxis: chart.YAxis{
			Name: "Score",
			Range: &chart.ContinuousRange{
				Min: 0,
				Max: 100,
			},
		},
		Bars:     bars,
		BarWidth: 100,
	}

	buffer := bytes.NewBuffer([]byte{})
	err := graph.Render(chart.PNG, buffer)
	return buffer.Bytes(), err
}
