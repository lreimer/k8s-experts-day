package main

import (
	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v2/go/kubernetes/apps/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v2/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v2/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		appLabels := pulumi.StringMap{
			"app": pulumi.String("nginx"),
		}

		deployment, err := appsv1.NewDeployment(ctx, "nginx", &appsv1.DeploymentArgs{
			Metadata: metav1.ObjectMetaArgs{
				Name: pulumi.String("nginx-deployment"),
			},
			Spec: appsv1.DeploymentSpecArgs{
				Selector: &metav1.LabelSelectorArgs{
					MatchLabels: appLabels,
				},
				Replicas: pulumi.Int(2),
				Template: &corev1.PodTemplateSpecArgs{
					Metadata: &metav1.ObjectMetaArgs{
						Labels: appLabels,
					},
					Spec: &corev1.PodSpecArgs{
						Containers: corev1.ContainerArray{
							corev1.ContainerArgs{
								Name:  pulumi.String("nginx"),
								Image: pulumi.String("nginx"),
							}},
					},
				},
			},
		})
		if err != nil {
			return err
		}

		service, err := corev1.NewService(ctx, "nginx", &corev1.ServiceArgs{
			Metadata: metav1.ObjectMetaArgs{
				Name: pulumi.String("nginx-service"),
			},
			Spec: corev1.ServiceSpecArgs{
				Type: pulumi.String("LoadBalancer"),
				Ports: corev1.ServicePortArray{
					corev1.ServicePortArgs{
						Port:     pulumi.Int(80),
						Protocol: pulumi.String("TCP"),
					}},
				Selector: appLabels,
			},
		})
		if err != nil {
			return err
		}

		ctx.Export("deploymentName", deployment.Metadata.Elem().Name())
		ctx.Export("serviceName", service.Metadata.Elem().Name())

		return nil
	})
}
