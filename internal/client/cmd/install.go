package cmd

import (
	"fmt"
	"os"

	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(installCmd)
}

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install dotlocal to /usr/local/bin",
	Run: func(cmd *cobra.Command, args []string) {
		fileName := lo.Must1(os.Executable())
		target := "/usr/local/bin/dotlocal"
		lo.Must0(os.MkdirAll("/usr/local/bin", 0755))
		_ = os.Remove(target)
		lo.Must0(os.Symlink(fileName, target))
		fmt.Fprintf(os.Stderr, "Installed to %s\n", target)
	},
}
