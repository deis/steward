# Steward

This is the Kubernetes-native service broker. Modeled after the [CloudFoundry Service Broker System](https://docs.cloudfoundry.org/services/overview.html),
it functions as a gateway from your cluster-aware applications to other services, both inside and outside your cluster.

Specifically, its high-level goals are to:

1. Decouple the provider of the service from its consumers
2. To allow operators to independently scale and manage applications and the services they depend on
3. Provides a standard way for for applications and developers to discover available services
4. Provides a standard way to provision or bind to resources
5. Provides a standard way for applications to consume those resources through configuration

# Concepts

Steward is a Go program that runs in 1 or more namespaces. It has a control loop that writes to
reads from a set of [`ThirdPartyResource`](https://github.com/kubernetes/kubernetes/blob/master/docs/design/extending-api.md)s
(called 3PRs hereafter).

On startup, it publishes data to a 3PR that indicates the availability of an available service
(`ServiceProvider` hereafter). It also publishes the set of plans the service provides (`ServicePlan` hereafter).

Steward can integrate with a wide variety of cloud systems, any standard CloudFoundry service broker,
and [Helm](https://github.com/kubernetes/helm) to provide services to consumers.

# Native Provider Mode

(Note: the features described in this section are not yet implemented)

Steward has built-in support for integrating with cloud services such as Amazon AWS and Google Cloud.
For example, a steward instance can publish a `ServiceProvider` called `mysql-rds` which represents
native integration with [Amazon RDS](https://aws.amazon.com/rds).
The `ServicePlan`s for would represent the different database engines (Postgresql, Mysql, Amazon
Aurora, etc...), database size, etc...

TODO: decide on startup command

TODO:

# CloudFoundry Broker Mode

(Note: the features described in this section are not yet implemented)

Steward can also integrate with any standard CloudFoundry service broker to provide all of the services
the broker provides to the Kubernetes cluster. When a steward is started with such a configuration,
it will immediately make a request to the broker's catalog (i.e. `GET /v2/catalog`), convert the
catalog results to a `ServiceProvider` and set of `ServicePlan`s, and write the results to the
appropriate 3PR.

To start Steward in CloudFoundry broker mode, run the following command:

```console
./steward --mode=cf --hostname=broker.domain.com --username=admin --password=foo`
```

# Helm mode

(Note: the features described in this section are not yet implemented)

TODO
