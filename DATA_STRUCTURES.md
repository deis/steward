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

# `ServicePlanClaim`

NOTE: `ServicePlanClaim`s are not yet supported because of [issue #17](https://github.com/deis/steward/issues/17) and its related issues

This object is submitted by the application as JSON in a [`ThirdPartyResource`][3pr]
when the application wants Steward to create a new service for its use. Steward then mutates the object
to communicate the status of the service creation operation. Applications may watch the event stream
for this object to watch progress of service creation.

- `serviceProvider` - the name of the `ServiceProvider` the application wants
- `servicePlan` - the name of the `ServicePlan` the application wants
- `claimID` - an application-generated [UUID][uuid]
- `action` - the application-specified action to take. Valid values are `create` and `delete`. Steward will never modify this value. Applications should submit a new `ServicePlanClaim` with `create` in this field. It is an error for this field to be empty
- `status` - the current status of the claim. Steward will modify this value, but will ignore any modifications by the application. Valid values are `Failed`, `Creating`, `Created`, and `Deleted`.
- `statusDescription` - a human-readable explanation of the current `Status`. Steward will modify this value, but will ignore any modifications by the application

# `ServicePlanCreation`

This object is returned in a [ConfigMap][configMap] by Steward after a service is successfully created (i.e.
the `status` field of the associated `ServicePlanClaim` object becomes `Created`).

- `claimID` - the claim ID that was submitted in the associated `ServicePlanClaim`
- `secretNames` - a list of names of secrets that were created along with this object

[3pr]: https://github.com/kubernetes/kubernetes/blob/master/docs/design/extending-api.md
[uuid]: https://tools.ietf.org/html/rfc4122
[configMap]: http://kubernetes.io/docs/user-guide/configmap/
