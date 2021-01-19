import * as k8s from "@pulumi/kubernetes";
import * as kx from "@pulumi/kubernetesx";

const appLabels = { app: "nginx" };

const deployment = new k8s.apps.v1.Deployment("nginx", {
    metadata: {
        name: "nginx-deployment",
    },
    spec: {
        selector: { matchLabels: appLabels },
        replicas: 2,
        template: {
            metadata: { labels: appLabels },
            spec: { containers: [{ name: "nginx", image: "nginx" }] }
        }
    }
});

const service = new k8s.core.v1.Service("nginx", {
    metadata: {
        name: "nginx-service",
    },
    spec: {
        type: "LoadBalancer",
        ports: [
            { port: 80, protocol: "TCP" }
        ],
        selector: appLabels
    }
});

export const deploymentName = deployment.metadata.name;
export const serviceName = service.metadata.name;
