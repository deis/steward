# Integration Testing

This document describes how to set up and run integration tests against Steward. Each section below explains the process for each respective steward mode.

All sections assume the following:

- You've installed [Deis Workflow](https://github.com/deis/workflow) onto your Kubernetes cluster and configured your `deis` command line client to talk to your Workflow controller
  - If you haven't, see the instructions at https://deis.com/docs/workflow/installing-workflow/
  - After you've set up your cluster, store the router IP in the `DEIS_ROUTER_IP` environment variable for later use. Find the router IP by running `kubectl get svc deis-router --namespace=deis`
- You've already deployed steward to your cluster
  - If you haven't, see the instructions at [INSTALLATION.md](./INSTALLATION.md)

# Cloud Foundry Service Broker Mode

Integration tests for service broker mode require a running service broker. To make it easy to run one, we've forked the CloudFoundry service broker from https://github.com/cloudfoundry-samples/github-service-broker-ruby and adapted it to run on Deis Workflow. The fork can be found at https://github.com/deis/cloudfoundry-github-service-broker-ruby.

## Preparation

To run a broker on Workflow, create the project:

```console
deis create --no-remote cf-sample-broker
```

Configure it with your github username and token:

```console
deis config:set GITHUB_USER=${GITHUB_USER} GITHUB_TOKEN=${GITHUB_TOKEN} -a cf-sample-broker
```

And finally, deploy it:

```console
deis pull quay.io/deis/cloudfoundry-github-service-broker:master
```

After you've set up the broker, make sure steward can talk to it by re-configuring the following configuration values on steward:

- `CF_BROKER_SCHEME` - `http`
- `CF_BROKER_HOSTNAME` - `cf-sample-broker.${DEIS_ROUTER_IP}.nip.io`
- `CF_BROKER_PORT` - `80`
- `CF_BROKER_USERNAME` - admin
- `CF_BROKER_PASSWORD` - password

## Running the tests

After you've set up your broker, compile the integration tests:

```console
make build-integration
```

Then run them:

```console
API_SCHEME=http API_HOST=steward.${DEIS_ROUTER_IP}.nip.io API_PORT=80 API_USER=deis API_PASS=steward TARGET_NAMESPACE=steward integration/integration
