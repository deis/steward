# Steward

[![Build Status](https://travis-ci.com/deis/steward.svg?token=UQsxfwHAz3NPyVqxkrrp&branch=master)](https://travis-ci.com/deis/steward)

This is the Kubernetes-native service broker. Modeled after the [CloudFoundry Service Broker System][cfbroker]
it functions as a gateway from your cluster-aware applications to other services, both inside and outside your cluster.

Specifically, its high-level goals are to:

1. Decouple the provider of the service from its consumers
2. To allow operators to independently scale and manage applications and the services they depend on
3. Provides a standard way for for applications and developers to discover available services
4. Provides a standard way to provision or bind to resources
5. Provides a standard way for applications to consume those resources through configuration

# Installation

Please see [INSTALLATION.md](./INSTALLATION.md) for full instructions on how to install steward.

# Concepts

Steward is a program written in Go that runs a control loop to watch the Kubernetes event stream
for a set of [`ThirdPartyResource`][3pr]s (called 3PRs hereafter), in one, some, or all available
namespaces. It uses these 3PRs to communicate with an application that requests a service.

A single Steward process is responsible for provisioning one specific _type_ of service. For each
service type you intend to support, you should run exactly one Steward in your cluster that is
configured to provision that service type.

## Available Services

On startup, a Steward publishes its data to a set of 3PRs that indicate the availability of a
service (`ServiceProvider` hereafter). It also publishes a set of plans availble for the service
provides (`ServicePlan` hereafter).  A plan is a specific "size" of a service like `mysql:small`,
`mysql:large`, `memcache:xlarge`.
Steward can be configured to publish this data to any subset of namespaces in the Kubernetes cluster, or all of them.

Once published, an application that needs some service should query these `ServiceProvider` and
`ServicePlan` 3PRs to determine what services are available in its namespace or cluster.

## Using a Service

Once an application has found a service and plan that it would like to use, it should submit a 3PR containing
a [`ServicePlanClaim`](./DATA_STRUCTURES.md) data structure.

Steward constantly watches for `ServicePlanClaim`s in its control loop. Upon finding a new `ServicePlanClaim`,
it executes the following algorithm:

1. Looks for the `ServiceProvider`/`ServicePlan` pair in the entire service catalog
  - If not found, sets the `status` field to `Failed` and adds an appropriate explanation to the `statusDescription` field
    field to a human-readable description of the error. Then stops processing
2. Creates the service according to its configuration (see below)
  - If creation failed, Steward sets the `status` field to `Failed` and adds an appropriate explanation to the `statusDescription` field
3. Creates a [ConfigMap][configMap] that contains the non-secret data describing how to use the service. The ConfigMap will contain a field called `metadata`, which is a [`ServicePlanCreation`][servicePlanCreation] object

# Backing Services

Steward can integrate with a wide variety of cloud systems, any standard CloudFoundry service broker,
and [Helm](https://github.com/kubernetes/helm) to provide services to consumers. It's started with a `--mode` argument to indicate which backing service it should use.

## Native Provider Mode

(Note: the features described in this section are not yet implemented)

Steward has built-in support for integrating with cloud services such as Amazon AWS and Google Cloud.
For example, a steward instance can publish a `ServiceProvider` called `mysql-rds` which represents
native integration with [Amazon RDS][rds].
The `ServicePlan`s for would represent the different database engines (Postgresql, Mysql, Amazon
Aurora, etc...), database size, etc...

TODO: decide on startup command

## CloudFoundry Broker Mode

(Note: the features described in this section are not yet implemented)

Steward can also integrate with any standard CloudFoundry service broker to provide all of the services
the broker provides to the Kubernetes cluster. When a steward is started with such a configuration,
it will immediately make a request to the broker's catalog (i.e. `GET /v2/catalog`), convert the
catalog results to a `ServiceProvider` and set of `ServicePlan`s, and write the results to the
appropriate 3PR.

To start Steward in CloudFoundry broker mode, run the following commands:

```console
export STEWARD_MODE=cf
export STEWARD_CF_HOSTNAME=broker.domain.com
export STEWARD_CF_USERNAME=admin
export STEWARD_CF_PASSWORD=foo
./steward
```

## Helm Mode

Steward can integrate with any standard [Helm Tiller](https://github.com/kubernetes/helm) server to provide a broker in front of [Helm charts](https://github.com/kubernetes/charts). See Helm's [quick start](https://github.com/kubernetes/helm/blob/master/docs/quickstart.md) documentation for a guide on installing both the Helm CLI and the Tiller server. Note that because of [Helm issue #1083](https://github.com/kubernetes/helm/issues/1083#issuecomment-243520610), you'll have to install the Tiller server v2.0.0-alpha3 with the following command:

```console
helm init --image gcr.io/kubernetes-helm/tiller:v2.0.0-alpha.3
```

### Configuration

Configure steward to run in helm mode by setting the `STEWARD_MODE` environment variable to `helm`. Then, configure its behavior with the following environment variables:

- `HELM_TILLER_IP` - the IP address of the Tiller server to talk to
- `HELM_TILLER_PORT` - the port of that the Tiller server at `HELM_TILLER_IP` is listening on
- `HELM_CHART_URL` - the URL of the chart to install
- `HELM_CHART_INSTALL_NAMESPACE` - the Kubernetes namespace in which to install charts
- `HELM_PROVISION_BEHAVIOR` - See the "Provision and Deprovision Operations" section below
- `HELM_SERVICE_ID` - the service ID to list in the service catalog for this steward instance
- `HELM_SERVICE_NAME` - the service name to list in the service catalog for this steward instance
- `HELM_SERVICE_DESCRIPTION` - the service description to list in the service catalog for this steward instance
- `HELM_PLAN_ID` - the plan ID to list in the service catalog for this steward instance
- `HELM_PLAN_NAME` - the plan name to list in the service catalog for this steward instance
- `HELM_PLAN_DESCRIPTION` - the plan description to list in the service catalog for this steward instance


### Provision and Deprovision Operations

Steward can be configured to take one of two different actions when it receives a provision or deprovision operation. Configure this behavior with one of the folowing two values for the `HELM_PROVISION_BEHAVIOR` env var:

- `noop` - steward will not install or uninstall the chart specified at `HELM_CHART_URL`. In this confifguration, steward expects the following state of the cluster when it starts:
  - The operator has installed the same chart as specfied in the `HELM_CHART_URL` environment variable
  - The ConfigMaps specified in the chart's `values.yaml` file (see below) exist and represent valid credentials for the bindable services the chart has started
  - The operator will not uninstall or otherwise modify the chart in such a way that bound consumers cannot properly interact with the chart's exposed services
- `active` - steward will download and install the chart specified at `HELM_CHART_URL` on provision operations, and uninstall it on deprovision operations

### Bind and Unbind Operations

On bind operations, steward reads a set of chart-specified `ConfigMap`s to get the credentials to return to the consumer. The namespace and name for each `ConfigMap` should be specified in the chart's top-level `values.yaml` file as such:

```yaml
stewardConfigMaps:
  - name: cm1
    namespace: ns1
  - name: cm2
    namespace: ns2
```

On unbind operations, steward will attempt to delete this same set of `ConfigMap`s.


## Jobs Mode

In jobs mode, Steward is capable of being configured to delegate all discreet operations such as `provision`, `bind`, `unbind`, and `deprovision` to a containerized binary or executable script.

In jobs mode, containerized brokers are not long lived. They run only for as long as is required to complete a discrete operation, then exit.

### Configuration

#### Environment Variables

Configure Steward to run in jobs mode by setting the `STEWARD_MODE` environment variable to `jobs`. Then configure its behavior with the following environment variables:

* `POD_NAMESPACE`, the namespace within which this Steward instance exists. It's best to configure this using Kubernetes [downward API](http://kubernetes.io/docs/user-guide/downward-api/).
* `JOBS_IMAGE`, the URL for a Docker image containing the binary or script to which discrete operations will be delegated.
* `JOBS_CONFIG_MAP`, the _optional_ name of a configmap within the same namespace as the Steward instance. This configmap may include (non-sensitive) configuration for the containers in which discrete operations will be executed.
* `JOBS_SECRET`, the _optional_ name of a secret within the same namespace as the Steward instance. This secret may include (sensitive) configuration for the containers in which discrete operations will be executed.

#### Configmap and/or Secret (optional)

If a broker requires specific configuration (e.g. RDBMS connection details or cloud platform credentials), these can be conveyed via a Kubernetes configmap and/or secret. When _referenced_ by name in Steward's own configuration (as mentioned in the previous section), Steward will make all such configuration values within available to a broker using _two_ different mechanisms. Broker authors may opt to utilize whichever they find more convenient.

##### Volume Mounts

1. An optional configmap referenced via Steward's `JOBS_CONFIG_MAP` environment variable is mounted to the location `/config` within the broker container. As per the usual case when mounting configmaps as volumes, the value of each key is written to a flat file at `/config/<key name>`.

1. An optional secret referenced via Steward's `JOBS_SECRET` environment variable is mounted to the location `/secret` within the broker container. As per the usual case when mounting secrets as volumes, the value of each key is written to a flat file at `/secret/<key name>`.

##### Downward API

Before continuing, make note that keys in Kubernetes configmaps and secrets are constrained by regular expressions that don't permit uppercase letters or underscores.

For each key in an optional configmap referenced via Steward's `JOBS_CONFIG_MAP` environment variable _and_ for each key in an optional secret referenced via Steward's `JOBS_SECRET` environment variable, Steward utilizes the Kubernetes downward API to expose an environment variable within the broker's container. In doing so, it upcases the key and replaces all `.` with `_`. For example, a key of `aws.key` within a configmap becomes the environment variable `AWS_KEY` and a key of `aws.secret.key` within a secret becomes the environment variable `AWS_SECRET_KEY`.

### Interface

Jobs mode brokers may be implemented in any compiled language or scripting language. To ensure interoperability with Steward, however, certain contractual requirements must be met.

#### Executable

The executable, whether a compiled program or a script, _must_ exist at `/bin/broker`.

#### Subcommands

The following subcommands must be honored. The precise meaning of each is contextual and broker-specific. Some subcommands may even be implemented as no-ops.

1. `catalog`-- list the services and plans offered by the broker.
1. `provision`-- generally, create or allocate services, resources, etc.
1. `bind`-- generally, create or share connection details, credentials, etc.
1. `unbind`-- generally, destroy credentials formerly used by the requestor.
1. `deprovision`-- generally, destroy or deallocate services, resources, etc.

#### Flags

Inputs to the subcommands listed above are specified using flags (to avoid requiring inputs be entered in any precise order-- which would easily be subject to programmer error).

Each flag must be accepted, even if the information it conveys is contextually unnecessary for a given broker. 

The following table summarizes the flags that must be supported by each subcommand:

| Subcommand / Flag | `--service-id` | `--plan-id`  | `--instance-id`  | `--binding-id` |
|-------------------|----------------|--------------|------------------|----------------|
| `catalog`         |                |              |                  |                |
| `provision`       | X              | X            | X                |                |
| `bind`            | X              | X            | X                | X              |
| `unbind`          | X              | X            | X                | X              |
| `deprovision`     | X              | X            | X                |                |


#### Output

The output of each subcommand listed above must be written to `STDOUT`, must be JSON, and must conform to the schema prescribed for the analogous operation in the [Cloud Foundry Service Broker API](https://docs.cloudfoundry.org/services/api.html).

| Subcommand    | URL                                                            |
|---------------|----------------------------------------------------------------|
| `catalog`     | https://docs.cloudfoundry.org/services/api.html#catalog-mgmt   |
| `provision`   | https://docs.cloudfoundry.org/services/api.html#provisioning   |
| `bind`        | https://docs.cloudfoundry.org/services/api.html#binding        |
| `unbind`      | https://docs.cloudfoundry.org/services/api.html#unbinding      |
| `deprovision` | https://docs.cloudfoundry.org/services/api.html#deprovisioning |


__Currently, no other output (including debug output) may be written to `STDOUT`, or else Steward will fail to parse responses.__


# Development & Testing

Steward is written in Go and tested with [Go unit tests](https://godoc.org/testing).

If you'd like to contribute to this project, simply fork the repository, make your changes, and submit a pull request. Please make sure to follow [these guidelines](https://deis.com/docs/workflow/contributing/submitting-a-pull-request/) when contributing.

[cfbroker]: https://docs.cloudfoundry.org/services/overview.html
[3pr]: https://github.com/kubernetes/kubernetes/blob/master/docs/design/extending-api.md
[rds]: https://aws.amazon.com/rds
[configMap]: http://kubernetes.io/docs/user-guide/configmap/
[servicePlanCreation]: ./DATA_STRUCTURES.md#serviceplancreation
