# Helm Mode

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
