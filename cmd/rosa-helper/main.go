package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/openshift-online/ocm-common/cmd/rosa-helper/create"

	"github.com/spf13/cobra"
)

var root = &cobra.Command{
	Use:   "rosa-helper",
	Short: "Command line tool for rosa support testing.",
	Long:  "Command line tool for rosa support to prepare data",
}

func init() {
	// Register the subcommands:
	root.AddCommand(create.Cmd)
}

func main() {
	// Execute the root command:
	root.SetArgs(os.Args[1:])
	err := root.Execute()
	if err != nil {
		if !strings.Contains(err.Error(), "Did you mean this?") {
			fmt.Fprintf(os.Stderr, "Failed to execute root command: %s\n", err)
		}
		os.Exit(1)
	}
}
