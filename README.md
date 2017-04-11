# Steward

[![Build Status](https://travis-ci.com/deis/steward.svg?token=UQsxfwHAz3NPyVqxkrrp&branch=master)](https://travis-ci.com/deis/steward)

*__Note__: This project is in __alpha status__. It is subject to API and architectural changes.
We will change or remove this notice when development status changes.*

This is the Kubernetes-native service broker. Modeled after the [CloudFoundry Service Broker System][cfbroker]
it functions as a gateway from your cluster-aware applications to other services, both inside and outside your cluster.

Specifically, its high-level goals are to:

1. Decouple the provider of the service from its consumers
2. Allow operators to independently scale and manage applications and the services they depend on
3. Provide a standard way for operators to:
  - Publish a catalog of services to operators or other interested parties
  - Provision a service
  - Bind an application to a service
  - Configure the application to consume the service through standard Kubernetes resources

## Glossary

* **Cloud Foundry Service Broker API**: API definition created by CloudFoundry, broadly describing provisioning/deprovisioning and binding/unbinding
* **Cloud Foundry Service Broker**: a concrete implementation of the CF Service Broker API, e.g. <https://github.com/cloudfoundry/cf-mysql-release>

* **Consumer**: An application and/or developer who would like access to some service provided by a third party. The consumer might not directly provision the service that it needs.
* **Requestable Service** (RS): is something that may be provisioned, created, or exposed on behalf of a Consumer. **Requestable Services** are not related to Kubernetes services. Example include:
    * account and access credentials for an off-cluster SaaS service like Sendgrid
    * access credentials for a relational data store like MySQL or Postgres
* **Service Plan**: is a specific "configuration" of a service, which may be expressed in terms like "small", "medium", or "large". Any specific quota or difference between plans is left up to the Service Provider to implement.
* **Service Plan Claim**: is a concrete Kubernetes `ConfigMap` which represents the desire of a **Consumer** to gain access to a **Requestable Service**. The **Service Plan Claim** references the **Requestable Service** and informs Steward where the consuming application expects to read **Service Credentials** that are created after processing the claim.
* **Service Catalog**: is a registry of **Requestable Services**  into a Kubernetes cluster.
* **Service Credentials/Configuration**: is the configuration (hostnames, usernames, passwords, etc.) meant for the **Consumer** to use for connection and authentication to the **Service Instance**.
* **Service Provider**: is a system that lives either on or off-cluster and holds implementation logic for a **Requestable Service**. In the case of a SaaS-based RS like Sendgrid, the **Service Provider** is the Sendgrid SaaS platform.
* **Service Instance**: is the entity or entities created or exposed on behalf of the **Consumer** placing a **ServicePlanClaim** and that claim being fulfilled by a **Service Provider**. Example **Service Instance**s include:
    * a provisioned AWS RDS service, a logical database and credentials
    * a logical database, username and password created on a shared RDBM system
* **Service Backend**: is a process that handles API calls as part of creating a Requestable Service, including `provision`, `deprovision`, `bind`, `unbind`. The **Service Backend** sits "in front" of the **Service Provider**. Example **Service Backend**s include:
    * a deployed Cloud Foundry Service Broker
    * a deployed Steward running in Helm Mode
    * a deployed Steward running in Command Mode


# Deploying Steward

Please see [INSTALLATION.md](./doc/INSTALLATION.md) for full instructions, including sample Kubernetes manifests, on how to deploy steward to your cluster.

Once deployed, you can view logs for each steward instance via the standard `kubectl logs` command:

```console
kubectl logs -f ${STEWARD_POD_NAME} --namespace=${STEWARD_NAMESPACE}
```

# Concepts

Steward runs a control loop to watch the Kubernetes event stream for a set of [`ThirdPartyResource`][3pr]s (called 3PRs hereafter) in one, some, or all available namespaces. It uses these 3PRs to communicate with an operator that requests a service.

A single Steward process is responsible for talking to a single **Service Provider**. If an operator wants to deploy multiple **Service Providers**, a cluster operator would deploy additional stewards. Each Steward process may run in one of three _modes_:

- CloudFoundry Broker Mode
- Helm Tiller Mode
- Command Mode

Below is an example deployment that that exemplifies multiple Steward processes each exposing a **Service Provider**:

- One Steward process configured to use CloudFoundry Broker A
- One Steward process configured to use CloudFoundry Broker B
- One Steward process configured to use Helm Tiller Server A
- One Steward Process configured to use Command A
- One Steward process configured to use Command B

Details on each Steward mode can be found below.

## Available Services

On startup, a Steward process publishes its service data as a set of `ServiceCatalogEntry` 3PRs that indicate the availability of each of a **Requestable Service**. Each **Service Catalog Entry** contains the name of the Steward instance (specified in configuration), the **Requestable Service**, and at least one **Service Plan**. Here are some example 3PRs:

- `firststeward-mysql-small`
- `secondsteward-mysql-large`
- `thirdsteward-memcache-xlarge`

Once published, an operator (or other interested party) will be able to see these `ServiceCatalogEntry` 3PRs to determine what services are available to applications in the cluster. Currently, we recommend simply using the `kubectl` command to list the catalog:

```console
kubectl get servicecatalogentry --namespace=steward
```

Please see [DATA_STRUCTURES.md](./doc/DATA_STRUCTURES.md) for a complete example of a `ServiceCatalogEntry`.

## Requesting a Service from the Catalog

Once an operator has found a service and plan they would like to use, they should submit a [ConfigMap][configMap] containing
a [`ServicePlanClaim`](./doc/DATA_STRUCTURES.md) data structure (just called `ServicePlanClaim`s hereafter).

Steward constantly watches for `ServicePlanClaim`s in its control loop. Upon finding a new `ServicePlanClaim`,
it executes the following algorithm:

  1. Looks for the `ServiceCatalogEntry` 3PR in the catalog
    - If not found, sets the `status` field to `Failed`, adds an appropriate explanation to the `statusDescription` field
    field to a human-readable description of the error, and stops processing
  2. Looks in the `action` field of the claim and takes the appropriate action
    - Valid values are `provision`, `bind`, `unbind`, `deprovision`, `create` and `delete`. See [`ServicePlanClaim` documentation](./doc/DATA_STRUCTURES.md#serviceplanclaim) for details on each value
    - If the action failed, Steward sets the `status` field to `Failed` and adds an appropriate explanation to the `statusDescription` field
  3. On success, writes values appropriate to the `action` that was submitted. See below for details on each `action`

### `provision`
- `status: provisioned`
- `instance-id: $UUID` (where `$UUID` is the instance ID returned by the provision operation)

### `bind`
- `status: bound`
- `bind-id: $UUID` (where `$UUID` is the bind ID returned by the bind operation)
- Also creates a [Secret][secrets] with the credentials data for the service. The Secret's name and namespace will be created according to the `target-name` and `target-namespace` fields passed in the `ServicePlanClaim`. See [`ServicePlanClaim` documentation](./doc/DATA_STRUCTURES.md#serviceplanclaim) for more information

### `unbind`
- `status: unbound`
- Removes the Secret created as a result of the `bind` action

### `deprovision`
- `status: deprovisioned`

### `create`

This action produces results equivalent to claims with `action: provision`, then `action: bind`

### `delete`

This action produces results equivalent to claims with `action: unbind`, then `action: deprovision`


# Backing Services

While steward provides the same consumer-facing interface (`ServiceCatalogEntry` and `ServicePlanClaim`), it can be configured to run in one of three modes to support different backing services. The modes are listed in the below list, and each link points to detailed documentation for that mode.

- [CloudFoundry Service Broker](./doc/CF_BROKER_MODE.md)
- [Helm](./doc/HELM_MODE.md) (see https://github.com/kubernetes/helm for more information on Helm)
- [Custom Docker Image](./doc/CMD_MODE.md). This mode is also called _cmd mode_

# Development & Testing

Steward is written in Go and tested with [Go unit tests](https://godoc.org/testing).

If you'd like to contribute to this project, simply fork the repository, make your changes, and submit a pull request. Please make sure to follow [these guidelines](CONTRIBUTING.md) when contributing.

[cfbroker]: https://docs.cloudfoundry.org/services/overview.html
[3pr]: https://github.com/kubernetes/kubernetes/blob/master/docs/design/extending-api.md
[rds]: https://aws.amazon.com/rds
[configMap]: http://kubernetes.io/docs/user-guide/configmap/
[secrets]: http://kubernetes.io/docs/user-guide/secrets/
[servicePlanCreation]: ./DATA_STRUCTURES.md#serviceplancreation
