package analyzer

import (
	"math"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/agnivo988/Repo-lyzer/internal/github"
)

// Hotspot represents a problematic file
type Hotspot struct {
	FilePath   string `json:"file_path"`
	Score      int    `json:"score"`       // 0-100
	ChurnScore int    `json:"churn_score"` // Based on commit frequency
	SizeScore  int    `json:"size_score"`  // Based on LOC/Size
	Complexity int    `json:"complexity"`  // Cyclomatic complexity
	Reason     string `json:"reason"`
}

// AnalyzeHotspots identifies the most problematic files in the repository
func AnalyzeHotspots(
	repo *github.Repo,
	commits []github.Commit,
	fileTree []github.TreeEntry,
	client *github.Client,
) ([]Hotspot, error) {

	// 1. Calculate Churn (File Change Frequency)
	// We need to fetch details for recent commits to see which files changed.
	// To avoid rate limits, we analyze the last 30 commits.
	churnMap := analyzeChurn(commits, repo, client)

	// 2. Analyze Size and Nesting Depth
	// 3. Combine metrics to identify potential hotspots
	var candidates []Hotspot

	for _, entry := range fileTree {
		if entry.Type != "blob" || !isSourceFile(entry.Path) {
			continue
		}

		size := entry.Size
		churn := churnMap[entry.Path]

		// Skip small, stable files
		if size < 1000 && churn == 0 {
			continue
		}

		// Calculate Preliminary Score
		// Size Score: Logarithmic scale (1KB ~ 10pts, 100KB ~ 100pts)
		sizeScore := int(math.Log10(float64(size)) * 20)
		if sizeScore > 100 {
			sizeScore = 100
		}

		// Churn Score: Linear (1 change ~ 10pts, 10 changes ~ 100pts)
		churnScore := churn * 10
		if churnScore > 100 {
			churnScore = 100
		}

		// Nesting Score
		nestingDepth := strings.Count(entry.Path, "/")
		nestingScore := nestingDepth * 5

		totalScore := (sizeScore * 30 / 100) + (churnScore * 50 / 100) + (nestingScore)

		candidates = append(candidates, Hotspot{
			FilePath:   entry.Path,
			Score:      totalScore,
			ChurnScore: churnScore,
			SizeScore:  sizeScore,
			Complexity: 0, // Calculated later for top candidates
		})
	}

	// Sort by score
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].Score > candidates[j].Score
	})

	// Top candidates (max 10)
	limit := 10
	if len(candidates) < limit {
		limit = len(candidates)
	}
	topCandidates := candidates[:limit]

	// 4. Calculate Cyclomatic Complexity for top candidates
	// Fetch content for these files
	var wg sync.WaitGroup
	for i := range topCandidates {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			hotspot := &topCandidates[idx]

			content, err := client.GetFileContent(repo.Owner.Login, repo.Name, hotspot.FilePath)
			if err == nil {
				hotspot.Complexity = calculateComplexity(content, hotspot.FilePath)

				// Update score with actual complexity
				// Map complexity 1-50 to 0-100
				compScore := hotspot.Complexity * 2
				if compScore > 100 {
					compScore = 100
				}

				// Refine total score
				// Churn: 40%, Size: 20%, Complexity: 40%
				hotspot.Score = (hotspot.ChurnScore * 40 / 100) +
					(hotspot.SizeScore * 20 / 100) +
					(compScore * 40 / 100)

				hotspot.Reason = generateReason(hotspot)
			}
		}(i)
	}
	wg.Wait()

	// Re-sort after complexity update
	sort.Slice(topCandidates, func(i, j int) bool {
		return topCandidates[i].Score > topCandidates[j].Score
	})

	return topCandidates, nil
}

func analyzeChurn(commits []github.Commit, repo *github.Repo, client *github.Client) map[string]int {
	churnMap := make(map[string]int)
	limit := 30
	if len(commits) < limit {
		limit = len(commits)
	}

	// Use a worker pool or simple concurrency to fetch commit details
	// Since we are CLI, we can do it with a semaphore to avoid 429s too fast,
	// but purely sequential is safer for public API. parallel with limit is best.

	// For now, let's do sequential to be safe with rate limits as per plan warning
	// Or semi-parallel with slight delay.
	// Actually, the plan said "analyze only last 30 commits".

	recentCommits := commits[:limit]
	var mu sync.Mutex
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 5) // Concurrency limit of 5

	for _, c := range recentCommits {
		wg.Add(1)
		go func(sha string) {
			defer wg.Done()
			semaphore <- struct{}{}        // Acquire
			defer func() { <-semaphore }() // Release

			detail, err := client.GetCommit(repo.Owner.Login, repo.Name, sha)
			if err != nil {
				// Ignore error for individual commit fetch failures
				return
			}

			mu.Lock()
			for _, f := range detail.Files {
				churnMap[f.Filename]++
			}
			mu.Unlock()

			time.Sleep(100 * time.Millisecond) // Slight delay to be nice
		}(c.SHA)
	}
	wg.Wait()

	return churnMap
}

func calculateComplexity(content, filename string) int {
	// Rudimentary complexity calculation
	// Count keywords: if, for, case, catch, while, &&, ||, ternary ?

	complexity := 1 // Base complexity

	keywords := []string{
		"if ", "if(",
		"for ", "for(", "foreach",
		"case ",
		"catch",
		"while ", "while(",
		"&&", "||",
		" ? ", " ?? ",
	}

	lines := strings.Split(content, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "//") || strings.HasPrefix(trimmed, "*") {
			continue // Skip comments (basic check)
		}

		for _, kw := range keywords {
			if strings.Contains(trimmed, kw) {
				complexity++
			}
		}
	}

	return complexity
}

func generateReason(h *Hotspot) string {
	var reasons []string
	if h.ChurnScore > 70 {
		reasons = append(reasons, "Frequently Changed")
	}
	if h.Complexity > 20 {
		reasons = append(reasons, "High Complexity")
	}
	if h.SizeScore > 80 {
		reasons = append(reasons, "Very Large File")
	}

	if len(reasons) == 0 {
		return "Moderate Complexity"
	}
	return strings.Join(reasons, ", ")
}
