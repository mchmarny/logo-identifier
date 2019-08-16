# logo-identifier

Simple logo identification service demonstrating use of Cloud SQL and Google Vision API in Cloud Run.

Demo: https://logoider-2gtouos2pq-uc.a.run.app

## Why

Setting up development and deployment pipeline for Cloud Run service backed by project-level authenticated services like Cloud Vision and connection string based authenticated ones like Cloud SQL can be complicated. While logo identification is probably not representative of a real production workload, this simple service does illustrate all of the aspects of local development and service deployment on Cloud Run.

## Pre-requirements

### GCP Project and gcloud SDK

If you don't have one already, start by creating new project and configuring [Google Cloud SDK](https://cloud.google.com/sdk/docs/). Similarly, if you have not done so already, you will have [set up Cloud Run](https://cloud.google.com/run/docs/setup).

## Setup

To setup this service you will need to:

* Configure service dependencies
* Build docker image from the source in this repo
* Deploy and configure service using the previously built image

To start, clone this repo:

```shell
git clone https://github.com/mchmarny/logo-identifier.git
```

And navigate into that directory:

```shell
cd logo-identifier
```

### Configure Dependencies

In this section you will:

* Enable required GCP APIs
* Create Cloud SQL instance and configure it with database and user
* Create a Service Account and configure it with required IAM policies
* Configure OAuth credentials with Cloud Run service callback URL

#### APIs

First, let's ensure all the required GCP APIs are enabled. To do that, run the [bin/setup](./bin/setup) script:

> Note, to keep this readme short, I will be asking you to execute scripts rather than listing here complete commands. You should really review each one of these scripts for content, and, to understand the individual commands so you can use them in the future.


```shell
bin/setup
```

#### DB

Before creating Cloud SQL instance, you'll need to define a couple of passwords: default DB user (`root`), and the application specific user (`logoider-db-user`). To do that we will use the `openssl` command line utility. If for some reason you do not have `openssl` configured, you can just set these value to your own secrets.

```shell
export DB_ROOT_SECRET=$(openssl rand -base64 16)
echo "root user password: ${DB_ROOT_SECRET}"
export DB_USER_SECRET=$(openssl rand -base64 16)
echo "app user password: ${DB_USER_SECRET}"
```

To create `logoider` CLoud SQL instance, run [bin/db](./bin/db) script:

```shell
bin/db
```

One Cloud SQL finished configuring, there number you can [connect to your newly created instance](https://cloud.google.com/sql/docs/mysql/external-connection-methods). The by-far easiest is Cloud Shell. Navigate to [Cloud Shell](https://console.cloud.google.com/) and Connect to your instance:

```shell
gcloud sql connect logoider --user=root --quiet
```

You will see a message about whitelisting your IP for incoming connection. Just wait until that's finished and when finally prompted for `root` password enter the value `$DB_ROOT_SECRET`.

> You can print it out in the console where you run original setup `echo DB_ROOT_SECRET`

That will get you to MySQL prompt in Cloud Shell:

```shell
MySQL [(none)]>
```

To create and configure your database schema, review and copy the entire [schema.ddl](sql/schema.ddl) file, paste it into the Cloud Shell window, and hit enter. You should see SQL output with 4 `Query OK` entries. You can now close the Cloud Shell window.

#### IAM

To ensure that your Cloud RUn service is able to do only the intended tasks and nothing more, you will create a service account which will be used to run the Cloud Run service and configure it with a few explicit roles:

* `run.invoker` - required to execute Cloud Run service
* `cloudsql.editor` - required to connect and write/read/delete on Cloud SQL
* `logging.logWriter` - required for Stackdriver logging
* `cloudtrace.agent` - required for Stackdriver tracing
* `monitoring.metricWriter` - required to write custom metrics to Stackdriver

To create and configure `logoider-sa` service account, run [bin/user](./bin/user) script:

```shell
bin/user
```

#### OAuth (chicken and an egg)

Cloud Run service URLs are not currently predictable. To deploy Cloud Run service you will need the OAuth credentials. To configure those credentials, you will need a callback URL that includes the Cloud Run service URL. You can see the chicken and an egg situation here. To work around it, we will deploy service first and then come back here and update the OAuth settings.

### Build Container Image

Cloud Run runs container images. To build one for this service we are going to use the included [Dockerfile](./Dockerfile) and submit it along with the source code as a build job to Cloud Build using [bin/image](./bin/image) script.

```shell
bin/image
```

### Deploy the Cloud Run Service

Once you have configured all the service dependencies, we can now deploy your Cloud Run service. To do that run [bin/service](./bin/service) script:

```shell
bin/service
```

The output of the script will include the URL by which you can access that service.

### OAuth

OAuth credentials have to be set up manually. In your Google Cloud Platform (GCP) project console navigate to the Credentials section. You can use the search bar, just type `Credentials` and select the option with "API & Services". To create new OAuth credentials:

* Click “Create credentials” and select “OAuth client ID”
* Select "Web application"
* Add authorized redirect URL at the bottom using the fully qualified domain we defined above and appending the `callback` path:
 * `https://[SERVICE_URL]/auth/callback`
* Click create and copy both `client id` and `client secret`
* CLICK `OK` to save

For ease of use, export the copied client `id` as `LI_OAUTH_CLIENT_ID` and `secret` as `LI_OAUTH_CLIENT_SECRET`

```shell
export LI_OAUTH_CLIENT_ID=""
export LI_OAUTH_CLIENT_SECRET=""
```

> You will also have to add the service URL to Authorized Domains section on the OAuth consent screen. More on that [here](https://support.google.com/cloud/answer/6158849?hl=en#authorized-domains)

### Deploy the Cloud Run Service (agin)

To finish configuring our service we will need to create a new Cloud Run service revision. The easiest way to do that is to redeploy the service again using the same [bin/service](./bin/service) script:

```shell
bin/service
```

## Usage

You can access the The Icon Identifier application by navigating to the deployed Cloud Run service URL. You can run the [bin/url](./bin/url) script to print it out again

```shell
bin/url
```

Once you access the deployed application:

* Click on "Sign in With Google" and follow the OAuth prompts
* Once authenticated, click on any of the provided logos (populates image URL)
* Click "Identify" button to see it's description

> You can also use any other publicly accessible logo image by pasting its URL

## Cleanup

To cleanup all resources created by this sample execute the [bin/cleanup](bin/cleanup) script.

```shell
bin/cleanup
```

## Disclaimer

This is my personal project and it does not represent my employer. I take no responsibility for issues caused by this code. I do my best to ensure that everything works, but if something goes wrong, my apologies is all you will get.
