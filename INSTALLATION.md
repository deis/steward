# Installing Steward

Installing steward to your Kuberentes cluster is done with [Deis Workflow](https://github.com/deis/workflow).

This section assumes you've installed Workflow onto your Kubernetes cluster and configured your `deis` command line client to talk to your Workflow controller. If you haven't, see the instructions at https://deis.com/docs/workflow/installing-workflow/.

Before installing steward, it requires a [ThirdPartyResource][3pr] called `ServiceCatalogEntry`. Install that using the following command:

```console
kubectl create -f manifests/service-catalog-entry.yaml
```

Then, create the `steward` application:

```console
deis create --no-remote steward
```

Once the application is created, it needs to be configured. Configuration differs by steward mode, and the mode is set with the `MODE` environment variable.

See the appropriate section below for the configuration values for each mode.

After you've configured steward, deploy it with:

```console
deis pull quay.io/deis/steward:devel -a steward
```

# CloudFoundry Broker Mode

The following environment variables are required to run steward in CloudFoundry broker mode (denoted `cf` in the `MODE` environment variable):

- `CF_BROKER_SCHEME` - the scheme with which to construct the URL to communicate with the backing broker. Can be either `http` or `https`
- `CF_BROKER_HOSTNAME` - the host name of the backing service broker
- `CF_BROKER_PORT` - the port of the backing service broker
- `CF_BROKER_USERNAME` - the username to use in the HTTP basic authentication to the backing service broker
- `CF_BROKER_PASSWORD` - the password to use in the HTTP basic authentication to the backing service broker

[3pr]: https://github.com/kubernetes/kubernetes/blob/master/docs/design/extending-api.md
