package cmd

import (
	"fmt"
	"os"

	"github.com/agnivo988/Repo-lyzer/internal/ui"
)

func RunMenu() {
	if err := ui.Run(); err != nil {
		fmt.Println("Error running application:", err)
		os.Exit(1)
	}
}
