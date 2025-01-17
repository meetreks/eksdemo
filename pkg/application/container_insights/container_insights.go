package container_insights

import (
	"eksdemo/pkg/application"
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/installer"
	"eksdemo/pkg/resource"
	"eksdemo/pkg/resource/irsa"
	"eksdemo/pkg/template"
)

// Docs:     https://docs.aws.amazon.com/AmazonCloudWatch/latest/monitoring/Container-Insights-setup-EKS-quickstart.html
// GitHub:   https://github.com/aws-samples/amazon-cloudwatch-container-insights, see Releases for versions
// Manifest: https://github.com/aws-samples/amazon-cloudwatch-container-insights/blob/master/k8s-deployment-manifest-templates/deployment-mode/daemonset/container-insights-monitoring/quickstart/cwagent-fluent-bit-quickstart.yaml
// Repo:     https://gallery.ecr.aws/cloudwatch-agent/cloudwatch-agent
// Version:  Latest is 1.247350.0b251780 aka k8s/1.3.9 (as of 03/22/22)

func NewApp() *application.Application {
	app := &application.Application{
		Command: cmd.Command{
			Name:        "container-insights",
			Description: "CloudWatch Container Insights",
			Aliases:     []string{"ci"},
		},

		Dependencies: []*resource.Resource{
			irsa.NewResourceWithOptions(&irsa.IrsaOptions{
				CommonOptions: resource.CommonOptions{
					Name: "container-insights-irsa",
				},
				PolicyType: irsa.PolicyARNs,
				Policy:     []string{"arn:aws:iam::aws:policy/CloudWatchAgentServerPolicy"},
			}),
		},

		Installer: &installer.KustomizeInstaller{
			ResourceTemplate: &template.TextTemplate{
				Template: containerInsightsManifestTemplate,
			},
			KustomizeTemplate: &template.TextTemplate{
				Template: kustomizeTemplate,
			},
		},
	}

	return addOptions(app)
}
