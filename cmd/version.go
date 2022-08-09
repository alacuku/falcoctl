package cmd

import (
	"github.com/falcosecurity/falcoctl/pkg/version"
	"github.com/spf13/cobra"
)

const falcoctlVersionLongHelp = `Print the falco CLI version.

The CLI version is embedded in the binary and directly displayed.

Examples:
$ falcoctl version
`

// newVersion creates the `version` command.
func newVersion() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "version",
		DisableFlagsInUseLine: true,
		Short:                 "Print the falco CLI version",
		Long:                  falcoctlVersionLongHelp,
		Args:                  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return version.Run()
		},
	}

	return cmd
}
