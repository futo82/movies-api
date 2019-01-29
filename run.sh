#!/bin/bash

AWS_ACCESS_KEY_ID=$(aws --profile default configure get aws_access_key_id)
AWS_SECRET_ACCESS_KEY=$(aws --profile default configure get aws_secret_access_key)
AWS_REGION=$(aws --profile default configure get region)
CLIENT_ID=$(grep CLIENT .env | cut -d= -f2)
ISSUER=$(grep ISSUER .env | cut -d= -f2)
PORT=8080

docker build . -t futo82/movies-api
docker run -i -t -p 8080:8080 -e AWS_ACCESS_KEY_ID=$AWS_ACCESS_KEY_ID -e AWS_SECRET_ACCESS_KEY=$AWS_SECRET_ACCESS_KEY -e AWS_REGION=$AWS_REGION -e CLIENT_ID=$CLIENT_ID -e ISSUER=$ISSUER -e PORT=$PORT futo82/movies-api
