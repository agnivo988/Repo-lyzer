package ui

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/agnivo988/Repo-lyzer/internal/predictive"
	"github.com/agnivo988/Repo-lyzer/internal/temporal"
	"github.com/charmbracelet/lipgloss"
)

type monthlyTrendPoint struct {
	Month            time.Time
	CommitCount      int
	ContributorCount int
	HealthScore      int
}

func (m DashboardModel) repositoryTrendsView() string {
	header := TitleStyle.Render(" Repository Trends ")
	series := m.buildMonthlyTrendSeries(4)

	if len(series) == 0 {
		return lipgloss.JoinVertical(
			lipgloss.Left,
			header,
			CardStyle.Render("No historical commit data available for trend analysis"),
		)
	}

	healthLines := make([]string, 0, len(series))
	contributorLines := make([]string, 0, len(series))
	for _, point := range series {
		healthLines = append(healthLines, fmt.Sprintf("%s: %d", point.Month.Format("Jan"), point.HealthScore))
		contributorLines = append(contributorLines, fmt.Sprintf("%s: %d", point.Month.Format("Jan"), point.ContributorCount))
	}

	healthTrend := trendLabelFromFirstLast(series, func(point monthlyTrendPoint) int { return point.HealthScore })
	contributorTrend := trendLabelFromFirstLast(series, func(point monthlyTrendPoint) int { return point.ContributorCount })

	forecastLines := m.forecastTrendCardLines(series)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		CardStyle.Render(lipgloss.JoinVertical(lipgloss.Left,
			lipgloss.NewStyle().Bold(true).Render("Health Score"),
			"",
			strings.Join(healthLines, "\n"),
			"",
			"Trend: "+healthTrend,
		)),
		CardStyle.Render(lipgloss.JoinVertical(lipgloss.Left,
			lipgloss.NewStyle().Bold(true).Render("Contributor Growth"),
			"",
			strings.Join(contributorLines, "\n"),
			"",
			"Trend: "+contributorTrend,
		)),
		CardStyle.Render(lipgloss.JoinVertical(lipgloss.Left,
			lipgloss.NewStyle().Bold(true).Render("Predicted Health Score"),
			"",
			strings.Join(forecastLines, "\n"),
		)),
	)
}

func (m DashboardModel) repositoryTrendsSummaryCard(series []monthlyTrendPoint) string {
	if len(series) == 0 {
		return ""
	}

	healthTrend := trendLabelFromFirstLast(series, func(point monthlyTrendPoint) int { return point.HealthScore })
	contributorTrend := trendLabelFromFirstLast(series, func(point monthlyTrendPoint) int { return point.ContributorCount })
	forecastLines := m.forecastTrendCardLines(series)
	forecastSummary := "Forecast unavailable"
	if len(forecastLines) >= 2 {
		forecastSummary = fmt.Sprintf("%s / %s", forecastLines[0], forecastLines[1])
	} else if len(forecastLines) == 1 {
		forecastSummary = forecastLines[0]
	}

	lines := []string{
		lipgloss.NewStyle().Bold(true).Render("Repository Trends"),
		"",
		fmt.Sprintf("Health: %s", healthTrend),
		fmt.Sprintf("Contributors: %s", contributorTrend),
		fmt.Sprintf("Forecast: %s", forecastSummary),
	}

	return CardStyle.Render(strings.Join(lines, "\n"))
}

