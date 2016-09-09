# Cloud Foundry Broker Mode

Steward can integrate with any standard CloudFoundry service broker to provide all of the services
the broker provides to the Kubernetes cluster. When a steward is started with such a configuration,
it will immediately make a request to the broker's catalog (i.e. `GET /v2/catalog`), convert the
catalog results to a `ServiceProvider` and set of `ServicePlan`s, and write the results to the
appropriate 3PR.

Configure steward to run in helm mode by setting the `STEWARD_MODE` environment variable to `cf`. Then, configure its behavior with the following environment variables:

- `CF_HOSTNAME` - the hostname or IP of the CloudFoundry service broker with which Steward should communicate
- `CF_USERNAME` - the username Steward should use to authenticate with the backend CloudFoundry service broker
- `CF_PASSWORD` - the password Steward should use to authenticate with the backend CloudFoundry service broker
