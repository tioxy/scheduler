# SCHEDULER

Scheduler is an API which abstracts the concepts of Kubernetes Jobs and CronJobs to easily schedule ad-hoc jobs in highly dynamic environments.

- [Installation](README.md#installation)
    - [Prerequisites](README.md#prerequisites)
- [Getting Started](README.md#getting-started)
    - [Baking Image](README.md#baking-image)
    - [Deploying Infrastructure](README.md#deploying-infrastructure)
    - [Deploying API](README.md#deploying-api)
- Local Development
    - Testing
        - Infrastructure
        - API

## [Installation](#installation)
-----

### [Prerequisites](#prerequisites)

These dependencies below are used to test, provision and deploy API and the infrastructure:
- [Docker](https://docs.docker.com/install/) >= 18.0.0
- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/) configured
- [Helm](https://helm.sh/docs/using_helm/#installing-helm) >= 2.13.1
- [Packer](https://www.packer.io/intro/getting-started/install.html#precompiled-binaries) >= 1.4.0
- [Golang](https://golang.org/doc/install#download) >= 1.12.0
- [Python](https://github.com/pyenv/pyenv) >= 3.7.0
- [pip](https://pip.pypa.io/en/stable/installing/) >= 19.0.0
- [AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/installing.html) configured

[Poetry](https://poetry.eustace.io/docs/) is used to manage all Python dependencies through the project. Install it by running:
```bash
pip install --upgrade poetry
```

<br>

## [Getting Started](#getting-started)

-----

Make sure that your Poetry environment is up and running:
```bash
$ poetry up
```

Import your Public Key to AWS to create an EC2 Keypair which will be used through all instances in the cluster:
```bash
$ make import-keypair AWS_KEYPAIR_NAME=mykey PUBLIC_KEY_FILE=/path/to/key.pub
```

If you already have a Kubernetes cluster running, you can jump to [Deploying API ](README.md#deploying-api)

<br>

### [Baking Image](#baking-image)
In this step, we will bake the *Kubernetes Base* AMI which will be used for Master and Workers nodes in the cluster.

Make sure your packer is installed by running:
```bash
$ packer version
Packer v1.4.0
```

There is a [Cloudformation template](infra/cloudformation/templates/packer.yml) to make things easier which creates:
- Packer IAM User
- Packer IAM Policy

To bootstrap Packer IAM credentials through the Cloudformation template, run:
```bash
$ make packer-creds CLOUDFORMATION_STACK_NAME=packer-creds
```

Get the AWS Access Key and Secret Key from the generated template by running:
```bash
$ make get-packer-creds CLOUDFORMATION_STACK_NAME=packer-creds
[
    {
        "Param": "PACKER_AWS_SECRET_KEY",
        "Value": "YOURSECRETKEY"
    },
    {
        "Param": "PACKER_AWS_ACCESS_KEY",
        "Value": "YOURACCESSKEY"
    }
]
```

With Packer installed and credentials set up, you can bake the AMI using the credentials from the previous output:
```bash
$ make build-ami PACKER_AWS_ACCESS_KEY="YOURACCESSKEY" PACKER_AWS_SECRET_KEY="YOURSECRETKEY"
```

A [python script](scripts/latest_base_ami.py) was created to check if your AMI was created successfully and always show the latest of them:
```bash
$ make get-ami
ami-y0uram1idw1llb3h3re
```

*OPTIONAL: With your AMI created, you can remove Packer Cloudformation stack:*
```bash
$ make clean-cf CLOUDFORMATION_STACK_NAME=packer-creds
```

<br>

### [Deploying Infrastructure](#deploying-infrastructure)

There is a [Cloudformation template](infra/cloudformation/templates/stack.yml) to bootstrap the entire cluster composed of:
- VPC (Subnets, Route Tables, Internet Gateway, NAT Gateway)
- IAM Role & Instance Profile for Nodes
- Master Node (helm and cluster-autoscaler configured automatically)
- Auto Scaling Group for Worker Nodes

```bash
$ make build-cluster AWS_KEYPAIR_NAME=mykey PACKER_BASE_AMI_ID=$(make get-ami)
```

This template support multiple parameters like *Master Instance Size* and *Kubeadm Token*, but the defaults should support most deployments. If you want to customize it, deploy the [Cloudformation template](infra/cloudformation/templates/stack.yml) manually through the console or using [aws cloudformation deploy](https://docs.aws.amazon.com/cli/latest/reference/cloudformation/deploy/index.html).

The average time to deploy the whole infrastructure is 13 minutes, made of:
- Cloudformation ≈ 7 min
- Master & Worker node ≈ 5 min
- Tiller & cluster-autoscaler ≈ 1 min

*OPTIONAL: If you want to clean the Cloudformation created from **make build-cluster**:*
```bash
$ make clean-cf
```

<br>

### [Deploying API](#deploying-api)

The scheduler is deployed using a Helm Chart which works on any Kubernetes cluster with Tiller successfully deployed. The chart folder is located at ```deployments/chart/scheduler/```. There is a [documentation](deployments/chart/scheduler/README.md) of Chart Values if you need to customize them.

Installing the scheduler and exposing it through a LoadBalancer *(Balancers are automatically created by your Cloud Provider if the cluster is configured properly)*:
```bash
$ helm upgrade --install scheduler deployments/chart/scheduler --namespace default --set service.type=LoadBalancer
```

If you want a raw installation to check if everything is running (recommended for Minikube installation):
```bash
$ helm upgrade --install scheduler deployments/chart/scheduler --namespace default
```

Testing locally with your Minikube requires Service port-forward using kubectl:
```bash
$ kubectl port-forward svc/scheduler 8080:8080 --namespace default

$ curl http://localhost:8080/healthz
{"message":"ok","status":200}
```

*OPTIONAL: With your kubectl configured, you can delete the Scheduler from your cluster:*
```bash
$ helm delete --purge scheduler
```

<br>


## [Local Development](#local-development)

-----

WIP
