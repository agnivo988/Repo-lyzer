package output

import (
	"fmt"
	"os"

	"Repo-lyzer/internal/github"

	"github.com/olekukonko/tablewriter"
)

func PrintRepo(r *github.Repo) {
	table := tablewriter.NewWriter(os.Stdout)
	table.Header([]string{"Repository", "Stars", "Forks", "Open Issues"})
	table.Append([]string{
		r.FullName,
		fmt.Sprint(r.Stars),
		fmt.Sprint(r.Forks),
		fmt.Sprint(r.OpenIssues),
	})

	table.Render()
}