func (m DashboardModel) buildMonthlyTrendSeries(limit int) []monthlyTrendPoint {
	if limit <= 0 {
		limit = 4
	}

	if len(m.data.Commits) == 0 {
		if m.data.HealthScore == 0 && len(m.data.Contributors) == 0 {
			return nil
		}
		return []monthlyTrendPoint{{
			Month:            time.Now().UTC(),
			CommitCount:      len(m.data.Commits),
			ContributorCount: len(m.data.Contributors),
			HealthScore:      m.data.HealthScore,
		}}
	}

	type monthBucket struct {
		month        time.Time
		commitCount  int
		contributors map[string]struct{}
	}

	buckets := make(map[string]*monthBucket)
	var earliest time.Time
	var latest time.Time
	for _, commit := range m.data.Commits {
		commitTime := commit.Commit.Author.Date
		if commitTime.IsZero() {
			continue
		}

		month := time.Date(commitTime.Year(), commitTime.Month(), 1, 0, 0, 0, 0, time.UTC)
		key := month.Format("2006-01")
		bucket, ok := buckets[key]
		if !ok {
			bucket = &monthBucket{month: month, contributors: make(map[string]struct{})}
			buckets[key] = bucket
		}
		bucket.commitCount++

		contributorKey := ""
		if commit.Author != nil && strings.TrimSpace(commit.Author.Login) != "" {
			contributorKey = commit.Author.Login
		}
		bucket.contributors[contributorKey] = struct{}{}

		if earliest.IsZero() || month.Before(earliest) {
			earliest = month
		}
		if latest.IsZero() || month.After(latest) {
			latest = month
		}
	}

	if earliest.IsZero() || latest.IsZero() {
		return nil
	}

	months := make([]time.Time, 0)
	for cursor := earliest; !cursor.After(latest); cursor = cursor.AddDate(0, 1, 0) {
		months = append(months, cursor)
	}

	if len(months) > limit {
		months = months[len(months)-limit:]
	}

	series := make([]monthlyTrendPoint, 0, len(months))
	for _, month := range months {
		bucket := buckets[month.Format("2006-01")]
		commitCount := 0
		contributorCount := 0
		if bucket != nil {
			commitCount = bucket.commitCount
			contributorCount = len(bucket.contributors)
		}
		series = append(series, monthlyTrendPoint{
			Month:            month,
			CommitCount:      commitCount,
			ContributorCount: contributorCount,
		})
	}

	maxCommits := 0
	maxContributors := 0
	for _, point := range series {
		if point.CommitCount > maxCommits {
			maxCommits = point.CommitCount
		}
		if point.ContributorCount > maxContributors {
			maxContributors = point.ContributorCount
		}
	}

	for i := range series {
		series[i].HealthScore = calculateMonthlyHealth(series[i].CommitCount, series[i].ContributorCount, maxCommits, maxContributors)
	}

	return series
}

func calculateMonthlyHealth(commitCount, contributorCount, maxCommits, maxContributors int) int {
	if maxCommits == 0 && maxContributors == 0 {
		return 0
	}

	score := 50
	if maxCommits > 0 {
		score += int(math.Round(25 * float64(commitCount) / float64(maxCommits)))
	}
	if maxContributors > 0 {
		score += int(math.Round(25 * float64(contributorCount) / float64(maxContributors)))
	}

	if score > 100 {
		score = 100
	}
	return score
}

func trendLabelFromFirstLast[T any](series []T, value func(T) int) string {
	if len(series) < 2 {
		return "Stable"
	}

	first := value(series[0])
	last := value(series[len(series)-1])

	switch {
	case last > first:
		return "Increasing"
	case last < first:
		return "Decreasing"
	default:
		return "Stable"
	}
}

func (m DashboardModel) forecastTrend(series []monthlyTrendPoint) *predictive.ForecastResult {
	if len(series) < 2 {
		return nil
	}

	timeline := temporal.NewTimeline("", "")
	for _, point := range series {
		snapshot := temporal.NewSnapshot(point.Month, nil)
		snapshot.Metrics.AverageHealth = point.HealthScore
		snapshot.Metrics.ContributorCount = point.ContributorCount
		_ = timeline.AddSnapshot(snapshot)
	}

	predictor := predictive.NewPredictor()
	predictor.ForecastHorizon = 3
	predictor.ConfidenceLevel = 0.95

	forecast, err := predictive.ForecastHealthFromTimeline(predictor, timeline, 3)
	if err != nil {
		return nil
	}
	return forecast
}

func formatForecastLinesFromForecast(forecast *predictive.ForecastResult) []string {
	if forecast == nil || len(forecast.Predictions) == 0 {
		return nil
	}

	lines := []string{fmt.Sprintf("30 Days: %.0f", forecast.Predictions[0].Value)}
	if len(forecast.Predictions) >= 3 {
		lines = append(lines, fmt.Sprintf("90 Days: %.0f", forecast.Predictions[2].Value))
	} else {
		lines = append(lines, fmt.Sprintf("90 Days: %.0f", forecast.Predictions[len(forecast.Predictions)-1].Value))
	}
	return lines
}

func (m DashboardModel) forecastTrendCardLines(series []monthlyTrendPoint) []string {
	forecast := m.forecastTrend(series)
	if forecast == nil {
		return []string{"Forecast unavailable"}
	}

	lines := formatForecastLinesFromForecast(forecast)
	if len(lines) == 0 {
		return []string{"Forecast unavailable"}
	}
	return lines
}
