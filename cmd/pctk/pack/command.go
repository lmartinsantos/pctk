package pack

import (
	"github.com/spf13/cobra"
)

var output string

var Command = &cobra.Command{
	Use:   "pack [src]",
	Short: "pack game resources into a file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		src := args[0]
		return do(src, output)
	},
}

func init() {
	Command.LocalFlags().StringVarP(
		&output, "output", "o", "resources", "output files (.idx/.dat suffixes will be added)",
	)
}
