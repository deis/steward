# Steward API

Much of the steward API conforms to the [CloudFoundry service broker API](https://docs.cloudfoundry.org/services/api.html). In all cases, the steward API acts as a proxy. Depending on its mode, it proxies API requests to a backing [Helm](https://github.com/kubernetes/helm) (Tiller) server, its own built-in logic, or a backing CF broker.

Steward does not implement 100% of the CF service broker API, but strives to be 100% compatible with the API endpoints that it does implement. See below for details on each API.

## Provisioning

Steward implements, and is 100% compatible with, the synchronous [provisioning API](https://docs.cloudfoundry.org/services/api.html#provisioning). Note, however, that it does not currently implement the `GET /v2/service_instances/:instance_id/last_operation` API endpoint, so is not compatible with asynchronous provisioning.

Also note that all key/value pairs in the `parameters` object, both in the request and response to this API call, must be strings.

## Deprovisioning

Steward does not yet implement the [deprovisioning](https://docs.cloudfoundry.org/services/api.html#deprovisioning) API endpoint.

## Updating a Service Instance

Steward does not yet implement the [updating a service instance](https://docs.cloudfoundry.org/services/api.html#updating_service_instance) API endpoint.

# Binding

Steward implements the [binding API](https://docs.cloudfoundry.org/services/api.html#binding). It uses a few extensions to do Kubernetes-specific tasks.

When steward binds a service to an application, it stores the service's information (i.e. the data that is returned from the broker in the `credentials` field of the JSON response body) in a [ConfigMap](http://kubernetes.io/docs/user-guide/configmap/) and set of [Secret](http://kubernetes.io/docs/user-guide/secrets/)s. It does not return any credentials information in the response body, but instead returns information on where to find those resources.

Therefore, steward requires a `target_namespace` parameter inside the standard `parameters` field. This value tells steward in which namespace to write the resulting ConfigMap and set of Secrets.

The steward broker API proxy does not return standard [binding credentials](https://docs.cloudfoundry.org/services/binding-credentials.html) in its response. Instead, it returns a JSON object that looks like the following:

```json
{
  "config_map_info": <qualifiedName>,
  "secrets_info": [<qualifiedName1>, <qualifiedName2>, ...]
}
```

Note that each `qualifiedName` instance above specifies the name and namespace of a kubernetes resource. For example:

```json
{
  "name": "mySecret1",
  "namespace": "myApp"
}
```

Finally, note that all key/value pairs in the `parameters` object, both in the request and response to this API call, must be strings.
