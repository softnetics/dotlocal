package cmd

import (
	"fmt"
	"os"

	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(uninstallCmd)
}

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall dotlocal from /usr/local/bin",
	Run: func(cmd *cobra.Command, args []string) {
		target := "/usr/local/bin/dotlocal"
		lo.Must0(os.Remove(target))
		fmt.Fprintf(os.Stderr, "Removed %s\n", target)
	},
}
