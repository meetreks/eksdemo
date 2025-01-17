package install

import (
	"eksdemo/pkg/application"
	"eksdemo/pkg/application/ack/s3_controller"
	"eksdemo/pkg/application/adot_operator"
	"eksdemo/pkg/application/appmesh_controller"
	"eksdemo/pkg/application/aws_fluentbit"
	"eksdemo/pkg/application/aws_lb"
	"eksdemo/pkg/application/cert_manager"
	"eksdemo/pkg/application/cluster_autoscaler"
	"eksdemo/pkg/application/container_insights"
	"eksdemo/pkg/application/container_insights_prom"
	"eksdemo/pkg/application/ebs_csi"
	"eksdemo/pkg/application/efs_csi"
	"eksdemo/pkg/application/external_dns"
	"eksdemo/pkg/application/fsx_lustre_csi"
	"eksdemo/pkg/application/grafana_amp"
	"eksdemo/pkg/application/karpenter"
	"eksdemo/pkg/application/keycloak_amg"
	"eksdemo/pkg/application/kube_prometheus"
	"eksdemo/pkg/application/metrics_server"
	"eksdemo/pkg/application/prometheus_amp"

	"github.com/spf13/cobra"
)

func NewUninstallCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "uninstall",
		Short:   "uninstall application and delete dependencies",
		Aliases: []string{"uninst"},
	}

	// Don't show flag errors for uninstall without a subcommand
	cmd.DisableFlagParsing = true

	cmd.AddCommand(NewUninstallAckCmd())
	for _, c := range NewUninstallAliasCmds(ack, "ack-") {
		cmd.AddCommand(c)
	}
	cmd.AddCommand(adot_operator.NewApp().NewUninstallCmd())
	cmd.AddCommand(appmesh_controller.NewApp().NewUninstallCmd())
	cmd.AddCommand(aws_fluentbit.NewApp().NewUninstallCmd())
	cmd.AddCommand(aws_lb.NewApp().NewUninstallCmd())
	cmd.AddCommand(cert_manager.NewApp().NewUninstallCmd())
	cmd.AddCommand(cluster_autoscaler.NewApp().NewUninstallCmd())
	cmd.AddCommand(container_insights.NewApp().NewUninstallCmd())
	cmd.AddCommand(container_insights_prom.NewApp().NewUninstallCmd())
	cmd.AddCommand(ebs_csi.NewApp().NewUninstallCmd())
	cmd.AddCommand(efs_csi.NewApp().NewUninstallCmd())
	cmd.AddCommand(NewUninstallExampleCmd())
	for _, c := range NewUninstallAliasCmds(exampleApps, "example-") {
		cmd.AddCommand(c)
	}
	cmd.AddCommand(external_dns.NewApp().NewUninstallCmd())
	cmd.AddCommand(grafana_amp.NewApp().NewUninstallCmd())
	cmd.AddCommand(fsx_lustre_csi.NewApp().NewUninstallCmd())
	cmd.AddCommand(NewUninstallIstioCmd())
	for _, c := range NewUninstallAliasCmds(istioApps, "istio-") {
		cmd.AddCommand(c)
	}
	cmd.AddCommand(karpenter.NewApp().NewUninstallCmd())
	cmd.AddCommand(keycloak_amg.NewApp().NewUninstallCmd())
	cmd.AddCommand(kube_prometheus.NewApp().NewUninstallCmd())
	cmd.AddCommand(metrics_server.NewApp().NewUninstallCmd())
	cmd.AddCommand(prometheus_amp.NewApp().NewUninstallCmd())

	cmd.AddCommand(s3_controller.NewApp().NewUninstallCmd())

	return cmd
}

// This creates alias commands for subcommands under INSTALL
func NewUninstallAliasCmds(appList []func() *application.Application, prefix string) []*cobra.Command {
	cmds := make([]*cobra.Command, 0, len(appList))

	for _, app := range appList {
		a := app()
		a.Command.Name = prefix + a.Command.Name
		a.Command.Hidden = true
		for i, alias := range a.Command.Aliases {
			a.Command.Aliases[i] = prefix + alias
		}
		cmds = append(cmds, a.NewUninstallCmd())
	}

	return cmds
}
