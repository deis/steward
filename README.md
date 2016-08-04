# Steward

This is the Kubernetes-native service broker. Modeled after the [CloudFoundry Service Broker System][cfbroker]
it functions as a gateway from your cluster-aware applications to other services, both inside and outside your cluster.

Specifically, its high-level goals are to:

1. Decouple the provider of the service from its consumers
2. To allow operators to independently scale and manage applications and the services they depend on
3. Provides a standard way for for applications and developers to discover available services
4. Provides a standard way to provision or bind to resources
5. Provides a standard way for applications to consume those resources through configuration

# Concepts

Steward is a program written in Go that runs a control loop to watch the Kubernetes event stream
for a set of [`ThirdPartyResource`][3pr]s (called 3PRs hereafter), in one, some, or all available
namespaces. It uses these 3PRs to communicate with an application that requests a service.

A single Steward process is responsible for provisioning one specific _type_ of service. For each
service type you intend to support, you should run exactly one Steward in your cluster that is
configured to provision that service type.

Each steward instance runs an API server that implements part of the [CloudFoundry Broker API](https://docs.cloudfoundry.org/services/api.html). See [API.md](./API.md) for more information.

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

## Helm mode

TODO

# Development & Testing

Steward is written in Go, and tested with [Go unit tests](https://godoc.org/testing) and an integration test suite. See [INTEGRATION_TESTING.md](./INTEGRATION_TESTING.md) for more information on the latter.

If you'd like to contribute to this project, simply fork the repository, make your changes, and submit a pull request. Please make sure to follow [these guidelines](https://deis.com/docs/workflow/contributing/submitting-a-pull-request/) when contributing.

[cfbroker]: https://docs.cloudfoundry.org/services/overview.html
[3pr]: https://github.com/kubernetes/kubernetes/blob/master/docs/design/extending-api.md
[rds]: https://aws.amazon.com/rds
[configMap]: http://kubernetes.io/docs/user-guide/configmap/
[servicePlanCreation]: ./DATA_STRUCTURES.md#serviceplancreation
