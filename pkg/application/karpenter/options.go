package karpenter

import (
	"eksdemo/pkg/application"
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/util"
	"errors"
	"fmt"
)

type KarpenterOptions struct {
	application.ApplicationOptions

	skipSubnetCheck bool
}

func NewOptions() (options *KarpenterOptions, flags cmd.Flags) {
	options = &KarpenterOptions{
		ApplicationOptions: application.ApplicationOptions{
			Namespace:      "karpenter",
			ServiceAccount: "karpenter",
			DefaultVersion: &application.LatestPrevious{
				Latest:   "v0.5.2",
				Previous: "v0.5.1",
			},
		},
	}

	flags = cmd.Flags{
		&cmd.BoolFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "skip-subnet-check",
				Description: "don't check subnets for required Karpenter tags",
			},
			Option: &options.skipSubnetCheck,
		},
	}
	return
}

func (o *KarpenterOptions) PreDependencies(action application.Action) error {
	if o.skipSubnetCheck || action == application.Uninstall {
		return nil
	}

	if err := util.CheckSubnets(o.ClusterName); err != nil {
		errMsg := err.Error()
		errMsg += fmt.Sprintf("\n\nKarpenter requires subnets tagged with %q to perform subnet discovery\n",
			fmt.Sprintf(util.K8stag, o.ClusterName))
		errMsg += fmt.Sprintf("Either run `eksdemo util tag-subnets -c %s` or use the `--skip-subnet-check` flag", o.ClusterName)
		return errors.New(errMsg)
	}

	return nil
}
