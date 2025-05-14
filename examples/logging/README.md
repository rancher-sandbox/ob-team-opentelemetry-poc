# install opensearch

```sh
helm repo add opensearch https://opensearch-project.github.io/helm-charts/
```

```sh
helm repo update
```

```sh
helm install opensearch opensearch/opensearch -f ./examples/logging/opensearch.yaml
```

```sh
helm install opensearch-dashboards opensearch/opensearch-dashboards
```

```sh
kubectl port-forward deploy/opensearch-dashboards 5601:5601