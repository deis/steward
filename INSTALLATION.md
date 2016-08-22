# Installing Steward

Installation of Steward into a running Kubernetes cluster is facilitated through a series of Make targets. Users wishing to familiarize themselves with the particulars of the deployment will want to examine artifacts in this repository's `manifests/` directory.

All subsequent sections of this document assume that you have a running Kubernetes cluster and that your `kubectl` client is properly configured to interact with that cluster.

## Prerequisites

### A Service Broker

Since Steward acts as a _service broker gateway_, a Steward instance is of no use unless backed by a service broker.

Before proceeding further, be sure that the service broker you wish to expose via Steward is available. Whether it runs on-or-off-cluster is inconsequential.

If you are trying Steward for the first time or are hacking on Steward, the Steward team has provided a trivial Cloud Foundry [sample broker][cf-sample-broker]. See that project's [README.md](https://github.com/deis/cf-sample-broker/blob/master/README.md) for installation instructions.

### The Namespace

Steward will be installed into the `steward` namespace. As such, it is important to first ensure the existence of this namespace. If it does not exist, it is easily create like so:

```
$ make install-namespace
```

### The ServiceCatalogEntry Third Party Resource

Steward requires a [ThirdPartyResource](https://github.com/kubernetes/kubernetes/blob/master/docs/design/extending-api.md) called `ServiceCatalogEntry` to be pre-defined within your Kubernetes cluster. This can be achieved easily as follows:

```
$ make install-3prs
```

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
$ make install-steward
```

For details on Steward's pure Kubernetes-based workflow, please refer to [README.md](./README.md).

[cf-sample-broker]: https://github.com/deis/cf-sample-broker
