package subnets

import (
	"os"
	"strings"

	logger "github.com/openshift-online/ocm-common/pkg/log"
	vpcClient "github.com/openshift-online/ocm-common/pkg/test/vpc_client"

	"github.com/spf13/cobra"
)

var args struct {
	region string
	zones  string
	vpcID  string
	tags   string
}

var Cmd = &cobra.Command{
	Use:   "subnets",
	Short: "Create subnets",
	Long:  "Create subnets.",
	Example: `  # Create a pair of subnets with prefix 'mysubnet-' in region 'us-east-2'
  rosa-helper create subnets --name-prefix=mysubnet --region us-east-2 --vpc-id <vpc id>`,

	Run: run,
}

func init() {
	flags := Cmd.Flags()
	flags.SortFlags = false
	flags.StringVarP(
		&args.region,
		"region",
		"",
		"",
		"Vpc region",
	)
	err := Cmd.MarkFlagRequired("region")
	if err != nil {
		logger.LogError(err.Error())
		os.Exit(1)
	}
	flags.StringVarP(
		&args.zones,
		"zones",
		"",
		"",
		"Zones to create subnets in",
	)
	err = Cmd.MarkFlagRequired("zones")
	if err != nil {
		logger.LogError(err.Error())
		os.Exit(1)
	}
	flags.StringVarP(
		&args.vpcID,
		"vpc-id",
		"",
		"",
		"ID of vpc to be created",
	)
	err = Cmd.MarkFlagRequired("vpc-id")
	if err != nil {
		logger.LogError(err.Error())
		os.Exit(1)
	}
}
func run(cmd *cobra.Command, _ []string) {
	vpc, err := vpcClient.GenerateVPCByID(args.vpcID, args.region)
	if err != nil {
		panic(err)
	}
	zones := strings.Split(args.zones, ",")
	for _, zone := range zones {
		subnetMap, err := vpc.PreparePairSubnetByZone(zone)
		if err != nil {
			panic(err)
		}
		for subnetType, subnet := range subnetMap {
			logger.LogInfo("ZONE %s %s SUBNET: %s", zone, strings.ToUpper(subnetType), subnet.ID)
		}
	}
}
