# Rewrite

This example demonstrates how to use the Rewrite annotations

## Prerequisites

You will need to make sure your Ingress targets exactly one Ingress controller by specifying the [ingress.class annotation](https://kubernetes.github.io/ingress-nginx/user-guide/multiple-ingress/), and that you have an ingress controller [running](https://kubernetes.github.io/ingress-nginx/deploy/) in your cluster.

## Deployment

Rewriting can be controlled using the following annotations:

| Name                                           | Description                                                  | Values |
| :--------------------------------------------- | :----------------------------------------------------------- | :----- |
| nginx.ingress.kubernetes.io/rewrite-target     | Target URI where the traffic must be redirected              | string |
| nginx.ingress.kubernetes.io/ssl-redirect       | Indicates if the location section is accessible SSL only (defaults to True when Ingress contains a Certificate) | bool   |
| nginx.ingress.kubernetes.io/force-ssl-redirect | Forces the redirection to HTTPS even if the Ingress is not TLS Enabled | bool   |
| nginx.ingress.kubernetes.io/app-root           | Defines the Application Root that the Controller must redirect if it's in '/' context | string |
| nginx.ingress.kubernetes.io/use-regex          | Indicates if the paths defined on an Ingress use regular expressions | bool   |

[更多信息](https://kubernetes.github.io/ingress-nginx/examples/rewrite/)

https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md#types-kinds