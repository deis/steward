# Integration Testing

This document describes how to run integration tests against Steward.

# Cloud Foundry Service Broker

This suite of tests is completely manual for now. Automated tests are forthcoming (see https://github.com/deis/steward/issues/20). See below for testing instructions.

1. Make sure you have a `steward` namespace:

  ```console
  kubectl create -f manifests/steward-ns.yaml
  ```

2. Create a service broker RC and service:

  ```console
  kubectl create -f manifests/sample-broker/rc.yaml --namespace=steward
  kubectl create -f manifests/sample-broker/svc.yaml --namespace=steward
  ```

3. Create the `ServiceCatalogEntry` [third party resource](https://github.com/kubernetes/kubernetes/blob/release-1.3/docs/design/extending-api.md), which describes the data type that holds available services:

  ```console
  kubectl create -f manifests/3pr/service-catalog-entry.yaml
  ```

  Note: do not create any other third party resoruces (including the one in `manifests/3pr/service-plan-claim.yaml`), or you'll hit https://github.com/deis/steward/issues/17
4. Start steward:

  ```console
  kubectl create -f manifests/steward-rc.yaml --namespace=steward
  ```

5. View logs (substitute your steward pod's name in for `<pod_name>`):

  ```console
  kubectl logs -f <pod_name> --namespace=steward
  ```

  They should look similar to the following:

  ```console
  2016-07-28 23:49:52.494438 I | steward version dev started
  2016-07-28 23:49:52.495583 I | starting in CloudFoundry mode with hostname 10.171.244.158:9292 and username admin
  2016-07-28 23:49:52.495608 I | CF client making request to http://admin:password@10.171.244.158:9292/v2/catalog
  2016-07-28 23:49:52.527358 I | making request to https://10.171.240.1:443/apis/steward.deis.com/v1/namespaces/steward/servicecatalogentries
  2016-07-28 23:49:52.561523 I | published 1 entries into the catalog:
  2016-07-28 23:49:52.561560 I | github-repo
  ```

  Note: if steward has already run and published its service catalog, it will try to publish a `ServiceCatalogEntry` that already exists, print an error and crash:

  ```console
  2016-07-29 00:05:37.049442 I | steward version dev started
  2016-07-29 00:05:37.052923 I | starting in CloudFoundry mode with hostname 10.171.250.49:9292 and username admin
  2016-07-29 00:05:37.052989 I | CF client making request to http://admin:password@10.171.250.49:9292/v2/catalog
  2016-07-29 00:05:37.065077 I | making request to https://10.171.240.1:443/apis/steward.deis.com/v1/namespaces/steward/servicecatalogentries
  2016-07-29 00:05:37.096350 I | error publishing the cloud foundry service catalog (duplicate service catalog entry: github-repo-public)
  ```

  To resolve this issue, delete the `ServiceCatalogEntry`:

  ```console
  kubectl delete servicecatalogentry github-repo-public
  ```

  and restart the server. https://github.com/deis/steward/issues/14 will remove the need for this extra step.

6. View the newly-created `ServiceCatalogEntry` that steward just published:

  ```console
  kubectl get servicecatalogentry github-repo-public --namespace=steward
  ```
