# SCHEDULER

Scheduler is an API which abstracts the concepts of Kubernetes Jobs and CronJobs to easily schedule ad-hoc jobs in highly dynamic environments.

- [Installation](README.md#installation)
    - [Prerequisites](README.md#prerequisites)
- [Getting Started](README.md#getting-started)
    - [Baking Image](README.md#baking-image)
    - [Deploying Infrastructure](README.md#deploying-infrastructure)
    - [Deploying API](README.md#deploying-api)
- [Concepts](README.md#concepts)
- [Local Development](README.md#local-development)
    - [Infrastructure Tests](README.md#infrastructure-tests)
    - [API Tests](README.md#api-tests)

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

#### 1. Make sure your packer is installed by running:
```bash
$ packer version
Packer v1.4.0
```

There is a [Cloudformation template](infra/cloudformation/templates/packer.yml) to make things easier which creates:
- Packer IAM User
- Packer IAM Policy

#### 2. To bootstrap Packer IAM credentials through the Cloudformation template, run:
```bash
$ make packer-creds CLOUDFORMATION_STACK_NAME=packer-creds
```

#### 3. Get the AWS Access Key and Secret Key from the generated template by running:
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

#### 4. With Packer installed and credentials set up, you can bake the AMI using the credentials from the previous output:
```bash
$ make build-ami PACKER_AWS_ACCESS_KEY="YOURACCESSKEY" PACKER_AWS_SECRET_KEY="YOURSECRETKEY"
```

#### 5. A [python script](scripts/latest_base_ami.py) was created to check if your AMI was created successfully and always show the latest of them:
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

#### 1. Create the cluster from the Cloudformation template:
```bash
$ make build-cluster AWS_KEYPAIR_NAME=mykey PACKER_BASE_AMI_ID=$(make get-ami)
```

This template supports multiple parameters like *Master Instance Size* and *Kubeadm Token*, but the defaults should support most deployments. If you want to customize it, deploy the [Cloudformation template](infra/cloudformation/templates/stack.yml) manually through the console or using [aws cloudformation deploy](https://docs.aws.amazon.com/cli/latest/reference/cloudformation/deploy/index.html).

The average time to deploy the whole infrastructure is **13 minutes**, made of:
- Cloudformation ≈ **7 min**
- Master & Worker node ≈ **5 min**
- Tiller & cluster-autoscaler ≈ **1 min**

*OPTIONAL: If you want to clean the Cloudformation created from **make build-cluster**:*
```bash
$ make clean-cf
```

OBS: To deploy the API you must have your kubectl configured. This Cloudformation only bootstraps the cluster, so it requires you to create your kubeconfig. You may need to connect via SSH to your Master Node to run the next commands, so check the Security Group Inbound rules.

<br>

### [Deploying API](#deploying-api)

The scheduler is deployed using a Helm Chart which works on any Kubernetes cluster with Tiller successfully deployed. The chart folder is located at ```deployments/chart/scheduler/```. There is a [documentation](deployments/chart/scheduler/README.md) of Chart Values if you need to customize them.

#### AWS
Balancers are automatically created by AWS if the cluster is configured properly, and may take some minutes to attach the nodes and expose the service.


1. Install the scheduler and exposing it through a LoadBalancer:
```bash
$ helm upgrade --install scheduler deployments/chart/scheduler --namespace default --set service.type=LoadBalancer
```

2. Get Load Balancer address by getting the Kubernetes service: 
```
$ kubectl get svc/scheduler -o wide
NAME        TYPE           CLUSTER-IP       EXTERNAL-IP                                                        PORT(S)
scheduler   LoadBalancer   100.99.98.97     yourLoadBalancerInAws.us-west-2.elb.amazonaws.com   80:31196/TCP   5m
```

3. Test the Health Check endpoint:
```bash
# For AWS Load Balancer
$ curl http://yourLoadBalancerInAws.us-west-2.elb.amazonaws.com:8080/healthz
{"message":"ok","status":200}
```

#### Minikube
1. Install the scheduler without exposing it:
```bash
$ helm upgrade --install scheduler deployments/chart/scheduler --namespace default
```

2. Expose the scheduler service using Kubernetes port-forward via kubectl:
```bash
$ kubectl port-forward svc/scheduler 8080:8080 --namespace default
```

3. Test the Health Check endpoint:
```bash
# For Minikube port-forward
$ curl http://localhost:8080/healthz
{"message":"ok","status":200}
```

#### Clean API deployment
*OPTIONAL: With your kubectl configured, you can delete the Scheduler from your cluster:*
```bash
$ helm delete --purge scheduler
```

<br>

## [Concepts](#concepts)

-----

### [Swagger API reference - WIP](https://app.swaggerhub.com/apis/tioxy/scheduler/1.0.0)

### SimpleJob
A SimpleJob is an abstraction on top of CronJobs and Jobs from Kubernetes. The SimpleJob can be separated in:
- **/jobs:** Run once and immutable
- **/scheduled:** Run multiple times(cron syntax) and mutable 

```go
type SimpleJob struct {
	Name       string         `json:"name"`
	Namespace  string         `json:"namespace"`
	Containers []v1.Container `json:"containers"`
	MaxRetries int32          `json:"maxRetries"`
	Cron       string         `json:"cron,omitempty"`
}

// v1.Container is the Container object used in Pods from Kubernetes API
```

An example SimpleJob to calculate first 2000 numbers of pi:
```json
{
    "name": "pi",
    "namespace": "default",
    "maxRetries": 4,
    "containers": [
        {
            "name": "pi",
            "image": "perl",
            "command": ["perl", "-Mbignum=bpi", "-wle", "print bpi(2000)"]
        }
    ]
}
```

Now the same SimpleJob, but scheduled to run every day at midnight:
```json
{
    "name": "pi",
    "namespace": "default",
    "maxRetries": 4,
    "cron": "0 0 * * *",
    "containers": [
        {
            "name": "pi",
            "image": "perl",
            "command": ["perl", "-Mbignum=bpi", "-wle", "print bpi(2000)"]
        }
    ]
}
```

OBS: The **"cron"** key inside your SimpleJob will only be used if you are interacting with **/scheduled** endpoints, otherwise it doesn't matter.


<br>

## [Local Development](#local-development)

-----

## [Infrastructure Tests](#infrastructure-tests)

**Ansible**

[Molecule](https://molecule.readthedocs.io/en/stable/) is used to test roles.

- Test single role
```bash
$ make test-role ANSIBLE_ROLE=roleName
```

- Test all roles 
```bash
$ make test-roles
```

**Cloudformation**

[Taskcat](https://github.com/aws-quickstart/taskcat) is used to test stacks.

- Test all stacks
```bash
$ make test-cf
```

## [API Tests](#api-tests)

- Unit test
```
make test
```

- E2E test with endpoint
```
make test-e2e SCHEDULER_ENDPOINT=http://localhost:8080
```
