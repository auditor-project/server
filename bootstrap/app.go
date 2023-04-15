package bootstrap

import (
	console "auditor.z9fr.xyz/server/cmd"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:              "auditor-server",
	Short:            "auditor-server",
	Long:             "auditor server to handle workload",
	TraverseChildren: true,
}

// App root of the application
type App struct {
	*cobra.Command
}

// NewApp creates new root command
func NewApp() App {
	cmd := App{
		Command: rootCmd,
	}
	cmd.AddCommand(console.GetSubCommands(CommonModules)...)
	return cmd
}

var RootApp = NewApp()
