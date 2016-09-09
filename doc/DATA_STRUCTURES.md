# `ServiceCatalogEntry`

This object is written to the `steward` namespace and represents a single (service, plan) pair that at least one steward can provision and bind. It has the following fields:

- `service_info` - a JSON object containing the information on the service. See below for a description of the fields in this object
- `service_plan` - a JSON object containing the information on the service's plan. See below for a description of the fields in this Object

## `service_info`

This object contains information for a service offered by steward. It has the below fields. Each is a string unless otherwise indicated.

- `name` - the name of the service
- `id` - the ID of the service
- `description` - the description of the service
- `plan_updateable` - a boolean indicating whether the service's plans are updateable

## `service_plan`

This object contains information on an individual plan for a service offered by steward. It has the below fields. Each is a string unless otherwise indicated.

- `id` - the ID of the plan
- `name` - the name of the plan
- `descripton` - the description of the plan
- `free` - a boolean indicating whether this plan is free

## Example

Below is an example of a `ServiceCatalogEntry`, in yaml format:

```yaml
apiVersion: steward.deis.io/v1
description: an alpine pod provisioned with a backing Tiller server (there is only
  1 plan for this chart)
kind: ServiceCatalogEntry
metadata:
  creationTimestamp: 2016-08-30T22:37:55Z
  name: helm-alpine-standard
  namespace: steward
  resourceVersion: "504405"
  selfLink: /apis/steward.deis.io/v1/namespaces/steward/servicecatalogentries/helm-alpine-standard
  uid: 649ed98f-6f02-11e6-9018-42010a800069
service_info:
  description: an alpine pod provisioned with a backing Tiller server
  id: helm-alpine
  name: alpine-server
  plan_updateable: false
service_plan:
  description: there is only 1 plan for this chart
  free: false
  id: standard
  name: standard
```

# `ServicePlanClaim`

This object is submitted by the application as JSON in a [`ConfigMap`][configMap] (to become a [`ThirdPartyResource`][3pr] after https://github.com/deis/steward/issues/17 is fixed) when the application wants Steward to create a new service for its use. Steward then mutates the object to communicate the status of the service creation operation. Applications may watch the event stream for this object to watch progress of service creation.

- `target-name` - the name of the [`ConfigMap`][configMap] that steward should write the resulting credentials
- `target-namespace` - the namespace that steward should write the [`ConfigMap`][configMap] with the resulting credentials
- `service-provider` - the name of the `ServiceProvider` the application wants
- `service-plan` - the name of the `ServicePlan` the application wants
- `claim-id` - an application-generated [UUID][uuid]
- `action` - the application-specified action to take. Valid values are `provision`,`bind`, `unbind`, `deprovision`, `create` and `delete`. A few more notes:
  - Steward will never modify this value
  - `create` will execute both the `provision` and `bind` actions, in order
  - `delete` will execute both the `unbind` and `deprovision` actions, in order
  - All new `ServicePlanClaim`s submitted by applications must have `action` set to `provision` or `create`
  - If steward encounters an error, the actions it has already taken will not be rolled back. See the following examples:
    - If you submit a claim with `action: create` and the bind step fails, the provision step will not be rolled back
    - If you submit a claim with `action: bind` with a `target-name` and/or `target-namespace` value that points to a ConfigMap that already exists, steward will execute the bind action on the backend but will fail to write the new credentials ConfigMap. The backend bind action will not be rolled back
  - It is an error for this field to be empty
- `status` - the current status of the claim. Steward will modify this value, but will ignore any modifications by the application. Valid values and short descriptions are listed below:
  - `provisioning` - immediately after `action` is set to `provision`
  - `provisioned` - after `action` is set to `provision` and the provisioning process succeeded
  - `binding` - immediately after `action` is set to `bind`
  - `bound` - after `action` is set to `bind` and the binding process succeeded
  - `unbinding` - immediately after `action` is set to `unbind`
  - `unbound` - after `action` is set to `unbind` and the unbinding process succeeded
  - `deprovisioning` - immediately after `action` is set to `deprovision`
  - `deprovisioned` - after `action` is set to `deprovision` and the deprovisioning process succeeded
  - `failed` - after any `action` failed
- `status-description` - a human-readable explanation of the current `status`. Steward will modify this value, but will ignore any modifications by the application
- `instance-id` - for internal use only. The application should not modify this field
- `bind-id` - for internal use only. The application should not modify this field
- `extra` - for internal use only. The application should not modify this field

[3pr]: https://github.com/kubernetes/kubernetes/blob/master/docs/design/extending-api.md
[uuid]: https://tools.ietf.org/html/rfc4122
[configMap]: http://kubernetes.io/docs/user-guide/configmap/
