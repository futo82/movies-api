# Movies RESTful API

This example project demonstrates how to write a RESTful API in Golang using the Gin framework. The API provides a create, retrieve, update, and delete operation against the movie data stored in DynamoDB. Okta's implicit flow is integrated to provide access control.

## Directory Structure

```
├── controllers
│   └── controller.go   // Handle Requests
├── db
│   └── dyanmodb.go     // Intialize Database 
├── middlewares
│   └── middleware.go   // Okta JWT Middleware
├── models
│   └── movie
│       └── movie.go    // Application Model
├── main.go
├── docker-compose.yml  // Startup API & DyanmoDB Containers
├── Dockerfile          // Build API Docker Image
└── run.sh
```

## Prerequisites

* [AWS CLI](https://aws.amazon.com/cli/)
* [Golang](https://golang.org/dl/)
* [Docker](https://www.docker.com/)

## Go Module Dependencies

* [Go Gin Web Framework](https://github.com/gin-gonic/gin)
* [AWS SDK Go](https://github.com/aws/aws-sdk-go)
* [Okta JWT Verifier](https://github.com/okta/okta-jwt-verifier-golang)

## API

#### GET /v1/api/movies/{movie-id}

```
curl -X GET \
  http://localhost:8080/v1/api/movies/19995 \
  -H 'Authorization: Bearer JWT_OKTA_TOKEN'
```

#### POST /v1/api/movies

```
curl -X POST \
  http://localhost:8080/v1/api/movies \
  -H 'Authorization: Bearer JWT_OKTA_TOKEN' \
  -H 'Content-Type: application/json' \
  -d '{
	"movie-id": "19995",
	"title": "Avatar",
	"budget": 237000000,
	"release-date": "2009-12-10T00:00:00Z",
	"revenue": 2787965087,
	"runtime": 162,
	"vote-average": 7.2,
	"vote-count": 11800
}'
```

#### PUT /v1/api/movies/{movie-id}

```
curl -X PUT \
  http://localhost:8080/v1/api/movies/19995 \
  -H 'Content-Type: application/json' \
  -d '{
	"title": "Avatar",
	"budget": 237000000,
	"release-date": "2009-12-10T00:00:00Z",
	"revenue": 2787965087,
	"runtime": 162,
	"vote-average": 7.2,
	"vote-count": 11800
}'
```

#### DELETE /v1/api/movies/{movie-id}

```
curl -X DELETE \
  http://localhost:8080/v1/api/movies/19995
```

## Setup Environment

Run the following command to set up your AWS CLI installation if it is not already done.

```
aws configure
```

Create a [Okta Developer Account](https://developer.okta.com/signup/) and setup a OpenID Connect Application.

1. From the Applications page, choose Add Application.
2. On the Create New Application page, select SPA.
3. Fill-in the Application Settings, then click Done.
4. Record the client id and the issuer values.

Create a .env file and add the following properties.

```
CLIENT_ID=<YOUR_OKTA_APPLICATION_CLIENT_ID>
ISSUER=<YOUR_OKTA_ISSUER>
```

## Obtaining an Okta Access Token

You will need to obtain an Okta access token and include it in the Authorization header when sending in a request to the API.

1. Obtain a Okta session token using the https://{baseUrl}/api/v1/authn endpoint.

```
curl -X POST \
  https://{baseUrl}/api/v1/authn \
  -H 'Content-Type: application/json' \
  -H 'cache-control: no-cache' \
  -d '{
  "username": "YOUR_USERNAME",
  "password": "YOUR_PASSWORD",
  "relayState": "YOUR_RELAY_STATE",
  "options": {
    "multiOptionalFactorEnroll": false,
    "warnBeforePasswordExpired": false
  }
}'
```

2. Obtain a Okta access token using the 
https://{baseUrl}/oauth2/default/v1/authorize endpoint.

```
curl -X GET \
  'https://{baseUrl}/oauth2/default/v1/authorize?client_id=YOUR_CLIENT_ID&response_type=token&scope=openid&redirect_uri={YOUR_REDIRECT_URI}&state=YOUR_STATE_VALUE&nonce=YOUR_NONCE_VALUE&sessionToken=YOUR_SESSION_TOKEN' \
  -H 'cache-control: no-cache'
```

## Build & Tag the API Docker Image

```
docker build . -t futo82/movies-api
```

## Run API locally that connects to AWS DynamoDB

This requires that you have an AWS account and your AWS CLI is configured with an access key and secret key that has permission to access AWS DyanmoDB.

Run the run.sh script and it will perform the following tasks.

```
./run.sh
```

* Use the AWS CLI to extract out the access key, secret key, and region from the default profile and set them as environment variables. 
* Extract out the Okta client id and issuer values from the .env file and set them as environment variables.
* Build and tag the API docker image.
* Startup the API docker container.

The API will listen on port 8080 and ready to accept requests. Make sure to include the Authorization header in the request with the Bearer token.

## Run API and DynamoDB locally

Before running the docker-compose command, the API docker image must be built and tagged (see above). After you run the docker-compose command, it will bring up the API and DynamoDB docker containers on your local machine. The API is configured to talk to the DyanmoDB container.

The API will listen on port 8080 and ready to accept requests. Make sure to include the Authorization header in the request with the Bearer token.

```
docker-compose up
```