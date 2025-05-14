Some example queries for opensearch.

From the top-left nav-bar, select dev tools. From there you can queries against the `ss4o_logs-default-namespace` table

```
GET ss4o_logs-default-namespace/_search
{
  "size": 10,
  "sort": [
    { "@timestamp": "desc" }
  ]
}
```

By namespace:
```
GET ss4o_logs-default-namespace/_search
{
  "query": {
    "match": {
      "resource.k8s.namespace.name": "default"
    }
  },
  "sort": [
    { "@timestamp": "desc" }
  ],
  "size": 4
}
```