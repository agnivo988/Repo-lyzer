package risk

import (
	"strings"

	"github.com/agnivo988/Repo-lyzer/internal/analyzer/core"
	"github.com/agnivo988/Repo-lyzer/internal/github"
)

// SecurityRiskAnalyzer assesses repository security hygiene and posture
type SecurityRiskAnalyzer struct {
	Engine *core.WeightedScoreEngine
}

// NewSecurityRiskAnalyzer initializes the security risk analyzer
func NewSecurityRiskAnalyzer() *SecurityRiskAnalyzer {
	thresholds := core.Thresholds{
		Warning:   40,
		Healthy:   0,
		Excellent: 80,
	}
	return &SecurityRiskAnalyzer{
		Engine: core.NewWeightedScoreEngine(100.0, thresholds),
	}
}

// AnalyzeSecurityPosture evaluates missing security files and basic CI pipelines
func (s *SecurityRiskAnalyzer) AnalyzeSecurityPosture(tree []github.TreeEntry) (float64, RiskCategory) {
	metrics := []core.Metric{}

	hasSecurityMD := false
	hasCIWorkflows := false
	hasIssueTemplates := false

	for _, entry := range tree {
		path := strings.ToLower(entry.Path)

		// Check for SECURITY.md
		if path == "security.md" || path == ".github/security.md" || path == "docs/security.md" {
			hasSecurityMD = true
		}

		// Check for CI workflows
		if strings.HasPrefix(path, ".github/workflows/") && strings.HasSuffix(path, ".yml") {
			hasCIWorkflows = true
		}

		// Check for Issue Templates
		if strings.HasPrefix(path, ".github/issue_template") {
			hasIssueTemplates = true
		}
	}

	// 1. Missing SECURITY.md Risk
	securityMDRisk := 100.0
	if hasSecurityMD {
		securityMDRisk = 0.0
	}
	metrics = append(metrics, core.Metric{
		Name:        "Missing SECURITY.md",
		Score:       securityMDRisk,
		Weight:      3.0,
		Description: "High risk if repository has no clear vulnerability reporting guidelines",
	})

	// 2. Missing CI Workflows Risk
	ciRisk := 100.0
	if hasCIWorkflows {
		ciRisk = 0.0
	}
	metrics = append(metrics, core.Metric{
		Name:        "Missing CI Workflows",
		Score:       ciRisk,
		Weight:      2.0,
		Description: "Risk of unverified code merging if no Continuous Integration is present",
	})

	// 3. Missing Issue Templates Risk
	templateRisk := 50.0 // Warning level risk
	if hasIssueTemplates {
		templateRisk = 0.0
	}
	metrics = append(metrics, core.Metric{
		Name:        "Missing Issue Templates",
		Score:       templateRisk,
		Weight:      1.0,
		Description: "Increases risk of poorly formatted bug reports",
	})

	score, _ := s.Engine.CalculateScore(metrics)

	var category RiskCategory
	if score >= 80.0 {
		category = RiskCritical
	} else if score >= 40.0 {
		category = RiskWarning
	} else {
		category = RiskStable
	}

	return score, category
}
