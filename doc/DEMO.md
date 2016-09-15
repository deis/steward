# Demo-ing Steward

This document outlines how to show a demo for the three modes of Steward.

# Cloud Foundry Mode

A Steward running in Cloud Foundry (cf hereafter) mode fulfills ServicePlanClaim actions by issuing standard [Cloud Foundry Broker API][cfbrokerapi] calls.

## Setup

To run a demo in Cloud Foundry mode, you need the following:

1. A Cloud Foundry broker running
2. A Steward running and configured to talk to the broker

### Broker Setup

First, broker setup and run instructions vary between brokers, but we've provided a [sample broker][sample-broker]. This broker is easily installable and requires no configuration, but does nothing useful. It's best used for testing the Steward-to-broker communication, but for real demos we recommend using a broker that does something more useful.

After you've successfully set up and run a Steward with this sample broker, we recommend using the Cloud Foundry community's [AWS RDS service broker](https://github.com/cloudfoundry-community/pe-rds-broker). Configure that broker according to its README, the run it in your Kubernetes cluster (it can be installed with Deis Workflow) and point Steward to it.

### Configuring and Running Steward

Once a broker is up and running, configuring and starting Steward is fairly simple. See the [pre-built cf mode manifest](https://github.com/deis/steward/blob/master/manifests/steward-template-cf.yaml) that contains all the configuration necessary to run Steward. Simply change each config field in the `spec.template.spec.containers[0].env` field and run `kubectl create -f manifests/steward-template-cf.yaml` to install.

See the [cf mode documentation][cf-mode] for more details.

After Steward starts, it will query the backend cf service broker's catalog, convert each catalog entry into a list of Kubernetes [Third Party Resource][3pr]s called `ServiceCatalogEntry` and write each into the `steward` namespace.

See the below "Running a Demo" section for instructions on how to inspect the catalog and make claims on the services therein.

# Helm Mode

A Steward running in helm/tiller (helm hereafter) mode fulfills ServicePlanClaim actions by issuing network RPC calls to a backing tiller server running in the same Kubernetes cluster.

## Setup

To run a demo in helm mode, you need the following:

1. A publically-accessible chart that conforms to Steward's standard
2. A running tiller server
3. A running Steward configured to communicate with the backing tiller server

### Installing a Publically-Accessible Chart

A publically-accessible chart is simply a gzipped-tarball of a standard Helm chart accessible via HTTP at any stable address. Steward expects charts to have two specific features that are not required in any helm chart:

- A list called `stewardConfigMaps` in the `values.yaml` file. See [documentation][helm-mode] for more information
- One or more Kubernetes ConfigMaps installed

See [documentation on chart repositories][chart-repo] for more information on how to create a helm chart, and [helm mode documentation][helm-mode] for more information on why Steward needs the special additions to the helm chart, and how it uses them.

### Installing a Tiller Server

Installing a Tiller server is fairly simple - simply download the `helm` CLI and run:

```console
helm init
```

The Tiller server will be installed in the `kube-system` namespace. See [helm initialization instructions][helm-init] for more details on getting the `helm` CLI and installing the tiller server.

However, `helm init` doesn't install a Kubernetes Service to proxy network traffic to the Tiller server (instead, the `helm` CLI tunnels traffic directly to the Tiller pod). We've provided a Service in `manifests/tiller-service.yaml` that you must install before launching a Steward. Install it with:

```console
kubectl install manifests/tiller-service.yaml
```

A few extra notes:

- This service is of type `ClusterIP`, so you'll only be able to access the tiller service from inside the cluster
- This service will install into the `kube-system` namespace, which is the default namespace tiller is installed to

### Configuring and Running Steward

Now that a tiller server is running, we can configure Steward to talk to it.

Once a broker is up and running, configuring and starting Steward is fairly simple. See the [pre-built helm mode manifest](https://github.com/deis/steward/blob/master/manifests/steward-template-helm.yaml) that contains all the configuration necessary to run Steward. Simply change each config field in the `spec.template.spec.containers[0].env` field and run the following to install:

```console
kubectl create -f manifests/steward-template-helm.yaml
```

A few notes on configuration:

- Since you installed `manifests/tiller-service.yaml` in the previous step, set `HELM_TILLER_IP` to the service's DNS name: `tiller.kube-system`
- The standard tiller port is 44134. Set `HELM_TILLER_PORT` to that value
- `HELM_SERVICE_ID`, `HELM_SERVICE_NAME`, `HELM_SERVICE_DESCRIPTION`, `HELM_PLAN_ID`, `HELM_PLAN_NAME` and `HELM_PLAN_DESCRIPTION` are used by Steward to write a single service catalog entry to the Kubernetes catalog.

See the [helm mode documentation][helm-mode] for more details on configuration.

When Steward starts it will create a new service catalog entry in the Kubernetes cluster which will match the given config options.

See the below "Running a Demo" section for instructions on how to inspect the catalog and make claims on the services therein.

# Command Mode

A Steward running in command mode fulfills ServicePlanClaim actions by launching a pod with an operator-specified image.

## Setup

To run a demo in command mode, you need the following:

1. An image that adheres to the standard Steward command mode interface
2. A running Steward configured in command mode with the aforementioned image

### Creating an Image that Adheres to the Command Mode Interface

Before installing Steward, make sure you have an image that conforms to the [command mode interface](https://github.com/deis/steward/blob/master/doc/CMD_MODE.md#interface). One such image is located at `quay.io/deisci/cmd-postgres-broker:devel`. See https://github.com/deis/cmd-postgres-broker for more details.

When claims are submitted, Steward will launch pods using this image, so ensure that your Kubernetes cluster can access it.

See [the "Interface" section](https://github.com/deis/steward/blob/master/doc/CMD_MODE.md#interface) for more details on the command mode interface.

### Configuring and Running Steward

Now that your image is ready, we can configure Steward in command mode to use it.

First, see the [pre-built command mode manifest](https://github.com/deis/steward/blob/master/manifests/steward-template-cmd.yaml) that contains all the configuration necessary to run Steward in command mode. Simply change each config field in the `spec.template.spec.containers[0].env` field and run the following to install:

```console
kubectl create -f manifests/steward-template-cmd.yaml
```

See the [command mode documentation][cmd-mode] for more details on configuration.

After Steward starts, it will write the configuration-specified service catalog entry into the Kubernetes catalog.

See the below "Running a Demo" section for instructions on how to inspect the catalog and make claims.

# Running a demo

Immediately after Steward starts, it will install a set of Kubernetes [Third Party Resource][3pr]s called `ServiceCatalogEntry` into the `steward` namespace. Depending on the mode, it gets the catalog entries from a different place. In cf mode, for example, it gets the catalog by querying the backend cf service broker.

## Inspect the Catalog
After a successful startup, view the catalog by running the following command:

```console
kubectl get servicecatalogentries --namespace=steward
```

You'll see a list of entries. Choose one and run the following command to see details on it:

```console
kubectl get servicecatalogentry $ENTRY_NAME --namespace=steward -o yaml
```

You should see some YAML output. Make note of the `service_info.id` and `plan_info.id` fields, as you'll use them in the service plan claim that you'll install next.

## Make a Claim

After Steward is properly configured and running, interact with it using service plan claims (claims hereafter). Claims are config maps with the following properties:

- A label with key `type` and value `service-plan-claims`
- `data` with the following keys:
  - `service-id` - the ID of a service that's present in the catalog
  - `plan-id` - the ID of a plan for the aforementioned service
  - `claim-id` - an operator-specified identifier for the claim. Currently, these need not follow any format or be unique
  - `action` - one of the actions in the below bulleted list
  - `target-name` - the name of the Secret that Steward will create and to which it will write bound credentials. The Secret will always be created in the same namespace as the claim itself. Note that Steward deletes the secret named here during an unbind operation, so this field should never be changed

Steward watches the `default` namespace for claims to be added or modified. Each addition or modification follows a state machine that defines the one of the following actions you can take:

- `provision` - the creation action of a requestable service
- `bind` - the process of getting credentials to access a provisioned requestable service
- `unbind` - the process of invalidating (or "giving up") credentials that were acquired by a bind operation
- `deprovision` - the deletion action of a requestable service that was previously provisioned
- `create` - the compound action of `provision`, then `bind`
- `delete` - the compound action of `unbind`, then `deprovision`

### Running a Demo

After you've installed Steward, verified that it's correctly written the service catalog, and identified the service ID (`$SERVICE_ID` hereafter) and one of that service's plans (`$PLAN_ID` hereafter), follow the following script to show a demonstration of Steward:

- List the service catalog: `kubectl get servicecatalogentry --namespace=steward`
- Configure the claim (this claim is already configured with the `create` action):
  - Open `manifests/sample-service-plan-claim-cm.yaml`
  - Change the `metadata.name` value to `"steward-demo"`
  - Change the `data.service-id` value to `"$SERVICE_ID"`
  - Change the `data.plan-id` field to `"$PLAN_ID"`
  - Change the `data.target-name` field to `"steward-demo"`
- Submit the claim: `kubectl create -f manifests/sample-service-plan-claim-cm.yaml`
- Inspect the created secret: `kubectl get secret steward-demo -o yaml`
- Delete the provisioned, bound resource:
  - `kubectl edit configmap claim-1`
  - Change the `data.action` field to `"delete"`
- Ensure that the previously-created secret is missing: `kubect get secret steward-demo`

A few extra notes:

- If you intend to demo all 3 Steward modes, run 3 different Steward instances (all in the `steward` namespace), and run the demo as follows:
  - Ensure that each Steward writes a unique service catalog entry (Steward will log an error but still start up properly if it tries to write a catalog entry that already exists)
  - Ensure that you submit 3 different claims. Each should do actions on its corresponding mode
  - Ensure that you are tailing logs (e.g. `kubectl logs -f $STEWARD_POD_NAME --namespace=steward`) on each Steward pod, so you can see the applicable Steward doing work when claims are submitted
- Recall that `create` and `delete` are compound actions (`provision`/`bind` and `unbind`/`deprovision`, respectively). Some audiences may want to see all 4 individual actions (`provision`, `bind`, `unbind`, `deprovision`). If that's the case, simply modify the above demo instructions to:
  - Change the `data.action` field to `"provision"` before submitting the claim
  - Call `kubectl edit configmap claim-1` three times to change the `action` field to `bind`, then `unbind`, then finally `deprovision`
  - Ensure that the appropriate resources are in the appropriate state after each action. For example, the appropriate command should be run to successful completion in command mode.

[sample-broker-installation]: https://github.com/deis/cf-sample-broker#installing-the-sample-broker
[sample-broker]: https://github.com/deis/cf-sample-broker
[cfbrokerapi]: https://docs.Cloud Foundry.org/services/api.html
[helm-init]: https://github.com/kubernetes/helm/blob/master/docs/quickstart.md#initialize-helm-and-install-tiller
[chart-repo]: https://github.com/kubernetes/helm/blob/master/docs/chart_repository.md
[cf-mode]: https://github.com/deis/steward/blob/master/doc/CF_BROKER_MODE.md
[helm-mode]: https://github.com/deis/steward/blob/master/doc/HELM_MODE.md
[cmd-mode]: https://github.com/deis/steward/blob/master/doc/CMD_MODE.md
