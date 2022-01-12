# Kubernetes Stuff

This directory contains files needed for running this service on our Kubernetes
cluster.

NOTE: to deploy the ingress on AWS, don't bother with the helm chart - it will
deploy a Classic load balancer to connect to the ingress which is really not
what we want. Classic load balancer is kind of funky to configure and I've had
multiple issues with the default setup deployed by the ingress-nginx helm chart.
Instead, we want a Network Load Balancer (NLB) - to deploy the ingress with this
kind of load balancer we need to create a `Service` resource with
`type: LoadBalancer`. The `ingress-nginx.yaml` file includes a complete example
of deploying this sort of setup.
