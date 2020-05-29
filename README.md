## Kenza

*Kenza* is an open source cloud-native system (moving from `Docker Swarm` to `Kubernetes` in 2020) for Machine Learning Continuous Integration and Delivery you can run in one command. It leverages containers and the cloud to provide basic mechanisms for building, training and deploying Machine Learning models on [AWS SageMaker](https://aws.amazon.com/sagemaker/).

It makes it easy to run training, batch prediction, hyperparameter tuning and deployment jobs on regular intervals or specific points in time e.g. running batch predictions every week or automatically redeploying a model in a QA or production environment every morning.

On top of its traditional "pipeline" features, the *Kenza* web app helps identify training (or other) jobs that are performing better e.g. have a better _accuracy_ or _precision_ score.

## Installation

Download the binary from the latest GitHub release:

```sh
# Linux
curl -L https://github.com/kenza-ai/kenza/releases/download/v0.0.1-alpha/kenza-linux-amd64 -o kenza
```

```sh
# macOS
curl -L https://github.com/kenza-ai/kenza/releases/download/v0.0.1-alpha/kenza-darwin-amd64 -o kenza
```

Move it under a PATH directory, we prefer `/usr/local/bin`:
```sh
chmod +x kenza
sudo mv ./kenza /usr/local/bin/kenza
```

Ensure you are on the expected version:

```sh
kenza info
```

You should see output similar to the following:

```sh
Kenza info

Version: v0.0.32
Built:   2019-12-09T18:57:05Z
Commit:  099415b5087d919d086b383da73afe1b99bf5k0a
```

## Getting Started

#### Starting Kenza
To start (or restart *Kenza*) run:

```sh
kenza start
```

> **Note:** The first run might take longer than subsequent runs due to the *Docker* images downloading for the first time.

> **Important:** The directory from which the `kenza` commands are run from is significant. `kenza start` creates a `kenza` directory in the directory the command was run from. If you run the command again in a different directory, a new `kenza` directory will be created there, essentially a separate `kenza` installation. 

After Kenza has started, it will attempt to navigate you to `http://localhost/#/signup` to create an account and get you started.

#### Checking current service status

You can check the status of *Kenza* and its services with:

```sh
kenza status
```

> If the output feels familiar, it's because *Kenza* is deployed as a *Docker stack*. Running `docker stack ps kenza` would generate the same output.


### Scaling down/up

Kenza runs "one job per worker"; workers are ephemeral in nature and only handle one job before shutting down. To run more than one jobs in parallel, simply add more workers:

```sh
kenza scale worker=5
```

> **Note:** Kenza workers do not need nearly as many resources as one may think (due to the nature of ML jobs) because the actual training takes place on the cloud. Kenza workers only clone the repos, prepare the job commands to be run and report on the status of the jobs as they progress through their lifetime.

### Cleaning up

You can stop Kenza _without any data loss_ with:

```sh
kenza stop
```


### Updating Kenza

To update to the latest available version, run:

```sh
kenza update
```
> **Note:** Currently, this only updates the Kenza executable, future work will stop Kenza, apply all necessary changes and restart the system to ensure all services are brought up to their latest versions, migrations are performed etc. For now, please run `kenza stop` before updating.

## Running Kenza on the Cloud

### Provisioning resources

*Kenza* leverages [containers](https://docs.docker.com/machine/overview/) (currently orchestrated with _Docker Swarm_, moving to _Kubernetes_ in 2020) to run on the cloud. Before starting *Kenza*, the required resources (manager server(s) / instances, security groups etc) need to be provisioned first.

Ensure your local *AWS* access levels (the profile or role you will be using when running `kenza provision` commands) meet the [*IAM* policy requirements for deploying a *Docker Machine*](https://github.com/docker/machine/issues/1655#issuecomment-409407523).

To provision a machine with the [default values](https://docs.docker.com/machine/drivers/aws/#options) on *AWS*, run:

```sh
kenza provision --driver amazonec2 --amazonec2-iam-instance-profile your-sagemaker-aware-intance-profile kenza-machine-1
```

Any other options you pass will be honored; all options are passed as-is to the corresponding `docker-machine` command. One would pass additional options to use a pre-existing _VPC_ or _Security Group_ to limit access to the instance to a specific office IP range for example. The full list of options available can be found [here](https://docs.docker.com/machine/drivers/aws/#options).

You can use any name for the `Docker Machine` (_kenza-machine-1_ in the example above) but the only `driver` supported for now is `"amazonec2"`.

> **Note**: It is highly recommended that the role assigned to the Kenza manager instance follows the [Principle of Least Privilege](https://docs.aws.amazon.com/IAM/latest/UserGuide/best-practices.html#grant-least-privilege) and only provides access to the services and resources that will actually be needed. To identify the exact permissions needed for your use cases use [this AWS reference](https://docs.aws.amazon.com/sagemaker/latest/dg/sagemaker-roles.html)
specific to _SageMaker_. If unsure, *AWS* has been aggressively adding [tools](https://aws.amazon.com/blogs/security/tag/access-advisor/) to make control of roles' more manageable. There are also open-source _Least Privilege Policy_ generators like _Saleforce's_ [Policy Centry](https://github.com/salesforce/policy_sentry/) you can use to ensure permissions are only as elevated as needed.


Run `docker machine ls` to verify the machine you just created is available.

You can also check the [EC2 Dashboard](https://console.aws.amazon.com/ec2/home) on your *AWS* account for the various resources created (e.g. an instance and a key pair matching the "name" parameter provided earlier to the `provision` command, the "docker-machine" security group and others).

To deploy *Kenza* on the newly created resources, we first need to ensure the `Docker Machine` we just created is [*active*](https://docs.docker.com/machine/reference/active/). To do this, run (substituting if needed `kenza-machine-1` with the name you provided to the `provision` command):

```sh
eval $(kenza env kenza-machine-1)
```

Verify `Docker` is now actually "forwarding all calls" to the remote machine:

```sh
docker-machine active
```

With the machine set up, all *Kenza* commands will now be run against the newly deployed infrastructure, not your local machine.

To start _Kenza_ on EC2, simply run (substituting if needed `kenza-machine-1` with the name you provided to the `provision` command):

```sh
kenza start --name kenza-machine-1 --github-secret webhooks-secret --apikey a-randomly-generated-key
```

After _Kenza_ starts, it will open your default browser to the URL / Public IP of the machine where the _Kenza_ web app can be reached.

Once launched, you can [associate your instance with a static IP or a domain name](https://docs.aws.amazon.com/Route53/latest/DeveloperGuide/routing-to-ec2-instance.html).

## Troubleshooting

#### Getting detailed service execution details

You can observe detailed log output for a service with:

```sh
kenza logs service_name
```

You can stop an individual service with:

```sh
kenza stop service_name
```

Valid service names:
- db
- api
- web
- worker
- pubsub
- progress
- scheduler

#### Restarting Kenza

Restarting _Kenza_ or _Docker_ can sometimes help when `Docker Swarm` seems to be "stuck".

For any other issue, please [raise an issue](https://github.com/Kenza-AI/kenza/issues/new).

## Component Overview

Kenza is composed of the following components:

- **API** - Service called by all other services, including the cli, to read / mutate Kenza related data (projects, jobs, schedules etc). 

> Note for contributors: API is the only service with direct access / dependency to the Kenza data store(s). All other services *MUST* go through the API.

- **Web** - *React.js* web application, the *Kenza UI*.

- **Worker** - Worker nodes, the container tasks actually running the jobs. Workers are ephemeral and strictly process one job and one job only before shutting themselves down.

- **Progress** - Listens for job updates published by the worker nodes and propagates them to the *API*.

- **Scheduler** - Listens for job arrivals (on-demand, webhooks and scheduled jobs) and schedules them accordingly to be picked up by workers for processing.

- **PubSub** - *RabbitMQ* exchanges and queues, used for async comms among services.

- **DB** - The *kenza* data store (currently *Postgres*). It can be a *Postgres* container (default option, provided by *Kenza* as a container) or an external resource e.g. an *AWS RDS*, *Heroku* or on-prem installation.

- **CLI** - The *Kenza* command line utility. Think *kubectl, systemctl*.

*Kenza* currently supports *Docker Swarm* environments. Support for *Kubernetes* is being added in 2020.


## Kenza UI

The Kenza web application is a [ReactJS](https://github.com/facebook/react) / [Redux](https://react-redux.js.org) Single Page Application (SPA). You can use the standard tooling e.g. React Tools ([Chrome](https://chrome.google.com/webstore/detail/react-developer-tools/fmkadmapgofadopljbjfkapdkoienihi?hl=en), [Firefox](https://addons.mozilla.org/en-US/firefox/addon/react-devtools/)) to troubleshoot  / report issues with specific browsers.

## Note on tests (or lack thereof)
_Kenza_ was originally built as a typical cloud based pipeline on _AWS_; tests will be being moved as they are getting adapted to the container-based world, probably starting with the ones that are the least impacted by the move e.g. the UI / web app.

