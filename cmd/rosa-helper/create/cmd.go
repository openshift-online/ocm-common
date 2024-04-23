package create

import (
	"github.com/openshift-online/ocm-common/cmd/rosa-helper/create/proxy"
	"github.com/openshift-online/ocm-common/cmd/rosa-helper/create/sg"
	"github.com/openshift-online/ocm-common/cmd/rosa-helper/create/subnets"
	"github.com/openshift-online/ocm-common/cmd/rosa-helper/create/vpc"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"add"},
	Short:   "Create a resource from stdin",
	Long:    "Create a resource from stdin",
}

func init() {
	Cmd.AddCommand(vpc.Cmd)
	Cmd.AddCommand(sg.Cmd)
	Cmd.AddCommand(subnets.Cmd)
	Cmd.AddCommand(proxy.Cmd)
}
