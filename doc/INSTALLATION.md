# Installing Steward

Installation of Steward into a running Kubernetes cluster is facilitated through a Make target. Users wishing to familiarize themselves with the particulars of the deployment will want to examine this repository's `Makefile` and artifacts in the `manifests/` directory.

All subsequent sections of this document assume that you have a running Kubernetes cluster and that your `kubectl` client is properly configured to interact with that cluster.

## Prerequisites

Pre-requisites for Steward vary based on the mode with which steward is run. Please see the below documents according to the mode with which you wish to configure and run steward.

- CloudFoundry Broker: https://github.com/deis/steward/blob/master/doc/CF_BROKER_MODE.md
  - Note that if you're trying Steward for the first time or are hacking on Steward, the Steward team has provided a trivial Cloud Foundry [sample broker][cf-sample-broker]. See that project's [README.md](https://github.com/deis/cf-sample-broker/blob/master/README.md) for installation instructions.
- Helm: https://github.com/deis/steward/blob/master/doc/HELM_MODE.md
- Command: https://github.com/deis/steward/blob/master/doc/CMD_MODE.md

## Deploy Steward

With all configuration now set, Steward can be deployed as follows:

### Cloud Foundry Mode

```
$ make deploy-cf
```

Or build and deploy from source using:

```
$ make dev-deploy-cf
```

### Helm Mode

```
$ make deploy-helm
```

Or build and deploy from source using:

```
$ make dev-deploy-helm
```

### CMD Mode

```
$ make deploy-cmd
```

Or build and deploy from source using:

```
$ make dev-deploy-cmd
```


For details on Steward's pure Kubernetes-based workflow, please refer to [README.md](./README.md).

[cf-sample-broker]: https://github.com/deis/cf-sample-broker
