# Steward API

Much of the steward API conforms to the [CloudFoundry service broker API](https://docs.cloudfoundry.org/services/api.html). In all cases, the steward API acts as a proxy. Depending on its mode, it proxies API requests to a backing [Helm](https://github.com/kubernetes/helm) (Tiller) server, its own built-in logic, or a backing CF broker (_NOTE: currently, only backend CF broker support is implemented_).

Steward does not implement 100% of the CF service broker API, but strives to be as compatible as possible with the API endpoints that it does implement. Since the proxy also interacts with Kubernetes primitives, 100% API compatibility is sometimes not possible.

See below for details on each API.

# Provisioning

Steward implements the synchronous [provisioning API](https://docs.cloudfoundry.org/services/api.html#provisioning). However, note the following restrictions and changes between Steward's implementation and the that according to the published CloudFoundry documentation:

- It does not currently implement the `GET /v2/service_instances/:instance_id/last_operation` API endpoint, so is not compatible with asynchronous provisioning.
- All key/value pairs in the `parameters` object, both in the request and response to this API call, must be strings.

# Deprovisioning

Steward implements, and is 100% compatible with, the synchronous [deprovisioning API](https://docs.cloudfoundry.org/services/api.html#deprovisioning). However, note the following restrictions and chnages between Steward's implementation and that according to the published CloudFoundry documentation:

- It does not currently implement the `GET /v2/service_instances/:instance_id/last_operation` API endpoint, so is not compatible with asynchronous deprovisioning.

# Updating a Service Instance

Steward does not yet implement the [updating a service instance](https://docs.cloudfoundry.org/services/api.html#updating_service_instance) API endpoint.

# Binding

Steward implements the [binding API](https://docs.cloudfoundry.org/services/api.html#binding). It uses a few extensions to do Kubernetes-specific tasks.

When steward binds a service to an application, it stores the service's information (i.e. the data that is returned from the broker in the `credentials` field of the JSON response body) in a [ConfigMap](http://kubernetes.io/docs/user-guide/configmap/) and returns _only_ the information on where to find the resulting ConfigMap. Note that it does not return any credentials information in the response body.

Steward requires that the application pass information to it indicating where to store the ConfigMap. This information must be passed in two key/value pairs in the standard `parameters` field of the request body:

- `target_namespace` - which namespace to write the resulting ConfigMap
- `target_name` - the name of the resuting ConfigMap

On success, steward will simply return the information on where it stored the ConfigMap:

```json
{
  "credentias": {
    "target_namespace": "namespace-that-was-given",
    "target_name": "name-that-was-given"
  }
}
```

A few final notes:

- All _other_ key/value pairs in the request and response body's `parameters` object must be strings
- All data contained in the resulting `ConfigMap` will be base64 encoded with [Go's ``(encoding/base64).StdEncoding` encoder](https://godoc.org/encoding/base64#pkg-variables)

# Unbinding

Steward implements the [unbinding API](https://docs.cloudfoundry.org/services/api.html#unbinding). It, however, requires two extra query string parameters to ensure it can clean up all Kubernetes resources it created during the bind API call:

- `target_namespace` - the namespace that was passed in the `parameters` object in the associated bind request
- `target_name` - the name that was passed in the `parameters` object in the associated bind request

Upon successful completion of an unbind API call, those resources will have been deleted.

As with the standard CF broker API for unbinding, be sure to call this API with the same `instance_id`, `binding_id`, `service_id` and `plan_id` paramters, as all of these parameters - which were used to name and create `ConfigMap`s and `Secret`s in the above bind API call - are used to locate and delete the created resources.
