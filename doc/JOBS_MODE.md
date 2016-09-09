# Jobs Mode

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
