# krumble

## The idea

Why this tool? Basically I use `kops` to build a Kubernetes cluster on AWS.
Once it is built I use a combination of custom deployments using `kubectl` and
[helm](https://helm.sh/) to install any additional tools I need such as
nginx-ingress, grafana, etc.

Until now I've been using a simple shell script which does do the job but I
thought I could improve it with a nice wrapper.

## The code

It's rubbish. This is my first project in Go and there are lots of things I need
improving if I really want this to work properly but I'm content it didn't
take me very long to get a working version.

## How it works?

You just need a `yaml` config file with all the params you'll be using for kops,
helm and kubectl and then call the program with the config file as an argument.

- The `kops` section will take any argument supported by `kops create cluster`.
- For `helm` it is used _as is_.
- `kubectl` you can use either URL or local file.

Until now I have *ONLY* tested this tool on a AWS VPC configured with two
subnets, one public and one private.

```yaml
global:
  cluster_name: &cluster_name my_cluster.mydomain.com
  domain: &domain mydomain.com
  bucket: &bucket s3://kops-my_cluster.mydomain.com
  environment: &environment dev
  provider: aws
  aws:
    region: eu-west-1
    availabilityZones:
      - eu-west-1a
      - eu-west-1b
      - eu-west-1c
    vpc_id:
      filters:
        key: "tag:Name"
        value: "dev"
    subnets:
      filters:
        key: "tag:Name"
        value: "dev_private*"
    utility-subnets:
      filters:
        key: "tag:Name"
        value: "dev_public*"

kops:
  snippets:
    cluster: cluster.conf.d
    node: node.conf.d
    master: master.conf.d
  name: *cluster_name
  state: *bucket
  node-count: 3
  node-size: "m5.large"
  master-size: "m5.large"
  master-count: 1

kubectl:
  - name: heapster
    url: https://raw.githubusercontent.com/kubernetes/heapster/master/deploy/kube-config/influxdb/heapster.yaml
    namespace: monitoring
  - name: kube-dashboard
    url: https://raw.githubusercontent.com/kubernetes/kops/master/addons/kubernetes-dashboard/v1.8.3.yaml
    namespace: default

# Possible options with take the same config parameters
#
# pre_exec:   -- runs firt thing
# post_exec:  -- runs last
# exec:       -- runs in the middle, currently after helm
exec:
  # sample running shell script
  - command: monitoring/install.sh
    env:
      - name: AWS_ZONE
        value: eu-west-1
  # another command but change directory before ejecting
  - command: ./my_command.sh
    rundir: /home/randomuser
    env:
      - name: CLUSTER_NAME
        value: *cluster_name

helm:
  repositories:
    - name: stable
      url: https://kubernetes-charts.storage.googleapis.com
    - name: incubator
      url: http://storage.googleapis.com/kubernetes-charts-incubator

  helmDefaults:
    wait: true
    timeout: 600

# if you use env variables you MUST escape the line with '
  releases:
    - name: nginx-ingress
      namespace: ingress
      chart: stable/nginx-ingress
      values:
      - '{{ requiredEnv "PWD" }}/nginx-ingress/values.yaml'
```

## Snippets

Something I find very useful is to be able to change some configurations on `kops`
at the time the cluster is coming up. Let's say for example you want to use Spot
instances. According to the
[documentation](https://github.com/kubernetes/kops/blob/master/docs/instance_groups.md#converting-an-instance-group-to-use-spot-instances)
you need editing the intance groups and add `maxPrice`. What I do instead is I
create an additional config (snippet) into `nodes.conf.d` and the
`krumble` process run will merge this config in before bringing up the cluster.

