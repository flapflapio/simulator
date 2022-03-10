#!/bin/bash

main() (
    # Set your KUBECONFIG like so
    # export KUBECONFIG="/home/ethanbenabou/git/490/simulator/PRIVATE/kubeconfig-flapflap.yml"

    # Var
    export CHART_NAME=simulator
    export NAMESPACE=flapflap
    export TAG=main-0.1.0-24

    # You will need to set this before running the script
    # export CLIENT_SECRET=xxxxxxxxxxxx

    helm upgrade \
        --namespace "$NAMESPACE" \
        --install \
        --atomic \
        --wait \
        --timeout '1200s' \
        --set imageTag="$TAG" \
        --set oauth2.clientSecret="${CLIENT_SECRET}" \
        "$CHART_NAME" \
        ./simulator-chart
)

[[ ${BASH_SOURCE[0]} == $0 ]] && main "$@"
