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

When steward binds a service to an application, it stores the service's information (i.e. the data that is returned from the broker in the `credentials` field of the JSON response body) in a [ConfigMap](http://kubernetes.io/docs/user-guide/configmap/) and set of [Secret](http://kubernetes.io/docs/user-guide/secrets/)s. It does not return any credentials information in the response body, but instead returns information on where to find those resources.

Therefore, steward requires a `target_namespace` parameter inside the standard `parameters` field. This value tells steward in which namespace to write the resulting ConfigMap and set of Secrets.

The steward broker API proxy does not return standard [binding credentials](https://docs.cloudfoundry.org/services/binding-credentials.html) in its response. Instead, it returns a JSON object that looks like the following:

```json
{
  "config_map_info": {"name": "name1", "namespace": "namespace1"},
  "secrets_info": [
    {"name": "name2", "namespace": "namespace2"},
    {"name": "name3", "namespace": "namespace3"}
  ]
}
```

Also, note that all key/value pairs in the `parameters` object, both in the request and response to this API call, must be strings.

Finally, note that all data contained in the resulting `ConfigMap` and `Secret`s will be base64 encoded with [Go's ``(encoding/base64).StdEncoding` encoder](https://godoc.org/encoding/base64#pkg-variables)

# Unbinding

Steward implements the [unbinding API](https://docs.cloudfoundry.org/services/api.html#unbinding). It, however, requires one extra query string parameter, called `target_namespace`, to function properly. This parameter will be used to locate the `ConfigMap` and `Secret` resources that it wrote to Kubernetes as a result of the bind API call (see above). Upon successful completion of an unbind API call, those resources will have been deleted.

As with the standard CF broker API for unbinding, be sure to call this API with the same `instance_id`, `binding_id`, `service_id` and `plan_id` paramters, as all of these parameters - which were used to name and create `ConfigMap`s and `Secret`s in the above bind API call - are used to locate and delete the created resources.
