package core

// Metric represents an individual scoring component
type Metric struct {
	Name        string
	Score       float64
	Weight      float64
	Description string
}

// ScoreCategory categorizes the final score
type ScoreCategory string

const (
	Critical  ScoreCategory = "Critical"
	Warning   ScoreCategory = "Warning"
	Healthy   ScoreCategory = "Healthy"
	Excellent ScoreCategory = "Excellent"
)

// Thresholds defines the bounds for each category
type Thresholds struct {
	Warning   float64
	Healthy   float64
	Excellent float64
}

// Engine defines the generic scoring engine
type Engine interface {
	CalculateScore(metrics []Metric) (float64, ScoreCategory)
	GetMaxScore() float64
}

// WeightedScoreEngine implements a reusable weighted scoring logic
type WeightedScoreEngine struct {
	MaxScore   float64
	Thresholds Thresholds
}

// NewWeightedScoreEngine initializes a new engine
func NewWeightedScoreEngine(maxScore float64, thresholds Thresholds) *WeightedScoreEngine {
	return &WeightedScoreEngine{
		MaxScore:   maxScore,
		Thresholds: thresholds,
	}
}

// CalculateScore computes the weighted sum and determines the category
func (e *WeightedScoreEngine) CalculateScore(metrics []Metric) (float64, ScoreCategory) {
	var totalScore float64
	var totalWeight float64

	for _, m := range metrics {
		totalScore += m.Score * m.Weight
		totalWeight += m.Weight
	}

	if totalWeight == 0 {
		return 0, Critical
	}

	finalScore := totalScore / totalWeight
	if finalScore > e.MaxScore {
		finalScore = e.MaxScore
	}
	if finalScore < 0 {
		finalScore = 0
	}

	return finalScore, e.GetCategory(finalScore)
}

// GetCategory returns the category string based on the thresholds
func (e *WeightedScoreEngine) GetCategory(score float64) ScoreCategory {
	if score >= e.Thresholds.Excellent {
		return Excellent
	} else if score >= e.Thresholds.Healthy {
		return Healthy
	} else if score >= e.Thresholds.Warning {
		return Warning
	}
	return Critical
}

// GetMaxScore returns the max score configured for the engine
func (e *WeightedScoreEngine) GetMaxScore() float64 {
	return e.MaxScore
}
