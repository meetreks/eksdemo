package ecr_controller

import (
	"eksdemo/pkg/application"
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/installer"
	"eksdemo/pkg/resource"
	"eksdemo/pkg/resource/irsa"
	"eksdemo/pkg/template"
)

// Docs:    https://aws-controllers-k8s.github.io/community/docs/community/overview/
// Docs:    https://aws-controllers-k8s.github.io/community/reference/
// GitHub:  https://github.com/aws-controllers-k8s/ecr-controller
// Helm:    https://github.com/aws-controllers-k8s/ecr-controller/tree/main/helm
// Repo:    https://gallery.ecr.aws/aws-controllers-k8s/ecr-controller
// Version: Latest is v0.0.11 (as of 12/3/21)

func NewApp() *application.Application {
	app := &application.Application{
		Command: cmd.Command{
			Name:        "ecr-controller",
			Description: "ACK ECR Controller",
			Aliases:     []string{"ecr"},
		},

		Dependencies: []*resource.Resource{
			irsa.NewResourceWithOptions(&irsa.IrsaOptions{
				CommonOptions: resource.CommonOptions{
					Name: "ack-ecr-controller-irsa",
				},
				// https://github.com/aws-controllers-k8s/ecr-controller/blob/main/config/iam/recommended-policy-arn
				PolicyType: irsa.PolicyARNs,
				Policy:     []string{"arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryFullAccess"},
			}),
		},

		Options: &application.ApplicationOptions{
			Namespace:      "ack-system",
			ServiceAccount: "ack-ecr-controller",
			DefaultVersion: &application.LatestPrevious{
				Latest:   "v0.0.11",
				Previous: "v0.0.11",
			},
		},

		Installer: &installer.HelmInstaller{
			ReleaseName:   "ack-ecr-controller",
			RepositoryURL: "oci://public.ecr.aws/aws-controllers-k8s/ecr-chart:v0.0.11",
			ValuesTemplate: &template.TextTemplate{
				Template: valuesTemplate,
			},
		},
	}
	return app
}

const valuesTemplate = `
image:
  tag: {{ .Version }}
fullnameOverride: ack-ecr-controller
aws:
  region: {{ .Region }}
serviceAccount:
  name: {{ .ServiceAccount }}
  annotations:
    {{ .IrsaAnnotation }}
`