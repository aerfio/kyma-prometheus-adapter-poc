# kyma-prometheus-adapter-poc

Let's see what needs to be done to introduce custom metrics in Kyma using Prometheus Adapter

# magic spells

You need [ko](https://github.com/google/ko) to apply deploy folder:

```bash
export KO_DOCKER_REPO=your_docker_name
ko apply -f ./deploy
```

then

```bash
helm upgrade --install my-release ./ --namespace aerfio --set prometheus.url=http://monitoring-prometheus.kyma-system.svc --set tls.enable=true
```

The name "my-release" is mandatory for POC, "aerfio" namespace is mandatory too -> it's because of hard coded certs.

Certs are already in values.yaml.
Basic configuration is also already in values.yaml

If you have knative-serving installed, then you need to remove it's apiservices.apiregistration.k8s.io v1beta1.custom.metrics.k8s.io, https://knative.dev/docs/serving/custom-metrics-api/

After all pods are up you should be able to see `request_number_total` in prometheus UI.

In order to see list of all custom metrics run:

```bash
k get --raw "/apis/custom.metrics.k8s.io/v1beta1/"
```

Generate load for that svc using:

```bash
kubectl port-forward svc/sample-metrics-8080 -n testing-monitoring 8080:8080
```

and for example fortio:

```bash
fortio load -a -qps 0 -t 1m -c 800 -r 0.0001 http://localhost:8080/data
```

And finally:

```bash
k get --raw "/apis/custom.metrics.k8s.io/v1beta1/namespaces/testing-monitoring/pods/*/request_number_per_second"
```

to get those metrics. If there were no problems pod should autoscale based on those metrics by now, after step with fortio.
