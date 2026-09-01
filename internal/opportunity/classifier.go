package opportunity

import "github.com/google/go-github/v55/github"

type IssueClassification string

const (
    ClassificationGoodFirstIssue IssueClassification = "good-first-issue"
    ClassificationHelpWanted    IssueClassification = "help-wanted"
    ClassificationBug          IssueClassification = "bug"
    ClassificationEnhancement  IssueClassification = "enhancement"
    ClassificationDocumentation IssueClassification = "documentation"
    ClassificationOther        IssueClassification = "other"
)

// Classify returns a classification based on issue labels.
func Classify(issue *github.Issue) IssueClassification {
    if issue == nil {
        return ClassificationOther
    }
    for _, lbl := range issue.Labels {
        if lbl == nil || lbl.Name == nil {
            continue
        }
        name := *lbl.Name
        switch name {
        case "good first issue", "good-first-issue":
            return ClassificationGoodFirstIssue
        case "help wanted", "help-wanted":
            return ClassificationHelpWanted
        case "bug":
            return ClassificationBug
        case "enhancement", "feature":
            return ClassificationEnhancement
        case "documentation", "docs":
            return ClassificationDocumentation
        }
    }
    return ClassificationOther
}
