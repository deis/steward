# Cloud Foundry Broker Mode

Steward can integrate with any standard Cloud Foundry service broker (broker hereafter) to provide all of the services
the broker provides to the Kubernetes cluster. When a steward is started with such a configuration,
it will immediately make a request to the broker's catalog (i.e. `GET /v2/catalog`), convert the
catalog results to a `ServiceProvider` and set of `ServicePlan`s, and write the results to the
appropriate 3PR.

Configure steward to run in Cloud Foundry mode by setting the `STEWARD_MODE` environment variable to `cf`. Then, configure its behavior with the following environment variables:

- `CF_BROKER_SCHEME` - the scheme (`http` or `https`) by which to access the backend broker with which steward should communicate
- `CF_BROKER_HOSTNAME` - the IP or DNS name of the broker
- `CF_BROKER_PORT` - the port of the broker
- `CF_BROKER_USERNAME` - the username steward should use to authenticate with the broker
- `CF_BROKER_PASSWORD` - the password steward should use to authenticate with the broker
- `HTTP_REQUEST_TIMEOUT_SEC` - the timeout after which steward should fail a request to the broker for any individual request

Note that steward will, by default, watch for claims submitted in the `default` namespaces. However, it can be configured to watch in other namespaces as well. To do so, set the `WATCH_NAMESPACES` environment variable to a comma-separated list of namespaces.

For example, to listen in the `steward` and `default` namespaces, set `WATCH_NAMESPACES="default,steward"`
