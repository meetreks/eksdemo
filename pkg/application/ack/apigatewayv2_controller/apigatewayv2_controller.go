package apigatewayv2_controller

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
// GitHub:  https://github.com/aws-controllers-k8s/apigatewayv2-controller
// Helm:    https://github.com/aws-controllers-k8s/apigatewayv2-controller/tree/main/helm
// Chart:   https://gallery.ecr.aws/aws-controllers-k8s/apigatewayv2-chart
// Repo:    https://gallery.ecr.aws/aws-controllers-k8s/apigatewayv2-controller
// Version: Latest is v0.0.17 (as of 03/21/22)

func NewApp() *application.Application {
	app := &application.Application{
		Command: cmd.Command{
			Name:        "apigatewayv2-controller",
			Description: "ACK API Gateway v2 Controller",
			Aliases:     []string{"apigatewayv2", "apigwv2"},
		},

		Dependencies: []*resource.Resource{
			irsa.NewResourceWithOptions(&irsa.IrsaOptions{
				CommonOptions: resource.CommonOptions{
					Name: "ack-apigatewayv2-controller-irsa",
				},
				// https://github.com/aws-controllers-k8s/apigatewayv2-controller/blob/main/config/iam/recommended-policy-arn
				PolicyType: irsa.PolicyARNs,
				Policy: []string{
					"arn:aws:iam::aws:policy/AmazonAPIGatewayAdministrator",
					"arn:aws:iam::aws:policy/AmazonAPIGatewayInvokeFullAccess",
				},
			}),
		},

		Options: &application.ApplicationOptions{
			Namespace:      "ack-system",
			ServiceAccount: "ack-apigatewayv2-controller",
			DefaultVersion: &application.LatestPrevious{
				Latest:   "v0.0.17",
				Previous: "v0.0.16",
			},
		},

		Installer: &installer.HelmInstaller{
			ReleaseName:   "ack-apigatewayv2-controller",
			RepositoryURL: "oci://public.ecr.aws/aws-controllers-k8s/apigatewayv2-chart:v0.0.17",
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
fullnameOverride: ack-apigatewayv2-controller
aws:
  region: {{ .Region }}
serviceAccount:
  name: {{ .ServiceAccount }}
  annotations:
    {{ .IrsaAnnotation }}
`
