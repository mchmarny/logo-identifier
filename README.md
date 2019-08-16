# logo-identifier

Simple Cloud Run service demonstrating use of Google Vision API and Cloud SQL as backing services for log identification with rate limiting and service usage metrics in Stackdriver.

## Why

Setting up development and deployment pipeline for Cloud Run service backed by project-level authenticated services like Cloud Vision and connection string based ones like Cloud Run can be complicated. While logo identification is probably not a real production workload, this sample aims to illustrate all aspects of local development and service deployment on GCP using these services.                       |

## Pre-requirements

### GCP Project and gcloud SDK

If you don't have one already, start by creating new project and configuring [Google Cloud SDK](https://cloud.google.com/sdk/docs/). Similarly, if you have not done so already, you will have [set up Cloud Run](https://cloud.google.com/run/docs/setup).

## Setup

To setup this service you will:

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

Before starting, we will need:

* Enable required GCP APIs
* Create Cloud SQL instance and configure it with database and user
* Create a Service Account and configure it with required IAM policies
* Configure OAuth credentials with Cloud Run service callback URL

Before we create the Cloud SQL instance and configure database, we need to define a couple of password: default DB user (`root`) and the application specific user (`logoider-db-user`). To do that we will use `openssl`. If for some reason you do not have `openssl` configured you can just set these value to your own secrets. Just don't re-use other secrets or make it too easy to guess.

```shell
export DB_ROOT_SECRET=$(openssl rand -base64 16)
echo "root user password: ${DB_ROOT_SECRET}"
export DB_USER_SECRET=$(openssl rand -base64 16)
echo "app user password: ${DB_USER_SECRET}"
```

Now we are ready to create dependencies by running [bin/setup](./bin/setup) script:

> You should review each one of the provided scripts for content to understand the individual commands

```shell
bin/setup
```

#### Cloud SQL (database schema)

One Cloud SQL finished configuring, there number of ways you can use to [connect to your newly created instance](https://cloud.google.com/sql/docs/mysql/external-connection-methods). The by-far easiest is Cloud Shell.

Navigate to [Cloud Shell](https://console.cloud.google.com/) and Connect to your instance:

```shell
gcloud sql connect logoider --user=root --quiet
```

You will see a message about whitelisting your IP for incoming connection. Just wait until that's finished and when finally prompted for `root` password enter the value `$DB_ROOT_SECRET`.

> You can print it out in the console where you run original setup `echo DB_ROOT_SECRET`

That will get you to MySQL prompt

```shell
MySQL [(none)]>
```

To create and configure your database schema, copy the entire [schema.ddl](sql/schema.ddl) file, paste it into the Cloud Shell window, and hit enter. You should see SQL output with 4 `Query OK` entries. You can now close the Cloud Shell window.

#### OAuth (chicken and an egg)

Cloud Run service URLs are not currently predictable. To deploy Cloud Run service will need the OAuth credentials and to configure those we will need a call back URL that includes the Cloud Run service URL. You can see the chicken and an egg situation here. To work around it, we will deploy service first and then come back here and update the OAuth settings.

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

>
To protect you and your users, Google only allows applications that authenticate using OAuth to use Authorized Domains. Your applications' links must be hosted on Authorized Domains. Learn more

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

To use the application:

* Click on "Sign in With Google" and follow the OAuth prompts
* Once authenticated
  * Click on any of the provided logos (populates image URL)
  * Click "Identify" button to see it's description

> You can also use any other publicly accessible logo image by pasting its URL

## Cleanup

To cleanup all resources created by this sample execute the [bin/cleanup](bin/cleanup) script.

```shell
bin/cleanup
```

## Disclaimer

This is my personal project and it does not represent my employer. I take no responsibility for issues caused by this code. I do my best to ensure that everything works, but if something goes wrong, my apologies is all you will get.
