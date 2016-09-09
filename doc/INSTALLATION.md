# Installing Steward

Installation of Steward into a running Kubernetes cluster is facilitated through a Make target. Users wishing to familiarize themselves with the particulars of the deployment will want to examine this repository's `Makefile` and artifacts in the `manifests/` directory.

All subsequent sections of this document assume that you have a running Kubernetes cluster and that your `kubectl` client is properly configured to interact with that cluster.

## Prerequisites

### A Service Broker

Since Steward acts as a _service broker gateway_, a Steward instance is of no use unless backed by a service broker.

Before proceeding further, be sure that the service broker you wish to expose via Steward is available. Whether it runs on-or-off-cluster is inconsequential.

If you are trying Steward for the first time or are hacking on Steward, the Steward team has provided a trivial Cloud Foundry [sample broker][cf-sample-broker]. See that project's [README.md](https://github.com/deis/cf-sample-broker/blob/master/README.md) for installation instructions.

## Installation Steps

### Configure Broker Details

For Steward instances running in Cloud Foundry mode, the following environment variables are required to describe the connection to and credentials for the backing broker:

- `CF_BROKER_NAME` - the name of the broker for which this Steward instance will be a gateway
- `CF_BROKER_SCHEME` - the scheme with which to construct the URL to communicate with the backing broker. Can be either `http` or `https`
- `CF_BROKER_HOSTNAME` - the host name of the backing service broker
- `CF_BROKER_PORT` - the port of the backing service broker
- `CF_BROKER_USERNAME` - the username to use in the HTTP basic authentication to the backing service broker
- `CF_BROKER_PASSWORD` - the password to use in the HTTP basic authentication to the backing service broker

Before proceeding, refer to the documentation for the backing broker you will be exposing via Steward.

If using Steward to expose the [cf-sample-broker], that broker's connection details are easily sourced into the current bash shell as described [here](https://github.com/deis/cf-sample-broker/blob/master/README.md#source-connection-details).

### Deploy Steward

With all configuration now set, Steward can be deployed as follows:

```
$ make deploy
```

For details on Steward's pure Kubernetes-based workflow, please refer to [README.md](./README.md).

[cf-sample-broker]: https://github.com/deis/cf-sample-broker
