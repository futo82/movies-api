#!/bin/bash

export AWS_ACCESS_KEY_ID=$(aws --profile default configure get aws_access_key_id)
export AWS_SECRET_ACCESS_KEY=$(aws --profile default configure get aws_secret_access_key)
export AWS_REGION=$(aws --profile default configure get region)
export DYNAMODB_ENDPOINT=http://localhost:8000

go run load_data.go